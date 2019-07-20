# 5.8 Interface và Table Driven Development

Trong dự án web, bạn sẽ thường bắt gặp sự thay đổi từ môi trường phụ thuộc bên ngoài, như là:

1. Hệ thống cũ dùng để lưu trữ dữ liệu của công ty đã bị hư hỏng trong một thời gian dài và hiện tại không có ai bảo trì nó. Hệ thống mới được xem là không thể chuyển giao trơn tru, những cuối cùng yêu cầu đưa ra là phải chuyển giao trong vòng N ngày.
2. Hệ thống cũ của platform department bị hư hỏng trong thời gian dài, và bây giờ không có ai bảo trì chúng. Đó là một câu chuyện buồn. Hệ thống mới không tương thích với interface cũ, nhưng cuối cùng nó cũng bị sụp đổ, và yêu cầu phải chuyển giao trong vòng N ngày.
3. Hệ thống hàng đợi tin tức của công ty bị hư hỏng. Những công nghệ mới không tương thích với nó, nhưng cuối cùng cũng phải thực hiện và chuyển giao trong vòng nửa năm.

## 5.8.1 Quy trình phát triển hệ thống doanh nghiệp

Miễn là công ty Internet tồn tại trong vòng ba năm, vấn đề chính mà những người kỹ sư phải đối mặt là mã nguồn phình to. Sau khi mã nguồn hệ thống bị lớn lên, những phần của hệ thống có thể không liên quan đến business process của chúng ta có thể được tháo rời và không đồng bộ. Đâu là những phần liên quan đến business, như là thống kê, chống gian lận, tiếp thị, tính toán giá, cập nhật trạng thái user, v,v.. Những yêu cầu đó, thường sẽ phụ thuộc vào dữ liệu trong main process nhưng chúng chỉ là một nhánh bên của main process, và chúng khép kín.

Trong thời gian này, chúng ta có thể tháo rời những nhánh bên và deploy, phát triển và bảo trì chúng trong một hệ thống độc lập. Sự trì hoãn của những quá trình xử lý song song này rất nhạy cảm. Ví dụ, nếu user click vào button trên giao diện, kết quả sẽ trả về ngay lặp tức (tính toán giá, thanh toán), sau đó RPC communication trong main process system được yêu cầu, và khi comunication faileds, kết quả có thể được trực tiếp trả về. Về phía user, nếu việc trì hoãn không quá nhạy cảm, như là hệ thống xổ số, và kết quả được công bố sau đó, hoặc hệ thống thống kê không cần phải theo thời gian thực, thì không cần phải hiện thực một tiến trình RPC cho mỗi system trong main process. Chúng ta chỉ cần package dữ liệu cần trong downstream vào trong một message và chuyển nó tới message queue. Những thứ tiếp theo sẽ không có gì để làm trong main process (dĩ nhiên, quá trình theo dõi người dùng vẫn cần phải được thực hiện).

Mặc dù, một số vấn đề có thể được giải quyết thông qua việc tháo gỡ và không đồng bộ, cũng không thể  giải quyết được tất cả. Trong quá trình phát triển business, những modules trong nguyên lý đơn trách nhiệm sẽ trở nên phức tạp hơn, chúng vẫn là một xu hướng không thể tránh khỏi. Nếu một thứ trở nên phức tạp, sau đó việc gỡ bỏ và không đồng bộ  không hoạt động. Chúng ta vẫn phải làm một số thứ nhất định để đóng gói sự trừu tượng trong bản thân nó.

## 5.8.2 Gói các quá trình business vào Functions

Trong hầu hết các quá trình package cơ bản, chúng ta đặt một số hành vi tương tự cùng nhau, và sau đó package chúng trong một hàm duy nhất, do đó mã nguồn của chúng ta trong rất bừa bộn như sau

```go
func BusinessProcess(ctx context.Context, params Params) (resp, error){
    ValidateLogin()
    ValidateParams()
    AntispamCheck()
    GetPrice()
    CreateOrder()
    UpdateUserStatus()
    NotifyDownstreamSystems()
}
```

Không quan tâm đến độ phức tạp của business, logic trong hệ thống có thể được chia ra thành `step1 -> step2 -> step3 ...` như là một tiến trình.

Sẽ có một số tiến trình trong mỗi bước, như là:

```go
func CreateOrder() {
    ValidateDistrict()
    ValidateVIPProduct()
    GetUserInfo()
    GetProductDesc()
    DecrementStorage()
    CreateOrderSnapshot()
    return CreateSuccess
}
```

Khi đọc business process code, chúng ta cần đọc tên function để biết được chúng làm gì trong tiến trình. Nếu bạn cần phải thay đổi một số chi tiết, và sau đó đi đến mỗi bước business để xem một process cụ thể. Một business process code được viết tốt sẽ đẩy tất cả các processes vào một số hàm, trả về hàng trăm hoặc hàng ngàn dòng functions. Kiểu spaghetti-style code này khi đọc hoặc bảo trì rất kinh khủng. Trong development process, một package đơn giản như trên sẽ được thực thi ngay lập tức nếu đó là một điều kiện.

## 5.8.3 Dùng interfaces để trừu tượng hóa

Trong thời gian đầu của quá trình phát triển doanh nghiệp, không phù hợp để đưa interfaces vào. Trong nhiều trường hợp, business process thay đổi rất nhanh. Việc đưa vào các interfaces quá sớm có thể làm tăng độ phức tạp của hệ thống businesses bằng việc thêm vào các phân tầng không cần thiết, kết quả dẫn đến sự phủ định hoàn toàn của mỗi sửa đổi.

Khi hệ thống business phát triển tới mức độ nhất định, và main process đã ổn định, interface có thể được dùng cho việc trừu tượng hóa. Tính ổn định có nghĩa là hầu hết các bước của business trong main process sẽ phải được xác định. Nếu như những sự thay đổi được tạo ra, thì sẽ không có thay đổi theo quy mô lớn, nhưng chỉ một phần nhỏ được sửa lại, hoặc chỉ thêm hoặc xóa một số bước business.

Nếu chúng ta đã packaged các business step tốt trong suốt quá trình phát triển, nó rất dễ để trừu tượng interface tại thời điểm này. Đây là mã giả:

```go
type OrderCreator interface {
    ValidateDistrict()
    ValidateVIPProduct()
    GetUserInfo()
    GetProductDesc()
    DecrementStorage()
    CreateOrderSnapshot()
}
```

Chúng ta có thể hoàn toàn trừu tượng hóa bằng việc đề cập tới các bước function signatures được viết ở trên.

Trước khi trừu tượng hóa, chúng ta cần phải hiểu rằng, việc giới thiệu interfaces sẽ có ý nghĩa đối với hệ thống, nó sẽ được phân tích theo ngữ cảnh. Nếu hệ thống chỉ phục vụ cho một product line, và mã nguồn bên trong chỉ được thay đổi cho những ngữ cảnh cụ thể, thì việc giới thiệu interface không thực sự mang lại ý nghĩa to lớn. Liệu rằng nó có thuận tiện để test, chúng ta sẽ bàn về chúng trong các phần sau.

Nếu hiện thực một platform system mà nó yêu cầu định nghĩa các uniform business processes và  business specifications, sau đó interface-based abstraction make sense. Ví dụ:

<div align="center">
	<img src="../images/ch5-interface-impl.uml.png">
	<br/>
	<span align="center">
		<i>Implementing a public interface</i>
	</span>
</div>
<br/>

Flatform cần phải phục vụ nhiều business khác nhau, nhưng dữ liệu được định nghĩa cần phải thống nhất. Về phía platform, chúng ta có thể định nghĩa một tập các interfaces tương tự như trên, và sau đó yêu cầu bên business access chúng phải hiện thực lại. Nếu interface có một số bước không mong muốn, chỉ cần trả về `nil`, hoặc phớt lờ chúng.

Khi business lặp đi lặp lại, platform không được thay đổi. Do đó, chúng ta sử dụng các services như là một plugin của platform đó. Điều gì xảy ra nếu chúng ta không có một interface?

```go
import (
    "sample.com/travelorder"
    "sample.com/marketorder"
)

func CreateOrder() {
    switch businessType {
    case TravelBusiness:
        travelorder.CreateOrder()
    case MarketBusiness:
        marketorder.CreateOrderForMarket()
    default:
        return errors.New("not supported business")
    }
}

func ValidateUser() {
    switch businessType {
    case TravelBusiness:
        travelorder.ValidateUserVIP()
    case MarketBusiness:
        marketorder.ValidateUserRegistered()
    default:
        return errors.New("not supported business")
    }
}

// ...
switch ...
switch ...
switch ...
```

Đúng vậy, nó kết thúc với  `switch`. Sau khi giới thiệu về interface, chúng tôi dùng `switch` chỉ để cần thực thi một lần trong business portal.

```go
type BusinessInstance interface {
    ValidateLogin()
    ValidateParams()
    AntispamCheck()
    GetPrice()
    CreateOrder()
    UpdateUserStatus()
    NotifyDownstreamSystems()
}

func entry() {
    var bi BusinessInstance
    switch businessType {
        case TravelBusiness:
            bi = travelorder.New()
        case MarketBusiness:
            bi = marketorder.New()
        default:
            return errors.New("not supported business")
    }
}

func BusinessProcess(bi BusinessInstance) {
    bi.ValidateLogin()
    bi.ValidateParams()
    bi.AntispamCheck()
    bi.GetPrice()
    bi.CreateOrder()
    bi.UpdateUserStatus()
    bi.NotifyDownstreamSystems()
}
```

Chương trình Interface-oriented, sẽ không quan tâm về việc hiện thực cụ thể. Nếu những service tương ứng được thay đổi, tất cả các logic sẽ hoàn toàn minh bạch ở phía platform.

## 5.8.4 Điểm mạnh và yếu của interface

Interface design được thường xuyên sử dụng trong ngôn ngữ Go. Modules không cần phải biết đến sự xuất hiện của những modules khác. Module A định nghĩa một interface và module B có thể hiện thực interface đó. Nếu không có kiểu dữ liệu được định nghĩa module A trong interface, thì sau đó module B sẽ không cần phải dùng `import A`. Ví dụ, trong thư viện chuẩn `io.Writer` :

```go
type Writer interface {
    Write(p []byte) (n int, err error)
}
```

Chúng tôi cần phải hiện thực interface `io.Writer` trong module của chúng ta:

```go
type MyType struct {}

func (m MyType) Write(p []byte) (n int, err error) {
    return 0, nil
}
```

Sau đó chúng ta truyền `MyType` vào hàm `io.Writer` mà nó được dùng như là một parameter, như là:

```go
package log

func SetOutput(w io.Writer) {
    output = w
}
```

Sau đó:

```go
package my-business

import "xy.com/log"

func init() {
    log.SetOutput(MyType)
}
```

Trong việc định nghĩa `MyType`, không cần phải `import "io"` để trực tiếp hiện thực `io.Writer` interface, chúng ta có thể kết hợp nhiều hàm để hiện thực các interfaces, trong khi phía interface không có thiết lập import các dependency được sinh ra. Do đó, nhiều người nghĩ rằng orthogonality của Go rất tốt để thiết kế.

Mặc dù sự thuận tiện, lợi ích mang lại bởi interface là hiển nhiên. Đầu tiên, dựa vào inversion, cái ảnh hưởng đến interface trên dự án phần mềm trong hầu hết các ngôn ngữ, trong việc thiết kế Go's orthogonal interface. Hoàn toàn có thể loại bỏ tất cả các dependencies; hai là bộ biên dịch sẽ giúp ta kiểm tra lỗi như "not fully implemented interfaces" tại thời điểm biên dịch, nếu business không hiện thực đủ các method trong Interface, nhưng sử dụng laị sử dụng nó.

```go
package main

type OrderCreator interface {
    ValidateUser()
    CreateOrder()
}

type BookOrderCreator struct{}

func (boc BookOrderCreator) ValidateUser() {}

func createOrder(oc OrderCreator) {
    oc.ValidateUser()
    oc.CreateOrder()
}

func main() {
    createOrder(BookOrderCreator{})
}
```

Những lỗi sau có thể được in ra

```
# command-line-arguments
./a.go:18:30: cannot use BookOrderCreator literal (type BookOrderCreator) as type OrderCreator in argument to createOrder:
    BookOrderCreator does not implement OrderCreator (missing CreateOrder method)
```

Do đó, interface có thể được xem như là một cách an toàn để kiểm tra kiểu tại thời điểm biên dịch.

## 5.8.5 Table Driven Development

Nếu trong hàm có sử dụng `if` hoặc `switch` thì sẽ làm phức tạp hơn. Có cách
```go
func entry() {
    var bi BusinessInstance
    switch businessType {
    case TravelBusiness:
        bi = travelorder.New()
    case MarketBusiness:
        bi = marketorder.New()
    default:
        return errors.New("not supported business")
    }
}
```

Có thể được sửa đổi thành:

```go
var businessInstanceMap = map[int]BusinessInstance {
    TravelBusiness : travelorder.New(),
    MarketBusiness : marketorder.New(),
}

func entry() {
    bi := businessInstanceMap[businessType]
}
```

Ở `Table-driven design`, nhiều thiết kế liên quan không dùng nó như một design pattern, nhưng chúng tôi nghĩ nó vẫn có ý nghĩa quan trọng để giúp chúng ta đơn giản mã nguồn.

Dĩ nhiên, `table-driven` không phải là một lựa chọn hoàn hảo, bởi vì bạn cần tính hash từ `key`. Trong trường hợp hiệu suất là quan trọng ta cần phải cân nhắc kĩ lưỡng khi sử dụng.
