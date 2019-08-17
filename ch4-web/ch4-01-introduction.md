# 4.1. Giới thiệu về Web Development

Phần này sẽ đề cập về cách xây dựng một chương trình web đơn giản bằng thư viện chuẩn của Go, sau đó giới thiệu các framework web trong cộng đồng Open source.

## 4.1.1 Dùng thư viện chuẩn net/http

Gói thư viện [net/http](https://golang.org/pkg/net/http/) đã cung cấp những hàm cơ bản cho việc routing URL, chúng ta sẽ dùng nó để viết một chương trình `http echo server`:

***echo.go:***

```go
package main
// các gói thư viện cần import
import (
    "io/ioutil"
    "log"
    "net/http"
)
// hàm routing echo, gồm hai params
// r *http.Request : dùng để đọc yêu cầu từ client
// wr http.ResponseWriter : dùng để ghi phản hồi về client
func echo(wr http.ResponseWriter, r *http.Request) {
    // đọc thông điệp mà client gửi tới trong r.Body
    msg, err := ioutil.ReadAll(r.Body)
    // phản hồi về client lỗi nếu có
    if err != nil {
        wr.Write([]byte("echo error"))
        return
    }
    // phản hồi về client chính thông điệp mà client gửi
    writeLen, err := wr.Write(msg)
    // nếu lỗi xảy ra, hoặc kích thước thông điệp phản hồi khác
    // kích thước thông điệp nhận được
    if err != nil || writeLen != len(msg) {
        log.Println(err, "write len:", writeLen)
    }
}
// hàm main của chương trình
func main() {
    // mapping url ứng với hàm routing echo
    http.HandleFunc("/", echo)
    // địa chỉ http://127.0.0.1:8080/
    err := http.ListenAndServe(":8080", nil)
    // log ra lỗi nếu bị trùng port
    if err != nil {
        log.Fatal(err)
    }
}
```

Kết quả khi chạy chương trình:

```sh
$ go run echo.go &
$ curl http://127.0.0.1:8080/ -d '"Hello, World"'
"Hello, World"
```

## 4.1.2 Dùng thư viện bên ngoài

Bởi vì gói thư viện chuẩn [net/http](https://golang.org/pkg/net/http/) của Golang chỉ hỗ trợ những hàm routing và hàm chức năng cơ bản. Cho nên trong cộng đồng Golang có ý tưởng là viết thêm các thư viện hỗ trợ routing khác ngoài `net/http`.

Thông thường, nếu các dự án routing HTTP của bạn có những đặc điểm sau: [URI](https://vi.wikipedia.org/wiki/URI) cố định, và tham số không truyền thông qua URI, thì nên dùng thư viện chuẩn là đủ. Nhưng với những trường hợp phức tạp hơn, thư viện chuẩn `net/http` vẫn còn thiếu các chức năng hỗ trợ. Ví dụ, xét các route sau:

```sh
GET     /card/:id
POST    /card/:id
DELETE  /card/:id
GET     /card/:id/name
GET     /card/:id/relations
```

Có thể thấy rằng, cùng là đường dẫn có chứa `/card/:id`, nhưng có phương thức khác nhau hoặc nhánh con khác nhau sẽ dẫn đến logic xử lý khác nhau, cách xử lý những đường dẫn trùng tên như vậy thường sẽ phức tạp. Khi đó chúng ta có thể nghĩ đến việc sử dụng một số framework routing bên ngoài từ cộng đồng Open source.

Framework web của Go có thể được chia thành hai loại sau:

1. Router framework ([HttpRouter](https://github.com/julienschmidt/httprouter), [Gin](https://github.com/gin-gonic/gin), [Gorilla](https://github.com/gorilla/mux),...)

2. MVC class framework ([Revel](https://github.com/revel/revel), [Beego](https://github.com/astaxie/beego), [Iris](https://github.com/kataras/iris),...)

Chúng ta có thể xem thống kê các framwork web phổ biến được dùng trong cộng đồng Golang [ở đây](https://github.com/mingrammer/go-web-framework-stars/blob/master/README.md).