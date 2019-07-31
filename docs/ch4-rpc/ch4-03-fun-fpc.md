# 4.3 RPC trong Golang

Trong những trường hợp khác nhau lại có nhu cầu về RPC khác nhau, vì vậy cộng đồng opensource đã tạo ra khá nhiều framework RPC. Trong phần này, chúng tôi sẽ sử dụng framework RPC tích hợp sẵn trong một số  tình huống đặc biệt.

## 4.3.1 Nguyên tắc hiện thực của RPC Client

Cách dễ nhất để sử dụng thư viện Go là dùng phương thức `Client.Call` để thực hiện lời gọi synchronous blocking. Phần hiện thực của phương thức này như sau:

```go
func (client *Client) Call(
    serviceMethod string, args interface{},
    reply interface{},
) error {
    call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
    return call.Error
}
```

Đầu tiên `client.Go` thực hiện một lời gọi bất đồng bộ và trả về một cấu trúc `Call`. Sau đó `Call` sẽ chờ  pipe `Done` trả về kết quả lời gọi.

Chúng  ta cũng có thể gọi `client.Go` tới service trước đó là `HelloService` theo kiểu bất đồng bộ bằng phương pháp sau:

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

Sau khi lệnh gọi bất đồng bộ được thực hiện, các tác vụ khác sẽ được thực thi, sau đó các tham số đầu vào và giá trị trả về của lời gọi bất đồng bộ có thể  nhận được thông qua biến `Call` trả về.

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

[net/rpc/client.go](https://golang.org/src/net/rpc/client.go)

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

Trong nhiều hệ thống, interface cho việc theo dõi (`Watch`) được cung cấp. Khi hệ thống gặp những điều kiện nhất định, phương thức `Watch` trả về kết quả của việc giám sát. Chúng ta có thể thử hiện thực hàm `Watch` cơ bản thông qua RPC framework. Như đã đề cập ở trên, vì `client.send` là thread-safe, ta cũng có thể gọi phương thức RPC theo kiểu đồng bộ blocking trong nhiều Goroutine khác nhau. Giám sát bằng cách gọi `Watch` trong những Goroutine riêng biệt.

Để minh họa, ta sẽ đi xây dựng cơ sở dữ liệu KV bộ nhớ đơn giản thông qua RPC. Đầu tiên xác định service như sau:

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

`m` thuộc kiểu map được sử dụng để lưu trữ dữ liệu KV. `filter` tương ứng với một danh sách các hàm lọc được xác định tại mỗi cuộc gọi. `mu` thuộc kiểu mutex để cung cấp sự bảo vệ cho các thành phần khác khi được truy cập và sửa đổi từ nhiều Goroutine cùng lúc.

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

Danh sách các filter được cung cấp trong phương thức `Watch`:

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

Tham số đầu vào của phương thức `Watch` là số giây timeout. Khóa được trả về khi có khóa thay đổi. Nếu không có khóa nào được sửa đổi sau khi timeout thì trả về error. Trong phần hiện thực của `Watch`, mỗi lời gọi được biểu thị bằng một id duy nhất và sau đó hàm filter tương ứng được thêm vào danh sách `p.filter` dựa theo id.

Quá trình đăng ký và khởi động service `KVStoreService` sẽ không được lặp lại. Hãy xem cách sử dụng phương thức `Watch` từ client:

```go
func doClientWork(client *rpc.Client) {
    go func() {
        var keyChanged string
        err := client.Call("KVStoreService.Watch", 30, &keyChanged)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println("watch:", keyChanged)
    } ()

    err := client.Call(
        "KVStoreService.Set", [2]string{"abc", "abc-value"},
        new(struct{}),
    )
    if err != nil {
        log.Fatal(err)
    }

    time.Sleep(time.Second*3)
}
```

Đầu tiên khởi chạy một Goroutine riêng biệt để giám sát khóa thay đổi. Một lời gọi `watch` đồng bộ sẽ block cho đến khi có khóa thay đổi hoặc timeout. Sau đó, khi giá trị KV được thay đổi bằng phương thức `Set`, server sẽ trả về khóa đã thay đổi thông qua phương thức `Watch`. Bằng cách này chúng ta có thể giám sát việc thay đổi trạng thái của khóa.

## 4.3.3 Reverse RPC

RPC bình thường dựa trên cấu trúc client-server. Server của RPC tương ứng với server của mạng và client của RPC cũng tương ứng với client mạng. Tuy nhiên, đối với một số trường hợp đặc biệt, chẳng hạn như khi cung cấp dịch vụ RPC trên mạng nội bộ, nhưng mạng bên ngoài không thể  kết nối với server mạng nội bộ. Trong trường hợp này, có thể sử dụng công nghệ tương tự như reverse proxy. Trước tiên chủ động kết nối với server TCP của mạng bên ngoài từ mạng nội bộ, sau đó cung cấp dịch vụ RPC cho mạng bên ngoài dựa trên kết nối TCP đó.

Sau đây là mã nguồn để khởi động một reverse RPC service:

```go
func main() {
    rpc.Register(new(HelloService))

    for {
        conn, _ := net.Dial("tcp", "localhost:1234")
        if conn == nil {
            time.Sleep(time.Second)
            continue
        }

        rpc.ServeConn(conn)
        conn.Close()
    }
}
```

Reverse RPC service sẽ không còn  cung cấp service lắng nghe TCP, thay vào đó nó  sẽ chủ động kết nối với server TCP của client. RPC service sau đó được cung cấp dựa trên mỗi liên kết TCP được thiết lập.

RPC client  cần cung cấp một service TCP có địa chỉ công khai để chấp nhận request từ RPC server:

```go
func main() {
    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("ListenTCP error:", err)
    }

    clientChan := make(chan *rpc.Client)

    go func() {
        for {
            conn, err := listener.Accept()
            if err != nil {
                log.Fatal("Accept error:", err)
            }

            clientChan <- rpc.NewClient(conn)
        }
    }()

    doClientWork(clientChan)
}
```

Khi mỗi đường link được thiết lập, đối tượng RPC client được khởi tạo dựa trên link đó và gửi tới pipeline `clientChan`.

Client thực hiện lời gọi RPC trong hàm `doClientWork`:

```go
func doClientWork(clientChan <-chan *rpc.Client) {
    client := <-clientChan
    defer client.Close()

    var reply string
    err := client.Call("HelloService.Hello", "hello", &reply)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(reply)
}
```

Đầu tiên nhận vào đối tượng RPC client từ pipeline và sử dụng câu lệnh `defer` để xác định đóng kết nối với client trước khi hàm exit. Kế tiếp là thực hiện lời gọi RPC bình thường.

## 4.3.4 RPC theo ngữ cảnh

Dựa trên ngữ cảnh chúng ta có thể cung cấp những RPC services thích hợp cho những client khác nhau. Ta có thể hỗ trợ các tính năng theo ngữ cảnh bằng cách cung cấp các RPC service cho từng link kết nối.

Đầu tiên thên vào thành phần `conn` ở `HelloService` cho link tương ứng:

```go
type HelloService struct {
    conn net.Conn
}
```

Sau đó bắt đầu một RPC service riêng cho từng link:

```go
func main() {
    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("ListenTCP error:", err)
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal("Accept error:", err)
        }

        go func() {
            defer conn.Close()

            p := rpc.NewServer()
            p.Register(&HelloService{conn: conn})
            p.ServeConn(conn)
        } ()
    }
}
```

Trong phương thức `Hello`, bạn có thể xác định lời gọi RPC cho các link khác nhau dựa trên biến `conn`:

```go
func (p *HelloService) Hello(request string, reply *string) error {
    *reply = "hello:" + request + ", from" + p.conn.RemoteAddr().String()
    return nil
}
```

Dựa vào thông tin ngữ cảnh mà  chúng ta có thể dễ dàng thêm vào một cơ chế xác minh trạng thái đăng nhập đơn giản cho RPC service:

```go
type HelloService struct {
    conn    net.Conn
    isLogin bool
}

func (p *HelloService) Login(request string, reply *string) error {
    if request != "user:password" {
        return fmt.Errorf("auth failed")
    }
    log.Println("login ok")
    p.isLogin = true
    return nil
}

func (p *HelloService) Hello(request string, reply *string) error {
    if !p.isLogin {
        return fmt.Errorf("please login")
    }
    *reply = "hello:" + request + ", from" + p.conn.RemoteAddr().String()
    return nil
}
```

Theo cách này, khi client kết nối tới RPC service, chức năng login sẽ được thực hiện trước, và các service khác có thể thực thi bình thường sau khi login thành công.
