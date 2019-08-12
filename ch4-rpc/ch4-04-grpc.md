# 4.4 Bắt đầu với gRPC

<div align="center">
	<img src="../images/grpc.png" width="800">
	<br/>
</div>
<br/>

[gRPC](https://grpc.io/) là một framework RPC opensource đa ngôn ngữ được Google phát triển dựa trên [Protobuf](https://developers.google.com/protocol-buffers/) và giao thức HTTP/2. Phần này sẽ giới thiệu một số cách sử dụng gRPC để xây dựng service đơn giản.

## 4.4.1 Kiến trúc gRPC

Kiến trúc gRPC trong Golang:

<div align="center">
    <img src="../images/ch4-1-grpc-go-stack.png" width="450">
	<br/>
	<span align="center">
		<i>gRPC technology stack</i>
	</span>
</div>
<br/>

Lớp dưới cùng là giao thức TCP hoặc Unix Socket. Ngay trên đấy là phần hiện thực của giao thức HTTP/2. Thư viện gRPC core cho Golang được xây dựng ở lớp kế. Stub code được tạo ra bởi chương trình thông qua plug-in gRPC giao tiếp với thư viện gRPC core.

## 4.4.2 Bắt đầu với gRPC

Từ quan điểm của Protobuf, gRPC không gì khác hơn là một trình tạo code cho interface service.

Tạo file *hello.proto* và định nghĩa interface `HelloService`:

```protobuf
syntax = "proto3";

package main;

message String {
    string value = 1;
}

service HelloService {
    rpc Hello (String) returns (String);
}
```

Tạo gRPC code sử dụng hàm dựng sẵn trong gRPC plugin từ protoc-gen-go:

```sh
$ protoc --go_out=plugins=grpc:. hello.proto
```

gRPC plugin tạo ra các interface khác nhau cho server và client:

```go
type HelloServiceServer interface {
    Hello(context.Context, *String) (*String, error)
}

type HelloServiceClient interface {
    Hello(context.Context, *String, ...grpc.CallOption) (*String, error)
}
```

gRPC cung cấp hỗ trợ context cho mỗi lệnh gọi phương thức thông qua tham số `context.Context`. Khi client gọi phương thức, nó có thể cung cấp thông tin context bổ sung thông qua các tham số tùy chọn của kiểu `grpc.CallOption`.

`HelloSercieServer` interface dựa trên server có thể reimplement service `HelloService`:

```go
type HelloServiceImpl struct{}

func (p *HelloServiceImpl) Hello(
    ctx context.Context, args *String,
) (*String, error) {
    reply := &String{Value: "hello:" + args.GetValue()}
    return reply, nil
}
```

Quá trình khởi động của gRPC service  tương tự như quá trình khởi động RPC service của thư viện chuẩn:

```go
func main() {
    // khởi tạo một đối tượng gRPC service
    grpcServer := grpc.NewServer()

    // đăng ký service với grpcServer (của gRPC plugin)
    RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

    // cung cấp gRPC service trên port `1234`
    lis, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal(err)
    }
    grpcServer.Serve(lis)
}
```

Tiếp theo bạn đã có thể kết nối tới gRPC service từ client:

```go
func main() {
    // thiết lập kết nối với gRPC service
    conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // xây dựng đối tượng `HelloServiceClient` dựa trên kết nối đã thiết lập
    client := NewHelloServiceClient(conn)
    reply, err := client.Hello(context.Background(), &String{Value: "hello"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(reply.GetValue())
}
```

Có một sự khác biệt giữa gRPC và framework RPC của thư viện chuẩn: gRPC không hỗ trợ các gọi asynchronous. Tuy nhiên, ta có thể chia sẻ  kết nối HTTP/2 trên nhiều Goroutines, vì vậy có thể mô phỏng các lời gọi bất đồng bộ bằng cách block các lời gọi trong Goroutine khác.

## 4.4.3 gRPC streaming

RPC là lời gọi hàm từ xa, vì vậy các tham số hàm và giá trị trả về của mỗi cuộc gọi không thể quá lớn, nếu không thời gian phản hồi của mỗi lời gọi sẽ bị ảnh hưởng nghiêm trọng. Do đó, các lời gọi phương thức RPC truyền thống không phù hợp để tải lên và tải xuống trong trường hợp khối lượng dữ liệu lớn. Để khắc phục điểm này, framework gRPC cung cấp các chức năng stream cho phía server và client tương ứng.

Ta viết thêm phương thức channel hỗ trợ luồng hai chiều (Bidirect Streaming) trong `HelloService`:

```protobuf
service HelloService {
    rpc Hello (String) returns (String);

    // nhận vào tham số một stream và trả về giá trị là một stream.
    rpc Channel (stream String) returns (stream String);
}
```

Tạo lại code để thấy định nghĩa mới được thêm vào phương thức kiểu channel trong interface:

```go
type HelloServiceServer interface {
    Hello(context.Context, *String) (*String, error)

    // tham số kiểu HelloService_ChannelServer được sử dụng
    // để liên lạc hai chiều với client.
    Channel(HelloService_ChannelServer) error
}

type HelloServiceClient interface {
    Hello(ctx context.Context, in *String, opts ...grpc.CallOption) (
        *String, error,
    )

    // trả về giá trị trả về thuộc kiểu `HelloService_ChannelClient`
    // có thể được sử dụng để liên lạc hai chiều với server.
    Channel(ctx context.Context, opts ...grpc.CallOption) (
        HelloService_ChannelClient, error,
    )
}
```

`HelloService_ChannelServer` và `HelloService_ChannelClient` thuộc interface:

```go
type HelloService_ChannelServer interface {
    Send(*String) error
    Recv() (*String, error)
    grpc.ServerStream
}

type HelloService_ChannelClient interface {
    Send(*String) error
    Recv() (*String, error)
    grpc.ClientStream
}
```

Có thể thấy các interface hỗ trợ server và client stream đều có định nghĩa phương thức `Send` và `Recv` cho giao tiếp hai chiều của dữ liệu streaming.

Bây giờ ta có thể hiện thực các streaming service:

```go
func (p *HelloServiceImpl) Channel(stream HelloService_ChannelServer) error {
    for {
        // Server nhận dữ liệu được gửi từ client
        // trong vòng lặp.
        args, err := stream.Recv()
        if err != nil {
            // Nếu gặp `io.EOF`, client stream sẽ đóng.
            if err == io.EOF {
                return nil
            }
            return err
        }

        reply := &String{Value: "hello:" + args.GetValue()}

        // Dữ liệu trả về được  gửi đến client
        // thông qua stream và việc gửi nhận
        // dữ liệu stream hai chiều là hoàn toàn độc lập
        err = stream.Send(reply)
        if err != nil {
            return err
        }
    }
}
```

Client cần gọi phương thức Channel để lấy đối tượng stream trả về:

```go
stream, err := client.Channel(context.Background())
if err != nil {
    log.Fatal(err)
}
```

Ở phía client ta thêm vào các thao tác gửi và nhận trong các Goroutine riêng biệt. Trước hết là để gửi dữ liệu tới server:

```go
go func() {
    for {
        if err := stream.Send(&String{Value: "hi"}); err != nil {
            log.Fatal(err)
        }
        time.Sleep(time.Second)
    }
}()
```

Kế đến là nhận dữ liệu trả về từ server trong vòng lặp:

```go
for {
    reply, err := stream.Recv()
    if err != nil {
        if err == io.EOF {
            break
        }
        log.Fatal(err)
    }
    fmt.Println(reply.GetValue())
}
```

Bạn có thể xem code cụ thể tại [link](../examples/ch4/ch4.4/3-grpc-flow/server/main.go)

## 4.4.4 Mô hình Publishing - Subscription

Trong phần trước chúng ta đã hiện thực phiên bản đơn giản của phương thức `Watch` dựa trên thư viện RPC dựng sẵn của Go. Ý tưởng đó có thể sử dụng cho hệ thống publish-subscribe, nhưng bởi vì RPC thiếu đi cơ chế streaming nên nó chỉ có thể trả về 1 kết quả trong 1 lần. Trong chế độ publish-subscribe, hành động publish đưa ra bởi *caller* giống với lời gọi hàm thông thường, trong khi subscriber bị động thì giống với *receiver* trong gRPC client stream một chiều. Bây giờ ta có thể thử xây dựng một hệ thống publish - subscribe dựa trên đặc điểm stream của gRPC.


<div align="center">

<img src="../images/pubsub.png" width="400">
<br/>
<span align="center"><i>Nhắc lại mô hình Pub/Sub</i></span>
    <br/>

</div>

Publishing - Subscription là một mẫu thiết kế thông dụng và đã có nhiều hiện thực của mẫu thiết kế này trong cộng đồng opensource. Đoạn code sau đây hiện thực cơ chế publish - subscription dựa trên package pubsub:

```go
import (
    "github.com/moby/moby/pkg/pubsub"
)

func main() {
    //  xây dựng một đối tượng để publish
    p := pubsub.NewPublisher(100*time.Millisecond, 10)

    // subscribe các topic "golang"
    golang := p.SubscribeTopic(func(v interface{}) bool {
        if key, ok := v.(string); ok {
            if strings.HasPrefix(key, "golang:") {
                return true
            }
        }
        return false
    })

    // subscribe các topic "docker"
    docker := p.SubscribeTopic(func(v interface{}) bool {
        if key, ok := v.(string); ok {
            if strings.HasPrefix(key, "docker:") {
                return true
            }
        }
        return false
    })

    go p.Publish("hi")
    go p.Publish("golang: https://golang.org")
    go p.Publish("docker: https://www.docker.com/")
    time.Sleep(1)

    go func() {
        fmt.Println("golang topic:", <-golang)
    }()
    go func() {
        fmt.Println("docker topic:", <-docker)
    }()

    <-make(chan bool)
}
```

Giờ  chúng ta thử cung cấp một hệ thống publishing-subscription khác mạng dựa trên gRPC và pubsub package. Đầu tiên định nghĩa một service publish subscription interface bằng protobuf:

```protobuf
service PubsubService {
    //  phương thức RPC thông thường
    rpc Publish (String) returns (String);

    // service server streaming
    rpc Subscribe (String) returns (stream String);
}
```

gRPC plugin sẽ tạo ra interface tương ứng cho server và client:

```go
type PubsubServiceServer interface {
    Publish(context.Context, *String) (*String, error)
    Subscribe(*String, PubsubService_SubscribeServer) error
}
type PubsubServiceClient interface {
    Publish(context.Context, *String, ...grpc.CallOption) (*String, error)
    Subscribe(context.Context, *String, ...grpc.CallOption) (
        PubsubService_SubscribeClient, error,
    )
}

type PubsubService_SubscribeServer interface {
    Send(*String) error
    grpc.ServerStream
}
```

Bởi vì `Subscribe` là stream 1 chiều phía server nên chỉ có phương thức `Send` được tạo ra trong interface `HelloService_SubscribeServer`.

Sau đó có thể hiện thực các service publish và subscribe như sau:

```go
type PubsubService struct {
    pub *pubsub.Publisher
}

func NewPubsubService() *PubsubService {
    return &PubsubService{
        pub: pubsub.NewPublisher(100*time.Millisecond, 10),
    }
}
```

Kế đến là các phương thức publishing và subscription:

```go
func (p *PubsubService) Publish(
    ctx context.Context, arg *String,
) (*String, error) {
    p.pub.Publish(arg.GetValue())
    return &String{}, nil
}

func (p *PubsubService) Subscribe(
    arg *String, stream PubsubService_SubscribeServer,
) error {
    ch := p.pub.SubscribeTopic(func(v interface{}) bool {
        if key, ok := v.(string); ok {
            if strings.HasPrefix(key,arg.GetValue()) {
                return true
            }
        }
        return false
    })

    for v := range ch {
        if err := stream.Send(&String{Value: v.(string)}); err != nil {
            return err
        }
    }

    return nil
}
```

Hàm `main` cho phép đăng mới thông tin từ client tới server:

```go
func main() {
    conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := NewPubsubServiceClient(conn)

    _, err = client.Publish(
        context.Background(), &String{Value: "golang: hello Go"},
    )
    if err != nil {
        log.Fatal(err)
    }
    _, err = client.Publish(
        context.Background(), &String{Value: "docker: hello Docker"},
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

Chi tiết: [clientpub](../examples/ch4/ch4.4/4-pubsub/clientpub/main.go).

Sau đó có thể subscribe thông tin đó từ một client khác:

```go
func main() {
    conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := NewPubsubServiceClient(conn)
    stream, err := client.Subscribe(
        context.Background(), &String{Value: "golang:"},
    )
    if err != nil {
        log.Fatal(err)
    }

    for {
        reply, err := stream.Recv()
        if err != nil {
            if err == io.EOF {
                break
            }
            log.Fatal(err)
        }

        fmt.Println(reply.GetValue())
    }
}
```

Chi tiết: [clientsub](../examples/ch4/ch4.4/4-pubsub/clientsub/main.go).

Cho đến giờ chúng ta đã hiện thực được service publishing và subscription khác mạng dựa trên gRPC. Trong phần kế tiếp chúng ta sẽ xét một số ứng dụng nâng cao hơn của Go trong gRPC.
