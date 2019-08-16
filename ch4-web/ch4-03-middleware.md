# 4.3 Middleware

Phần này sẽ phân tích tình huống dẫn tới việc sử dụng middleware, sau đó trình bày cách hiện thực một middleware đơn giản để tách biệt mã nguồn business và non-business. Trước hết chúng ta cùng nhắc lại khái niệm middleware là gì? Có thể nói ngắn gọn middleware là những đoạn mã trung gian nằm ở giữa request và response của ứng dụng web của chúng ta.

<div align="center">
	<img src="../images/ch4-middleware.png">
	<br/>
    <br/>
</div>

Middleware thường được dùng trong một số trường hợp chúng ta muốn ghi log hoạt động của hệ thống, báo cáo thời gian thực thi, xác thực,..

## 4.3.1 Tình huống đặt ra

Hãy nhìn vào đoạn mã nguồn sau:

***main.go:***

```go
package main

func hello(wr http.ResponseWriter, r *http.Request) {
    wr.Write([]byte("hello"))
}

func main() {
    http.HandleFunc("/", hello)
    err := http.ListenAndServe(":8080", nil)
    //...
}
```

Bây giờ có một số nhu cầu mới, chúng tôi muốn tính được thời gian xử lý của Hello service được viết ở trên, nhu cầu này rất đơn giản, chúng tôi sẽ làm một số thay đổi nhỏ trên chương trình ở trên.

```go
var logger = log.New(os.Stdout, "", 0)

func hello(wr http.ResponseWriter, r *http.Request) {
    timeStart := time.Now()
    wr.Write([]byte("hello"))
    timeElapsed := time.Since(timeStart)
    logger.Println(timeElapsed)
}
```

Đoạn mã nguồn thêm vào ở trên đã giải quyết được yêu cầu đặt ra, tuy nhiên trong quá trình phát triển, số lượng API ngày một tăng lên như sau:

***main.go:***

```go
package main

func helloHandler(wr http.ResponseWriter, r *http.Request) {
    // ...
}

func showInfoHandler(wr http.ResponseWriter, r *http.Request) {
    // ...
}

func showEmailHandler(wr http.ResponseWriter, r *http.Request) {
    // ...
}

func showFriendsHandler(wr http.ResponseWriter, r *http.Request) {
    timeStart := time.Now()
    wr.Write([]byte("your friends is tom and alex"))
    timeElapsed := time.Since(timeStart)
    logger.Println(timeElapsed)
}

func main() {
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/info/show", showInfoHandler)
    http.HandleFunc("/email/show", showEmailHandler)
    http.HandleFunc("/friends/show", showFriendsHandler)
    // ...
}
```

Mỗi handler có một đoạn mã nguồn để ghi lại thời gian được đề cập từ trước. Mỗi lần chúng tôi thêm vào một route mới, cần phải sao chép những mã nguồn tương tự tới nơi chúng ta ta cần, bởi vì số lượng route ít, nên không phải là vấn đề lớn khi hiện thực.

Bây giờ, chúng ta cần bản báo cáo về dữ liệu thời gian chạy service cho hệ thống metrics. Nên cần phải thay đổi mã nguồn và gửi thời gian đến hệ thống metrics thông qua HTTP POST. Hãy thay đổi nó `helloHandler()`:

***main.go:***

```go
func helloHandler(wr http.ResponseWriter, r *http.Request) {
    timeStart := time.Now()
    wr.Write([]byte("hello"))
    timeElapsed := time.Since(timeStart)
    logger.Println(timeElapsed)
    // Thêm phần upload thời gian
    metrics.Upload("timeHandler", timeElapsed)
}
```

Mỗi khi thay đổi, chúng ta có thể dễ dàng thấy rằng công việc phát triển sẽ rơi vào bế tắc. Bất kể nhu cầu phi chức năng hoặc thống kê trên hệ thống Web trong tương lai, các sửa đổi sẽ ảnh hưởng tới toàn bộ. Cũng như khi ta thêm một nhu cầu thống kê đơn giản, chúng ta cần phải thêm hàng tá những mã nguồn độc lập với business. Mặc dù dường như chúng không có lỗi trong thời gian đầu, nhưng sẽ thấy rõ hơn khi business càng phát triển.

## 4.3.2 Hiện thực middleware

Thực tế, vấn đề là chúng ta gây ra là đặt mã nguồn business và non-business cùng nhau. Trong hầu hết trường hợp, những yêu cầu non-business thường là làm một thứ gì đó trước khi xử lý HTTP request, và làm một thứ gì đó ngay sau khi chúng hoàn thành. Ý tưởng ở đây là tái cấu trúc lại mã nguồn để tách riêng mã nguồn của non-business riêng, như sau:

***main.go (version2):***

```go
// hàm business logic của chúng ta
func hello(wr http.ResponseWriter, r *http.Request) {
    wr.Write([]byte("hello"))
}
// đây là một hàm thực hiện việc đóng gói hàm truyền vào
// để ghi nhận thời gian thực thi service
// hàm này được xem như là một Middleware
func timeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
        // ghi nhận thời gian trước khi chạy
        timeStart := time.Now()

        // next là hàm business logic được truyền vào
        next.ServeHTTP(wr, r)
        // tính toán thời gian thực thi
        timeElapsed := time.Since(timeStart)
        // log ra thời gian thực thi
        logger.Println(timeElapsed)
    })
}
// Trong đó, http.Handler:
// type Handler interface {
//    ServeHTTP(ResponseWriter, *Request)
// }
func main() {
    http.Handle("/", timeMiddleware(http.HandlerFunc(hello)))
    err := http.ListenAndServe(":8080", nil)
    // ...
}
```

Bất cứ hàm nào định nghĩa `ServeHTTP`, cũng đều là một đối tượng `http.Handler`. Những gì mà middleware làm là nhận vào hàm handler và trả về hàm handler khác kèm theo non-business logic. Chúng ta có thể dùng các middleware lồng vào nhau như bên dưới:

```go
customizedHandler = logger(timeout(ratelimit(helloHandler)))
```

Ngữ cảnh của chuỗi các hàm trong quá trình thực thi có thể được thể hiện dưới đây:

<div align="center">
	<img src="../images/ch5-03-middleware_flow.png" width="600">
	<br/>
    <br/>
</div>

Tuy nhiên cách dùng middleware lồng vào nhau như trên còn khá phức tạp.

## 4.3.3 Cách viết middleware thanh lịch hơn

Trong phần trước, sự tách biệt về mã nguồn hàm business và non-business function được giải quyết. Nhưng nếu bạn cần phải thay đổi thứ tự của những hàm đó, hoặc thêm, hoặc xóa middleware vẫn còn một số khó khăn, phần này chúng ta sẽ thực hiện việc tối ưu như sau.

***Ví dụ:***

```go
r = NewRouter()
r.Use(logger)
r.Use(timeout)
r.Use(ratelimit)
r.Add("/", helloHandler)
```

Qua nhiều bước thiết lập, chúng ta có một chuỗi thực thi các hàm tương tự như ví dụ trước. Cách làm này giúp chúng ta dễ hiểu hơn. Nếu bạn muốn thêm hoặc xóa middleware, đơn giản thêm và xóa dòng ứng với lời gọi `Use()`.

Từ góc nhìn về framework, làm sao để viết được hàm như vậy?

***Hiện thực:***

```go
// định nghĩa kiểu interface
type middleware func(http.Handler) http.Handler
// cấu trúc Router
type Router struct {
    // slice gồm các hàm middleware
    middlewareChain [] middleware
    // mapping cấu trúc routing với name
    mux map[string] http.Handler
}
func NewRouter() *Router{
    return &Router{}
}
// mỗi khi gọi Use là thêm hàm middleware vào slice
func (r *Router) Use(m middleware) {
    r.middlewareChain = append(r.middlewareChain, m)
}
// mỗi khi gọi Add là thêm phần routing trong đó, áp dụng các middleware vào
func (r *Router) Add(route string, h http.Handler) {
    var mergedHandler = h
    // duyệt theo thứ tự ngược lại để apply middleware
    for i := len(r.middlewareChain) - 1; i >= 0; i-- {
        mergedHandler = r.middlewareChain[i](mergedHandler)
    }
    // cuối cùng register hàm handler vào name route tương ứng
    r.mux[route] = mergedHandler
}
```
