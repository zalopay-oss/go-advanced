# 2.2. CGO Foundation

Để sử dụng tính năng CGO, bạn cần cài đặt compiler C/C++ trên macOS và Linux là  `GCC` còn trên Windows là `MinGW`. Đồng thời, cần đảm bảo rằng biến môi trường `CGO_ENABLED` được đặt thành 1.

## 2.2.1. Lệnh `import "C"`

Lệnh  `import "C"` xuất hiện trong code nói cho compiler rằng tính năng CGO sẽ được sử dụng. Code trong cặp `/* */` trước lệnh đó là cú pháp để Go nhận ra code của C. Lúc này, ta có thể thêm các file code của C/C++ tương ứng trong thư mục hiện tại.

Ví dụ:

```go
package main

/*
#include <stdio.h>

void printint(int v) {
    printf("printint: %d\n", v);
}
*/
import "C"

func main() {
    v := 42

    // int(v) chuyển đổi giá trị kiểu int trong Go đến giá trị
    // kiểu int trong ngôn ngữ C
    C.printint(C.int(v))
}
```

Tất cả các thành phần ngôn ngữ C trong file header sau khi được  include sẽ được thêm vào package "C" ảo. Cần lưu ý rằng câu lệnh `import "C"` yêu cầu một dòng riêng và không thể được import cùng với các package khác.

Cần lưu ý rằng Go là ngôn ngữ ràng buộc kiểu mạnh, do đó tham số được truyền phải đúng kiểu khai báo, và phải được chuyển đổi sang kiểu trong C bằng các hàm chuyển đổi trước khi truyền, không thể truyền trực tiếp bằng kiểu của Go.

Các ký hiệu của C được import thông qua package C thì *không cần phải viết hoa*, không cần phải tuân theo quy tắc của Go.

Ví dụ tiếp theo ta định nghĩa kiểu `CChar` tương ứng với con trỏ char của C trong Go và sau đó thêm phương thức `GoString` để trả về chuỗi trong ngôn ngữ Go:

```go
package cgo_helper

//#include <stdio.h>
import "C"

type CChar C.char

// nhận vào CChar của C và trả về
// chuỗi string của Go
func (p *CChar) GoString() string {
    return C.GoString((*C.char)(p))
}

func PrintCString(cs *C.char) {
    C.puts(cs)
}
```

Bây giờ có thể ta muốn sử dụng hàm này trong các package Go khác:

```go
package main

//static const char* cs = "hello";
import "C"
import "./cgo_helper"

func main() {

    // C.cs là kiểu của package C xây dựng
    // trên kiểu *char (thực ra là *main.C.char)
    cgo_helper.PrintCString(C.cs)
    // kiểu *C.type của hàm PrintCString trong
    // cgo_helper là *cgo_helper.C.char
    // 2 kiểu này không giống nhau
    // --> Code lỗi
}
```

Các tham số được chuyển đổi kiểu rồi mới truyền vào có được không?

Câu trả lời là không, bởi vì các tham số của `cgo_helper.PrintCString` là kiểu `*C.char` được định nghĩa trong package riêng của nó và không thể truy cập trực tiếp từ bên ngoài.

Nói cách khác, nếu một package sử dụng trực tiếp một kiểu thuộc package C ảo (tương tự như `*C.char`)  trong một interface chung thì các package khác sẽ không thể sử dụng trực tiếp các kiểu đó trừ khi package kia cũng cung cấp hàm tạo `*C.chartype`.

## 2.2.2. Lệnh `#cgo`

Trước dòng `import "C"` ta có thể đặt các tham số cho quá trình biên dịch (compile phase) và quá trình liên kết (link phase) thông qua các lệnh `#cgo`.

- Các tham số của quá trình biên dịch chủ yếu được sử dụng để xác định các macro liên quan và đường dẫn truy xuất file header đã chỉ định.
- Các tham số của quá trình liên kết chủ yếu là để xác định đường dẫn truy xuất file thư viện và file thư viện sẽ được liên kết.

```go
// Định nghĩa macro PNG_DEBUG, giá trị là 1
// #cgo CFLAGS: -DPNG_DEBUG=1 -I./include
// #cgo LDFLAGS: -L/usr/local/lib -lpng
// #include <divng.h>
import "C"
```

Trong đoạn mã trên:

- Phần `CFLAGS`: `-D` định nghĩa macro `PNG_DEBUG`, giá trị là 1
- `-I` xác định thư mục tìm kiếm có trong file header.
- Phần `DFLAGS`: `-L` chỉ ra thư mục truy xuất các file thư viện, `-l` chỉ định thư viện png là bắt buộc.

Do các vấn đề mà C/C ++ để lại, đường dẫn truy xuất file header C có thể là *relative path*, nhưng đường dẫn truy xuất file thư viện bắt buộc phải là *absolute path*. Absolute path của thư mục package hiện tại có thể được biểu diễn bằng  biến `${SRCDIR}` trong thư mục truy xuất các file thư viện:

```c
// #cgo LDFLAGS: -L${SRCDIR}/libs -lfoo
```

Đoạn code trên sẽ được phân giải trong link phase và trở thành:

```c
// #cgo LDFLAGS: -L/go/src/foo/libs -lfoo
```

Lệnh `#cgo` chủ yếu ảnh hưởng đến một số biến môi trường của trình biên dịch như CFLAGS, CPPFLAGS, CXXFLAGS, FFLAGS và LDFLAGS. LDFLAGS được sử dụng để đặt tham số của liên kết, CFLAGS được sử dụng để đặt tham số biên dịch cho code ngôn ngữ C.

Đối với người dùng sử dụng C và C++ trong môi trường CGO, có thể có ba tùy chọn biên dịch khác nhau:

- CFLAGS cho các tùy chọn biên dịch theo ngôn ngữ C.
- CPPFLAGS cho các tùy chọn biên dịch cụ thể C++.
- CXXFLAGS cho các biên dịch C và C++.

Các lệnh `#cgo` cũng hỗ trợ  tùy chọn biên dịch hoặc liên kết với các hệ điều hành hoặc một kiểu kiến trúc CPU khác nhau:

```go
// tuỳ chọn cho Windows
// #cgo windows CFLAGS: -DX86=1

// tuỳ chọn cho non-windows platforms
// #cgo !windows LDFLAGS: -lm
```

Một ví dụ để xác định hệ thống nào đang chạy CGO:

```go
package main

/*
#cgo windows CFLAGS: -DCGO_OS_WINDOWS=1
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo linux CFLAGS: -DCGO_OS_LINUX=1

#if defined(CGO_OS_WINDOWS)
    const char* os = "windows";
#elif defined(CGO_OS_DARWIN)
    const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
    const char* os = "linux";
#else
#    error(unknown os)
#endif
*/
import "C"

func main() {
    print(C.GoString(C.os))
}
```

Bằng cách này, chúng ta có thể biết được hệ thống mà code đang vận hành, nhờ đó áp dụng các kĩ thuật riêng cho các nền tảng khác nhau.

## 2.2.3. Biên dịch với tag

Build tag là một comment đặc biệt ở đầu file C/C++ trong môi trường Go hoặc CGO. Biên dịch có điều kiện tương tự như sử dụng macro `#cgo` để xác định các nền tảng khác nhau (ví dụ trên). Code được build sau khi macro của nền tảng tương ứng được xác định.

Ví dụ trình bày một cách khác khi các file nguồn sau sẽ chỉ được tạo khi debug build flag  được thiết lập:

```go
// +build debug

package main

var buildMode = "debug"
```

Có thể được build bằng lệnh sau:

```sh
go build -tags="debug"
go build -tags="windows debug"
```

Chúng ta có thể dùng `-tags` chỉ định nhiều build flag cùng một lúc thông qua các đối số dòng lệnh.

Ví dụ các build flag sau chỉ ra rằng việc build chỉ được thực hiện trong kiến trúc "linux/386" hoặc "non-cgo environment" trong nền tảng darwin.

```go
// +build linux,386 darwin,!cgo
```

Trong đó, dấu phẩy (`,`) nghĩa là **và**. Khoảng trắng (`  `) nghĩa là **hoặc**.

[Tiếp theo](ch2-03-type-conversion.md)