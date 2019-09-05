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

Ngôn ngữ Go được biết đến như là ngôn ngữ C của thời đại Internet. Ngày nay, hệ thống Internet không phải là thời đại khi các hệ thống trước đó đã làm mọi thứ rồi. Các dịch vụ trong kỷ nguyên Internet đã và đang sử dụng một số lượng lớn các hệ thống phân tán. Đồng thời, sự trỗi dậy của các nhà cung cấp dịch vụ đám mây như Alibaba Cloud và Tencent Cloud đã đánh dấu sự xuất hiện của kỷ nguyên đám mây. Lập trình phân tán (distributed programming) trong kỷ nguyên đám mây sẽ trở thành một kỹ năng cơ bản. Các hệ thống Docker và K8s được xây dựng trên ngôn ngữ Go đã thúc đẩy sự xuất hiện sớm của kỷ nguyên đám mây này.

Ở chương này đối với các hệ thống phân tán đã được phát triển tốt, chúng tôi sẽ trình bày ngắn gọn về cách sử dụng chúng để cải thiện hiệu quả công việc. Đối với các hệ thống không có giải pháp sẵn có, chúng tôi sẽ đề xuất một giải pháp dựa trên các nhu cầu kinh doanh.
