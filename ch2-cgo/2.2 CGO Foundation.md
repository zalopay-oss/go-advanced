# 2.2 CGO Foundation

Để sử dụng tính năng CGO, bạn cần cài đặt công cụ C/C++. Theo macOS và Linux, bạn cần cài đặt `GCC`. Trong Windows, bạn cần cài đặt công cụ `MinGW`. Đồng thời, bạn cần đảm bảo rằng biến môi trường `CGO_ENABLED` được đặt thành 1.

## 2.2.1 Lệnh `import "C"`

Nếu lệnh import `import "C"` xuất hiện trong code, nó có nghĩa là tính năng CGO được sử dụng. Comment trước lệnh đó là cú pháp để Go nhận ra code của C. Khi CGO được bật, bạn có thể thêm các file code của C/C++ tương ứng trong thư mục hiện tại.

Ví dụ đơn giản nhất:

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
    C.printint(C.int(v))
}
```

[>> mã nguồn](../examples/ch2/ch2.2/1-simplest-go/main.go)

Ví dụ này cho thấy việc sử dụng CGO cơ bản. Phần đầu của comment khai báo hàm C sẽ được gọi và file header được liên kết. Tất cả các thành phần ngôn ngữ C trong file header sau khi được đưa vào sẽ được thêm vào gói "C" ảo. `Cần lưu ý rằng câu lệnh import "C" yêu cầu một dòng riêng và không thể được import cùng với các gói khác`. Truyền tham số cho hàm C cũng rất đơn giản và nó có thể được chuyển đổi trực tiếp thành một loại ngôn ngữ C tương ứng. Trong ví dụ trên, `C.int(v)` được sử dụng để chuyển đổi giá trị kiểu int trong Go đến giá trị kiểu int trong ngôn ngữ C, sau đó gọi hàm printint được xác định bằng ngôn ngữ C để in.

Cần lưu ý rằng Go là loại ngôn ngữ mạnh (chặt chẽ), do đó tham số được truyền phải chính xác với khai báo, và phải được chuyển đổi sang kiểu trong C bằng các hàm chuyển đổi trước khi truyền, không thể truyền trực tiếp bằng kiểu của Go. Đồng thời, các ký hiệu của C được import thông qua gói C thì `không cần phải viết hoa`, không cần phải tuân theo quy tắc của Go.

Cgo đặt các ký hiệu ngôn ngữ C được tham chiếu bởi gói hiện tại vào gói C ảo. Đồng thời, các gói Go khác mà gói hiện tại phụ thuộc cũng có thể giới thiệu các gói C ảo tương tự thông qua cgo, nhưng các gói Go khác nhau giới thiệu các gói ảo. Các kiểu giữa các gói C không phải là toản thể. Ràng buộc này có thể có một tác động nhỏ đến khả năng tự xây dựng một số chức năng cgo.

Ví dụ chúng tôi muốn định nghĩa kiểu `CChar` tương ứng với con trỏ char của C trong Go và sau đó thêm phương thức GoString để trả về chuỗi ngôn ngữ Go

```go
package cgo_helper

//#include <stdio.h>
import "C"

type CChar C.char

func (p *CChar) GoString() string {
    return C.GoString((*C.char)(p))
}

func PrintCString(cs *C.char) {
    C.puts(cs)
}
```

[>> mã nguồn](../examples/ch2/ch2.2/1-cchar/cgo_helper/cgo_helper.go)

Bây giờ có thể ta muốn sử dụng hàm này trong các package Go khác:

```go
package main

//static const char* cs = "hello";
import "C"
import "./cgo_helper"

func main() {
    cgo_helper.PrintCString(C.cs)
}
```

[>> mã nguồn](../examples/ch2/ch2.2/1-cchar/main/main.go)

Nhưng đoạn code này sẽ không chạy được. Vì biến `C.cs` được đề cập trong gói main hiện tại là kiểu của gói ảo C được xây dựng trên `*char` (*C.char, chính xác hơn là *main.C.char), còn kiểu `*C.type` được đề cập đến trong gói `cgo_helper` (`*cgo_helper.C.char`) là khác nhau. Trong ngôn ngữ Go, các phương thức phụ thuộc vào kiểu. Các kiểu được gói C ảo được đề cập trong các gói Go khác nhau là khác nhau (`main.C` không giống `cgo_helper.C`) chính là nguyên nhân khiến các kiểu Go được mở rộng từ chúng thành các kiểu khác nhau (`*main.C.char` khác `*cgo_helper.C.char`). Điều này cuối cùng đã khiến đoạn code đó không hoạt động được.

Người dùng có kinh nghiệm với ngôn ngữ Go có thể đề xuất rằng các tham số được truyền vào sau khi chuyển đổi. Nhưng phương pháp này dường như không khả thi, bởi vì các tham số của `cgo_helper.PrintCString` là kiểu `*C.char` được đề cập bởi gói riêng của nó và nó không thể truy cập trực tiếp từ bên ngoài. Nói cách khác, nếu một gói trực tiếp sử dụng loại C ảo tương tự `*C.char`  trong một interface chung, các gói Go khác không thể sử dụng trực tiếp các loại này trừ khi gói Go cũng cung cấp hàm tạo `*C.chartype`. Do nhiều yếu tố này, nếu bạn muốn kiểm tra các kiểu được export trực tiếp bởi cgo trong môi trường thử nghiệm đi, sẽ có những hạn chế tương tự.

## 2.2.2 Lệnh `#cgo`

Trong dòng ghi chú `import "C"` phía trước các lệnh. Bạn có thể đặt các tham số cho quá trình biên dịch và quá trình liên kết thông qua các lệnh `#cgo`. Các tham số của quá trình biên dịch chủ yếu được sử dụng để xác định các macro liên quan và đường dẫn truy xuất file header đã chỉ định. Các tham số của quá trình liên kết chủ yếu là để xác định đường dẫn truy xuất file thư viện và file thư viện sẽ được liên kết.

```go
// #cgo CFLAGS: -DPNG_DEBUG=1 -I./include   // Định nghĩa macro PNG_DEBUG, giá trị là 1
// #cgo LDFLAGS: -L/usr/local/lib -lpng
// #include <png.h>
import "C"
```

Trong đoạn mã trên, phần CFLAGS, -Dpart định nghĩa macro PNG_DEBUG, giá trị là 1; -I xác định thư mục tìm kiếm có trong file header. Trong phần LDFLAGS, -L thư mục truy xuất file thư viện được -l chỉ định khi liên kết và thư viện png liên kết là bắt buộc khi liên kết được chỉ định.

Do các vấn đề mà C/C ++ để lại, thư mục truy xuất file header C có thể là một thư mục tương đối, nhưng thư mục truy xuất tệp thư viện yêu cầu một đường dẫn tuyệt đối. ${SRCDIR} Đường dẫn tuyệt đối của thư mục gói hiện tại có thể được biểu diễn bằng các biến trong thư mục truy xuất của tệp thư viện:

```go
// #cgo LDFLAGS: -L${SRCDIR}/libs -lfoo
```

Đoạn mã trên sẽ được mở rộng khi được liên kết:

```go
// #cgo LDFLAGS: -L/go/src/foo/libs -lfoo
```

Lệnh `#cgo` tuyên bố chủ yếu ảnh hưởng đến một số biến môi trường của trình biên dịch như CFLAGS, CPPFLAGS, CXXFLAGS, FFLAGS và LDFLAGS. LDFLAGS được sử dụng để đặt tham số của liên kết, ngoài một số biến được sử dụng để thay đổi tham số xây dựng của giai đoạn biên dịch (CFLAGS được sử dụng để đặt tham số biên dịch cho mã ngôn ngữ C).

Đối với người dùng sử dụng C và C++ trong môi trường cgo, có thể có ba tùy chọn biên dịch khác nhau: CFLAGS cho các tùy chọn biên dịch theo ngôn ngữ C, CXXFLAGS cho các tùy chọn biên dịch cụ thể C++ và CPPFLAGS cho các biên dịch C và C++. Tuy nhiên, trong giai đoạn liên kết, các tùy chọn liên kết C và C++ là chung, do đó không còn sự khác biệt giữa C và C++ tại thời điểm này và các target file của chúng cùng loại.

Các lệnh `#cgo` cũng hỗ trợ lựa chọn có điều kiện và các tùy chọn biên dịch hoặc liên kết tiếp theo có hiệu lực khi một hệ điều hành hoặc một kiểu kiến trúc CPU nhất định được đáp ứng. Ví dụ sau đây là các tùy chọn biên dịch và liên kết cho các nền tảng windows và non-windows:

```go
// #cgo windows CFLAGS: -DX86=1
// #cgo !windows LDFLAGS: -lm
```

Trong nền tảng windows, macro X86 được định nghĩa trước là 1 trước khi biên dịch, dưới nền tảng không phải là window, thư viện toán học được yêu cầu phải được liên kết trong pha liên kết. Việc sử dụng này hữu ích cho các tình huống trong đó chỉ có một vài khác biệt trong các tùy chọn biên dịch trên các nền tảng khác nhau.

Nếu cgo tương ứng với mã c khác nhau trong các hệ thống khác nhau, chúng tôi có thể sử dụng `#cgoinemony` để xác định các macro ngôn ngữ C khác nhau, sau đó sử dụng macro để phân biệt các mã khác nhau:

```go
package main

/*
#cgo windows CFLAGS: -DCGO_OS_WINDOWS=1
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo linux CFLAGS: -DCGO_OS_LINUX=1

#if defined(CGO_OS_WINDOWS)
    const char* os = "windows";
#elif defined(CGO_OS_DARWIN)
    static const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
    static const char* os = "linux";
#else
#    error(unknown os)
#endif
*/
import "C"

func main() {
    print(C.GoString(C.os))
}
```

[>> mã nguồn](../examples/ch2/ch2.2/2-cgo-statement/main.go)

Bằng cách này, chúng ta có thể sử dụng các kỹ thuật thường được sử dụng trong C để xử lý mã nguồn khác biệt giữa các nền tảng khác nhau.

## 2.2.3 Biên dịch có điều kiện

Build tag là một comment đặc biệt ở đầu file C/C++ trong môi trường Go hoặc cgo. Biên dịch có điều kiện tương tự như `#cgomacro` được định nghĩa cho các nền tảng khác nhau. Mã tương ứng chỉ được build sau khi macro của nền tảng tương ứng được xác định. Tuy nhiên, #cgo có một hạn chế trong việc xác định các macro theo chỉ thị. Nó chỉ có thể dựa trên các hệ điều hành được hỗ trợ bởi Go, chẳng hạn như windows, darwin và linux. Nếu chúng ta muốn xác định một macro cho cờ DEBUG, các #cgo hướng dẫn sẽ bất lực. Có thể dễ dàng thực hiện tính năng biên dịch có điều kiện của build tag được cung cấp bởi ngôn ngữ Go.

Ví dụ các tệp nguồn sau sẽ chỉ được tạo khi cờ xây dựng gỡ lỗi được đặt:

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

Chúng ta có thể dùng `-tags` chỉ định nhiều cờ xây dựng cùng một lúc thông qua các đối số dòng lệnh, được phân tách bằng dấu cách.

Khi có nhiều build tag, chúng ta kết hợp nhiều cờ thông qua các quy tắc hoạt động hợp lý. Ví dụ, các cờ xây dựng sau chỉ ra rằng việc xây dựng chỉ được thực hiện trong "linux/386" hoặc "non-cgo environment" trong nền tảng darwin.

```go
// +build linux,386 darwin,!cgo
```

Trong đó, Dấu phẩy "," nghĩa là `và`. Khoản trắng nghĩa là `hoặc`.
