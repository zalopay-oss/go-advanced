# 1.7. Error và Exceptions

Error handling (xử lý lỗi) là một chủ đề quan trọng được đề cập trong mỗi ngôn ngữ lập trình. Go có một cơ chế xử lý lỗi đơn giản khác với `try catch` trên các ngôn ngữ lập trình khác dựa vào giá trị lỗi trả về của hàm. Ngoài ra, package `errors` giúp chúng ta định nghĩa các lỗi.

## 1.7.1. Ngữ cảnh thường gặp

Có một số hàm trong chương trình luôn yêu cầu phải chạy thành công. Ví dụ `strconv.Itoa` chuyển một số nguyên thành string, đọc và ghi phần tử từ array hoặc slice, đọc một phần tử tồn tại trong `map` và tương tự.

Những tác vụ như thế sẽ khó để có lỗi trong thời gian chạy trừ phi có `bug` trong chương trình hoặc những tình huống không thể đoán trước được như là memory leak tại thời điểm chạy. Nếu bạn thực sự bắt gặp một tình huống khác thường như thế, chương trình sẽ ngừng thực thi.

Khi một chương trình bị lỗi và dừng, chúng ta có thể cân nhắc một số khả năng xảy ra. Đối với các hàm xử lý lỗi tốt, nó sẽ trả về thêm một giá trị phụ, thường thì giá trị này được dùng để chứa thông điệp lỗi. Nếu chỉ có một lý do dẫn đến lỗi, giá trị thêm vào này có thể là một biến Boolean, thường đặt tên là `ok`. Ví dụ như bên dưới :

```go
if v, ok := m["key"]; ok {
    return v
}
```

Nhưng thông thường sẽ có nhiều hơn một nguyên nhân gây ra lỗi, và nhiều khi user muốn biết nhiều thêm về lỗi đó. Nếu bạn chỉ sử dụng một biến boolean, thì bạn sẽ không giải quyết được yêu cầu trên. Trong ngôn ngữ C, một số nguyên `errno` được sử dụng mặc định để thể hiện lỗi, do đó bạn có thể định nghĩa nhiều loại error theo nhu cầu. Trong ngôn ngữ Go, có thể gọi `syscall.Errno` giống như với mã lỗi `errno` trong ngôn ngữ C.

Ví dụ, khi chúng ta dùng `syscall.Chmod` để thay đổi `mode` của một file, chúng ta có thể thấy thông tin về lỗi qua biến `err` như bên dưới :

```go
err := syscall.Chmod(":invalid path:", 0666)
if err != nil {
    log.Fatal(err.(syscall.Errno))
}
```


Trong ngôn ngữ Go, errors được xem xét như là một kết quả đã được đoán trước; exceptions là một kết quả không thể đoán trước được, và một ngoại lệ có thể chỉ ra rằng một bug trong chương trình hoặc một vấn đề nào đó không được kiểm soát. Ngôn ngữ Go đề xuất dùng hàm `recover` để chuyển đổi exceptions thành error handling, chúng cho phép users thực sự quan tâm về những lỗi liên quan đến business.

Nếu một interface đơn giản ném tất cả những lỗi thông thường như là một ngoại lệ, chúng sẽ làm thông báo lỗi lộn xộn và không có giá trị.  Như `main` trực tiếp bao gồm mọi thứ trong một hàm, nó không mang lại ý nghĩa gì.

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

## 1.7.2. Chiến lược xử lý lỗi

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

Khi đoạn code trên chạy, sẽ tìm ẩn rủi ro. Nếu đầu tiên `os.Open` gọi thành công, nhưng lệnh gọi thứ hai `os.Create` gọi bị failed, nó sẽ trả về  ngay lặp tức mà không giải phóng tài nguyên file.

Mặc dù chúng ta có thể  fix bug bằng việc gọi `src.Close()` trước lệnh return về mệnh đề return thứ hai; nhưng khi code trở nên phức tạp hơn, những vấn đề tương tự sẽ khó để tìm thấy và giải quyết. Chúng ta có thể sử dụng mệnh đề `defer` để đảm bảo rằng một file bình thường khi được mở cũng sẽ được đóng.


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

Mệnh đề `defer` được thực thi khi ra khỏi tầm vực của hàm, chúng ta nghĩ về làm cách nào để đóng một file ngay khi mở file đó. Bất kể làm thế nào hàm được trả về, bởi về mệnh đề close có thể luôn luôn được thực thi. Cùng một thời điểm, mệnh đề defer sẽ đảm bảo rằng `io.Copy` file có thể được đóng an toàn nếu một ngoại lệ xảy ra.

Như chúng ta đã đề cập trước đó, hàm export trong ngôn ngữ Go sẽ thông thường ném ra một ngoại lệ, và một ngoại lệ không được kiểm soát có thể xem là một bug trong một chương trình. Nhưng với những framework Web services, chúng thường cần sự truy cập từ bên thứ ba ở middleware.

Bởi vì thư viện middleware thứ ba có bug, khi mà một ngoại lệ ném một exception, web framework bản thân nó không chắc chắn. Để cải thiện sự bền vững của hệ thống, web framework thường thu hồi chính xác nhất có thể những ngoại lệ trong luồng thực thi của chương trình và sau đó sẽ gây exception về bằng cách return error thông thường.

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

Ngôn ngữ Go có cách hiện thực thư viện như vậy; mặc dù sử dụng package `panic`, chúng sẽ có thể được chuyển đổi đến một giá trị lỗi cụ thể khi một hàm được export.

## 1.7.3. Trường hợp dẫn đến lỗi sai

Thỉnh thoảng rất dễ cho những upper user hiểu rằng bên dưới sự hiện thực sẽ đóng gói lại error như là một loại error mới và trả kết quả về cho user.

```go
if _, err := html.Parse(resp.Body); err != nil {
    return nil, fmt.Errorf("parsing %s as HTML: %v", url,err)
}
```

Khi một upper user bắt gặp một lỗi, nó có thể dễ dàng để hiểu rằng lỗi đó được gây ra trong thời gian chạy từ cấp business. Nhưng rất khó để có cả hai. Khi một upper user nhận được một sai sót mới, chúng ta cũng mất những error type bên dưới (chỉ những thông tin về mô tả sẽ bị mất).

Để ghi nhận thông tin về kiểu lỗi, chúng ta thông thường sẽ định nghĩa một hàm `WrapError` chúng bọc lấy lỗi gốc. Để tạo điều kiện cho những vấn đề như vậy, và để ghi nhận lại trạng thái của hàm khi một lỗi xảy ra, chúng ta sẽ muốn lưu trữ toàn bộ thông tin về hàm thực thi khi một lỗi xảy ra. Lúc này, để hỗ trợ transition như là RPC, chúng ta cần phải serialize error thành những dữ liệu tương tự như  định dạng JSON, và sau đó khôi phục lại err từ việc decoding dữ liệu.

Để làm việc đó, chúng ta sẽ phải tự định nghĩa cấu trúc lỗi riêng ví dụ như [github.com/chai2010/errors](https://github.com/chai2010/errors) với những kiểu cơ bản sau:

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

`New` dùng để tạo ra một loại error mới tương tự như `errors.New` trong thư viện chuẩn, nhưng với việc thêm vào thông tin gọi hàm tại thời điểm gây ra error. `FromJson` được dùng để khôi phục một kiểu đối tượng error từ chuỗi JSON. `NewWithCode` nó cũng gây dựng một error với một mã error, chúng cũng có thể chứa thông tin về gọi hàm call stack. `Wrap` và `WrapWithCode` là một hàm error secondary wrapper chúng sẽ gói error như là một error mới, nhưng sẽ giữ lại message error gốc. Đối tượng error sẽ trả về từ đây và có thể trực tiếp gọi `json.Marshal`  để encode error như là JSON string.

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

Ở ví dụ trên, error sẽ được bao bọc trong hai lớp. chúng ta có thể duyệt quy trình đóng gói và bỏ qua.

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

Nếu chúng ta cần truyền một error thông qua network. chúng ta có thể encode `errors.ToJson(err)` như là JSON string.

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

Cho web service dựa trên http protocol, chúng ta cũng có thể kết hợp trạng thái http với error.

```go
err := errors.NewWithCode(404, "http error code")

fmt.Println(err)
fmt.Println(err.(errors.Error).Code())
```

Trong ngôn ngữ Go, error handling cũng có một coding style duy nhất. Sau khi kiểm tra nếu chức năng phụ bị failed, chúng ta thường đặt logic code tại sao chúng failed vào process trước khi code process thành công. Nếu một error gây ra function return, sau đó logic code về thành công sẽ không được đặt trên mệnh đề `else`, nó nên được đặt trực tiếp trong body của function.

```go
f, err := os.Open("filename.ext")
if err != nil {
    // Trong trường hợp thất bại, trả về lỗi ngay lặp tức
}

// Tiếp tục xử lý nếu không có lỗi
```

Cấu trúc code của hầu hết các hàm trong ngôn ngữ Go cũng tương tự, bắt đầu bới một chuỗi khởi tạo việc kiểm tra để ngăn chặn lỗi xảy ra, theo sau bởi những logic thực sự trong function.

## 1.7.4. Trả về kết quả sai

Error trong ngôn ngữ Go là một kiểu interface. Thông tin về interface sẽ chứa kiểu dữ liệu nguyên mẫu, và kiểu dữ liệu gốc. Giá trị của interface chỉ tương ứng nếu như cả kiểu interface và giá trị gốc cả hai đều empty `nil`. Thực tế, khi kiểu của interface là empty, kiểu gốc sẽ tương ứng với interface sẽ không cần thiết phải empty.

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

## 1.7.5. Phân tích ngoại lệ

`Panic` là một hàm dựng sẵn được dùng để dừng luồng thực thi thông thường và bắt đầu `panicking`. Khi hàm `F` gọi `panic`, hàm F sẽ dừng thực thi, bất cứ hàm liên quan tới F sẽ thực thi một cách bình thường, và sau đó lệnh return F sẽ được gọi.

`Panic` được hỗ trợ để ném ra một kiểu ngoại lệ tùy ý (không chỉ là kiểu `error`), `recover` sẽ trả về một giá trị của lời gọi hàm và `panic` cũng như thông tin về kiểu tham số của hàm và những nguyên mẫu của hàm sẽ như sau:

```go
func panic (interface{})
func recover() interface{}
```

Luồng thông thường trong ngôn ngữ Go là kết quả trả về của việc thực thi lệnh return. Đó không phải là một exception trong luồng, do đó luồng thực thi của ngoại lệ `recover` sẽ catch function trong process sẽ luôn luôn trả về  `nil`. Cái khác là ngoại lệ exception. Khi một lời gọi `panic` sẽ ném ra một ngoại lệ, function sẽ kết thúc việc thực thi lệnh con, nhưng vì lời gọi registered `defer` sẽ vấn được thực thi một cách bình thường và sau đó trả về caller. Caller trong hàm hiện tại, bởi vì trạng thái xử lý ngoại lệ chưa được bắt, panic sẽ tương tự như hành vi gọi hàm một cách trực tiếp. Khi một ngoại lệ xảy ra, nếu `defer` được thực thi lời gọi `recover`, nó có thể được bắt bằng việc trigger tham số  panic, và trả về luồng thực thi bình thường.

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

Cả hai lời gọi trên không có thể catch exceptions. Khi lời gọi recover đầu tiên được thực thi, hàm sẽ phải được trong một thứ tự thực thi bình thường, tại một điểm mà recover có thể trả về `nil`. Khi mà một exception xảy ra, lời gọi recover thứ hai sẽ không làm thay đổi việc thực thi, bởi vị lệnh gọi `panic` sẽ gây ra `defer` hàm sẽ trả về ngay lặp tức sau khi thực thi registered function.

Trong thực tế, hàm `recover` sẽ có những yêu cầu nghiêm ngặt: chúng ta phải gọi lệnh `defer` để gọi chúng một cách trực tiếp từ hàm `recover`. Nếu hàm wrapper `defer` được gọi, `recover` sẽ catchup ngoại lệ sẽ bị fail. Ví dụ, thông thường chúng ta sẽ muốn gói hàm `MyRecover` và thêm những log cần thiết những thông tin bên trong, và sau đó gọi hàm `recover`. Đây là một hướng tiếp cận sai.

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

Nó sẽ phải tách biệt từ stack frame với một ngoại lệ bởi stack frame, do đó hàm `recover` sẽ có thể ném một ngoại lệ một cách bình thường. Hay nói cách khác, hàm `recover` sẽ bắt ngoại lệ  của mức trên gọi hàm stack frame (chỉ là một layer `defer` function).

Dĩ nhiên, để tránh việc gọi `recover` không nhận ra được ngoại lệ, chúng ta nên tránh ném ra ngoại lệ `nil` như là một tham số.

```go
func main() {
    defer func() {
        if r := recover(); r != nil { ... }
    }()

    panic(nil)
}
```

Khi chúng ta muốn trả về việc ném ngoại lệ vào error, nếu chúng ta muốn trung thành trả về thông tin gốc, bạn sẽ phải cần sử lý chúng một cách rời rạc cho những kiểu khác nhau.

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

Nhưng làm như vậy sẽ đi ngược lại với triết lý lập trình đơn giản và dễ hiểu của Go.

[Tiếp theo](../ch2-cgo/ch2-01-quick-start.md)