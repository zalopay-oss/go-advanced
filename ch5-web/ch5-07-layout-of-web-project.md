# 5.7 Layout Common large web project layering (bố cục thông thường của các dự án web lớn)

MVC frameworks là những frameworks rất phổ biến trong việc phát triển web. Khái niệm MVC được đề xuất đầu tiên bởi **Trygve Reenskaug** vào năm 1978. Để thuận tiện cho việc mở rộng ứng dụng GUI (graphical user interface), chương trình được chia thành:

1. **Controller** - Có vai trò tiếp nhận và xử lý những requests.
2. **View** - Giao diện đồ họa được thiết kế để tương tác với người dùng.
3. **Model** - Programmer viết các hàm mà chương trình cần phải có (hiện thực thuật toán, v,v), quản lý cơ sở dữ liệu (thêm, xóa, sửa, truy vấn, v,v), thiết kế cơ sở dữ liệu.

Trải qua quá trình phát triển, phần front-end của chương trình ngày càng phức tạp. Để phần kỹ thuật tốt hơn, những phần như thế sẽ thường phân chia ra thành nhiều kiến trúc con. Có thể nhìn thấy rằng, trước và sau khi phân chia lớp V (view) từ mô hình MVC thành các thành phần, một back-end project thường chỉ có lớp C và M. Phần front và back sẽ tương tác lẫn nhau thông qua ajax. Thỉnh thoảng, ta cần giải quyết vấn đề cross-domain, và đã có những giải pháp sẵn rồi. Hình 5.13 sẽ là một lưu đồ của hệ thống từ front tới back.

![](../images/ch5-08-frontend-backend.png)

*Hình 5-13  Separation interaction diagram*

**Vue** và **React** trong hình là hai frameworks front-end phổ biến trên thế giới, bởi vì chúng ta không tập trung nói về nó, do đó, cấu trúc front-end của project không được nhấn mạnh trên lưu đồ. Thực tế trong vài projects đơn giản, ngành công nghiệp không hoàn toàn tuân theo mô hình MVC, đặc biệt là phần M và C. Có nhiều công ty mà project của họ có rất nhiều phần logic bên trong lớp Controller, và chỉ quản lý phần lưu trữ dữ liệu ở lớp Model. Điều đó thường dẫn đến việc hiểu sai ý nghĩa của lớp Model. Về nghĩa đen, lớp này sẽ được đối xử với một vài modeling, và cái gì là Model? nó là dữ liệu!

Cách hiểu này hiển nhiên có vấn đề. Một business process cũng thuộc một loại "model". Nó là một model của hành vi người dùng trong thế giới thực hoặc là những quá trình đã tồn tại. Nó không chỉ là cách tổ chức dữ liệu được định dạng mà gọi là model. Tuy nhiên, theo như ý tưởng của người sáng lập MVC, nếu chúng ta đặt mã nguồn thao tác với dữ liệu và business projects vào lớp M của MVC, thì lớp M sẽ quá cồng kềnh. Cho những projects phức tạp, một lớp C hoặc M hiển nhiên là không đủ. Có nhiều phần pure back-end API thường dùng những phương pháp phân chia sau: 

1. **Controller** tương tự như ở trên, là một điểm đầu vào của service, chịu trách nhiệm để xử lý logic routing, kiểm tra tham số, chuyển tiếp request.
2. **Logic/Service**  là lớp logical (service), nó thường là một điểm vào của business logic. Có thể xem rằng tất cả những tham số request sẽ phải được hợp lệ từ đây, Business logic và business processes cũng nằm trong lớp này. Nó thường được gọi là Business Rules trong những thiết kế thường thấy.
3. **DAO/Responsitory**, lớp này thường có vai trò chính để thao tác với data (dữ liệu) và storage (vùng nhớ). Về cơ bản phần storage được gửi đến lớp Logic để dùng trong các hàm đơn giản, từ interface. Có trách nhiệm nhất quán dữ liệu.

Mỗi lớp sẽ thực thi công việc của nó, sau đó xây dựng lên cấu trúc của các phần parameters để truyền cho các lớp kế tiếp bằng việc tạo request từ context hiện tại
, và sau đó gọi hàm để thực thi lớp tiếp theo. Sau khi công việc hoàn thành, kết quả của quá trình sẽ được trả về lớp đầu vào gọi nó.

![](../images/ch5-08-controller-logic-dao.png)

*Hình 5-14  Request processing flow*

Sau khi chia ra ba lớp của CLD, chúng ta cần phải hỗ trợ nhiều giao thức tại cùng một lúc trong lớp C.  Thrift, gRPC và http được đề cập từ những chương trước, và chúng ta chỉ cần một trong số đó để đảm nhận công việc này. Thỉnh thoảng, chúng ta cần hỗ trợ hai trong số chúng, như là cùng một interface. Chúng ta cần cả hai efficient thrift và http hooks cho việc debugging. Do đó, trong CLD, một số lớp giao thức được phân tách có nhiệm vụ xử lý chi tiết trong nhiều giao thức tương tác đa dạng. Quá trình xử lý requesting sẽ như hình 5.15

![](../images/ch5-08-control-flow.png)

Do đó, entry function trong Controller sẽ như sau

```go
func CreateOrder(ctx context.Context, req *CreateOrderStruct) (
    *CreateOrderRespStruct, error,
) {
    // ...
}
```

CreateOrder có hai parameters. ctx được dùng để truyền vào tham số toàn cục ví dụ `trace_id`, nó yêu cầu một chuỗi các request. Req lưu những thông tin về input mà chúng ta cần để tạo ra một order. Kết quả trả về là cấu trúc response và một error. Có thể nói rằng sau khi mã nguồn trên thực thi trên lớp Controller, sẽ không có mã nguồn nào liên kết với "protocol". Bạn không thể tìm nó ở đây, không thể tìm thấy `http.Request`, hoặc `http.ResponseWriter` hoặc bất cứ gì liên quan đến thrift hoặc gRPC.

Tại lớp protocol, một mã nguồn sẽ thao tác với giao thức http như sau

```go
// defined in protocol layer
type CreateOrderRequest struct {
    OrderID int64 `json:"order_id"`
    // ...
}

// defined in controller
type CreateOrderParams struct {
    OrderID int64
}

func HTTPCreateOrderHandler(wr http.ResponseWriter, r *http.Request) {
    var req CreateOrderRequest
    var params CreateOrderParams
    ctx := context.TODO()
    // bind data to req
    bind(r, &req)
    // map protocol binded to protocol-independent
    map(req, params)
    logicResp,err := controller.CreateOrder(ctx, &params)
    if err != nil {}
    // ...
}
```

Theo giả thuyết, chúng ta có thể dùng cùng một cấu trúc request để kết hợp nhiều tags để đạt được cấu trúc tái sử dụng cho nhiều giao thức. Không may, trong thrift, một cấu trúc request sẽ tự động được sinh ra từ IDL. Nội dụng ở trong sẽ tự động sinh ra trong file `ttypes.go`. Về mặt cấu trúc, gRPC cũng tương tự. Một phần của mã nguồn vẫn cần thiết.

Người đọc thông minh có thể nhìn thấy rằng, chi tiết của việc thao tác với protocol thực sự là một quá trình lặp đi lặp lại công việc. Việc xử lý mỗi interface trong lớp protocol không gì hơn là đọc dữ liệu từ một cấu trúc protocol cụ thể (ví dụ `http.Request`, thrift wrapped out) và kết hợp chúng với cấu trúc protocol độc lập của chúng ta, và map cấu trúc này tới Controller entry. Mã nguồn sẽ thực sự trông giống nau. Ngoại trừ tuân theo một mẫu nào đó, sau đó có thể đơn giản tính trừu tượng của các mẫu và dùng hàm để sinh ra mã nguồn để trích xuất những giao thức xử lý phức tạp từ nội dung công việc.

Hãy nhìn vào cấu trúc của HTTP, cấu trúc này sẽ tương ứng với thrift, và một cấu trúc protocol độc lập khác của chúng ta.

```go
// http 请求结构体
type CreateOrder struct {
    OrderID   int64  `json:"order_id" validate:"required"`
    UserID    int64  `json:"user_id" validate:"required"`
    ProductID int    `json:"prod_id" validate:"required"`
    Addr      string `json:"addr" validate:"required"`
}

// thrift 请求结构体
type FeatureSetParams struct {
    DriverID  int64  `thrift:"driverID,1,required"`
    OrderID   int64  `thrift:"OrderID,2,required"`
    UserID    int64  `thrift:"UserID,3,required"`
    ProductID int    `thrift:"ProductID,4,required"`
    Addr      string `thrift:"Addr,5,required"`
}

// controller input struct
type CreateOrderParams struct {
    OrderID int64
    UserID int64
    ProductID int
    Addr string
}
```

Để sinh ra HTTP và thrift entry code, chúng ta cần thông qua một cấu trúc mã nguồn. Nhìn vào ba cấu trúc được định nghĩa ở trên, thực tế, chúng ta có thể dùng cấu trúc để sinh ra IDL của thrift, và "IDL của HTTP service(cấu trúc định nghĩa với json hoặc form related tages)". Cấu trúc ban đầu này có thể được đặt thêm vào HTTP tags và thrift tags.


```go
type FeatureSetParams struct {
    DriverID  int64  `thrift:"driverID,1,required" json:"driver_id"`
    OrderID   int64  `thrift:"OrderID,2,required" json:"order_id"`
    UserID    int64  `thrift:"UserID,3,required" json:"user_id"`
    ProductID int    `thrift:"ProductID,4,required" json:"prod_id"`
    Addr      string `thrift:"Addr,5,required" json:"addr"`
}
```

Sau đó mã nguồn thrift được sinh ra từ IDL và HTTP requests được sinh ra từ cấu trúc.

![](../examples/ch5-08-code-gen.png)

*Hình 5-16 Creating a project entry through the Go code definition structure*

Đối với phương tiện để tạo, bạn có thể đọc mã nguồn Go trong tệp văn bản thông qua Trình phân tích cú pháp được xây dựng bằng ngôn ngữ Go, sau đó tạo mã đích theo AST hoặc đơn giản là biên dịch cấu trúc nguồn và mã Trình tạo với nhau. Bạn có thể có cấu trúc làm tham số đầu vào cho Trình tạo (sẽ đơn giản hơn).

Dĩ nhiên, ý tưởng này không phải là lựa chọn duy nhất. Chúng ta có thể sinh ra một tập các cấu trúc HTTP interface bằng việc parsing IDL của thrift. Nếu chúng ta làm như vậy, toàn bộ quá trình sẽ như hình bên dưới

![](../images/ch5-08-code-gen-2.png)

*Hình 5-17 can also generate other parts from thrift*

Quy trình này trông có vẻ mượt mà hơn trước, nhưng nếu chúng ta chọn nó để hiện thực, bạn cần phải parse IDL của thrift trước, nó sẽ tương tự với việc Parser sẽ phải viết IDL bằng tay, mặc dù **Antlr** hoặc **peg** có thể giúp bạn. Đơn giản hơn việc viết những Parser, nhưng ở bước "parsing" chúng tôi không muốn giới thiệu quá nhiều, do đó chúng ta có thể thực hiện nó.

Bây giờ, workflow đã được định hình, chúng ta có thể nhận ra làm thế nào để toàn bộ quy trình trở nên thân thiện hơn.

Ví dụ, trong môi trường sinh mã đã được giới thiệu ở chương Web, cũng như user có thể sinh ra SDK với vài cú click chuột, người đọc có thể tự tìm hiểu.


Mặc dù chúng ta đã thành công trong việc cho phép projects của chúng ta hỗ trợ nhiều giao thức tại portal, chúng vẫn có một số vấn đề cần được giải quyết. Việc phân lớp được mô tả trong chương này không dùng middleware để phân lớp project. Nếu chúng ta xem xét middleware, đâu là quá trình requesting? Nhìn vào hình 5-18 bên dưới.

![](../images/ch5-08-control-flow-2.png)

*Hình 5-18 Control flow after adding middleware*

Ở phần middleware trước mà chúng ta đã tìm hiểu, nó liên hệ chặt chẽ đến giao thức HTTP. Không may là không có middleware trong thrift có thể giải quyết những vấn đề về non-functional logic code dupplication problems với HTTP.

Đây là một vấn đề thực sự được bắt gặp trong những dự án doanh nghiệp. Không may, cộng đồng opensource đã không có một giải pháp multi-protocol middleware thuận tiện. Dĩ nhiên, như chúng ta đã nói từ trước, trong nhiều trường hợp, HTTP interface của chúng ta chỉ được dùng cho việc debugging, và không được đề xuất ra bên ngoài. trong trường hợp đó, đó là vấn đề non-functional code cần được giải quyết trong thrift code.

