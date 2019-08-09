# Chương 4: RPC và Protobuf

<div align="center">
	<img src="../images/ch4.png">
	<br/>
	<span align="center">
		<i></i>
	</span>
</div>
<br/>

>*“One of the reasons I enjoy working with Go is that I can mostly hold the spec in my head - and when I do misremember parts it’s a few seconds' work to correct myself. It’s quite possibly the only non-trivial language I’ve worked with where this is the case.” – Eleanor McHugh*

RPC là viết tắt của Remote Procedure Call (lời gọi hàm từ xa), nó thường là một hàm nằm ở remote, có thể là một hàm khác nằm trong cùng một file, hoặc có thể là một function nằm trong một process khác trên cùng một máy, hoặc sẽ là một phương thức bí mật trên Sao Hỏa. Chương này sẽ bàn về việc sử dụng RPC, thiết kế RPC services của chúng ta như thế nào trong nhiều ngữ cảnh khác nhau, và hệ sinh thái RPC lớn được xây dựng dựa trên Protobuf.

