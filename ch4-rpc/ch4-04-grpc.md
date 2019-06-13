# 4.4 Bắt đầu với gRPC

gRPC là một framework RPC mã nguồn mở đa ngôn ngữ được Google phát triển dựa trên Protobuf. Được thiết kế dựa trên giao thức HTTP/2, gRPC có thể cung cấp nhiều dịch vụ dựa trên liên kết HTTP/2, giúp cho framework này thân thiện hơn với thiết bị di động. Phần này sẽ giới thiệu một số cách sử dụng gRPC đơn giản.

## 4.4.1 Kiến trúc gRPC

Kiến trúc gRPC trong Golang được trình bày trong hình 4-1

<p align="center">

<img src="../images/ch4-1-grpc-go-stack.png">
<span align="center">Hình 4-1 gRPC technology stack</span>

</p>

Lớp dưới cùng là giao thức TCP hoặc Unix Socket. Trên đấy phần hiện thực của giao thức HTTP/2. Thư viện gRPC core cho Golang được xây dựng ở lớp kế. Stub code được tạo ra bởi chương trình thông qua plug-in gRPC giao tiếp với thư viện gRPC core.

## 4.4.2 Bắt đầu với gRPC

Từ quan điểm của Protobuf, gRPC không gì khác hơn là một trình tạo code cho interface service. Bây giờ chúng ta sẽ tìm hiểu cách sử dụng gRPC.

Tạo file *hello.proto* và định nghĩa interface `HelloService`:

[>> mã nguồn](../examples/ch4/ch4.4/2-grpc/example-1/hello.proto)

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

```shell
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

gRPC cung cấp hỗ trợ ngữ cảnh cho mỗi lệnh gọi phương thức thông qua tham số `context.Context`. Khi client gọi phương thức, nó có thể cung cấp thông tin ngữ cảnh bổ sung thông qua các tham số tùy chọn của kiểu `grpc.CallOption`.

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

Quá trình khởi động của gRPC  service  tương tự như quá trình khởi động RPC   service của thư viện chuẩn:

```go
func main() {
    grpcServer := grpc.NewServer()
    RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

    lis, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal(err)
    }
    grpcServer.Serve(lis)
}
```

[>> mã nguồn](../examples/ch4/ch4.4/2-grpc/example-1/server/main.go)

Dòng đầu tiên để khởi tạo một đối tượng gRPC service, kế đó phần hiện thực của `HelloServiceImpl` service được đăng ký với grpcServer thông qua  hàm `RegisterHelloServiceServer` (của gRPC plugin). Cuối cùng `grpcServer.Serve(lis)` cung cấp gRPC service trên port `1234`.

Tiếp theo bạn đã có thể kết nối tới gRPC service từ client:

```go
func main() {
    conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := NewHelloServiceClient(conn)
    reply, err := client.Hello(context.Background(), &String{Value: "hello"})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(reply.GetValue())
}
```

[>> mã nguồn](../examples/ch4/ch4.4/2-grpc/example-1/client/main.go)

Trong đó `grpc.Dial` chịu trách nhiệm thiết lập kết nối với dịch vụ gRPC và sau đó hàm `NewHelloServiceClient` xây dựng một đối tượng `HelloServiceClient` dựa trên kết nối đã thiết lập. Client được trả về  là một đối tượng thuộc interface `HelloServiceClient`. Phương thức được xác định bởi interface này có thể gọi phương thức được cung cấp bởi dịch vụ gRPC tương ứng ở server.

Có một sự khác biệt giữa gRPC và framework RPC của thư viện chuẩn: Framework được tạo bởi gRPC không hỗ trợ các cuộc gọi bất đồng bộ. Tuy nhiên, ta có thể chia sẻ  kết nối HTTP/2 cơ bản một cách an toàn  giữa các gRPC trên nhiều Goroutines, vì vậy có thể mô phỏng các lời gọi bất đồng bộ bằng cách block các lời gọi trong Goroutine khác.

## 4.4.3 gRPC flow

RPC là lời gọi hàm từ xa, vì vậy các tham số hàm và giá trị trả về của mỗi cuộc gọi không thể quá lớn, nếu không thời gian phản hồi của mỗi lời gọi sẽ bị ảnh hưởng nghiêm trọng. Do đó, các lời gọi phương thức RPC truyền thống không phù hợp để tải lên và tải xuống trong trường hợp khối lượng dữ liệu lớn. Đồng thời RPC truyền thống không áp dụng cho các mô hình đăng ký và phát hành không chắc chắn về thời gian. Để khắc phục điểm này, framework gRPC cung cấp các  flow cho server và client tương ứng.

Flow một chiều của server hoặc client là trường hợp đặc biệt của flow hai chiều. Chúng tôi thêm phương thức channel hỗ trợ luồng hai chiều trong `HelloService`:

```protobuf
service HelloService {
    rpc Hello (String) returns (String);

    rpc Channel (stream String) returns (stream String);
}```
