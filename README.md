# Go Language Advanced Programming

- [Go Language Advanced Programming](#go-language-advanced-programming)
  - [Giới thiệu](#gi%e1%bb%9bi-thi%e1%bb%87u)
  - [Tại sao chúng tôi thực hiện bộ tài liệu này ?](#t%e1%ba%a1i-sao-ch%c3%bang-t%c3%b4i-th%e1%bb%b1c-hi%e1%bb%87n-b%e1%bb%99-t%c3%a0i-li%e1%bb%87u-n%c3%a0y)
  - [Đối tượng sử dụng](#%c4%90%e1%bb%91i-t%c6%b0%e1%bb%a3ng-s%e1%bb%ad-d%e1%bb%a5ng)
  - [Tài liệu tham khảo](#t%c3%a0i-li%e1%bb%87u-tham-kh%e1%ba%a3o)
  - [Mục lục](#m%e1%bb%a5c-l%e1%bb%a5c)
  - [Phương thức đọc](#ph%c6%b0%c6%a1ng-th%e1%bb%a9c-%c4%91%e1%bb%8dc)
  - [Tham gia phát triển](#tham-gia-ph%c3%a1t-tri%e1%bb%83n)
  - [Nhóm phát triển](#nh%c3%b3m-ph%c3%a1t-tri%e1%bb%83n)
  - [Cơ hội nghề nghiệp tại ZaloPay](#c%c6%a1-h%e1%bb%99i-ngh%e1%bb%81-nghi%e1%bb%87p-t%e1%ba%a1i-zalopay)
  - [Liên hệ](#li%c3%aan-h%e1%bb%87)
  
  
## Giới thiệu

<div align="center">
	<img src="./images/background-book/ver1.3.0.png">
	<br/>
	<span align="center">
		<i></i>
	</span>
</div>
<br/>

Ngôn ngữ [Golang](https://golang.org/) không còn quá xa lạ trong giới lập trình nữa. Đây là một ngôn ngữ dễ học, các bạn có thể tự học Golang cơ bản ở trang [Go by Example](https://gobyexample.com/). Đa phần các tài liệu về Golang từ cơ bản hay đến nâng cao đều do các nhà lập trình viên nước ngoài biên soạn. Bộ tài liệu [Advanced Go Programming](#Go-Language-Advanced-Programming-Advanced-Go-Programming) được chúng tôi biên soạn hoàn toàn bằng Tiếng Việt sẽ trình bày về những chủ đề nâng cao trong Golang như CGO, RPC framework, Web framework, Distributed systems,... và kèm theo các ví dụ minh họa cụ thể theo từng chủ đề. Chúng tôi rất mong bộ tài liệu này sẽ giúp các bạn lập trình viên có thêm nhiều kiến thức mới và nâng cao kỹ năng lập trình Golang cho bản thân.

## Tại sao chúng tôi thực hiện bộ tài liệu này ?

Chúng tôi thực hiện bộ tài liệu nhằm:

- Tạo ra bộ tài liệu về Go cho nội bộ ZaloPay sử dụng.
- Đây là cơ hội để mọi người biết tới technical stack của ZaloPay.
- Public ra bên ngoài để cộng đồng Golang Việt Nam có bộ tài liệu tiếng Việt do chính người Việt Nam biên soạn. 
- Đồng thời tạo ra sân chơi mới có cơ hội giao lưu mở rộng mối quan hệ với các bạn có cùng đam mê lập trình.

## Đối tượng sử dụng

Tất cả các bạn có đam mê lập trình Golang và đã nắm được cơ bản về lập trình Golang. Ngoài ra, trong bộ tài liệu này chúng tôi cũng có nhắc lại vài điểm cơ bản trong lập trình Golang.

## Tài liệu tham khảo

Bộ tài liệu này được chúng tôi biên soạn dựa trên kinh nghiệm và kiến thức tích luỹ trong quá trình làm việc tại ZaloPay. Đồng thời chúng tôi có tham khảo các tài liệu bên ngoài như: 
 - [Advanced Go Programming](https://github.com/chai2010/advanced-go-programming-book).
 - [Khoá học Distributed Systems của Princeton](https://www.cs.princeton.edu/courses/archive/fall18/cos418/schedule.html).

## Mục lục

Xem mục lục chính của bộ tài liệu [ở đây](./SUMMARY.md).

## Phương thức đọc

- Đọc online: [GitBook](https://zalopay-oss.github.io/go-advanced/).
- Tải file pdf: <a href="./pdf/advanced-go-book.pdf" download>pdf</a>
- Tải file epub: <a href="./epub/advanced-go-book.epub" download>epub</a>
- Tải file mobi: <a href="./mobi/advanced-go-book.mobi" download>mobi</a>

## Tham gia phát triển

Chúng tôi biết tài liệu này còn nhiều hạn chế. Để trở nên hoàn chỉnh hơn trong tương lai, chúng tôi rất vui khi nhận được sự đóng góp từ mọi người.

Các bạn có thể đóng góp bằng cách:

- [Liên hệ](#li%C3%AAn-h%E1%BB%87) với chúng tôi.
- Trả lời các câu hỏi trong [issues](https://github.com/zalopay-oss/go-advanced/issues).
- Tạo các issues gặp phải trên [issues](https://github.com/zalopay-oss/go-advanced/issues).
- Tạo pull request trên repository của chúng tôi.
- ...

## Nhóm phát triển

Dự án này được phát triển bởi các thành viên sau đây. 

| [<img src="https://avatars1.githubusercontent.com/u/38773351?s=460&v=4" width="100px;"/><br /><sub><b>phamtai97</b></sub>](https://github.com/phamtai97) | [<img src="https://avatars1.githubusercontent.com/u/26034284?s=460&v=4" width="100px;"/><br /><sub><b>thinhdang</b></sub>](https://github.com/thinhdanggroup) | [<img src="https://avatars2.githubusercontent.com/u/23535926?s=460&v=4" width="100px;"/><br /><sub><b>quocanh1897</b></sub>](https://github.com/quocanh1897) | [<img src="https://avatars2.githubusercontent.com/u/32214488?s=400&v=4" width="100px;"/><br /><sub><b>thoainguyen</b></sub>](https://github.com/thoainguyen) | [<img src="https://avatars1.githubusercontent.com/u/3270746?s=460&v=4" width="100px;"/><br /><sub><b>anhldbk</b></sub>](https://github.com/anhldbk) |
| :---------------------------------------------------------------------------------------------------------------------------------------------------: | :---------------------------------------------------------------------------------------------------------------------------------------------------------: | :--------------------------------------------------------------------------------------------------------------------------------------------------: | :-------------------------------------------------------------------------------------------------------------------------------------------------------: | :-----------------------------------------------------------------------------------------------------------------------------------------------------------------: |

## Cơ hội nghề nghiệp tại ZaloPay

<div align="center">
	<img src="./images/qc-zalopay.png" width="600">
	<br/>
	<span align="center">
		<i></i>
	</span>
</div>
<br/>

[ZaloPay](https://zalopay.vn/) là một trong những ví điện tử được ưa chuộng hiện nay với nhiều tính năng và tiện ích hấp dẫn, giúp chúng ta giao dịch tài chính nhanh chóng hơn thông qua ứng dụng ZaloPay. Chúng tôi luôn mong muốn có thêm các thành viên mới gia nhập đội ngũ engineering, cùng giải quyết các bài toán hóc búa về high performance, fault tolerant và distributed transaction. `Java`, `Golang` và `Rust` là ngôn ngữ chính của chúng tôi.

## Liên hệ

- Github: [ZaloPay Open Source](https://github.com/zalopay-oss)
  
- Facebook: [ZaloPay Engineering](https://www.facebook.com/zalopay.engineering/)

- Blog: [Medium ZaloPay Engineering](https://medium.com/zalopay-engineering)

- Bugs report: [issues](https://github.com/zalopay-oss/go-advanced/issues)
