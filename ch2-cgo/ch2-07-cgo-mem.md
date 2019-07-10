# 2.7 Mô hình bộ nhớ CGO

CGO là cầu nối giữa Go và C. Nó cho phép khả năng tương tác ở cấp độ nhị phân, nhưng chúng ta nên chú ý đến các vấn đề có thể phát sinh do sự khác biệt về mô hình bộ nhớ giữa hai ngôn ngữ. Nếu việc chuyển con trỏ có liên quan đến lệnh gọi hàm khác ngôn ngữ  được xử lý bởi CGO, có thể có trường hợp trong đó ngôn ngữ Go và ngôn ngữ C chia sẻ một segment bộ nhớ nhất định. Chúng ta biết rằng bộ nhớ của ngôn ngữ C là cố định sau khi cấp phát, nhưng ngôn ngữ Go có thể giữ địa chỉ bộ nhớ trong stack    cho mục đích dynamic scaling của hàm (đây là sự khác biệt lớn nhất giữa mô hình bộ nhớ Go và C). Nếu ngôn ngữ C giữ con trỏ Go trước khi di chuyển, việc truy cập đối tượng Go bằng con trỏ cũ sẽ khiến chương trình bị sập.

## 2.7.1 Truy cập bộ nhớ C

Bộ nhớ của không gian ngôn ngữ C khá ổn định, miễn là nó không bị release trước, thì không gian ngôn ngữ Go có thể được sử dụng rất tự tin. Truy cập bộ nhớ C trong Go là trường hợp đơn giản nhất và chúng ta đã thấy nó nhiều lần trong các ví dụ trước.

Do những hạn chế của hiện thực, chúng ta không thể tạo các slice lớn hơn 2GB bộ nhớ trong Go. Nhưng với sự trợ giúp của cgo, chúng ta có thể tạo ra hơn 2GB bộ nhớ trong môi trường ngôn ngữ C, sau đó chuyển sang slice của ngôn ngữ Go:

```go
package main

/*
#include <stdlib.h>

void* makeslice(size_t memsize) {
    return malloc(memsize);
}
*/
import "C"
import "unsafe"

func makeByteSlize(n int) []byte {
    p := C.makeslice(C.size_t(n))
    return ((*[1 << 31]byte)(p))[0:n:n]
}

func freeByteSlice(p []byte) {
    C.free(unsafe.Pointer(&p[0]))
}

func main() {
    s := makeByteSlize(1<<32+1)
    s[len(s)-1] = 255
    print(s[len(s)-1])
    freeByteSlice(s)
}
```

[>> mã nguồn](../examples/ch2/ch2.7/1-c-mem-access/main.go)

Trong ví dụ này, chúng ta sử dụng `makeByteSlize` để tạo ra một slice lớn hơn kích thước bộ nhớ 4G, do đó bỏ qua giới hạn trong hiện thực. Hàm hỗ trợ `freeByteSlice` được sử dụng để giải phóng slice được tạo từ hàm ngôn ngữ C.

Do không gian bộ nhớ ngôn ngữ C ổn định, các slice dựa trên cấu trúc bộ nhớ ngôn ngữ C cũng hoàn toàn ổn định và sẽ không bị di chuyển do những thay đổi trong stack của  Go.

## 2.7.2 C Truy cập tạm thời vào bộ nhớ Go

Một yếu tố chính trong sự tồn tại của cgo là tạo điều kiện cho việc sử dụng một lượng lớn tài nguyên phần mềm được xây dựng bằng ngôn ngữ Go bằng phần mềm ngôn ngữ C/C++ đã tồn tại trong vài thập kỷ qua. Nhiều thư viện trong C/C++ cần xử lý trực tiếp dữ liệu trong bộ nhớ Go thông qua các con trỏ. Do đó có nhiều trường hợp ứng dụng CGO cần chuyển bộ nhớ Go vào các hàm ngôn ngữ C.

Giả sử một trường hợp cực đoan: sau khi chúng ta truyền một hàm ngôn ngữ Go trên một stack của goroutinue, chúng ta sẽ truyền hàm ngôn ngữ C. Trong quá trình thực thi hàm ngôn ngữ C này, stack của goroutinue này được mở rộng do không đủ không gian, dẫn đến Bộ nhớ ngôn ngữ Go ban đầu sẽ được chuyển đến một vị trí mới. Nhưng tại thời điểm này, hàm ngôn ngữ C không biết rằng bộ nhớ ngôn ngữ Go đã di chuyển vị trí, vẫn sử dụng địa chỉ trước đó để vận hành bộ nhớ - điều này sẽ dẫn đến bộ nhớ vượt ngoài giới hạn. Trên đây là một hệ quả (có một số khác biệt trong tình huống thực tế), nghĩa là việc truy cập từ C vào bộ nhớ Go có thể không an toàn!

Tất nhiên, người dùng có kinh nghiệm với các cuộc gọi thủ tục từ xa (RPC) có thể xử lý bằng cách truyền hoàn toàn bằng giá trị: với đặc tính ổn định bộ nhớ của ngôn ngữ C, trước tiên khởi tạo cùng một lượng bộ nhớ trong không gian ngôn ngữ C, sau đó điền bộ dữ liệu từ Go vào bộ nhớ đó của C. Lúc trả về cũng được xử lý như vậy. Ví dụ sau đây là một triển khai cụ thể của ý tưởng này:

```go
package main

/*
void printString(const char* s) {
    printf("%s", s);
}
*/
import "C"

func printString(s string) {
    cs := C.CString(s)
    defer C.free(unsafe.Pointer(cs))

    C.printString(cs)
}

func main() {
    s := "hello"
    printString(s)
}
```

[>> mã nguồn](../examples/ch2/ch2.7/2-go-mem-access/example-1/main.go)

Khi bạn cần truyền chuỗi của Go sang ngôn ngữ C, trước tiên hãy sao chép dữ liệu bộ nhớ tương ứng với chuỗi `C.CString` ngôn ngữ Go sang không gian bộ nhớ ngôn ngữ C mới được tạo. Mặc dù ví dụ trên là an toàn, nhưng nó cực kỳ kém hiệu quả (vì phải phân bổ bộ nhớ nhiều lần và sao chép từng phần tử một), và nó cực kỳ cồng kềnh.

Để đơn giản hóa và xử lý hiệu quả vấn đề chuyển bộ nhớ ngôn ngữ Go sang ngôn ngữ C, CGO xác định một quy tắc đặc biệt cho trường hợp này: trước khi hàm ngôn ngữ C được CGO gọi trả về, CGO đảm bảo rằng bộ nhớ ngôn ngữ Go không tồn tại trong giai đoạn này. Khi thay đổi địa chỉ diễn ra, hàm ngôn ngữ C giờ có thể mạnh dạn sử dụng bộ nhớ ngôn ngữ Go!

Theo các quy tắc mới, chúng ta có thể truyền trực tiếp vào bộ nhớ của chuỗi Go:

```go
package main

/*
#include<stdio.h>

void printString(const char* s, int n) {
    int i;
    for(i = 0; i < n; i++) {
        putchar(s[i]);
    }
    putchar('\n');
}
*/
import "C"

func printString(s string) {
    p := (*reflect.StringHeader)(unsafe.Pointer(&s))
    C.printString((*C.char)(unsafe.Pointer(p.Data)), C.int(len(s)))
}

func main() {
    s := "hello"
    printString(s)
}
```

[>> mã nguồn](../examples/ch2/ch2.7/2-go-mem-access/example-2/main.go)

Việc xử lý giờ đã đơn giản hơn và tránh phân bổ thêm bộ nhớ. Một lời giải hoàn hảo?

Chúng ta giả định rằng hàm ngôn ngữ C được gọi cần phải thực thi trong một thời gian dài, điều này sẽ khiến ngôn ngữ Go được hàm tham chiếu tới không thể thay đổi (địa chỉ con trỏ bộ nhớ) trước khi ngôn ngữ C trả về, điều này có thể gián tiếp khiến goroutine tương ứng với stack bộ nhớ của Go không thể thực hiện dynamic scale. Điều này cũng tương đương với làm cho groutine này bị block. Do đó, trong các hàm ngôn ngữ C phải thực thi trong một thời gian dài (đặc biệt là trong các hoạt động CPU thuần túy, nhưng vì phải chờ các tài nguyên khác trong khoảng thời gian không xác định hoàn thành), bạn cần cẩn thận để xử lý bộ nhớ ngôn ngữ Go.

Phải thận trọng khi truyền vào hàm ngôn ngữ C ngay sau khi nhận được bộ nhớ Go. Bạn không thể lưu biến tạm và sau đó truyền  gián tiếp cho hàm ngôn ngữ C được vì CGO chỉ có thể đảm bảo rằng bộ nhớ ngôn ngữ Go được truyền vào sau khi gọi hàm C không thay đổi (địa chỉ trên bộ nhớ), không đảm bảo rằng bộ nhớ sẽ không thay đổi trước khi hàm C được truyền.

Đoạn code sau đây là sai:

```go
tmp := uintptr(unsafe.Pointer(&x))
pb := (*int16)(unsafe.Pointer(tmp))
*pb = 42
```

Vì `tmp` không thuộc kiểu con trỏ, nên đối tượng `x` có thể thay đổi địa chỉ sau khi nó nhận được địa chỉ của đối tượng Go, nhưng vì nó không thuộc kiểu con trỏ, nên nó sẽ không được cập nhật thành địa chỉ của bộ nhớ mới khi ngôn ngữ Go thực thi. Giữ địa chỉ của đối tượng Go trong kiểu `tmp` không có con trỏ có tác dụng tương tự như giữ địa chỉ của đối tượng Go trong ngôn ngữ C: nếu địa chỉ bộ nhớ của đối tượng Go ban đầu đã thay đổi, runtime của Go sẽ không cập nhật đồng bộ.

## 2.7.3 Giữ đối tượng con trỏ dài hạn trong C

Là một lập trình viên Go nên khi sử dụng CGO chúng ta sẽ luôn nghĩ rằng Go gọi các hàm C. Trên thực tế, trong CGO, các hàm ngôn ngữ C cũng có thể gọi lại các hàm được hiện thực bởi Go. Cụ thể, chúng ta có thể viết một thư viện động bằng ngôn ngữ Go và export interface đặc tả ngôn ngữ C cho người dùng khác. Khi hàm ngôn ngữ C gọi hàm ngôn ngữ Go, hàm ngôn ngữ C sẽ trở thành người gọi chương trình và vòng đời của đối tượng bộ nhớ Go được trả về bởi hàm ngôn ngữ Go hoàn toàn nằm ngoài sự quản lý của runtime trong Go. **Nói tóm lại, chúng ta không thể sử dụng bộ nhớ của đối tượng ngôn ngữ Go trực tiếp trong hàm ngôn ngữ C.**

Nếu cần truy cập các đối tượng bộ nhớ ngôn ngữ Go bằng ngôn ngữ C, chúng ta có thể ánh xạ các đối tượng bộ nhớ ngôn ngữ Go trong không gian ngôn ngữ Go sang kiểu int `id`, sau đó gián tiếp truy cập và điều khiển các đối tượng ngôn ngữ Go thông qua `id` này.

Đoạn code sau được sử dụng để ánh xạ đối tượng Go sang `ObjectId` của kiểu int. Sau khi sử dụng, bạn cần gọi thủ công phương thức free để giải phóng nó:

```go
package main

import "sync"

type ObjectId int32

var refs struct {
    sync.Mutex
    objs map[ObjectId]interface{}
    next ObjectId
}

func init() {
    refs.Lock()
    defer refs.Unlock()

    refs.objs = make(map[ObjectId]interface{})
    refs.next = 1000
}

func NewObjectId(obj interface{}) ObjectId {
    refs.Lock()
    defer refs.Unlock()

    id := refs.next
    refs.next++

    refs.objs[id] = obj
    return id
}

func (id ObjectId) IsNil() bool {
    return id == 0
}

func (id ObjectId) Get() interface{} {
    refs.Lock()
    defer refs.Unlock()

    return refs.objs[id]
}

func (id *ObjectId) Free() interface{} {
    refs.Lock()
    defer refs.Unlock()

    obj := refs.objs[*id]
    delete(refs.objs, *id)
    *id = 0

    return obj
}
```

Chúng ta sử dụng `map` để quản lý ánh xạ giữa các đối tượng ngôn ngữ Go và đối tượng `id`. `NewObjectId` được sử dụng để tạo `id` liên kết với đối tượng và phương thức của đối tượng `id` có thể được sử dụng để decode đối tượng Go ban đầu và cũng có thể được sử dụng để kết thúc liên kết của `id` và đối tượng Go ban đầu.

Tập hợp các hàm sau đây được export dưới dạng các interface C và có thể được gọi bằng các hàm ngôn ngữ C:

```go
package main

/*
extern char* NewGoString(char* );
extern void FreeGoString(char* );
extern void PrintGoString(char* );

static void printString(const char* s) {
    char* gs = NewGoString(s);
    PrintGoString(gs);
    FreeGoString(gs);
}
*/
import "C"

//export NewGoString
func NewGoString(s *C.char) *C.char {
    gs := C.GoString(s)
    id := NewObjectId(gs)
    return (*C.char)(unsafe.Pointer(uintptr(id)))
}

//export FreeGoString
func FreeGoString(p *C.char) {
    id := ObjectId(uintptr(unsafe.Pointer(p)))
    id.Free()
}

//export PrintGoString
func PrintGoString(s *C.char) {
    id := ObjectId(uintptr(unsafe.Pointer(p)))
    gs := id.Get().(string)
    print(gs)
}

func main() {
    C.printString("hello")
}
```

Trong hàm `printString`, chúng ta tạo một đối tượng chuỗi Go tương ứng thông qua `NewGoString`, giá trị trả về thực sự là một id và không thể được sử dụng trực tiếp. Chúng ta sử dụng hàm `PrintGoString` để  parse id này thành chuỗi ngôn ngữ Go. Chuỗi hoàn toàn mở rộng việc quản lý bộ nhớ của ngôn ngữ Go trong hàm ngôn ngữ C. Ngay cả khi địa chỉ chuỗi Go thay đổi quá trình stack scaling gây ra trước lệnh gọi `PrintGoString`, nó vẫn có thể hoạt động bình thường, vì id tương ứng với chuỗi ổn định. Chuỗi thu được bằng cách decode id trong không gian ngôn ngữ Go cũng hợp lệ.

## 2.7.4 Export các hàm C

Golang phân bổ bộ nhớ từ một không gian địa chỉ ảo cố định. Bộ nhớ được phân bổ bởi ngôn ngữ C không thể sử dụng không gian bộ nhớ ảo được dành riêng cho ngôn ngữ Go. Trong môi trường CGO, runtime của Go luôn kiểm tra theo mặc định liệu bộ nhớ được trả về  do lệnh export  có được phân bổ bởi ngôn ngữ Go hay không và nếu có sẽ ném ra runtime exception.

Sau đây là một ví dụ về  runtime exception trong CGO:

```go
/*
extern int* getGoPtr();

static void Main() {
    int* p = getGoPtr();
    *p = 42;
}
*/
import "C"

func main() {
    C.Main()
}

//export getGoPtr
func getGoPtr() *C.int {
    return new(C.int)
}
```

`GetGoPtr` trả về một con trỏ của kiểu trong C, nhưng bộ nhớ được phân bổ từ hàm Go của ngôn ngữ Go, là bộ nhớ được quản lý bởi runtime của Go. Sau đó, chúng ta gọi hàm `getGoPtr` trong hàm main của C và runtime exception sẽ ném ra theo mặc định:

```sh
$ go run main.go
panic: runtime error: cgo result has Go pointer

goroutine 1 [running]:
main._cgoexpwrap_cfb3840e3af2_getGoPtr.func1(0xc420051dc0)
command-line-arguments/_obj/_cgo_gotypes.go:60 +0x3a
main._cgoexpwrap_cfb3840e3af2_getGoPtr(0xc420016078)
command-line-arguments/_obj/_cgo_gotypes.go:62 +0x67
main._Cfunc_Main()
command-line-arguments/_obj/_cgo_gotypes.go:43 +0x41
main.main()
/Users/chai/go/src/github.com/chai2010 \
/advanced-go-programming-book/examples/ch2-xx \
/return-go-ptr/main.go:17 +0x20
exit status 2
```

Exception chỉ ra rằng kết quả được trả về bởi hàm cgo chứa con trỏ tới ngôn ngữ Go. Hoạt động kiểm tra của con trỏ xảy ra trong phiên bản ngôn ngữ C của hàm `getGoPtr`, đây là hàm của ngôn ngữ C và ngôn ngữ Go được CGO tạo ra.

Sau đây là phiên bản ngôn ngữ C của hàm `getGoPtr` được CGO tạo ra (`_cgo_export.c` được xác định trong file được CGO tạo ra):

```c
int* getGoPtr()
{
    __SIZE_TYPE__ _cgo_ctxt = _cgo_wait_runtime_init_done();
    struct {
        int* r0;
    } __attribute__((__packed__)) a;
    _cgo_tsan_release();
    crosscall2(_cgoexp_95d42b8e6230_getGoPtr, &a, 8, _cgo_ctxt);
    _cgo_tsan_acquire();
    _cgo_release_context(_cgo_ctxt);
    return a.r0;
}
```

Trong đó `_cgo_tsan_acquire` là hàm scan con trỏ bộ nhớ được chuyển đổi từ project LLVM, kiểm tra xem kết quả được trả về bởi hàm CGO có chứa con trỏ Go hay không.

Cần lưu ý rằng việc kiểm tra mặc định của con trỏ đối với kết quả được trả về là khá tốn kém, đặc biệt nếu kết quả được trả về bởi hàm CGO là một cấu trúc dữ liệu phức tạp, sẽ mất nhiều thời gian hơn nữa. Nếu bạn đã đảm bảo rằng các kết quả được trả về bởi hàm CGO là an toàn, bạn có thể  tắt thao tác kiểm tra con trỏ bằng cách đặt biến môi trường `GODEBUG=cgocheck=0`.

```sh
$ GODEBUG=cgocheck=0 go run main.go
```

Sau khi tắt `cgocheck` và chạy đoạn code trên, exception  sẽ không được ném ra. Tuy nhiên, cần lưu ý rằng nếu bộ nhớ tương ứng trong ngôn ngữ C được release bởi runtime của Go, nó sẽ gây ra sự cố nghiêm trọng hơn. Giá trị mặc định của `cgocheck` là 1, tương ứng với phiên bản detection đơn giản hoá. Nếu bạn cần hàm detection đầy đủ, bạn có thể đặt `cgocheck` thành 2.

Để biết mô tả chi tiết về các hàm CGO rumtime pointer detection hãy tham khảo tài liệu chính thức của Golang: [package runtime - GoDoc](https://godoc.org/runtime#hdr-Environment_Variables).
