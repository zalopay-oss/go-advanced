# 4.7. Bố cục thông thường của các dự án web lớn

Phần này sẽ trình bày mô hình MVC và đi vào chi tiết các lớp trong một dự án web.

## 4.7.1 Kiến trúc MVC

MVC frameworks là những framework rất phổ biến trong việc phát triển web, khái niệm MVC được đề xuất đầu tiên bởi [Trygve Reenskaug](https://en.wikipedia.org/wiki/Trygve_Reenskaug) vào năm 1978, chương trình MVC gồm ba thành phần:

* **Model** : là nơi chứa những nghiệp vụ tương tác với dữ liệu hoặc hệ quản trị cơ sở dữ liệu (mysql, mssql… ); nó sẽ bao gồm các class/function xử lý nhiều nghiệp vụ như kết nối database, truy vấn dữ liệu, thêm – xóa – sửa dữ liệu...
* **View** : là nơi có những giao diện như một nút bấm, khung nhập, menu, hình ảnh… nó đảm nhiệm nhiệm vụ hiển thị dữ liệu và giúp người dùng tương tác với hệ thống.
* **Controller** : là nơi tiếp nhận những yêu cầu xử lý được gửi từ người dùng, nó sẽ gồm những class/ function xử lý nhiều nghiệp vụ logic giúp lấy đúng dữ liệu thông tin cần thiết nhờ các nghiệp vụ lớp Model cung cấp và hiển thị dữ liệu đó ra cho người dùng nhờ lớp View.

Trải qua quá trình phát triển, phần back-end của chương trình ngày càng phức tạp. Để quản lý tốt hơn, những phần như thế sẽ thường phân chia ra thành nhiều kiến trúc con. Hình sau là một lưu đồ của hệ thống từ front-end tới back-end:

<div align="center">
	<img src="../images/ch5-07-frontend-backend.png" width="600">
	<br/>
	<span align="center">
		<i>Kiến trúc một dự án web</i>
	</span>
</div>
<br/>

**Vue** và **React** trong hình là hai frameworks front-end phổ biến trên thế giới, bởi vì chúng ta không tập trung nói về nó, do đó cấu trúc front-end của dự án không được nhấn mạnh trên lưu đồ. Thực tế trong vài dự án đơn giản, ngành công nghiệp không hoàn toàn tuân theo mô hình MVC, đặc biệt là phần M và C. Có nhiều công ty mà dự án của họ có rất nhiều phần logic bên trong lớp Controller, và chỉ quản lý phần lưu trữ dữ liệu ở lớp Model.

Tuy nhiên, theo như ý tưởng của người sáng lập MVC thì một business process cũng thuộc một loại "model". Nếu chúng ta đặt mã nguồn thao tác với dữ liệu và business process vào lớp M của MVC, thì lớp M sẽ quá cồng kềnh. Trong những dự án phức tạp, một lớp C hoặc M hiển nhiên là không đủ mà phải có nhiều lớp pure back-end API bên dưới nữa:

## 4.7.2 Bên dưới Controller và Model

Các lớp pure back-end API bên dưới có thể phân chia như sau:

1. **Controller** tương tự như ở trên, là một điểm đầu vào của service, chịu trách nhiệm để xử lý logic routing, kiểm tra tham số, chuyển tiếp request.
2. **Logic/Service**  là lớp logical (service), nó thường là một điểm vào của business logic. Có thể xem rằng tất cả những tham số request sẽ phải được hợp lệ từ đây, Business logic và business processes cũng nằm trong lớp này. Nó thường được gọi là Business Rules trong những thiết kế thường thấy.
3. **DAO/Responsitory** lớp này thường có vai trò chính là thao tác với data (dữ liệu bền vững) và storage (vùng nhớ).

Mỗi lớp sẽ thực thi công việc của nó, sau đó xây dựng lên cấu trúc các phần parameters để truyền cho các lớp kế tiếp bằng việc tạo request từ context hiện tại, và sau đó gọi hàm để thực thi lớp tiếp theo. Sau khi công việc hoàn thành, kết quả của quá trình sẽ được trả về lớp ban đầu gọi nó.

<div align="center">
	<img src="../images/ch5-07-controller-logic-dao.png" width="800">
	<br/>
	<span align="center">
		<i>Flow xử lý request</i>
	</span>
</div>
<br/>

Sau khi chia ra ba lớp CLD (Controller-Logic-DAO), chúng ta cần phải hỗ trợ nhiều giao thức trong lớp Controller như Thrift, gRPC hoặc HTTP đóng vai trò như là một interface. Quá trình xử lý request có thể như sơ đồ sau:

<div align="center">
	<img src="../images/ch4-7-cld-layout.png" width="350">
	<br/>
	<span align="center"><i></i></span>
	<br/>
</div>

Tuy nhiên, bên cạnh các yêu cầu business logic, hệ thống web còn phải hiện thực phần non-business logic như xác thực, gửi số liệu báo cáo, để tách riêng hai phần này, các lớp middleware của từng giao thức được thêm vào, do đó sơ đồ cuối cùng có thể như sau:

<div align="center">
	<img src="../images/ch4-7-cld-layout-md.png">
	<br/>
	<span align="center">
		<i>Flow sau khi thêm phần middleware</i>
	</span>
</div>
<br/>
