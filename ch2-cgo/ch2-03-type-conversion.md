# 2.3 Chuyển đổi kiểu

Ban đầu, CGO được tạo ra để thuận lợi cho việc sử dụng các hàm trong C (các hàm hiện thực khai báo Golang  trong C) để sử dụng lại các tài nguyên của  C (vì ngôn ngữ C cũng liên quan đến các hàm callback, dĩ nhiên nó liên quan đến việc gọi các hàm trong Go từ các hàm của C (các hàm thực hiện khai báo ngôn ngữ C trong Go)). Ngày nay, CGO đã phát triển thành cầu nối giao tiếp hai chiều giữa C và Go. Để tận dụng tính năng của CGO, việc hiểu các quy tắc chuyển đổi giữa hai loại ngôn ngữ là điều quan trọng. Đây là vấn đề sẽ được thảo luận trong phần này.

## 2.3.1 Các kiểu dữ liệu số học

Khi ta sử dụng các ký hiệu của C trong Golang, thường nó sẽ truy cập thông qua package "C" ảo, chẳng hạn như kiểu `int` tương ứng với `C. int`. Một số kiểu trong C bao gồm nhiều từ khóa, nhưng khi truy cập chúng thông qua package "C" ảo, phần tên không thể có ký tự khoảng trắng, ví dụ `unsigned int` không thể  truy cập trực tiếp `C.unsigned int`. Do đó, CGO cung cấp quy tắc chuyển đổi tương ứng cho các kiểu trong C cơ bản, ví dụ như `C.uint` tương ứng trong  C là `unsigned int`.

Kiểu dữ liệu số học và kiểu dữ liệu trong C của Golang về cơ bản là tương tự nhau. Bảng 2-1 thể hiện sự tương tự này.

| Kiểu trong C           | Kiểu trong CGO | Kiểu trong Go |
| ---------------------- | -------------- | ------------- |
| char                   | C.char         | byte          |
| singed char            | C.schar        | int8          |
| unsigned char          | C.uchar        | uint8         |
| short                  | C.short        | int16         |
| unsigned short         | C.ushort       | uint16        |
| int                    | C.int          | int32         |
| unsigned int           | C.uint         | uint32        |
| long                   | C.long         | int32         |
| unsigned long          | C.ulong        | uint32        |
| long long int          | C.longlong     | int64         |
| unsigned long long int | C.ulonglong    | uint64        |
| float                  | C.float        | float32       |
| double                 | C.double       | float64       |
| size_t                 | C.size_t       | uint          |

_Bảng 2-1 So sánh kiểu trong các ngôn ngữ Go và C_

Cần lưu ý rằng mặc dù kích thước của những kiểu không được chỉ rõ kích thước trong C như `int`, `short` v.v., kích thước  của chúng đều được xác định trong CGO. Trong CGO, kiểu `int` và `uint` của C đều có kích thước 4 byte, kiểu `size_t` có thể được coi là kiểu số nguyên không dấu  `uint` của ngôn ngữ Go .

Mặc dù kiểu `int` và `uint` của C đều có kích thước cố định, nhưng với GO thì  `int` và `uint` có thể là 4 byte hoặc 8 byte. Nếu cần sử dụng đúng kiểu `int` của C trong Go, bạn có thể  sử dụng kiểu `GoInt` được xác định trong file header  `_cgo_export.h` được tạo ra bởi công cụ CGO.Trong file header này, mỗi kiểu giá trị Go cơ bản sẽ  xác định kiểu tương ứng trong C  có tiền tố "Go". Ví dụ sau trong hệ thống 64-bit, có file header `_cgo_export.h` được CGO định nghĩa các kiểu giá trị, nơi mà `GoInt` và `GoUint` lần lượt là `GoInt64` và `GoUint64`:

```go
typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef float GoFloat32;
typedef double GoFloat64;
```

Bên cạnh `GoInt` và `GoUint`, chúng tôi không khuyên bạn nên sử dụng trực tiếp  `GoInt32`, `GoInt64` và các kiểu khác. Cách tiếp cận tốt hơn là khai báo file header <stdint.h> thông qua tiêu chuẩn C99 của C. Để cải thiện tính linh hoạt của C, không chỉ mỗi kiểu số học được xác định kích thước rõ ràng trong file mà còn chúng còn sử dụng các tên phù hợp với tên kiểu tương ứng trong Golang. So sánh các kiểu tương ứng trong <stdint.h> được trình bày trong Bảng 2-2.

| Kiểu trong C | Kiểu trong CGO | Kiểu trong Go |
| ------------ | -------------- | ------------- |
| int8_t       | C.int8_t       | int8          |
| uint8_t      | C.uint8_t      | uint8         |
| int16_t      | C.int16_t      | int16         |
| uint16_t     | C.uint16_t     | uint16        |
| int32_t      | C.int32_t      | int32         |
| uint32_t     | C.uint32_t     | uint32        |
| int64_t      | C.int64_t      | int64         |
| uint64_t     | C.uint64_t     | uint64        |

_Bảng 2-2 So sánh kiểu trong `stdint.h` _

Như đã đề cập trước đó, nếu kiểu trong C bao gồm nhiều từ khóa, nó không thể được sử dụng trực tiếp thông qua package "C" ảo (ví dụ: `unsigned short` không thể được truy cập trực tiếp `C.unsigned short`). Tuy nhiên, sau khi định nghĩa lại kiểu trong <stdint.h> bằng cách sử dụng `typedef`, chúng ta có thể truy cập tới kiểu gốc. Đối với các kiểu trong C phức tạp hơn thì nên sử dụng `typedef` để đặt lại tên cho nó, thuận tiện hơn cho việc truy cập trong CGO.

## 2.3.2 Go Strings và Slices

Trong  file header `_cgo_export.h` được tạo ra bởi CGO, kiểu trong C tương ứng cũng được tạo cho các kiểu của Go như string, slice, dictionary, interface và pipe:

```go
typedef struct { const char *p; GoInt n; } GoString;
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;
```

Tuy nhiên, cần lưu ý rằng chỉ các string và slice là có giá trị sử dụng trong CGO, vì CGO tạo ra các phiên bản ngôn ngữ C cho một số  hàm trong Go, vì vậy cả hai đều có thể gọi các hàm C trong Go, điều này được thực hiện lặp tức và CGO không cung cấp các hàm hỗ trợ liên quan cho các kiểu khác, đồng thời  mô hình bộ nhớ dành riêng cho ngôn ngữ Go ngăn chúng ta duy trì các kiểu con trỏ tới bộ nhớ này quản lý bởi  Go, vì vậy mà môi trường ngôn ngữ C của các kiểu đó không có giá trị sử dụng.

Trong hàm C đã export, chúng ta có thể trực tiếp sử dụng các string và slice trong Go. Giả sử bạn có hai hàm export sau:

```go
//export helloString
func helloString(s string) {}

//export helloSlice
func helloSlice(s []byte) {}
```

File header `_cgo_export.h` được tạo bởi CGO sẽ chứa khai báo hàm sau:

```go
extern void helloString(GoString p0);
extern void helloSlice(GoSlice p0);
```

Nhưng lưu ý rằng nếu bạn sử dụng kiểu `GoString` thì sẽ phụ thuộc vào file header `_cgo_export.h` và tập file này có output động.

Phiên bản Go1.10 thêm một chuỗi kiểu  `_GoString_` định nghĩa trước, có thể làm giảm xuống code có rủi ro phụ thuộc file header `_cgo_export.h`. Chúng ta có thể điều chỉnh khai báo ngôn ngữ C của hàm `helloString` thành:

```go
extern void helloString(_GoString_ p0);
```

## 2.3.3 Struct, Union, Enumerate

Các kiểu struct, Union và Enumerate của ngôn ngữ C không thể được nhúng dưới dạng thành phần ẩn danh vào struct của ngôn ngữ Go. Trong Go, chúng ta có thể truy cập các kiểu struct như  `struct xxx` tương ứng là `C.struct_xxx` trong ngôn ngữ C. Bố cục bộ nhớ của struct tuân theo các quy tắc căn chỉnh (alignment) chung của ngôn ngữ C. Trong môi trường ngôn ngữ Go 32 bit, struct của C cũng tuân theo quy tắc căn chỉnh 32 bit và môi trường ngôn ngữ Go 64 bit tuân theo quy tắc căn chỉnh 64 bit. Đối với các struct có quy tắc căn chỉnh đặc biệt được chỉ định, chúng không thể được truy cập trong CGO.

Cách sử dụng struct đơn giản như sau:

```go
*/
struct A {
    int i;
    float f;
};
*/
import "C"
import "fmt"

func main() {
    var a C.struct_A
    fmt.Println(a.i)
    fmt.Println(a.f)
}
```

Nếu tên thành phần của struct tình cờ là một từ khóa trong ngôn ngữ Go, bạn có thể truy cập nó bằng cách thêm một dấu gạch dưới ở đầu tên thành viên:

```go
/*
struct A {
    int type; // type là một từ khóa trong Golang
};
*/
import "C"
import "fmt"

func main() {
    var a C.struct_A
    fmt.Println(a._type) // _type tương ứng với type
}
```

Nhưng nếu có 2 thành phần: một thành phần được đặt tên theo từ khóa của Go và phần kia là trùng khi thêm vào dấu gạch dưới, thì các thành phần được đặt tên theo từ khóa ngôn ngữ Go sẽ không thể truy cập (bị chặn):

```go
/*
struct A {
    int   type;  // type là một từ khóa trong Go
    float _type; // chặn CGO truy cập type trên kia
};
*/
import "C"
import "fmt"

func main() {
    var a C.struct_A
    fmt.Println(a._type) // _type tương ứng với _type
}
```

Các thành phần tương ứng với trường bit (biến được định nghĩa với giá trị độ lớn cho sẵn) trong cấu trúc ngôn ngữ C không thể được truy cập bằng ngôn ngữ Go. Nếu bạn cần thao tác với các thành phần này, bạn cần xác định hàm hỗ trợ trong ngôn ngữ C.

Đối với các thành phần của mảng có độ dài bằng 0, các phần tử của mảng không thể truy cập trực tiếp trong Go, nhưng vẫn có thể truy cập phần bù vị trí (offset) của phần tử trong mảng có độ dài bằng 0 thông qua `unsafe.Offsetof(a.arr)`.

```go
/*
struct A {
    int   size: 10; // Trường bit không thể truy cập
    float arr[];    // Mạng có độ dài bằng 0 cũng không thể truy cập được
};
*/
import "C"
import "fmt"

func main() {
    var a C.struct_A
    fmt.Println(a.size) // Lỗi không thể truy cập trường bit
    fmt.Println(a.arr)  // Lỗi mảng có độ dài bằng 0
}
```

Trong ngôn ngữ C, chúng ta không thể truy cập trực tiếp vào kiểu struct được xác định bởi ngôn ngữ Go.

Đối với các kiểu union, chúng ta có thể truy cập các kiểu `union xxx`  tương ứng là  `C.union_xxx` trong ngôn ngữ C. Tuy nhiên, các kiểu union trong C không được hỗ trợ trong Go và chúng được chuyển đổi thành các mảng byte có kích thước tương ứng.

```go
/*
#include <stdint.h>

union B1 {
    int i;
    float f;
};

union B2 {
    int8_t i8;
    int64_t i64;
};
*/
import "C"
import "fmt"

func main() {
    var b1 C.union_B1;
    fmt.Printf("%T\n", b1) // [4]uint8

    var b2 C.union_B2;
    fmt.Printf("%T\n", b2) // [8]uint8
}
```

Nếu bạn cần thao tác biến kiểu lồng nhau trong C (union), thường có ba phương pháp: cách thứ nhất là xác định hàm hỗ trợ  trong  C, cách thứ hai là giải mã thủ công các thành phần thông qua "encoding/binary" của ngôn ngữ Go (không phải vấn đề big endian), thứ ba là sử dụng package `unsafe` để chuyển sang kiểu tương ứng (đây là cách tốt nhất để thực hiện). Sau đây cho thấy cách truy cập các thành viên kiểu union thông qua package `unsafe`:

```go
/*
#include <stdint.h>

union B {
    int i;
    float f;
};
*/
import "C"
import "fmt"

func main() {
    var b C.union_B;
    fmt.Println("b.i:", *(*C.int)(unsafe.Pointer(&b)))
    fmt.Println("b.f:", *(*C.float)(unsafe.Pointer(&b)))
}
```

Mặc dù truy cập bằng package `unsafe` là cách dễ nhất và tốt nhất về hiệu suất, nó có thể làm phức tạp vấn đề với các tình huống mà trong đó các kiểu union lồng nhau được xử lý. Đối với các kiểu này ta nên xử lý chúng bằng cách xác định các hàm hỗ trợ trong ngôn ngữ C.

Đối với các kiểu liệt kê (enum), chúng ta có thể truy cập các kiểu `enum xxx` tương ứng là `C.enum_xxx` trong C.

```go
/*
enum C {
    ONE,
    TWO,
};
*/
import "C"
import "fmt"

func main() {
    var c C.enum_C = C.TWO
    fmt.Println(c)
    fmt.Println(C.ONE)
    fmt.Println(C.TWO)
}
```

Trong ngôn ngữ C, kiểu `int` bên dưới kiểu liệt kê hỗ trợ giá trị âm. Chúng ta có thể truy cập trực tiếp các giá trị liệt kê được xác định bằng `C.ONE`, `C.TWO`, v.v.

## 2.3.4 Array, String và Slice

Trong C, tên mảng thực sự tương ứng với một con trỏ tới một phần bộ nhớ có độ dài cụ thể của một loại cụ thể, nhưng con trỏ này không thể được sửa đổi, khi chuyển tên mảng cho một hàm, nó thực sự chuyển phần tử đầu tiên của mảng. Địa chỉ. Để thảo luận, chúng tôi sẽ đề cập đến một độ dài nhất định của bộ nhớ là một mảng. Chuỗi ngôn ngữ C là một mảng kiểu char và độ dài của chuỗi cần được xác định theo vị trí của ký tự NULL cho biết kết thúc. Không có loại lát trong ngôn ngữ C.

Trong Go, một mảng là một loại giá trị và độ dài của mảng là một phần của loại mảng. Chuỗi ngôn ngữ Go tương ứng với một độ dài 

```go
// Go string to C string
// The C string is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be
// freed, such as by calling C.free (be sure to include stdlib.h
// if C.free is needed).
func C.CString(string) *C.char

// Go []byte slice to C array
// The C array is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be
// freed, such as by calling C.free (be sure to include stdlib.h
// if C.free is needed).
func C.CBytes([]byte) unsafe.Pointer

// C string to Go string
func C.GoString(*C.char) string

// C data with explicit length to Go string
func C.GoStringN(*C.char, C.int) string

// C data with explicit length to Go []byte
func C.GoBytes(unsafe.Pointer, C.int) []byte
```

Đối với C.CStringchuỗi Go đầu vào, sao chép chuỗi định dạng ngôn ngữ C, chuỗi trả về được mallocgán bởi chức năng ngôn ngữ C và cần được freephát hành bởi chức năng ngôn ngữ C khi không sử dụng . C.CBytesHàm và C.CStringcác hàm tương tự được sử dụng để sao chép phiên bản ngôn ngữ C của một mảng byte từ lát byte ngôn ngữ Go đầu vào. C.GoStringĐược sử dụng để sao chép chuỗi ngôn ngữ C từ chuỗi ngôn ngữ C kết thúc NULL. C.GoStringNLà một mảng nhân vật khác chức năng nhân bản. C.GoBytesĐược sử dụng để sao chép một lát byte ngôn ngữ Go từ một mảng ngôn ngữ C.

Tập hợp các hàm trợ giúp này được chạy trong chế độ sao chép. Khi string và slice ngôn ngữ Go được chuyển đổi thành C, bộ nhớ nhân bản được malloccấp phát bởi chức năng ngôn ngữ C và cuối cùng có thể được freegiải phóng bởi chức năng. Khi một chuỗi hoặc mảng ngôn ngữ C được chuyển đổi thành Go, bộ nhớ nhân bản được quản lý bởi ngôn ngữ Go. Với bộ chức năng chuyển đổi này, bộ nhớ trước chuyển đổi và sau chuyển đổi vẫn ở trong các địa phương tương ứng của chúng và chúng không trải rộng các ngôn ngữ Go và C. Ưu điểm của chuyển đổi chế độ nhân bản là quản lý giao diện và bộ nhớ rất đơn giản. Nhược điểm là nhân bản cần phân bổ bộ nhớ mới và các hoạt động sao chép sẽ dẫn đến chi phí bổ sung.

reflectCó các định nghĩa cho string và slice trong package:

```go
type StringHeader struct {
    Data uintptr
    Len  int
}

type SliceHeader struct {
    Data uintptr
    Len  int
    Cap  int
}
```

Nếu bạn không muốn phân bổ bộ nhớ riêng, bạn có thể truy cập trực tiếp vào không gian bộ nhớ của ngôn ngữ C bằng ngôn ngữ Go:

```go
/*
static char arr[10];
static char *s = "Hello";
*/
import "C"
import "fmt"

func main() {
    // 通过 reflect.SliceHeader 转换
    var arr0 []byte
    var arr0Hdr = (*reflect.SliceHeader)(unsafe.Pointer(&arr0))
    arr0Hdr.Data = uintptr(unsafe.Pointer(&C.arr[0]))
    arr0Hdr.Len = 10
    arr0Hdr.Cap = 10

    // 通过切片语法转换
    arr1 := (*[31]byte)(unsafe.Pointer(&C.arr[0]))[:10:10]

    var s0 string
    var s0Hdr = (*reflect.StringHeader)(unsafe.Pointer(&s0))
    s0Hdr.Data = uintptr(unsafe.Pointer(C.s))
    s0Hdr.Len = int(C.strlen(C.s))

    sLen := int(C.strlen(C.s))
    s1 := string((*[31]byte)(unsafe.Pointer(&C.s[0]))[:sLen:sLen])
}
```

Vì chuỗi ngôn ngữ Go là chỉ đọc, người dùng cần đảm bảo rằng nội dung của chuỗi C bên dưới sẽ không thay đổi trong quá trình sử dụng chuỗi Go và bộ nhớ sẽ không được giải phóng trước.

Trong CGO, phiên bản ngôn ngữ C của cấu trúc tương ứng với cấu trúc trên được tạo cho string và slice:

```go
typedef struct { const char *p; GoInt n; } GoString;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;
```

Trong ngôn ngữ C có thể GoStringvà GoSliceđể truy cập chuỗi và cắt Go ngôn ngữ. Nếu nó là một kiểu mảng trong Go, bạn có thể chuyển đổi mảng thành một lát và sau đó chuyển đổi nó. Nếu không gian bộ nhớ cơ bản tương ứng với một chuỗi hoặc lát được quản lý bởi thời gian chạy của ngôn ngữ Go, thì đối tượng bộ nhớ Go có thể được lưu trong một thời gian dài trong ngôn ngữ C.

Chi tiết về mô hình bộ nhớ CGO sẽ được thảo luận chi tiết hơn trong các chương sau.

## 2.3.5 Chuyển đổi giữa các con trỏ

Trong ngôn ngữ C, các loại con trỏ khác nhau có thể được chuyển đổi rõ ràng hoặc ngầm định. Nếu nó ẩn, nó sẽ chỉ đưa ra một số thông tin cảnh báo tại thời điểm biên dịch. Nhưng ngôn ngữ Go rất nghiêm ngặt đối với các loại chuyển đổi khác nhau và mọi thông báo cảnh báo có thể xuất hiện trong ngôn ngữ C có thể sai trong ngôn ngữ Go! Con trỏ là linh hồn của ngôn ngữ C và việc chuyển đổi miễn phí giữa các con trỏ cũng là vấn đề quan trọng đầu tiên thường được giải quyết trong mã cgo.

Trong ngôn ngữ Go, hai con trỏ hoàn toàn giống nhau và có thể được sử dụng trực tiếp mà không cần chuyển đổi. Nếu một loại con trỏ được xây dựng bên trên một loại con trỏ khác bằng lệnh loại, nói cách khác, hai con trỏ bên dưới là các con trỏ có cùng cấu trúc, sau đó chúng ta có thể chuyển đổi giữa các con trỏ bằng cú pháp truyền trực tiếp. Tuy nhiên, cgo thường phải đối phó với việc chuyển đổi giữa hai loại con trỏ hoàn toàn khác nhau. Về nguyên tắc, thao tác này bị nghiêm cấm trong mã ngôn ngữ thuần túy.

Một trong những mục đích của cgo là phá vỡ sự cấm đoán của ngôn ngữ Go và khôi phục các hoạt động chuyển đổi và con trỏ miễn phí của các con trỏ mà ngôn ngữ C nên có. Đoạn mã sau trình bày cách chuyển đổi một con trỏ loại X thành một con trỏ loại Y:

```go
var p *X
var q *Y

q = (*Y)(unsafe.Pointer(p)) // *X => *Y
p = (*X)(unsafe.Pointer(q)) // *Y => *X
```

Để chuyển đổi con trỏ loại X thành con trỏ loại Y, chúng ta cần unsafe.Pointerthực hiện chuyển đổi giữa các loại con trỏ khác nhau như một loại cầu nối trung gian. unsafe.Pointerkiểu con trỏ tương tự với ngôn ngữ C void*kiểu của con trỏ.

Sau đây là sơ đồ quy trình chuyển đổi giữa các con trỏ:

![x](../images/ch2-1-x-ptr-to-y-ptr.uml.png)

*Hình 2-1 Con trỏ loại X đến con trỏ loại Y*

Bất kỳ loại con trỏ nào cũng có thể được chuyển sang unsafe.Pointerloại con trỏ để loại bỏ thông tin loại ban đầu, sau đó gán lại một loại con trỏ mới để đạt được mục đích chuyển đổi giữa các con trỏ.

![x2](../images/ch2-2-int32-to-char-ptr.uml.png)

*Hình 2-2 Int32 và `char` chuyển đổi con trỏ*

Việc chuyển đổi được chia thành nhiều giai đoạn và một mục tiêu nhỏ được thực hiện ở mỗi giai đoạn: đầu tiên là kiểu int32 sang uintptr, sau đó là uintptr thành unsafe.Pointrloại con trỏ và cuối cùng là unsafe.Pointrloại con trỏ thành *C.charloại.

## 2.3.7 Chuyển đổi giữa kiểu slice

Mảng cũng là một loại con trỏ trong ngôn ngữ C, vì vậy việc chuyển đổi giữa hai loại mảng khác nhau về cơ bản tương tự như chuyển đổi giữa các con trỏ. Tuy nhiên, trong ngôn ngữ Go, lát tương ứng với một mảng hoặc một mảng không còn là loại con trỏ, vì vậy chúng ta không thể chuyển đổi trực tiếp giữa các loại lát khác nhau.

Tuy nhiên, package phản chiếu của ngôn ngữ Go cung cấp cấu trúc cơ bản của loại lát cắt []Xvà []Ychuyển đổi lát có thể được thực hiện và nhập kết hợp với kỹ thuật chuyển đổi con trỏ được thảo luận ở trên giữa các loại khác nhau :

```go
var p []X
var q []Y

pHdr := (*reflect.SliceHeader)(unsafe.Pointer(&p))
qHdr := (*reflect.SliceHeader)(unsafe.Pointer(&q))

pHdr.Data = qHdr.Data
pHdr.Len = qHdr.Len * unsafe.Sizeof(q[0]) / unsafe.Sizeof(p[0])
pHdr.Cap = qHdr.Cap * unsafe.Sizeof(q[0]) / unsafe.Sizeof(p[0])
```

Ý tưởng chuyển đổi giữa các loại lát khác nhau là trước tiên xây dựng một lát đích trống, sau đó điền vào lát cắt đích với dữ liệu cơ bản của lát cắt gốc. Nếu loại X và Y có kích thước khác nhau, bạn cần đặt lại thuộc tính Len và Cap. Cần lưu ý rằng nếu X hoặc Y là loại null, mã trên có thể gây ra lỗi chia cho 0 và mã thực tế cần được xử lý khi thích hợp.

Sau đây cho thấy luồng cụ thể của chuyển đổi giữa các lát:

![slicexy](../images/ch2-3-x-slice-to-y-slice.uml.png)

*Hình 2-3 kiểu cắt X thành lát cắt Y*

Đối với các tính năng thường được sử dụng trong CGO, tác giả package gọn package "github.com/chai2010/cgo", cung cấp các chức năng chuyển đổi cơ bản. Để biết chi tiết, hãy tham khảo mã thực hiện.
