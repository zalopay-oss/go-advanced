# 2.9. Thư viện tĩnh và động

Có ba cách để sử dụng mã nguồn C/C++ trong **CGO**:
  1. Dùng trực tiếp mã nguồn (thêm dòng `import "C"` và chú thích mã nguồn C phía trên, hoặc bao gộp mã nguồn `C/C++` trong package hiện tại).
  2. Liên kết tĩnh mã nguồn (khai báo thư viện liên kết trong cờ `LDFLAGS`).
  3. Liên kết động mã nguồn.

<div align="center">
	<img src="../images/ch2-9-static-dynamic-lib.gif" width="400">
</div>

Chi tiết về sự khác biệt giữa thư viện tĩnh và động bạn đọc có thể xem thêm tại đây: [What-is-the-difference-between-static-and-dynamic-linking](https://www.quora.com/What-is-the-difference-between-static-and-dynamic-linking).

Sau đây chúng ta sẽ đi vào cách dùng thư viện tĩnh và thư viện động trong CGO.

## 2.9.1. Dùng thư viện C tĩnh

Nếu mã nguồn C/C++ được dùng trong CGO có kích thước nhỏ thì cách đưa trực tiếp chúng vào chương trình là một ý tưởng phổ biến nhất, nhưng nhiều lúc chúng ta không tự xây dựng mã nguồn, hoặc quá trình xây dựng mã nguồn C/C++ rất phức tạp thì đây là lúc thư viện C tĩnh phát huy thế mạnh của mình.

Ở ví dụ đầu tiên, chúng ta sẽ xây dựng một thư viện tĩnh đơn giản được gọi là `number`, chỉ có một hàm `number_add_mod` trong thư viện dùng để lấy modulo của một tổng hai số cho một số thứ ba, những files của thư viện `number` đặt trong cùng một thư mục:

***number/number.h*** : header chứa prototype của hàm:

```c
int number_add_mod(int a, int b, int mod);
```

***number/number.c***: phần hiện thực hàm như sau:

```c
#include "number.h"

int number_add_mod(int a, int b, int mod) {
    return (a+b)%mod;
}
```

Bởi vì CGO dùng lệnh GCC để biên dịch và liên kết mã nguồn C và Go lại, do đó thư viện tĩnh cần phải tương thích với GCC compiler.

Thư viện tĩnh `libnumber.a` có thể được sinh ra bằng lệnh sau:

```sh
// di chuyển tới thư mục mã nguồn
$ cd ./number
// biên dịch ra file object từ file mã nguồn
$ gcc -c -o number.o number.c
// lệnh tạo ra thư viện tĩnh libnumber.a từ file object
// chi tiết về lệnh ar có thể xem tại https://linux.die.net/man/1/ar
$ ar rcs libnumber.a number.o
```

Sau khi sinh ra thư viện tĩnh mang tên `libnumber.a`, chúng ta dùng nó trong CGO.

***main.go*** được tạo ra như sau:

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

Hai lệnh `#cgo` trên dùng để biên dịch và liên kết mã nguồn với nhau:
  * Cờ `CFLAGS -I./number` : khai báo đường dẫn đến thư mục mã nguồn.
  * Cờ `LDFLAGS: -L${SRCDIR}/number -lnumber` : khai báo đường dẫn đến thư viện tĩnh `libnumber.a` với search path `-lnumber`.

Chú ý rằng: đường dẫn trong liên kết không thể dùng [relative path](https://support.dtsearch.com/webhelp/dtsearch/relative_paths.htm) mà phải dùng một [absolute path](http://www.linfo.org/absolute_pathname.html), ngoài ra đường dẫn không được chứa bất kỳ khoảng trắng nào.

Ví dụ : `LDFLAGS: -L/home/mypc/number -lnumber`

Kết quả như sau:

```sh
$ go run main.go
3
```

Nếu chúng ta sử dụng thư viện tĩnh từ bên thứ ba, chúng ta cần phải tải chúng và cài đặt thư viện tĩnh đến một nơi phù hợp, sau đó đặc tả location của header files và libraries qua cờ `CFLAGS` và `LDFLAGS` trong lệnh `#cgo`.

Trong môi trường Linux, có một lệnh [pkg-config](https://linux.die.net/man/1/pkg-config) được dùng để truy vấn các tham số compile và link khi dùng các thư viện tĩnh, chúng ta có thể dùng lệnh pkg-config trực tiếp trong lệnh `#cgo` để generate compilation và linking parameters, bạn có thể customize lệnh pkg-config với biến môi trường `PKG_CONFIG`.

## 2.9.2. Sử dụng thư viện C động

Ý tưởng của thư viện động là shared library, các process khác nhau có thể chia sẻ trên cùng một tài nguyên bộ nhớ trên RAM hoặc đĩa cứng, nhưng hiện nay giá thành đĩa cứng và RAM cũng tương đối rẻ, nên hai vai trò sẽ trở nên không đáng quan tâm, do đó đâu là giá trị của thư viện động ở đây?

Từ góc nhìn của việc phát triển thư viện, thư viện động có thể tách biệt nhau và giảm thiểu rủi ro của việc xung đột trong khi liên kết, với những nền tảng như Windows, thư viện động là một cách khả thi để mở rộng các nền tảng biên dịch như `gcc`.

Trong môi trường `gcc` dưới MacOS hoặc Linux, chúng ta có thể sinh ra thư viện động của một số thư viện với những lệnh sau:

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

`CGO` sẽ tự động tìm `libnumber.a` hoặc `libnumber.so` ở bước liên kết trong thời gian biên dịch.

Với nền tảng Windows, chúng ta có thể dùng công cụ [VC](https://en.wikipedia.org/wiki/Microsoft_Visual_C%2B%2B) để sinh ra thư viện động (sẽ có một số thư viện Windows phức tạp chỉ có thể được build với `VC`). Đầu tiên, chúng ta phải tạo một file định nghĩa cho `number.dll` để quản lý các kí hiệu dùng để export thư viện động.

Nội dung của file `number.def` như sau:

```
LIBRARY number.dll

EXPORTS
number_add_mod
```

Dòng đầu tiên `LIBRARY` sẽ chỉ ra tên của file và tên của thư viện động, và sau đó là mệnh đề  `EXPORTS` theo sau bởi một danh sách các tên dùng để export.

Giờ đây, chúng ta có thể dùng những lệnh sau để tạo ra thư viện động (cần dùng `VC` tools).

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

Nên chú ý rằng, tại thời điểm thực thi, thư viện động cần được đặt ở cùng nơi để system có thể thấy. Trên Windows, bạn có thể đặt dynamic library và executable program trong cùng một thư mục, hoặc thêm đường dẫn tuyệt đối của thư mục trong khi dynamic library được đưa vào biến môi trường PATH. Trong MacOS, bạn cần phải thiết lập biến môi trường DYLD_LIBRARY_PATH. Trong hệ thống Linux, bạn cần thiết lập biến LD_LIBRARY_PATH.

## 2.9.3.  Exporting thư việc C tĩnh

CGO không chỉ được dùng trong thư viện C tĩnh, các export functions được hiện thực bởi Go hoặc C static libraries. Chúng ta có thể dùng Go để hiện thực modulo addition function như phần trước ở ví dụ như sau đây.

Tạo `number.go` với nội dung như sau:

```go
package main

import "C"

func main() {}

//export number_add_mod
func number_add_mod(a, b, mod C.int) C.int {
    return (a + b) % mod
}
```

Theo như mô tả của tài liệu CGO, chúng ta cần export C function trong main package. Với cách xây dựng thư viện C tĩnh, hàm main trong main package được bỏ qua, và hàm C sẽ đơn giản được export. Xây dựng các lệnh sau:

```
$ go build -buildmode=c-archive -o number.a
```

Khi sinh ra thư viện tĩnh `number.a`, cgo cũng sẽ sinh ra file `number.h`.

Nội dung của `number.h` sẽ như sau:

```c
#ifdef __cplusplus
extern "C" {
#endif

extern int number_add_mod(int p0, int p1, int p2);

#ifdef __cplusplus
}
#endif
```

Phần ngữ pháp `extern "C"` được thêm vào để  dùng đồng thời trên cả ngôn ngữ C và C++. Nội dung của phần lõi sẽ định nghĩa hàm number_add_mod để được export.

Sau đó chúng ta tạo ra file _test_main.c để kiểm tra việc sinh ra C static library (bên dưới là phần prefix dùng cho việc xây dựng C static library để bỏ qua file đó).

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
 

## 2.9.4. Exporting thư viện C động

Quá trình exporting một thư viện động bằng CGO sẽ tương tự như một thư viện tĩnh, ngoại trừ build mode sẽ thay đổi c-shared và output file name được đổi thành `number.so`.

```
$ go build -buildmode=c-shared -o number.so
```

Nội dung của file _test_main.c sẽ không thay đổi, sau đó biên dịch và chạy chúng với các lệnh sau:

```
$ gcc -o a.out _test_main.c number.so
$ ./a.out
```

[Tiếp theo](ch2-10-link.md)