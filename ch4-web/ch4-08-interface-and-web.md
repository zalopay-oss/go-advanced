# 4.8 Interface và Table Driven Development 

Trong các dự án web thực tế chúng ta thường phải thay đổi mã nguồn (thêm, loại bớt) do các yếu tố bên ngoài, như là:

1. Hệ thống cũ dùng để lưu trữ dữ liệu của công ty đã bị hư hỏng trong một thời gian dài và hiện tại không có ai bảo trì nó. Hệ thống mới được xem là không thể chuyển giao trơn tru, những cuối cùng yêu cầu đưa ra là phải chuyển giao trong vòng N ngày.
2. Hệ thống cũ của platform department bị hư hỏng trong thời gian dài, và bây giờ không có ai bảo trì chúng. Đó là một câu chuyện buồn. Hệ thống mới không tương thích với interface cũ, nhưng cuối cùng nó cũng bị sụp đổ, và yêu cầu phải chuyển giao trong vòng N ngày.
3. Hệ thống hàng đợi tin tức của công ty bị hư hỏng. Những công nghệ mới không tương thích với nó, nhưng cuối cùng cũng phải thực hiện và chuyển giao trong vòng nửa năm.

## 4.8.1 Quy trình phát triển hệ thống doanh nghiệp

Các công ty Internet tồn tại trong vòng khoảng ba năm thì cácn mã nguồn của hệ thống dần phình to và gây khó khăn cho các kỹ sư lập trình. Sau khi mã nguồn hệ thống bị lớn lên, có một số phần của hệ thống có thể được tách rời thành các service nhỏ hơn. Các service được tách rời giúp chúng ta dễ dàng deploy, phát triển và bảo trì chúng.

Mặc dù, một số vấn đề có thể được giải quyết thông qua việc tách rời service, cũng không thể  giải quyết được tất cả. Trong quá trình phát triển business, những service này cũng dần trở nên phức tạp hơn, chúng vẫn là một xu hướng không thể tránh khỏi. Vậy nên cách tốt nhất là chúng ta sẽ sử dụng interface khi lập trình để tách rời sự phụ thuộc giữa các thành phần trong mã nguồn cũng như giúp chúng ta dễ dàng mở rộng chúng.

## 4.8.2 Đóng gói các business vào functions

Trong hầu hết các package cơ bản, chúng ta đặt một số hành vi xử ly logic tương tự cùng nhau, và sau đó đóng gói chúng trong một hàm duy nhất, do đó mã nguồn của chúng ta trong rất bừa bộn như sau:

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

Khi đọc business process code, chúng ta cần đọc tên function để biết được chức năng xử lý của chúng. Nếu chúng ta cần phải thay đổi một số chi tiết, chúng ta sẽ vào các function đó và thêm/sửa/xoá các dòng code. Kiểu [spaghetti-style code](https://en.wikipedia.org/wiki/Spaghetti_code) này khi đọc hoặc bảo trì rất khó.

## 4.8.3 Dùng interfaces để trừu tượng hóa

Trong thời gian đầu của quá trình phát triển hệ thống doanh nghiệp, không phù hợp để sử dụng interfaces. Trong nhiều trường hợp, khi business process thay đổi rất nhanh, việc sử dụng các  interfaces quá sớm có thể làm tăng độ phức tạp của hệ thống. Khi hệ thống phát triển tới mức độ nhất định, và có một business ổn định, đây là thời điểm tốt để áp dụng interface vào mã nguồn hệ thống.

Nếu chúng ta đã đóng gói các business step tốt trong suốt quá trình phát triển, nó rất dễ để áp dụng  interface tại thời điểm này. Đây là mã giả:

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

Trước khi trừu tượng hóa, chúng ta cần phải hiểu rằng, việc áp dụng interfaces sẽ có ý nghĩa đối với hệ thống tuỳ theo ngữ cảnh. Nếu hệ thống xác định có một business cố định và mã nguồn bên trong không có sự thay đổi thường xuyên, thì việc áp dụng  interface không thực sự mang lại ý nghĩa to lớn.

Nếu hiện thực một platform system mà nó yêu cầu định nghĩa các business sau:

<div align="center">
	<img src="../images/ch5-interface-impl.uml.png">
	<br/>
	<span align="center">
		<i>Implementing a public interface</i>
	</span>
</div>
<br/>

Flatform cần phải phục vụ nhiều business khác nhau, nhưng dữ liệu được định nghĩa cần phải thống nhất. Về phía platform, chúng ta có thể định nghĩa một tập các interfaces tương tự như trên, sau đó tuỳ theo yêu cầu của các business cụ thể của chúng cần hiện thực lại. Nếu interface có một số bước không mong muốn, chỉ cần trả về `nil`, hoặc có thể bỏ qua chúng.

Điều gì xảy ra nếu chúng ta không có một interface?

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

Chúng ta phải sử dụng   `switch-case` rất nhiều. Sau khi áp dụng interface, chúng ta dùng `switch-case` chỉ để xác định loại business nào cần thực hiện.

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

Chương trình trên sẽ dễ mở rộng và minh bạch hơn rất nhiều. Hàm `BusinessProcess` sẽ không quan tâm đầu vào là loại business nào. Các business khácn nhau của chúng ta chỉ cần hiện thực các chức năng trong interface ban đầu.

## 4.8.4 Điểm mạnh và yếu của interface

Thiết kế interface được sử dụng thường xuyên trong ngôn ngữ Go. Ví dụ, trong thư viện chuẩn `io.Writer` :

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

Sau đó khi sử dụng chúng ta chỉ cần truyền biến có kiểu `MyType`:

```go
package my-business

import "xy.com/log"

func init() {
    log.SetOutput(MyType)
}
```

Trong việc định nghĩa `MyType`, không cần phải `import "io"` để trực tiếp hiện thực `io.Writer` interface, chúng ta có thể kết hợp nhiều hàm để hiện thực các interfaces.

Mặc dù sự thuận tiện, lợi ích mang lại bởi interface là hiển nhiên. Đầu tiên, chúng ta có thể hoàn toàn loại bỏ tất cả các phụ thuộc lẫn nhau trong mã nguồn. Thứ hai là khi biên dịch sẽ giúp ta kiểm tra lỗi như **not fully implemented interfaces** tại thời điểm biên dịch, nếu chúng ta không hiện thực đủ các hàm trong interface, nhưng lại sử dụng nó. Ví dụ như trong trường hợp này:

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

Những lỗi sau có thể được in ra:

```
# command-line-arguments
./a.go:18:30: cannot use BookOrderCreator literal (type BookOrderCreator) as type OrderCreator in argument to createOrder:
    BookOrderCreator does not implement OrderCreator (missing CreateOrder method)
```

Do đó, interface có thể được xem như là một cách an toàn để kiểm tra kiểu tại thời điểm biên dịch.

## 4.8.5 Table Driven Development

Nếu trong hàm nếu chúng ta có sử dụng `if` hoặc `switch` thì sẽ làm mã nguồn trông phức tạp hơn. Ví dụ:

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

Chúng ta có thể được sửa đổi thành:

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
