# 1.6. Mô hình thực thi đồng thời

Một điểm mạnh của Golang là tích hợp sẵn cơ chế xử lý đồng thời (concurrency). Lý thuyết về hệ thống tương tranh của Go là CSP (Communicating Sequential Process) được đề xuất bởi Hoare vào năm 1978. CSP  được áp dụng lần đầu cho máy tính đa dụng T9000 mà Hoare có tham gia. Từ NewSqueak, Alef, Limbo đến Golang hiện tại, Rob Pike, người có hơn 20 năm kinh nghiệm thực tế với CSP, rất quan tâm  đến tiềm năng áp dụng CSP vào ngôn ngữ lập trình đa dụng. Khái niệm cốt lõi của lý thuyết CSP cũng được áp dụng vào lập trình concurrency trong Go.

<div align="center">

<img src="../images/gophercomplex5.jpg" width="800">
<br/>
<span align="center"><i>Go concurrency</i></span>
    <br/>

</div>

Ở phần này chúng ta cùng xem qua các cách dùng goroutine cũng như các cách xử lý goroutine mà chúng ta thường gặp trong lúc lập trình.

## 1.6.1. Phiên bản concurrency với Hello World

Trong hầu hết các ngôn ngữ hiện đại, vấn đề chia sẻ tài nguyên được giải quyết bằng cơ chế đồng bộ hóa như khóa (lock) nhưng Golang có cách tiếp cận riêng là chia sẻ giá trị thông qua channel.

<div align="center">
<img src="../images/channel.jpg" width="400">
    <br/>
    <span align="center"><i>Goroutine trao đổi giá trị qua channel</i></span>
    <br/>

</div>

Trên thực tế khi nhiều thread thực thi độc lập chúng hiếm khi chủ động chia sẻ tài nguyên. Tại bất kỳ thời điểm nào, tốt nhất là chỉ Goroutine sở hữu tài nguyên của chính mình. Golang có một triết lý được thể hiện bằng slogan:

> *Do not communicate by sharing memory; instead, share memory by communicating.*
>
> *Do not communicate through shared memory, but share memory through communication.*

Mặc dù các vấn đề tương tranh đơn giản như   tham chiếu đến biến đếm có thể được hiện thực bằng  `atomic operations` hoặc `mutex lock`, nhưng việc kiểm soát truy cập thông qua Channel giúp cho code của chúng ta clean và "Golang" hơn.

### Áp dụng Mutex

Xem xét đoạn code sau:

```go
func main() {
    var mu sync.Mutex

    // khởi chạy một goroutine chạy đồng thời với main
    go func(){
        fmt.Println("Hello World")

        // không thể đảm bảo hàm này sẽ chạy trước `mu.Unlock()` phía dưới
        mu.Lock()
    }()

    // xảy ra lỗi vì `mu` đang ở trạng thái unlocked
    mu.Unlock()
}
```

Ở đây, `mu.Lock()` và `mu.Unlock()` không ở trong cùng một Goroutine, vì vậy nó không đáp ứng được mô hình bộ nhớ nhất quán tuần tự (sequential consistency memory model).

Sửa lại đoạn code trên như sau:

```go
func main() {
    var mu sync.Mutex

    // lệnh này khiến main thread bị block
    mu.Lock()

    go func(){
        fmt.Println("Hello World")

        // lệnh này Unlock cho main thread
        mu.Unlock()
    }()

    mu.Lock()
}
```

### Áp dụng Channel

Đồng bộ hóa với mutex là một cách tiếp cận ở mức độ tương đối đơn giản. Bây giờ ta sẽ sử dụng một unbuffered channel để hiện thực đồng bộ hóa:  

```go
func main() {
    done := make(chan int)

    go func(){
        fmt.Println("Hello World")

        // chỉ sau khi goroutine này hoàn thành thao tác nhận
        // thì thao tác gửi ở main thread mới được kết thúc
        <-done
    }()

    // thao tác gửi qua channel `done`
    done <- 1
}
```

Cách này gặp bất cập với buffered channel  vì lúc đó không có gì đảm bảo rằng goroutine sẽ in ra trước khi thoát `main`. Cách tiếp cận tốt hơn là hoán đổi hướng gửi và nhận của channel để tránh các sự kiện đồng bộ hóa bị ảnh hưởng bởi kích thước buffer của nó:  

```go
func main() {
    done := make(chan int, 1)

    go func(){
        fmt.Println("Hello World")

        // gửi 1 giá trị vào channel thông báo kết thúc goroutine này
        done <- 1
    }()

    // main thread nhận giá trị từ channel và thoát khỏi
    // trạng thái block
    <-done
}
```

Dựa trên buffered channel, chúng ta có thể dễ dàng mở rộng thread print đến N. Ví dụ sau là mở 10 goroutine để in riêng biệt:  

```go
func main() {
    done := make(chan int, 10)

    // mở ra N goroutine
    for i := 0; i < cap(done); i++ {
        go func(){
            fmt.Println("Hello World")
            done <- 1
        }()
    }

    // đợi cả 10 goroutine hoàn thành
    for i := 0; i < cap(done); i++ {
        <-done
    }
}
```

### Sử dụng sync.WaitGroup thay cho Channel

Một cách đơn giản hơn là sử dụng `sync.WaitGroup` để chờ một tập các sự kiện:

```go
func main() {
    var wg sync.WaitGroup

    // mở N goroutine
    for i := 0; i < 10; i++ {
        // tăng số lượng sự kiện chờ, hàm này phải được
        // đảm bảo thực thi trước khi bắt đầu 1 goroutine chạy nền
        wg.Add(1)

        go func() {
            fmt.Println("Hello World")

            // cho biết hoàn thành một sự kiện
            wg.Done()
        }()
    }

    // đợi N goroutine hoàn thành
    wg.Wait()
}
```

## 1.6.2. Tác vụ Atomic

[Tác vụ atomic](https://preshing.com/20130618/atomic-vs-non-atomic-operations/) trên một vùng nhớ chia sẻ thì đảm bảo rằng vùng nhớ đó chỉ có thể được truy cập bởi một Goroutine tại một thời điểm. Để đạt được điều này ta có thể dùng [sync.Mutex](https://golang.org/pkg/sync/#Mutex).

### Sử dụng sync.Mutex

```go
import (
// package cần dùng
    "sync"
)

// total là một atomic struct
var total struct {
    sync.Mutex
    value int
}

func worker(wg *sync.WaitGroup) {
    // thông báo hoàn thành khi ra khỏi hàm
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
    // in ra kết quả thực thi
    fmt.Println(total.value)
}
```

Trong một chương trình đồng thời, ta cần có cơ chế để `lock` và `unlock` trước và sau khi truy cập vào vùng [critical section](https://en.wikipedia.org/wiki/Critical_section). Nếu không có sự bảo vệ biến `total` , kết quả cuối cùng có thể bị sai khác do sự truy cập đồng thời của nhiều thread.

### Sử dụng sync/atomic

Thay vì dùng mutex, chúng ta cũng có thể dùng package [sync/atomic](https://golang.org/pkg/sync/atomic/), đây là giải pháp hiệu quả hơn đối với một biến số học.

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

Để ghi và đọc atomic trên những đối tượng phức tạp hơn thì ta dùng kiểu [atomic.Value](https://golang.org/pkg/sync/atomic/#Value), ví dụ:

```go
package main

import (
    "sync/atomic"
    "time"
)

func loadConfig() map[string]string {
    return make(map[string]string)
}

func requests() chan int {
    return make(chan int)
}

func main() {
    // nắm giữ thông tin cấu hình của server
    var config atomic.Value
    // khởi tạo giá trị ban đầu
    config.Store(loadConfig())
    go func() {
    // cập nhật thông tin sau mỗi 10 giây
        for {
            time.Sleep(10 * time.Second)
            config.Store(loadConfig())
        }
    }()
    // tạo nhiều worker sử lý request
    // dùng thông tin cấu hình gần nhất
    for i := 0; i < 10; i++ {
        go func() {
            for r := range requests() {
                c := config.Load()
                // xử lý request với cấu hình c
                _, _ = r, c
            }
        }()
    }
}

```


## 1.6.3. Mô hình Producer Consumer

<div align="center">

<img src="../images/producer-consumer.png" width="800">
<br/>
<span align="center"><i>Mô hình Producer - Consumer</i></span>
    <br/>

</div>

Ví dụ phổ biến nhất về lập trình concurrency là mô hình Producer Consumer, giúp tăng tốc độ xử lý chung của chương trình bằng cách cân bằng sức mạnh của các thread "sản xuất" (produce) và "tiêu thụ" (consume).

Producer tạo ra một số dữ liệu và sau đó đưa nó vào hàng đợi, cùng lúc đó consumer cũng lấy dữ liệu từ hàng đợi này ra để xử lý. Điều này làm cho produce và consume trở thành hai quá trình bất đồng bộ. Khi không có dữ liệu trong hàng đợi kết quả, consumer sẽ chờ đợi ở trạng thái "đói", còn khi dữ liệu trong hàng đợi bị đầy, producer phải đối mặt với vấn đề mất mát dữ liệu khi CPU phải loại bỏ bớt dữ liệu trong đó để nạp thêm.

Golang hiện thực cơ chế này khá đơn giản:

```go
// producer: liên tục tạo ra một chuỗi số nguyên dựa trên bội số factor và đưa vào channel
func Producer(factor int, out chan<- int) {
    for i := 0; ; i++ {
        out <- i*factor
    }
}

// consumer: liên tục lấy các số từ channel ra để print
func Consumer(in <-chan int) {
    for v := range in {
        fmt.Println(v)
    }
}
func main() {
    // hàng đợi
    ch := make(chan int, 64)

    // tạo một chuỗi số với bội số 3
    go Producer(3, ch)

    // tạo một chuỗi số với bội số 5
    go Producer(5, ch)

    // tạo consumer
    go Consumer(ch)

    // thoát ra sau khi chạy trong một khoảng thời gian nhất định
    time.Sleep(5 * time.Second)
}
```

Chúng ta có thể để hàm `main` giữ trạng thái block mà không thoát và chỉ  thoát khỏi chương trình khi người dùng gõ `Ctrl-C`:

```go
func main() {
    // hàng đợi
    ch := make(chan int, 64)

    go Producer(3, ch)
    go Producer(5, ch)
    go Consumer(ch)

    // Ctrl+C để thoát
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    fmt.Printf("quit (%v)\n", <-sig)
}
```

Có 2 producer trong ví dụ trên và không có sự kiện đồng bộ nào giữa hai producer mà chúng concurrency. Do đó, thứ tự của chuỗi output ở consumer là không xác định.

## 1.6.4. Mô hình Publish Subscribe

Mô hình publish-and-subscribe thường được viết tắt là mô hình pub/sub. Trong mô hình này, producer trở thành publisher và consumer  trở thành subscriber, đồng thời producer:consumer là mối quan hệ M:N.

<div align="center">

<img src="../images/pubsub.png" width="800">
<br/>
<span align="center"><i>Mô hình Publish - Subscribe</i></span>
    <br/>

</div>

Trong mô hình producer-consumer truyền thống, thông điệp được gửi đến hàng đợi và mô hình publish-subscription sẽ publish thông điệp đến một topic.

Để hiện thực mô hình này ta implement package `pubsub`:

```go
// package pubsub implements a simple multi-topic pub-sub library.
package pubsub

import (
    "sync"
    "time"
)

type (
    // subscriber thuộc kiểu channel
    subscriber chan interface{}

    // topic là một filter
    topicFunc func(v interface{}) bool
)

type Publisher struct {
    // Read/Write Mutex
    m sync.RWMutex

    // kích thước  hàng đợi
    buffer int

    // timeout cho việc publishing
    timeout time.Duration

    // subscriber đã subscribe vào topic nào
    subscribers map[subscriber]topicFunc
}

// constructor với timeout và độ dài hàng đợi
func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
    return &Publisher{
        buffer:      buffer,
        timeout:     publishTimeout,
        subscribers: make(map[subscriber]topicFunc),
    }
}

// thêm subscriber mới, đăng ký hết tất cả topic
func (p *Publisher) Subscribe() chan interface{} {
    return p.SubscribeTopic(nil)
}

// thêm subscriber mới, subscribe các topic đã được filter lọc
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
    ch := make(chan interface{}, p.buffer)
    p.m.Lock()
    p.subscribers[ch] = topic
    p.m.Unlock()
    return ch
}

// hủy subscribe
func (p *Publisher) Evict(sub chan interface{}) {
    p.m.Lock()
    defer p.m.Unlock()

    delete(p.subscribers, sub)
    close(sub)
}

// publish ra 1 topic
func (p *Publisher) Publish(v interface{}) {
    p.m.RLock()
    defer p.m.RUnlock()

    var wg sync.WaitGroup
    for sub, topic := range p.subscribers {
        wg.Add(1)
        go p.sendTopic(sub, topic, v, &wg)
    }
    wg.Wait()
}

// đóng 1 đối tượng publisher và đóng tất cả các subscriber
func (p *Publisher) Close() {
    p.m.Lock()
    defer p.m.Unlock()

    for sub := range p.subscribers {
        delete(p.subscribers, sub)
        close(sub)
    }
}

// gửi 1 topic có thể duy trì trong thời gian chờ wg
func (p *Publisher) sendTopic(
    sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup,
) {
    defer wg.Done()
    if topic != nil && !topic(v) {
        return
    }

    select {
    case sub <- v:
    case <-time.After(p.timeout):
    }
}
```

Trong ví dụ sau đây, 2 subscriber đăng ký hết tất cả các topic với "golang":

```go
import (
    "./pubsub"
    "time"
    "strings"
    "fmt"
)
func main() {
    // khởi tạo 1 publisher
    p := pubsub.NewPublisher(100*time.Millisecond, 10)

    // để đảm bảo p được đóng trước khi exit
    defer p.Close()

    // `all` subscribe hết tất cả topic
    all := p.Subscribe()

    // subscribe các topic có "golang"
    golang := p.SubscribeTopic(func(v interface{}) bool {
        if s, ok := v.(string); ok {
            return strings.Contains(s, "golang")
        }
        return false
    })

    // publish ra 2 topic
    p.Publish("hello,  world!")
    p.Publish("hello, golang!")

    // print những gì subscriber `all` nhận được
    go func() {
        for  msg := range all {
            fmt.Println("all:", msg)
        }
    } ()

    // print những gì subscriber `golang` nhận được
    go func() {
        for  msg := range golang {
            fmt.Println("golang:", msg)
        }
    } ()

    // thoát ra sau khi chạy 3 giây
    time.Sleep(3 * time.Second)
}
```

Trong mô hình pub/sub, mỗi thông điệp được gửi tới nhiều subscriber. Publisher thường không biết hoặc không quan tâm subscriber nào nhận được thông điệp. Subscriber và publisher có thể được thêm vào động ở thời điểm thực thi, cho phép các hệ thống phức tạp có thể phát triển theo thời gian. Trong thực tế, những ứng dụng như dự báo thời tiết có thể áp dụng mô hình concurrency này.

## 1.6.5. Kiểm soát số lượng goroutine

Goroutine là một tính năng mạnh mẽ của Go, mất chi phí rất ít để sử dụng, những tất nhiên nếu dùng với số lượng quá lớn sẽ chiếm gây nhiều lãng phí và cần có một cơ chế để kiểm soát. Một cách thông dụng để đạt được mục đích trên là dùng worker pool.

<div align="center">

<img src="../images/worker-pool.png" width="500">
<br/>
<span align="center"><i>Mô hình Worker pool</i></span>
    <br/>

</div>

Đầu tiên tạo ra các worker:

```go
func worker(queue chan int, worknumber int, done chan bool) {
    for j := range queue {
        fmt.Println("worker", worknumber, "finished job", j)
        done <- true
    }
}
```

Sau đó có thể áp dụng như sau:

```go
func main() {

    // queue of jobs
    q := make(chan int)

    // done channel lấy ra kết quả của jobs
    done := make(chan bool)

    // số lượng worker trong pool
    numberOfWorkers := 4
    for i := 0; i < numberOfWorkers; i++ {
        go worker(q, i, done)
    }

    // đưa job vào queue
    numberOfJobs := 17
    for j := 0; j < numberOfJobs; j++ {
        go func(j int) {
            q <- j
        }(j)
    }

    // chờ nhận đủ kết quả
    for c := 0; c < numberOfJobs; c++ {
        <-done
    }
}
```

## 1.6.6. Dọn dẹp Goroutine

Sau khi job queue rỗng, ta sẽ phải dừng tất cả worker. Goroutine dù khá nhẹ nhưng vẫn không phải miễn phí, nhất là với các hệ thống lớn, dù chỉ là các chi phí nhỏ nhất cũng có thể trở nên khác biệt lớn nếu thay đổi.

Cách đơn giản là dùng kill channel để phát ra tín hiệu ngừng cho goroutine.

```go
func main() {
    // channel để terminate các worker
    killsignal := make(chan bool)

    // queue các jobs
    q := make(chan int)
    // done channel nhận vào kết quả của các job
    done := make(chan bool)

    // số lượng worker trong pool
    numberOfWorkers := 4
    for i := 0; i < numberOfWorkers; i++ {
        go worker(q, i, done, killsignal)
    }

    // đưa job vào queue
    numberOfJobs := 17
    for j := 0; j < numberOfJobs; j++ {
        go func(j int) {
            q <- j
        }(j)
    }

    // chờ để nhận đủ kết quả
    for c := 0; c < numberOfJobs; c++ {
        <-done
    }

    // dọn dẹp các worker
    close(killsignal)
    time.Sleep(2 * time.Second)
}
```

Trong đó các worker được thiết kế như sau:

```go
func worker(queue chan int, worknumber int, done, ks chan bool) {
    for true {
        // dùng select để chờ cùng lúc trên cả 2 channel
        select {
        // xử lý job trong channel queue
        case k := <-queue:
            fmt.Println("doing work!", k, "worknumber", worknumber)
            done <- true

        // nếu nhận được kill signal thì return
        case <-ks:
            fmt.Println("worker halted, number", worknumber)
            return
        }
    }
}
```

## 1.6.7. Sàng số nguyên tố

Trong phần ***1.2***, chúng tôi đã trình bày việc triển khai phiên bản concurrency của sàng số nguyên tố để chứng minh tính concurrency của Newsqueak. Nguyên tắc "sàng số nguyên tố" như sau:

<div align="center">

<img src="../images/ch1-13-prime-sieve.png">
<br/>
<span align="center"><i>Sàng số nguyên tố</i></span>
    <br/>

</div>

Chúng ta cần khởi tạo một chuỗi các số tự nhiên `2, 3, 4, ...` (không bao gồm 0, 1):

```go
// trả về channel tạo ra chuỗi số: 2, 3, 4, ...
func GenerateNatural() chan int {
    ch := make(chan int)
    go func() {
        for i := 2; ; i++ {
            ch <- i
        }
    }()
    return ch
}
```

Tiếp theo xây dựng một sàng cho mỗi số nguyên tố: đề xuất một số là bội số của số nguyên tố trong chuỗi đầu vào và trả về một chuỗi mới, đó là một channel mới.

```go
// bộ lọc: xóa các số có thể chia hết cho số nguyên tố
func PrimeFilter(in <-chan int, prime int) chan int {
    out := make(chan int)
    go func() {
        for {
            if i := <-in; i%prime != 0 {
                out <- i
            }
        }
    }()
    return out
}
```

Bây giờ ta có thể sử dụng bộ lọc này trong hàm `main`:

```go
func main() {
    // chuỗi số: 2, 3, 4, ...
    ch := GenerateNatural()
    for i := 0; i < 100; i++ {
        // số nguyên tố mới
        prime := <-ch

        // Bộ lọc dựa trên số nguyên tố mới
        fmt.Printf("%v: %v\n", i+1, prime)
        ch = PrimeFilter(ch, prime)

        // dựa trên chuỗi số còn lại trong channel để lọc
        // các số nguyên tố tiếp theo với các số được
        // trích xuất dưới dạng filter. Các channel tương ứng
        // với các sàng số nguyên tố khác nhau được kết nối liên tiếp nhau.
    }
}
```

## 1.6.8. Kẻ thắng làm vua

Có nhiều động lực để lập trình concurrency nhưng tiêu biểu là vì lập trình concurrency có thể đơn giản hóa các vấn đề. Lập trình concurrency cũng có thể cải thiện hiệu năng. Mở hai thread trên CPU đa lõi thường nhanh hơn mở một thread.  Trên thực tế về mặt cải thiện hiệu suất, chương trình không chỉ đơn giản là chạy nhanh, mà trong nhiều trường hợp chương trình có thể đáp ứng yêu cầu của người dùng một cách nhanh chóng là điều quan trọng nhất. Khi không có yêu cầu từ người dùng cần xử lý, nên xử lý một số tác vụ nền có độ ưu tiên thấp.

Giả sử chúng ta muốn nhanh chóng tìm kiếm các chủ đề liên quan đến "golang", có thể mở nhiều công cụ tìm kiếm như Bing, Google hoặc Yahoo. Khi tìm kiếm trả về kết quả trước, ta có thể đóng các trang tìm kiếm khác. Do ảnh hưởng của môi trường mạng và thuật toán của công cụ tìm kiếm mà một số công cụ tìm kiếm có thể trả về kết quả tìm kiếm nhanh hơn. Chúng ta có thể sử dụng một chiến lược tương tự để viết chương trình này:  

```go
func main() {
    // tạo ra một channel với buffer đủ lớn để đảm bảo
    // không bị block do kích thước của buffer.
    ch := make(chan string, 32)

    // chạy nhiều goroutine dưới nền và gửi yêu cầu
    // tìm kiếm đến các công cụ tìm kiếm khác nhau
    go func() {
        ch <- searchByBing("golang")
    }()
    go func() {
        ch <- searchByGoogle("golang")
    }()
    go func() {
        ch <- searchByBaidu("golang")
    }()

    // khi bất kỳ công cụ tìm kiếm nào có kết quả
    // nó sẽ ngay lập tức gửi kết quả đến channel
    // ta chỉ lấy kết quả đầu tiên từ channel
    fmt.Println(<-ch)
}
```

Áp dụng ý tưởng trên có thể giúp cải thiện hiệu suất bằng cách chọn lấy kẻ chiến thắng trong cuộc đua thời gian.

## 1.6.9. Context package

Ở thời điểm phát hành Go1.7, thư viện tiêu chuẩn đã thêm một package context để đơn giản hóa hoạt động của dữ liệu, thời gian chờ và thoát giữa nhiều Goroutines. Package context định nghĩa kiểu Context, chứa deadline, cancelation signal và các giá trị request-scope giữa các API và giữa các process.

Chúng ta có thể sử dụng package context để hiện thực lại cơ chế kiểm soát timeout:

```go
func worker(ctx context.Context, wg *sync.WaitGroup) error {
    defer wg.Done()

    for {
        select {
        default:
            fmt.Println("hello")
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func main() {
    // nhận vào context parent (Background) và trả về context child (ctx) và hàm cancel
    // deadline 10 secs
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    var wg sync.WaitGroup

    // đơn giản hoá worker pool cho ngắn gọn
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go worker(ctx, &wg)
    }

    time.Sleep(time.Second)

    // mặc dù ctx sẽ expire theo timeout đã set trước đó
    // ta vẫn gọi cancel để đóng context child và các children của nó
    // để tránh giữ chúng tồn tại không cần thiết
    cancel()

    // sử dụng waitGroup thay cho done channel
    wg.Wait()
}
```

Golang tự động lấy lại   bộ nhớ, do đó bộ nhớ thường không bị rò rỉ (memory leak). Trong ví dụ trước về sàng số nguyên tố, một Goroutine mới  được đưa vào bên trong hàm `GenerateNatural` và Goroutine nền `PrimeFilter` có nguy cơ bị leak khi hàm `main` không còn sử dụng channel. Chúng ta có thể tránh vấn đề này với package context. Dưới đây là phần triển khai sàng số nguyên tố được cải thiện:  

```go
// trả về channel có chuỗi số: 2, 3, 4, ...
func GenerateNatural(ctx context.Context) chan int {
    ch := make(chan int)
    go func() {
        for i := 2; ; i++ {
            select {
            case <- ctx.Done():
                return
            case ch <- i:
            }
        }
    }()
    return ch
}

// bộ lọc: xóa các số có thể chia hết cho số nguyên tố
func PrimeFilter(ctx context.Context, in <-chan int, prime int) chan int {
    out := make(chan int)
    go func() {
        for {
            if i := <-in; i%prime != 0 {
                select {
                case <- ctx.Done():
                    return
                case out <- i:
                }
            }
        }
    }()
    return out
}

func main() {
    // kiểm soát trạng thái Goroutine nền thông qua context
    ctx, cancel := context.WithCancel(context.Background())

    ch := GenerateNatural(ctx) // chuỗi số: 2, 3, 4, ...
    for i := 0; i < 100; i++ {
        prime := <-ch // số nguyên tố mới
        fmt.Printf("%v: %v\n", i+1, prime)
        ch = PrimeFilter(ctx, ch, prime) // Bộ lọc dựa trên số nguyên tố mới
    }

    cancel()
}
```

Khi hàm `main` kết thúc hoạt động, nó được thông báo bằng lệnh `cancel()` gọi đến Goroutine nền để thoát, do đó tránh khỏi việc leak Goroutine.

Concurrency là một chủ đề rất lớn, và ở đây chúng tôi chỉ đưa ra một vài ví dụ về lập trình concurrency rất cơ bản. Tài liệu chính thức cũng có rất nhiều cuộc thảo luận về lập trình concurrency, có khá nhiều  cuốn sách  thảo luận cụ thể về lập trình concurrency trong Golang. Độc giả có thể tham khảo các tài liệu liên quan theo nhu cầu của mình.

[Tiếp theo](ch1-07-error-and-panic.md)