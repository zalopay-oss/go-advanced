# Chapter 2: CGO Programming

<<<<<<< HEAD
<div align="center">
	<img src="../images/ch2.png">
	<br/>
	<span align="center">
		<i></i>
	</span>
</div>
<br/>

=======
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376
Sau nhiều thập kỷ phát triển, với số lượng lớn phần mềm được viết bằng C/C++,, nhiều trong số đó được kiểm thử và tối ưu hiệu năng. Ngôn ngữ Go nên tận dụng ưu thế to lớn đó của C/C++.

Vì C là ngôn ngữ cơ bản, nhiều thư viện sẽ chọn cung cấp API tương thích với C, sau đó được hiện thực lại bằng một ngôn ngữ khác.

Go hỗ trợ các lời gọi hàm từ C thông qua một công cụ gọi là `CGO`, do đó ta có thể sử dụng Go để đưa thư viện động của C (C dynamic libraby) sang các ngôn ngữ khác. Chương này sẽ đi sâu vào tìm hiểu các vấn đề liên quan đến lập trình với CGO.
