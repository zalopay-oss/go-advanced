# 1.1 Nguồn gốc của ngôn ngữ Go

Ngôn ngữ **Go** ban đầu được thiết kế và phát triển bởi một nhóm kĩ sư Google bao gồm **Robert Griesemer**, **Ken Thompson** và **Rob Pike** vào năm 2007. Mục đích của việc thiết kế ngôn ngữ mới bắt nguồn từ một số phản hồi về tính chất phức tạp của C++11, cuối cùng là nhu cầu thiết kế lại ngôn ngữ C trong môi trường network và multi-core (đa lõi). Vào giữa năm 2008, hầu hết các tính năng được thiết kế trong ngôn ngữ được hoàn thành, họ bắt đầu hiện thực trình biên dịch và bộ thực thi với **Russ Cox** là nhà phát triển chính. Trước năm 2010, ngôn ngữ Go dần dần được hoàn thiện. Vào tháng 9 cùng năm, ngôn ngữ Go chính thức được công bố dưới dạng open source (mã nguồn mở).

Ngôn ngữ Go thường được mô tả là "Ngôn ngữ tựa C" hoặc là "Ngôn ngữ C của thế kỉ 21". Từ nhiều khía cạnh, ngôn ngữ Go thừa hưởng những ý tưởng từ ngôn ngữ C, như là cú pháp, cấu trúc điều khiển, kiểu dữ liệu cơ bản, thủ tục gọi - trả về, con trỏ, v,v.., hoàn toàn kế thừa và phát triển ngôn ngữ C - một ngôn ngữ cấp thấp, hình bên dưới mô tả sự liên quan của ngôn ngữ Go với các ngôn ngữ khác. Theo đó chúng ta có thể thấy những ngôn ngữ ảnh hưởng tới Go.

<p align="center">
<img src="../images/ch1-1-go-family-tree.png" width="600"/>
</p>

*Hình 1-1 Cây phả hệ của ngôn ngữ Go*

Đầu tiên, quan sát phía bên trái của sơ đồ, có thể được nhìn thấy rõ ràng rằng tính chất **concurrency** (đồng thời) của ngôn ngữ **Go** được phát triển từ học thuyết **CSP** được công bố bởi **Bell Labs' Hoare** vào năm 1978. Sau đó, mô hình **CSP** concurrency dần dần được tinh chế và được ứng dụng thực tế trong một số ngôn ngữ lập trình như là **Squeak/NewSqueak** và **Alef**. Những thực tiễn thiết kế mô hình **CSP** đó cuối cùng được hấp thu bởi ngôn ngữ Go. Mô hình concurrency của thế hệ ngôn ngữ Erlang là một hiện thực khác của học thuyết **CSP**.

Chính giữa của sơ đồ chủ yếu thể hiện tính chất hướng đối tượng và đóng gói của **Go** được kế thừa từ **Pascal** - ngôn ngữ được thiết kế bởi Niklaus Wirth và những ngôn ngữ liên quan khác dẫn xuất từ chúng. Những cú pháp `package concept`, `package import` và `declaration` chủ yếu đến từ ngôn ngữ Modula-2. Cú pháp của các phương thức hỗ trợ tính hướng đối tượng đến từ ngôn ngữ Oberon, ngôn ngữ Go được phát triển có thêm những tính chất đặc trưng như là `implicit interface` để chúng hỗ trợ mô hình `duck object-oriented`.

Cuối cùng, cột bên phải của sơ đồ là ngôn ngữ C. Ngôn ngữ Go được thừa hưởng hầu hết
 từ C, không chỉ về cú pháp mà là nhiều thứ khác, quan trọng nhất là tách rời tính linh hoạt khỏi sự nguy hiểm khi dùng `pointer` (con trỏ). Ngoài ra, ngôn ngữ Go như là một bản thiết kế lại với sự ưu tiên về việc hỗ trợ ít toán tử hơn C, được đánh bóng cũng như thay đổi vẻ ngoài. Dĩ nhiên, tính chất trực tiếp can thiệp vào hệ thống như C cũng được phát triển trong Go (ngôn ngữ Go có tổng cộng 25 từ khóa, và những đặc tả của ngôn ngữ chỉ vỏn vẹn 50 trang).

Một vài những tính năng khác của ngôn ngữ Go đến từ một số ngôn ngữ khác, ví dụ là cú pháp `iota` được mượn từ ngôn ngữ **APL**, những đặc điểm như là `lexical scope` và `nested functions` đến từ Scheme. Cũng có những ý tưởng khác được thiết kế và đưa vào. Ví dụ, Go hỗ trợ `slice` để truy cập phần tử nhanh như mảng tĩnh, đồng thời nó có thể được tăng giảm kích thước bằng cơ chế chia sẻ vùng nhớ tương tự linked list, mệnh đề `defer` có trong Go (phát minh của Ken) cũng rất hữu ích.

## 1.1.1 Di truyền từ Bell Labs

Khả năng lập trình concurrency của Go đến từ một nghiên cứu ít người biết tới và được công bố bởi Tony Hoarce tại Bell Labs vào năm 1978 -  Commutative sequential processes (CSP). Về bài báo khoa học nói về CSP, chương trình chỉ là một tập hợp các tiến trình được chạy song song, mà không có sự chia sẻ về trạng thái, sử dụng `pipes` cho việc giao tiếp và điều khiển đồng bộ. Mô hình Tony Hoare's CSP concurrency chỉ là một ngôn ngữ mô tả cho những khái niệm cơ bản về concurrency (tính đồng thời), nó cũng không hẳn là một ngôn ngữ lập trình.

Ví dụ kinh điển của việc áp dụng mô hình CSP concurrent là ngôn ngữ  **Erlang**, được phát triển bởi **Ericsson**. Tuy nhiên, trong khi Erlang sử dụng học thuyết CSP trong mô hình lập trình concurrency, Rob Pike là người cũng đến từ Bell Labs và đồng nghiệp của ông cũng thử giới thiệu mô hình CSP concurrency vào việc phát triển một ngôn ngữ mới tại thời điểm đó. Lần đầu tiên họ cố gắng giới thiệu mô hình CSP concurency trong một ngôn ngữ được gọi là Squeak, đây là một ngôn ngữ xử lý sự kiện từ chuột và bàn phím trong những pipe được khởi tạo tĩnh. Sau đó có một phiên bản cải thiện là NewSqueak, cú pháp và mệnh đề của chúng cũng tương tự như C, và Pascal. NewSqueak thì chỉ là một ngôn ngữ lập trình hàm với cơ chế thu gom vùng nhớ thừa tự động từ các sự kiện của bàn phím, chuột, và màn hình. Tuy nhiên trong ngôn ngữ Newsquek, pipeline đã thực sự được khởi tạo động, pipeline là kiểu giá trị đầu tiên có thể lưu trữ trong biến. Sau đó, ngôn ngữ Alef (đó cũng là một ngôn ngữ được ưa thích bởi Ritchie - cha đẻ của ngôn ngữ C). Alef là sự chuyển đổi của Newsqueak thành một ngôn ngữ lập trình hệ thống, nhưng cực kỳ khó khăn để có mô hình concurrency trong một ngôn ngữ thiếu cơ chế thu gom vùng nhớ tự động (trong C, ta phải gọi hàm `free()` thủ công để làm việc này). Có một ngôn ngữ khác tên là Limbo sau ngôn ngữ Alef, nó là một ngôn ngữ nhúng chạy trên máy ảo. Limbo là thế hệ gần nhất với ngôn ngữ Go, nó có những cú pháp tương tự như Go. Qua việc thiết kế Go, Rob Pike đã tổng hợp nhiều thập kỉ trong việc thiết kế mô hình CSP concurrent. Tính chất lập trình concurrency trong Go hơi phức tạp, chúng sẽ được đề cập trong bộ tài liệu này.

*Hình 1-2 chỉ ra lịch sử phát triển của ngôn ngữ Go qua codebase logs (Git is git log --before={2008-03-03} --reverseviewed with commands).*

<p align="center">
<img src="../images/ch1-2-go-log4.png" width="600"/>
</p>

*Hình 1-2 Go language development log*

Có thể nhìn thấy từ những submission log rằng ngôn ngữ Go được dần phát triển từ ngôn ngữ B - được phát minh bởi **Ken Thompson** và ngôn ngữ C được phát triển bởi **Dennis M.Ritchie**. Đó là thế hệ ngôn ngữ C đầu tiên, do đó nhiều người gọi Go là ngôn ngữ lập trình C của thế kỉ 21.

Hình 1-3 chỉ ra cuộc cách mạng của các thế hệ ngôn ngữ lập trình từ Bell Labs và đến Go:

<p align="center">
<img src="../images/ch1-3-go-history.png" width="600"/>
</p>

*Hình 1-3 Lịch sử phát triển của lập trình concurrency trong Go*

Trong suốt quá trình phát triển ngôn ngữ lập trình từ Bell Labs, từ B đến C, NewSqueak, Alef, Limbo, ngôn ngữ Go thừa hưởng một nửa thế kỉ của việc thiết kế từ những thế hệ trước, cuối cùng hoàn thành sứ mệnh tạo ra một thế hệ ngôn ngữ tựa C mới. Trong vòng những năm gần đâu, Go trở thành một ngôn ngữ lập trình vô cùng quan trọng trong `cloud computing` và `cloud storage`.

## 1.1.2 Hello World

Để bắt đầu, chương trình đầu tiên thường in ra dòng chữ "Hello World", đoạn code bên dưới là chương trình này.

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```

Lưu đoạn code trên thành file `hello.go`. Sau đó, chuyển tới thư mục chứa file `hello.go`, nơi mà file được lưu trữ, Bây giờ chúng ta có thể sử dụng Go như là ngôn ngữ scripting bằng cách gõ `go run hello.go` đó là một câu lệnh command line trực tiếp cho ra kết quả là dòng chữ "Hello World".

Bây giờ, giới thiệu ngắn về chương trình trên, tất cả những chương trình Go sẽ được thập hợp thành những đơn vị cơ bản là hàm và biến, Một hàm và biến được tổ chức thành các mã nguồn (source file). Những source đó được tổ chức thành một package phù hợp theo ý định của tác giả. Cuối cùng, những package đó cũng được tổ chức thành một khối thống nhất, chúng cấu thành chương trình Golang. Function được sử dụng chứa những chuỗi statements (mệnh đề) và những biến lưu trữ dữ liệu. Tên của hàm khởi nguồn toàn chương trình được gọi là hàm main. Mặc dù không có nhiều quy định về việc đặt tên hàm trong Go, hàm main phải được đặt trong package main và là điểm khởi đầu của toàn chương trình. Package được sử dụng để đóng gói những hàm, biến, hằng có liên quan và sử dụng cú pháp import để khai báo package, ví dụ chúng ta có thể sử dụng hàm `Println` trong package `fmt`.

Hai dấu nháy kép chứa chuỗi "Hello World" là một kí tự để biểu diễn hằng `string` trong ngôn ngữ Go. Không giống như `string` trong C, nội dung `string` trong Go không thể chỉnh sửa. Khi truyền tham số `string` tới hàm fmt.Println, nội dung của `string` ko được sao chép - mà thực tế là địa chỉ và chiều dài của string được truyền vào (chúng được định nghĩa trong cấu trúc `reflect.StringHeader`). Trong ngôn ngữ Go, tham số hàm được truyền vào như là một bản sao chép (không hỗ trợ tham khảo).
