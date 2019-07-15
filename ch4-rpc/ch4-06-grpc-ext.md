# 4.6 gRPC và Protobuf extensions

Hiện nay, cộng đồng Open Source đã phát triển rất nhiều extensions xung quanh Protobuf và gRPC, tạo thành một hệ sinh thái to lớn.

## 4.6.1 Validator

Cho đến nay, chúng tôi đã giới thiệu Protobuf phiên bản thứ ba. Ở phiên bản thứ hai của Protobuf có một thuộc tính `default` ở các trường nhằm định nghĩa giá trị mặc định cho nó là một giá trị thuộc kiểu string hoặc kiểu số.

Chúng ta sẽ tạo ra file proto sử dụng phiên bản Protobuf thứ hai:

```proto
syntax = "proto2";

package main;

message Message {
    optional string name = 1 [default = "gopher"];
    optional int32 age = 2 [default = 10];
}
```

Cú pháp dựng sẵn này sẽ được hiện thực thông qua phần mở rộng tính năng của Protobuf. Giá trị mặc định không còn được hỗ trợ trong Protobuf phiên bản thứ ba, nhưng chúng ta có thể mô phỏng giá trị mặc định của chúng bởi một phần mở rộng của option.

Sau đây là phần viết lại của file proto trên với phần mở rộng thuộc cú pháp proto3

```go
syntax = "proto3";

package main;

import "google/protobuf/descriptor.proto";

extend google.protobuf.FieldOptions {
    string default_string = 50000;
    int32 default_int = 50001;
}

message Message {
    string name = 1 [(default_string) = "gopher"];
    int32 age = 2[(default_int) = 10];
}
```

Trong dấu đóng mở ngoặc vuông sau mỗi thành viên là cú pháp mở rộng. Chúng ta sẽ sinh lại mã nguồn Go, chúng sẽ chứa những thông tin liên quan đến phần mở rộng của options.

```go
var E_DefaultString = &proto.ExtensionDesc{
    ExtendedType:  (*descriptor.FieldOptions)(nil),
    ExtensionType: (*string)(nil),
    Field:         50000,
    Name:          "main.default_string",
    Tag:           "bytes,50000,opt,name=default_string,json=defaultString",
    Filename:      "helloworld.proto",
}

var E_DefaultInt = &proto.ExtensionDesc{
    ExtendedType:  (*descriptor.FieldOptions)(nil),
    ExtensionType: (*int32)(nil),
    Field:         50001,
    Name:          "main.default_int",
    Tag:           "varint,50001,opt,name=default_int,json=defaultInt",
    Filename:      "helloworld.proto",
}
```



Chúng ta có thể parse out phần mở rộng của option được định nghĩa trong mỗi thành viên của Message tại thời điểm thực thi bởi kiểu `reflection`, và sau đó parse out gía trị mặc định mà chúng ta đã định nghĩa sẵn từ những thông tin liên quan khác cho phần mở rộng.

Trong cộng đồng Open Source, `github.com/mwitkow/go-proto-validators` đã hiện thực hàm validator rất mạnh mẽ dựa trên phần mở rộng tự nhiên của Protobuf. Để sử dụng validator đầu tiên ta cần phải tải plugin sinh mã nguồn bên dưới

```
$ go get github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
```

Sau đó thêm phần `validation rules` vào các thành viên của Message dựa trên rules của go-proto-validators validator.

```proto
syntax = "proto3";

package main;

import "github.com/mwitkow/go-proto-validators/validator.proto";

message Message {
    string important_string = 1 [
        (validator.field) = {regex: "^[a-z]{2,5}$"}
    ];
    int32 age = 2 [
        (validator.field) = {int_gt: 0, int_lt: 100}
    ];
}
```


Ở phần mở rộng của thành viên được biểu diễn bởi dấu ngoặc vuông, `validator.field` chỉ ra rằng phần mở rộng là một trường tùy chọn trong gói `validator`. Kiểu của `validator.field` là cấu trúc `FieldValidator` được imported trong file `validator.proto`.

Tất cả những validation rules được định nghĩa bởi `FieldValidator` trong file `validator.proto`

```proto
syntax = "proto2";
package validator;

import "google/protobuf/descriptor.proto";

extend google.protobuf.FieldOptions {
    optional FieldValidator field = 65020;
}

message FieldValidator {
    // Uses a Golang RE2-syntax regex to match the field contents.
    optional string regex = 1;
    // Field value of integer strictly greater than this value.
    optional int64 int_gt = 2;
    // Field value of integer strictly smaller than this value.
    optional int64 int_lt = 3;

    // ... more ...
}
```


Từ những chú thích được định nghĩa trong `FieldValidator` chúng ta có thể thấy vài cú pháp cho phần mở rộng của `validator`, ở chỗ `regex` biểu diễn `regular expression` (biểu thức chính quy) cho phần string validation, và `int_gt`, `int_lt` biểu diễn giới hạn của giá trị

Chúng ta dùng lệnh sau để sinh ra mã nguồn của hàm validation

```
$ protoc  \
    --proto_path=${GOPATH}/src \
    --proto_path=${GOPATH}/src/github.com/google/protobuf/src \
    --proto_path=. \
    --govalidators_out=. --go_out=plugins=grpc:.\
    hello.proto
```

> windows: Thay thế ${GOPATH} thành %GOPATH% .

Lệnh trên sẽ gọi chương trình `protoc-gen-govalidators` sau đó sinh ra file với tên `hello.validator.pb.go`

```go
var _regex_Message_ImportantString = regexp.MustCompile("^[a-z]{2,5}$")

func (this *Message) Validate() error {
    if !_regex_Message_ImportantString.MatchString(this.ImportantString) {
        return go_proto_validators.FieldError("ImportantString", fmt.Errorf(
            `value '%v' must be a string conforming to regex "^[a-z]{2,5}$"`,
            this.ImportantString,
        ))
    }
    if !(this.Age > 0) {
        return go_proto_validators.FieldError("Age", fmt.Errorf(
            `value '%v' must be greater than '0'`, this.Age,
        ))
    }
    if !(this.Age < 100) {
        return go_proto_validators.FieldError("Age", fmt.Errorf(
            `value '%v' must be less than '100'`, this.Age,
        ))
    }
    return nil
}
```

Mã nguồn được sinh ra sẽ thêm phương thức `Validate` vào cấu trúc Message để xác định rằng những thành viên được định nghĩa trong nó sẽ thỏa mãn điều kiện ràng buộc trong Protobuf. Bất kể kiểu dữ liệu như thế nào, tất cả những phương thức Validate sẽ dùng chung một signature, do đó cùng một `authentication interface` có thể được chấp nhận.

Thông qua hàm `validation` được sinh ra, chúng sẽ được kết hợp với `gRPC interceptor`, chúng ta có thể dễ dàng thẩm định giá trị của tham số đầu vào và kết quả trả về của mỗi hàm.

## 4.6.2 REST interface

gRPC service thường được dùng trong việc giao tiếp giữa các cluster trong hệ thống. Nếu một service được yêu cầu phải giao tiếp với bên ngoài, thì thường một REST interface sẽ được sinh ra để làm việc đó. Để thuận tiện, phía front-end qua Javascript và phía back-end sẽ giao tiếp với nhau thông qua REST interface. Cộng đồng opensource đã hiện thực project với tên gọi là grpc-gateway, nó giúp chúng ta chuyển các yêu cầu REST HTTP thành các yêu cầu gRPC HTTP2.

Nguyên tắc hoạt động bên dưới của grpc-gateway sẽ như sau:

![](../images/ch4-2-grpc-gateway.png)

Hình 4-2: gRPC-Gateway workflow

Chúng ta có thể **sinh grpc-gateway dựa trên docker** theo [hướng dẫn](https://medium.com/zalopay-engineering/buildingdocker-grpc-gateway-e2efbdcfe5c) hoặc theo cách thông thường như sau:

Bằng việc thêm vào những thông tin liên quan trong phần routing trong file Protobuf, những quá trình liên quan đến routing sẽ được sinh ra thông qua file được tùy biến trong mã nguồn plugin, và cuối cùng REST request sẽ chuyển tiếp đến những service back-end được chạy bên dưới.

Phần mở rộng của routing sẽ được cung cấp thông qua metadata của Protobuf như sau

```proto
syntax = "proto3";

package main;

import "google/api/annotations.proto";

message StringMessage {
  string value = 1;
}

service RestService {
    rpc Get(StringMessage) returns (StringMessage) {
        option (google.api.http) = {
            get: "/get/{value}"
        };
    }
    rpc Post(StringMessage) returns (StringMessage) {
        option (google.api.http) = {
            post: "/post"
            body: "*"
        };
    }
}
```

Đầu tiên chúng ta sẽ định nghĩa các phương thức POST và GET cho gRPC, và sau đó chúng ta sẽ thêm vào phần thông tin liên quan đến routing trong phương thức tương ứng qua cú pháp meta-extension. Đường dẫn "/get/{value}" sẽ tương ứng với phương thức GET và `{value}` tương ứng với một số thành viên trong parameter, kết quả có thể được trả về theo định dạng json. Phương thức POST sẽ tương ứng với đường dẫn "/post" và phần thân chứa thông tin về request cũng định dạng theo kiểu json.

Sau đó chúng ta cài đặt plugin `protoc-gen-grpc-gateway` với những lệnh sau:

```
$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
```

Sau đó chúng ta sinh ra mã nguồn xử lý routing cần cho grpc-gateway thông qua plugin sau:

```
$ protoc -I/usr/local/include -I. \
    -I$GOPATH/src \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --grpc-gateway_out=. --go_out=plugins=grpc:.\
    hello.proto
```

> windows: Thay thế ${GOPATH} với %GOPATH%.

Plugin sẽ sinh ra hàm RegisterRestServiceHandlerFromEndpoint tương ứng tới RestService service.

```go
func RegisterRestServiceHandlerFromEndpoint(
    ctx context.Context, mux *runtime.ServeMux, endpoint string,
    opts []grpc.DialOption,
) (err error) {
    ...
}
```

Hàm `RegisterRestServiceHandlerFromEndpoint` được dùng để chuyển tiếp những request được định nghĩa trong REST interface đến gRPC service thực sự. Sau khi registering the route handle, chúng ta sẽ bắt đầu web service.

```go
func main() {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := runtime.NewServeMux()

    err := RegisterRestServiceHandlerFromEndpoint(
        ctx, mux, "localhost:5000",
        []grpc.DialOption{grpc.WithInsecure()},
    )
    if err != nil {
        log.Fatal(err)
    }

    http.ListenAndServe(":8080", mux)
}
```


Bắt đầu gRPC service tại port 5000

```go
type RestServiceImpl struct{}

func (r *RestServiceImpl) Get(ctx context.Context, message *StringMessage) (*StringMessage, error) {
    return &StringMessage{Value: "Get hi:" + message.Value + "#"}, nil
}

func (r *RestServiceImpl) Post(ctx context.Context, message *StringMessage) (*StringMessage, error) {
    return &StringMessage{Value: "Post hi:" + message.Value + "@"}, nil
}
func main() {
    grpcServer := grpc.NewServer()
    RegisterRestServiceServer(grpcServer, new(RestServiceImpl))
    lis, _ := net.Listen("tcp", ":5000")
    grpcServer.Serve(lis)
}
```

Đầu tiên, chúng ta tạo ra route handler thông qua hàm thực thi `runtime.NewServeMux()`, và sau đó chuyển đổi REST interface liên quan đến RestService service đến phần subsequent gRPC service thông qua hàm RegisterRestServiceHandlerFromEndpoint. Lớp `runtime.ServeMux` sẽ được hỗ trợ bởi `grpc-gateway` chúng được hiện thực bên dưới interface `http.Handler`, do đó chúng ta có thể dùng những hàm liên quan được cung cấp trong thư viện chuẩn

Sau khi tất cả các gRPC và REST services được khởi động, chúng ta có thể khởi tạo request REST service với curl:

```
$ curl localhost:8080/get/gopher
{"value":"Get: gopher"}

$ curl localhost:8080/post -X POST --data '{"value":"grpc"}'
{"value":"Post: grpc"}
```

Khi chúng ta publishing REST interface, chúng ta thông thường sẽ cung cấp một swagger file theo định dạng interface được mô tả bên dưới.

```
$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

$ protoc -I. \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=. \
  hello.proto
```

Sau đó file `hello.swagger.json` sẽ được sinh ra. Trong trường hợp này, chúng ta có thể dùng `swagger-ui project` để cung cấp tài liệu `REST interface` và testing dưới dạng web pages.


## 4.6.3 Nginx

Phiên bản [Nginx](https://www.nginx.com/) hỗ trợ `gRPC`. Bên dưới back-end của nhiều gRPC services có thể được tổng hợp trong một Nginx service thông qua Nginx. Cùng một thời điểm, Nginx sẽ hỗ trợ khả năng register nhiều back-end tới cùng gRPC service, chúng sẽ làm cho việc hỗ trợ load balancing (cân bằng tải) dễ dàng hơn. Những extension của Nginx's gRPC là một chủ đề lớn, tốt hơn chúng ta nên tham khảo đến những tài liệu liên quan nói về chúng.