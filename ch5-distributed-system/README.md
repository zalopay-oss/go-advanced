# Chapter 5 Hệ thống phân tán

<div align="center">
	<img src="../images/ch6.png">
	<br/>
	<span align="center">
		<i></i>
	</span>
</div>
<br/>

>*“We used it to write our own simple distributed computing software after realizing hadoop was too complicated (and thus bug prone) for our embarrassingly parallel needs. It took us less time to get the system written, stable and up and running then it had to get hadoop setup” – micro_cam in Hacker News*

Ngôn ngữ Go ngày càng trở nên phổ biển và đang dần thay thế các ngôn ngữ truyền thống vì các ưu điểm vượt trội của nó. Bên cạnh đó, sự phát triển mạnh mẽ của điện toán đám mây như AWS, Azure,... đã mang lại nhiều lợi ích cho các doanh nghiệp. Và không thể không kể đến các hệ thống như Docker, Kubernetes được xây dựng bằng Go, nhờ chúng mà kỷ nguyên đám mây đã phát triển mạnh mẽ và nhanh chóng. Cùng với đó, nó lại kéo theo là các mô hình thiết kế hiện đại ra đời như serverless, microservices, ..., nơi mà phần cứng đã đạt giới hạn của nó, ta không thể tiếp tục mở rộng theo kiểu `vertical` mà thay vào đó phải tập trung theo `horizontal` hay chia nhỏ vấn đề lớn thành nhiều vấn đề nhỏ để giải quyết. Lúc này, ta sẽ gặp lại các câu hỏi quen thuộc trong hệ thống phân tán như:

- Làm sao tạo một ID duy nhất trong hệ thống phân tán?
- Làm sao tạo ra một `lock` phân tán trên nhiều hệ thống?
- Làm sao cân bằng tải trong hệ thống?
- Tính thống nhất khi nhiều hệ thống sử dụng chung cấu hình?
- Thu thập dữ liệu lớn?

Các vấn đề trên tuy không mới nhưng ở đây, chúng ta sẽ giải quyết chúng theo cách của Go.

## Liên kết
* Phần tiếp theo: [Distributed ID generator](./ch5-01-dist-id.md)
* Phần trước: [Chương 4: Lời nói thêm](../ch4-web/ch4-08-ext.md)
* [Mục lục](../SUMMARY.md)