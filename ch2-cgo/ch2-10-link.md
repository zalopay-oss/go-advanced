# 2.10  Biên dịch và liên kết các tham số

Biên dịch và liên kết các parameters thường là một vấn đề mà các lập trình viên `C/C++` phải đối mặt. Xây dựng một ứng dụng `C/C++` yêu cầu hai bước biên dịch và liên kết cũng như trong `CGO`. Trong phần này, chúng tôi sẽ giới thiệu ngắn các bước biên dịch và link parameters thường được dùng trong `CGO`.

## 2.10.1 Biên dịch Parameters : `CFLAGS/CPPFLAGS/CXXFLAGS`

**Compilation parameters** chủ yếu sẽ duyệt qua đường dẫn chứa header file, các mascros, và các parameters khác. Theo lý thuyết, C và C++ hoàn toàn là hai ngôn ngữ lập trình độc lập, và chúng có các tham số biên dịch riêng. Nhưng bởi vì ngôn ngữ C++ tương thích rất chặt đối với ngôn ngữ C, C++ có thể được hiểu là một tập hợp cha của ngôn ngữ C, do đó C và C++ sẽ chia sẻ một lượng lớn các compilation parameters. Do đó, CGO sẽ cùng cấp ba tham số `CFLAGS/CPPFLAGS/CXXFLAGS`, trong đó `CFLAGS` sẽ ứng với việc biên dịch ngôn ngữ C (`.c`), `CPPFLAGS` sẽ ứng với C/C++ (`.c`, `.cc`, `.cpp`, `.cxx`), và `CXXFLAGS` ứng với C++ thuần (`.cc`, `.cpp`, `*.cxx`).

## 2.10.2 Link parameters: `LDFLAGS`

**Link parameters** sẽ chủ yếu chứa **search directory** của thư viện cần liên kết và tên của thư viện được liên kết. Link library sẽ không hỗ trợ relative paths, chúng ta phải đặc tả một absolute path cho link library. **${SRCDIR}** trong `cgo` là một `absolute path` trong thư mục hiện tại. Các file định dạng của đối tượng C và C++ sau khi được biên dịch là như nhau, do đó `LDFLAGS` sẽ ứng với các `C/C++` link parameters phổ biến.

## 2.10.3 `pkg-config`

Bằng việc hỗ trợ biên dịch và link parameters cho các thư viện C/C++ khác nhau là một công việc kinh khủng, do đó `cgo` đã hỗ trợ công cụ `pkg-config`. Chúng ta có thể dùng `#cgo pkg-config xxx` để sinh ra lệnh compile và link parameter cho thư viện `xxx` yêu cầu, bằng việc gọi `pkg-config xxx --cflags` sẽ sinh ra `compiler parameters`, dùng `pkg-config xxx --libs` để sinh ra  links parameters. Có thể thấy  linked parameter được sinh ra bởi công cụ `pkg-config` là rất bình thường với `C/C++` và không thể  phân biệt nhiều.

Mặc dù công cụ `pkg-config` là rất thuận tiện, vẫn có những thư viện `C/C++` non-standard mà chúng không hỗ trợ. Khi đó, chúng ta có thể thủ công hiện thực compilation và link parameter.

Ví dụ, có một thư viện `C/C++` gọi là `xxx`, chúng ta có thể tạo ra `/usr/local/lib/pkgconfig/xxx.bc` một cách thủ công.

```
Name: xxx
Cflags:-I/usr/local/include
Libs:-L/usr/local/lib –lxxx2
```

Trong đó trường `Name` là tên của thư viện, và dòng `Cflags` và `Libs` sẽ ứng với compilation và link parameters được yêu cầu cho thư viện `xxx`. Nếu file `bc` nằm trong một đường dẫn khác, công cụ  search directory `pkg-config` sẽ tìm thấy chúng thông qua biến môi trường `PKG_CONFIG_PATH`.

Trong `cgo`, chúng ta cũng có thể viết một chương trình custom  `pkg-config` qua biến `PKG_CONFIG`. Nếu bạn tự hiện thực một chương trình CGO cụ thể `pkg-config`, chỉ cần hiện thực hai tham số `--cflags` và `--libs`.

Chương trình bên dưới sẽ xây dựng một link parameters để sinh ra `Python3` dưới macos system:

```go
// py3-config.go
func main() {
    for _, s := range os.Args {
        if s == "--cflags" {
            out, _ := exec.Command("python3-config", "--cflags").CombinedOutput()
            out = bytes.Replace(out, []byte("-arch"), []byte{}, -1)
            out = bytes.Replace(out, []byte("i386"), []byte{}, -1)
            out = bytes.Replace(out, []byte("x86_64"), []byte{}, -1)
            fmt.Print(string(out))
            return
        }
        if s == "--libs" {
            out, _ := exec.Command("python3-config", "--ldflags").CombinedOutput()
            fmt.Print(string(out))
            return
        }
    }
}
```

Sau đó build và dùng công cụ `pkg-config` đã được custom với các lệnh sau

```
$ go build -o py3-config py3-config.go
$ PKG_CONFIG=./py3-config go build -buildmode=c-shared -o gopkg.so main.go
```

## 2.10.4 `go get chain`

Lệnh `go get` package sẽ liên kết các package phụ thuộc liên quan. Một chuỗi các package phụ thuộc sau như A phụ thuộc B, B phụ thuộc C, C phụ thuộc D `pkgA -> pkgB -> pkgC -> pkgD -> ...`. Sau khi `go get A_package`, chúng sẽ get B,C,D package. Nếu việc build hỏng sau khi get package_B, nó sẽ dẫn tới việc hỏng toàn bộ chuỗi , kết quả là lệnh get package_A bị hỏng.

Sẽ có một số lý do cho việc hỏng chuỗi get, đây là một số lý do phổ biến :

* Một số hệ thống không hỗ trợ, việc biên dịch sẽ hỏng
* Phụ thuộc vào `cgo`, user không có `gcc` được cài đặt sẵn.
* Phụ thuộc vào `cgo`, những thư viện độc lập không được cài đặt.
* Phụ thuộc vào `pkg-config`, không được cài đặt trên windows
* Phụ thuộc vào `pkg-config`, không có file `bc` được tìm thấy tương ứng.
* Tin cậy vào `custom pkg-config`, yêu cầu một số thiết lập thêm
* Phụ thuộc vào `swig`, user chưa cài đặt `swig`, hoặc phiên bản không tương thích.

Việc phân tích kĩ càng sẽ cho thấy vấn đề nảy sinh với `CGO` đa số sẽ dẫn đến hư hỏng. Đây không phải là một hiện tượng tình cờ. Việc tự động xây dựng mã nguồn `C/C++` luôn là vấn đề của thế giới. Hiện nay, không có một chuẩn thống nhất việc quản lý công cụ C/C++ để mọi người có thể nhận diện được.

Bằng việc dùng `cgo`, `gcc` và các công cụ build khác phải được cài đặt, và cố gắng để hỗ trợ mainstream system.

Nếu một gói C/C++ độc lập là nhỏ và mã nguồn luôn có sẵn, người ta thường thích tự xây dựng hơn.

Ví dụ, gói thư viện `github.com/chai2010/webp` hiện thực `zero configuration dependencies` bằng việc tạo ra một key files cho mỗi file mã nguồn C/C++ trong gói hiện tại.

```
// z_libwebp_src_dec_alpha.c
#include "./internal/libwebp/src/dec/alpha.c"
```

Do đó `z_libwebp_src_dec_alpha.c`, một mã nguồn `libweb native` có thể được biên dịch khi file vừa biên dịch xong. Các dependencies là các relative directories, tính nhất quán tối đa có thể được duy trì trong một số platform khác nhau.

## 2.10.5 Export các hàm C trong nhiều gói non-main packages

Một tài liệu chính thức đã viết rằng, exported Go function nên có một main package, nhưng sự thật là những Go export function của các packages khác cũng hợp lệ. Bởi vì exported Go function có thể được dùng như là một C function, nó có thể hợp lệ. Nhưng các Go functions được exported bởi các packages khác nhau có thể có chung một global namespace, do đó chúng ta cần phải cẩn trọng trong việc tránh trùng tên. Nếu bạn export một hàm Go từ các package khác nhau vào một không gian C language, sau đó `cgo` sẽ tự động sinh ra file `_cgo_export.h` và không chứa tất cả các lời định nghĩa hàm ở mọi nơi, chúng ta sẽ phải export tất cả các hàm thông qua cách tự viết header file.

