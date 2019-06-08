# 6.2 Distributed lock

Khi một chương trình đồng thời hoặc song song sửa đổi biến toàn cục, hành vi sửa đổi cần phải được lock để tạo một vùng tranh chấp. Tại sao bạn cần phải lock? Hãy xem điều gì xảy ra khi trong bài toán đếm số một cách đồng thời mà không lock:

```go
Package main

Import (
  "sync"
)

// biến toàn cục
Var counter int

Func main() {
  Var wg sync.WaitGroup
  For i := 0; i < 1000; i++ {
    wg.Add(1)
    Go func() {
    Defer wg.Done()
      Counter++
    }()
  }

  wg.Wait()
  Println(counter)
}
```

Khi ta chạy nhiều lần, các kết quả sẽ khác nhau:

```shell
❯❯❯ go run local_lock.go
945
❯❯❯ go run local_lock.go
937
❯❯❯ go run local_lock.go
959
```

## 6.2.1 Lock quá trình đang thực hiện

Để có kết quả chính xác, lock phần code thực thi của bộ đếm:

```go
// ... bỏ qua phần trước
Var wg sync.WaitGroup
Var l sync.Mutex
For i := 0; i < 1000; i++ {
  wg.Add(1)
  Go func() {
    Defer wg.Done()
    l.Lock()
    Counter++
    l.Unlock()
  }()
}

wg.Wait()
Println(counter)
// ... after omitting the part
```

Kết quả tính toán sẽ ổn định:

```shell
❯❯❯ go run local_lock.go
1000
```

## 6.2.2 Trylock

Trong một số tình huống, chúng ta chỉ muốn một tiến trình thực thi một nhiệm vụ. Ở ví dụ đếm số ở trên, tất cả goroutines đều thực hiện thành công. Giả sử có goroutine thất bại trong khi thực hiện, chúng ta cần phải bỏ qua tiến trình của nó. Đây là lúc cần `trylock`.

Trylock, như tên của nó, cố gắng lock và nếu lock thành công thì thực hiện các công việc tiếp theo. Nếu lock bị lỗi, nó sẽ không bị chặn lại mà sẽ trả về kết quả lock. Trong lập trình Go, chúng ta có thể mô phỏng một trylock với kênh có kích thước 1:

```go
Package main

Import (
 "sync"
)

// Lock try lock
Type lock struct {
 c chan struct{}
}

// NewLock generate a try lock
Func NewLock() Lock {
 Var l Lock
 Lc = make(chan struct{}, 1)
 Lc <- struct{}{}
 Return l
}

// Lock try lock, return lock result
Func (l Lock) Lock() bool {
 lockResult := false
 Select {
 Case <-lc:
  lockResult = true
 Default:
 }
 Return lockResult
}

// Unlock , Unlock the try lock
Func (l Lock) Unlock() {
 Lc <- struct{}{}
}

Var counter int

Func main() {
 Var l = NewLock()
 Var wg sync.WaitGroup
 For i := 0; i < 10; i++ {
  wg.Add(1)
  Go func() {
   Defer wg.Done()
   If !l.Lock() {
    // log error
    Println("lock failed")
    Return
   }
   Counter++
   Println("current counter", counter)
   l.Unlock()
  }()
 }
 wg.Wait()
}
```

Bởi vì logic của chúng ta giới hạn mỗi con goroutine chỉ thực hiện logic sau khi nó `Lock` thành công. Còn đối với `Unlock`, nó đảm bảo rằng kênh của Lock ở đoạn code trên phải trống, nên nó sẽ không bị chặn hoặc thất bại giữa chừng. Đoạn code trên sử dụng kênh có kích thước 1 để mô phỏng một tryLock. Về lý thuyết, bạn có thể sử dụng CAS trong thư viện chuẩn để đạt được chức năng tương tự với chi phí thấp hơn. Bạn có thể thử dùng nó.

Trong một hệ thống đơn, trylock không phải là một lựa chọn tốt. Bởi vì khi có một lượng lớn khóa goroutine có thể gây lãng phí tài nguyên trong CPU một cách vô nghĩa. Có một danh từ thích hợp được sử dụng để mô tả kịch bản khóa này: `livelock`.

`livelock` nghĩa là chương trình trông có vẻ đang thực thi bình thường, nhưng trên thực tế, CPU bị lãng phí khi phải lo lấy lock thay vì thực thi tác vụ, do đó việc thực thi của chương trình không hiệu quả. Vấn đề của livelock sẽ gây ra rất nhiều hậu quả xấu. Do đó, trong ngữ cảnh máy đơn, không nên sử dụng loại khóa này.

## 6.2.3 Redis dựa trên setnx

Trong ngữ cảnh phân tán, chúng ta cũng cần một loại logic "ưu tiên". Làm sao để có được nó? Chúng ta có thể sử dụng lệnh `setnx` do Redis cung cấp:

```go
Package main

Import (
 "fmt"
 "sync"
 "time"

 "github.com/go-redis/redis"
)

Func incr() {
 Client := redis.NewClient(&redis.Options{
  Addr: "localhost:6379",
  Password: "", // no password set
  DB: 0, // use default DB
 })

 Var lockKey = "counter_lock"
 Var counterKey = "counter"

 // lock
 Resp := client.SetNX(lockKey, 1, time.Second*5)
 lockSuccess, err := resp.Result()

 If err != nil || !lockSuccess {
  fmt.Println(err, "lock result: ", lockSuccess)
  Return
 }

 // counter ++
 getResp := client.Get(counterKey)
 cntValue, err := getResp.Int64()
 If err == nil || err == redis.Nil {
  cntValue++
  Resp := client.Set(counterKey, cntValue, 0)
  _, err := resp.Result()
  If err != nil {
   // log err
   Println("set value error!")
  }
 }
 Println("current counter is ", cntValue)

 delResp := client.Del(lockKey)
 unlockSuccess, err := delResp.Result()
 If err == nil && unlockSuccess > 0 {
  Println("unlock success!")
 } else {
  Println("unlock failed", err)
 }
}

Func main() {
 Var wg sync.WaitGroup
 For i := 0; i < 10; i++ {
  wg.Add(1)
  Go func() {
   Defer wg.Done()
   Incr()
  }()
 }
 wg.Wait()
}
```

Nhìn vào kết quả khi chạy:

```shell
❯❯❯ go run redis_setnx.go
<nil> lock result: false
<nil> lock result: false
<nil> lock result: false
<nil> lock result: false
<nil> lock result: false
<nil> lock result: false
<nil> lock result: false
<nil> lock result: false
<nil> lock result: false
Current counter is 2028
Unlock success!
```

Thông qua code và kết quả chạy thực tế, chúng ta có thể thấy rằng khi gọi `setnx` từ xa thì nó thực sự rất giống với một trylock. Nếu khóa bị lỗi, logic của tác vụ liên quan sẽ không được thực thi.

`setnx` tuyệt vời cho các ngữ cảnh cần tính đồng thời cao và được sử dụng để giành một số tài nguyên "duy nhất". Ví dụ: trong hệ thống kiểm tra giao dịch, người bán tạo một đơn đặt hàng và nhiều người mua sẽ giành lấy nó. Trong ví dụ này, chúng ta không có cách nào dựa vào thời gian cụ thể để phán đoán thứ tự, bởi vì, dù đó là thời gian của thiết bị người dùng hay thời gian của mỗi máy trong hệ thống phân tán, không có cách nào để đảm bảo thời gian chính xác sau khi tổng hợp lại. Ngay cả khi chúng ta đã đặ các máy trong cùng một phòng, thời gian hệ thống của các máy khác nhau vẫn có sự khác biệt nhỏ.

Do đó, chúng ta cần dựa vào thứ tự của các yêu cầu này để node Redis để thực hiện thao tác khóa chính xác. Nếu môi trường mạng của người dùng tương đối kém, thì họ chỉ cần tạo thêm yêu cầu.

## 6.2.4 Dựa trên ZooKeeper

```go
Package main

Import (
 "time"

 "github.com/samuel/go-zookeeper/zk"
)

Func main() {
 c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
 If err != nil {
  Panic(err)
 }
 l := zk.NewLock(c, "/lock", zk.WorldACL(zk.PermAll))
 Err = l.Lock()
 If err != nil {
  Panic(err)
 }
 Println("lock succ, do your business logic")

 time.Sleep(time.Second * 10)

 // do some thing
 l.Unlock()
 Println("unlock succ, finish business logic")
}
```

Lock dựa trên ZooKeeper khác với lock dựa trên Redis ở chỗ nó sẽ chặn cho đến khi lấy lock thành công, tương tự như `mutex.Lock` trong ngữ cảnh máy đơn.

Nguyên tắc này cũng dựa trên node Thứ tự tạm thời và quan sát API. Ví dụ, chúng ta sử dụng node `/lock`. Các Lock sẽ chèn giá trị của chính nó vào danh sách node bên dưới node này. Khi các node con ở dưới node này thay đổi, nó sẽ thông báo cho tất cả các chương trình quan sát giá trị của node. Lúc này, chương trình sẽ kiểm tra xem id của node con gần node hiện tại nhất có giống với giá trị của chính nó không. Nếu chúng giống nhau, lock thành công.

This kind of distributed blocking lock is more suitable for distributed task scheduling scenarios, but it is not suitable for stealing scenarios with high frequency locking time. According to Google's Chubby paper, locks based on strong consistent protocols apply to the "coarse-grained" locking operation. The coarse grain size here means that the lock takes a long time. We should also consider whether it is appropriate to use it in our own business scenarios.

Loại khóa chặn phân tán này phù hợp hơn cho các ngữ cảnh định thời tác vụ phân tán, nhưng nó không phù hợp trong các ngữ cảnh thường xuyên cần lock trong thời gian lâu. Theo bài báo Chubby của Google, các lock dựa trên các giao thức nhất quán cao áp dụng cho loại khóa "coarse-grained". Loại "coarse-grained" có nghĩa là khóa mất nhiều thời gian. Chúng ta nên xem xét liệu kháo này có phù hợp để sử dụng với mục đích của chúng ta hay không.

## 6.2.5 Dựa trên etcd

Etcd là một thành phần của một hệ thống phân tán có chức năng giống với ZooKeeper và đã trở nên "hot" hơn trong hai năm qua. Dựa trên ZooKeeper, chúng tôi đã triển khai khóa chặn phân tán. Với etcd, chúng ta cũng có thể thực hiện các chức năng tương tự:

```go
Package main

Import (
 "log"

 "github.com/zieckey/etcdsync"
)

Func main() {
 m, err := etcdsync.New("/lock", 10, []string{"http://127.0.0.1:2379"})
 If m == nil || err != nil {
  log.Printf("etcdsync.New failed")
  Return
 }
 Err = m.Lock()
 If err != nil {
  log.Printf("etcdsync.Lock failed")
  Return
 }

 log.Printf("etcdsync.Lock OK")
 log.Printf("Get the lock. Do something here.")

 Err = m.Unlock()
 If err != nil {
  log.Printf("etcdsync.Unlock failed")
 } else {
  log.Printf("etcdsync.Unlock OK")
 }
}
```

Không có node Thứ tự như ZooKeeper trong etcd. Vì vậy, việc thực hiện lock của nó khác với ZooKeeper. Quá trình lock cho etcdsync của đoạn code mẫu ở trên cụ thể như sau:

1. Kiểm tra xem có giá trị nào trong đường dẫn `/lock` không. Nếu có một giá trị, khóa đã bị người khác lấy.
2. Nếu không có giá trị, nó sẽ ghi giá trị của chính nó vào. Khi ghi giá trị thành công thì lock đã thành công. Giả sử có một node đang ghi thì node khác đến ghi, điều này khiến khóa bị lỗi. Tiếp bước 3.
3. Kiểm tra sự kiện trong `/lock`, cái mà đang bị kẹt tại lúc này.
4. Khi một sự kiện xảy ra trong đường dẫn `/lock`, process hiện tại được đánh thức. Kiểm tra xem sự kiện xảy ra có phải là sự kiện xóa không (chứng tỏ rằng lock đang được mở) hoặc sự kiện đã hết hạn (cho biết khóa đã hết hạn). Nếu vậy, quay lại bước 1 và thực hiện tuần tự các bước.

Điều đáng nói là trong etcdv3 API đã chính thức cung cấp API lock mà có thể sử dụng trực tiếp. Các bạn có thể tham khảo tài liệu etcdv3 để nghiên cứu thêm.

## 6.2.7 Làm sao để chọn đúng loại lock

Khi vấn đề kinh doanh quan trọng việc hoạt động như trên một máy đơn, thì nên sử dụng lock.

Nếu bạn phát triển một dịch vụ phân tán, nhưng quy mô kinh doanh không lớn, qps thì nhỏ, thì sử dụng lock nào cũng như nhau. Nếu bạn có một cụm ZooKeeper, etcd hoặc Redis có sẵn trong công ty, hãy sử dụng cái có sẵn để đáp ứng nhu cầu kinh doanh của bạn mà không cần các công nghệ mới.

Nếu doanh nghiệp phát triển đến một mức độ nhất định, thì chúng ta cần xem xét ở nhiều khía cạnh. Đầu tiên là hãy xem xét liệu lock của bạn có cho phép mất dữ liệu trong bất kỳ điều kiện nào không. Nếu không, thì đừng sử dụng khóa `setnx` của Redis.

Nếu độ tin cậy của dữ liệu lock là rất cao, thì chỉ có khóa etcd hoặc ZooKeeper đảm bảo độ tin cậy của dữ liệu thông qua giao thức kết hợp. Nhưng mặt trái của sự đáng tin cậy là throughput thấp hơn và latency cao hơn. Cần phải có những bài kiểm tra kỹ lưỡng theo từng cấp độ kinh doanh để đảm bảo rằng các lock phân tán bằng cụm etcd hoặc ZooKeeper có thể chịu được áp lực của các yêu cầu kinh doanh thực tế. Cần lưu ý sẽ không có cách nào để cải thiện hiệu suất của cụm etcd và Zookeeper bằng cách thêm các node. Để mở rộng quy mô theo chiều ngang, bạn chỉ có thể tăng số lượng cụm để hỗ trợ nhiều yêu cầu hơn. Điều này sẽ tăng thêm các chi phí vận hành, bảo trì và giám sát. Nhiều cụm có thể cần phải thêm proxy. Nếu không có proxy, dịch vụ cần được phân phối theo một ID nhất định. Nếu dịch vụ đã được mở rộng, bạn cũng nên xem xét việc di chuyển dữ liệu động. Đây không phải là điều dễ dàng.

Khi chọn một kế hoạch cụ thể, bạn cần suy nghĩ nhiều hơn và đưa ra các dự đoán rủi ro sớm.
