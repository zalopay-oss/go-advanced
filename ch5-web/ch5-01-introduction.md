# 5.1 Giới thiệu về Web Development

Bởi vì gói thư viện `net/http` của Golang chỉ hỗ trợ cơ bản những hàm routing và hàm chức năng. Cho nên trong cộng đồng Golang có một ý tưởng phổ biến là viết thêm các APIs. Theo ý kiến của tôi, nếu các project routing của bạn có những đặc điểm sau: URI cố định, và tham số không truyền thông qua URI, thì nên dùng thư viện chuẩn. Nhưng với những ngữ cảnh phức tạp hơn, thư viện chuẩn `http` vẫn còn một vài điểm yếu. Ví dụ, xét route sau:

```
GET   /card/:id
POST  /card/:id
DELTE /card/:id
GET   /card/:id/name
...
GET   /card/:id/relations
```

Có thể thấy rằng, khi nào framework được dùng hoặc khi nào những vấn đề cụ thể được phân tích

Framework web của Go có thể được chia thành hai phần như sau

1. Router framework
2. MVC class framework

Khi chọn một framework, trong nhiều trường hợp sẽ theo tham khảo của cá nhân về những công nghệ mà công ty đang sử dụng. Ví dụ, nếu công ty có nhiều người làm về `PHP`, thì chúng ta nên chọn framework `beego`, nhưng nếu công ty có nhiều lập trình viên `C`, thì những ý tưởng của họ sẽ đơn giản hết mức có thể. Ví dụ, nhiều lập trình viên C trong những công ty lớn sẽ dùng ngôn ngữ C để viết một chương trình `CGI` nhỏ. Họ không thể sẵn sàng để học `MVC` hoặc nhiều framework Web phức tạp khác. Tất cả họ cần những route đơn giản. Mặc dù route thì không cần, trong khi chỉ cần một thư viện xử lý giao thức HTTP cơ bản để giúp anh ta tiết kiệm công sức làm việc thủ công.

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

Hãy đào sâu mã nguồn trên, dự án Kafka monitoring của một công ty nổi tiếng linkedin. Nếu chúng không dùng bất cứ route framework nào và chỉ dùng `net/http`. Nhìn lại mã nguồn trên dường như chúng rất đẹp, chỉ có nằm URIs đơn giản trong project của chúng ta, do đó service chúng ta hỗ trợ như sau

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

Bởi vì gói thư viện `net/http`mặc định `mux` routing không hỗ trợ `arguments`, do đó mã nguồn trên dùng những kĩ thuật rất nhảm nhí như `Split` và bừa bộn như `switch case` để đạt được mục tiêu, nhưng điều đó thực sự làm chúng ta tập trung vào việc quản lý logic routing trở nên phức tạp. Nhìn qua hệ thống thật khó để bảo trì và quản lý. Nếu bạn đọc mã nguồn cẩn thận, bạn sẽ thấy hàm phức tạp nhất là `handleKafka()`. Nhưng thực tế, hệ thống của chúng ta luôn luôn tập hợp nhiều những hàm gây phức tạp như vậy, và cuối cùng sẽ rất khó để làm sạch chúng.

Về kinh nghiệm của tôi rất đơn giản, những route có một số `parameters` và số lượng APIs cho dự án đó đạt đến tối đa là 10, thì đừng `net/http` là thư viện route mặc định. Thư viện route được dùng rộng rãi nhất trong cộng đồng opensource Go là httpRouter, và nhiều opensource router framework dựa trên httpRouter để chắc chắn đạt được một sự biến đổi về độ. Nguyên tắc của httpRouter được giải thích chi tiết ở phần router của chương này.

Nhìn lại khi bắt đầu bài viết, có một vài framework trong thế giới opensource. Đầu tiên là một `wrap httpRouter` đơn giản và sau đó hỗ trợ để custom middleware và tích hợp một số tiện ích đơn giản như `gin`, nó nhẹ, dễ học, và hiệu năng cao. Thứ hai là học từ một số phong cách của các ngôn ngữ lập trình khác. Chỉ có một số framework mạnh mẽ, ngoại trừ thiết kế  `database schema`, trong khi hầu hết mã nguồn là sinh ra một cách trực tiếp, như là `goa`. Theo những framework đó, sẽ phụ thuộc vào nền tảng của người lập trình.

Thêm vào đó, những nguyên tắc về router và middleware, nội dung của chương này sẽ kết hợp với Go để giải thích một số ví dụ. Tôi hy vọng có thể giúp ích được người đọc chưa biết về những nội dung đó.

