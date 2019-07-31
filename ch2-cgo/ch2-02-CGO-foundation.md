# 2.2 CGO Foundation

<<<<<<< HEAD
Để sử dụng tính năng CGO, bạn cần cài đặt compiler C/C++: Trên macOS và Linux là  `GCC` còn trên Windows là `MinGW`. Đồng thời, cần đảm bảo rằng biến môi trường `CGO_ENABLED` được đặt thành 1.

## 2.2.1 Lệnh `import "C"`

Lệnh  `import "C"` xuất hiện trong code nói cho compiler rằng tính năng CGO sẽ được sử dụng. Code trong cặp `/* */` trước lệnh đó là cú pháp để Go nhận ra code của C. Lúc này, ta có thể thêm các file code của C/C++ tương ứng trong thư mục hiện tại.

Ví dụ:
=======
Để sử dụng tính năng CGO, bạn cần cài đặt công cụ C/C++. Trên macOS và Linux, bạn cần cài đặt `GCC`. Trên Windows, bạn cần cài đặt công cụ `MinGW`. Đồng thời, bạn cần đảm bảo rằng biến môi trường `CGO_ENABLED` được đặt thành 1.

## 2.2.1 Lệnh `import "C"`

Nếu lệnh import `import "C"` xuất hiện trong code, nó có nghĩa là tính năng CGO được sử dụng. Comment trước lệnh đó là cú pháp để Go nhận ra code của C. Khi CGO được bật, bạn có thể thêm các file code của C/C++ tương ứng trong thư mục hiện tại.

Ví dụ đơn giản nhất:
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376

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
<<<<<<< HEAD

    // int(v) chuyển đổi giá trị kiểu int trong Go đến giá trị
    // kiểu int trong ngôn ngữ C
=======
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376
    C.printint(C.int(v))
}
```

<<<<<<< HEAD
Tất cả các thành phần ngôn ngữ C trong file header sau khi được  include sẽ được thêm vào package "C" ảo. Cần lưu ý rằng câu lệnh `import "C"` yêu cầu một dòng riêng và không thể được import cùng với các package khác.

Cần lưu ý rằng Go là ngôn ngữ ràng buộc kiểu mạnh, do đó tham số được truyền phải đúng kiểu khai báo, và phải được chuyển đổi sang kiểu trong C bằng các hàm chuyển đổi trước khi truyền, không thể truyền trực tiếp bằng kiểu của Go.

Các ký hiệu của C được import thông qua package C thì *không cần phải viết hoa*, không cần phải tuân theo quy tắc của Go.

Ví dụ tiếp theo ta định nghĩa kiểu `CChar` tương ứng với con trỏ char của C trong Go và sau đó thêm phương thức `GoString` để trả về chuỗi trong ngôn ngữ Go:
=======
Ví dụ này cho thấy việc sử dụng CGO cơ bản. Phần đầu của comment khai báo hàm C sẽ được gọi và file header được liên kết. Tất cả các thành phần ngôn ngữ C trong file header sau khi được đưa vào sẽ được thêm vào package "C" ảo. `Cần lưu ý rằng câu lệnh import "C" yêu cầu một dòng riêng và không thể được import cùng với các package khác`. Truyền tham số cho hàm C cũng rất đơn giản và nó có thể được chuyển đổi trực tiếp thành một kiểu trong C tương ứng. Trong ví dụ trên, `C.int(v)` được sử dụng để chuyển đổi giá trị kiểu int trong Go đến giá trị kiểu int trong ngôn ngữ C, sau đó gọi hàm printint được xác định bằng ngôn ngữ C để in.

Cần lưu ý rằng Go là ngôn ngữ ràng buộc kiểu mạnh, do đó tham số được truyền phải chính xác với khai báo, và phải được chuyển đổi sang kiểu trong C bằng các hàm chuyển đổi trước khi truyền, không thể truyền trực tiếp bằng kiểu của Go. Đồng thời, các ký hiệu của C được import thông qua package C thì `không cần phải viết hoa`, không cần phải tuân theo quy tắc của Go.

CGO đặt các ký hiệu ngôn ngữ C được tham chiếu bởi package hiện tại vào package C ảo. Đồng thời, các package Go khác mà package hiện tại phụ thuộc cũng có thể giới thiệu các package C ảo tương tự thông qua CGO, nhưng các package Go khác nhau giới thiệu các package ảo. Các kiểu giữa các package C không phải là toàn thể. Ràng buộc này có thể có một tác động nhỏ đến khả năng tự xây dựng một số chức năng CGO.

Ví dụ chúng tôi muốn định nghĩa kiểu `CChar` tương ứng với con trỏ char của C trong Go và sau đó thêm phương thức GoString để trả về chuỗi ngôn ngữ Go
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376

```go
package cgo_helper

//#include <stdio.h>
import "C"

type CChar C.char

<<<<<<< HEAD
// nhận vào CChar của C và trả về
// chuỗi string của Go
=======
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376
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
<<<<<<< HEAD

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

## 2.2.2 Lệnh `#cgo`

Trước dòng `import "C"` ta có thể đặt các tham số cho quá trình biên dịch và quá trình liên kết (compile phase và link phase) thông qua các lệnh `#cgo`.

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

- CFLAGS cho các tùy chọn biên dịch theo ngôn ngữ C,
- CPPFLAGS cho các tùy chọn biên dịch cụ thể C++,
- CXXFLAGS cho các biên dịch C và C++.

Các lệnh `#cgo` cũng hỗ trợ  tùy chọn biên dịch hoặc liên kết với các hệ điều hành hoặc một kiểu kiến trúc CPU khác nhau:

```go
// tuỳ chọn cho Windows
// #cgo windows CFLAGS: -DX86=1

// tuỳ chọn cho non-windows platforms
// #cgo !windows LDFLAGS: -lm
```

Một ví dụ để xác định hệ thống nào đang chạy CGO:
=======
    cgo_helper.PrintCString(C.cs)
}
```

Nhưng đoạn code này sẽ không chạy được. Vì biến `C.cs` được đề cập trong package main hiện tại là kiểu của package ảo C được xây dựng trên `*char` (*C.char, chính xác hơn là *main.C.char), còn kiểu `*C.type` được đề cập đến trong package `cgo_helper` (`*cgo_helper.C.char`) là khác nhau. Trong ngôn ngữ Go, các phương thức phụ thuộc vào kiểu. Các kiểu được package C ảo được đề cập trong các package Go khác nhau là khác nhau (`main.C` không giống `cgo_helper.C`) chính là nguyên nhân khiến các kiểu Go được mở rộng từ chúng thành các kiểu khác nhau (`*main.C.char` khác `*cgo_helper.C.char`). Điều này cuối cùng đã khiến đoạn code đó không hoạt động được.

Người dùng có kinh nghiệm với ngôn ngữ Go có thể đề xuất rằng các tham số được truyền vào sau khi chuyển đổi. Nhưng phương pháp này dường như không khả thi, bởi vì các tham số của `cgo_helper.PrintCString` là kiểu `*C.char` được đề cập bởi package riêng của nó và nó không thể truy cập trực tiếp từ bên ngoài. Nói cách khác, nếu một package trực tiếp sử dụng kiểu C ảo tương tự `*C.char`  trong một interface chung, các package Go khác không thể sử dụng trực tiếp các kiểu này trừ khi package Go cũng cung cấp hàm tạo `*C.chartype`. Do nhiều yếu tố này, nếu bạn muốn kiểm tra các kiểu được export trực tiếp bởi CGO trong môi trường thử nghiệm đi, sẽ có những hạn chế tương tự.

## 2.2.2 Lệnh `#cgo`

Trong dòng ghi chú `import "C"` phía trước các lệnh. Bạn có thể đặt các tham số cho quá trình biên dịch và quá trình liên kết thông qua các lệnh `#cgo`. Các tham số của quá trình biên dịch chủ yếu được sử dụng để xác định các macro liên quan và đường dẫn truy xuất file header đã chỉ định. Các tham số của quá trình liên kết chủ yếu là để xác định đường dẫn truy xuất file thư viện và file thư viện sẽ được liên kết.

```go
// #cgo CFLAGS: -DPNG_DEBUG=1 -I./include   // Định nghĩa macro PNG_DEBUG, giá trị là 1
// #cgo LDFLAGS: -L/usr/local/lib -lpng
// #include <png.h>
import "C"
```

Trong đoạn mã trên, phần CFLAGS, -Dpart định nghĩa macro PNG_DEBUG, giá trị là 1; -I xác định thư mục tìm kiếm có trong file header. Trong phần LDFLAGS, -L thư mục truy xuất file thư viện được -l chỉ định khi liên kết và thư viện png liên kết là bắt buộc khi liên kết được chỉ định.

Do các vấn đề mà C/C ++ để lại, thư mục truy xuất file header C có thể là một thư mục tương đối, nhưng thư mục truy xuất tệp thư viện yêu cầu một đường dẫn tuyệt đối. ${SRCDIR} Đường dẫn tuyệt đối của thư mục package hiện tại có thể được biểu diễn bằng các biến trong thư mục truy xuất của tệp thư viện:

```go
// #cgo LDFLAGS: -L${SRCDIR}/libs -lfoo
```

Đoạn mã trên sẽ được mở rộng khi được liên kết:

```go
// #cgo LDFLAGS: -L/go/src/foo/libs -lfoo
```

Lệnh `#cgo` chủ yếu ảnh hưởng đến một số biến môi trường của trình biên dịch như CFLAGS, CPPFLAGS, CXXFLAGS, FFLAGS và LDFLAGS. LDFLAGS được sử dụng để đặt tham số của liên kết, ngoài một số biến được sử dụng để thay đổi tham số xây dựng của giai đoạn biên dịch (CFLAGS được sử dụng để đặt tham số biên dịch cho mã ngôn ngữ C).

Đối với người dùng sử dụng C và C++ trong môi trường CGO, có thể có ba tùy chọn biên dịch khác nhau: CFLAGS cho các tùy chọn biên dịch theo ngôn ngữ C, CXXFLAGS cho các tùy chọn biên dịch cụ thể C++ và CPPFLAGS cho các biên dịch C và C++. Tuy nhiên, trong giai đoạn liên kết, các tùy chọn liên kết C và C++ là chung, do đó không còn sự khác biệt giữa C và C++ tại thời điểm này và các target file của chúng cùng kiểu.

Các lệnh `#cgo` cũng hỗ trợ lựa chọn có điều kiện và các tùy chọn biên dịch hoặc liên kết tiếp theo có hiệu lực khi một hệ điều hành hoặc một kiểu kiến trúc CPU nhất định được đáp ứng. Ví dụ sau đây là các tùy chọn biên dịch và liên kết cho các nền tảng windows và non-windows:

```go
// #cgo windows CFLAGS: -DX86=1
// #cgo !windows LDFLAGS: -lm
```

Trên nền tảng windows, macro X86 được định nghĩa trước là 1 trước khi biên dịch, dưới nền tảng không phải là window, thư viện toán học được yêu cầu phải được liên kết trong pha liên kết. Việc sử dụng này hữu ích cho các tình huống trong đó chỉ có một vài khác biệt trong các tùy chọn biên dịch trên các nền tảng khác nhau.

Nếu CGO tương ứng với mã c khác nhau trong các hệ thống khác nhau, chúng tôi có thể sử dụng `#cgoinemony` để xác định các macro ngôn ngữ C khác nhau, sau đó sử dụng macro để phân biệt các mã khác nhau:
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376

```go
package main

/*
#cgo windows CFLAGS: -DCGO_OS_WINDOWS=1
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo linux CFLAGS: -DCGO_OS_LINUX=1

#if defined(CGO_OS_WINDOWS)
    const char* os = "windows";
#elif defined(CGO_OS_DARWIN)
<<<<<<< HEAD
    const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
    const char* os = "linux";
=======
    static const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
    static const char* os = "linux";
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376
#else
#    error(unknown os)
#endif
*/
import "C"

func main() {
    print(C.GoString(C.os))
}
```

<<<<<<< HEAD
Bằng cách này, chúng ta có thể biết được hệ thống mà code đang vận hành, nhờ đó áp dụng các kĩ thuật riêng cho các nền tảng khác nhau.

## 2.2.3 Biên dịch với tag

Build tag là một comment đặc biệt ở đầu file C/C++ trong môi trường Go hoặc CGO. Biên dịch có điều kiện tương tự như sử dụng macro `#cgo` để xác định các nền tảng khác nhau (ví dụ trên). Code được build sau khi macro của nền tảng tương ứng được xác định.

Ví dụ trình bày một cách khác khi các file nguồn sau sẽ chỉ được tạo khi debug build flag  được thiết lập:
=======
Bằng cách này, chúng ta có thể sử dụng các kỹ thuật thường được sử dụng trong C để xử lý mã nguồn khác biệt giữa các nền tảng khác nhau.

## 2.2.3 Biên dịch có điều kiện

Build tag là một comment đặc biệt ở đầu file C/C++ trong môi trường Go hoặc CGO. Biên dịch có điều kiện tương tự như `#cgomacro` được định nghĩa cho các nền tảng khác nhau. Mã tương ứng chỉ được build sau khi macro của nền tảng tương ứng được xác định. Tuy nhiên, #cgo có một hạn chế trong việc xác định các macro theo chỉ thị. Nó chỉ có thể dựa trên các hệ điều hành được hỗ trợ bởi Go, chẳng hạn như windows, darwin và linux. Nếu chúng ta muốn xác định một macro cho cờ DEBUG, các #cgo hướng dẫn sẽ bất lực. Có thể dễ dàng thực hiện tính năng biên dịch có điều kiện của build tag được cung cấp bởi ngôn ngữ Go.

Ví dụ các tệp nguồn sau sẽ chỉ được tạo khi cờ xây dựng gỡ lỗi được đặt:
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376

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

<<<<<<< HEAD
Chúng ta có thể dùng `-tags` chỉ định nhiều build flag cùng một lúc thông qua các đối số dòng lệnh.

Ví dụ các build flag sau chỉ ra rằng việc build chỉ được thực hiện trong kiến trúc "linux/386" hoặc "non-cgo environment" trong nền tảng darwin.
=======
Chúng ta có thể dùng `-tags` chỉ định nhiều cờ xây dựng cùng một lúc thông qua các đối số dòng lệnh, được phân tách bằng dấu cách.

Khi có nhiều build tag, chúng ta kết hợp nhiều cờ thông qua các quy tắc hoạt động hợp lý. Ví dụ, các cờ xây dựng sau chỉ ra rằng việc xây dựng chỉ được thực hiện trong "linux/386" hoặc "non-cgo environment" trong nền tảng darwin.
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376

```go
// +build linux,386 darwin,!cgo
```

<<<<<<< HEAD
Trong đó, Dấu phẩy (`,`) nghĩa là **và**. Khoảng trắng (`  `) nghĩa là **hoặc**.
=======
Trong đó, Dấu phẩy "," nghĩa là `và`. Khoản trắng nghĩa là `hoặc`.
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376
