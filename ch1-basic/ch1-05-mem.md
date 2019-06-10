# 1.5 Điều khiển tuần tự cấu trúc vùng nhớ

Vào thời gian đầu, CPU sẽ thực thi những câu lệnh máy trên một lõi duy nhất. Thế hệ tổ tiên C của Go là tiêu biểu cho ngôn ngữ lập trình tuần tự. Thứ tự thực hiện của ngôn ngữ có nghĩa là cách mà những lệnh được thực thi, và chỉ có duy nhất một CPU được thực thi lệnh tại một thời điểm.

Với sự phát triển của công nghệ bộ xử lý, kiến trúc đơn lõi sẽ bắt gặp một điểm bottlenecks (thắt cổ chai) theo cách làm tăng tần số của bộ xử lý để cải thiện tốc độ tính toán. Hiện nay nhiều hệ thống CPU thường có tần số rơi vào khoảng 3GHZ. Sự trì trệ của việc phát triển CPU đơn lõi sẽ đem đến cơ hội phát triển CPU đa lõi. Theo đó, ngôn ngữ lập trình đầu tiên sẽ được phát triển theo hướng song song. Ngôn ngữ Go là một ngôn ngữ tuyệt vời để hỗ trợ tính đồng thời trong ngữ cảnh đa lõi và networking.

Những ngôn ngữ lập trình thông thường sẽ có nhiều mô hình đa dạng, đa luồng, thông điệp,v,v.. Theo lý thuyết, đa luồng (muti-threading) và message-based concurrent programming (lập trình đồng thời dựa trên truyền thông điệp) là tương đương. Do mô hình `muti-thread concurrency` (đa luồng) về mặt tự nhiên tương ứng với `muti-core processors` (đa xử lý) và những hệ điều hành đó sẽ cung cấp mức hệ thống để lập trình `muti-threading` và khái niệm về multi-threading dường như sẽ trực quan hơn, từ đó mô hình lập trình `muti-threading` sẽ dần dần được hấp thu. Go là một ngôn ngữ lập trình hệ thống cung cấp những thư viện tính năng chuyên biệt. Ngôn ngữ lập trình hệ thống sẽ hỗ trợ ít hơn mô hình lập trình message-based. Có thể tích hợp lý thuyết lập trình CSP đồng bộ. Nó có thể dễ dàng để bắt đầu một Goroutine với từ khóa `go`. Không giống như `Erlang`, ngôn ngữ `Go` sẽ chia sẻ vùng nhớ với nhau.

## 1.5.1 Goroutine và thread system

Goroutine là một đơn vị đồng thời của ngôn ngữ Go. Việc khởi tạo goroutine sẽ nhẹ hơn thread nhiều và thông qua từ khóa `go`. Trong việc hiện thực ngôn ngữ Go, goroutines và system thread không tương đương nhau. Sự khác nhau đó chỉ trên cơ sở định lượng, chính về sự thay đổi về lượng này sẽ dẫn đến sự thay đổi về chất của ngôn ngữ lập trình Go.

Đầu tiên kernel thread sẽ có một kích thước vùng nhớ stack cố định (thông thường vào khoảng 2MB). Vùng nhớ stack chủ yếu được dùng để lưu trữ những tham số và biến cục bộ khi chúng ta gọi đệ quy. Kích thước cố định của stack sẽ dẫn đến hai vấn đề, một là phần lớn vùng nhớ bị lãng phí khi khởi tạo nhiều thread nhưng thực tế chỉ cần một không gian stack nhỏ và vấn đề khác là rủi ro của việc stack overflow trong khi một số ít thread cần một lượng lớn không gian stack. Giải pháp cho vấn đề này chính là hoặc giảm kích thước không gian stack cấp phát cho mỗi thread và tăng không gian vùng nhớ sử dụng, hoặc tăng kích thước của stack để cho phép hàm đệ quy được gọi sâu hơn, nhưng cả hai cách này không có thể được kết hợp cùng lúc. Thay vì đó một gouroutine sẽ được bắt đầu bằng một vùng nhớ nhỏ (khoảng 2KB hoặc 4KB), và khi nó chạm ngưỡng đệ quy sâu hơn, không gian stack hiện tại là không đủ khoảng trống. Goroutine sẽ tự động tăng giảm không gian stack khi cần (kích thước tối đa của stack có thể được đạt tới 1GB). Bởi vì chi phí của việc khởi tạo là nhỏ, chúng ta có thể dễ dàng giải phóng hàng ngàn goroutines.

Bộ thực thi Go có một bộ đồng thời cho riêng nó, nó dùng một số kĩ thuật để ghép kênh M Goroutines trên N thread của hệ thống. Cơ chế định thời goroutine tương tự với kernel, nhưng bộ định thời  chỉ tập trung vào việc Goroutines trên các chương trình Go riêng biệt. Goroutine sử dụng cơ chế semi-preemptive cooperative scheduling, nó có thể gây ra việc định thời khi mà chương trình goroutines hiện tại bị block, Cùng thời điểm đó. nó sẽ switch sang user mode. Bộ định thời chỉ lưu trữ những thanh ghi cần thiết cho một vài hàm đặc biệt, và chi phí chuyển ngữ cảnh của các goroutines sẽ thấp hơn nhiều so với việc chuyển ngữ cảnh của thread hệ thống. Bộ runtime có một biến là `runtime.GOMAXPROCS` nó sẽ điều khiển số lượng system thread hiện thời chạy trên cơ chế non-blocking Goroutine thông thường.

Bắt đầu một goroutine trong go không chỉ là gọi một hàm, mà là kèm theo chi phí của việc định thời giữa các Goroutines. Những đặc điểm đó có sự ảnh hưởng lớn đến sự phổ biến và phát triển của lập trình đồng thời.

### 1.5.2 Toán tử Atomic

Tác vụ atomic là những tác vụ nhỏ nhất và không thể song song được trong lập trình đồng thời. Về mặt chung, nếu nhiều tác vụ được thực thi đồng thời trên cùng một tài nguyên là atomic, sau đó nhiều nhất một thực thể có thể truy cập vào một tài nguyên. Từ góc độ thread, những thread khác không thể cùng truy cập vào tài nguyên. Tác vụ atomic trong mô hình lập trình đồng thời sẽ không khác nhau nhiều với mô hình single thread, và sự tương thích này đối với việc chia sẻ resource sẽ được đảm bảo.
Thông thường sẽ có một vài lệnh CPU đặc biệt giúp bảo vệ vùng nhớ này. chúng ta có thể dùng `sync.Mutex` để đạt được điều đó.

```go
import (
    "sync"
)

var total struct {
    sync.Mutex
    value int
}

func worker(wg *sync.WaitGroup) {
    defer wg.Done()

    for i := 0; i <= 100; i++ {
        total.Lock()
        total.value += i
        total.Unlock()
    }
}

func main() {
    var wg sync.WaitGroup
    wg.Add(2)
    go worker(&wg)
    go worker(&wg)
    wg.Wait()

    fmt.Println(total.value)
}
```

[>> mã nguồn](../examples/chapter1/ch1.5/2-atomic-operation/example-1/main.go)

Trong vòng lặp của `worker`, theo thứ tự sẽ đảm bảo `total.value+=i` được đơn nguyên, chúng ta dùng `sync.Mutex` đẻ đảm bảo rằng mệnh đề chỉ được truy cập  bởi một thread trong cùng một thời điểm bằng cơ chế locking và unlocking. Trong chương trình với mô hình mutithread, rất cần thiết để lock và unlock trước và sau khi truy nhập vào vùng critical section. Với không có sự bảo vệ biến `total` , kết quả cuối cùng có thể bị sai khác do sự truy nhập đồng thời của nhiều thread.

Sử dụng mutex để bảo vệ biến số học được chia sẻ chung là một cách cồng kềnh và không hiệu quả. Thư viện chuẩn `sync/atomic` sẽ cùng cấp một gói với sự hỗ trợ cho toán tử đơn nguyên. Chúng ta có thể hiện thực lại đoạn code trên như sau.

```go
import (
    "sync"
    "sync/atomic"
)

var total uint64

func worker(wg *sync.WaitGroup) {
    defer wg.Done()

    var i uint64
    for i = 0; i <= 100; i++ {
        atomic.AddUint64(&total, i)
    }
}

func main() {
    var wg sync.WaitGroup
    wg.Add(2)

    go worker(&wg)
    go worker(&wg)
    wg.Wait()
}
```
[>> mã nguồn](../examples/chapter1/ch1.5/2-atomic-operation/example-2/main.go)


Hàm `atomic.AddUint64` khi được gọi sẽ đảm bảo rằng biến `total` được đọc và cập nhật và lưu trữ như một tác vụ đơn nguyên, do đó việc truy cập bởi nhiều thread được an toàn.


Tác vụ đơn nguyên với mutex có thể đạt được một cách hiệu quả trong một chế độ duy nhất. Phi phí của mutex sẽ cao hơn nhiều so với một biến interger bình thường. bạn có thể cộng một số numeric với một hiệu suất cao, để thay thể hiệu suất bằng việc làm giảm số lượng mutex cùng lock bởi việc bảo vệ tác vụ đơn nguyên.

```go
type singleton struct {}

var (
    instance    *singleton
    initialized uint32
    mu          sync.Mutex
)

func Instance() *singleton {
    if atomic.LoadUint32(&initialized) == 1 {
        return instance
    }

    mu.Lock()
    defer mu.Unlock()

    if instance == nil {
        defer atomic.StoreUint32(&initialized, 1)
        instance = &singleton{}
    }
    return instance
}
```

[>> mã nguồn](../examples/chapter1/ch1.5/2-atomic-operation/example-3/main.go)

Chúng ta có thể  trích xuất phần code trên trở thành `sync.One` bằng việc hiện thực lại thư viện chuẩn như sau.

```go
type Once struct {
    m    Mutex
    done uint32
}

func (o *Once) Do(f func()) {
    if atomic.LoadUint32(&o.done) == 1 {
        return
    }

    o.m.Lock()
    defer o.m.Unlock()

    if o.done == 0 {
        defer atomic.StoreUint32(&o.done, 1)
        f()
    }
}
```

[>> mã nguồn](../examples/chapter1/ch1.5/2-atomic-operation/example-4/main.go)

Dựa trên `sync.One` chúng ta sẽ hiện thực lại chế độ single piece như sau:

```go
var (
    instance *singleton
    once     sync.Once
)

func Instance() *singleton {
    once.Do(func() {
        instance = &singleton{}
    })
    return instance
}
```

[>> mã nguồn](../examples/chapter1/ch1.5/2-atomic-operation/example-5/main.go)

`sync/atomic` package này sẽ hỗ trợ những tác vụ atomic cho những kiểu cơ bản và cho việc đọc và ghi một đối tượng phức tạp, `atomic.Value` sẽ hỗ trợ hai hàm `Load` và `Store` hai hàm làm việc load và save dữ liệu, trả về giá trị và tham số là `interface{}` nó có thể được sử dụng trong một vài kiểu đặc biệt.

```go
var config atomic.Value // Toán tử atomic

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

[>> mã nguồn](../examples/chapter1/ch1.5/2-atomic-operation/example-6/main.go)

Đó là một mô hình producer và comsumer. Bên dưới thread sẽ sinh ra thông tin cấu hình gần nhất; ở phía front-end sẽ có nhiều worker thread để lấy thông tin cấu hình gần nhất.

## 1.5.3 Mô hình thống nhất chuỗi vùng nhớ

Nếu bạn muốn đồng bộ dữ liệu giữa các thread, tác vụ atomic sẽ cung câp một vài cơ chế đồng bộ để giúp cho người lập trình, Tuy nhiên, sự đảm bảo đó cũng có một tiền đề: một chuỗi mô hình  consistency memory (` sequential consistency memory model`) . Để hiểu thứ tự của chúng, hãy làm một ví dụ nhỏ

```go
var a string
var done bool

func setup() {
    a = "hello, world"
    done = true
}

func main() {
    go setup()
    for !done {}
    print(a)
}
```

[>> mã nguồn](../examples/chapter1/ch1.5/2-atomic-operation/example-7/main.go)

Chúng ta sẽ tạo ra một set up thread cho việc khởi tạo chuỗi ban đầu khởi tạo cờ `done` theo sau tác vụ khởi tạo là true. Trong thread main, nơi mà hàm được lưu giữ, khi mà câu lệnh `for !done{}` kiểm tra biến done có thể chuyển thành true, nó có thể được xem xét như tác vụ khởi tạo string được hoàn thành, sau đó một kí tự trong string sẽ được in ra.

Tuy nhiên, ngôn ngữ Go sẽ không đảm bảo rằng việc ghi trong main sẽ được xem xét là `done` xảy ra sau khi phép toán ghi của string `a`, do đó chương trình sẽ như là in ra một chuỗi rỗng. Để làm vấn đề tệ hơn, bởi vì không có một cơ chế đồng bộ sự kiện giữa hai thread, và hàm main sẽ rơi vào trạng thái lặp không giới hạn.


Trong ngôn ngữ Go, một cấu trúc vùng nhớ liên tục sẽ được đảm bảo trong cùng một Goroutine thread. tuy nhiên giữa những Goroutine khác nhau, tính chất đồng bộ của chuỗi nhớ sẽ không được đảm bảo. và một cách định nghĩa đồng bộ sự kiện sẽ cần thiết để tăng tối đa tính song song, bộ biên dịch Go sẽ biên dịch và bộ xử lý sẽ sắp xếp lại thứ thự các lệnh mà không ảnh hưởng đến những quy luật trên (CPU sẽ biểu diễn một vài lệnh ngoài thứ tự đó)

Do đó, nếu `a=1;b=2` hai mệnh đề trên sẽ được thực hiện tuần tự trong goroutine, mặc dù `a=1` hay là `b=2` được thực thi trước. Những sự thay đổi đó không theo dự đoán trước. Nếu chương trình đồng bộ không thể được xác đinh dựa vào thứ tự các mối liên hệ của sự kiện, kết quả của chương trình sẽ không chắn chắn, ví dụ bên dưới

```go
func main() {
    go println("你好, 世界");
}
```

Theo đặc tả của ngôn ngữ Go, hàm main sẽ kết thúc và khi hàm kết thúc nó sẽ không đợi bất kỳ quá trình nào chạy nền bên dưới. Bởi vì việc thực thi goroutine trong hàm main sẽ trả về một sự kiện là đồng thời, bất cứ phần nào cũng có thể chạy trước. Do đó, khi in ra màn hình, bất cứ khi nào chúng in ra là không biết.

Sử dụng tác vụ atomic trước không giúp giải bài toán trên bởi vì chúng ta không xác định thứ tự của hai phép toán atomic, Hướng giải quyết của vấn đề này là cụ thể cho chúng chạy theo thứ tự nhờ vào việc cơ chế bên dưới,

```go
func main() {
    done := make(chan int)

    go func(){
        println("你好, 世界")
        done <- 1
    }()

    <-done
}
```

[>> mã nguồn](../examples/chapter1/ch1.5/2-atomic-operation/example-8/main.go)

Khi mà `<-done` được thực thi, thì những yêu cầu không thể thay thế `done <- 1` sẽ được hiện thực. Theo như trong cùng một goroutine sẽ thỏa mãn quy luật nhất quán. Chúng ta có thể nói rằng khi `done <- 1` được thực thi, thì mệnh đề `println()` sẽ được thực thi trước rồi,  Do đó chương trình hiện tại sẽ có kết quả được in ra màn hình bình thường.

Dĩ nhiên, cơ chế đồng bộ của `sync.Mutex` sẽ có thể đạt được thông qua `Mutex`

```go
func main() {
    var mu sync.Mutex

    mu.Lock()
    go func(){
        println("你好, 世界")
        mu.Unlock()
    }()

    mu.Lock()
}
```

[>> mã nguồn](../examples/chapter1/ch1.5/3-sequence-consistency-mem-model/example-9/main.go)


Có thể xác định rằng, bên dưới việc thực thi `mutex.UnLock()` sẽ phải là `println("你好, 世界")` hoàn thành trước. (một số thread thỏa mãn thứ tự nhất quán), và trong main, hàm thứ hai sẽ `mu.Lock()` sẽ phải là `mu.UnLock()` xảy ra bên dưới background thread (được đảm bảo bởi `sync.Mutex`) và bên dưới nền sẽ in ra công việc được hoàn thành một cách thành công.

### 1.5.4 Khởi tạo chuỗi

Trong chương trước, chúng ta đã được giới thiệu ngắn gọn về việc khởi tạo một chuỗi trong chương trình, nó là một số đặc điểm đặt biệt của ngôn ngữ Go theo mô hình vùng nhớ đồng thời.

Việc khởi tạo và thực thi trong chương trình Go luôn luôn bắt đầu bằng hàm `main.main`. Tuy nhiên nếu package `main` import các package khác vào, chúng sẽ được import theo thứ tự của string của trên file và tên thư mục) Nếu một package được import nhiều lần, nó chỉ được import và thực thi đúng một lần. Khi mà một package được import, nếu nó cũng import những package khác nữa, thì đầu tiên sẽ bao gồm package khác, sau đó tạo ra và khởi tạo biến và hằng của package. Sau đó hàm `init` trong package, nêu một package có nhiều hàm `init` thì việc hiện thực sẽ gọi chúng theo thứ tự file name, nhiều hàm init trong cùng một file được gọi theo thứ tự chúng xuất hiện (`init` không phải là một hàm thông thường, chúng có thể được định nghĩa nhiều lần, chúng sẽ không được gọi từ những hàm khác). Cuối cùng, package `main` biến và hằng được khai báo và khởi tạo, và hàm `init` sẽ được thực thi trước khi hàm thực thi `main.main`. Chương trình bắt đầu thực thi một cách bình thường, theo sau là một sơ đồ ngữ nghĩa của việc khởi động hàm Go bên dưới.




<p align="center" width="600">
<img src="../images/ch1-12-init.ditaa.png">
<br/>
<span>Hình 1-12 Quá trình khởi tạo package</span>
</p>

Nên chú ý rằng `main.main` trong những mã nguồn sẽ được thực thi trong cùng Goroutine trong cùng một hàm mà nó thực thi, và nó cũng là việc chạy trong main thread của chương trình. Nếu hàm `init` giải phóng một Goroutine mới với từ khóa `go`, thì Goroutine và `main.main` sẽ được thực thi một cách tuần tự.

Bởi vì tất cả hàm `init` và hàm `main` sẽ được hoàn thành trong cùng một thread, nó cũng sẽ thoả mãn thứ tự về mô hình nhất quán.


### 1.5.5 Khởi tạo một Goroutine

Mệnh đề đứng trước từ khóa `go` sẽ tạo ra một Goroutine mới trước khi trả về một goroutine hiện tại, ví dụ :

```go
var a string

func f() {
    print(a)
}

func hello() {
    a = "hello, world"
    go f()
}
```

[>> mã nguồn](../examples/chapter1/ch1.5/5-create-go-routine/example-10/main.go)

Việc thực thi của `go f()` sẽ tạo ra một Goroutine, và hàm `hello` sẽ thực thi cùng lúc với Goroutine. Theo thứ tự của các statement được viết, nó có thể được xác định bằng một khi việc khởi tạo Goroutine được xảy ra, nó có thể không được sắp xếp. Nó là việc đồng thời. Việc gọi hello sẽ in ra tại một số điểm trong tương lai "hello,world", hoặc có thể là `hello` được in ra sao khi hàm đã thực thi xong

### 1.5.6 Giao tiếp thông qua kênh Channel

Giao tiếp thông qua channel là một phương pháp chính trong việc đồng bộ giữa các goroutine. Mỗi lần thực hiện thao tác gửi trên một `unbufferred Channel` thường đi đôi với tác vụ nhận. Tác vụ gửi và nhận thường xảy ra ở những Goroutine khác nhau (hai tác vụ diễn ra trên cùng một goroutine có thể dễ dàng dẫn đến deadlocks). **Tác vụ gửi trên một unbufferred Channel luôn luôn xảy ra trước khi tác vụ nhận hoàn thành**.

```go
var done = make(chan bool)
var msg string

func aGoroutine() {
    msg = "你好, 世界"
    done <- true
}

func main() {
    go aGoroutine()
    <-done
    println(msg)
}
```

[>> mã nguồn](../examples/chapter1/ch1.5/6-channel-base-com/example-11/main.go)


Cũng đảm bảo rằng, khi in dòng "hello, world". Vì thread nền sẽ tiếp nhận trước khi bắt đầu `main` thread là `done <- true` trước khi gửi `<-done`, sẽ đảm bảo rằng `msg = "hello, world"` được thực thi, do đó chuỗi `println(msg)` sẽ được gán rồi. Tóm lại, bên thread nền sẽ đầu tiên ghi vào biến `msg`, sau đó sẽ nhận tín hiệu từ `done`, theo sau bởi `main` là một thread để truyền tín hiệu tương ứng với lần thực thi hàm `println(msg)` kết thúc. Tuy nhiên, nếu Channel được buffered (ví dụ, `done = make(chan bool, 1)` ), main thread sẽ nhận tác vụ `done <- true` sẽ blocked cho đến khi thread nền nhận, và chương trình sẽ không đảm bảo in ra dòng chữ "hello, world".

Với `buffered Channel`, đầu tiên sẽ hoàn toàn nhận `K` tác vụ trên channel xảy ra trước khi `K+C` tác vụ gửi được hoàn thành, với `C` là kích thước của buffer Channel, trước khi truyền đến Channel được hoàn thành.

Chúng ta có thể diều khiển số Gouroutine chạy đồng thời dựa trên kích thước của bộ nhớ đệm control channel, ví dụ như sau

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

[>> mã nguồn](../examples/chapter1/ch1.5/6-channel-base-com/example-13/main.go)


Dòng `select{}` cuối cùng là một mệnh đề lựa chọn một empty pipe sẽ làm cho main thread bị block, ngăn chặn chương trình kết thúc sớm. Tương tự `for{}` và `<- make(chan int)` nhiều hàm khác sẽ đạt được kết quả tương tự. Bởi vì thread main sẽ bị blocked. nó có thể là `os.Exit(0)` được hiện thực nếu chương trình cần kết thúc một cách thông thường.

### 1.5.7 Tác vụ đồng bộ không tin cậy

Như chúng ta phân tích trước, đoạn code sau sẽ không đảm bảo thứ tự in ra kết quả bình thường. Việc chạy thực sự bên dưới sẽ có một xác suất lớn kết quả sẽ không bình thường.

```go
func main(){
    go println("Hello, World")
}
```

Chỉ liên hệ với Go, bạn có thể  đảm bảo rằng kết quả sẽ xuất ra bình thường bởi việc thêm vào thời gian sleep như sau

```go
func main(){
    go println("hello, world")
    time.Sleep(time.Second)
}
```

Bởi vì thread main sleep một giây, chương trình sẽ có xác suất lớn rằng kết quả được in ra một cách bình thường. Do đó, nhiều người sẽ cảm thấy rằng chương trình sẽ không còn là một vấn đề. Nhưng chương trình này sẽ không ổn đi và đó sẽ vẫn dẫn đến failure. Đầu tiên hãy giả sử rằng chương trình có thể được ổn định kết quả đầu ra. Bởi vì việc bắt đầu thực thi thì thread Go sẽ không bị blocking, thread `main`
 sẽ cụ thể sleep một giây và chương trình sẽ kết thúc. Chúng ta có thể giả sử rằng chương trình sẽ thực nhiều hơn một giây. Bây giờ giả sử hàm `println` sẽ sleep lâu hơn main thread bị sleep. Nó có thể dẫn đến hai mặt đối lập sau: do bên dưới thread nền main thread sẽ kết thúc trước khi việc in ra hoàn thành, thời gian thực thi sẽ nhỏ hơn thời gian thực thi của thread chính. Dĩ nhiên điều đó là hoàn toàn có thể.

Tính chất đúng đắn của của việc thực thi chương trình đồng thời nghiêm ngặt không nên phụ thuộc vào các yếu tố không đáng tin cậy như tốc độ thực thi CPU và thời gian ngủ. Đồng thời, cũng có thể lấy được kết quả tĩnh, theo tính chất nhất quán  của đơn hàng trong luồng, kết hợp với khả năng sắp xếp của các sự kiện đồng bộ hóa kênh, hoặc đồng bộ hóa sự kiện dẫn xuất. Nếu hai sự kiện không thể được sắp xếp theo quy tắc đó, sau đó là thực thi đồng thời, do đó việc thực thi sẽ không tin cậy.

Ý tưởng của việc giải quyết thực thi đồng thời cũng giống nhau: cụ thể sử dụng cơ chế đồng bộ.

