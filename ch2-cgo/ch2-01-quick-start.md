# 2.1 Quick Start

Trong phần này, chúng tôi sẽ hướng dẫn cách sử dụng cơ bản của CGO thông qua một loạt các ví dụ nhỏ từ dễ đến khó.

## 2.1.1 Chương trình CGO đơn giản

Các chương trình CGO thực tế thường phức tạp hơn. Nhưng chúng ta có thể đi từ chương trình dễ đến khó. Để xây dựng một chương trình CGO đơn giản, trước tiên hãy bỏ qua một số tính năng CGO phức tạp. Đây là chương trình CGO đơn giản nhất mà chúng tôi xây dựng:

[>> mã nguồn](../examples/ch2/ch2.1/1-simplest-cgo/main.go)

```go
// main.go
package main

import "C"//Sẵn sàng lập trình CGO

func main() {
    println("hello cgo")
}
```

Chúng ta import package CGO thông qua câu lệnh `import "C"`. Chương trình trên chưa thực hiện bất kì thao tác nào với CGO, chỉ mới thông báo sẵn sàng cho việc lập trình với CGO. Mặc dù chúng ta chưa sử dụng gì đến CGO nhưng lệnh `go build` vẫn sẽ gọi trình biên dịch `gcc` trong suốt quá trình biên dịch do nó đã là một chương trình CGO hoàn chỉnh.

## 2.1.2 Xuất chuỗi dựa trên thư viện chuẩn của C

```go
// main.go
package main

//#include <stdio.h>
import "C"

func main() {
    C.puts(C.CString("Hello World\n"))  
}
```

[>> mã nguồn](../examples/ch2/ch2.1/2-cputs/main.go)

Chúng ta `import package "C"` để thực hiện các chức năng của CGO và include thư viện <stdio.h> của ngôn ngữ C. Tiếp theo, chuỗi string trong `C.CString` của ngôn ngữ Go được chuyển đổi thành chuỗi string trong ngôn ngữ C bằng phương thức `C.puts` của gói CGO. Cuối cùng phương thức của package CGo được gọi để in ra kết quả.

So với các ngôn ngữ khác trên thế giới khi in câu "Hello World", điểm khác biệt lớn nhất của chương trình CGO là chương trình chúng ta sẽ không giải phóng trước khi chương trình kết thúc việc tạo chuỗi bằng lệnh C.CString. Ở đó chúng ta chuyển phương thức `puts` để in sang đầu ra tiêu chuẩn (stdout) trước khi áp dụng việc in bằng `fputs`.

Việc lỗi xảy ra khi không giải phóng chuỗi được tạo bằng C.CString của ngôn ngữ C sẽ dẫn đến việc rò rỉ bộ nhớ. Nhưng đối với chương trình nhỏ trên điều này không đáng lo ngại, bởi vì hệ điều hành sẽ tự động lấy lại các tài nguyên của chương trình sau khi chương trình kết thúc.

## 2.1.3 Sử dụng hàm C tự khai báo

Phần trên chúng tôi đã sử dụng các chức năng đã có trong thư viện tiêu chuẩn. Bây giờ, chúng tôi sẽ tùy chỉnh một hàm `SayHello` của ngôn ngữ C. Chức năng hàm này là in ra chuỗi chúng ta truyền vào hàm. Sau đó gọi hàm `SayHello` trong hàm main.

```go
// main.go
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

[>> mã nguồn](../examples/ch2/ch2.1/3a-cfunction/main.go)

Chúng ta có thể cài đặt hàm `SayHello` trong file nguồn với đuôi tệp là `.c`. Bởi vì hàm `SayHello` được viết bởi một tệp riêng biệt, để có thể sử dụng hàm `SayHello` chúng ta cần loại bỏ các dấu `*/`

Chúng ta tạo file hello.go và cài đặt như sau:

```C
// hello.c
#include <stdio.h>

void SayHello(const char* s) {
    puts(s);
}
```

[>> mã nguồn](../examples/ch2/ch2.1/3b-cfunction/hello.c)

Sau đó bên file main.go chúng chỉ cần khai báo hàm `SayHello` trong phần CGO như bên dưới.

```Go
// main.go
package main

//void SayHello(const char* s);
import "C"

func main() {
    C.SayHello(C.CString("Hello World\n"))
}
```

[>> mã nguồn](../examples/ch2/ch2.1/3b-cfunction/main.go)

`Lưu ý`: thay vì chạy lệnh `go run main.go` hoặc `go build main.go`, chúng ta phải sử dụng `go run "your/package"` hoặc `go build "your/package"`. Nếu bạn đang đứng trong thư mục chứa mã nguồn thì bạn có thể chạy chương trình bằng lệnh `go run .` hoặc `go build .`

Vì `SayHello` được đặt trong file riêng, ta có thể biên dịch thành các thư viện tĩnh hoặc động để sử dụng. Nếu sử dụng dưới dạng thư viện, file nguồn (`hello.c`) cần được đưa ra ngoài thư mục hiện tại (CGO tự động build các file nguồn của C, gây ra xung đột tên hàm). Chi tiết sẽ được đề cập sau.

## 2.1.4 Module hóa C code

Trừu tượng và module hóa là cách để đơn giản hóa các vấn đề phức tạp trong lập trình. Khi code dài hơn, ta có thể đưa các lệnh tương tự nhau vào chung một hàm. Khi có nhiều hàm hơn, ta chia chúng vào các file hoặc module. Cốt lõi của việc này là lập trình theo `interface` (interface không phải là khái niệm interface trong ngôn ngữ Go mà là khái niệm về API).

Trong ví dụ trước, ta trừu tượng hóa một module tên là `hello` và tất cả các interface của module đó được khai báo trong file header `hello.h`:

```h
// hello.h
void SayHello(const char* s);
```

[>> mã nguồn](../examples/ch2/ch2.1/4-modularization/hello.h)

Và chỉ có 1 khai báo cho hàm `SayHello` nhưng ta có thể an tâm sử dụng mà không phải lo lắng về việc hiện thực cụ thể  hàm đó. Khi hiện thực hàm `SayHello`, ta chỉ cần đáp ứng đúng đặc tả của khai báo hàm trong file header. Ví dụ sau là hiện thực hàm `SayHello` trong file `hello.c`:

```c
// hello.c

#include "hello.h" // Đảm bảo việc hiện thực hàm thỏa mãn interface của module.
#include <stdio.h>

void SayHello(const char* s) {
    puts(s);
}
```

[>> mã nguồn](../examples/ch2/ch2.1/4-modularization/hello.c)

Trong file `hello.c` chúng ta include file `hello.h` và sau đó cài đặt hàm SayHello đúng như đặc tả ở file `hello.h`.

File interface `hello.h` chỉ là thỏa thuận giữa người hiện thực và người sử dụng của module `hello`. Ta có thể hiện thực nó bằng ngôn ngữ C hoặc C++.

```cpp
// hello.cpp

#include <iostream>

extern "C" {
    #include "hello.h"
}

void SayHello(const char* s) {
    std::cout << s;
}
```

[>> mã nguồn](../examples/ch2/ch2.1/4-modularization/hello.cpp)

Tuy nhiên, để đảm bảo rằng hàm SayHello được hiện thực bởi C++ đáp ứng đặc tả hàm bởi file header của ngôn ngữ C, ta cần phải thêm lệnh `extern "C"` để chỉ ra rằng mối liên hệ đó ([hello.h](../examples/ch2/ch2.1/4-modularization/hello.h) và [hello.cpp](../examples/ch2/ch2.1/4-modularization/hello.cpp)) vẫn tuân theo quy tắc của C.

Với việc lập trình C bằng API interface, ta có thể hiện thực module bằng bất kỳ ngôn ngữ nào, miễn là đáp ứng được API. Ta có thể hiện thực SayHello bằng C, C++, Go hoặc kể cả Assembly.

## 2.1.5 Sử dụng Go để hiện thực hàm trong C

Trong thực tế, CGO không chỉ được sử dụng để gọi các hàm của C bằng ngôn ngữ Go mà còn được dùng để xuất các hàm (viết bằng) ngôn ngữ Go sang các lời gọi hàm của C.

Trong ví dụ trước, chúng ta đã trừu tượng hóa một module có tên hello và tất cả các chức năng interface của module được xác định trong tệp header hello.h:

```h
// hello.h
void SayHello(/*const*/ char* s);
```

[>> mã nguồn](../examples/ch2/ch2.1/5-implement-function-go/hello.h)

Bây giờ, chúng ta tạo một tệp hello.go và hiện thực lại chức năng SayHello của interface ngôn ngữ C bằng ngôn ngữ Go:

```go
// hello.go
package main

import "C"

import "fmt"

//export SayHello
func SayHello(s *C.char) {
    fmt.Print(C.GoString(s))
}
```

[>> mã nguồn](../examples/ch2/ch2.1/5-implement-function-go/hello.go)

Ta sử dụng chỉ thị `//export SayHello` của CGO để xuất hàm được hiện thực bằng Go sang hàm sử dụng được cho C. Tuy nhiên để đáp ứng được các hàm của ngôn ngữ C được hiện thực bằng Go, ta cần bỏ `const` trong file header. Vậy nên cần chú ý là ta sẽ có hai phiên bản `SayHello`: một là trong môi trường cục bộ của Go, hai là của C. Phiên bản SayHello của C được sinh ra bởi CGO cuối cùng cũng sẽ gọi phiên bản SayHello của Go thông qua `bridge code`.

Với việc lập trình ngôn ngữ C qua inteface, ta có thể tự do hiện thực và đơn giản hóa việc sử dụng hàm. Bây giờ ta có thể dùng SayHello như là một thư viện:

```go
package main

//#include <hello.h>
import "C"

func main() {
    C.SayHello(C.CString("Hello World\n"))
}
```

[>> mã nguồn](../examples/ch2/ch2.1/5-implement-function-go/main.go)

## 2.1.6 Sử dụng Go để lập trình interface cho C

Trong ví dụ trên, tất cả đoạn mã CGO của chúng ta đều nằm trong tệp Go. Sau đó, SayHello được chia thành các tệp C khác nhau bằng kỹ thuật lập trình interface C và hàm main vẫn là tệp Go. Sau đó, hàm SayHello của interface ngôn ngữ C được thực hiện lại bằng hàm trong Go. Nhưng đối với ví dụ hiện tại chỉ có một chức năng và việc chia thành ba tệp khác nhau thì hơi cồng kềnh.

Nếu các bạn làm những project lớn thì nên chia rõ ràng các file như ví dụ trên, còn ở đây chúng tôi sẽ gộp lại thành một file main.go duy nhất như ví dụ dưới đây.

```go
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

[>> mã nguồn](../examples/ch2/ch2.1/6a-go-programming/main.go)

Tỉ lệ đoạn mã C trong chương trình bây giờ ít hơn. Tuy nhiên vẫn phải sử dụng chuỗi trong C thông qua hàm `C.CString` chứ không thể dùng trực tiếp chuỗi của Go. Trong `Go1.10`, CGO đã thêm một loại `_GoString_pred` xác để thể hiện chuỗi trong ngôn ngữ Go. Đây là mã nguồn được cải tiến.

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

[>> mã nguồn](../examples/ch2/ch2.1/6b-go-programming/main.go)

Mặc dù có vẻ như tất cả đều được viết bằng ngôn ngữ Go, nhưng việc triển khai từ hàm `main()` của ngôn ngữ Go đến phiên bản ngôn ngữ C đã tự động tạo ra hàm SayHello, và cuối cùng cũng trở lại môi trường ngôn ngữ Go. Đoạn mã này vẫn chứa bản chất của lập trình CGO và người đọc cần hiểu sâu về nó.
