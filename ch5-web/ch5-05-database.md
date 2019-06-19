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

Interface `Driver` trong package `sql`:

```go
type Driver interface {
    Open(name string) (Conn, error)
}
```

`sql.Open()` trả về đối tượng `db` từ lời gọi hàm `Conn`

```go
type Conn interface {
    Prepare(query string) (Stmt, error)
    Close() error
    Begin() (Tx, error)
}
```

Đó cũng là một interface. Trong thực thế, nếu nhìn vào code của `database/sql/driver/driver.go` sẽ thấy rằng tất cả  các thành phần trong file đều là interface cả. Để thực thi được các kiểu này, bạn sẽ phải gọi tới những phương thức `driver` phù hợp.

Ở phía người dùng, trong process sử dụng package `databse/sql`, ta  có thể sử dụng các hàm được cung cấp trong những interface này,  hãy nhìn vào một ví dụ hoàn chỉnh sử dụng `database/sql` và `go-sql-driver/mysql`:

```go
package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // db 是一个 sql.DB 类型的对象
    // 该对象线程安全，且内部已包含了一个连接池
    // 连接池的选项可以在 sql.DB 的方法中设置，这里为了简单省略了
    db, err := sql.Open("mysql",
        "user:password@tcp(127.0.0.1:3306)/hello")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    var (
        id int
        name string
    )
    rows, err := db.Query("select id, name from users where id = ?", 1)
    if err != nil {
        log.Fatal(err)
    }

    defer rows.Close()

    // 必须要把 rows 里的内容读完，或者显式调用 Close() 方法，
    // 否则在 defer 的 rows.Close() 执行之前，连接永远不会释放
    for rows.Next() {
        err := rows.Scan(&id, &name)
        if err != nil {
            log.Fatal(err)
        }
        log.Println(id, name)
    }

    err = rows.Err()
    if err != nil {
        log.Fatal(err)
    }
}
```

Nếu bạn đọc muốn biết `database/sql` chi tiết hơn, có thể xem tại <http://go-database-sql.org/>.

Một vài hiện thực bao gồm các hàm, giới thiệu, cách sử dụng, các cảnh báo và các phản trực quan (counter-intuition) về thư viện (ví dụ như `sql.DB`, các truy vấn trong cùng goroutine có thể ở trên nhiều connections) đều được đề cập, và chúng sẽ không được nhắc tới nữa trong chương này.

Bạn có thể đã cảm thấy một vài chỗ tiêu cực từ đoạn code thủ tục ngắn trên. Hàm cung cấp `db` của thư viện chuẩn quá đơn giản. Chúng ta có cần phải viết cùng một đoạn code mỗi lần ta truy cập database để đọc dữ liệu? Hoặc nếu đối tượng là struct, việc binding `sql.Rows` với đối tượng trở nên trùng lặp và nhàm chán.

Câu trả lời là Có, cho nên cộng dộng sẽ có rất nhiều  các SQL Builder và ORM khác nhau.

## 5.5.2 ORM và SQL Builder để cải thiện hiệu suất
