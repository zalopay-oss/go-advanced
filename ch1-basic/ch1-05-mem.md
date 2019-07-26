# 1.5. Concurrent-oriented memory model

Thời gian đầu, CPU chỉ có một lõi duy nhất, các ngôn ngữ khi đó sẽ theo mô hình lập trình tuần tự, điển hình là ngôn ngữ C. Ngày nay, với sự phát triển của công nghệ đa xử lý, để tận dụng tối đa sức mạnh của CPU, mô hình lập trình song song hay [multi-threading](https://en.wikipedia.org/wiki/Multithreading_(computer_architecture)) thường thấy trên các ngôn ngữ lập trình ra đời. Ngôn ngữ Go cũng phát triển mô hình lập trình song song rất hiệu quả với khái niệm Goroutines.

<div align="center">
	
Lập trình tuần tự|Lập trình song song
---|---
![](../images/ch1-5-sequence-programming.png) | ![](../images/ch1-5-parallelprograming.png)


</div>
<br/>


## 1.5.1. Goroutines và system threads

Goroutines là một đơn vị concurrency của ngôn ngữ Go. Việc khởi tạo goroutines sẽ ít tốn chi phí hơn khởi tạo `thread` nhiều và đơn giản thông qua từ khóa `go`. Về góc nhìn hiện thực, `goroutines` và `system thread` không giống nhau.

Đầu tiên, system thread sẽ có một kích thước vùng nhớ stack cố định (thông thường vào khoảng 2MB). Vùng nhớ stack chủ yếu được dùng để lưu trữ những tham số, biến cục bộ và địa chỉ trả về khi chúng ta gọi hàm.

Kích thước cố định của stack sẽ dẫn đến hai vấn đề:
  * Lãng phí vùng nhớ đối với chương trình đơn giản
  * StackOverflow với những chương trình gọi hàm phức tạp.

Giải pháp cho vấn đề này chính là cấp phát linh hoạt vùng nhớ stack:
  * Một Goroutines sẽ được bắt đầu bằng một vùng nhớ nhỏ (khoảng 2KB hoặc 4KB).
  * Khi gọi đệ quy sâu (không gian stack hiện tại là không đủ) Goroutines sẽ tự động tăng không gian stack (kích thước tối đa của stack có thể được đạt tới 1GB)
  * Bởi vì chi phí của việc khởi tạo là nhỏ, chúng ta có thể dễ dàng giải phóng hàng ngàn goroutines.

Bộ thực thi (runtime) Go có riêng cơ chế định thời cho Goroutines, nó dùng một số kỹ thuật để ghép M Goroutines trên N thread của hệ thống. Cơ chế định thời Goroutines tương tự với cơ chế định thời của `kernel` nhưng chỉ ở mức chương trình. Biến `runtime.GOMAXPROCS` quy định số lượng system thread hiện thời chạy trên các Goroutines.

## 1.5.2. Tác vụ Atomic

Tác vụ atomic là những tác vụ nhỏ nhất và không thể chạy song song được trong lập trình concurrency. Tác vụ atomic trên một vùng nhớ chia sẻ đảm bảo vùng nhớ chỉ có thể được truy cập bởi một Goroutine tại một thời điểm. Để đạt được điều này ta có thể dùng `sync.Mutex`.

```go
import (
// package cần dùng
    "sync"
)

// total là một struct atomic
var total struct {
    sync.Mutex
    value int
}

func worker(wg *sync.WaitGroup) {
    // thông báo tôi đã hoàn thành khi ra khỏi hàm
    defer wg.Done()

    for i := 0; i <= 100; i++ {
        // chặn các Goroutines khác vào
        total.Lock()
        // bây giờ, lệnh total.value += i được đảm bảo là atomic (đơn nguyên)
        total.value += i
        // bỏ chặn
        total.Unlock()
    }
}

func main() {
    // khai báo wg để main Goroutine dừng chờ các Goroutines khác trước khi kết thúc chương trình
    var wg sync.WaitGroup
    // wg cần chờ 2 Goroutines khác
    wg.Add(2)
    // thực thi Goroutines thứ nhất
    go worker(&wg)
    // thực thi Goroutines thứ hai
    go worker(&wg)
    // wg bắt đầu đợi để 2 Goroutines kia xong
    wg.Wait()
    // in ra kết quả thực thi để xem chính xác không
    fmt.Println(total.value)
}
```

Trong chương trình với mô hình `multithread`, rất cần thiết để `lock` và `unlock` trước và sau khi truy cập vào vùng [critical section](https://en.wikipedia.org/wiki/Critical_section). Nếu không có sự bảo vệ biến `total` , kết quả cuối cùng có thể bị sai khác do sự truy nhập đồng thời của nhiều thread.

Sử dụng `mutex` chỉ để bảo vệ một biến số học là cách làm phức tạp và không hiệu quả, thay vào đó chúng ta nên dùng package `sync/atomic`:

```go
import (
    "sync"
    // khai báo biến gói sync/atomic
    "sync/atomic"
)

// biến total được truy cập đồng thời
var total uint64

func worker(wg *sync.WaitGroup) {
    // wg thông báo hoàn thành khi ra khỏi hàm
    defer wg.Done()

    var i uint64
    for i = 0; i <= 100; i++ {
        // lệnh cộng atomic.AddUint64 total được đảm bảo là atomic (đơn nguyên)
        atomic.AddUint64(&total, i)
    }
}

func main() {
    // wg được dùng để dừng hàm main đợi các Goroutines khác
    var wg sync.WaitGroup
    // wg cần đợi hai Goroutines gọi lệnh Done() mới thực thi tiếp
    wg.Add(2)
    // tạo Goroutines thứ nhất
    go worker(&wg)
    // tạo Goroutines thứ hai
    go worker(&wg)
    // bắt đầu việc đợi
    wg.Wait()
    // in ra kết quả
    fmt.Println(total)
}
```

Ví dụ bên dưới minh họa cho việc sử dụng `mutex` và `sync/atomic` để hiện thực singleton pattern.

```go
// khai báo một struct singleton
type singleton struct {}

var (
    // khai báo một đối tượng singleton
    instance    *singleton
    // khai báo một số atomic
    initialized uint32
    // dùng mutex để lock và unlock
    mu          sync.Mutex
)

func Instance() *singleton {
    // nếu giá trị của initialized là 1, tức đối tượng đã được khởi tạo trước đó thì trả về nó
    if atomic.LoadUint32(&initialized) == 1 {
        return instance
    }
    // lock vùng critical section
    mu.Lock()
    // unlock khi ra khỏi hàm
    defer mu.Unlock()
    // bằng nil là chưa được khởi tạo, khác nil thì có Goroutines khởi tạo rồi
    if instance == nil {
        // lưu initialized là 1 để đánh dấu đã khởi tạo
        defer atomic.StoreUint32(&initialized, 1)
        // khởi tạo duy nhất 1 lần từ nay trở về sau
        instance = &singleton{}
    }
    // trả về instance đã được khởi tạo
    return instance
}
```

Chúng ta có thể  refactor phần code trên thành `sync.One` như sau:

```go
// Once là một struct atomic
type Once struct {
    m    Mutex
    done uint32
}
// Hàm Do đảm bảo f được thực thi một lần duy nhất
func (o *Once) Do(f func()) {
    // nếu giá trị o.done là 1, ta trả về ngay
    if atomic.LoadUint32(&o.done) == 1 {
        return
    }
    // lock các Goroutines khác
    o.m.Lock()
    // unlock nếu Goroutines hiện tại thực thi xong
    defer o.m.Unlock()
    // nếu o.done là 0 giống như mới khởi tạo
    if o.done == 0 {
        // lưu trữ o.done là 1 để đánh dấu
        defer atomic.StoreUint32(&o.done, 1)
        // thực thi hàm f()
        f()
    }
}
```

Dựa trên `sync.One` chúng ta sẽ hiện thực lại chế độ single piece như sau:

```go
var (
    instance *singleton
    once     sync.Once
)

func Instance() *singleton {
    // thủ tục được truyền vào once.Do sẽ thực thi một lần duy nhất
    once.Do(func() {
        instance = &singleton{}
    })
    return instance
}
```

* Package `sync/atomic` sẽ hỗ trợ những tác vụ atomic cho những kiểu cơ bản.
* Cho việc đọc và ghi một đối tượng phức tạp, `atomic.Value` sẽ hỗ trợ hai hàm `Load` và `Store` để load và save dữ liệu, trả về giá trị và tham số là `interface{}` nó có thể được sử dụng trong một vài kiểu đặc biệt.

```go
var config atomic.Value
// Lưu giá trị vào atomic
config.Store(loadConfig())

// định thời sau mỗi giây sẽ cập nhật lại
go func() {
    for {
        time.Sleep(time.Second)
        config.Store(loadConfig())
    }
}()

// Giải phóng 10 thread để lấy giá trị
for i := 0; i < 10; i++ {
    go func() {
        for r := range requests() {
            c := config.Load()
        }
    }()
}
```

Đó là một mô hình producer và comsumer. Bên dưới thread sẽ sinh ra thông tin cấu hình gần nhất. Phía front-end sẽ có nhiều worker thread để lấy thông tin cấu hình gần nhất.

## 1.5.3. Mô hình thực thi tuần tự nhất quán

Ví dụ bên dưới minh họa cho việc điều khiển thứ tự thực thi giữa các Goroutines:

```go
// biến string
var a string
// cờ done cho biết trạng thái thực thi
var done bool
// hàm này phải được chạy trước
func setup() {
    a = "hello world"
    // biến done để đánh dấu thực thi xong
    done = true
}

func main() {
    // chạy setup
    go setup()
    // thực thi vòng lặp busy waiting để chờ setup thực thi xong
    for !done {}
    // in ra giá trị
    print(a)
}
```

Tuy nhiên, cách làm này có vấn đề khi không đảm bảo rằng việc ghi trong main sẽ được xem xét là `done` xảy ra sau khi phép toán ghi của string `a`, bởi vì cấu trúc vùng nhớ liên tục được đảm bảo trong cùng một Goroutines nhưng không được đảm bảo khi khác Goroutines

Trong ngôn ngữ Go, một cấu trúc vùng nhớ liên tục sẽ được đảm bảo trong cùng một Goroutine thread. tuy nhiên giữa những Goroutine khác nhau, tính chất đồng bộ của chuỗi nhớ sẽ không được đảm bảo. và một cách định nghĩa đồng bộ sự kiện sẽ cần thiết để tăng tối đa tính song song, bộ biên dịch Go sẽ biên dịch và bộ xử lý sẽ sắp xếp lại thứ thự các lệnh mà không ảnh hưởng đến những quy luật trên (CPU sẽ biểu diễn một vài lệnh ngoài thứ tự đó)

Do đó, nếu `a=1;b=2` hai mệnh đề trên sẽ được thực hiện tuần tự trong goroutine, mặc dù `a=1` hay là `b=2` được thực thi trước. Những sự thay đổi đó không theo dự đoán trước. Nếu chương trình đồng bộ không thể được xác đinh dựa vào thứ tự các mối liên hệ của sự kiện, kết quả của chương trình sẽ không chắn chắn, ví dụ bên dưới

```go
func main() {
    go println("Hello World");
}
```

Theo đặc tả của ngôn ngữ Go, hàm main sẽ kết thúc và khi hàm kết thúc nó sẽ không đợi bất kỳ quá trình nào chạy nền bên dưới. Bởi vì việc thực thi goroutine trong hàm main sẽ trả về một sự kiện là concurrency, bất cứ phần nào cũng có thể chạy trước. Do đó, khi in ra màn hình, bất cứ khi nào chúng in ra là không biết.

Sử dụng tác vụ atomic trước không giúp giải bài toán trên bởi vì chúng ta không xác định thứ tự của hai phép toán atomic, Hướng giải quyết của vấn đề này là cụ thể cho chúng chạy theo thứ tự nhờ vào việc cơ chế bên dưới,

```go
func main() {
    done := make(chan int)

    go func(){
        println("Hello World")
        done <- 1
    }()

    <-done
}
```

Khi mà `<-done` được thực thi, thì những yêu cầu không thể thay thế `done <- 1` sẽ được hiện thực. Theo như trong cùng một goroutine sẽ thỏa mãn quy luật nhất quán. Chúng ta có thể nói rằng khi `done <- 1` được thực thi, thì mệnh đề `println()` sẽ được thực thi trước rồi,  Do đó chương trình hiện tại sẽ có kết quả được in ra màn hình bình thường.

Dĩ nhiên, cơ chế đồng bộ của `sync.Mutex` sẽ có thể đạt được thông qua `Mutex`

```go
func main() {
    var mu sync.Mutex

    mu.Lock()
    go func(){
        println("Hello World")
        mu.Unlock()
    }()

    mu.Lock()
}
```

Có thể xác định rằng, bên dưới việc thực thi `mutex.UnLock()` sẽ phải là `println("Hello World")` hoàn thành trước. (một số thread thỏa mãn thứ tự nhất quán), và trong main, hàm thứ hai sẽ `mu.Lock()` sẽ phải là `mu.UnLock()` xảy ra bên dưới background thread (được đảm bảo bởi `sync.Mutex`) và bên dưới nền sẽ in ra công việc được hoàn thành một cách thành công.

## 1.5.4 Khởi tạo chuỗi

Trong chương trước, chúng ta đã được giới thiệu ngắn gọn về việc khởi tạo một chuỗi trong chương trình, nó là một số đặc điểm đặt biệt của ngôn ngữ Go theo mô hình vùng nhớ concurrency.

Việc khởi tạo và thực thi trong chương trình Go luôn luôn bắt đầu bằng hàm `main.main`. Tuy nhiên nếu package `main` import các package khác vào, chúng sẽ được import theo thứ tự của string của trên file và tên thư mục) Nếu một package được import nhiều lần, nó chỉ được import và thực thi đúng một lần. Khi mà một package được import, nếu nó cũng import những package khác nữa, thì đầu tiên sẽ bao gồm package khác, sau đó tạo ra và khởi tạo biến và hằng của package. Sau đó hàm `init` trong package, nêu một package có nhiều hàm `init` thì việc hiện thực sẽ gọi chúng theo thứ tự file name, nhiều hàm init trong cùng một file được gọi theo thứ tự chúng xuất hiện (`init` không phải là một hàm thông thường, chúng có thể được định nghĩa nhiều lần, chúng sẽ không được gọi từ những hàm khác). Cuối cùng, package `main` biến và hằng được khai báo và khởi tạo, và hàm `init` sẽ được thực thi trước khi hàm thực thi `main.main`. Chương trình bắt đầu thực thi một cách bình thường, theo sau là một sơ đồ ngữ nghĩa của việc khởi động hàm Go bên dưới.

<div align="center" width="600">
<img src="../images/ch1-12-init.ditaa.png">
<br/>
<span  align="center"><i>Quá trình khởi tạo package</i></span>
</div>
<br/>

Nên chú ý rằng `main.main` trong những mã nguồn sẽ được thực thi trong cùng Goroutine trong cùng một hàm mà nó thực thi, và nó cũng là việc chạy trong main thread của chương trình. Nếu hàm `init` giải phóng một Goroutine mới với từ khóa `go`, thì Goroutine và `main.main` sẽ được thực thi một cách tuần tự.

Bởi vì tất cả hàm `init` và hàm `main` sẽ được hoàn thành trong cùng một thread, nó cũng sẽ thoả mãn thứ tự về mô hình nhất quán.


## 1.5.5 Khởi tạo một Goroutine

Mệnh đề đứng trước từ khóa `go` sẽ tạo ra một Goroutine mới trước khi trả về một goroutine hiện tại, ví dụ :

```go
var a string

func f() {
    print(a)
}

func hello() {
    a = "hello world"
    go f()
}
```

Việc thực thi của `go f()` sẽ tạo ra một Goroutine, và hàm `hello` sẽ thực thi cùng lúc với Goroutine. Theo thứ tự của các statement được viết, nó có thể được xác định bằng một khi việc khởi tạo Goroutine được xảy ra, nó có thể không được sắp xếp. Nó là việc concurrency. Việc gọi hello sẽ in ra tại một số điểm trong tương lai "hello,world", hoặc có thể là `hello` được in ra sao khi hàm đã thực thi xong

## 1.5.6 Giao tiếp thông qua kênh Channel

Giao tiếp thông qua channel là một phương pháp chính trong việc đồng bộ giữa các goroutine. Mỗi lần thực hiện thao tác gửi trên một `unbufferred Channel` thường đi đôi với tác vụ nhận. Tác vụ gửi và nhận thường xảy ra ở những Goroutine khác nhau (hai tác vụ diễn ra trên cùng một goroutine có thể dễ dàng dẫn đến deadlocks). **Tác vụ gửi trên một unbufferred Channel luôn luôn xảy ra trước khi tác vụ nhận hoàn thành**.

```go
var done = make(chan bool)
var msg string

func aGoroutine() {
    msg = "Hello World"
    done <- true
}

func main() {
    go aGoroutine()
    <-done
    println(msg)
}
```

Cũng đảm bảo rằng, khi in dòng "hello world". Vì thread nền sẽ tiếp nhận trước khi bắt đầu `main` thread là `done <- true` trước khi gửi `<-done`, sẽ đảm bảo rằng `msg = "hello world"` được thực thi, do đó chuỗi `println(msg)` sẽ được gán rồi. Tóm lại, bên thread nền sẽ đầu tiên ghi vào biến `msg`, sau đó sẽ nhận tín hiệu từ `done`, theo sau bởi `main` là một thread để truyền tín hiệu tương ứng với lần thực thi hàm `println(msg)` kết thúc. Tuy nhiên, nếu Channel được buffered (ví dụ, `done = make(chan bool, 1)` ), main thread sẽ nhận tác vụ `done <- true` sẽ blocked cho đến khi thread nền nhận, và chương trình sẽ không đảm bảo in ra dòng chữ "hello world".

Với `buffered Channel`, đầu tiên sẽ hoàn toàn nhận `K` tác vụ trên channel xảy ra trước khi `K+C` tác vụ gửi được hoàn thành, với `C` là kích thước của buffer Channel, trước khi truyền đến Channel được hoàn thành.

Chúng ta có thể diều khiển số Gouroutine chạy concurrency dựa trên kích thước của bộ nhớ đệm control channel, ví dụ như sau

```go
var limit = make(chan int, 3)

func main() {
    for _, w := range work {
        go func() {
            limit <- 1
            w()
            <-limit
        }()
    }
    select{}
}
```

Dòng `select{}` cuối cùng là một mệnh đề lựa chọn một empty pipe sẽ làm cho main thread bị block, ngăn chặn chương trình kết thúc sớm. Tương tự `for{}` và `<- make(chan int)` nhiều hàm khác sẽ đạt được kết quả tương tự. Bởi vì thread main sẽ bị blocked. nó có thể là `os.Exit(0)` được hiện thực nếu chương trình cần kết thúc một cách thông thường.

## 1.5.7 Tác vụ đồng bộ không tin cậy

Như chúng ta phân tích trước, đoạn code sau sẽ không đảm bảo thứ tự in ra kết quả bình thường. Việc chạy thực sự bên dưới sẽ có một xác suất lớn kết quả sẽ không bình thường.

```go
func main(){
    go println("Hello World")
}
```

Chỉ liên hệ với Go, bạn có thể  đảm bảo rằng kết quả sẽ xuất ra bình thường bởi việc thêm vào thời gian sleep như sau

```go
func main(){
    go println("hello world")
    time.Sleep(time.Second)
}
```

Bởi vì thread main sleep một giây, chương trình sẽ có xác suất lớn rằng kết quả được in ra một cách bình thường. Do đó, nhiều người sẽ cảm thấy rằng chương trình sẽ không còn là một vấn đề. Nhưng chương trình này sẽ không ổn đi và đó sẽ vẫn dẫn đến failure. Đầu tiên hãy giả sử rằng chương trình có thể được ổn định kết quả đầu ra. Bởi vì việc bắt đầu thực thi thì thread Go sẽ không bị blocking, thread `main` sẽ cụ thể sleep một giây và chương trình sẽ kết thúc. Chúng ta có thể giả sử rằng chương trình sẽ thực nhiều hơn một giây. Bây giờ giả sử hàm `println` sẽ sleep lâu hơn main thread bị sleep. Nó có thể dẫn đến hai mặt đối lập sau: do bên dưới thread nền main thread sẽ kết thúc trước khi việc in ra hoàn thành, thời gian thực thi sẽ nhỏ hơn thời gian thực thi của thread chính. Dĩ nhiên điều đó là hoàn toàn có thể.

Tính chất đúng đắn của của việc thực thi chương trình concurrency nghiêm ngặt không nên phụ thuộc vào các yếu tố không đáng tin cậy như tốc độ thực thi CPU và thời gian ngủ. concurrency, cũng có thể lấy được kết quả tĩnh, theo tính chất nhất quán  của đơn hàng trong luồng, kết hợp với khả năng sắp xếp của các sự kiện đồng bộ hóa kênh, hoặc đồng bộ hóa sự kiện dẫn xuất. Nếu hai sự kiện không thể được sắp xếp theo quy tắc đó, sau đó là thực thi concurrency, do đó việc thực thi sẽ không tin cậy.

Ý tưởng của việc giải quyết thực thi concurrency cũng giống nhau: cụ thể sử dụng cơ chế đồng bộ.
