# 2.4 Lời gọi hàm

Thông qua công cụ CGO, chúng ta không chỉ có thể gọi hàm của ngôn ngữ C bằng Go mà còn có thể export hàm của Go để sử dụng như là hàm ngôn ngữ C.

## 2.4.1 Go gọi hàm C

Đối với chương trình kích hoạt các tính năng CGO, CGO xây dựng một package C ảo. Hàm ngôn ngữ C có thể được gọi qua package C ảo này.

```go
/*
static int add(int a, int b) {
    return a+b;
}
*/
import "C"

func main() {
    C.add(1, 1)
}
```

Code CGO ở trên trước tiên xác định hàm `add` hiển thị trong file hiện tại và sau đó chuyển sang `C.add`.

## 2.4.2 Giá trị trả về của hàm C

Đối với hàm của C có giá trị trả về, chúng ta có thể nhận giá trị trả về bình thường.

```go
/*
static int div(int a, int b) {
    return a/b;
}
*/
import "C"
import "fmt"

func main() {
    v := C.div(6, 3)
    fmt.Println(v)
}
```

Hàm `div` ở trên thực hiện một phép toán chia số nguyên và trả về kết quả của phép chia.

Tuy nhiên, không có cách xử lý đặc biệt nào cho trường hợp số chia là 0. Vì ngôn ngữ C không hỗ trợ trả về nhiều kết quả, thư viện chuẩn <errno.h> cung cấp macro `errno` để trả về trạng thái lỗi. Nếu bạn muốn trả về lỗi khi số chia là 0 còn những lần khác trả về kết quả bình thường. Chúng ta có thể xem  `errno` là một biến toàn cục thread-safe có thể được sử dụng để ghi lại mã trạng thái của lỗi đây đây nhất.

Hàm `div` cải tiến được hiện thực như sau:

```c
#include <errno.h>

int div(int a, int b) {
    if(b == 0) {
        errno = EINVAL;
        return 0;
    }
    return a/b;
}
```

CGO cũng có hỗ trợ đặc biệt cho các macro `errno`  thuộc thư viện tiêu chuẩn <errno.h>: nếu có hai giá trị trả về khi CGO gọi hàm C thì giá trị trả về thứ hai sẽ tương ứng với trạng thái lỗi `errno`.

```go
/*
#include <errno.h>

static int div(int a, int b) {
    if(b == 0) {
        errno = EINVAL;
        return 0;
    }
    return a/b;
}
*/
import "C"
import "fmt"

func main() {
    v0, err0 := C.div(2, 1)
    fmt.Println(v0, err0)

    v1, err1 := C.div(1, 0)
    fmt.Println(v1, err1)
}
```

Thực thi đoạn code trên sẽ cho output như sau:

```sh
2 <nil>
0 invalid argument
```

Chúng ta có thể xem hàm `div` tương ứng với một hàm trong Go như sau:

```go
func C.div(a, b C.int) (C.int, error)
```

Tham số thứ hai trả về (giá trị error) có thể bỏ qua, được hiện thực bên dưới là kiểu `syscall.Errno`.

## 2.4.3 Giá trị trả về của hàm void

Trong C cũng có hàm không trả về kiểu giá trị (thay vào đó trả về void). Chúng ta không thể nhận được giá trị trả về của hàm kiểu void vì đó thực sự không phải giá trị!?.

Như đã đề cập trong ví dụ trước, CGO hiện thực một phương pháp đặc biệt cho errno và có thể nhận về  trạng thái lỗi của ngôn ngữ C thông qua giá trị trả về thứ hai. Tính năng này vẫn hợp lệ cho các hàm kiểu void. Đoạn code sau để lấy mã trạng thái lỗi của hàm không có giá trị trả về:

```go
//static void noreturn() {}
import "C"
import "fmt"

func main() {
    _, err := C.noreturn()
    fmt.Println(err)
}

// kết quả: <nil>
```

Lúc này, chúng ta bỏ qua giá trị trả về đầu tiên và chỉ nhận được mã lỗi tương ứng với giá trị trả về thứ hai.

Chúng ta cũng có thể thử lấy giá trị trả về đầu tiên, cũng chính là kiểu tương ứng với kiểu void trong ngôn ngữ C:

```go
//static void noreturn() {}
import "C"
import "fmt"

func main() {
    v, _ := C.noreturn()
    fmt.Printf("%#v", v)
}
```

Chạy code này sẽ thu được kết quả:

```sh
main._Ctype_void{}
```

Chúng ta có thể thấy rằng kiểu void của ngôn ngữ C tương ứng với kiểu trong package main  `_Ctype_void`. Trong thực tế, hàm `noreturn` của ngôn ngữ C cũng được coi là một hàm với kiểu trả về `_Ctype_void`, do đó bạn có thể trực tiếp nhận giá trị trả về của hàm kiểu void:

```go
//static void noreturn() {}
import "C"
import "fmt"

func main() {
    fmt.Println(C.noreturn())
}
```

Chạy code này sẽ cho ra kết quả sau:

```sh
[]
```

Trong thực tế, trong code được CGO tạo ra, kiểu `_Ctype_void` tương ứng với kiểu mảng có độ dài 0 (`[0]byte`), do đó output `fmt.Println` là một cặp dấu ngoặc vuông biểu thị một giá trị null.

## 2.4.4 C gọi hàm do Go export

CGO có một tính năng mạnh mẽ là export các hàm Go thành các hàm ngôn ngữ C. Trong trường hợp này, chúng ta có thể định nghĩa interface bằng C và sau đó triển khai nó thông qua Go. Trong phần đầu tiên của chương này cũng đã có một số ví dụ về hàm ngôn ngữ C gọi hàm do Go export.

Nhắc lại một chút về hàm `add` trong chương đầu:

```go
import "C"

//export add
func add(a, b C.int) C.int {
    return a+b
}
```

Tên hàm `add` bắt đầu bằng một chữ cái viết thường và là một hàm private trong package của Go. Nhưng theo cái nhìn của ngôn ngữ C thì hàm `add` là hàm ngôn ngữ C có thể được truy cập toàn cục. Nếu có một hàm `add` cùng tên được export dưới dạng hàm ngôn ngữ C trong hai package ngôn ngữ Go khác nhau, vấn đề về trùng tên sẽ xảy ra trong link phase.

Chúng ta có thể include file header `_cgo_export.h` để thêm tham chiếu đến các hàm export. Nếu bạn muốn sử dụng ngay lập tức hàm `add` của C được export trong file CGO hiện tại, bạn không thể tham khảo đến file `_cgo_export.h`,  bởi vì việc tạo file `_cgo_export.h`  cần phụ thuộc vào file hiện tại mà trong file hiện tại lại tham khảo tới `_cgo_export.h` (file chưa được tạo) thì sẽ gây ra lỗi.

```c
#include "_cgo_export.h"

void foo() {
    add(1, 1);
}
```

Khi export interface của ngôn ngữ C, ta cần đảm bảo rằng các tham số hàm và kiểu giá trị trả về là kiểu "thân thiện" với C (xem lại [2.3](./ch2-03-type-conversion.md)) đồng thời giá trị trả về không được trực tiếp hoặc gián tiếp chứa con trỏ vào không gian bộ nhớ ngôn ngữ Go.

[Tiếp theo](ch2-05-internal-mechanisms.md)