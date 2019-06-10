# 1.6 Các chế độ đồng thời thông dụng (Concurrency Mode)

Khía cạnh hấp dẫn nhất của ngôn ngữ Go có tích hợp sẵn cơ chế xử lý đồng thời (concurrency). Lý thuyết về hệ thống tương tranh của Go là CSP (Communicating Sequential Process) được đề xuất bởi CAR Hoare vào năm 1978. CSP  được áp dụng lần đầu cho máy tính đa dụng T9000 mà Hoare có tham gia. Từ NewSqueak, Alef, Limbo đến ngôn ngữ Go hiện tại, Rob Pike, người có hơn 20 năm kinh nghiệm thực tế với CSP, rất quan tâm  đến tiềm năng áp dụng CSP vào ngôn ngữ lập trình đa dụng. Khái niệm cốt lõi của lý thuyết CSP cũng là  của lập trình đồng thời trong Go: giao tiếp đồng bộ (synchronous communication). Chủ đề về giao tiếp đồng bộ đã được đề cập trong phần trước. Trong phần này chúng tôi sẽ giới thiệu ngắn gọn về các mẫu đồng thời phổ biến trong ngôn ngữ Go.

Điều đầu tiên cần làm rõ là khái niệm: "đồng thời" không phải "song song" (concurrency is not parallel). Đồng thời quan tâm nhiều hơn ở cấp độ thiết kế của chương trình. Các chương trình đồng thời có thể được thực thi tuần tự và chỉ trên các CPU đa lõi thực sự mới có thể chạy cùng một lúc. Song song quan tâm nhiều hơn đến cấp độ thực thi của chương trình. Song song cơ bản là lặp lại một số lần rất lớn các vòng lặp đơn giản. Ví dụ như một số lượng lớn các thao tác  song song về xử lý hình ảnh được thực thi trong GPU. Với mục đích viết chương trình chạy đồng thời hiệu quả hơn, từ khi bắt đầu thiết kế ngôn ngữ Go đã tập trung vào cách thiết kế một mô hình trừu tượng đơn giản, an toàn và hiệu quả ở cấp độ ngôn ngữ lập trình, cho phép các lập trình viên tập trung vào giải quyết vấn đề và kết hợp các giải thuật mà không phải quá chú tâm vào việc quản lý các thread và tín hiệu.

Trong lập trình đồng thời, việc truy cập đúng vào tài nguyên được chia sẻ đòi hỏi sự kiểm soát chính xác. Trong hầu hết các ngôn ngữ hiện đại, vấn đề khó khăn này được giải quyết bằng cơ chế đồng bộ hóa như khóa (lock) nhưng ngôn ngữ Go có cách tiếp cận riêng là chia sẻ Giá trị (Value) được truyền qua channel (trên thực tế, nhiều thread thực thi độc lập hiếm khi chủ động chia sẻ tài nguyên). Tại bất kỳ thời điểm nào, tốt nhất là chỉ có một Goroutine để sở hữu tài nguyên. Cạnh tranh dữ liệu đã được loại bỏ từ cấp độ thiết kế. Để thúc đẩy lối suy nghĩ này, ngôn ngữ Go đã triết lý hóa chương trình đồng thời của nó thành một slogan:

> Đừng giao tiếp bằng cách chia sẻ bộ nhớ, thay vào đó hãy chia sẻ bộ nhớ bằng cách giao tiếp.
> Đừng giao tiếp thông qua bộ nhớ chia sẻ, nhưng hãy chia sẻ bộ nhớ thông qua giao tiếp.

Đây là một cấp độ cao hơn của triết lý lập trình đồng thời (truyền các giá trị qua pipeline   luôn được Go khuyến nghị). Mặc dù các vấn đề tương tranh đơn giản như   tham chiếu đến biến đếm có thể được hiện thực bằng  *atomic operations* hoặc  *mutex lock*, nhưng việc kiểm soát truy cập thông qua Channel giúp cho code của chúng ta đúng và "sạch" hơn.

## 1.6.1 Phiên bản đồng thời của *Hello world*

Trước tiên, ta in ra "Hello world" trong Goroutine mới và đợi cho output của thread thread nền (chứ Goroutine này) kế thúc và thoát, chương trình với cơ chế đồng thời đơn giản này sẽ được thực thi.

Khái niệm cốt lõi của lập trình đồng thời là giao tiếp đồng bộ, nhưng có nhiều cách để đồng bộ hóa. Trước tiên chúng ta dùng `sync.Mutex` để giao tiếp đồng bộ với cơ chế mutex quen thuộc . Theo document, chúng ta không thể trực tiếp mở khóa `sync.Mutex` ở trạng thái đã mở khóa , điều này có thể gây ra runtime exceptions. Cách sau đây sẽ không thực thi bình thường được: [(source)](../examples/ch1.6/1-hello-world-concurrent-ver/example-1/main.go)

```go
func main() {
    var mu sync.Mutex

    go func(){
        fmt.Println("Hello World")
        mu.Lock()
    }()

    mu.Unlock()
}
```

Bởi vì `mu.Lock()` và `mu.Unlock()` không ở trong cùng một Goroutine, vì vậy nó không đáp ứng được mô hình bộ nhớ nhất quán tuần tự. Đồng thời, chúng không có sự kiện đồng bộ hóa nào khác để tham chiếu tới. Hai sự kiện này vì thế không thể thực thi đồng thời. Bởi vì khi chúng đồng thời, rất có khả `mu.Unlock()` trong `main` sẽ thực thi trước và tại thời điểm này, mutex `mu` vẫn ở trạng thái mở khóa, điều này sẽ gây ra runtime exceptions.

Sau đây là đoạn code đã sửa:  [(source)](../examples/ch1.6/1-hello-world-concurrent-ver/example-2/main.go)

```go
func main() {
    var mu sync.Mutex

    mu.Lock()
    go func(){
        fmt.Println("Hello World")
        mu.Unlock()
    }()

    mu.Lock()
}
```

Cách khắc phục là thực hiện 2 lần `mu.Lock()` trong hàm `main`. Khi khóa thứ hai bị block, nó sẽ bị block vì khóa đã bị chiếm (không phải là khóa đệ quy). Trạng thái block của hàm `main` khiến thread nền tiếp tục thực thi. Khi thread nền được mở khóa `mu.Unlock()`, công việc `print` được hoàn thành và việc mở khóa sẽ khiến trạng thái block của `mu.Lock()` thứ hai được hủy. Tại thời điểm này, thread nền và thread `main` không có tham chiếu sự kiện đồng bộ hóa nào khác, và sự kiện thoát của chúng sẽ là đồng thời: Khi hàm `main` thoát và chương trình thoát, thread nền có thể đã thoát hoặc không thể thoát. Mặc dù không thể xác định khi nào hai thread sẽ thoát, công việc `print` vẫn có thể được thực hiện chính xác.

Đồng bộ hóa với mutex là một cách tiếp cận ở mức độ tương đối thấp. Bây giờ ta sẽ sử dụng một pipeline không được  cache để đạt được đồng bộ hóa:  [(source)](../examples/ch1.6/1-hello-world-concurrent-ver/example-3/main.go)

```go
func main() {
    done := make(chan int)

    go func(){
        fmt.Println("Hello World")
        <-done
    }()

    done <- 1
}
```

Theo đặc tả mô hình bộ nhớ ngôn ngữ Go, thao tác nhận từ một channel không phải buffer thực hiện trước khi việc truyền tới channel được hoàn thành (kết thúc việc truyền). Do đó, sau khi thread nền hoàn thành thao tác nhận `<-done`, thao tác gửi  của thread `main` là `done <- 1` mới được hoàn thành (do đó thoát khỏi `main` rồi thoát khỏi chương trình) và công việc in được thực thi xong.

Mặc dù đoạn code trên có thể được đồng bộ đúng đắng, nhưng nó quá nhạy cảm với kích thước cache của  pipeline: nếu pipeline có cache, không có gì đảm bảo rằng thread nền sẽ in đúng trước khi thoát `main`. Cách tiếp cận tốt hơn là hoán đổi hướng gửi và nhận của pipeline để tránh các sự kiện đồng bộ hóa bị ảnh hưởng bởi kích thước cache của pipeline:  [(source)](../examples/ch1.6/1-hello-world-concurrent-ver/example-4/main.go)

```go
func main() {
    done := make(chan int, 1) // pipeline cache

    go func(){
        fmt.Println("Hello World")
        done <- 1
    }()

    <-done
}
```

Đối với buffered channel, thao tác nhận thứ `K` cho channel xảy ra trước khi hoàn thành thao tác truyền thứ `K + C`, trong đó `C` là kích thước cache của channel. Mặc dù pipeline được lưu vào cache, việc thread `main` tiếp nhận  được hoàn thành tại thời điểm khi thread nền gửi nhưng chưa hoàn thành, và thao tác `print` được hoàn thành.

Dựa trên pipeline với cache, chúng ta có thể dễ dàng mở rộng thread print đến N. Ví dụ sau là mở 10 thread nền để in riêng biệt:  [(source)](../examples/ch1.6/1-hello-world-concurrent-ver/example-5/main.go)

```go
func main() {
    done := make(chan int, 10)

    // Mở ra N thread
    for i := 0; i < cap(done); i++ {
        go func(){
            fmt.Println("Hello World")
            done <- 1
        }()
    }

    // Đợi cả 10 thread hoàn thành
    for i := 0; i < cap(done); i++ {
        <-done
    }
}
```

Một cách đơn giản để làm điều này là đợi N thread hoàn thành trước khi tiến hành thao tác đồng bộ hóa tiếp theo, đó là sử dụng `sync.WaitGroup` để chờ một tập các sự kiện: [(source)](../examples/ch1.6/1-hello-world-concurrent-ver/example-6/main.go)

```go
func main() {
    var wg sync.WaitGroup

    // Mở N thread
    for i := 0; i < 10; i++ {
        wg.Add(1)

        go func() {
            fmt.Println("Hello World")
            wg.Done()
        }()
    }

    // Đợi N thread hoàn thành
    wg.Wait()
}
```

Trong đó `wg.Add(1)` sử dụng để tăng số lượng sự kiện chờ, hàm này phải được đảm bảo thực thi trước khi bắt đầu thread nền (nếu nó được thực thi trong thread nền, nó không được đảm bảo có thể thực thi bình thường). Khi thread nền kết thúc việc in, lời gọi `wg.Done()` cho biết hoàn thành một sự kiện. Hàm `wg.Wait()` trong `main` để chờ tất cả các sự kiện hoàn thành.

## 1.6.2 Mô hình Producer Consumer

Ví dụ phổ biến nhất về lập trình đồng thời là mô hình Producer Consumer, giúp tăng tốc độ xử lý chung của chương trình bằng cách cân bằng sức mạnh làm việc của các thread "sản xuất" (produce) và "tiêu thụ" (consume). Nói một cách đơn giản, producer tạo ra một số dữ liệu và sau đó đưa nó vào hàng đợi kết quả, cùng lúc đó consumer cũng lấy dữ liệu từ hàng đợi này. Điều này làm cho sản xuất và tiêu thụ trở thành hai quá trình không đồng bộ. Khi không có dữ liệu trong hàng đợi kết quả, consumer sẽ chờ đợi ở trạng thái "đói", còn khi dữ liệu trong hàng đợi kết quả bị đầy, producer phải đối mặt với vấn đề CPU sẽ loại bỏ dữ liệu trong hàng đợi để nạp thêm.

Golang hiện thực cơ chế này rất đơn giản: [(source)](../examples/ch1.6/2-producer-consumer/example-1/main.go)

```go
// Producer: tạo ra một chuỗi số nguyên dựa trên bội số factor
func Producer(factor int, out chan<- int) {
    for i := 0; ; i++ {
        out <- i*factor
    }
}

// Consumer
func Consumer(in <-chan int) {
    for v := range in {
        fmt.Println(v)
    }
}
func main() {
    ch := make(chan int, 64) // hàng đợi kết quả

    go Producer(3, ch) // Tạo một chuỗi số với bội số 3
    go Producer(5, ch) // Tạo một chuỗi số với bội số 5
    go Consumer(ch)    // Tạo consumer

    // Thoát ra sau khi chạy trong một khoảng thời gian nhất định
    time.Sleep(5 * time.Second)
}
```

Chúng ta đã mở hai producer để tạo ra hai chuỗi bội số của 3 và 5. Sau đó mở một consumer và in ra kết quả. Ta  cho phép các producer và consumer làm việc trong một khoảng thời gian nhất định bằng hàm `time.Sleep`. Như đã đề cập trong phần trước, chế độ ngủ này không đảm bảo output ổn định.

Chúng ta có thể để hàm `main` giữ trạng thái block mà không thoát và chỉ  thoát khỏi chương trình khi người dùng gõ `Ctrl-C`: [(source)](../examples/ch1.6/2-producer-consumer/example-2/main.go)

```go
func main() {
    ch := make(chan int, 64) // hàng đợi kết quả

    go Producer(3, ch) // Tạo một chuỗi số với bội số 3
    go Producer(5, ch) // Tạo một chuỗi số với bội số 5
    go Consumer(ch)    // Tạo consumer

    // Ctrl+C để thoát
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    fmt.Printf("quit (%v)\n", <-sig)
}
```

Có 2 producer trong ví dụ trên và không có sự kiện đồng bộ nào giữa hai producer mà chúng đồng thời. Do đó, thứ tự của chuỗi output ở consumer là không xác định, tuy nhiên điều này không có vấn đề gì, và producer và consumer vẫn có thể làm việc cùng nhau.

## 1.6.3 Mô hình Publish Subscribe

Mô hình publish-and-subscribe thường được viết tắt là mô hình pub/sub. Trong mô hình này, producer trở thành publisher và consumer  trở thành subscriber, đồng thời producer:consumer là mối quan hệ M:N. Trong mô hình producer-consumer truyền thống, thông điệp được gửi đến hàng đợi và mô hình publish-subscription sẽ publish thông điệp đến một topic.

Để làm điều này, chúng tôi đã xây dựng một package  hỗ trợ mô hình pub/sub  tên là `pubsub`: [(source)](../examples/ch1.6/3-pubsub/example-1/main.go)

```go
// Package pubsub implements a simple multi-topic pub-sub library.
package pubsub

import (
    "sync"
    "time"
)

type (
    subscriber chan interface{}         // subscriber kiểu channel
    topicFunc  func(v interface{}) bool // topic là một filter
)

type Publisher struct {
    m           sync.RWMutex             // khóa đọc ghi
    buffer      int                      // kích thước  hàng đợi subscribe
    timeout     time.Duration            // hết thời gian publish
    subscribers map[subscriber]topicFunc // thông tin subscriber
}

// constructor với timeout và độ dài hàng đợi
func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
    return &Publisher{
        buffer:      buffer,
        timeout:     publishTimeout,
        subscribers: make(map[subscriber]topicFunc),
    }
}

// Thêm subscriber mới, đăng ký hết tất cả topic
func (p *Publisher) Subscribe() chan interface{} {
    return p.SubscribeTopic(nil)
}

// Thêm subscriber mới, đăng ký các filter được lọc theo topic
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
    ch := make(chan interface{}, p.buffer)
    p.m.Lock()
    p.subscribers[ch] = topic
    p.m.Unlock()
    return ch
}

// hủy đăng ký
func (p *Publisher) Evict(sub chan interface{}) {
    p.m.Lock()
    defer p.m.Unlock()

    delete(p.subscribers, sub)
    close(sub)
}

// đăng 1 topic
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

// Đóng 1 đối tượng publisher và đóng tất cả các subscriber
func (p *Publisher) Close() {
    p.m.Lock()
    defer p.m.Unlock()

    for sub := range p.subscribers {
        delete(p.subscribers, sub)
        close(sub)
    }
}

// Gửi 1 topic có thể duy trì trong thời gian chờ wg
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
    p := pubsub.NewPublisher(100*time.Millisecond, 10)
    defer p.Close()

    all := p.Subscribe()
    golang := p.SubscribeTopic(func(v interface{}) bool {
        if s, ok := v.(string); ok {
            return strings.Contains(s, "golang")
        }
        return false
    })

    p.Publish("hello,  world!")
    p.Publish("hello, golang!")

    go func() {
        for  msg := range all {
            fmt.Println("all:", msg)
        }
    } ()

    go func() {
        for  msg := range golang {
            fmt.Println("golang:", msg)
        }
    } ()

    // Thoát ra sau khi chạy 3 giây
    time.Sleep(3 * time.Second)
}
```

Trong mô hình pub/sub, mỗi thông điệp được gửi tới nhiều subscriber. Publisher thường không biết hoặc không quan tâm subscriber nào nhận được thông điệp. Subscriber và publisher có thể được thêm vào động ở thời điểm thực thi, một quan hệ không chặt cho phép hệ thống phức tạp có thể phát triển theo thời gian. Trong thực tế, những ứng dụng như dự báo thời tiết có thể áp dụng mô hình đồng thời này.

## 1.6.4 Kiểm soát Concurrency Numbers

Nhiều người dùng có xu hướng viết các chương trình đồng thời nhất có thể sau khi thích ứng với các tính năng đồng thời mạnh mẽ của ngôn ngữ Go, vì điều này dường như cung cấp một hiệu suất tối đa. 

Tuy nhiên trong thực tế chúng ta cần kiểm soát mức độ đồng thời ở mức thích hợp, bởi vì nó không chỉ có thể bỏ bớt các ứng dụng/task, dự trữ một lượng tài nguyên của CPU, ta cũng có thể giảm mức tiêu thụ năng lượng để giảm bớt áp lực cho pin.

Trong chương trình Godoc của ngôn ngữ Go,  package `vfs`  tương ứng với hệ thống tập tin ảo. Package phụ `gatefs` trong package 'vfs' với mục đích  là kiểm soát số lượng truy cập đồng thời tối đa vào hệ thống tập tin ảo. Ứng dụng của  package `gatefs`  rất đơn giản:

```go
import (
    "golang.org/x/tools/godoc/vfs"
    "golang.org/x/tools/godoc/vfs/gatefs"
)

func main() {
    fs := gatefs.New(vfs.OS("/path"), make(chan bool, 8))
    // ...
}
```

Trong trường hợp các cấu trúc hệ thống tập tin local  dựa trên một hệ thống tập tin ảo `vfs.OS("/path")`,  một cơ chế đồng thời `gatefs.New` sẽ kiểm soát hệ thống tập tin ảo dựa trên cấu trúc hệ thống tập tin ảo đang tồn tại. Nguyên tắc kiểm soát tương tranh đã được thảo luận ở phần trước, đó là để đạt được block đồng thời tối đa bằng cách gửi và nhận các rule với pipeline cache: [(source)](../examples/ch1.6/4-controlling-concurrency/example-1/main.go)

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

Ta bổ sung thêm phương thức `enter` và `leave` tương ứng để nhập vào và rời đi. Khi vượt quá số lượng giới hạn đồng thời, phương thức `enter` sẽ chặn cho đến khi số lượng đồng thời giảm xuống.

```go
type gate chan bool

func (g gate) enter() { g <- true }
func (g gate) leave() { <-g }
```

Hệ thống tập tin ảo mới `gatefs` được đóng gói là để thêm  lời gọi các phương thức `enter` và `leave`  cần kiểm soát đồng thời :

```go
type gatefs struct {
    fs vfs.FileSystem
    gate
}

func (fs gatefs) Lstat(p string) (os.FileInfo, error) {
    fs.enter()
    defer fs.leave()
    return fs.fs.Lstat(p)
}
```

Chúng ta không chỉ có thể kiểm soát số lượng đồng thời tối đa mà còn xác định tốc độ đồng thời của chương trình đang chạy bằng tỷ lệ sử dụng và dung lượng tối đa của channel được lưu trữ. Khi pipeline trống, nó có thể được coi như ở trạng thái không hoạt động. Khi pipeline đầy, tác vụ bận. Đây là giá trị tham chiếu cho hoạt động của một số tác vụ cấp thấp trong nền.

## 1.6.5 Kẻ thắng làm vua

Có nhiều động lực để lập trình đồng thời nhưng tiêu biểu là vì lập trình đồng thời có thể đơn giản hóa các vấn đề. Lập trình đồng thời cũng có thể cải thiện hiệu năng. Mở hai thread trên CPU đa lõi thường nhanh hơn mở một thread.  Trên thực tế về mặt cải thiện hiệu suất, chương trình không chỉ đơn giản là chạy nhanh, mà trong nhiều trường hợp chương trình có thể đáp ứng yêu cầu của người dùng một cách nhanh chóng là điều quan trọng nhất. Khi không có yêu cầu từ người dùng cần xử lý, nên xử lý một số tác vụ nền có độ ưu tiên thấp.

Giả sử chúng ta muốn nhanh chóng tìm kiếm các chủ đề liên quan đến "golang", có thể mở nhiều công cụ tìm kiếm như Bing, Google hoặc Yahoo. Khi tìm kiếm trả về kết quả trước, ta có thể đóng các trang tìm kiếm khác. Do ảnh hưởng của môi trường mạng và thuật toán của công cụ tìm kiếm mà một số công cụ tìm kiếm có thể trả về kết quả tìm kiếm nhanh hơn. Chúng ta có thể sử dụng một chiến lược tương tự để viết chương trình này:  [(source)](../examples/ch1.6/5-winner-is-king/example-1/main.go)

```go
func main() {
    ch := make(chan string, 32)

    go func() {
        ch <- searchByBing("golang")
    }()
    go func() {
        ch <- searchByGoogle("golang")
    }()
    go func() {
        ch <- searchByBaidu("golang")
    }()

    fmt.Println(<-ch)
}
```

Đầu tiên,   ta   tạo ra một pipeline   với cache đủ lớn để đảm bảo rằng không bị block không cần thiết do kích thước của cache. Sau đó, ta mở nhiều goroutine dưới nền và gửi yêu cầu tìm kiếm đến các công cụ tìm kiếm khác nhau. Khi bất kỳ công cụ tìm kiếm nào có kết quả đầu tiên, nó sẽ ngay lập tức gửi kết quả đến pipeline (vì pipeline có đủ bộ đệm, quá trình này sẽ không block). Nhưng cuối cùng,  chỉ cần lấy kết quả đầu tiên từ pipeline, đó là kết quả trả về đầu tiên.

Ta luôn có thể áp dụng nhiều cách giải quyết cho vấn đề theo các hướng khác nhau và cuối cùng cải thiện hiệu suất bằng cách chọn lấy cách giành chiến thắng trong cuộc đua thời gian.

## 1.6.6 Sàng số nguyên tố

Trong phần ***1.2***, chúng tôi đã trình bày việc triển khai phiên bản đồng thời của sàng số nguyên tố để chứng minh sự đồng thời của Newsqueak. Phiên bản đồng thời của Prime Screen là một ví dụ  cổ điển giúp chúng ta hiểu sâu hơn về các tính năng về tương tranh của Go. Nguyên tắc "sàng số nguyên tố" như sau:

![prime-sieve](../images/ch1-13-prime-sieve.png)
Hình 1-13 Sàng số nguyên tố

Chúng ta cần khởi tạo một chuỗi các số tự nhiên `2, 3, 4, ...` (không bao gồm 0, 1):

```go
// Trả về channel tạo ra chuỗi số: 2, 3, 4, ...
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

Tiếp theo xây dựng một sàng cho mỗi số nguyên tố: đề xuất một số là bội số của số nguyên tố trong chuỗi đầu vào và trả về một chuỗi mới, đó là một pipeline mới.

```go
// Bộ lọc: xóa các số có thể chia hết cho số nguyên tố
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

Bây giờ ta có thể sử dụng bộ lọc này trong hàm `main`: [(source)](../examples/ch1.6/6-prime-sieve/example-1/main.go)

```go
func main() {
    ch := GenerateNatural() // chuỗi số: 2, 3, 4, ...
    for i := 0; i < 100; i++ {
        prime := <-ch // số nguyên tố mới
        fmt.Printf("%v: %v\n", i+1, prime)
        ch = PrimeFilter(ch, prime) // Bộ lọc dựa trên số nguyên tố mới
    }
}
```

Đầu tiên chúng ta gọi `GenerateNatural()` để tạo ra chuỗi số tự nhiên nguyên thủy nhất bắt đầu bằng 2. Sau đó bắt đầu một chu kỳ 100 lần lặp. Ở đầu mỗi lần lặp lặp, số đầu tiên trong channel phải là số nguyên tố. Ta đọc và in ra số này  trước. Sau đó, dựa trên chuỗi còn lại trong channel và lọc các số nguyên tố tiếp theo với các số nguyên tố hiện được trích xuất dưới dạng sàng. Các channel tương ứng với các sàng số nguyên tố khác nhau được kết nối thành chuỗi.

## 1.6.7 Thoát khỏi  an toàn quá trình đồng thời

Đôi khi chúng ta cần  goroutine   ngăn chặn những gì nó đang làm, đặc biệt là khi nó đang làm việc sai hướng. Ngôn ngữ Go không cung cấp cách chấm dứt trực tiếp Goroutine, vì điều này sẽ khiến biến chung được chia sẻ giữa các   goroutine ở trạng thái không xác định. Nhưng điều gì sẽ xảy ra nếu chúng ta muốn loại hai hoặc nhiều Goroutines?

Goroutines khác nhau trong ngôn ngữ Go dựa vào các channel để giao tiếp và đồng bộ hóa. Để xử lý việc gửi hoặc nhận nhiều channel cùng một lúc, chúng ta cần sử dụng từ khóa `select` (từ khóa này hoạt động giống như một hàm `select` trong lập trình mạng ). Khi có nhiều nhánh khác nhau, `select` sẽ chọn một nhánh có sẵn ngẫu nhiên. Nếu không có nhánh có sẵn, nó sẽ chọn default, nếu không thì trạng thái block luôn được giữ.

Timeout dựa trên hiện thực của pipeline:

```go
select {
case v := <-in:
    fmt.Println(v)
case <-time.After(time.Second):
    return // hết giờ
}
```

Thông qua `select` nhánh `default` được gửi hoặc nhận nonblocking:

```go
select {
case v := <-in:
    fmt.Println(v)
default:
    // ...
}
```

Dùng `select` để block `main` không thoát:

```go
func main() {
    // do some thins
    select{}
}
```

Khi có nhiều channel có thể được thực thi, một channel sẽ được chọn ngẫu nhiên. Dựa trên tính năng này, ta có thể  thực hiện một chương trình tạo ra một chuỗi các số ngẫu nhiên: [(source)](../examples/ch1.6/7-concurrent-exit/example-1/main.go)

```go
func main() {
    ch := make(chan int)
    go func() {
        for {
            select {
            case ch <- 0:
            case ch <- 1:
            }
        }
    }()

    for v := range ch {
        fmt.Println(v)
    }
}
```

Chúng ta có thể dễ dàng thực hiện kiểm soát thoát Goroutine thông qua   nhánh `select` và  nhánh `default`: [(source)](../examples/ch1.6/7-concurrent-exit/example-2/main.go)

```go
func worker(cannel chan bool) {
    for {
        select {
        default:
            fmt.Println("hello")
            // thực hiện bình thường
        case <-cannel:
            // thoát
        }
    }
}

func main() {
    cannel := make(chan bool)
    go worker(cannel)

    time.Sleep(time.Second)
    cannel <- true
}
```

Tuy nhiên, các hoạt động gửi và nhận của channel là một đối một. Nếu ta muốn dừng nhiều Goroutines, ta  cần phải tạo ra cùng một số lượng channel. Điều này quá tốn kém. Trên thực tế, chúng ta có thể đạt được hiệu quả của việc broadcast bằng cách đóng một channel bằng `close`. Tất cả các hoạt động nhận được từ channel sẽ nhận được giá trị bằng 0 và cờ lỗi tùy chọn.  [(source)](../examples/ch1.6/7-concurrent-exit/example-3/main.go)

```go
func worker(cannel chan bool) {
    for {
        select {
        default:
            fmt.Println("hello")
            // hoạt động bình thường
        case <-cannel:
            // thoát
        }
    }
}

func main() {
    cancel := make(chan bool)

    for i := 0; i < 10; i++ {
        go worker(cancel)
    }

    time.Sleep(time.Second)
    close(cancel)
}
```

Chúng ta sử dụng channel `cancel` để phát chỉ thị `close` đến nhiều Goroutine. Tuy nhiên, chương trình này vẫn chưa đủ mạnh: khi mỗi Goroutine nhận được lệnh thoát để thoát, nó thường thực hiện một số công việc dọn dẹp, nhưng việc dọn dẹp của exit không được đảm bảo hoàn thành, vì thread `main` không có cơ chế chờ mỗi công việc Goroutine thoát khỏi công việc của chúng. Ta có thể kết hợp `sync.WaitGroup` để cải thiện điều này:  [(source)](../examples/ch1.6/7-concurrent-exit/example-4/main.go)

```go
func worker(wg *sync.WaitGroup, cannel chan bool) {
    defer wg.Done()

    for {
        select {
        default:
            fmt.Println("hello")
        case <-cannel:
            return
        }
    }
}

func main() {
    cancel := make(chan bool)

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go worker(&wg, cancel)
    }

    time.Sleep(time.Second)
    close(cancel)
    wg.Wait()
}
```

Bây giờ việc tạo, thực thi, đình chỉ và thoát khỏi quá trình đồng thời của mỗi thread worker nằm dưới sự kiểm soát bảo mật của hàm `main`.

## 1.6.8 Context package

Ở thời điểm phát hành Go1.7, thư viện tiêu chuẩn đã thêm một package context để đơn giản hóa hoạt động của dữ liệu, thời gian chờ và thoát giữa nhiều Goroutines. Chúng ta có thể sử dụng package context để hiện thực lại cơ chế kiểm soát thoát thread-safe hoặc kiểm soát timeout:  [(source)](../examples/ch1.6/8-context-package/example-1/main.go)

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
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go worker(ctx, &wg)
    }

    time.Sleep(time.Second)
    cancel()

    wg.Wait()
}
```

Khi cơ thể đồng thời hết thời gian hoặc `main` chủ động dừng  Goroutine worker, mỗi worker có thể hủy bỏ công việc một cách an toàn.

Golang tự động lấy lại   bộ nhớ, do đó bộ nhớ thường không bị rò rỉ (memory leak). Trong ví dụ trước về sàng số nguyên tố, một Goroutine mới  được đưa vào bên trong hàm `GenerateNatural` và Goroutine nền `PrimeFilter` có nguy cơ bị rò rỉ khi hàm `main` không còn sử dụng channel. Chúng ta có thể tránh vấn đề này với package context. Dưới đây là phần triển khai sàng số nguyên tố được cải thiện:  [(source)](../examples/ch1.6/7-concurrent-exit/example-2/main.go)

```go
// Trả về channel có chuỗi số: 2, 3, 4, ...
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

// Bộ lọc: xóa các số có thể chia hết cho số nguyên tố
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
    // Kiểm soát trạng thái Goroutine nền thông qua context
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

Khi hàm `main` kết thúc hoạt động, nó được thông báo bằng lệnh `cancel()` gọi đến Goroutine nền để thoát, do đó tránh khỏi việc rò rỉ Goroutine.

Đồng thời là một chủ đề rất lớn, và ở đây chúng tôi chỉ đưa ra một vài ví dụ về lập trình đồng thời rất cơ bản. Tài liệu chính thức cũng có rất nhiều cuộc thảo luận về lập trình đồng thời, có khá nhiều  cuốn sách  thảo luận cụ thể về lập trình đồng thời trong Golang. Độc giả có thể tham khảo các tài liệu liên quan theo nhu cầu của mình.
