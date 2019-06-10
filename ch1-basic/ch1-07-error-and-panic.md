# 1.7 Error và Exceptions

Error handling là một chủ đề quan trọng để cân nhắc lựa chọn một ngôn ngữ lập trình. Trong Go error handling, errors là một phần quan trọng của gói API và giao diện người dùng.

Đó cũng là một số hàm trong chương trình luôn luôn yêu cầu phải chạy thành công. Ví dụ `strconv.Itoa` chuyển một số nguyên thành string, đọc và ghi phần tử từ array đến slice, đọc một phần tử tồn tại trong `map` và tương tự. Những tác vụ như thế sẽ khó để có lỗi trong thời gian chạy trừ phi có bug trong chương trình và những tình huống không đoán trước được, như là memory leak tại thời điểm chạy. Nếu bạn thực sự bắt gặp một tình huống khác thường như thế, chúng ta có thể đơn giản là ngừng thực thi chương trình.

Ngoại trừ trường hợp bất thường, nếu chương trình bị lỗi và dừng, có thể cần nhắc đến vài kết quả không mong đợi. Những hàm sẽ xem xét function như mong đợi, chúng ta sẽ trả về một kết quả phụ, thông thường kết quả cuối cùng có thể truyền vào thông điệp lỗi. Nếu chúng chỉ có một nguyên nhân gây ra lỗi, thông tin thêm đó có thể đơn giản là một giá trị Boolean, thông thường đặt tên là ok. cho một ví dụ, khi kết quả truy vấn `map`, bạn có thể sử dụng thêm một giá trị Boolean để xác định chúng có thành công hay không:

```go
if v, ok := m["key"]; ok {
    return v
}
```

Nhưng thông thường sẽ có nhiều hơn một nguyên nhân gây ra lỗi, và nhiều lần user muốn biến nhiều về lỗi đó. Nếu bạn chỉ sử dụng một biến boolean, thì bạn sẽ không giải quyết được yêu cầu trên. trong ngôn ngữ C, một số nguyên `errno` được sử dụng mặc định để truyền tải lỗi, do đó bạn có thể định nghĩa nhiều loại error theo nhu cầu. Trong ngôn ngữ Go, có thể gọi `syscall.Errno` là một `errno` ứng với mã lỗi trong ngôn ngữ C. `syscall` interface trong package, nếu có một error trả về, bên dưới cũng phải là `syscall.Errno` kiểu của error.

Ví dụ, khi chúng ta sửa đổi `syscall` để thay đổi chế độ của một file thông qua interface của package đó, nếu chúng ta bắt gặp một error, chúng ta có thể xử lý chúng bởi việc gây ra `err` trong phần `assertion` như là `syscall.Errno` là một kiểu error.

```go
err := syscall.Chmod(":invalid path:", 0666)
if err != nil {
    log.Fatal(err.(syscall.Errno))
}
```

Chúng ta có thể xa hơn chứa true error type thông qua một loại truy vấn kiểu hoặc assertions, do đó chúng ta có thể lấy nhiều thông tin về loại error. Tuy nhiên, tổng quát, chúng ta không quan tâm về cách mà error được thể hiện bên dưới. Chúng ta có thể chỉ cần biết rằng đó là một lỗi. Khi chúng ta trả về một giá trị error không phải `nil`, chúng ta có thể lấy một thông điệp error bởi việc gọi error interface type hoặc phương thức `Error`.

Trong ngôn ngữ Go, errors được xem xét như là một kết quả đã được đoán trước; ngoại lệ là một kết quả không thể đoán trước được, và một ngoại lệ có thể chỉ ra rằng một bug trong chương trình hoặc một vấn đề nào đó không được kiểm soát, nó sẽ cho phép user có thể quan tâm về những vấn đề về business liên quan đến việc xử lý lỗi.

Nếu một interface đơn giản ném tất cả những lỗi thông thường như là một ngoại lệ, chúng sẽ làm thông báo lỗi lộn xộn và không có giá trị. Chỉ như `main` bao gồm mọi thứ trực tiếp trong một hàm, nó không mang lại ý nghĩa gì.

```go
func main() {
    defer func() {
        if r := recover(); r != nil {
            log.Fatal(r)
        }
    }()

}
```

Bao bọc một mã lỗi không phải là một kết quả cuối cùng. Nếu một ngoại lệ không thể đoán trước được, trực tiếp gây ra một ngoại lệ là một cách tốt nhất để xử lý chúng.


### 1.7.1 Chiến lược xử lý lỗi

Hãy minh họa cho ví dụ về sao chép file: một hàm cần phải mở hai file và sau đó sao chép toàn bộ nội dung của một file nào đó về một file khác.

```go
func CopyFile(dstName, srcName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }

    dst, err := os.Create(dstName)
    if err != nil {
        return
    }

    written, err = io.Copy(dst, src)
    dst.Close()
    src.Close()
    return
}
```

Khi đoạn code trên chạy, nhưng bỏ qua một bug. Nếu đầu tiên `os.Open` gọi thành công, nhưng lệnh gọi thứ hai `os.Create` gọi bị failed, nó sẽ trả về  ngay lặp tức mà không giải phóng tài nguyên file đầu tiên. Mặc dù chúng ta có thể gọi `src.Close()` để fix bug bằng việc thêm vào lệnh gọi đó trước lệnh return về mệnh đề return thứ hai; nhưng khi code trở nên phức tạp hơn, những vấn đề tương tự sẽ khó để tìm thấy và giải quyết. Chúng ta có thể sử dụng mệnh đề `defer` để đảm bảo rằng một file bình thường được mở sẽ được đóng bình thường.


```go
func CopyFile(dstName, srcName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.Create(dstName)
    if err != nil {
        return
    }
    defer dst.Close()

    return io.Copy(dst, src)
}
```

Mệnh đề `defer` sẽ cho phép chúng ta nghĩ về làm cách nào để đóng một file ngay sau khi mở file đó. Bất kể làm thế nào hàm được trả về, bởi về mệnh đề close có thể luôn luôn được thực thi. Cùng một thời điểm, mệnh đề defer sẽ đảm bảo rằng `io.Copy` file có thể được đóng an toàn nếu một ngoại lệ xảy ra.

Như chúng ta đã đề cập trước đó, hàm exported trong ngôn ngữ Go sẽ thông thường ném ra một ngoại lệ, và một ngoại lệ không được kiểm soát có thể xem là một bug trong một chương trình.

Nhưng với những framework chúng cung cấp những web service tương tự, chúng thường cần sự truy cập từ bên thứ ba ở middleware. Bởi vì thư viện middleware thứ ba có bug, khi mà một ngoại lệ ném một exception, web framework bản thân nó không chắc chắn. Để cải thiện sự bền vững của hệ thống, web framework thường thu hồi chính xác nhất có thể những ngoại lệ trong luồng thực thi của chương trình và sau đó sẽ gây exception về bằng cách return error thông thường.

Chúng ta hãy xem JSON parse là một ví dụ minh họa cho việc dùng ngữ cảnh của việc phục hồi. Cho một hệ thống JSON parser phức tạp, mặc dù một ngôn ngữ parse có thể làm việc một cách phù hợp, có một điều không chắc chắn rằng nó không có lỗ hỏng. Do đó, khi một ngoại lệ xảy ra, chúng ta sẽ không chọn cách crash parser. Thay vì thế chúng ta sẽ làm việc với ngoại lệ panic nhưng là một lỗi parsing thông thường và đính kèm chúng với một thông tin thêm để thông báo cho user biết mà báo cáo lỗi.

```go
func ParseJSON(input string) (s *Syntax, err error) {
    defer func() {
        if p := recover(); p != nil {
            err = fmt.Errorf("JSON: internal error: %v", p)
        }
    }()
    // ...parser...
}
```

Gói `json` trong một thư viện chuẩn, nếu chúng gặp phải một error khi đệ quy parsing dữ liệu JSON bên trong, chúng sẽ nhanh chóng nhảy về mức cao nhất ở phía ngoài, và sau đó sẽ trả về thông điệp lỗi tương ứng.

Ngôn ngữ Go có cách hiện thực thư viện như vậy; mặc dù sử dụng package panic, chúng sẽ có thể được chuyển đổi đến một giá trị lỗi cụ thể khi một  hàm được export.

### 1.7.2 Getting the wrong context

Thỉnh thoảng rất đễ cho những user có cấp độ cao được hiểu, bên dưới sự hiện thực sẽ đóng gói lại error như là một loai error mới và trả kết quả về cho user.

```go
if _, err := html.Parse(resp.Body); err != nil {
    return nil, fmt.Errorf("parsing %s as HTML: %v", url,err)
}
```

Khi một upper user bắt gặp một lỗi, nó có thể dễ dàng để hiểu rằng lỗi đó được gây ra trong thời gian chạy từ cấp business. Nhưng rất khó để có cả hai. Khi một upper user nhận được một sai sót mới, chúng ta cũng mất những error type bên dưới (chỉ những thông tin về mô tả sẽ bị mất).

Để cần ghi nhận thông tin về kiểu lỗi trong package transition, chúng ta có thể định nghĩa một hàm `WrapError` nó sẽ gói những lỗi gốc khi bảo vệ toàn kiểu error. Để tạo điều kiện cho việc định vị vấn đến và để ghi nhận lại trạng thái lời gọi hàm, khi xảy ra lỗi, chúng ta thường muốn lưu trữ toàn bộ thông tin về lời gọi hàm hiện tại khi có lỗi xảy ra. Cùng lúc đó, để hỗ trợ network transmission như là RPC, chúng ta sẽ phải cần serialize error thành những dữ liệu tương tự như  định dạng JSON, và sau đó khôi phục lại error decoding từ dữ liệu.

Để làm việc đó, chúng ta sẽ phải tự định nghĩa cấu trúc lỗi riêng ví dụ như `github.com/chai2010/errors` với những kiểu cơ bản sau:

```go
type Error interface {
    Caller() []CallerInfo
    Wraped() []error
    Code() int
    error

    private()
}

type CallerInfo struct {
    FuncName string
    FileName string
    FileLine int
}
```

Trong số những `Error`, interface `error` là một mở rộng của kiểu interface, nó sẽ được dùng để thêm thông tin lời gọi hàm vào call stack, và hỗi trợ wrong muti-level gói lồng và hỗ trợ định dạng code. Cho tính dễ sử dụng, chúng ta có thể định nghĩa một số hàm giúp ích như sau:

```go
func New(msg string) error
func NewWithCode(code int, msg string) error

func Wrap(err error, msg string) error
func WrapWithCode(code int, err error, msg string) error

func FromJson(json string) (Error, error)
func ToJson(err error) string
```

`New` dùng để xây dựng một loại error mới tương tự như `errors.New` trong thư viện chuẩn, nhưng với việc thêm vào thông tin gọi hàm tại thời điểm gây ra error. `FromJson` được dùng để khôi phục một kiểu đối tượng error từ chuỗi JSON. `NewWithCode` nó cũng gây dựng một error với một mã error, chúng cũng có thể chứa thông tin về gọi hàm call stack. `Wrap` và `WrapWithCode` là một hàm error secondary wrapper chúng sẽ gói error như là một error mới, nhưng sẽ giữ lại message error gốc. Đối tượng error sẽ trả về từ đây và có thể trực tiếp gọi `json.Marshal`  để encode error như là JSON string.

Chúng ta có thể sử dụng wrapper function như sau:

```go
import (
    "github.com/chai2010/errors"
)

func loadConfig() error {
    _, err := ioutil.ReadFile("/path/to/file")
    if err != nil {
        return errors.Wrap(err, "read failed")
    }

    // ...
}

func setup() error {
    err := loadConfig()
    if err != nil {
        return errors.Wrap(err, "invalid config")
    }

    // ...
}

func main() {
    if err := setup(); err != nil {
        log.Fatal(err)
    }

    // ...
}
```

Ở ví dụ trên, error sẽ được bao bọc trong hai lớp. chúng ta có thể duyệt quy trình đóng gói và bỏ qua

```go
for i, e := range err.(errors.Error).Wraped() {
    fmt.Printf("wraped(%d): %v\n", i, e)
}
```

Chúng ta có thể lấy thông tin gọi hàm cho mỗi wrapper error:

```go
for i, x := range err.(errors.Error).Caller() {
    fmt.Printf("caller:%d: %s\n", i, x.FuncName)
}
```

Nếu chúng ta cần truyền một error thông qua network. chúng ta có thể encode `errors.ToJson(err)` như là JSON string

```go
// Gửi lỗi dưới dạng JSON
func sendError(ch chan<- string, err error) {
    ch <- errors.ToJson(err)
}

//  nhận lỗi dưới dạng JSON
func recvError(ch <-chan string) error {
    p, err := errors.FromJson(<-ch)
    if err != nil {
        log.Fatal(err)
    }
    return p
}
```

Cho web service dựa trên http protocol, chúng ta cũng có thể kết hợp trạng thái http với error

```go
err := errors.NewWithCode(404, "http error code")

fmt.Println(err)
fmt.Println(err.(errors.Error).Code())
```

Trong ngôn ngữ Go, error handling cũng có một coding style duy nhất. Sau khi kiểm tra chức năng phụ bị failed, chúng ta thường đặt logic code tại sao chúng failed vào process trước khi code process thành công. Nếu một error gây ra function return, sau đó logic code về thành công sẽ không được đặt trên mệnh đề `else`, nó nên được đặt trực tiếp trong body của function.

```go
f, err := os.Open("filename.ext")
if err != nil {
    // Trong trường hợp thất bại, trả về lỗi ngay lặp tức
}

// Tiếp tục xử lý nếu không có lỗi
```

Cấu trúc code của hầu hết các hàm trong ngôn ngữ Go cũng tương tự, bắt đầu bới một chuỗi khởi tạo việc kiểm tra để ngăn chặn lỗi xảy ra, theo sau bởi những logic thực sự trong function.

### 1.7.3 Incorrect error return

Error trong ngôn ngữ Go là một kiểu interface. Thông tin về interface sẽ chứa kiểu dữ liệu nguyên mâu, và kiểu dữ liệu gốc. Giá trị của interface chỉ tương ứng nếu như cả kiểu interface và giá trị gốc cả hai đều empty `nil`. Thực tế, khi kiểu của interface là empty, kiểu gốc sẽ tương ứng với interface sẽ không cần thiết phải empty.

Ví dụ sau, tôi thử cố gắng trả về một custom error type và trả về chỉ khi không có errors `nil`:


```go
func returnsError() error {
    var p *MyError = nil
    if bad() {
        p = ErrBad
    }
    return p // Will always return a non-nil error.
}
```

Tuy nhiên, kết quả trả về cuối cùng sẽ thực sự không phải `nil`, nó là một lỗi thông thường, giá trị sai là `MyError` type của con trỏ null. Sau đây là một sự cải thiện của `returnsError` :

```go
func returnsError() error {
    if bad() {
        return (*MyError)(err)
    }
    return nil
}
```

Do đó, khi đối mặt với giá trị error được return về, giá trị error return sẽ thích hợp khi được gán trực tiếp thành `nil`.

Ngôn ngữ Go sẽ có một kiểu dữ liệu mạnh, và cụ thể chuyển đổi sẽ được thực hiện giữa những kiểu khác nhau (và sẽ phải bên dưới cùng kiểu dữ liệu). Tuy nhiên, `interface` là một ngoại lệ của ngôn ngữ Go: non-interface kiểu đến kiểu interface, hoặc chuyển đổi từ interface type là cụ thể. Nó cũng sẽ hỗ trợ ducktype, dĩ nhiên, chúng sẽ thỏa mãn cấp độ 3 về bảo mật.


### 1.7.4 Parsing Exception

`panic` Support sẽ ném một `panic` liên quan đến việc ném ngoại lệ (không chỉ là error thông thường), `recover` sẽ trả về một giá trị của lời gọi hàm và `panic` cũng như thông tin về kiểu tham số của hàm
và những nguyên mẫu của hàm sẽ như sau:

```go
func panic (interface{})
func recover() interface{}
```

Luồng thông thường trong ngôn ngữ Go là kết quả trả về của việc thực thi lệnh return. Đó không phải là một exception trong luồng, do đó lường thực thi của ngoại lệ `recover` sẽ catch function trong process sẽ luôn luôn trả về  `nil`. Cái khác là ngoại lệ exception. Khi một lời gọi `panic` sẽ ném ra một ngoại lệ, function sẽ kết thúc việc thực thi lệnh con, nhưng vì lời gọi registered `defer` sẽ vấn được thực thi một cách bình thường và sau đó trả về caller. Caller trong hàm hiện tại, bởi vì trạng thái xử lý ngoại lệ chưa được bắt, `panic` sẽ tương tự như hành vi gọi hàm một cách trực tiếp. Khi một ngoại lệ xảy ra, nếu `defer` được thực thi lời gọi `recover`, nó có thể được bắt bằng việc trigger tham số  `panic, và trả về luồng thực thi bình thường.

`defer` sẽ thực hiện lệnh gọi `recover` nó thường gây khó khăn cho những người mới bắt đầu.

```go
func main() {
    if r := recover(); r != nil {
        log.Fatal(r)
    }

    panic(123)

    if r := recover(); r != nil {
        log.Fatal(r)
    }
}
```

Không trong hai lời gọi trên sẽ có thể catch exceptions. Khi lời gọi recover đầu tiên được thực thi, hàm sẽ phải được trong một thứ tự thực thi bình thường, tại một điểm mà recover có thể trả về `nil`. Khi mà một exception xảy ra, lời gọi recover thứ hai sẽ không làm thay đổi việc thực thi, bởi vị lệnh gọi `panic` sẽ gây ra `defer` hàm sẽ trả về ngay lặp tức sau khi thực thi registered function.

Trong thực tế, hàm `recover` sẽ có những yêu cầu nghiêm ngặt; chúng ta phải gọi lệnh `defer` để gọi chúng một cách trực tiếp từ hàm `recover`. Nếu hàm wrapper `defer` được gọi, `recover` sẽ catchup ngoại lệ sẽ bị fail. Ví dụ, thông thường chúng ta sẽ muốn gói hàm `MyRecover` và thêm những log cần thiết những thông tin bên trong, và sau đó gọi hàm `recover`. Đây là một hướng tiếp cận sai.

```go
func main() {
    defer func() {
        // Không thể bắt ngoại lệ
        if r := MyRecover(); r != nil {
            fmt.Println(r)
        }
    }()
    panic(1)
}

func MyRecover() interface{} {
    log.Println("trace...")
    return recover()
}
```

Một cách tương tự, nếu chúng ta gọi `defer` trong hàm nested, `recover` sẽ cũng sẽ gây ra một ngoại lệ có thể được bắt.

```go
func main() {
    defer func() {
        defer func() {
            // Không thể bắt ngoại lệ
            if r := recover(); r != nil {
                fmt.Println(r)
            }
        }()
    }()
    panic(1)
}
```

`defer` sẽ trực tiếp gọi two level nested function giống như wrapper function `recover` ở lớp một `defer` function. `MyRecover` là hàm trực tiếp trong statement sẽ work again.

```go
func MyRecover() interface{} {
    return recover()
}

func main() {
    // có thể bắt ngoại lệ bình thường
    defer MyRecover()
    panic(1)
}
```

Tuy nhiên, nếu defer statement trực tiếp gọi hàm `recover`, thì ngoại lệ sẽ không được bắt một cách phù hợp.

```go
func main(){
    defer recover()
    panic(1)
}
```

Nó sẽ phải tách biết từ stack frame với một ngoại lệ bởi stack frame, do đó hàm `recover` sẽ có thể ném một ngoại lệ một cách bình thường. Hay nói cách khác, hàm `recover` sẽ bắt ngoại lệ  của mức trên gọi hàm stack frame (chỉ là một layer `defer` function)

Dĩ nhiên, để tránh việc gọi `recover` không nhận ra được ngoại lệ, chúng ta nên tránh ném ra ngoại lệ `nil` như là một tham số.

```go
func main() {
    defer func() {
        if r := recover(); r != nil { ... }
    }()

    panic(nil)
}
```

Khi chúng ta muốn trả về việc ném ngoại lệ vào error, nếu chúng ta muốn trung thành trả về thông tin gốc, bạn sẽ phải cần sử lý chúng một cách rời rạc cho những kiểu khác nhau

```go
func foo() (err error) {
    defer func() {
        if r := recover(); r != nil {
            switch x := r.(type) {
            case string:
                err = errors.New(x)
            case error:
                err = x
            default:
                err = fmt.Errorf("Unknown panic: %v", r)
            }
        }
    }()

    panic("TODO")
}
```

Dựa trên mẫu code trên, chúng ta có thể mô phỏng nhiều kiểu exception. Bởi việc định nghĩa những kiểu khác nhau của việc bảo vệ interface, chúng ta có thể phân biệt kiểu của ngoại lệ.


```go
func main {
    defer func() {
        if r := recover(); r != nil {
            switch x := r.(type) {
            case runtime.Error:
                // ngoại lệ do quá trình chạy
            case error:
                // ngoại lệ do lỗi thông thường
            default:
                // ngoại lệ khác
            }
        }
    }()

    // ...
}
```

Nhưng làm như vậy chạy ngược lại triết lý lập trình đơn giản và dễ hiểu của Go.