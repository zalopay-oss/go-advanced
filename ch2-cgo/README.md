# Chapter 2: CGO Programming

<div align="center">
	<img src="../images/ch2.png">
	<br/>
	<span align="center">
		<i></i>
	</span>
</div>
<br/>

>*“I think it goes back to the Unix C traditions back to basics and other compiled languages and it remedies other deficiencies in C, I don’t think C++ was an improvement but I do think Go is a definite improvement on C and we’ve got Kernighan and things in the background there and obviously they’ve got wonderful experience on building languages. It’s very nicely engineered and actually when it even came out impressive documentation, and all this stuff that you need. Even when it first came out it has a level of maturity that you would think would actually have been there for many years, so it is very impressive actually.” – Joe Armstrong, co-inventor of Erlang*

Sau nhiều thập kỷ phát triển, với số lượng lớn phần mềm được viết bằng C/C++, nhiều trong số đó được kiểm thử và tối ưu hiệu năng. Ngôn ngữ Go nên tận dụng ưu thế to lớn đó của C/C++.

Go hỗ trợ các lời gọi hàm từ C thông qua một công cụ gọi là `CGO`, do đó ta có thể sử dụng Go để đưa thư viện động của C (C dynamic libraby) sang các ngôn ngữ khác. Chương này sẽ đi sâu vào tìm hiểu các vấn đề liên quan đến lập trình với CGO.

Tuy nhiên ta cũng không nên lạm dụng CGO vì vấn đề hiệu suất, ví dụ nếu so sánh với Rust: `rustgo` chỉ chậm hơn gọi hàm Go trực tiếp khoảng 11%, nhưng nó lại nhanh hơn gấp 15 lần `CGO`, một vài con số cụ thể từ [link sau](https://blog.filippo.io/rustgo/):

name|        time/op
--- | ---
CallOverhead/Inline | 1.67ns ± 2%
CallOverhead/Go     | 4.49ns ± 3%
CallOverhead/rustgo | 4.58ns ± 3%
CallOverhead/cgo    | 69.4ns ± 0%
