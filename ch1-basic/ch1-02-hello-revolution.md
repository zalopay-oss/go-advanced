# 1.2. Sự tiến hóa của "Hello World"

Trong phần trước, chúng ta đã cùng tìm hiểu sơ lược về các ngôn ngữ cùng họ với Go và các ngôn ngữ lập trình khác được Bell Labs phát triển. Ở phần này, chúng ta sẽ nhìn lại dòng thời gian phát triển của từng ngôn ngữ và xem cách mà chương trình "Hello World" phát triển thành phiên bản của ngôn ngữ Go hiện tại và hoàn thiện những sự thay đổi mang tính cách mạng của nó.

<div align="center">
	<img src="../images/ch1-4-go-history.png">
	<br/>
	<span align="center">
		<i>Lịch sử tiến hóa của ngôn ngữ Go</i>
	</span>
</div>


## 1.2.1. Ngôn ngữ B - Ken Thompson, 1972

B là một ngôn ngữ lập trình đa dụng được phát triển bởi Ken Thompson thuộc Bell Labs, cha đẻ của ngôn ngữ Go, được thiết kế để hỗ trợ phát triển hệ thống UNIX. Tuy nhiên, B khá thiếu sự linh hoạt trong hệ thống kiểu khiến cho nó rất khó sử dụng.

Phiên bản "Hello World" sau đây nằm trong *A Tutorial Introduction to the Language B*  được viết bởi Brian W. Kernighan (là người commit đầu tiên vào mã code của Go), chương trình như sau :

```B
main() {
    extrn a, b, c;
    putchar(a); putchar(b); putchar(c);
    putchar('!*n');
}
a 'hell';
b 'o  w';
c 'orld';
```

Vì thiếu sự linh hoạt của kiểu dữ liệu trong B, các nội dung `a/b/c` cần in ra chỉ có thể được định nghĩa bằng các biến toàn cục, đồng thời chiều dài của mỗi biến phải được căn chỉnh (aligned) về 4 bytes (cảm giác giống như viết ngôn ngữ assembly vậy). Sau đó hàm `putchar` được gọi nhiều lần để làm nhiệm vụ output, lần gọi cuối với `!*n` để xuất ra một dòng mới.

Từ khi B được thay thế (bởi C), nó chỉ còn xuất hiện trong một số tài liệu và trở thành lịch sử.

## 1.2.2. Ngôn ngữ C - Dennis Ritchie, 1974 ~ 1989

C được phát triển bởi Dennis Ritchie trên nền tảng của B, trong đó thêm các kiểu dữ liệu phong phú hơn và đạt được mục tiêu lớn là viết lại UNIX. Có thể nói C chính là nền tảng phần mềm quan trọng nhất của ngành CNTT hiện đại. Hiện tại, gần như tất cả các hệ điều hành chính thống đều được phát triển bằng C, cũng như rất nhiều phần mềm cơ bản cũng được phát triển bằng C. Các ngôn ngữ lập trình của họ C đã thống trị trong nhiều thập kỷ và vẫn sẽ còn sức ảnh hưởng trong hơn nửa thế kỷ nữa.

Trong hướng dẫn giới thiệu ngôn ngữ C được viết bởi Brian W. Kernighan vào khoảng năm 1974, phiên bản ngôn ngữ C đầu tiên của chương trình "Hello World" đã xuất hiện. Điều này cung cấp quy ước cho chương trình đầu tiên với "Hello World" cho hầu hết các hướng dẫn ngôn ngữ lập trình sau này.

```c
// hàm không trả về kiểu giá trị một cách tường minh,
// mặc định sẽ trả về kiểu `int`
main()
{
    //'prinf' không cần import khai báo hàm mà mặc định có thể được sử dụng
    printf("Hello World");

    // không cần một câu lệnh return nhưng mặc định sẽ trả về giá trị 0
}
```

Ví dụ này cũng xuất hiện trong bản đầu tiên của **_C Programming Language_** xuất bản năm 1978 bởi Brian W. Kerninghan và Dennis M. Ritchie (K&R).

Năm 1988, 10 năm sau khi giới thiệu hướng dẫn của K&R, phiên bản thứ 2 của **_C Programming Language_** cuối cùng cũng được xuất bản. Thời điểm này, việc chuẩn hóa ngôn ngữ ANSI C đã được hoàn thành sơ bộ, nhưng phiên bản chính thức của document vẫn chưa được công bố.

```c
// thêm '#include <stdio.h>' là header file chứa câu lệnh đặc tả
// dùng để khai báo hàm `printf`
#include <stdio.h>

main()
{
    printf("Hello World\n");
}

```

Đến năm 1989, tiêu chuẩn quốc tế đầu tiên cho ANSI C được công bố, thường được nhắc tới với tên C89. C89 là tiêu chuẩn phổ biến nhất của ngôn ngữ C và vẫn còn được sử dụng rộng rãi. Phiên bản thứ 2 của **_C Programming Language_** cũng được in lại bản mới:

```c
#include <stdio.h>
// 'void' được thêm vào danh sách các tham số hàm,
// chỉ ra rằng không có tham số đầu vào
main(void)
{
    printf("Hello World\n");
}
```

Tại thời điểm này, sự phát triển của ngôn ngữ C về cơ bản đã hoàn thành. C92/C99/C11 về sau chỉ hoàn thiện một số chi tiết trong ngôn ngữ. Do các yếu tố lịch sử khác nhau, C89 vẫn là tiêu chuẩn được sử dụng rộng rãi nhất.

## 1.2.3. Newsqueak - Rob Pike, 1989

Newsqueak là thế hệ thứ 2 của ngôn ngữ chuột do Rob Pike sáng tạo ra, ông dùng nó để thực hành mô hình CSP lập trình song song. Newsqueak nghĩa là ngôn ngữ squeak mới, với "squeak" là tiếng của con chuột, hoặc có thể xem là giống tiếng click của chuột. Ngôn ngữ lập trình squeak cung cấp các cơ chế xử lý sự kiện chuột và bàn phím. Phiên bản nâng cấp của Newsqueak có cú pháp câu lệnh giống như của C và các biểu thức có cú pháp giống như Pascal. Newsqueak là một ngôn ngữ chức năng (function language) thuần túy với bộ thu thập rác tự động cho các sự kiện bàn phím, chuột và cửa sổ.

Newsqueak tương tự như một ngôn ngữ kịch bản có chức năng in tích hợp. Chương trình "Hello World" của nó không có gì đặc biệt:

```c
// hàm 'print' có thể hỗ trợ nhiều tham số
print("Hello ", "World", "\n");
```

Bởi vì các tính năng liên quan đến ngôn ngữ Newsqueak và ngôn ngữ Go chủ yếu là đồng thời (concurrency) và pipeline nên ta sẽ xem xét các tính năng này thông qua phiên bản concurrency của thuật toán "sàng số nguyên tố". Nguyên tắc "sàng số nguyên tố" như sau:

<div align="center">
	<img src="../images/ch1-5-prime-sieve.png">
	<br/>
	<span align="center">
		<i>Sàng số nguyên tố</i>
	</span>
</div>


Chương trình "sàng số nguyên tố" cho phiên bản concurrency của ngôn ngữ Newsqueak như sau:

```go
// 'counter' dùng để xuất ra chuỗi gốc gồm các số tự nhiên vào các channel
counter := prog(c:chan of int) {
    i := 2;
    for(;;) {
        c <-= i++;
    }
};

// Mỗi hàm 'filter' tương ứng với mỗi channel lọc số nguyên tố mới.
// Những channel lọc số nguyên tố này lọc các chuỗi input theo
// sàng số nguyên tố hiện tại và đưa kết quả tới channel đầu ra.
filter := prog(prime:int, listen, send:chan of int) {
    i:int;
    for(;;) {
        if((i = <-listen)%prime) {
            send <-= i;
        }
    }
};

// Dòng đầu tiên của mỗi channel phải là số nguyên tố
// sau đó xây dựng sàng nguyên tố  dựa trên số nguyên tố mới này
sieve := prog() of chan of int {
    // 'mk(chan of int)' tạo 1 channel, tương tự như 'make(chan int)' trong Go.
    c := mk(chan of int);

    begin counter(c);
    prime := mk(chan of int);
    begin prog(){
        p:int;
        newc:chan of int;
        for(;;){
            prime <-= p =<- c;
            newc = mk();

            // 'begin filter(p,c,newc)' bắt đầu một hàm concurrency,
            // giống với câu lệnh 'go filter(p,c,newc)' trong Go.
            begin filter(p, c, newc);
            c = newc;
        }
    }();
    
    // 'become' dùng để trả về kết quả của hàm, tương tự như 'return'
    become prime;
};

// kết quả là các số nguyên tố còn lại trên sàng
prime := sieve();
```

Cú pháp xử lý concurrency và channel  trong ngôn ngữ Newsqueak khá tương tự với Go, ngay cả cách khai báo kiểu dạng hậu tố của 2 ngôn ngữ này cũng giống nhau.

## 1.2.4. Alef - Phil Winterbottom, 1993

Trước khi xuất hiện ngôn ngữ Go, ngôn ngữ Alef có thể xem là ngôn ngữ xử lý concurrency hoàn hảo, hơn nữa cú pháp và runtime của Alef về cơ bản tương thích hoàn hảo với ngôn ngữ C.  Tuy nhiên, do thiếu cơ chế phục hồi bộ nhớ tự động, việc quản lý tài nguyên bộ nhớ của cơ chế concurrency là vô cùng phức tạp. Hơn nữa, ngôn ngữ Alef chỉ cung cấp hỗ trợ ngắn hạn trong hệ thống Plan9 và các hệ điều hành khác không có môi trường phát triển Alef thực tế. Ngôn ngữ Alef chỉ có hai tài liệu công khai: **_Alef Language Specification_** và **_the Alef Programming Wizard_**. Do đó, không có nhiều thảo luận về ngôn ngữ Alef ngoài Bell Labs.

Hình sau đây là trạng thái concurrency của Alef:

<div align="center">
	<img src="../images/ch1-6-alef.png">
	<br/>
	<span align="center">
		<i>Mô hình concurrency trong Alef</i>
	</span>
</div>


Chương trình "Hello World" cho phiên bản concurrency của ngôn ngữ Alef:

```c
// Khai báo thư viện runtime chứa
// ngôn ngữ Alef
#include <alef.h>

void receive(chan(byte*) c) {
    byte *s;
    s = <- c;
    print("%s\n", s);
    terminate(nil);
}

void main(void) {
    chan(byte*) c;

    // tạo ra một channel chan(byte*)
    // tương tự make(chan []byte) của Go
    alloc c;

    // receive khởi động hàm trong proc và thread
    // tương ứng.
    proc receive(c);
    task receive(c);
    c <- = "hello proc or task";
    c <- = "hello proc or task";
    print("done\n");

    // kết thúc bằng lệnh terminate
    terminate(nil);
}
```

Ngữ pháp của Alef về cơ bản giống như ngôn ngữ C. Nó có thể được coi là ngôn ngữ C ++ dựa trên ngữ pháp của ngôn ngữ C.

## 1.2.5. Limbo - Sean Dorward, Phil Winterbottom, Rob Pike, 1995

Limbo (Hell) là ngôn ngữ lập trình để phát triển các ứng dụng phân tán chạy trên máy tính nhỏ. Nó hỗ trợ lập trình mô-đun, kiểm tra kiểu mạnh vào thời gian biên dịch và thời gian chạy, liên lạc bên trong process thông qua channel, có bộ thu gom rác tự động. Có các loại dữ liệu trừu tượng đơn giản. Limbo được thiết kế để hoạt động an toàn ngay cả trên các thiết bị nhỏ mà không cần bảo vệ bộ nhớ phần cứng. Ngôn ngữ Limbo chạy chủ yếu trên hệ thống Inferno.

Phiên bản  Limbo của chương trình "Hello World" như sau:

```c
// tương tự 'package Hello' trong go
implement Hello;

// import các module khác
include "sys.m"; sys: Sys;
include "draw.m";

Hello: module
{
    // cung cấp hàm khởi tạo và kiểu khai báo dạng hậu tố
    // khác với Go không có tham số
    init: fn(ctxt: ref Draw->Context, args: list of string);
};

init(ctxt: ref Draw->Context, args: list of string)
{
    sys = load Sys Sys->PATH;
    sys->print("Hello World\n");
}
```

## 1.2.6. Ngôn ngữ Go - 2007 ~ 2009

Bell Labs sau khi trải qua nhiều biến động dẫn tới việc nhóm phát triển ban đầu của dự án Plan9 (bao gồm Ken Thompson) cuối cùng đã gia nhập Google. Sau khi phát minh ra ngôn ngữ tiền nhiệm là Limbo hơn 10 năm sau, vào cuối năm 2007, cảm thấy khó chịu với các tính năng "khủng khiếp" của C, ba tác giả gốc của ngôn ngữ Go đã tập hợp lại quyết định dùng 20% thời gian rảnh của mình để tạo ngôn ngữ một ngôn ngữ mới, chống lại sự thống trị của C/C++ ở Google lúc bấy giờ.

Đặc tả ngôn ngữ Go ban đầu được viết vào tháng 3 năm 2008 và chương trình Go gốc được biên dịch trực tiếp vào C và sau đó được dịch thành mã máy. Tháng 5 năm 2008, các nhà lãnh đạo Google cuối cùng đã phát hiện ra tiềm năng to lớn của ngôn ngữ Go và bắt đầu hỗ trợ cho dự án, cho phép các tác giả dành toàn bộ thời gian của mình để hoàn thiện ngôn ngữ. Sau khi phiên bản đầu tiên của đặc tả ngôn ngữ Go được hoàn thành, trình biên dịch ngôn ngữ Go cuối cùng có thể tạo ra mã máy trực tiếp (mà không phải thông qua C).

### hello.go - Tháng 6 năm 2008

```go
package main

func main() int {
    // vẫn còn dấu ';' cuối câu
    print "Hello World\n";

    // cần câu lệnh return để trả về giá trị
    // một cách tường minh
    return 0;
}
```

### hello.go - 27 tháng 6 năm 2008

```go
package main

func main() {
    print "Hello World\n";

    // loại bỏ câu lệnh return
    // chương trình trả về mặc định
    // bằng lệnh gọi 'exit(0)'
}
```

### hello.go - 11 tháng 8 năm 2008

```go
package main

func main() {
    // hàm built-in 'print' được đổi thành dạng hàm thông thường
    print("Hello World\n");
}
```

### hello.go - 24 tháng 10 năm 2008

```go
package main

import "fmt"

func main() {
    // 'printf' có thể định dạng chuỗi giống trong C
    // và được đặt trong package 'fmt' (viết tắt cho 'format')
    // phần đầu của tên hàm vẫn là chữ thường, lúc này tính năng export
    // vẫn chưa xuất hiện
    fmt.printf("Hello World\n");
}
```

### hello.go - 15 tháng 1 năm 2009

```go
package main

import "fmt"

func main() {
    // chữ 'P' viết hoa chỉ ra rằng hàm được export
    // các chữ viết thường chỉ ra hàm được dùng trong
    // nội bộ package
    fmt.Printf("Hello World\n");
}
```

### hello.go - 11 tháng 12 năm 2009

```go
package main

import "fmt"

func main() {
    // dấu ';' cuối cùng cũng được loại bỏ
    fmt.Printf("Hello World\n")
}
```

## 1.2.7. Hello World! - V2.0

Sau nửa thế kỷ phát triển, ngôn ngữ Go không chỉ có thể in được phiên bản Unicode của "Hello World", mà còn có thể cung cấp service tương tự cho người dùng trên toàn thế giới. Phiên bản sau đây in ra kí tự tiếng Việt "Xin chào" và thời gian hiện tại của mỗi client truy cập vào service.

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func main() {
    fmt.Println("Please visit http://127.0.0.1:12345/")

    // sử dụng giao thức http để in ra chuỗi bằng lệnh 'fmt.Fprintf'
    // thông qua log package
    http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        s := fmt.Sprintf("Xin chào - Thời gian hiện tại: %s", time.Now().String())
        fmt.Fprintf(w, "%v\n", s)
        log.Printf("%v\n", s)
    })

    // khởi động service http
    if err := http.ListenAndServe(":12345", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
```

Lúc này, Go cuối cùng đã hoàn thành việc chuyển đổi từ ngôn ngữ C của kỷ nguyên đơn lõi sang một ngôn ngữ lập trình đa dụng của môi trường đa lõi trong kỷ nguyên Internet thế kỷ 21.

[Tiếp theo](ch1-03-array-string-and-slice.md)