# 5.5 Database và giao tiếp với Database

Phần này sẽ thực hiện một số phân tích cơ bản về thư viện `db/sql` tiêu chuẩn và giới thiệu một số ORM và SQL Builder opensource được sử dụng rộng rãi. Đứng ở góc độ phát triển ứng dụng doanh nghiệp, sẽ phù hợp hơn để phân tích kiến trúc công nghệ nào phù hợp cho các ứng dụng doanh nghiệp hiện đại.

## 5.5.1 Bắt đầu từ database/sql

Go cung cấp một package `database/sql` để làm việc với cơ sở dữ liệu cho người dùng. Trên thực tế, thư viện `database/sql` chỉ cung cấp một bộ interface và thông số kỹ thuật để vận hành cơ sở dữ liệu, như SQL trừu tượng `prep`  (chuẩn bị), quản lý nhóm kết nối, liên kết dữ liệu (data binding), giao dịch, xử lý lỗi, và nhiều hơn nữa. Golang chính thức không cung cấp hỗ trợ giao thức cụ thể cho việc hiện thực cơ sở dữ liệu nhất định.

Để giao tiếp với một cơ sở dữ liệu nhất định như MySQL, bạn phải cung cấp driver MySQL như sau:

```go
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

db, err := sql.Open("mysql", "user:password@/dbname")
```

Ở dòng thứ hai thực sự là một hàm `init` gọi tới package `mysql`:

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

Đó cũng là một interface. Trong thực tế, nếu nhìn vào code của `database/sql/driver/driver.go` sẽ thấy rằng tất cả pprof các thành phần trong file đều là interface cả. Để thực thi được các kiểu này, bạn sẽ phải gọi tới những phương thức `driver` phù hợp.

Ở phía người dùng, trong process sử dụng package `databse/sql`, ta  có thể sử dụng các hàm được cung cấp trong những interface này,  hãy nhìn vào một ví dụ hoàn chỉnh sử dụng `database/sql` và `go-sql-driver/mysql`:

```go
package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    // db là một đối tượng của kiểu sql.DB
    // đối tượng là một thread-safe chứa kết nối
    // Tùy chọn kết nối có thể được đặt trong phương thức sql.DB, ở đây bỏ qua
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

    // Đọc nội dung các rows rồi gọi Close()
    // kết nối sẽ không được giải phóng cho đến khi defer rows.Close() thực thi
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

Bạn có thể cảm thấy một vài chỗ tiêu cực từ đoạn code thủ tục ngắn trên. Hàm cung cấp `db` của thư viện chuẩn quá đơn giản. Chúng ta có cần phải viết cùng một đoạn code mỗi lần ta truy cập database để đọc dữ liệu? Hoặc nếu đối tượng là struct, việc binding `sql.Rows` với đối tượng trở nên trùng lặp và nhàm chán.

Câu trả lời là Có, cho nên cộng đồng sẽ có rất nhiều các SQL Builder và ORM khác nhau.

## 5.5.2 ORM và SQL Builder để cải thiện hiệu suất

Hãy xem định nghĩa của ORM trên wikipedia:

> Object-relational mapping (ORM, O/RM, and O/R mapping tool) trong khoa học máy tính là một kĩ thuật lập trình cho phép chuyển đổi dữ liệu giữa các hệ thống kiểu không tương thích bằng ngôn ngữ hướng đối tượng. Điều này tạo ra một "cơ sở dữ liệu đối tượng ảo" có thể được sử dụng từ trong ngôn ngữ lập trình.

Thông thường ORM thực hiện việc mapping từ database tới các class hoặc struct của chương trình. Ví dụ như chương có thể mapping các class của mình từ các bảng trong MySQL. Đầu tiên hãy xem cách ORM trong các ngôn ngữ lập trình khác được viết như thế nào:

```sql
>>> from blog.models import Blog
>>> b = Blog(name='Beatles Blog', tagline='All the latest Beatles news.')
>>> b.save()
```

Không còn dấu vết nào của cơ sở dữ liệu ở đây nữa. Mục đích của ORM là che chắn lớp DB khỏi người sử dụng. Trên thực tế, ORM của nhiều ngôn ngữ chỉ định nghĩa class hoặc struct, sau đó sử dụng một cú pháp cụ thể để tạo ra struct tương ứng 1-1. Sau đó, ta có thể thực hiện các thao tác khác nhau trên các đối tượng đã map từ các bảng trong cơ sở dữ liệu như lưu, tạo, truy xuất và xóa. Đối với những gì ORM đã thực hiện ẩn bên dưới, ta không cần phải rõ ràng. Khi sử dụng ORM, chúng ta có xu hướng quên đi cơ sở dữ liệu. Ví dụ: ta có nhu cầu hiển thị cho người dùng danh sách sản phẩm mới nhất, giả định rằng hàng hóa và doanh nghiệp có mối quan hệ 1:1, có thể thể hiện bằng đoạn code sau:

```go
# mã giả
shopList := []
for product in productList {
    shopList = append(shopList, product.GetShop)
}
```

Công cụ như ORM là để bảo vệ cơ sở dữ liệu ngay từ điểm bắt đầu, cho phép chúng ta vận hành cơ sở dữ liệu gần hơn với cách suy nghĩ của con người. Vì vậy, nhiều lập trình viên mới tiếp xúc với ORM cũng có thể viết được code như trên.

Đoạn code trên sẽ phóng to yêu cầu đọc cơ sở dữ liệu theo hệ số của N. Nói cách khác, nếu danh sách sản phẩm có 15 SKU (Stock-Keeping Unit), mỗi lần người dùng mở trang, ít nhất 1 (danh sách mục truy vấn) + 15 (yêu cầu thông tin cửa hàng liên quan đến truy vấn) là bắt buộc. Ở đây N là 16. Nếu trang danh sách khá lớn, giả sử 600 mục, thì ta phải thực hiện ít nhất 1 + 600 truy vấn. Nếu số lượng truy vấn đơn giản lớn nhất mà cơ sở dữ liệu có thể chịu được là 120 000 QPS và truy vấn trên chỉ là truy vấn được sử dụng phổ biến nhất, thì khả năng service có thể cung cấp là bao nhiêu? 200 QPS! Một trong những nguyên tắc cấm kỵ của hệ thống Internet là sự khuếch đại số lượng thao tác đọc không cần thiết này.

Tất nhiên bạn có thể nói rằng đó không phải là vấn đề của ORM. Nếu viết bằng sql ta vẫn có thể viết được một chương trình giống vậy, hãy nhìn vào demo sau:

```go
o := orm.NewOrm()
num, err := o.QueryTable("cardgroup").Filter("Cards__Card__Name", cardName).All(&cardgroups)
```

Nhiều ORM cung cấp kiểu truy vấn `Filter` này, nhưng trên thực tế, đằng sau ORM còn ẩn nhiều thao tác chi tiết khác, chẳng hạn như tạo ra câu lệnh SQL tự động `limit 1000`.

Có lẽ một số bạn đọc sẽ thấy ngạc nhiên với thao tác đó. Thực ra trong tài liệu chính thức của ORM đã nói qua rằng tất cả các truy vấn sẽ tự động `limit 1000` mà không cần chỉ định rõ, chính vì vậy mà điều này trở nên khó khăn đối với nhiều người chưa đọc tài liệu hoặc đọc mã nguồn của ORM. Những lập trình viên thích ngôn ngữ ràng buộc kiểu mạnh thường không thích những gì ngôn ngữ tự thực hiện ngầm định, chẳng hạn như chuyển đổi kiểu ngầm của các ngôn ngữ khác nhau trong thao tác gán để rồi mất đi độ chính xác trong chuyển đổi, điều này chắc chắn sẽ khiến họ đau đầu. Vì vậy, càng có ít thứ mà thư viện làm ẩn bên dưới thì càng tốt. Nếu ta cần thực hiện điều gì hãy thực hiện nó ở một nơi dễ thấy. Trong ví dụ trên, tốt hơn hết là loại bỏ hành vi tự hành động ngầm định này hoặc là bắt buộc người dùng phải truyền vào tham số `limit`.

Ngoài vấn đề giới hạn, chúng ta hãy xem truy vấn này dưới đây:

```go
num, err := o.QueryTable("cardgroup").Filter("Cards__Card__Name", cardName).All(&cardgroups)
```

Bạn có thấy rằng `Filter` này là một thao tác JOIN không? Rất khó để nhận ra vì ORM đã che giấu quá nhiều chi tiết khỏi thiết kế. Cái giá của sự tiện lợi là những hoạt động ẩn đằng sau nó hoàn toàn nằm ngoài kiểm soát. Một dự án như vậy sẽ trở nên ngày càng khó theo dõi và bảo trì chỉ sau một vài lần nâng cấp.

Tất nhiên chúng ta không thể phủ nhận được tầm quan trọng của ORM. Mục đích ban đầu của nó là loại bỏ việc triển khai cụ thể các hoạt động với database và lưu trữ dữ liệu. Nhưng một số công ty đã dần xem ORM có thể là một thiết kế thất bại vì các chi tiết quan trọng bị ẩn giấu. Các chi tiết quan trọng ẩn rất quan trọng đối với sự phát triển của các hệ thống cần mở rộng quy mô.

So sánh với ORM, SQL Builder đạt được sự cân bằng tốt hơn giữa SQL và khả năng bảo trì của dự án. Đầu tiên, sql builder không ẩn quá nhiều chi tiết như ORM. Thứ hai, từ góc độ phát triển, SQL Builder cũng có thể hoàn thiện rất hiệu quả chỉ sau vài thao tác đóng gói đơn giản:

```go
where := map[string]interface{} {
    "order_id > ?" : 0,
    "customer_id != ?" : 0,
}
limit := []int{0,100}
orderBy := []string{"id asc", "create_time desc"}

orders := orderModel.GetList(where, limit, orderBy)
```

Việc code SQL Builder và đọc nó đều không gặp khó khăn gì. Chuyển đổi những dòng code này thành sql cũng không cần quá nhiều nỗ lực. Thông qua code ta có thể lấy được database index trên truy vấn này, xem kết quả thông qua index và nó liệu có thể phân tích với index đã JOIN hay không.

Nói một cách dễ hiểu, SQL Builder là một cách biểu diễn ngôn ngữ đặc biệt của sql trong mã. Nếu bạn không có DBA, nhưng R & D có khả năng phân tích và tối ưu hóa sql hoặc DBA của công ty bạn không phản đối  các kiểu ngôn ngữ sql như thế này thì bạn sử dụng SQL Builder là một lựa chọn rất tốt.

Ngoài ra, trong một số trường hợp không yêu cầu can thiệp DBA, ta cũng có thể sử dụng SQL Builder. Ví dụ: nếu muốn tạo một bộ các thao tác với database và bảo trì hệ thống, coi MySQL là một thành phần trong hệ thống, thì QPS của hệ thống sẽ không cao đồng thời truy vấn cũng không quá phức tạp.

Khi bạn đang thực hiện một hệ thống trực tuyến OLTP (On-line transactional processing) high-concurrency và bạn muốn tối đa hóa rủi ro của hệ thống với sự phân công thực thi rõ ràng thì việc sử dụng SQL Builder là không phù hợp.

## 5.5.3 Database mỏng manh

Cả ORM và SQL Builder đều có một lỗ hổng nghiêm trọng là không cách nào thực hiện việc kiểm toán trước-sql trên hệ thống. Mặc dù nhiều ORM và SQL Builder cũng cung cấp chức năng in sql khi chạy, nhưng nó chỉ có thể là đầu ra sau khi truy vấn. Chức năng được cung cấp bởi SQL Builder và ORM quá linh hoạt, nó khiến ta không thể liệt kê tất cả các sql có thể được thực thi trực tuyến bằng cách test. Ví dụ: ta có thể sử dụng SQL Builder để viết đoạn mã sau:

```go
where := map[string]interface{} {
    "product_id = ?" : 10,
    "user_id = ?" : 1232 ,
}

if order_id != 0 {
    where["order_id = ?"] = order_id
}

res, err := historyModel.GetList(where, limit, orderBy)
```

Nếu bạn có nhiều ví dụ tương tự trong hệ thống của mình, rất khó để bao quát tất cả các kết hợp SQL có thể có thông qua các test case.

Một hệ thống như vậy khi release sẽ sinh ra rủi ro rất lớn.

Đối với các công ty Internet có dịch vụ 24/7, việc đảm bảo dịch vụ luôn vận hành là vấn đề rất quan trọng. Mặc dù kiến trúc công nghệ của tầng lưu trữ đã trải qua nhiều năm phát triển, nó vẫn là phần dễ bị hư hại nhất trong toàn hệ thống. Thời gian chết của hệ thống cũng chính là tổn thất kinh tế trực tiếp cho các công ty cung cấp kiểu này.

Từ góc độ phân ngành, các công ty Internet ngày nay có chức vụ DBA toàn thời gian, hầu hết họ không nhất thiết phải có khả năng code. Từ quan điểm của DBA, chúng tôi vẫn hy vọng sẽ có một cơ chế kiểm toán SQL đặc biệt để có được tất cả nội dung SQL của hệ thống với chi phí thấp, thay vì đọc code SQL Builder do bộ phận bussiness development viết.

Ngày nay, core online business của các công ty lớn sẽ cung cấp SQL ở vị trí nổi bật trong code để DBA review, ví dụ:

```go
const (
    getAllByProductIDAndCustomerID = `select * from p_orders where product_id in (:product_id) and customer_id=:customer_id`
)

// GetAllByProductIDAndCustomerID
// @param driver_id
// @param rate_date
// @return []Order, error
func GetAllByProductIDAndCustomerID(ctx context.Context, productIDs []uint64, customerID uint64) ([]Order, error) {
    var orderList []Order

    params := map[string]interface{}{
        "product_id" : productIDs,
        "customer_id": customerID,
    }

    // getAllByProductIDAndCustomerID là một chuỗi thuộc kiểu sql
    sql, args, err := sqlutil.Named(getAllByProductIDAndCustomerID, params)
    if err != nil {
        return nil, err
    }

    err = dao.QueryList(ctx, sqldbInstance, sql, args, &orderList)
    if err != nil {
        return nil, err
    }

    return orderList, err
}
```

Khá thuận tiện khi lấy phần const của lớp DAO (Data Access Object) đặt trực tiếp vào DBA để review trước khi online. `Sqlutil.Named` ở đây tương tự như hàm được đặt tên trong sqlx và hỗ trợ các toán tử so sánh trong biểu thức `where`.
