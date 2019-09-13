# 3.6. Công cụ grpcurl

Bản thân Protobuf đã có chức năng phản chiếu (reflection) lại file Proto của đối tượng khi thực thi. gRPC cũng cung cấp một package reflection để thực hiện các truy vấn cho  gRPC service. Mặc dù gRPC có một hiện thực bằng C++ của công cụ `grpc_cli`, có thể được sử dụng để truy vấn danh sách gRPC hoặc gọi phương thức gRPC, nhưng bởi vì phiên bản đó cài đặt khá  phức tạp nên ở đây chúng ta sẽ dùng công cụ `grpcurl` được hiện thực thuần bằng Golang. Phần này ta sẽ cùng tìm hiểu cách sử dụng công cụ này.

## 3.6.1 Khởi động một reflection service

Chỉ có duy nhất hàm `Register` trong package reflection, hàm này dùng để đăng ký `grpc.Server` với reflection service. Trong document của package có hướng dẫn như sau:

```go
import (
    "google.golang.org/grpc/reflection"
)

func main() {
    s := grpc.NewServer()
    pb.RegisterYourOwnServer(s, &server{})

    // đăng ký reflection service trên gRPC server.
    reflection.Register(s)

    s.Serve(lis)
}
```

Nếu gRPC reflection service được khởi chạy thì các gRPC service   có thể được truy vấn hoặc gọi ra bằng reflection service do package reflection cung cấp.

## 3.6.2 Xem danh sách service

Grpcurl là công cụ được cộng đồng Open source của Golang phát triển, quá trình cài đặt như sau:

```sh
$ go get github.com/fullstorydev/grpcurl
$ go install github.com/fullstorydev/grpcurl/cmd/grpcurl
```

Sử dụng phổ biến nhất trong grpcurl là lệnh list, được sử dụng để lấy danh sách các service hoặc các phương thức trong service. Ví dụ, `grpcurl localhost:1234 list` là lệnh sẽ nhận được một danh sách các service gRPC trên port 1234 ở localhost.

Khi sử dụng grpcurl với giao thức TLS ta cần chỉ định các đường dẫn tới public key `-cert` và private key `-key`. Đối với  gRPC service không có giao thức TLS, quy trình xác minh chứng chỉ TLS có thể bỏ qua bằng   tham số `-plaintext`. Nếu đó là giao thức Unix Socket, cần chỉ định tham số `-unix`.

Nếu các file public và private key chưa được cấu hình và quá trình xác minh chứng chỉ bị bỏ qua, ta có thể sẽ gặp lỗi như sau:

```sh
$ grpcurl localhost:1234 list
Failed to dial target host "localhost:1234": tls: first record does not \
look like a TLS handshake
```

Nếu gRPC service bình thường nhưng được khởi động reflection service thì  sẽ  có thông báo lỗi:

```sh
$ grpcurl -plaintext localhost:1234 list
Failed to list services: server does not support the reflection API
```

Giả định rằng gRPC sercie đã được kích hoạt reflection service, file Protobuf của service như sau:

```protobuf
syntax = "proto3";

package HelloService;

message String {
    string value = 1;
}

service HelloService {
    rpc Hello (String) returns (String);
    rpc Channel (stream String) returns (stream String);
}
```

Kết quả với lệnh `list`:

```sh
$ grpcurl -plaintext localhost:1234 list
HelloService.HelloService
grpc.reflection.v1alpha.ServerReflection
```

Trong đó `HelloService.HelloService` là service được định nghĩa trong file protobuf. `ServerReflection` là reflection service được package reflection đăng ký. Thông qua service này chúng ta có thể truy vấn thông tin của tất cả các gRPC service bao gồm chính nó.

## 3.6.3 Danh sách các phương thức của service

Nếu tiếp tục sử dụng lệnh `list` ta có thể xem được cả danh sách các phương thức trong `HelloService`:

```sh
$ grpcurl -plaintext localhost:1234 list HelloService.HelloService
Channel
Hello
```

Từ kết quả cho thấy service này cung cấp 2 phương thức là `Channel` và `Hello`, tương ứng với các định nghĩa trong file Protobuf.

Nếu muốn biết chi tiết của từng phương thức, ta có thể sử dụng câu lệnh `describe`:

```sh
$ grpcurl -plaintext localhost:1234 describe HelloService.HelloService
HelloService.HelloService is a service:
service HelloService {
  rpc Channel ( stream .HelloService.String ) returns ( stream .HelloService.String );
  rpc Hello ( .HelloService.String ) returns ( .HelloService.String );
```

Kết quả là danh sách các phương thức có trong service cùng với mô tả các tham số input cũng như giá trị trả về tương ứng của chúng.

## 3.6.4 Lấy thông tin kiểu dữ liệu

Sau khi có được danh sách các phương thức và kiểu của giá tị trả về, chúng ta có thể tiếp tục xem thông tin chi tiết hơn về kiểu của các biến này. Sau đây là các sử dụng lệnh `describe` để xem thông tin của tham số `HelloService.String`:

```sh
$ grpcurl -plaintext localhost:1234 describe HelloService.String
HelloService.String is a message:
message String {
  string value = 1;
}
```

Kết quả trả về đúng với mô tả trong file protobuf của service.

## 3.6.5 Lệnh gọi phương thức

Ta có thể gọi phương thức gRPC bằng cách truyền thêm tham số `-d` và một chuỗi json như input của hàm và gọi tới phương thức `Hello` trong `HelloService`, chi tiết như sau:

```sh
$ grpcurl -plaintext -d '{"value": "gopher"}' localhost:1234 HelloService.HelloService/Hello

{
  "value": "hello:gopher"
}
```

Nếu có tham số `-d`, `@` nghĩa là đọc vào tham số dạng json từ input chuẩn (stdin), cách này thường dùng để test các phương thức stream.

Ví dụ sau đây kết nối tới phương thức stream tên `Channel` và đọc tham số input stream từ input chuẩn:

```sh
$ grpcurl -plaintext -d @ localhost:1234 HelloService.HelloService/Channel
{"value":"gopher-vn"}
{
  "value": "hello:gopher-vn"
}

{"value": "vietnamese-vng"}
{
  "value": "hello:vietnamese-vng"
}
```
