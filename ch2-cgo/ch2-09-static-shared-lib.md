# 2.9 Thư viện tĩnh và động

Có ba cách để dùng tài nguyên `C/C++` trong `CGO`: dùng trực tiếp mã nguồn, liên kết tĩnh, liên kết động. Cách trực tiếp sử dụng mã nguồn là thêm dòng `import "C"` và phần comment mã nguồn  C phía trên, hoặc bao gộp mã nguồn `C/C++` trong package hiện tại. Cách dùng liên kết tĩnh và động thư viện cũng tương tự, bằng việc đặc tả thư viện liên kết trong cờ `LDFLAGS`. Trong phần này, chúng ta sẽ tập trung vào việc dùng thư viện tĩnh và động trong `CGO` như thế nào.

## 2.9.1 Dùng thư viện C tĩnh

Nếu tài nguyên C/C++ được đưa vào trong CGO có mã nguồn, và kích thước mã nguồn là nhỏ, cách dùng mã nguồn trực tiếp là một ý tưởng phổ biến nhất, nhưng nhiều khi chúng ta không có mã nguồn, hoặc quá trình xây dựng mã nguồn `C/C++` rất phức tạp. Đây là lúc mà việc dùng thư viện C tĩnh là lựa chọn tốt nhất. Bởi vì thư viện tĩnh là liên kết tĩnh, phần chương trình đích sẽ không sinh ra thêm phần dependencies trong khi chạy, hoặc không có một thư viện động cụ thể nào có thể đảm bảo kiểm soát lỗi phát sinh giữa các thành phần liên kết trong quá trình chạy. Tuy nhiên, thư viện tĩnh cũng có một yêu cầu nhất định trong pha liên kết: thư viện tĩnh thường chứa ba thành phần mã nguồn, nó có một lượng lớn các ký tự. Nếu một ký tự bị xung đột trong quá trình liên kết tĩnh, thì toàn bộ `link` sẽ hỏng.

Đầu tiên, chúng ta sẽ xây dựng một thư viện tĩnh đơn giản với ngôn ngữ C thuần. Thư viện tĩnh mà chúng ta xây dựng được gọi là `number`. Chỉ có một hàm `number_add_mod` trong thư viện dùng để lấy modulo của một tổng hai số cho một số thứ ba. Những files của thư viện `number` đặt trong cùng một thư mục.

Trong file `number/number.h` sẽ định nghĩa phần header chứa prototype của hàm

```c
int number_add_mod(int a, int b, int mod);
```

File `number/number.c` là phần hiện thực hàm như sau

```c
#include "number.h"

int number_add_mod(int a, int b, int mod) {
    return (a+b)%mod;
}
```

Bởi vì CGO dùng lệnh `GCC` để biên dịch và liên kết mã nguồn `C` và `Go` lại. Do đó, thư viện tĩnh cũng phải tương thích theo định dạng `GCC`.

Một thư viện tĩnh sẽ gọi `libnumber.a` có thể được sinh ra bằng lệnh sau

```
$ cd ./number
$ gcc -c -o number.o number.c
$ ar rcs libnumber.a number.o
```

Sau khi sinh ra thư viện tĩnh mang tên `libnumber.a`, chúng ta dùng nó trong `CGO`.

Tạo ra file  main.go như sau


```go
package main

//#cgo CFLAGS: -I./number
//#cgo LDFLAGS: -L${SRCDIR}/number -lnumber
//
//#include "number.h"
import "C"
import "fmt"

func main() {
    fmt.Println(C.number_add_mod(10, 5, 12))
}
```

Có hai lệnh `#cgo`, nó sẽ biên dịch và liên kết các tham số với nhau. Cờ `CFLAGS -I./number` sẽ thêm vào thư mục chứa các thư viện ứng với file header. Cờ `LDFLAGS` sẽ thể hiện liên kết tới thư viện tĩnh `libnumber.a` bằng cách thêm vào trường `-L${SRCDIR}/number`, nó sẽ đưa thư viện tĩnh `number` được biên dịch xong vào liên kết qua search path `-lnumber`. Nên chú ý rằng, phần search path trong liên kết không thể dùng trong các relative path (được giới hạn bởi mã nguồn `C/C++` linker). Chúng ta phải mở rộng thư mục hiện tại `${SRCDIR}` tương ứng với file mã nguồn đến một absolute path qua biến `cgo-specific` (cũng trên windows).  Absolute paths trong platform không thể chứa kí tự trống.

Bởi vì chúng ta đã có tất cả các mã nguồn cho thư viện `number`, chúng ta có thể dùng trình tạo mã để sinh ra thư viện tĩnh, hoặc dùng Makefiles để xây dựng thư viện tĩnh. Do đó, khi chúng ta publishing mã nguồn package CGO, chúng ta sẽ không cần phải biên dịch C static library trước.

Bởi vì sẽ có nhiều hơn một bước biên dịch thư viện tĩnh, Gói Go dùng để custom static library đã chứa tất cả các mã nguồn static library và không có thể cài đặt trực tiếp với `go get`. Tuy nhiên, chúng ta có vẫn có thể tải chúng xuống bằng `go get`, và dùng `go` để sinh ra điểm gọi thư viện tĩnh được xây dựng, và cuối cùng chúng ta sẽ `go install` để hoàn thành việc cài đặt.

Để hỗ trợ lệnh `go get` cho việc download và install một cách trực tiếp, ngôn ngữ C của chúng ta sẽ có cú pháp `#include` dùng để liên kết file mã nguồn của thư viện number đến gói hiện tại.

Tạo ra file `z_link_number_c.c` như sau

```c
#include "./number/number.c"
```

Sau đó thực thi lệnh `go get` hoặc `go build`, `CGO` sẽ tự động biên dịch mã nguồn ứng với thư viện `number`. Kĩ thuật này sẽ chuyển đổi thư viện tĩnh thành mã nguồn để references mà không cần thay đổi kết cấu tổ chức của mã nguồn thư viện tĩnh. Gói `CGO` thật hoàn hảo.


Nếu chúng ta sử dụng thư viện tĩnh từ bên thứ ba, chúng ta cần phải tải chúng và cài đặt thư viện tĩnh đến một nơi phù hợp. Sau đó đặc tả location của header files và libraries qua cờ `CFLAGS` và `LDFLAGS` trong lệnh `#cgo`. Trong các hệ điều hành khác nhau, hoặc các phiên bản khác nhau của cùng hệ điều hành, installation paths của những thư viện có thể khác nhau, do đó làm cách nào để có thể biết được các thay đổi trong mã nguồn?

Trong môi trường Linux, có một lệnh `pkg-config` được dùng để truy vấn các tham số `compile` và `link` khi dùng các thư viện động/tĩnh. Chúng ta có thể dùng lệnh `pkg-config` trực tiếp trong lệnh `#cgo` để generate compilation và linking parameters. Bạn có thể customize lệnh `pkg-config` với biến môi trường `PKG_CONFIG`. Bởi vì các hệ điều hành khác nhau có thể hỗ trợ lệnh `pkg-config` theo cách khác nhau, thật khó để làm tương thích các build parameters cho các hệ điều hành khác nhau. Tuy nhiên, trong hệ điều hành cụ thể là Linux, lệnh `pkg-config` chỉ đơn giản quản lý các build parameters. Chi tiết của việc dùng `pkg-config` không được nói ở đây, do đó bạn có thể thấy chúng trong các tài liệu liên quan khác.

## 2.9.2 Sử dụng thư viện C động

Ý định ban đầu của thư viện động là chia sẻ cùng một thư viện, các processes khác nhau có thể chia sẻ trên cùng một memory hoặc disk resources. Nhưng đĩa cứng và RAM hiện nay có giá rẻ, hai vai trò sẽ trở nên không đáng quan tâm, do đó đâu là giá trị của thư viện động ở đây? Từ góc nhìn của việc phát triển thư viện, thư viện động có thể tách biệt nhau và giảm thiểu rủi ro của việc xung đột trong khi linking. Và với những platform như windows, thư viện động là một cách khả thi để mở rộng các nền tảng biên dịch như `VC` và `GCC`.

Trong `CGO`, việc dùng thư viện động và tĩnh là như nhau, bởi vì thư viện động sẽ phải có một static export library nhỏ dùng cho việc liên kết (Linux có thể trực tiếp liên kết các files, nhưng cũng tạo ra dll bên dưới file `.a` dùng cho liên kết). Chúng ta có thể dùng thư viện `number` ở phần trước như là một ví dụ minh họa cho việc dùng thư viện động.

Trong môi trường `gcc` dưới macOS hoặc Linux, chúng ta có thể sinh ra thư viện động của một số thư viện với những lệnh sau: 

```
$ cd number
$ gcc -shared -o libnumber.so number.c
```

Bởi vì, base names của thư viện động và tĩnh là `libnumber`, chỉ phần hậu tố là sẽ khác. Do đó trong mã nguồn của ngôn ngữ Go sẽ giống chính xác với phiên bản thư viện tĩnh.

```go
package main

//#cgo CFLAGS: -I./number
//#cgo LDFLAGS: -L${SRCDIR}/number -lnumber
//
//#include "number.h"
import "C"
import "fmt"

func main() {
    fmt.Println(C.number_add_mod(10, 5, 12))
}
```

`CGO` sẽ tự động tìm `libnumber.a` hoặc `libnumber.so` trong ở bước linking trong thời gian biên dịch.


Với windows platform, chúng ta có thể dùng công cụ `VC` để sinh ra thư viện động (sẽ có một số thư viện Windows phức tạp chỉ có thể được build với `VC`). Đầu tiên, chúng ta phải tạo một file định nghĩa cho `number.dll` để quản lý các kí hiệu dùng để exported thư viện động.

Nội dung của file `number.def` như sau:

```
LIBRARY number.dll

EXPORTS
number_add_mod
```

Dòng đầu tiên `LIBRARY` sẽ chỉ ra tên của file và tên của thư viện động, và sau đó là mệnh đề  `EXPORTS` theo sau bởi một danh sách các tên dùng để exported.

Giờ đây, chúng ta có thể dùng những lệnh sau để tạo ra thư viện động (cần dùng `VC` tools)

```
$ cl /c number.c
$ link /DLL /OUT:number.dll number.obj number.def
```

Vào lúc này, một export library `number.lib` sẽ sinh ra `dll` cùng lúc. Nhưng trong `CGO`, chúng ta không thể dùng `link library` trong định dạng `lib`.

Để sinh ra định dạng `.a` cho việc export library cần dùng `mingw Toolbox dlltool command`:

```
$ dlltool -dllname number.dll --def number.def --output-lib libnumber.a
```

Một khi `libnumber.a` được sinh ra, có thể  dùng `-lnumber` thông qua các  link parameters.

Nên chú ý rằng, tại thời điểm thực thi, thư viện động cần được đặt ở cùng nơi để system có thể thấy. Trên windows, bạn có thể đặt `dynamic library` và `executable program` trong cùng một thư mục, hoặc thêm một absolute path của directory trong khi dynamic library được đưa vào biến môi trường `PATH`. Trong macOS, bận cần phải thiết lập biến môi trường `DYLD_LIBRARY_PATH`. Trong hệ thống Linux, bạn cần thiết lập biến `LD_LIBRARY_PATH`.

## 2.9.3  Exporting C Static Libraries

`CGO` không chỉ được dùng trong thư viện C tĩnh, các export functions được hiện thực bởi Go hoặc  C static libraries. Chúng ta có thể dùng Go để hiện thực modulo addition function như phần trước như sau

Tạo `number.go` với nội dung như sau

```go
package main

import "C"

func main() {}

//export number_add_mod
func number_add_mod(a, b, mod C.int) C.int {
    return (a + b) % mod
}
```

Theo như mô tả của tài liệu `CGO`, chúng ta cần export C function trong main package. Với cách xây dựng thư viện C tĩnh, main function trong main package được phớt lờ, và C function sẽ đơn giản được exported. Xây dựng các lệnh sau

```
$ go build -buildmode=c-archive -o number.a
```

Khi sinh ra thư viện tĩnh `number.a`, cgo cũng sẽ sinh ra file `number.h`

Nội dung của `number.h` sẽ như sau, (để dễ hiển thị, nội dung sẽ được sắp xếp hợp lý)


```c
#ifdef __cplusplus
extern "C" {
#endif

extern int number_add_mod(int p0, int p1, int p2);

#ifdef __cplusplus
}
#endif
```

Khi phần ngữ pháp `extern "C"` được thêm vào để  dùng đồng thời trên cả ngôn ngữ `C` và `C++`. Nội dung của phần lõi sẽ định nghĩa hàm `number_add_mod` để được exported.

Sau đó chúng ta tạo ra file `_test_main.c` để kiểm tra việc sinh ra C static library (bên dưới là phần prefix dùng cho việc xây dựng C static library để bỏ qua file đó)

```c
#include "number.h"

#include <stdio.h>

int main() {
    int a = 10;
    int b = 5;
    int c = 12;

    int x = number_add_mod(a, b, c);
    printf("(%d+%d)%%%d = %d\n", a, b, c, x);

    return 0;
}
```

Biên dịch và chạy chúng với những lệnh sau:


```c
$ gcc -o a.out _test_main.c number.a
$ ./a.out
```

Quá trình sinh ra static library dùng `CGO` thự sự đơn giản.

## 2.9.4 Exporting C Dynamic Library

Quá trình exporting một thư viện động bằng `CGO` sẽ tương tự như một thư viện tĩnh, ngoại trừ build mode sẽ thay đổi `c-shared` và  output file name được đổi thành `number.so` 

```
$ go build -buildmode=c-shared -o number.so
```

Nội dung của file `_test_main.c` sẽ không thay đổi, sau đó biên dịch và chạy chúng với các lệnh sau

```
$ gcc -o a.out _test_main.c number.so
$ ./a.out
```

## 2.9.5 Exporting Functions of Non-Main Packages

Lệnh `go help buildmode` có thể được dùng để xem cấu trúc lệnh của thư viện C tĩnh và thư viện C động.


```
-buildmode=c-archive
    Build the listed main package, plus all packages it imports,
    into a C archive file. The only callable symbols will be those
    functions exported using a cgo //export comment. Requires
    exactly one main package to be listed.

-buildmode=c-shared
    Build the listed main package, plus all packages it imports,
    into a C shared library. The only callable symbols will
    be those functions exported using a cgo //export comment.
    Requires exactly one main package to be listed.

```

Phần tài liệu đã nói rằng `exported C function`  phải được `exported` trong main package trước khi sinh ra `header file` chứa những declared statement. Nhưng nhiều khi chúng ta phải đề cập đến việc tổ chức các kiểu khác nhau cho việc export functions đến các Go packages và sau đó export chúng như là thư viện động/tĩnh.

Để hiện thực hàm C từ một package khác `main` package, hoặc để export C function từ nhiều packages (bởi vì chỉ có thể có một main package), chúng ta cần cung cấp header file ứng với hàm C được export ( bởi vì `CGO` không thể là một non-main package) Export một function để sinh ra header file.

Hỗ trợ cho viện tạo ra một số subpackage chúng là một hàm modular addition function


```go
package number

import "C"

//export number_add_mod
func number_add_mod(a, b, mod C.int) C.int {
    return (a + b) % mod
}
```

Sau đó tạo ra một main package

```go
package main

import "C"

import (
    "fmt"

    _ "./number"
)

func main() {
    println("Done")
}

//export goPrintln
func goPrintln(s *C.char) {
    fmt.Println("goPrintln:", C.GoString(s))
}

```

[>>  mã nguồn](../examples/ch2/ch2.9/5-modular-func/main.go)

Trong số đó, chúng ta phải import một số sub-package, có một exported C function `number_add_mod` trong `number sub-package`, và chúng ta cũng phải export hàm `goPrintln` trong main package.

Tạo một `C static library` với lệnh sau

```
$ go build -buildmode=c-archive -o main.a
```

Giờ đây, trong khi sinh ra thư viện tĩnh `main.a`, một `main.h` header file cũng được sinh ra. Tuy nhiên, header file `main.h` chỉ việc định nghĩa  hàm `goPrintln` từ main package, và không có bất kì định nghĩa nào của `number subpackage export function`. Thực tế, hàm `number_add_mod` tồn tại trong khi sinh thư viện C tĩnh, chúng ta có thể dùng chúng một cách trực tiếp.

Tạo ra một `_test_main.c`  file sẽ theo như sau

```c
#include <stdio.h>

void goPrintln(char*);
int number_add_mod(int a, int b, int mod);

int main() {
    int a = 10;
    int b = 5;
    int c = 12;

    int x = number_add_mod(a, b, c);
    printf("(%d+%d)%%%d = %d\n", a, b, c, x);

    goPrintln("done");
    return 0;
}
```

Chúng sẽ không bao gộp header file `main.h` được tự động sinh ra bởi `CGO`, nhưng chúng ta có thể định nghĩa thủ công hai `export functions` goPrintln và number_add_mod. Cách này làm chúng ta sẽ phải hiện thực một export C functions từ nhiều Go package.

