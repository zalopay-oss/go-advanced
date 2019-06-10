# 4.3 RPC vui vẻ

Trong những trường hợp khác nhau lại có nhu cầu về RPC khác nhau, vì vậy cộng đồng mã nguồn mở đã tạo ra khá nhiều framework RPC. Trong phần này, chúng tôi sẽ sử dụng framework RPC tích hợp sẵn trong một số  tình huống đặc biệt.

## 4.3.1 Nguyễn tắc hiện thực của RPC Client

Cách dễ nhất để sử dụng thư viện Go là dùng phương thức `Client.Call` để thực hiện lời gọi đồng bộ blocking. Phần hiện thực của phương thức này như sau:

```go
func (client *Client) Call(
    serviceMethod string, args interface{},
    reply interface{},
) error {
    call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
    return call.Error
}
```

Đầu tiên `Client.Go` thực hiện một lời gọi không đồng bộ và trả về một cấu trúc `Call` đại diện cho nó. Sau đó `Call` sẽ chờ đợi pipe `Done` trả về kết quả lời gọi.

Chúng có cũng có thể gọi `Client.Go` tới service trước đó là `HelloService` theo cách bất đồng bộ bằng phương pháp sau:

```go
func doClientWork(client *rpc.Client) {
    helloCall := client.Go("HelloService.Hello", "hello", new(string), nil)

    // do some thing

    helloCall = <-helloCall.Done
    if err := helloCall.Error; err != nil {
        log.Fatal(err)
    }

    args := helloCall.Args.(string)
    reply := helloCall.Reply.(string)
    fmt.Println(args, reply)
}
```

Sau khi lệnh gọi không đồng bộ được thực hiện, các tác vụ khác sẽ được thực thi, do đó các tham số đầu vào và giá trị trả về của lời gọi không đồng bộ có thể  nhận được thông qua biến `Call` trả về.

Phương thức `Client.Go` thực thi một lời gọi bất đồng bộ được hiện thực như sau:

```go
func (client *Client) Go(
    serviceMethod string, args interface{},
    reply interface{},
    done chan *Call,
) *Call {
    call := new(Call)
    call.ServiceMethod = serviceMethod
    call.Args = args
    call.Reply = reply
    call.Done = make(chan *Call, 10) // buffered.

    client.send(call)
    return call
}
```

Phần đầu để khởi tạo một biến lời gọi đại diện cho cuộc lời gọi hiện thời, sau đó `client.send` gửi đi tham số đầy đủ của lời gọi đến RPC framework. Phương thức gọi `client.send` là thread-safe cho nên lệnh gọi có thể gửi từ nhiều Goroutine đồng thời tới cùng một đường link RPC.

Khi lời gọi hoàn thành hoặc có lỗi xuất hiện, phương thức thông báo `call.done` được gọi để hoàn thành:

```go
    select {
    case call.Done <- call:
        // ok
    default:
        // We don't want to block here. It is the caller's responsibility to make
        // sure the channel has enough buffer space. See comment in Go().
    }
}
```

Từ phần hiện thực của phương thức `Call.done`, có thể thấy rằng pipeline `call.Done` sẽ trả về lời gọi đã xử lý.

## 4.3.2 Hiện thực chức năng theo dõi dựa trên RPC

Trong nhiều hệ thống, interface cho việc theo dõi `Watch` được cung cấp. Khi hệ thống gặp những điều kiện nhất định, phương thức `Watch` trả về kết quả của việc giám sát. Chúng ta có thể thử hiện thực hàm `Watch` cơ bản thông qua RPC framework. Như đã đề cập ở trên, vì `client.send` là thread-safe, ta cũng có thể gọi phương thức RPC theo kiểu đồng thời blocking trong nhiều Goroutine khác nhau. Giám sát bằng cách gọi `Watch` trong những Goroutine riêng biệt.

Với mục đích dùng cho mô tả, chúng tôi dự định xây dựng cơ sở dữ liệu KV bộ nhớ đơn giản thông qua RPC. Đầu tiên xác định service như sau:

```go
type KVStoreService struct {
    m      map[string]string
    filter map[string]func(key string)
    mu     sync.Mutex
}

func NewKVStoreService() *KVStoreService {
    return &KVStoreService{
        m:      make(map[string]string),
        filter: make(map[string]func(key string)),
    }
}
```

`m` thuộc kiểu map được sử dụng để lưu trữ dữ liệu KV. `filter` tương ứng với một danh sách các hàm lọc được xác định tại mỗi cuộc gọi. `mu` thuộc kiểu mutex để cung cấp bảo vệ cho các thành phần khác khi được truy cập và sửa đổi từ nhiều Goroutine cùng lúc.

Sau đây là phương thức Get và Set:

```go
func (p *KVStoreService) Get(key string, value *string) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if v, ok := p.m[key]; ok {
        *value = v
        return nil
    }

    return fmt.Errorf("not found")
}

func (p *KVStoreService) Set(kv [2]string, reply *struct{}) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    key, value := kv[0], kv[1]

    if oldValue := p.m[key]; oldValue != value {
        for _, fn := range p.filter {
            fn(key)
        }
    }

    p.m[key] = value
    return nil
}
```

Trong phương thức Set, tham số đầu vào là một mảng của khóa và giá trị, struct rỗng ẩn danh được sử dụng để bỏ qua các tham số đầu ra. Mỗi hàm lọc được gọi khi giá trị tương ứng với khóa được sửa đổi.

Danh sách các bộ lọc được cung cấp trong phương thức `Watch`:

```go
func (p *KVStoreService) Watch(timeoutSecond int, keyChanged *string) error {
    id := fmt.Sprintf("watch-%s-%03d", time.Now(), rand.Int())
    ch := make(chan string, 10) // buffered

    p.mu.Lock()
    p.filter[id] = func(key string) { ch <- key }
    p.mu.Unlock()

    select {
    case <-time.After(time.Duration(timeoutSecond) * time.Second):
        return fmt.Errorf("timeout")
    case key := <-ch:
        *keyChanged = key
        return nil
    }

    return nil
}
```

Tham số đầu vào của phương thức `Watch` là số giây timeout. Khóa được trả về dưới dạng giá trị trả về khi có khóa thay đổi. Nếu không có khóa nào được sửa đổi sau khi timeout, lỗi thời gian chờ được trả về. Trong quá trình triển khai Đồng hồ, mỗi cuộc gọi Đồng hồ được biểu thị bằng một id duy nhất và sau đó chức năng bộ lọc tương ứng được đăng ký vào p.filterdanh sách theo id .

