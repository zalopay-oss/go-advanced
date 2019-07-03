# 2.7 Mô hình bộ nhớ CGO

CGO là cầu nối giữa Go và C. Nó cho phép khả năng tương tác ở cấp độ giao diện nhị phân, nhưng chúng ta nên chú ý đến các vấn đề có thể phát sinh do sự khác biệt về mô hình bộ nhớ giữa hai ngôn ngữ. Nếu việc chuyển con trỏ có liên quan đến lệnh gọi chức năng ngôn ngữ chéo được xử lý bởi CGO, có thể có các cảnh trong đó ngôn ngữ Go và ngôn ngữ C chia sẻ một phân đoạn bộ nhớ nhất định. Chúng tôi biết rằng bộ nhớ của ngôn ngữ C ổn định sau khi cấp phát, nhưng ngôn ngữ Go có thể khiến địa chỉ bộ nhớ trong ngăn xếp di chuyển do tỷ lệ động của ngăn xếp chức năng (đây là sự khác biệt lớn nhất giữa các mô hình bộ nhớ Go và C). Nếu ngôn ngữ C giữ con trỏ Go trước khi di chuyển, việc truy cập đối tượng Go bằng con trỏ cũ sẽ khiến chương trình bị sập.

## 2.7.1 Truy cập bộ nhớ C

Bộ nhớ của không gian ngôn ngữ C ổn định, miễn là nó không được con người phát hành trước, thì không gian ngôn ngữ Go có thể được sử dụng một cách tự tin. Truy cập bộ nhớ C trong Go là trường hợp đơn giản nhất và chúng tôi đã thấy nó nhiều lần trong các ví dụ trước.

Do những hạn chế của việc triển khai Go, chúng tôi không thể tạo các lát lớn hơn 2GB bộ nhớ trong Go (xem code triển khai makeice để biết chi tiết). Nhưng với sự trợ giúp của công nghệ cgo, chúng ta có thể tạo ra hơn 2GB bộ nhớ trong môi trường ngôn ngữ C, sau đó chuyển sang lát ngôn ngữ Go:

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

Trong ví dụ này, chúng tôi sử dụng makeByteSlize để tạo ra một lát lớn hơn kích thước bộ nhớ 4G, do đó bỏ qua giới hạn triển khai ngôn ngữ Go (yêu cầu xác thực code). Hàm trợ giúp freeByteSlice được sử dụng để giải phóng lát cắt được tạo từ hàm ngôn ngữ C.

Do không gian bộ nhớ ngôn ngữ C ổn định, các lát dựa trên cấu trúc bộ nhớ ngôn ngữ C cũng hoàn toàn ổn định và sẽ không bị di chuyển do những thay đổi trong ngăn xếp ngôn ngữ Go.

## 2.7.2 C Truy cập tạm thời vào bộ nhớ đến

Một yếu tố chính trong sự tồn tại của cgo là tạo điều kiện thuận lợi cho việc sử dụng một lượng lớn tài nguyên phần mềm được xây dựng bằng ngôn ngữ Go bằng phần mềm ngôn ngữ C / C ++ trong vài thập kỷ qua. Nhiều thư viện trong C / C ++ cần xử lý trực tiếp dữ liệu bộ nhớ đến thông qua các con trỏ. Do đó, có nhiều kịch bản ứng dụng trong Cgo cần chuyển bộ nhớ Go vào các chức năng ngôn ngữ C.

Giả sử một kịch bản cực đoan: sau khi chúng ta chuyển một hàm ngôn ngữ Go trên một ngăn xếp của goroutinue, chúng ta sẽ truyền hàm ngôn ngữ C. Trong quá trình thực thi chức năng ngôn ngữ C này, ngăn xếp của goroutinue này được mở rộng do không đủ không gian, dẫn đến Bộ nhớ ngôn ngữ Go ban đầu đã được chuyển đến một vị trí mới. Nhưng tại thời điểm này, chức năng ngôn ngữ C không biết rằng bộ nhớ ngôn ngữ Go đã di chuyển vị trí, vẫn sử dụng địa chỉ trước đó để vận hành bộ nhớ - điều này sẽ dẫn đến bộ nhớ ngoài giới hạn. Trên đây là một hệ quả (có một số khác biệt trong tình huống thực tế), có nghĩa là việc truy cập C vào bộ nhớ Go đến có thể không an toàn!

Tất nhiên, người dùng có kinh nghiệm với các cuộc gọi thủ tục từ xa RPC có thể xem xét xử lý bằng cách chuyển hoàn toàn các giá trị: với tính năng ổn định bộ nhớ ngôn ngữ C, trước tiên hãy mở cùng một lượng bộ nhớ trong không gian ngôn ngữ C, sau đó điền bộ nhớ Go vào bộ nhớ C. Không gian, bộ nhớ trả về được xử lý như vậy. Ví dụ sau đây là một triển khai cụ thể của ý tưởng này:

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

Khi bạn cần chuyển chuỗi của Go sang ngôn ngữ C, trước tiên hãy C.CStringsao chép dữ liệu bộ nhớ tương ứng với chuỗi ngôn ngữ Go sang không gian bộ nhớ ngôn ngữ C mới được tạo. Mặc dù ví dụ trên là an toàn, nhưng nó cực kỳ kém hiệu quả (vì nó phải phân bổ bộ nhớ nhiều lần và sao chép từng phần tử một), và nó cực kỳ cồng kềnh.

Để đơn giản hóa và xử lý hiệu quả vấn đề chuyển bộ nhớ ngôn ngữ Go này sang ngôn ngữ C, cgo xác định một quy tắc đặc biệt cho kịch bản này: trước khi hàm ngôn ngữ C được CGO gọi trở lại, cgo đảm bảo rằng bộ nhớ ngôn ngữ Go không tồn tại trong giai đoạn này. Di chuyển xảy ra, chức năng ngôn ngữ C có thể mạnh dạn sử dụng bộ nhớ ngôn ngữ Go!

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

Việc xử lý hiện tại đơn giản hơn và tránh phân bổ thêm bộ nhớ. Giải pháp hoàn hảo!

Khi bất kỳ công nghệ hoàn hảo nào bị lạm dụng, các quy tắc dường như hoàn hảo của CGO cũng bị ẩn đi. Chúng tôi giả định rằng hàm ngôn ngữ C được gọi cần phải chạy trong một thời gian dài, điều này sẽ khiến ngôn ngữ Go được anh ta tham chiếu không thể di chuyển trước khi ngôn ngữ C trở lại, điều này có thể gián tiếp khiến goroutine tương ứng với ngăn xếp bộ nhớ Go không tự động điều chỉnh bộ nhớ ngăn xếp. Đó là, nó có thể khiến con goroutine này bị chặn. Do đó, trong các chức năng ngôn ngữ C cần chạy trong một thời gian dài (đặc biệt là trong các hoạt động CPU thuần túy, nhưng cũng vì phải chờ các tài nguyên khác và cần thời gian không chắc chắn để hoàn thành chức năng), bạn cần cẩn thận để xử lý bộ nhớ ngôn ngữ Go đến.

Tuy nhiên, bạn cần cẩn thận để chuyển chức năng ngôn ngữ C ngay sau khi nhận được bộ nhớ Go. Bạn không thể lưu biến tạm thời và sau đó chuyển gián tiếp chức năng ngôn ngữ C. Vì CGO chỉ có thể đảm bảo rằng bộ nhớ ngôn ngữ Go được truyền vào sau khi gọi hàm C không di chuyển, nó không đảm bảo rằng bộ nhớ sẽ không thay đổi trước khi chức năng C được truyền.

Các code sau đây là sai:

```go
// 错误的代码
tmp := uintptr(unsafe.Pointer(&x))
pb := (*int16)(unsafe.Pointer(tmp))
*pb = 42
```

Vì tmp không phải là loại con trỏ, nên đối tượng x có thể được di chuyển sau khi nó nhận được địa chỉ đối tượng Go, nhưng vì nó không phải là loại con trỏ, nên nó sẽ không được cập nhật thành địa chỉ của bộ nhớ mới khi ngôn ngữ Go được chạy. Giữ địa chỉ của đối tượng Go trong loại tmp không có con trỏ có tác dụng tương tự như giữ địa chỉ của đối tượng Go trong ngôn ngữ C: nếu bộ nhớ của đối tượng Go ban đầu đã di chuyển, thời gian chạy ngôn ngữ Go sẽ không cập nhật chúng đồng bộ.

## 2.7.3 C giữ đối tượng con trỏ dài hạn

Là một lập trình viên Go, khi sử dụng CGO, tiềm thức sẽ luôn nghĩ rằng Go gọi các hàm C. Trên thực tế, trong CGO, các chức năng ngôn ngữ C cũng có thể gọi lại các chức năng được thực hiện bởi Go. Cụ thể, chúng ta có thể viết một thư viện động bằng ngôn ngữ Go và xuất giao diện của đặc tả ngôn ngữ C cho người dùng khác. Khi chức năng ngôn ngữ C gọi chức năng ngôn ngữ Go, chức năng ngôn ngữ C sẽ trở thành người gọi chương trình và vòng đời của bộ nhớ đối tượng Go được trả về bởi chức năng ngôn ngữ Go hoàn toàn nằm ngoài sự quản lý của thời gian chạy ngôn ngữ Go. Nói tóm lại, chúng ta không thể sử dụng bộ nhớ của đối tượng ngôn ngữ Go trực tiếp trong chức năng ngôn ngữ C.

Mặc dù ngôn ngữ Go cấm giữ lâu dài các đối tượng con trỏ Go trong các chức năng ngôn ngữ C, nhưng yêu cầu này là hữu hình. Nếu bạn cần truy cập các đối tượng bộ nhớ ngôn ngữ Go bằng ngôn ngữ C, chúng ta có thể ánh xạ các đối tượng bộ nhớ ngôn ngữ Go trong không gian ngôn ngữ Go sang id kiểu int, sau đó gián tiếp truy cập và điều khiển các đối tượng ngôn ngữ Go thông qua id này.

Đoạn code sau được sử dụng để ánh xạ đối tượng Go sang ObjectId của kiểu số nguyên. Sau khi sử dụng, bạn cần gọi thủ công phương thức miễn phí để giải phóng ID đối tượng:

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

Chúng tôi sử dụng bản đồ để quản lý ánh xạ giữa các đối tượng ngôn ngữ Go và đối tượng id. NewObjectId được sử dụng để tạo id liên kết với đối tượng và phương thức của đối tượng id có thể được sử dụng để giải code đối tượng Go ban đầu và cũng có thể được sử dụng để kết thúc liên kết của id và đối tượng Go ban đầu.

Tập hợp các hàm sau đây được xuất dưới dạng thông số kỹ thuật giao diện C và có thể được gọi bằng các hàm ngôn ngữ C:

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

Trong hàm printString, chúng tôi tạo một đối tượng chuỗi Go tương ứng thông qua NewGoString, lợi nhuận thực sự là một id, không thể được sử dụng trực tiếp. Chúng tôi sử dụng chức năng PrintGoString để phân tích id thành chuỗi ngôn ngữ Go. Chuỗi hoàn toàn mở rộng việc quản lý bộ nhớ của ngôn ngữ Go trong chức năng ngôn ngữ C. Ngay cả khi địa chỉ chuỗi Go thay đổi do tỷ lệ ngăn xếp gây ra trước lệnh gọi PrintGoString, nó vẫn có thể hoạt động bình thường, vì id tương ứng với chuỗi ổn định. Chuỗi thu được bằng cách giải code id trong không gian ngôn ngữ Go cũng hợp lệ.

## 2.7.4 Xuất các hàm C không thể trả về bộ nhớ

Trong Go, Go phân bổ bộ nhớ từ một không gian địa chỉ ảo cố định. Bộ nhớ được phân bổ bởi ngôn ngữ C không thể sử dụng không gian bộ nhớ ảo được dành riêng bởi ngôn ngữ Go. Trong môi trường CGO, thời gian chạy ngôn ngữ Go kiểm tra theo mặc định liệu bộ nhớ được trả về xuất khẩu có được phân bổ bởi ngôn ngữ Go hay không và nếu có, một ngoại lệ thời gian chạy sẽ bị ném.

Sau đây là một ví dụ về ngoại lệ thời gian chạy CGO:

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

GetGoPtr trả về một con trỏ của loại ngôn ngữ C, nhưng bộ nhớ được phân bổ từ chức năng Go của ngôn ngữ Go, là bộ nhớ được quản lý bởi thời gian chạy ngôn ngữ Go. Sau đó, chúng ta gọi hàm getGoPtr trong hàm chính C và ngoại lệ thời gian chạy sẽ được gửi theo mặc định:

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

Ngoại lệ chỉ ra rằng kết quả được trả về bởi hàm cgo chứa con trỏ tới gán ngôn ngữ Go. Hoạt động kiểm tra của con trỏ xảy ra trong phiên bản ngôn ngữ C của hàm getGoPtr, đây là chức năng của ngôn ngữ C và ngôn ngữ Go được tạo bởi cgo.

Sau đây là chi tiết cụ thể của phiên bản ngôn ngữ C của hàm getGoPtr được tạo bởi cgo ( _cgo_export.cđược xác định trong tệp được tạo bởi cgo ):

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

Một trong số đó _cgo_tsan_acquirelà hàm quét con trỏ bộ nhớ được chuyển từ dự án LLVM, kiểm tra xem kết quả được trả về bởi hàm cgo có chứa con trỏ Go hay không.

Cần lưu ý rằng việc kiểm tra mặc định của con trỏ đối với kết quả được trả về là tốn kém, đặc biệt nếu kết quả được trả về bởi hàm cgo là một cấu trúc dữ liệu phức tạp, sẽ mất nhiều thời gian hơn. Nếu bạn đã đảm bảo rằng các kết quả được trả về bởi hàm cgo là an toàn, bạn có thể GODEBUG=cgocheck=0 tắt hành vi kiểm tra con trỏ bằng cách đặt biến môi trường .

```sh
$ GODEBUG=cgocheck=0 go run main.go
```

Sau khi đóng chức năng cgocheck và sau đó chạy đoạn code trên, ngoại lệ trên sẽ không xảy ra. Tuy nhiên, cần lưu ý rằng nếu bộ nhớ tương ứng trong ngôn ngữ C được phát hành bởi thời gian chạy Go, nó sẽ gây ra sự cố nghiêm trọng hơn. Giá trị mặc định của cgocheck là 1, tương ứng với việc phát hiện phiên bản đơn giản hóa. Nếu bạn cần chức năng phát hiện đầy đủ, bạn có thể đặt cgocheck thành 2.

Để biết mô tả chi tiết về các chức năng của phát hiện con trỏ thời gian chạy, hãy tham khảo tài liệu chính thức của ngôn ngữ Go.
