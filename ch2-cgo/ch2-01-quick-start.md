# 2.1. Quick Start

Trong phần này, chúng ta sẽ tìm hiểu cách sử dụng CGO cơ bản  thông qua loạt ví dụ từ đơn giản đến phức tạp.

## 2.1.0. Cài đặt

***Linux, Mac***

- CGO được tích hợp sẵn, người dùng Linux và Mac chỉ cần cài đặt Go.
  
***Windows***

- CGO không tương thích với Cygwin
- Để dùng được CGO ban5 phải cài đặt MingGW và sử dụng cmd của Windows

## 2.1.1. Chương trình CGO đơn giản

Đầu tiên là một chương trình CGO đơn giản nhất:

```go
//main.go
package main

import "C"

func main() {
    println("hello cgo")
}
```

Chúng ta import package CGO thông qua câu lệnh `import "C"`. Chương trình trên chưa thực hiện bất kì thao tác nào với CGO, chỉ mới thông báo sẵn sàng cho việc lập trình với CGO. Mặc dù chưa sử dụng gì đến CGO nhưng lệnh `go build` vẫn sẽ gọi trình biên dịch `gcc` trong suốt quá trình biên dịch do đây được là một chương trình CGO hoàn chỉnh.

## 2.1.2. Xuất chuỗi dựa trên thư viện chuẩn của C

```go
//main.go
package main

//#include <stdio.h>
import "C"

func main() {
    C.puts(C.CString("Hello World\n"))  
}

// kết quả:
// Hello World
```

`import package "C"` để thực hiện các chức năng của CGO và include thư viện <stdio.h> của ngôn ngữ C (trong comment code là cú pháp của CGO, không phải phần comment như thông thường). Tiếp theo, chuỗi string trong `C.CString` của ngôn ngữ Go được chuyển đổi thành chuỗi string trong ngôn ngữ C bằng phương thức `C.puts` của package CGO.

Việc lỗi xảy ra khi không giải phóng chuỗi được tạo bằng C.CString của ngôn ngữ C sẽ dẫn đến rò rỉ bộ nhớ (memory leak). Nhưng đối với chương trình nhỏ ở trên điều này không đáng lo ngại  vì hệ điều hành sẽ tự động lấy lại các tài nguyên của chương trình sau khi chương trình kết thúc.

## 2.1.3. Sử dụng hàm C tự khai báo

Phần trên chúng ta đã sử dụng các hàm đã có trong `stdio`. Bây giờ ta sẽ sử dụng một hàm `SayHello` của ngôn ngữ C. Chức năng hàm này là in ra chuỗi chúng ta truyền vào hàm. Sau đó gọi hàm `SayHello` trong hàm main:

```go
//main.go
package main

/*
#include <stdio.h>
// Khai báo hàm trong ngôn ngữ C
static void SayHello(const char* s) {
    puts(s);
}
*/
import "C"

func main() {
    C.SayHello(C.CString("Hello World\n"))
}
```

Hoặc có thể đặt hàm `SayHello` trong file `hello.c` như sau:

```C
//hello.c
#include <stdio.h>

void SayHello(const char* s) {
    puts(s);
}
```

Sau đó bên file main.go chúng chỉ cần khai báo hàm `SayHello` trong phần CGO như bên dưới.

```Go
//main.go
package main

//void SayHello(const char* s);
import "C"

func main() {
    C.SayHello(C.CString("Hello World\n"))
}
```

Hiện tại thư mục làm việc có cấu trúc như sau:

```sh
.
├── hello.c
└── main.go
```

`Lưu ý`: thay vì chạy lệnh `go run main.go` hoặc `go build main.go`, chúng ta phải sử dụng `go run "tên/của/package"` hoặc `go build "tên/của/package"`. Nếu đang đứng trong thư mục chứa mã nguồn thì bạn có thể chạy chương trình bằng lệnh `go run .` hoặc `go build .`

## 2.1.4. Module hóa C code

Trừu tượng và module hóa là cách để đơn giản hóa các vấn đề trong lập trình:

- Khi code quá dài, ta có thể đưa các lệnh tương tự nhau vào chung một hàm.
- Khi có nhiều hàm hơn, ta chia chúng vào các file hoặc module.

Trong ví dụ trước, ta trừu tượng hóa một module tên là `hello` và tất cả các interface của module đó được khai báo trong file header `hello.h`:

```c
//hello.h
void SayHello(const char* s);
```

Và hiện thực hàm `SayHello` trong file `hello.c`:

```c
//hello.c
#include "hello.h"
#include <stdio.h>

// Đảm bảo việc hiện thực hàm thỏa mãn interface của module.
void SayHello(const char* s) {
    puts(s);
}
```

Ngoài ra ta có thể hiện thực hàm này bằng C++ cũng được:

```c
//hello.cpp
#include <iostream>

// extern giúp function C++ có được các liên kết (linkage)
// của C. Chỉ ra rằng liên hệ giữa hello.cpp và hello.h
// vẫn theo quy tắc của C
extern "C" {
    #include "hello.h"
}

void SayHello(const char* s) {
    std::cout << s;
}
```

Trong hàm main của Go ta gọi file header như sau:

```go
//main.go
package main

//#include <hello.h>
import "C"

func main() {
    C.SayHello(C.CString("Hello World"))
}
```

Với việc lập trình C thông qua interface, ta có thể hiện thực module bằng nhiều ngôn ngữ khác nhau, miễn là đáp ứng được interface: SayHello có thể được viết bằng C, C++, Go hoặc kể cả Assembly.

## 2.1.5. Sử dụng Go để hiện thực hàm trong C

Trong thực tế, CGO không chỉ được sử dụng để gọi các hàm của C bằng ngôn ngữ Go mà còn có thể cho phép C gọi các hàm Go trong code của mình.

Trong ví dụ trước, chúng ta đã trừu tượng hóa một module có tên hello và tất cả các chức năng interface của module được xác định trong file header `hello.h`:

```c
//hello.h
void SayHello(const char* s);
```

Bây giờ, chúng ta tạo một file `hello.go` và hiện thực lại hàm `SayHello` của interface bằng ngôn ngữ Go:

```go
//hello.go
package main

import "C"

import "fmt"

//export SayHello
func SayHello(s *C.char) {
    fmt.Print(C.GoString(s))
}
```

Sử dụng chỉ thị `//export SayHello` của CGO để export hàm được hiện thực bằng Go sang hàm sử dụng được cho C.

Cần chú ý là ta sẽ có hai phiên bản `SayHello`: một là trong môi trường cục bộ của Go, hai là của C. Phiên bản `SayHello` của C được sinh ra bởi CGO cuối cùng cũng sẽ gọi phiên bản `SayHello` của Go thông qua `bridge code`.

Với việc lập trình ngôn ngữ C qua inteface, ta có thể tự do hiện thực và đơn giản hóa việc sử dụng hàm. Bây giờ ta có thể dùng SayHello như là một thư viện:

```go
package main

//#include <hello.h>
import "C"

func main() {
    C.SayHello(C.CString("Hello World\n"))
}
```

Cấu trúc thư mục lúc này:

```sh
.
├── hello.go
├── hello.h
└── main.go
```

## 2.1.6. Sử dụng Go để lập trình interface cho C

Để cho đơn giản chúng ta sẽ gộp tất cả thành một file `main.go` duy nhất như ví dụ dưới đây.

```go
//main.go
package main

//void SayHello(char* s);
import "C"

import (
    "fmt"
)

func main() {
    C.SayHello(C.CString("Hello World\n"))
}

//export SayHello
func SayHello(s *C.char) {
    fmt.Print(C.GoString(s))
}
```

Tỉ lệ code C trong chương trình bây giờ ít hơn. Tuy nhiên vẫn phải sử dụng chuỗi trong C thông qua hàm `C.CString` chứ không thể dùng trực tiếp chuỗi của Go. Trong `Go1.10`, CGO đã thêm kiểu `_GoString_pred` để thể hiện chuỗi trong ngôn ngữ Go. Sau đây là code đã được cải tiến:

```go
// +build go1.10

package main

//void SayHello(_GoString_ s);
import "C"

import (
    "fmt"
)

func main() {
    C.SayHello("Hello World\n")
}

//export SayHello
func SayHello(s string) {
    fmt.Print(s)
}
```

Có vẻ như tất cả đều được viết bằng Go, nhưng việc triển khai từ hàm `main()` của ngôn ngữ Go đến phiên bản ngôn ngữ C đã tự động tạo ra hàm `SayHello`, rồi cuối cùng trở lại môi trường ngôn ngữ Go.

[Tiếp theo](ch2-01-quick-start.md)