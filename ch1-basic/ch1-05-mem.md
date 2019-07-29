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

## 1.5.3. Mô hình thực thi tuần tự

Xem xét ví dụ như sau:

```go
func main() {
    // in ra Hello World trên một Goroutine khác
    // tuy nhiên thông thường main sẽ chạy và kết thúc chương trình trước Goroutine
    // dẫn tới dòng chữ Hello World không được in ra
    go println("Hello World");
}
```

Có thể đồng bộ thứ tự thực thi nhờ vào channel như bên dưới:

```go
func main() { // Goroutine (1)
    // done ban đầu là một channel chưa có giá trị
    done := make(chan int)
    // chạy anynomous func trên một Goroutine khác
    go func(){ // Goroutine (2)
        // in ra dòng chữ Hello World
        println("Hello World")
        // đưa 1 vào channel, và bắt đầu dừng tới khi Goroutine (1) lấy 1 ra
        done <- 1
    }()
    // lấy ra giá trị, nên phải dừng cho tới khi channel done có giá trị để lấy ra
    <-done
}
```

## 1.5.4 Chuỗi khởi tạo

* Việc khởi tạo và thực thi chương trình Go luôn luôn bắt đầu bằng hàm `main.main`.
* Tuy nhiên nếu package `main` import các package khác vào, chúng sẽ được import theo thứ tự của string của trên file và tên thư mục.
* Nếu một package được import nhiều lần, thì những lần import sau được bỏ qua.
* Nếu các package lại import các package khác nữa, chúng sẽ import vào khởi tạo các biến theo chiều sâu như hình bên dưới:

<div align="center" width="600">
<img src="../images/ch1-12-init.ditaa.png">
<br/>
<span  align="center"><i>Quá trình khởi tạo package</i></span>
</div>
<br/>

* Nếu hàm `init` giải phóng một Goroutine mới với từ khóa `go`, thì Goroutine và `main.main` sẽ được thực thi một cách tuần tự. Bởi vì tất cả hàm `init` và hàm `main` sẽ được hoàn thành trong cùng một thread, nó cũng sẽ thoả mãn thứ tự về mô hình nhất quán.

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

Đảm bảo rằng, khi in dòng "hello world" được in ra, thread nền sẽ tiếp nhận trước khi bắt đầu `main` thread là `done <- true` trước khi gửi `<-done`, sẽ đảm bảo rằng `msg = "hello world"` được thực thi, do đó chuỗi `println(msg)` sẽ được gán rồi.


Với `buffered Channel`, đầu tiên sẽ hoàn toàn nhận `K` tác vụ trên channel xảy ra trước khi `K+C` tác vụ gửi được hoàn thành, với `C` là kích thước của buffer Channel, trước khi truyền đến Channel được hoàn thành.

Chúng ta có thể diều khiển số Gouroutine chạy concurrency dựa trên kích thước của bộ nhớ đệm control channel, ví dụ như sau

```go
// tạo ra một buffer channel có số lượng 3
var limit = make(chan int, 3)
func main() {
    // lặp qua một dãy các công việc
    for _, w := range work {
        go func() {
            // gửi 1 tới limit
            limit <- 1
            // thực thi hàm w()
            w()
            // lấy giá trị ra khỏi limit
            <-limit
        }()
    }
    // select{} sẽ dừng chương trình main
    select{}
}
```

Dòng `select{}` cuối cùng là một mệnh đề lựa chọn một empty channel sẽ làm cho main thread bị block, ngăn chặn chương trình kết thúc sớm. Tương tự `for{}` và `<- make(chan int)` nhiều hàm khác sẽ đạt được kết quả tương tự.

## 1.5.7 Tác vụ đồng bộ không tin cậy

Như chúng ta phân tích trước, đoạn code sau sẽ không đảm bảo thứ tự in ra kết quả bình thường. Việc chạy thực sự bên dưới sẽ có một xác suất lớn kết quả sẽ không bình thường.

```go
func main(){
    go println("Hello World")
}
```

Với Go, bạn có thể  đảm bảo rằng kết quả sẽ xuất ra bình thường bởi việc thêm vào thời gian sleep như sau

```go
func main(){
    go println("hello world")
    time.Sleep(time.Second)
}
```

Thread main sleep một giây sẽ đảm bảo đủ để Goroutine in ra dòng chữ "hello world" ra màn hình trước khi main thread thực thi xong. Tuy nhiên cũng không đảm bảo chắc chắn việc Goroutine thực thi xong sau một giây dẫn tới kết quả không như mong đợi.

Tính chất đúng đắn của của việc thực thi chương trình concurrency nghiêm ngặt không nên phụ thuộc vào các yếu tố không đáng tin cậy như tốc độ thực thi CPU và thời gian sleep.
