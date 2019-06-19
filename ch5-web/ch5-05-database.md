# Database và giao tiếp với Database

Phần này sẽ thực hiện một số phân tích cơ bản về thư viện `db/sql` tiêu chuẩn và giới thiệu một số ORM và SQL Builder mã nguồn mở được sử dụng rộng rãi. Đứng ở góc độ phát triển ứng dụng doanh nghiệp, sẽ phù hợp hơn để phân tích kiến trúc công nghệ nào phù hợp cho các ứng dụng doanh nghiệp hiện đại.

## 5.5.1 Bắt đầu từ database/sql

Go cung cấp một package `database/sql` để làm việc với cơ sở dữ liệu cho người dùng. Trên thực tế, thư viện `database/sql` chỉ cung cấp một bộ interface và thông số kỹ thuật để vận hành cơ sở dữ liệu, như SQL trừu tượng `prep`  (chuẩn bị), quản lý nhóm kết nối, liên kết dữ liệu (data binding), giao dịch, xử lý lỗi, và nhiều hơn nữa. Golang chính thức không cung cấp hỗ trợ giao thức cụ thể cho việc hiện thực cơ sở dữ liệu nhất định.

Để giao tiếp với một cơ sở dữ liệu nhất định, như MySQL, bạn phải cung cấp driver MySQL như sau:

```go
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

db, err := sql.Open("mysql", "user:password@/dbname")
```

Ở dòng thứ hai thực sư là một hàm `init` gọi tới package `mysql`:

```go
func init() {
    sql.Register("mysql", &MySQLDriver{})
}
```

Trong global package `sql` 
