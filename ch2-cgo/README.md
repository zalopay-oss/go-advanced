# Chapter 2: CGO Programming

<div align="center">
	<img src="../images/ch2.png">
	<br/>
	<span align="center">
		<i></i>
	</span>
</div>
<br/>

>*“C makes it easy to shoot yourself in the foot; C++ makes it harder, but when you do it blows your whole leg off (Bjarne Stroustrup)*

Sau nhiều thập kỷ phát triển, với số lượng lớn phần mềm được viết bằng C/C++, nhiều trong số đó được kiểm thử và tối ưu hiệu năng. Ngôn ngữ Go nên tận dụng ưu thế to lớn đó của C/C++.

Go hỗ trợ các lời gọi hàm từ C thông qua một công cụ gọi là `CGO`, do đó ta có thể sử dụng Go để đưa thư viện động của C (C dynamic libraby) sang các ngôn ngữ khác. Chương này sẽ đi sâu vào tìm hiểu các vấn đề liên quan đến lập trình với CGO.

Tuy nhiên ta cũng không nên lạm dụng CGO vì vấn đề hiệu suất, ví dụ nếu so sánh với Rust: `rustgo` chỉ chậm hơn gọi hàm Go trực tiếp khoảng 11%, nhưng nó lại nhanh hơn gấp 15 lần `CGO`, một vài con số cụ thể từ [link sau](https://blog.filippo.io/rustgo/):

<div align="center">
	<img src="../images/ch2-table-per.png" width="350>
	<br/>
	<span align="center">
		<i></i>
	</span>
</div>
<br/>
