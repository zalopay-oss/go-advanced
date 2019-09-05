# 1.1. Nguồn gốc của ngôn ngữ Go

Ngôn ngữ [Go](https://golang.org/) ban đầu được thiết kế và phát triển bởi một nhóm kĩ sư Google bao gồm **Robert Griesemer**, **Ken Thompson** và **Rob Pike** vào năm 2007. Mục đích của việc thiết kế ngôn ngữ mới bắt nguồn từ một số phản hồi về tính chất phức tạp của C++11 và nhu cầu thiết kế lại ngôn ngữ C trong môi trường network và multi-core.

Vào giữa năm 2008, hầu hết các tính năng được thiết kế trong ngôn ngữ được hoàn thành, họ bắt đầu hiện thực trình biên dịch (compiler) và Go runtime với **Russ Cox** là nhà phát triển chính. Trước năm 2010, ngôn ngữ Go dần dần được hoàn thiện. Vào tháng 9 cùng năm, ngôn ngữ Go chính thức được công bố dưới dạng Open source.

Ngôn ngữ Go thường được mô tả là "Ngôn ngữ tựa C" hoặc là "Ngôn ngữ C của thế kỉ 21". Từ nhiều khía cạnh, ngôn ngữ Go thừa hưởng những ý tưởng từ ngôn ngữ C, như là cú pháp, cấu trúc điều khiển, kiểu dữ liệu cơ bản, thủ tục gọi, trả về, con trỏ, v,v.., hoàn toàn kế thừa và phát triển ngôn ngữ C, hình bên dưới mô tả sự liên quan của ngôn ngữ Go với các ngôn ngữ khác.

<div align="center">
	<img src="../images/ch1-1-go-family-tree.png" width="500">
	<br/>
	<span align="center">
		<i>Cây phả hệ của ngôn ngữ Go</i>
	</span>
</div>

Phía bên trái sơ đồ thể hiện tính chất **concurrency** của ngôn ngữ **Go** được phát triển từ học thuyết [CSP](https://en.wikipedia.org/wiki/Communicating_sequential_processes) công bố bởi **Tony Hoare** vào năm 1978. Học thuyết **CSP** dần dần được tinh chế và được ứng dụng thực tế trong một số ngôn ngữ lập trình như là **Squeak/NewSqueak** và **Alef**, cuối cùng là **Go**.

Chính giữa sơ đồ cho thấy tính chất hướng đối tượng và đóng gói của **Go** được kế thừa từ **Pascal** và những ngôn ngữ liên quan khác dẫn xuất từ chúng. Những từ khóa `package`, `import` đến từ ngôn ngữ Modula-2. Cú pháp hỗ trợ tính hướng đối tượng đến từ ngôn ngữ Oberon, ngôn ngữ Go được phát triển có thêm những tính chất đặc trưng như là `implicit interface` để chúng hỗ trợ mô hình [duck typing](https://en.wikipedia.org/wiki/Duck_typing).

Phía bên phải sơ đồ cho thấy ngôn ngữ **Go** kế thừa và cải tiến từ **C**, Cũng như **C**, **Go** là ngôn ngữ lập trình cấp thấp, nó cũng hỗ trợ con trỏ (pointer) nhưng ít nguy hiểm hơn **C**.

> "Go is the result of C programmers designing a new programming language, and Rust is the result of C++ programmers designing a new programming language" - [link](https://drewdevault.com/2019/03/25/Rust-is-not-a-good-C-replacement.html)

Một vài những tính năng khác của ngôn ngữ Go đến từ một số ngôn ngữ khác:
  * Cú pháp `iota` được mượn từ ngôn ngữ **APL**.
  * Những đặc điểm như là `lexical scope` và `nested functions` đến từ **Scheme**.
  * Go hỗ trợ `slice` để truy cập phần tử nhanh và có thể tự động tăng giảm kích thước.
  * Mệnh đề `defer` trong Go.

## 1.1.1. Khởi nguồn từ Bell Labs

Tính chất concurrency của Go đến từ học thuyết [Commutative sequential processes (CSP)](https://www.cs.cmu.edu/~crary/819-f09/Hoare78.pdf) được công bố bởi Tony Hoare tại Bell Labs vào năm 1978. Bài báo khoa học về CSP nói rằng chương trình chỉ là một tập hợp các tiến trình được chạy song song, mà không có sự chia sẻ về trạng thái, sử dụng `channel` cho việc giao tiếp và điều khiển đồng bộ.

Học thuyết CSP của Tony Hoare chỉ là một mô hình lập trình với những khái niệm cơ bản về concurrency (tính đồng thời), nó cũng không hẳn là một ngôn ngữ lập trình. Qua việc thiết kế Go, Rob Pike đã tổng hợp nhiều thập kỷ trong việc ứng dụng học thuyết CSP trong việc xây dựng ngôn ngữ lập trình.

Ngôn ngữ **Erlang** là một hiện thực khác của học thuyết **CSP**, bạn có thể tìm kiếm thông tin về ngôn ngữ này trên [trang chủ Erlang](https://www.erlang.org/).

Hình dưới chỉ ra lịch sử phát triển của ngôn ngữ Go qua codebase logs.

<div align="center">
	<img src="../images/ch1-2-go-log4.png" width="600">
	<br/>
	<span align="center">
		<i>Go language development log</i>
	</span>
</div>


Có thể nhìn thấy từ những submission log rằng ngôn ngữ Go được dần phát triển từ ngôn ngữ B - được phát minh bởi **Ken Thompson** và ngôn ngữ C được phát triển bởi **Dennis M.Ritchie**. Đó là thế hệ ngôn ngữ C đầu tiên, do đó nhiều người gọi Go là ngôn ngữ lập trình C của thế kỉ 21.

<div align="center">
	<img src="../images/ch1-3-go-history.png" width="500">
	<br/>
	<span align="center">
		<i>Lịch sử phát triển của lập trình concurrency trong Go</i>
	</span>
</div>


Trong vòng những năm gần đây, Go là một ngôn ngữ được ưa chuộng khi viết các chương trình Micro Services, vì những đặc tính nhỏ gọn, biên dịch nhanh, import thư viện từ github, cú pháp đơn giản nhưng hiện đại.

## 1.1.2. Hello World

Việc đầu tiên là cài đặt chương trình Go lang theo hướng dẫn trên trang chủ [golang.org](https://golang.org/).

Để bắt đầu, chương trình đầu tiên thường in ra dòng chữ "Hello World", đoạn code bên dưới là chương trình này.

```go
// package main chứa điểm thực thi đầu tiên của toàn chương trình
package main

// import gói thư viện "fmt" hỗ trợ in ra màn hình
import "fmt"

// main là hàm đầu tiên được chạy
func main() {

    // in ra dòng chữ "Hello World"
    fmt.Println("Hello World")
}
```

Lưu file trên thành `hello.go` và chạy bằng lệnh sau.

```sh
$ go run hello.go
Hello World
// hoặc có thể biên dịch ra file thực thi
$ go build
$ ./hello
Hello World
```
[Tiếp theo](ch1-02-hello-revolution.md)