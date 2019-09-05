# Chương 4: Remote Procedure Call

<div align="center">
	<img src="../images/ch4.png">
	<br/>
	<span align="center">
		<i></i>
	</span>
</div>
<br/>

>*“Go is not meant to innovate programming theory. It’s meant to innovate programming practice.” – Samuel Tesla*

RPC - Remote Procedure Call (lời gọi hàm từ xa), là một kỹ thuật cho phép chúng ta gọi hàm nằm trong một process khác trên cùng một máy hoặc ở hai máy khác nhau. Mục tiêu chính của phương pháp này là giúp lời gọi RPC tương tự như lời gọi thủ tục thông thường và ẩn đi việc truyền dữ liệu đi về qua mạng. Chương này sẽ trình bày về cách sử dụng RPC, thiết kế RPC service, và hệ sinh thái RPC được xây dựng dựa trên nền tảng Protobuf của Google.
