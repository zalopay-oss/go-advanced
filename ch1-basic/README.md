# Chapter 1

## 1.2 Sự tiến hóa của "Hello, World"

Trong phần trước, chúng tôi đã giới thiệu sơ lược về các ngôn ngữ cùng họ với Go, đồng thời là các ngôn ngữ lập trình song song được phát triển bởi Bell Labs. Cuối cùng là phiên bản Go với chương trình "Hello, World" được trình bày. Trên thực tế, chương trình "Hello, World" là ví dụ điển hình nhất cho thấy các tính năng của những ngôn ngữ khác nhau. Trong phần này, chúng ta sẽ nhìn lại dòng thời gian phát triển của từng ngôn ngữ và xem cách mà  chương trình "Hello, World"  phát triển thành ngôn ngữ Go hiện tại và hoàn thành sứ mệnh cách mạng của nó.

![go-history](../images/ch1-4-go-history.png)

Hình 1-4. Lịch sử tiến hóa của ngôn ngữ Go

### 1.2.1 Ngôn ngữ B - Ken Thompson, 1972

Đầu tiên là ngôn ngữ B, là một ngôn ngữ lập trình đa dụng được phát triển bởi Ken Thompson thuộc Bell Labs, cha đẻ của ngôn ngữ Go, được thiết kế để hỗ trợ phát triển hệ thống UNIX. Tuy nhiên, do thiếu sự linh hoạt trong hệ thống kiểu khiến cho B rất khó sử dụng. Sau đó, đồng nghiệp của Ken Thompson là Denis Ritchie phát triển ngôn ngữ C dựa trên B. C cung cấp cơ chế kiểu đa dạng, giúp tăng khả năng diễn đạt của ngôn ngữ. Cho đến ngày nay C vẫn là một trong những ngôn ngữ lập trình được sử dụng phổ biến nhất trên thế giới. Từ khi B được thay thế, nó chỉ còn xuất hiện trong một số tài liệu và trở thành lịch sử.

Phiên bản "Hello, World" sau đây là từ hướng dẫn giới thiệu ngôn ngữ B được viết bởi Brian W. Kernighan (là người commit đầu tiên vào  mã code của Go), chương trình như sau :

```B
main() {
    extrn a, b, c;
    putchar(a); putchar(b); putchar(c);
    putchar('!*n');
}
a 'hell';
b 'o, w';
c 'orld';
```

Vì thiếu sự linh hoạt của kiểu dữ liệu trong B, các nội dung `a/b/c` cần in ra chỉ có thể được định nghĩa bằng các biến toàn cục, đồng thời chiều dài của mỗi biến phải được căn chỉnh (aligned) về 4 bytes (cảm giác giống như viết ngôn ngữ assembly vậy). Sau đó  hàm `putchar` được gọi nhiều lần để làm nhiệm vụ output, lần gọi cuối với `!*n` để xuất ra một dòng mới.

Nhìn chung, B là ngôn ngữ với các chức năng tương đối  đơn giản.

### 1.2.2 C - Dennis Ritchie, 1974 ~ 1989

Ngôn ngữ C được phát triển bởi Dennis Ritchie trên nền tảng của B, trong đó thêm các kiểu dữ liệu phong phú và cuối cùng đạt được mục tiêu lớn là viết lại UNIX. Có thể nói C chính là nền tảng phần mềm quan trọng nhất của ngành CNTT hiện đại. Hiện tại, gần như tất cả các hệ điều hành chính thống đều được phát triển bằng C, cũng như rất nhiều phần mềm cơ bản cũng được phát triển bằng C. Các ngôn ngữ lập trình của họ C đã thống trị trong nhiều thập kỷ và vẫn sẽ còn sức ảnh hưởng trong hơn nửa thế kỷ nữa.

Trong hướng dẫn giới thiệu ngôn ngữ C được viết bởi Brian W. Kernighan vào khoảng năm 1974, phiên bản ngôn ngữ C đầu tiên của chương trình "Hello World" đã xuất hiện. Điều này cung cấp quy ước cho chương trình đầu tiên với "Hello World" cho hầu hết các hướng dẫn ngôn ngữ lập trình sau này.

```c
main()
{
    printf("hello, world");
}
```

Một số điểm cần lưu ý về chương trình này:

- Đầu tiên là hàm `main` không trả về kiểu giá trị một cách tường minh, mặc định sẽ trả về kiểu `int`;
- Thứ hai, hàm `printf` không cần import khai báo hàm mà mặc định có thể được sử dụng.
- Cuối cùng, hàm `main` không cần một câu lệnh return nhưng mặc định sẽ trả về giá trị 0. Khi chương trình này xuất hiện, ngôn ngữ C khác hẳn tiêu chuẩn trước đó. Những gì chúng ta thấy là cú pháp ngôn ngữ C những ngày đầu: hàm không cần ghi giá trị trả về, các tham số hàm có thể bị bỏ qua và printf không cần include file header.

Ví dụ này cũng xuất hiện trong bản đầu tiên của ***C Programming Language*** xuất bản năm 1978 bởi Brian W. Kerninghan và Dennis M. Ritchie (K&R)

Năm 1988, 10 năm sau khi giới thiệu hướng dẫn của K&R, phiên bản thứ 2 của ***C Programming Language*** cuối cùng cũng được xuất bản. Thời điểm này, việc chuẩn hóa ngôn ngữ ANSI C đã được hoàn thành sơ bộ, nhưng phiên bản chính thức của document vẫn chưa được công bố. Tuy nhiên, chương trình "Hello World" trong cuốn sách đã thêm `#include <stdio.h>` là header file chứa câu lệnh đặc tả mới, dùng để khai báo hàm `printf` (trong tiêu chuẩn C89, chỉ riêng với hàm `printf`, có thể được dùng trực tiếp mà không cần khai báo hàm).

```c
#include <stdio.h>

main()
{
    printf("hello, world\n");
}

```

Đến năm 1989, tiêu chuẩn quốc tế đầu tiên cho ANSI C được công bố, thường được nhắc tới với tên C89. C89 là tiêu chuẩn phổ biến nhất của ngôn ngữ C và vẫn còn được sử dụng rộng rãi. Phiên bản thứ 2 của ***C Programming Language*** cũng được in lại bản mới, đối với đặc tả C89 mới này, `void` đã được thêm vào danh sách các tham số hàm, chỉ ra rằng không có tham số đầu vào.

```c
#include <stdio.h>

main(void)
{
    printf("hello, world\n");
}
```

Tại thời điểm này, sự phát triển của ngôn ngữ C về cơ bản đã hoàn thành. C92 / C99 / C11 về sau chỉ hoàn thiện một số chi tiết trong ngôn ngữ. Do các yếu tố lịch sử khác nhau, C89 vẫn là tiêu chuẩn được sử dụng rộng rãi nhất.

### 1.2.3 Newsqueak - Rob Pike, 1989

Newsqueak là thế hệ thứ 2 của ngôn ngữ chuột do Rob Pike sáng tạo ra, ông dùng nó để thực hành mô hình CSP lập trình song son. Newsqueak nghĩa là ngôn ngữ squeak mới, với "squeak" là tiếng của con chuột, hoặc có thể xem là giống tiếng click của chuột. Ngôn ngữ lập trình squeak cung cấp các cơ chế  xử lý sự kiện chuột và bàn phím. Phiên bản nâng cấp của Newsqueak có cú pháp câu lệnh giống như của C và các biểu thức có cú pháp giống như Pascal. Newsqueak là một ngôn ngữ chức năng (function language) thuần túy với bộ thu thập rác tự động cho các sự kiện bàn phím, chuột và cửa sổ.

Newsqueak tương tự như một ngôn ngữ kịch bản có chức năng in tích hợp. Chương trình "Hello World" của nó không có gì đặc biệt:

```c
print("Hello,", "World", "\n");
```

Từ chương trình trên, ngoài hàm `print` có thể hỗ trợ nhiều tham số, rất khó để thấy các tính năng liên quan đến ngôn ngữ Newsqueak. Bởi vì các tính năng liên quan đến ngôn ngữ Newsqueak và ngôn ngữ Go chủ yếu là đồng thời (concurrency) và pipeline. Do đó, ta sẽ xem xét các tính năng của ngôn ngữ Newsqueak thông qua phiên bản đồng thời của thuật toán "sàng số nguyên tố". Nguyên tắc "sàng số nguyên tố" như sau:

![prime-sieve](../images/ch1-5-prime-sieve.png)

Hình 1-5. Sàng số nguyên tố

Chương trình "sàng số nguyên tố" cho phiên bản đồng thời của ngôn ngữ Newsqueak như sau:

```go
// xuất 1 chuỗi số int từ 2 vào pipeline 
counter := prog(c:chan of int) {
    i := 2;
    for(;;) {
        c <-= i++;
    }
};

// Đối với chuỗi thu được từ pipeline listen, lọc ra các số là bội số của số nguyên tố
// gửi kết quả cho pipeline send
filter := prog(prime:int, listen, send:chan of int) {
    i:int;
    for(;;) {
        if((i = <-listen)%prime) {
            send <-= i;
        }
    }
};

// chức năng chính
// Dòng đầu tiên của mỗi pipeline phải là số nguyên tố
// sau đó xây dựng sàng nguyên tố  dựa trên số nguyên tố mới này
sieve := prog() of chan of int {
    c := mk(chan of int);
    begin counter(c);
    prime := mk(chan of int);
    begin prog(){
        p:int;
        newc:chan of int;
        for(;;){
            prime <-= p =<- c;
            newc = mk();
            begin filter(p, c, newc);
            c = newc;
        }
    }();
    become prime;
};

// kết quả là các số nguyên tố còn lại trên sàng
prime := sieve();
```

- Hàm `counter` dùng để xuất ra chuỗi gốc gồm các số tự nhiên vào các "đường ống" (pipeline). Mỗi hàm `filter` tương ứng với mỗi đường ống lọc số nguyên tố mới. Những đường ống lọc số nguyên tố này lọc các chuỗi đến theo sàng số nguyên tố hiện tại và đưa kết quả ra đường ống đầu ra. `mk(chan of int)` dùng để tạo 1 đường ống, tương tự như `make(chan int)` trong Go.
- Từ khóa `begin filter(p,c,newc)` bắt đầu một hàm đồng thời, giống với câu lệnh `go filter(p,c,newc)` trong Go.
- `become` dùng để trả về kết quả của hàm, tương tự như `return`.

Cú pháp xử lý đồng thời (concurrency) và đường ống (pipeline) trong ngôn ngữ Newsqueak khá tương tự với Go, ngay cả cách khai báo kiểu phía sau biến của 2 ngôn ngữ này cũng giống nhau.

### 1.2.4 Alef - Phil Winterbottom, 1993

Trước khi xuất hiện ngôn ngữ Go, ngôn ngữ Alef là ngôn ngữ xử lý đồng thời  hoàn hảo trong tâm trí của tác giả, hơn nữa cú pháp và thời gian chạy của Alef về cơ bản tương thích hoàn hảo với ngôn ngữ C. Hỗ trợ threads và process trong Alef là `proc receive(c)` dùng để bắt đầu một process và `task receive(c)` bắt đầu một thread với `c` để có thể giao tiếp qua pipes. Tuy nhiên, do thiếu cơ chế phục hồi bộ nhớ tự động, việc quản lý tài nguyên bộ nhớ của cơ chế đồng thời là vô cùng phức tạp. Hơn nữa, ngôn ngữ Alef chỉ cung cấp hỗ trợ ngắn hạn trong hệ thống Plan9 và các hệ điều hành khác không có môi trường phát triển Alef thực tế. Ngôn ngữ Alef chỉ có hai tài liệu công khai: ***Alef Language Specification*** và  ***the Alef Programming Wizard***. Do đó, không có nhiều thảo luận về ngôn ngữ Alef ngoài Bell Labs.

Vì ngôn ngữ Alef hỗ trợ cả thread và process trong cơ chế đồng thơi, và nhiều tiến trình đồng thời có thể bắt đầu đồng thơi, cho nên  trạng thái đồng thời của Alef là cực kỳ phức tạp. Cùng với đó, Alef cũng không có cơ chế thu gom rác tự động (Alef có tính năng con trỏ linh hoạt dành riêng cho ngôn ngữ C, điều này cũng khiến cơ chế thu gom rác tự động khó thực hiện).

Các tài nguyên khác nhau bị ngập giữa các thread và process khác nhau, ảnh hưởng lớn đến tài nguyên bộ nhớ đồng thời. Việc quản lý chúng sẽ vô cùng phức tạp. Ngôn ngữ Alef kế thừa cú pháp của ngôn ngữ C và có thể được coi là ngôn ngữ C tăng cường thêm cú pháp đồng thời. Hình ảnh sau đây là trạng thái đồng thời trong tài liệu ngôn ngữ Alef:

![concurrent](../images/ch1-6-alef.png)

Hình 1-6. Mô hình đồng thời trong Alef


Chương trình "Hello World" cho phiên bản đồng thời của ngôn ngữ Alef:

```c
#include <alef.h>

void receive(chan(byte*) c) {
    byte *s;
    s = <- c;
    print("%s\n", s);
    terminate(nil);
}

void main(void) {
    chan(byte*) c;
    alloc c;
    proc receive(c);
    task receive(c);
    c <- = "hello proc or task";
    c <- = "hello proc or task";
    print("done\n");
    terminate(nil);
}
```

Câu lệnh `#include <alef.h>` ở đầu chương trình dùng để khai báo thư viện runtime có chứa ngôn ngữ Alef. `receive` là một hàm bình thường, chương trình sử dụng như hàm nhập cho mỗi hàm đồng thời. câu lệnh `alloc c` trong hàm `main` trước tiên tạo ra một `chan(byte*)` loại đường ống, tương tự như `make(chan []byte)` của  Go , sau đó `receive` khởi  động hàm trong process và thread tương ứng, sau khi bắt đầu quá trình đồng thời, hàm `main` gửi đi  hai dữ liệu chuỗi tới đường ống, hàm `receive` chạy trong process và thread nhận dữ liệu từ đường ống theo thứ tự không xác định, sau đó in riêng các chuỗi, cuối cùng mỗi chuỗi đồng thời  tự kết thúc bằng cách gọi `terminate(nil)`.

Ngữ pháp của Alef về cơ bản giống như ngôn ngữ C. Nó có thể được coi là ngôn ngữ C ++ dựa trên ngữ pháp của ngôn ngữ C.