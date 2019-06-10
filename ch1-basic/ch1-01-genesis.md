
## 1.1 Nguồn gốc của ngôn ngữ Go

Ngôn ngữ **Go** ban đầu được thiết kế và phát triển bởi một nhóm kĩ sư Google bao gồm Robert Griesemer, Ken Thompson và Rob Pike vào năm 2007. Mục đích của việc thiết kế ngôn ngữ mới bắt nguồn từ một số phản hồi về tính chất phức tạp của C++11, cuối cùng là nhu cầu thiết kế lại ngôn ngữ C trong môi trường network và muti-core (đa lõi). Vào giữa năm 2008, hầu hết các tính năng được thiết kế trong ngôn ngữ được hoàn thành, họ bắt đầu hiện thực trình biên dịch và runtime với Russ Cox là nhà phát triển chính. Trước năm 2010, ngôn ngữ Go dần dần được hoàn thiện. Vào tháng 9 cùng năm, ngôn ngữ Go chính thức được công bố trên danh nghĩa open sourced.

Ngôn ngữ Go thường được mô tả là "Ngôn ngữ tựa C" hoặc là "Ngôn ngữ C của thế kỉ 21". Từ nhiều khía cạnh, ngôn ngữ Go thừa hưởng những ý tưởng từ ngôn ngữ C, như là cú pháp, cấu trúc điều khiển, kiểu dữ liệu cơ bản, lời gọi - trả về hàm, con trỏ, v,v.. là phiên bản kế thừa và phát triển của ngôn ngữ C - một ngôn ngữ cấp thấp, hình bên dưới mô tả sự liên quan của ngôn ngữ Go với các ngôn ngữ khác. Theo đó chúng ta có thể thấy những ngôn ngữ ảnh hưởng tới Go.

<p align="center">
<img src="../images/ch1-1-go-family-tree.png" width="600"/>
</p>

*Figure 1-1 Cây phả hệ của ngôn ngữ Go*

Đầu tiên, quan sát phía bên trái của sơ đồ, có thể được nhìn thấy rõ ràng rằng tính chất **concurrency** của ngôn ngữ **Go** được phát triển từ lý thiết **CSP** được công bố bởi Bell Labs' Hoare vào năm 1978. Sau đó, mô hình **CSP** concurrency dần dần được tinh chế và được ứng dụng trong thực tế trong một số ngôn ngữ lập trình như là **Squeak/NewSqueak** và **Alef**. Những thực tiễn thiết kế mô hình **CSP** đó cuối cùng được hấp thu bởi ngôn ngữ Go. Mô hình concurrency của thế hệ ngôn ngữ Erlang là một hiện thực khác của lý thiết **CSP**.

Chính giữa của sơ đồ chủ yếu thể hiện tính chất hướng đối tượng và đóng gói của **Go** được kế thừa từ **Pascal** - ngôn ngữ được thiết kế bởi Niklaus Wirth và những ngôn ngữ liên quan khác dẫn xuất từ chúng. Những cú pháp `package concept`, `package import` và `declaration` chủ yếu đến từ ngôn ngữ Modula-2. Sự định nghĩa cú pháp của `methods` để hỗ trợ tính hướng đối tượng đến từ ngôn ngữ Oberon, ngôn ngữ Go được phát triển có thêm những tính chất đặc trưng như là `implicit interface` để chúng hỗ trợ mô hình `duck object-oriented`.

Cuối cùng cột bên phải của sơ đồ gene là ngôn ngữ C. Ngôn ngữ Go được thừa hưởng hầu hết của ngôn ngữ C, không chỉ về cú pháp mà là nhiều thứ khác. Điều quan trọng nhất là tách rời tính linh hoạt khỏi sự nguy hiểm của `pointer` (con trỏ). Ngoài ra, ngôn ngữ Go như là một bản thiết kế lại với sự ưu tiên về việc hỗ trợ ít toán tử hơn C, và được đánh bóng cũng như thay đổi vẻ ngoài. Dĩ nhiên, trực tiếp can thiệp vào hệ thống như C cũng được phát triển trong Go (ngôn ngữ Go có tổng cộng 25 từ khóa, và những đặc tả của ngôn ngữ vỏn vẹn 25 trang).

Một vài những tính năng khác của ngôn ngữ Go đến từ một số ngôn ngữ khác, ví dụ là cú pháp `iota` được mượn từ **APL**, những đặc điểm như là lexical scope và nested functions đến từ Scheme. Cũng có những ý tưởng khác được thiết kế và đưa vào. Ví dụ, Go hỗ trợ slice để truy cập phần tử nhanh như mảng tĩnh đồng thời nó có thể được tăng giảm kích thước và chia sẻ với cơ chế tương tự linked list, mệnh dề `defer` có trong Go (phát minh của Ken) cũng rất hữu ích.

### 1.1.1 Duy truyền từ Bell Labs

Tính biểu tượng về concurent programming của Go đến từ một nghiên cứu ít biết được công bố bởi Tony Hoarce của Bell Labs vào năm 1978 -  Commutative sequential processes (CSP). Vào báo cáo khoa học về CSP, chương trình chỉ là một tập hợp các process được chạy song song, mà không có sự chia sẻ về trạng thái, sử dụng pipes cho việc giao tiếp vào điểu khiển đồng bộ. Mô hình Tony Hoare's CSP concurrency chỉ là một ngôn ngữ mô tả cho những ý tưởng cơ bản của concurrency, và nó không phải mô hình lập trình thường thấy.

Ví dụ kinh điển của việc áp dụng mô hình CSP concurrent là ngôn Erlang được phát triển bởi Ericsson. Tuy nhiên, trong khi Erlang sử dụng giả thiết CSP trong mô hình concurreny, Rob Pike là người cũng đến từ Bell Labs và đồng nghiệp của ông cũng liên tục cố gắng giới thiệu mô hình CSP concurrency vào sự phát triển ngôn ngữ mới trong thời gian đó. Lần đầu tiên họ cố gắng giới thiệu mô hình CSP concurency trong một ngôn ngữ được gọi là Squeak, đây là một ngôn ngữ xử lý sự kiện từ chuột và bàn phím trong những pipe được khởi tạo tĩnh. Sau đó có một phiên bản cải thiện là NewSqueak, cú pháp và mệnh đề của chúng cũng tương tự như C, và Pascal. NewSqueak thì chỉ là một ngôn ngữ lập trình hàm với cơ chế thu gom vùng nhớ thừa tự động từ các sự kiện của bàn phím, chuột, và màn hình. Tuy nhiên trong ngôn ngữ Newsquek, pipeline đã thực sự được khởi tạo động,pipeline là kiểu giá trị đầu tiên có thể lưu trữ trong biến. Sau đó, ngôn ngữ Alef (đó cũng là một ngôn ngữ được ưa thích bởi Ritchie - cha đẻ của ngôn ngữ C). Alef được dùng để chuyển đổi Newsqueak thành một ngôn ngữ lập trình hệ thống, nhưng cực kỳ khó khăn để có mô hình concurrency trong một ngôn ngữ thiếu cơ chế thu gom vùng nhớ tự động như C. Có một ngôn ngữ khác tên là Limbo sau ngôn ngữ Alef, nó là một ngôn ngữ scripting chạy trên máy ảo. Limbo là ngôn ngữ gần nhất với ngôn ngữ Go, nó có những cú pháp tương tự như Go. Qua việc thiết kế Go, Rob Pike đã tổng hợp nhiều thập kỉ trong việc thiết kế mô hình CSP concurrent. Ý tưởng của concurrent programing trong Go hơi phức tạp, và sự xuất phát từ một ngôn ngữ mới là cũng là vấn đề của khóa học.


*Figure 1-2 shows the most straightforward evolution of the Go code library's early codebase logs (Git is git log --before={2008-03-03} --reverseviewed with commands).*

<p align="center">
<img src="../images/ch1-2-go-log4.png" width="600"/>
</p>

*Figure 1-2 Go language development log*

Có thể nhìn thấy từ những submission log rằng ngôn ngữ Go được dần phát triển từ ngôn ngữ B - được phát minh bởi Ken Thomposon và ngôn ngữ C được phát triển bởi Dennis M.Ritchie. Đó là thế hệ ngôn ngữ C đầu tiên, do đó nhiều người gọi Go là ngôn ngữ lập trình C của thế kỉ 21.

*Figure 1-3 shows the evolution of the unique concurrent programming genes from Bell Labs in Go:*

<p align="center">
<img src="../images/ch1-3-go-history.png" width="600"/>
</p>

*Figure 1-3 Go language concurrent evolution history*

Trong suốt quá trình phát triển ngôn ngữ lập trình từ Bell Lab, từ B đến C, NewSqueak, Alef, Limbo, ngôn ngữ Go thừa hưởng một nửa thế kỉ của việc thiết kế những thế hệ trước, cuối cùng hoàn thành sứ mệnh tạo ra một thế hệ ngôn ngữ tựạ C mới. Trong vòng những năm gần đâu, Go trở thành một ngôn ngữ lập trình vô cùng quan trọng trong `cloud computing` và `cloud storage`.

### 1.1.2 Hello, the World

Để thuận tiện, chương trình đầu tiên thường in ra dòng chữ "Hello, world", đoạn code bên dưới là chương trình này.

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World")
}
```

Lưu đoạn code trên thành file `hello.o`. Sau đó  chuyển tới thư mục chứa file `hello.go`, nơi mà file được lưu trữ, Bây giờ chúng ta có thể sử dụng ngôn ngữ Go như là ngôn ngữ scripting bằng cách gõ `go run hello.go` đó là một câu lệnh command line trực tiếp chuyển thành kết quả là dòng chữ "Hello World!".

Bây giờ, giới thiệu ngắn về chương trình Golang, tất cả những chương trình Go sẽ được thập hợp thành những đơn vị cơ bản là hàm và biến, Một hàm và biến được tổ chức thành các mã nguồn (source file). Những source đó được tổ chức thành một package phù hợp theo ý định của tác giả. Cuối cùng, những package đó cũng được tổ chức thành một khối thống nhất, chúng cấu thành chương trình Golang. Function được sử dụng chứa những chuỗi statements và những biến lưu trữ dữ liệu. Tên của hàm khởi nguồn toàn chương trình được gọi là hàm main. Mặc dù không có nhiều quy định về việc đặt tên hàm trong Go, hàm main phải được đặt trong package main và là điểm khởi đầu của toàn chương trình. Package được sử dụng để đóng gói những hàm, biến, hằng có liên quan và sử dụng cú pháp import để khai báo package, ví dụ chúng ta có thể sử dụng hàm Println trong package fmt.

Hai dấu nháy kép chứa chuỗi "Hello world" là một kí tự để biểu diễn hằng string trong ngôn ngữ Go Không gống như string trong C, nội dung string trong go là bất biến. Khi chúng truyền một tham số string tới hàm fmt.Println, nội dung của string ko được sao chép - mà thực tế địa chỉ và chiều dài của string được truyền vào. Trong ngôn ngữ Go, tham số hàm được truyền vào như là một bản sao chép (không hỗ trợ tham khảo).
