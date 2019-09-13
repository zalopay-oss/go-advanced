# 4.2. Routing trong Web

Trong phần trước, chúng ta đã tìm hiểu cách dùng thư viện chuẩn [http/net](https://golang.org/pkg/net/http/) để hiện thực hàm routing đơn giản. Tuy nhiên một framework web sẽ có nhiều thành phần khác ngoài việc định tuyến như xử lý tham số URI, phương thức, mã lỗi.

## 4.2.1 RESTful API

[RESTful](https://restfulapi.net/) là một tiêu chuẩn thiết kế API trong ngành công nghiệp web hiện đại. Ngoài những phương thức GET, POST thì RESTful cũng định nghĩa vài phương thức khác trong giao thức HTTP bao gồm:

***Phương thức HTTP:***

```go
const (
    MethodGet     = "GET"
    MethodHead    = "HEAD"
    MethodPost    = "POST"
    MethodPut     = "PUT"
    MethodPatch   = "PATCH" // RFC 5789
    MethodDelete  = "DELETE"
    MethodConnect = "CONNECT"
    MethodOptions = "OPTIONS"
    MethodTrace   = "TRACE"
)
```

Nhìn vào những đường dẫn RESTful API sau:

```sh
// mỗi API sẽ có một phương thức tương ứng
// tham số được truyền vào thông qua URI

GET /repos/:owner/:repo/comments/:id/reactions
POST /projects/:project_id/columns
PUT /user/starred/:owner/:repo
DELETE /user/starred/:owner/:repo

```

Nếu hệ thống web của chúng ta cần có những API tương tự trên, việc sử dụng thư viện chuẩn net/http hiển nhiên là không đủ. Những API chứa parameters như trên của Github có thể được hỗ trợ hiện thực bởi thư viện [HttpRouter](https://github.com/julienschmidt/httprouter).

## 4.2.1 Tìm hiểu thư viện HttpRouter

Nhiều Open source web framework phổ biến của Go thường được xây dựng dựa trên [HttpRouter](https://github.com/julienschmidt/httprouter) như là [Gin](https://github.com/gin-gonic) framework, hoặc hỗ trợ cho routing dựa trên những biến thể của HttpRouter. Khi sử dụng các framework đó, chúng ta cần phải tránh một số trường hợp mà nó dẫn đến xung đột routing khi thiết kế.

***Ví dụ:***

```sh
// xung đột trong trường hợp đặc biệt id là 'info'
// vì cùng phương thức nên cùng nằm trên một 'cây định tuyến'
// 'cây định tuyến' được nói ở phần sau
GET /user/info/:name
GET /user/:id

// không xung đột vì khác phương thức
// nên sẽ tạo ra hai 'cây định tuyến' cho hai phương thức khác nhau
GET /user/info/:name
POST /user/:id

// các lỗi trên sẽ bị bắt lỗi panic trong HttpRouter
```

HttpRouter hỗ trợ kí tự đặc biệt `*` trong đường dẫn.

***Ví dụ:***

```sh
Pattern: /src/*filepath

/src/                     filepath = ""
/src/somefile.go          filepath = "somefile.go"
/src/subdir/somefile.go   filepath = "subdir/somefile.go"

// thiết kế này thường dùng để xây dựng một static file server
```

HttpRouter cũng hiện thực tùy chỉnh hàm callback trong một vài trường hợp đặc biệt như là lỗi 404:

***Ví dụ:***

```go
r := httprouter.New()
r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("oh no, not found"))
})
```

Hoặc tùy chỉnh hàm callback khi panic bên trong:

***Ví dụ:***

```go
r.PanicHandler = func(w http.ResponseWriter, r *http.Request, c interface{}) {
    log.Printf("Recovering from panic, Reason: %#v", c.(error))
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(c.(error).Error()))
}
```

## 4.2.2 Cấu trúc dữ liệu trong HttpRouter

Cấu trúc dữ liệu được dùng bởi HttpRouter và nhiều framework routing dẫn xuất khác là [Radix Tree](https://en.wikipedia.org/wiki/Radix_tree). Cây Radix thường được dùng để truy xuất chuỗi, để xem chúng có nằm trong cây hay không và lấy thông tin gắn với chuỗi đó, phương pháp tìm kiếm theo chiều sâu sẽ bắt đầu từ node gốc, và thời gian xấp xỉ là `O(n)`, và n là chiều sâu của cây.


<div align="center">
	<img src="../images/ch4-2-Patricia_trie.svg.png" width="500">
	<br/>
	<span align="center">
		<i>Cây Radix tree</i>
	</span>
	<br/>
</div>

Kiểu chuỗi không phải là một kiểu số học nên không thể so sánh trực tiếp như kiểu số, và thời gian xấp xỉ của việc so sánh hai chuỗi là phụ thuộc vào độ dài của chuỗi, và sau đó dùng giải thuật như là [binary search](https://vi.wikipedia.org/wiki/Tìm_kiếm_nhị_phân) để tìm kiếm, độ phức tạp về thời gian có thể cao. Dùng [Radix tree](https://en.wikipedia.org/wiki/Radix_tree) để lưu trữ và truy xuất chuỗi là một cách đảm bảo tối ưu về thời gian, mỗi phần trong đường dẫn được xem là một chuỗi và được lưu trữ trong cây Radix như ví dụ sau:

<div align="center">
	<img src="../images/ch5-02-radix.png" width="500">
	<br/>
	<span align="center"><i></i></span>
	<br/>
</div>

## 4.2.3 Xây dựng Radix tree

Hãy xét quy trình của một Radix tree trong HttpRouter. Phần thiết lập routing có thể như sau:

```sh
PUT /user/installations/:installation_id/repositories/:repository_id

GET /marketplace_listing/plans/
GET /marketplace_listing/plans/:id/accounts
GET /search
GET /status
GET /support

GET /marketplace_listing/plans/ohyes
```

### 4.2.3.1 Khởi tạo

Radix tree có thể được lưu trữ trong cấu trúc của Router trong HttpRouter sử dụng một số cấu trúc dữ liệu sau:

```go
// Router struct
type Router struct {
  // ...
  trees map[string]*node
  // Trong đó,
  // key: GET, HEAD, OPTIONS, POST, PUT, PATCH hoặc DELETE
  // value: node cha của cây Radix
  // ...
}
```

Mỗi phương thức sẽ tương ứng với một Radix tree độc lập và không chia sẻ dữ liệu với các cây khác. Đặc biệt đối với route chúng ta dùng ở trên, `PUT` và `GET` là hai Radix tree thay vì một. Đầu tiên, chèn route `PUT` vào Radix tree:

```go
r := httprouter.New()
r.PUT("/user/installations/:installation_id/repositories/:reposit", Hello)
```

<div align="center">
	<img src="../images/ch5-02-radix-put.png" width="900">
	<br/>
	<span align="center">
		<i>Một cây từ điển nén được insert vào route</i>
	</span>
	<br/>
</div>

Kiểu của mỗi node trong Radix tree là `*httprouter.node`, trong đó, một số trường mang ý nghĩa sau:

```sh
path: // đường dẫn ứng với node hiện tại
wildChild: // cho dù là nút con tham số, nghĩa là nút có ký tự đại diện hoặc :id
nType:    // loại nút có bốn giá trị liệt kê static/root/param/catchAll
  static  // chuỗi bình thường cho các node không gốc
  root    // nút gốc
  param   // nút tham số ví dụ :id
  catch   // các nút ký tự đại diện, chẳng hạn như * anyway
indices:
```

Tiếp theo, chúng ta chèn các route GET còn lại trong ví dụ để giải thích về quy trình chèn vào một node con.

### 4.2.3.2 Chèn các route khác

Khi chúng ta chèn `GET /marketplace_listing/plans`, quá trình này sẽ tương tự như trước nhưng ở một Radix tree khác:

<div align="center">
	<img src="../images/ch5-05-radix-get-1.png" width="300">
	<br/>
	<span align="center">
		<i>Chèn node đầu tiên vào Radix tree</i>
	</span>
	<br/>
</div>


Sau đó chèn đường dẫn `GET /marketplace_listing/plans/:id/accounts` cấu trúc Radix tree được hoàn thành sẽ như sau:

<div align="center">
	<img src="../images/ch5-02-radix-get-2.png" width="400">
	<br/>
	<span align="center">
		<i>Chèn node thứ hai vào Radix tree</i>
	</span>
</div>
<br/>


### 4.2.3.3 Phân nhánh

Tiếp theo chúng ta chèn `GET /search`, sau đó sẽ sinh ra split tree như hình 5.6:

<div align="center">
	<img src="../images/ch5-02-radix-get-3.png" width="800">
	<br/>
	<span align="center">
		<i>Chèn vào node thứ ba sẽ gây ra việc phân nhánh</i>
	</span>
	<br/>
</div>

Node gốc bây giờ sẽ bắt đầu từ ký tự `/`, chuỗi truy vấn phải bắt đầu từ node gốc chính, sau đó một route là `search` được phân nhánh từ gốc. Tiếp theo chèn  `GET /status` và `GET /support` vào Redix tree. Lúc này, sẽ dẫn đến node `search` bị tách một lần nữa, và kết quả cuối cùng được nhìn thấy ở hình dưới:

<div align="center">
	<img src="../images/ch5-02-radix-get-4.png" width="800">
	<br/>
	<span align="center">
		<i>Sau khi chèn tất cả các node</i>
	</span>
	<br/>
	<br/>
</div>

<div style="display: flex; justify-content: space-around;">
<span> <a href="ch4-01-intro.md">&lt Phần 4.1</a>
</span>
<span><a href="../SUMMARY.md"> Mục lục</a>  </span>
<span> <a href="ch4-03-middleware.md">Phần 4.3 &gt</a> </span>
</div>
