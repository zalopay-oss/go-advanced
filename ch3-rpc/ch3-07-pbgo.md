# 3.7 Framework dựa trên Protobuf: pbgo

[Pbgo](https://github.com/chai2010/pbgo) là một framework nhỏ gọn dựa trên cú pháp mở rộng của Protobuf để sinh ra mã nguồn `REST` cho RPC service, trong phần này, chúng ta sẽ cùng tìm hiểu Pbgo.

## 3.7.1 Cú pháp mở rộng của Protobuf

Cú pháp mở rộng của Protobuf được dùng trong rất nhiều dự án opensource xung quanh nó. Ở phần trước, chúng ta đã đề cập [validator](https://github.com/mwitkow/go-proto-validators), một plugin dùng để validate các trường theo các rules được định nghĩa trong phần mở rộng của trường tương ứng.

Trong dự án [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway), việc hỗ trợ REST interface đạt được bằng cách thêm thông tin HTTP vào phần mở rộng cho mỗi hàm RPC của service. Tương tự, các phần cú pháp mở rộng của Pbgo được định nghĩa:

***[pbgo/pbgo.proto](https://github.com/chai2010/pbgo/blob/master/pbgo.proto):***

```go
// các .proto file khác phải import file này khi sử dụng pbgo framework
syntax = "proto3";
package pbgo;
// định nghĩa package được sinh ra
option go_package = "github.com/chai2010/pbgo;pbgo";
// import cấu trúc mô tả của Protobuf
import "google/protobuf/descriptor.proto";

// định nghĩa một phần mở rộng có tên rest_api
// với cấu trúc HttpRule
extend google.protobuf.MethodOptions {
    HttpRule rest_api = 20180715;
}
// các phương thức Http được định nghĩa trong HttpRule
message HttpRule {
    string get = 1;
    string put = 2;
    string post = 3;
    string delete = 4;
    string patch = 5;
}
// xem thêm tại: https://github.com/chai2010/pbgo/blob/master/pbgo.proto
```

Sau khi extension đã được định nghĩa, chúng ta có thể import nó vào những file Protobuf khác, ví dụ là file hello.proto:

***hello.proto :***

```go
syntax = "proto3";
package hello_pb;
// import file pbgo.proto được định nghĩa ở trên
import "github.com/chai2010/pbgo/pbgo.proto";
// định nghĩa message truyền nhận
message String {
    string value = 1;
}
// định nghĩa HelloService
service HelloService {
    rpc Hello (String) returns (String) {
        // cú pháp mở rộng giống với grpc-gateway
        option (pbgo.rest_api) = {
            // get là tên phương thức HTTP
            // :value là giá trị trong String message
            // "/hello/:value" là uri trong httprouter
            get: "/hello/:value"
        };
    }
}
```

## 3.7.2. Đọc thông tin mở rộng của plugin

Phần trước, chúng ta đã định nghĩa plugin trong Protobuf, bây giờ để sinh ra mã nguồn cho RPC từ plugin. Đầu tiên, định nghĩa interface:

***Interface generator.Plugin :***

```go
type Plugin interface {
    // Name() trả về tên của plugin.
    Name() string
    // Init() được gọi sau khi data structures built 
    // xong và trước khi quá trình sinh code bắt đầu.
    Init(g *Generator)
    // Generate() là hàm sinh ra mã nguồn vào file
    Generate(file *FileDescriptor)
    // Hàm này được gọi sau khi Generate().
    GenerateImports(file *FileDescriptor)
}
```

***Hiện thực hàm Generate() :***

```go
// pbgoPlugin là đối tượng chính của framework pbgo
func (p *pbgoPlugin) Generate(file *generator.FileDescriptor) {
    // duyệt qua tất cả các service được định nghĩa trong file .proto
    for _, svc := range file.Service {
        // duyệt qua tất cả các hàm trong mỗi service
        for _, m := range svc.Method {
            // lấy cấu trúc httpRule được định nghĩa trong phần mở rộng
            // phương thức getServiceMethodOption được custom sẽ nói sau.
            httpRule := p.getServiceMethodOption(m)
            ...
        }
    }
}
```

Trước khi chúng ta nói về phương thức getServiceMethodOption(), định nghĩa phần extension cho phương thức.

***Extension :***

```go
extend google.protobuf.MethodOptions {
    // rest_api là tên extension
    HttpRule rest_api = 20180715;
    // từ rest_api sẽ sinh ra `pbgo.E_RestApi` được dùng để lưu 
    // thông tin mở rộng do người dùng định nghĩa
}
```



Bên dưới là phần hiện thực phương thức getServiceMethodOption().

***getServiceMethodOption() :***

```go
func (p *pbgoPlugin) getServiceMethodOption(
    m *descriptor.MethodDescriptorProto,
) *pbgo.HttpRule {
    if m.Options != nil && proto.HasExtension(m.Options, pbgo.E_RestApi) {
        // lấy thông tin mở rộng qua hàm GetExtension()
        ext, _ := proto.GetExtension(m.Options, pbgo.E_RestApi)
        if ext != nil {
            if x, _ := ext.(*pbgo.HttpRule); x != nil {
                return x
            }
        }
    }
    return nil
}
```

Với thông tin về extension trên, chúng ta có thể sinh ra mã nguồn REST bằng việc tham khảo đến cách mà mã nguồn RPC được sinh ra ở phần hai.

## 3.7.3 Sinh ra REST code

Framework **pbgo** cũng hỗ trợ một số plugin cho việc sinh ra mã nguồn REST. Tuy nhiên, mục tiêu của chúng ta là học được quy trình thiết kế framework pbgo, do đó đầu tiên chúng ta phải viết mã nguồn REST ứng với phương thức Hello, và sau đó phần mã nguồn được plugin tự động được sinh ra dựa trên một template được định nghĩa sẵn.

HelloService chỉ có một phương thức là Hello, phương thức Hello chỉ định nghĩa một REST interface.

***hello.proto:***

```go
message String {
    string value = 1;
}

service HelloService {
    rpc Hello (String) returns (String) {
        option (pbgo.rest_api) = {
            get: "/hello/:value"
        };
    }
}
```

Để người dùng cuối dễ dàng sử dụng, chúng ta cần xây dựng một route đến `HelloService`. Do đó, chúng ta sẽ có một hàm giống như `HelloServiceHandler` để sinh ra mã nguồn route handler dựa trên interface của service `HelloServiceInterface`.

***Mã nguồn Route handler (1):***

```go
type HelloServiceInterface interface {
    Hello(in *String, out *String) error
}

func HelloServiceHandler(svc HelloServiceInterface) http.Handler {
    var router = httprouter.New()
    _handle_HelloService_Hello_get(router, svc)
    return router
}
```

Mã nguồn chọn một opensource [httprouter](https://github.com/julienschmidt/httprouter) nổi tiếng để hiện thực. Hàm `_handle_HelloService_Hello_get` được dùng để register hàm `Hello` cho `route handler`.

***Mã nguồn Route handler (2):***

```go
func _handle_HelloService_Hello_get(
    router *httprouter.Router, svc HelloServiceInterface,
) {
    router.Handle("GET", "/hello/:value",
        func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
            var protoReq, protoReply String

            err := pbgo.PopulateFieldFromPath(&protoReq, fieldPath, 
            // ps.ByName("value") sẽ load giá trị parameter URL
            ps.ByName("value"))
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            // gọi hàm RPC Hello lưu giá trị vào protoReply
            if err := svc.Hello(&protoReq, &protoReply); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            // trả về protoReply cho user theo kiểu Json
            if err := json.NewEncoder(w).Encode(&protoReply); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
        },
    )
}
```

Sau khi thiết lập cấu trúc mã nguồn, bạn có thể xây dựng một template cho việc sinh ra mã nguồn plugin cơ bản. Toàn bộ plugin code và template nằm trong file [protoc-gen-pbgo/pbgo.go](https://github.com/chai2010/pbgo/blob/master/protoc-gen-pbgo/pbgo/pbgo.go).

## 3.7.4 Sử dụng Pbgo

Mặc dù quá trình để xây dựng một `pbgo` framework từ ban đầu hơi phức tạp, việc sử dụng `pbgo` để xây dựng một REST service lại cực kì đơn giản.

Đầu tiên định nghĩa file hello.proto:

***proto/hello.proto:***

```go
syntax = "proto3";
package hello_pb;

import "github.com/chai2010/pbgo/pbgo.proto";

message String {
    string value = 1;
}
service HelloService {
    rpc Hello (String) returns (String) {
        option (pbgo.rest_api) = {
            get: "/hello/:value"
        };
    }
}
```

Sinh ra mã nguồn ***hello.pb.go*** bằng lệnh:

```sh
$ protoc -I=. -I=$GOPATH/src --pbgo_out=. proto/hello.proto
```

Định nghĩa RPC Server của Hello Service:

***hello/hello.go:***

```go
package main

import (
	"log"
	"net/http"

	hello_pb "../proto"
)

type HelloService struct{}

// định nghĩa hàm Hello RPC bên phía server
func (p *HelloService) Hello(request *hello_pb.String, reply *hello_pb.String) error {
	reply.Value = "hello:" + request.GetValue()
	return nil
}

// hàm main để register HelloService và lắng nghe yêu cầu trên port 8080
func main() {
	router := hello_pb.HelloServiceHandler(new(HelloService))
	log.Fatal(http.ListenAndServe(":8080", router))
}

```

Sau đó chạy REST Service bằng lệnh:

```sh
$ go run hello/hello.go
```

Kiểm tra service với lệnh:

```sh
$ curl localhost:8080/hello/vietnam
{"value":"hello:vietnam"}
```

Bạn đọc có thể xem thêm các ví dụ tại [đây](https://github.com/chai2010/pbgo/blob/master/README.md).
