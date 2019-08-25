# 3.6. gRPC và Protobuf extensions

Hiện nay, cộng đồng Open source đã phát triển rất nhiều extensions xung quanh Protobuf và gRPC, tạo thành một hệ sinh thái to lớn. Ở phần này sẽ trình bày về một số extensions thông dụng.

## 3.6.1 Validator

Cho đến nay, Protobuf đã có phiên bản thứ ba. Ở phiên bản thứ hai của Protobuf có một thuộc tính `default` ở các trường nhằm định nghĩa giá trị mặc định cho nó là một giá trị thuộc kiểu string hoặc kiểu số.

Chúng ta sẽ tạo ra file proto sử dụng phiên bản Protobuf thứ hai:

***hello.proto (proto2):***

```go
// phiên bản protobuf
syntax = "proto2";
// định nghĩa tên package được sinh ra
package main;
// định nghĩa đối tượng dữ liệu
message Message {
    // nếu không khởi trị, thì giá trị mặc định của name là "gopher"
    string name = 1 [default = "gopher"];
    // tương tự, giá trị mặc định của age là 10
    int32 age = 2 [default = 10];
}
```

Cú pháp này sẽ được hiện thực thông qua phần mở rộng tính năng của Protobuf. Giá trị mặc định không còn được hỗ trợ trong Protobuf phiên bản thứ ba, nhưng chúng ta có thể mô phỏng giá trị mặc định của chúng bởi một phần mở rộng của option.

Sau đây là phần viết lại của file proto trên với phần mở rộng thuộc cú pháp proto3:

***hello.proto (proto3):***

```go
// phiên bản hiện tại là proto3
syntax = "proto3";
package main;
// import phần mở rộng của protobuf
import "google/protobuf/descriptor.proto";
// định nghĩa một số trường trong phần mở rộng
extend google.protobuf.FieldOptions {
    // những con số như: 50000, 50001 là duy nhất cho mỗi trường
    string default_string = 50000;
    int32 default_int = 50001;
}
// định nghĩa nội dung message
message Message {
    // default_string là giá trị mặc định cho name
    string name = 1 [(default_string) = "gopher"];
    // tương tự, age sẽ có giá trị 10 nếu không khởi trị
    int32 age = 2[(default_int) = 10];
}
```

Trong dấu đóng mở ngoặc vuông sau mỗi trường trong message là một cú pháp mở rộng. Chúng ta sẽ tạo lại mã nguồn Go dựa trên những thông tin liên quan đến phần mở rộng của options. Phần mã nguồn sinh ra có một số nội dung dựa trên phần mở rộng như sau:

***hello.pb.go:***

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

Chúng ta có thể parse out phần mở rộng của option được định nghĩa trong mỗi thành viên của Message tại thời điểm thực thi bởi kiểu `reflection`, và sau đó parse out giá trị mặc định mà chúng ta đã định nghĩa sẵn từ những thông tin liên quan khác cho phần mở rộng.

Trong cộng đồng Open source, thư viện [go-proto-validators](github.com/mwitkow/go-proto-validators) là một extension của protobuf có chức năng validator rất mạnh mẽ dựa trên phần mở rộng tự nhiên của Protobuf. Để sử dụng validator đầu tiên ta cần phải tải plugin sinh mã nguồn bên dưới:

```sh
$ go get github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
```

Sau đó thêm phần `validation rules` vào các thành viên của Message dựa trên rules của go-proto-validators validator.

***hello.proto:*** (dùng thư viện [validator](https://github.com/mwitkow/go-proto-validators))
```go
syntax = "proto3";

package main;
// import file validator.proto 
import "github.com/mwitkow/go-proto-validators/validator.proto";
// định nghĩa message
message Message {
    // dấu ngoặc vuông mang ý nghĩa là phần tùy chọn
    string important_string = 1 [
        // regex sẽ validate trường important_string đúng theo syntax hay không
        (validator.field) = {regex: "^[a-z]{2,5}$"}
    ];
    int32 age = 2 [
        // tương tự, giá trị của a sẽ được validate lớn hơn 0 và nhỏ hơn 100
        (validator.field) = {int_gt: 0, int_lt: 100}
    ];
}
```

Tất cả những validation rules được định nghĩa trong message `FieldValidator` trong file [validator.proto](https://github.com/mwitkow/go-proto-validators/blob/master/validator.proto). Trong đó ta sẽ thấy một số trường được dùng ở ví dụ trên như sau:

***mwitkow/go-proto-validators/validator.proto:***

```go
syntax = "proto2";
package validator;

import "google/protobuf/descriptor.proto";

extend google.protobuf.FieldOptions {
    optional FieldValidator field = 65020;
}

message FieldValidator {
    // sử dụng Golang RE2-syntax regex để match với nội dung các field
    optional string regex = 1;
    // giá trị của biến integer bình thường lớn hơn giá trị này.
    optional int64 int_gt = 2;
    // giá trị của biến integer bình thường nhỏ hơn giá trị này.
    optional int64 int_lt = 3;

    // ...
}
```

Phần chú thích của mỗi trường ở trên sẽ cho chúng ta thông tin về chức năng của chúng. Sau khi chọn được các chức năng validate cần thiết, chúng ta dùng lệnh sau để sinh ra mã nguồn validator:

```sh
$ protoc  \
    --proto_path=${GOPATH}/src \
    --proto_path=${GOPATH}/src/github.com/google/protobuf/src \
    --proto_path=. \
    --govalidators_out=. --go_out=plugins=grpc:.\
    hello.proto

// Trong đó:
// - proto_path: đường dẫn đến tất cả các file .proto được sử dụng
// - govalidators_out: plugin sinh ra mã nguồn validator
// Chú ý:
// - Trong Windows, ta thay thế ${GOPATH} thành %GOPATH%
```

Lệnh trên sẽ gọi chương trình `protoc-gen-govalidators` để sinh ra file với tên `hello.validator.pb.go`, nội dung của nó sẽ như sau:

***hello.validator.pb.go:***

```go
// định nghĩa chuỗi regex
var _regex_Message_ImportantString = regexp.MustCompile("^[a-z]{2,5}$")
// hàm Validate() sẽ chạy các rules và bắt lỗi nếu có
func (this *Message) Validate() error {
    // rule 1 kiểm tra ImportantString có theo regex hay không, nếu có lỗi sẽ ném ra
    if !_regex_Message_ImportantString.MatchString(this.ImportantString) {
        return go_proto_validators.FieldError("ImportantString", fmt.Errorf(
            `value '%v' must be a string conforming to regex "^[a-z]{2,5}$"`,
            this.ImportantString,
        ))
    }
    // rule 2 kiểm tra Age > 0 hay không, nếu có lỗi sẽ ném ra
    if !(this.Age > 0) {
        return go_proto_validators.FieldError("Age", fmt.Errorf(
            `value '%v' must be greater than '0'`, this.Age,
        ))
    }
    // rule 3 kiểm tra Age < 100 hay không, nếu có lỗi sẽ ném ra
    if !(this.Age < 100) {
        return go_proto_validators.FieldError("Age", fmt.Errorf(
            `value '%v' must be less than '100'`, this.Age,
        ))
    }
    // trả về nil nếu kiểm tra tất cả các rules trên đều hợp lệ
    return nil
}
```

Thông qua hàm Validate() được sinh ra, chúng có thể được kết hợp với `gRPC interceptor`, chúng ta có thể dễ dàng validate giá trị của tham số đầu vào và kết quả trả về của mỗi hàm.

## 3.6.2 REST interface

Hiện nay RESTful JSON API vẫn là sự lựa chọn hàng đầu cho các ứng dụng web hay mobile. Vì tính tiện lợi và dễ dùng của RESTful API nên chúng ta vẫn sử dụng nó để frontend có thể giao tiếp với hệ thống backend. Nhưng khi chúng ta sử dụng framework gRPC của Google để xây dựng các service. Các service sử dụng gRPC thì dễ dàng trao đổi dữ liệu với nhau dựa trên giao thức HTTP/2 và protobuf, nhưng ở phía frontend lại sử dụng [RESTful API](https://restfulapi.net/) API hoạt động trên giao thức HTTP/1. Vấn đề đặt ra là chúng ta cần phải chuyển đổi các yêu cầu RESTful API thành các yêu cầu gRPC để hệ thống các service gRPC có thể hiểu được.

Cộng đồng Open source đã hiện thực một project với tên gọi là [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway), nó sẽ sinh ra một proxy có vai trò chuyển các yêu cầu REST HTTP thành các yêu cầu gRPC HTTP2.


<div align="center">
	<img src="../images/ch4-2-grpc-gateway.png">
	<br/>
	<span align="center">
		<i>gRPC-Gateway workflow</i>
	</span>
</div>

Trong file Protobuf (chỉ có ở proto3), chúng ta sẽ thêm thông tin phần routing ứng với các hàm trong gRPC service, để dựa vào đó grpc-gateway sẽ sinh ra mã nguồn proxy tương ứng.

***rest_service.proto:***

```go
// phiên bản proto3
syntax = "proto3";
// tên package được sinh ra
package main;
// chú ý: import annotations.proto để dùng chức năng grpc-gateway
import "google/api/annotations.proto";
// định nghĩa message trao đổi
message StringMessage {
  string value = 1;
}
// định nghĩa RestService 
service RestService {
    // định nghĩa hàm RPC Get trong service
    rpc Get(StringMessage) returns (StringMessage) {
        // nội dung phần option trong này định nghĩa Rest API ra bên ngoài
        option (google.api.http) = {
            // get: là tên phương thức được sử dụng
            get: "/get/{value}"
            // "/get/{value}" : là đường dẫn uri,
            // trong đó {value} được pass vào uri là nội dung StringMessage request
        };
    }
    // định nghĩa hàm RPC Post trong service
    rpc Post(StringMessage) returns (StringMessage) {
        option (google.api.http) = {
            // dùng phương thức post
            post: "/post"
            // StringMessage sẽ dưới dạng chuỗi Json khi gửi Request (vd: '{"value":"Hello, World"}')
            body: "*"
        };
    }
}
```

Chúng ta cài đặt plugin protoc-gen-grpc-gateway với những lệnh sau:

```sh
$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
```

Sau đó chúng ta sinh ra mã nguồn routing cho grpc-gateway thông qua plugin sau:

```sh
$ protoc -I/usr/local/include -I. \
    -I$GOPATH/src \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --grpc-gateway_out=. --go_out=plugins=grpc:.\
    hello.proto

// Trong windows: Thay thế ${GOPATH} với %GOPATH%.
```

Plugin sẽ sinh ra hàm RegisterRestServiceHandlerFromEndpoint() cho RestService service như sau:

```go
func RegisterRestServiceHandlerFromEndpoint(
    ctx context.Context, mux *runtime.ServeMux, endpoint string,
    opts []grpc.DialOption,
) (err error) {
    ...
}
```

Hàm RegisterRestServiceHandlerFromEndpoint được dùng để chuyển tiếp những request được định nghĩa trong REST interface đến gRPC service. Sau khi registering các Route handle, chúng ta sẽ chạy proxy web service trong hàm main như sau:

***proxy/main.go:***

```go
func main() {
    // khai báo biến context để xử lý signal kết thúc goroutine
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    // hàm cancel() sẽ kích hoạt ctx.Done()
    defer cancel()
    // mux được dùng cho việc routing
    mux := runtime.NewServeMux()
    // gọi hàm để đăng kí RestService cho proxy
    err := RegisterRestServiceHandlerFromEndpoint(
        // truyền vào biến ctx, mux, và địa chỉ gRPC service
        ctx, mux, "localhost:5000",
        []grpc.DialOption{grpc.WithInsecure()},
    )
    // in ra lỗi nếu có
    if err != nil {
        log.Fatal(err)
    }
    // bắt đầu lắng nghe http client trên port 8080
    http.ListenAndServe(":8080", mux)
}

// $ go run proxy/main.go
```

Tiếp theo ta sẽ chạy gRPC service:

***restservice/main.go:***

```go
// khai báo struct hiện thực RestService
type RestServiceImpl struct{}
// hàm Get RPC được hiện thực như sau
func (r *RestServiceImpl) Get(ctx context.Context, message *StringMessage) (*StringMessage, error) {
    return &StringMessage{Value: "Get hi:" + message.Value + "#"}, nil
}
// tương tự với hàm Post RPC được hiện thực với
func (r *RestServiceImpl) Post(ctx context.Context, message *StringMessage) (*StringMessage, error) {
    return &StringMessage{Value: "Post hi:" + message.Value + "@"}, nil
}
// hàm main của gRPC service
func main() {
    // khởi tạo một grpc Server mới
    grpcServer := grpc.NewServer()
    // register grpc Server với đối tượng hiện thực các hàm RPC
    RegisterRestServiceServer(grpcServer, new(RestServiceImpl))
    // listen gRPC Service trên port 5000, bỏ qua lỗi trả về nếu có
    lis, _ := net.Listen("tcp", ":5000")
    grpcServer.Serve(lis)
}

// $ go run restservice/main.go
```

Sau khi chạy hai chương trình gRPC và REST services, chúng ta có thể tạo request REST service với lệnh [curl](https://thoainguyen.github.io/2019-06-15-using-curl/):

```sh
// gọi service Get
$ curl localhost:8080/get/gopher
{"value":"Get: gopher"}
// gọi service Post
$ curl localhost:8080/post -X POST --data '{"value":"grpc"}'
{"value":"Post: grpc"}
```

Khi chúng ta publishing REST interface thông qua [Swagger](http://swagger.io), một swagger file có thể được sinh ra nhờ vào công cụ grpc-gateway bằng lệnh bên dưới:

```sh
// chạy lệnh sau để cài đặt nếu chưa có sẵn
$ go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
// lệnh sinh ra swagger file
$ protoc -I. \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=. \
  hello.proto

// Trong đó,
// - --swagger_out=.: dùng plugin swagger để sinh ra swagger file tại thư mục hiện tại
```

File `hello.swagger.json` sẽ được sinh ra sau đó. Trong trường hợp này, chúng ta có thể dùng `swagger-ui project` để cung cấp tài liệu `REST interface` và testing dưới dạng web pages.

## 3.6.3 Dùng Docker grpc-gateway
Với những lập trình viên phát triển gRPC Services trên các ngôn ngữ không phải Golang như Java, C++, ... có nhu cầu sinh ra grpc gateway cho các services của họ nhưng gặp khá nhiều khó khăn từ việc cài đặt môi trường Golang, protobuf, các lệnh generate,v,v.. Có một giải pháp đơn giản hơn đó là sử dụng Docker để xây dựng grpc-gateway theo bài hướng dẫn chi tiết sau [buildingdocker-grpc-gateway](https://medium.com/zalopay-engineering/buildingdocker-grpc-gateway-e2efbdcfe5c).

## 3.6.4 Nginx
Những phiên bản [Nginx](https://www.nginx.com/) về sau cũng đã hỗ trợ `gRPC` với khả năng register nhiều gRPC service instance giúp load balancing (cân bằng tải) dễ dàng hơn. Những extension của Nginx về gRPC là một chủ đề lớn, ở đây chúng tôi không trinhf bày hết được, các bạn có thể tham khảo các tài liệu trên trang chủ của Nginx như [ở đây](https://www.nginx.com/blog/nginx-1-13-10-grpc/).