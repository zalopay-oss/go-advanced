# 5.1 Giới thiệu về Web Development

Bởi vì gói thư viện `net/http` của Golang chỉ hỗ trợ những hàm routing và hàm chức năng cơ bản. Cho nên trong cộng đồng Golang có một ý tưởng phổ biến là viết thêm các API hỗ trợ routing khác ngoài `net/http`. Theo ý kiến của tôi, nếu các project routing của bạn có những đặc điểm sau: URI cố định, và tham số không truyền thông qua URI, thì nên dùng thư viện chuẩn là đủ. Nhưng với những ngữ cảnh phức tạp hơn, thư viện chuẩn `http` vẫn còn một vài điểm yếu. Ví dụ, xét các route sau:
 
```
GET   /card/:id
POST  /card/:id
DELETE /card/:id
GET   /card/:id/name
...
GET   /card/:id/relations
```

Có thể thấy rằng, cũng là đường dẫn có chứa `/card/:id`, nhưng có phương thức khác nhau và nhánh con khác nhau sẽ dẫn đến cách xử lý khác nhau, logic xử lý những đường dẫn trùng tên như vậy thường sẽ phức tạp.

Framework web của Go có thể được chia thành hai thể loại như sau:

1. Router framework
2. MVC class framework

Khi chọn một framework, trong nhiều trường hợp chúng ta sẽ tham khảo những công nghệ mà công ty đang sử dụng. Ví dụ, nếu công ty có nhiều người làm về `PHP`, thì chúng ta nên chọn framework `beego`, nhưng nếu công ty có nhiều lập trình viên `C`, thì hầu hết những suy nghĩ của họ sẽ đơn giản hết sức có thể. Ví dụ, nhiều lập trình viên C trong những công ty lớn sẽ dùng ngôn ngữ C để viết một chương trình `CGI` nhỏ. Họ không thể sẵn sàng để học `MVC` hoặc nhiều framework Web phức tạp khác. Tất cả những gì họ cần là một route đơn giản, mặc dù họ có thể tự xử lý được nhưng chỉ cần một thư viện xử lý giao thức HTTP cơ bản để giúp anh ta tiết kiệm công sức làm việc thủ công.

Gói thư viện `net/http` đã cung cấp những hàm chức năng cơ bản, và viết một `http echo server` chỉ mất khoảng 30 giây.

```go
//brief_intro/echo.go
package main
import (...)

func echo(wr http.ResponseWriter, r *http.Request) {
    msg, err := ioutil.ReadAll(r.Body)
    if err != nil {
        wr.Write([]byte("echo error"))
        return
    }

    writeLen, err := wr.Write(msg)
    if err != nil || writeLen != len(msg) {
        log.Println(err, "write len:", writeLen)
    }
}

func main() {
    http.HandleFunc("/", echo)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
```

[>> mã nguồn](../examples/ch5/ch5.1/brief-intro.go) 

Nếu bạn không thể hoàn thành chương trình trên trong vòng 30 giây, hãy kiểm tra việc bạn gõ phím quá chậm. Đó là ví dụ để minh họa viết một chương trình định tuyến HTTP trong Go sẽ đơn giản như thế nào. Nếu bạn bắt gặp một trường hợp phức tạp hơn, như là một ứng dụng doanh nghiệp cần một tá interfaces, `net/http` sẽ không phù hợp nếu dùng trực tiếp.

Hãy nhìn một dự án Kafka monitoring trong cộng đồng opensource

```go
//Burrow: http_server.go
func NewHttpServer(app *ApplicationContext) (*HttpServer, error) {
    ...
    server.mux.HandleFunc("/", handleDefault)

    server.mux.HandleFunc("/burrow/admin", handleAdmin)

    server.mux.Handle("/v2/kafka", appHandler{server.app, handleClusterList})
    server.mux.Handle("/v2/kafka/", appHandler{server.app, handleKafka})
    server.mux.Handle("/v2/zookeeper", appHandler{server.app, handleClusterList})
    ...
}
```

Hãy đào sâu mã nguồn trên, dự án Kafka monitoring của một công ty nổi tiếng linkedin. Nếu chúng không dùng bất cứ route framework nào và chỉ dùng `net/http`. Nhìn lại mã nguồn trên dường như chúng rất đẹp, chỉ có 5 URIs đơn giản trong dự án của chúng ta, do đó service chúng ta hỗ trợ như sau

```sh
/
/burrow/admin
/v2/kafka
/v2/kafka/
/v2/zookeeper
```

Nếu bạn thực sự nghĩ vậy, bạn đã bị lừa. Hãy xem trong hàm `handleKafka()` được định nghĩa như thế nào

```go
func handleKafka(app *ApplicationContext, w http.ResponseWriter, r *http.Request) (int, string) {
    pathParts := strings.Split(r.URL.Path[1:], "/")
    if _, ok := app.Config.Kafka[pathParts[2]]; !ok {
        return makeErrorResponse(http.StatusNotFound, "cluster not found", w, r)
    }
    if pathParts[2] == "" {
        // Allow a trailing / on requests
        return handleClusterList(app, w, r)
    }
    if (len(pathParts) == 3) || (pathParts[3] == "") {
        return handleClusterDetail(app, w, r, pathParts[2])
    }

    switch pathParts[3] {
    case "consumer":
        switch {
        case r.Method == "DELETE":
            switch {
            case (len(pathParts) == 5) || (pathParts[5] == ""):
                return handleConsumerDrop(app, w, r, pathParts[2], pathParts[4])
            default:
                return makeErrorResponse(http.StatusMethodNotAllowed, "request method not supported", w, r)
            }
        case r.Method == "GET":
            switch {
            case (len(pathParts) == 4) || (pathParts[4] == ""):
                return handleConsumerList(app, w, r, pathParts[2])
            case (len(pathParts) == 5) || (pathParts[5] == ""):
                // Consumer detail - list of consumer streams/hosts? Can be config info later
                return makeErrorResponse(http.StatusNotFound, "unknown API call", w, r)
            case pathParts[5] == "topic":
                switch {
                case (len(pathParts) == 6) || (pathParts[6] == ""):
                    return handleConsumerTopicList(app, w, r, pathParts[2], pathParts[4])
                case (len(pathParts) == 7) || (pathParts[7] == ""):
                    return handleConsumerTopicDetail(app, w, r, pathParts[2], pathParts[4], pathParts[6])
                }
            case pathParts[5] == "status":
                return handleConsumerStatus(app, w, r, pathParts[2], pathParts[4], false)
            case pathParts[5] == "lag":
                return handleConsumerStatus(app, w, r, pathParts[2], pathParts[4], true)
            }
        default:
            return makeErrorResponse(http.StatusMethodNotAllowed, "request method not supported", w, r)
        }
    case "topic":
        switch {
        case r.Method != "GET":
            return makeErrorResponse(http.StatusMethodNotAllowed, "request method not supported", w, r)
        case (len(pathParts) == 4) || (pathParts[4] == ""):
            return handleBrokerTopicList(app, w, r, pathParts[2])
        case (len(pathParts) == 5) || (pathParts[5] == ""):
            return handleBrokerTopicDetail(app, w, r, pathParts[2], pathParts[4])
        }
    case "offsets":
        // Reserving this endpoint to implement later
        return makeErrorResponse(http.StatusNotFound, "unknown API call", w, r)
    }

    // If we fell through, return a 404
    return makeErrorResponse(http.StatusNotFound, "unknown API call", w, r)
}
```

Bởi vì mặc định gói thư viện `net/http` hỗ trợ `mux` routing nhưng không hỗ trợ `arguments`, do đó mã nguồn trên dùng những kĩ thuật rất nhảm nhí như `Split` và bừa bộn như `switch case` để đạt được mục tiêu, điều đó thực sự làm chúng ta tập trung nhiều thời gian vào việc xử lý logic routing hơn là logic business. Nhìn qua hệ thống, thật khó để bảo trì và quản lý. Nếu bạn đọc mã nguồn cẩn thận, bạn sẽ thấy hàm phức tạp nhất là `handleKafka()`. Nhưng thực tế, hệ thống của chúng ta luôn luôn tập hợp nhiều những hàm gây phức tạp như vậy, và cuối cùng sẽ rất khó để làm sạch chúng.

Về kinh nghiệm của tôi rất đơn giản, những route có một số `parameters` và số lượng APIs cho dự án đó đạt đến tối đa là 10, thì đừng dùng `net/http` là thư viện route mặc định. Thư viện route được dùng rộng rãi nhất trong cộng đồng opensource Go là `httpRouter`, và nhiều opensource framework router khác cũng dựa trên httpRouter. Nguyên tắc của httpRouter được giải thích chi tiết ở phần router của chương này.

Nhìn lại phần đầu bài viết, có một vài frameworks trong thế giới opensource. Đầu tiên là một `wrap httpRouter` đơn giản có hỗ trợ để custom middleware và tích hợp một số tiện ích đơn giản như `gin`, nó nhẹ, dễ học, và hiệu nâng cao. Thứ hai là học mô hình MVC của các framework từ các ngôn ngữ lập trình khác.
Thêm vào đó là những nguyên tắc về router và middleware, nội dung của chương này phần lớn là những ví dụ cụ thể bằng mã nguồn Go.
