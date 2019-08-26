# 5.2 Lock phân tán (Distributed lock)

<div align="center">
	<img src="../images/ch6-dis-lock.png" width="300">
	<br/>
	<span align="center">
		<i>Distributed lock</i>
	</span>
</div>

Trước khi đi vào nội dung chính là `Distributed lock` chúng ta cùng xem xét bài toán sau. Khi một chương trình chạy đồng thời hoặc song song sửa đổi biến toàn cục, hành vi sửa đổi này cần phải được `lock` để tránh trường hợp [race conditions](https://en.wikipedia.org/wiki/Race_condition). Tại sao bạn cần phải lock ? Hãy xem điều gì xảy ra khi trong bài toán đếm số một cách đồng thời mà không lock dưới đây.

***main.go***
```go
package main

import (
  "sync"
  "fmt"
)

// biến đếm toàn cục
var counter int

func main() {
  var wg sync.WaitGroup
  for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
    defer wg.Done()
      // tăng biến counter lên một đơn vị
      counter++
    }()
  }

  wg.Wait()
  fmt.Println("Counter: ", counter)
}
```

Khi ta chạy nhiều lần, các kết quả sẽ khác nhau:

```sh
$ go run local_lock.go
Counter: 945
$ go run local_lock.go
Counter 937
$ go run local_lock.go
Counter: 959
```

## 5.2.1 Lock quá trình đang thực hiện

Để có kết quả chính xác, lock phần code thực thi của bộ đếm như ví dụ dưới đây.

```go
// ... bỏ qua phần trước
var wg sync.WaitGroup
var l sync.Mutex
for i := 0; i < 1000; i++ {
  wg.Add(1)
  go func() {
    defer wg.Done()
    // lấy lọck
    l.Lock()
    counter++
    // trả lock
    l.Unlock()
  }()
}

wg.Wait()
fmt.Println("Counter: ", counter)
```

Các lần chạy đều cho ra cùng một kết quả:

```shell
$ go run local_lock.go
1000
```

## 5.2.2 Sử dụng Trylock

Trong một số tình huống, chúng ta chỉ muốn một tiến trình thực thi một nhiệm vụ. Ở ví dụ đếm số ở trên, tất cả goroutines đều thực hiện thành công. Giả sử có goroutine thất bại trong khi thực hiện, chúng ta cần phải bỏ qua tiến trình của nó. Đây là lúc cần `trylock`.

Trylock, như tên của nó, cố gắng lock và nếu lock thành công thì thực hiện các công việc tiếp theo. Nếu lock bị lỗi, nó sẽ không bị chặn lại mà sẽ trả về kết quả lock. Trong lập trình Go, chúng ta có thể mô phỏng một trylock với channel có kích thước buffer là 1.

***main.go***
```go
package main

import (
  "sync"
)

// Lock try lock
type Lock struct {
  c chan struct{}
}

// tạo một lock
func NewLock() Lock {
  var l Lock
  l.c = make(chan struct{}, 1)
  l.c <- struct{}{}
  return l
}

// Lock try lock, trả về kết qủa lock là true/false
func (l Lock) Lock() bool {
  lockResult := false
  select {
  case <-l.c:
    lockResult = true
  default:
  }
  return lockResult
}

// Unlock , giải phóng lock
func (l Lock) Unlock() {
  l.c <- struct{}{}
}

// biến đếm toàn cục
var counter int

func main() {
  var l = NewLock()
  var wg sync.WaitGroup
  for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
      defer wg.Done()
      if !l.Lock() {
        // log error
        fmt.Println("lock failed")
        return
      }
      counter++
      fmt.Println("current counter", counter)
      l.Unlock()
    }()
  }
  wg.Wait()
}

// output
// lock failed
// lock failed
// lock failed
// lock failed
// lock failed
// lock failed
// current counter 1
// lock failed
// lock failed
// lock failed
```

Bởi vì logic của chúng ta giới hạn bởi mỗi goroutine chỉ thực hiện logic sau khi nó Lock thành công. Còn đối với Unlock, nó đảm bảo rằng channel của Lock ở đoạn code trên phải trống, nên nó sẽ không bị chặn hoặc thất bại giữa chừng. Đoạn code trên sử dụng channel có kích thước 1 để mô phỏng một tryLock. Về lý thuyết, bạn có thể sử dụng [CAS](https://en.wikipedia.org/wiki/Compare-and-swap) trong thư viện chuẩn để đạt được chức năng tương tự với chi phí thấp hơn. Bạn có thể thử dùng nó.

Trong một hệ thống đơn, trylock không phải là một lựa chọn tốt. Bởi vì khi có một lượng lớn lock goroutine có thể gây lãng phí tài nguyên trong CPU một cách vô nghĩa. Có một danh từ thích hợp được sử dụng để mô tả kịch bản khóa này: `livelock`.

`livelock` nghĩa là chương trình trông có vẻ đang thực thi bình thường, nhưng trên thực tế, CPU bị lãng phí khi phải lo lấy lock thay vì thực thi tác vụ, do đó việc thực thi của chương trình không hiệu quả. Vấn đề của livelock sẽ gây ra rất nhiều hậu quả xấu. Do đó, trong ngữ cảnh máy đơn, không nên sử dụng loại lock này.

## 5.2.3 Redis dựa trên setnx

Trong ngữ cảnh phân tán, chúng ta cũng cần một loại logic "ưu tiên". Làm sao để có được nó? Chúng ta có thể sử dụng lệnh [setnx](https://github.com/go-redis/redis) do Redis cung cấp:

***main.go***
```go
package main

import (
  "fmt"
  "sync"
  "time"

  "github.com/go-redis/redis"
)

func incr() {
  client := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    // không cần khởi tạo password
    Password: "", 
    // không dùng database
    DB:       0,  
  })

  var lockKey = "counter_lock"
  var counterKey = "counter"

  // lấy lock
  resp := client.SetNX(lockKey, 1, time.Second*5)
  lockSuccess, err := resp.Result()

  if err != nil || !lockSuccess {
    fmt.Println(err, "lock result: ", lockSuccess)
    return
  }

  // tăng cntValue lên một đơn vị
  getResp := client.Get(counterKey)
  cntValue, err := getResp.Int64()
  if err == nil || err == redis.Nil {
    cntValue++
    resp := client.Set(counterKey, cntValue, 0)
    _, err := resp.Result()
    if err != nil {
      // log err
      fmt.Println("set value error!")
    }
  }
  fmt.Println("current counter is ", cntValue)

  delResp := client.Del(lockKey)
  // giải phóng lock
  unlockSuccess, err := delResp.Result()
  if err == nil && unlockSuccess > 0 {
    fmt.Println("unlock success!")
  } else {
    fmt.Println("unlock failed", err)
  }
}

func main() {
  var wg sync.WaitGroup
  for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
      defer wg.Done()
      incr()
    }()
  }
  wg.Wait()
}
```

Nhìn vào kết quả khi chạy:

```shell
$ go run redis_setnx.go
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

Thông qua mã nguồn và kết quả chạy thực tế, chúng ta có thể thấy rằng khi gọi `setnx` thực sự rất giống với một trylock. Nếu khóa bị lỗi, logic của tác vụ liên quan sẽ không được thực thi.

`setnx` tuyệt vời cho các ngữ cảnh cần tính đồng thời cao và được sử dụng để giành một số tài nguyên "duy nhất". Ví dụ: trong hệ thống kiểm tra giao dịch, người bán tạo một đơn đặt hàng và nhiều người mua sẽ giành lấy nó. Trong ví dụ này, chúng ta không có cách nào dựa vào thời gian cụ thể để phán đoán thứ tự. Bởi vì, dù đó là thời gian của thiết bị người dùng hay thời gian của mỗi server trong hệ thống phân tán, không có cách nào để đảm bảo thời gian chính xác sau khi tổng hợp lại. Ngay cả khi chúng ta đã đặt các service chung một server thì vẫn có sự khác biệt nhỏ.

Do đó, chúng ta cần dựa vào thứ tự của các yêu cầu này để node Redis thực hiện thao tác khóa chính xác. Nếu môi trường mạng của người dùng tương đối kém, thì họ chỉ cần tạo thêm yêu cầu. Các bạn có thể xem chi tiết phần `Distributed lock` bằng Redis [ở đây](https://redis.io/topics/distlock).

## 5.2.4 Sử dụng ZooKeeper

Chúng ta cùng xem qua một ví dụ về dùng lock trong [ZooKeeper](https://github.com/samuel/go-zookeeper).

***main.go***
```go
package main

import (
  "fmt"
  "time"

  "github.com/samuel/go-zookeeper/zk"
)

func main() {
  c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
  if err != nil {
    panic(err)
  }

  // tạo lock
  l := zk.NewLock(c, "/lock", zk.WorldACL(zk.PermAll))
  // sử dụng lock
  err = l.Lock()
  if err != nil {
    panic(err)
  }
  fmt.Println("lock succ, do your business logic")

  time.Sleep(time.Second * 10)

  // ...xử lý logic của bạn

  // giải phóng lock
  l.Unlock()
  fmt.Println("unlock succ, finish business logic")
}
```

Lock dựa trên ZooKeeper khác với lock dựa trên Redis ở chỗ nó sẽ chặn cho đến khi lấy lock thành công, tương tự như `mutex.Lock`.

Nguyên tắc này cũng dựa trên node (một server trong hệ thống phân tán) thứ tự tạm thời và quan sát API. Ví dụ, chúng ta sử dụng node `/lock`. Các Lock sẽ chèn giá trị của chính nó vào danh sách node bên dưới node này. Khi các node con ở dưới node này thay đổi, nó sẽ thông báo cho tất cả các chương trình quan sát giá trị của node. Lúc này, chương trình sẽ kiểm tra xem ID của node con gần node hiện tại nhất có giống với giá trị của chính nó không. Nếu chúng giống nhau thì việc lock diễn ra thành công.

Loại lock phân tán này phù hợp hơn cho các ngữ cảnh định thời tác vụ phân tán, nhưng nó không phù hợp trong các ngữ cảnh thường xuyên cần lock trong thời gian lâu. Theo bài báo [Chubby](https://static.googleusercontent.com/media/research.google.com/en//archive/chubby-osdi06.pdf) của Google, các lock dựa trên các giao thức nhất quán cao áp dụng cho loại khóa "coarse-grained". Loại "coarse-grained" có nghĩa là khóa mất nhiều thời gian. Chúng ta nên xem xét liệu lock này có phù hợp để sử dụng với mục đích của chúng ta hay không.

## 5.2.5 Sử dụng etcd

[Etcd](https://github.com/etcd-io/etcd) là một thư viện thường được dùng trong các hệ thống phân tán, nó có chức năng gần giống với ZooKeeper và đã trở nên "hot" hơn trong hai năm qua. Dựa trên ZooKeeper, chúng tôi đã triển khai distributed lock. Với etcd, chúng ta cũng có thể thực hiện các chức năng tương tự:

***main.go***
```go
package main

import (
  "log"

  "github.com/zieckey/etcdsync"
)

func main() {
  // khởi tạo lock
  m, err := etcdsync.New("/lock", 10, []string{"http://127.0.0.1:2379"})
  if m == nil || err != nil {
    log.Printf("etcdsync.New failed")
    return
  }
  // lock 
  err = m.Lock()
  if err != nil {
    log.Printf("etcdsync.Lock failed")
    return
  }

  log.Printf("etcdsync.Lock OK")
  log.Printf("Get the lock. Do something here.")

  // giải phóng lock
  err = m.Unlock()
  if err != nil {
    log.Printf("etcdsync.Unlock failed")
  } else {
    log.Printf("etcdsync.Unlock OK")
  }
}
```

Không có node thứ tự như ZooKeeper trong etcd. Vì vậy, việc thực hiện lock của nó khác với ZooKeeper. Quá trình lock cho etcdsync của đoạn code mẫu ở trên cụ thể như sau:

1. Kiểm tra xem có giá trị nào trong đường dẫn `/lock` không. Nếu có một giá trị, khóa đã bị người khác lấy.
2. Nếu không có giá trị, nó sẽ ghi giá trị của chính nó vào. Khi ghi giá trị thành công thì lock đã thành công. Giả sử có một node đang ghi thì node khác đến ghi, điều này khiến khóa bị lỗi. Tiếp bước 3.
3. Kiểm tra sự kiện trong `/lock`, cái mà đang bị kẹt tại lúc này.
4. Khi một sự kiện xảy ra trong đường dẫn `/lock`, process hiện tại được đánh thức. Kiểm tra xem sự kiện xảy ra có phải là sự kiện xóa không (chứng tỏ rằng lock đang được mở) hoặc sự kiện đã hết hạn (cho biết khóa đã hết hạn). Nếu vậy, quay lại bước 1 và thực hiện tuần tự các bước.

Điều đáng nói là trong [etcdv3 API](https://github.com/etcd-io/etcd/releases/tag/v3.3.13) đã chính thức cung cấp API lock mà có thể sử dụng trực tiếp. Các bạn có thể tham khảo tài liệu etcdv3 để nghiên cứu thêm.

## 5.2.7 Làm sao để chọn đúng loại lock

Nếu bạn phát triển một dịch vụ phân tán, nhưng quy mô dịch vụ kinh doanh không lớn, thì sử dụng lock nào cũng như nhau. Nếu bạn có một cụm ZooKeeper, etcd hoặc Redis có sẵn trong công ty, hãy sử dụng cái có sẵn để đáp ứng nhu cầu kinh doanh của bạn mà không cần các công nghệ mới.

Nếu doanh nghiệp phát triển đến một mức độ nhất định, thì chúng ta cần xem xét ở nhiều khía cạnh. Đầu tiên là hãy xem xét liệu lock của bạn có cho phép mất dữ liệu trong bất kỳ điều kiện nào không. Nếu không, thì đừng sử dụng khóa `setnx` của Redis.

Nếu độ tin cậy của dữ liệu lock là rất cao, thì chỉ có khóa etcd hoặc ZooKeeper đảm bảo độ tin cậy của dữ liệu. Nhưng mặt trái của sự đáng tin cậy là throughput thấp hơn và latency cao hơn. Cần phải có những bài kiểm tra kỹ lưỡng theo từng cấp độ kinh doanh để đảm bảo rằng các lock phân tán bằng cụm etcd hoặc ZooKeeper có thể chịu được áp lực của các yêu cầu kinh doanh thực tế. 

`Lưu ý:` sẽ không có cách nào để cải thiện hiệu suất của cụm etcd và Zookeeper bằng cách thêm các node. Để mở rộng quy mô theo chiều ngang, bạn chỉ có thể tăng số lượng cụm để hỗ trợ nhiều yêu cầu hơn. Điều này sẽ tăng thêm các chi phí vận hành, bảo trì và giám sát. Nhiều cụm có thể cần phải thêm proxy. Nếu không có proxy, dịch vụ cần được phân phối theo một ID nhất định. Nếu dịch vụ đã được mở rộng, bạn cũng nên xem xét việc di chuyển dữ liệu động. Đây không phải là điều dễ dàng.

[Tiếp theo](ch5-03-delay-job.md)