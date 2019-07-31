# Chương 4: RPC và Protobuf

<div align="center">
	<img src="../images/ch4.png">
	<br/>
	<span align="center">
		<i></i>
	</span>
</div>
<br/>

>Điều gì là quan trọng trong khi học lập trình? Thực hành nhiều, xem nhiều, thực hành nhiều! Sau khi học qua các ngôn ngữ, mastering cú pháp cơ bản và những đặc điểm của ngôn ngữ, trận chiến sẽ bắt đầu, rất nhanh! __khlipeng

RPC là viết tắt của Remote Procedure Call (lời gọi hàm từ xa), nó thường là một hàm nằm ở remote, có thể là một hàm khác nằm trong cùng một file, hoặc có thể là một function nằm trong một process khác trên cùng một máy, hoặc sẽ là một phương thức bí mật trên Sao Hỏa. Bởi vì hàm được gọi trong RPC có thể ở rất xa, xa như chúng ta nói những ngôn ngữ khác nhau, ngôn ngữ dường như là rào cản trong khi giao tiếp giữa hai phía. Protobuf hỗ trợ nhiều ngôn ngữ khác nhau (những ngôn ngữ chưa hỗ trợ cũng sẽ được mở rộng để hỗ trợ), những tính năng của chúng cũng rất thuận tiện để mô tả interface cho service (đó là một danh sách các method), do đó rất phù hợp để có một interface communication language của thế giới RPC. Chương này sẽ bàn về việc sử dụng RPC, thiết kế RPC services của chúng ta như thế nào trong nhiều ngữ cảnh khác nhau, và hệ sinh thái RPC lớn được xây dựng dựa trên Protobuf.

