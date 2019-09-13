# 3.2. Protobuf

[Protobuf](https://developers.google.com/protocol-buffers/) hay Protocols Buffer là một ngôn ngữ dùng để mô tả các cấu trúc dữ liệu, chúng ta dùng protoc để biên dịch chúng thành mã nguồn của các ngôn ngữ lập trình khác nhau có chức năng serialize và deserialize các cấu trúc dữ liệu này thành dạng binary stream. So với dạng XML hoặc JSON thì dữ liệu đó nhỏ gọn gấp 3-10 lần và được xử lý rất nhanh.

<div align="center">
	<img src="../images/ch4-2-size.png" width="580">
	<br/>
    <img src="../images/ch4-2-speed.png" width="580">
    <br/>
</div>

*Xem thêm: [Benchmarking Protocol Buffers, JSON and XML in Go](https://medium.com/@shijuvar/benchmarking-protocol-buffers-json-and-xml-in-go-57fa89b8525)*.

Bạn đọc có thể cài đặt và làm quen với các ví dụ Protobuf trên [trang chủ](https://developers.google.com/protocol-buffers/docs/gotutorial) trước khi đi vào nội dung chính.

## 3.2.1 Kết hợp Protobuf với RPC

Đầu tiên chúng ta tạo file `hello.proto` chứa kiểu String được dùng trong RPC HelloService.

***hello.proto:***

```go
// phiên bản proto3
syntax = "proto3";
// tên package được sinh ra
package main;
// message là một đơn vị dữ liệu trong Protobuf
message String {
    // chuỗi string được truyền vào hàm RPC
    string value = 1;
}
```

Để sinh ra mã nguồn Go từ file `hello.proto` ở trên, đầu tiên là cài đặt bộ biên dịch `protoc` qua liên kết [ở đây](https://github.com/google/protobuf/releases), sau đó là cài đặt một plugin cho Go thông qua lệnh:

```sh
$ go get github.com/golang/protobuf/protoc-gen-go
```

Chúng ta sẽ sinh ra mã nguồn Go bằng lệnh sau:

```sh
$ protoc --go_out=. hello.proto
// Trong đó,
// protoc: chương trình sinh mã nguồn
// go_out: chỉ cho protoc tải plugin protoc-gen-go, (cũng có java_out, python_out,..)
// --go_out=.: sinh ra mã nguồn tại thư mục hiện tại
// hello.proto: file Protobuf
```

Sẽ có một file `hello.pb.go` được sinh ra, trong đó cấu trúc String được định nghĩa là:

***hello.pb.go:***

```go
type String struct {
    Value string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
    //...
}

func (m *String) Reset()         { *m = String{} }
func (m *String) String() string { return proto.CompactTextString(m) }
func (*String) ProtoMessage()    {}
func (*String) Descriptor() ([]byte, []int) {
    return fileDescriptor_hello_069698f99dd8f029, []int{0}
}
//...
func (m *String) GetValue() string {
    if m != nil {
        return m.Value
    }
    return ""
}
//...
```

Ở phần [3.1](./ch3-01-rpc-go.md) chúng ta đã xây dựng một RPC HelloService đơn giản dựa trên thư viện chuẩn [net/rpc](https://godoc.org/net/rpc) có kiểu dữ liệu request, reply do người dùng tự định nghĩa, bây giờ dựa trên kiểu String mới được sinh ra từ Protobuf, chúng ta có thể viết lại RPC HelloService như sau:

***hello.go:***

```go
// RPC struct
type HelloService struct{}
// định nghĩa hàm Hello RPC, với tham số là kiểu String vừa định nghĩa trong Protobuf
func (p *HelloService) Hello(request *String, reply *String) error {
    // các hàm như .GetValue() đã được tạo ra trong file hello.pb.go
    reply.Value = "Hello, " + request.GetValue()
    // trả về nil khi thành công
    return nil
}
```

Chúng ta vẫn phải tự xây dựng hàm **Hello(request, reply)** bằng cách tự viết. Khi sử dụng Protobuf chúng ta có thể tự định nghĩa luôn service mình có những hàm rpc nào, nhận vào request và trả về reply ra sao. Chúng ta định nghĩa HelloService trong file proto như sau:

***hello.proto***

```go
// ...
// định nghĩa service
service HelloService {
    // định nghĩa lời gọi hàm RPC
    rpc Hello (String) returns (String);
}
```

Chúng ta cần có một plugin để sinh ra mã nguồn service tương ứng với định nghĩa ở trên. Hiện nay Google đã phát triển bộ [gRPC plugin](https://github.com/golang/protobuf/blob/master/protoc-gen-go/grpc/grpc.go) giúp tạo ra mã nguồn tương ứng với file proto. Ở phần dưới sẽ trình bày cách xây dựng một plugin dựa trên mã nguồn gRPC plugin, chi tiết về gRPC chúng tôi sẽ đề cập ở các phần sau.

## 3.2.2 Viết plugin sinh mã nguồn RPC service

Từ mã nguồn [gRPC plugin](https://github.com/golang/protobuf/blob/master/protoc-gen-go/grpc/grpc.go), chúng ta có thể thấy hàm `generator.RegisterPlugin` được dùng để đăng kí `plugin` đó, Interface của một plugin sẽ như sau:

```go
type Plugin interface {
    // Name() trả về tên của plugin.
    Name() string
    // Init() được gọi sau khi data structures built xong
    // và trước khi quá trình sinh code bắt đầu.
    Init(g *Generator)
    // Generate() là hàm sinh ra mã nguồn vào file
    Generate(file *FileDescriptor)
    // Hàm này được gọi sau khi Generate().
    GenerateImports(file *FileDescriptor)
}
```

Do đó, chúng ta có thể  xây dựng một plugin mang tên `netrpcPlugin` để sinh ra mã nguồn RPC service cho Go từ file Protobuf.

```go
import (
    // import gói thư viện để sinh ra plugin
    "github.com/golang/protobuf/protoc-gen-go/generator"
)
// định nghĩa struct netrpcPlugin xây dựng interface Plugin
type netrpcPlugin struct{ *generator.Generator }
// định nghĩa Name() function
func (p *netrpcPlugin) Name() string                { return "netrpc" }
// định nghĩa Init() function
func (p *netrpcPlugin) Init(g *generator.Generator) { p.Generator = g }
// định nghĩa GenerateImports()
func (p *netrpcPlugin) GenerateImports(file *generator.FileDescriptor) {
    if len(file.Service) > 0 {
        p.genImportCode(file)
    }
}
// định nghĩa Generate()
func (p *netrpcPlugin) Generate(file *generator.FileDescriptor) {
    for _, svc := range file.Service {
        p.genServiceCode(svc)
    }
}
```

Hiện tại, phương thức `genImportCode` và `genServiceCode` tạm thời như sau:

```go
func (p *netrpcPlugin) genImportCode(file *generator.FileDescriptor) {
    p.P("// TODO: import code")
}

func (p *netrpcPlugin) genServiceCode(svc *descriptor.ServiceDescriptorProto) {
    p.P("// TODO: service code, Name = " + svc.GetName())
}
```

Để sử dụng plugin, chúng ta cần phải đăng kí plugin đó với hàm `generator.RegisterPlugin`, chúng có thể được xây dựng chúng nhờ vào hàm `init()`.

```go
func init() {
    generator.RegisterPlugin(new(netrpcPlugin))
}
```

Bởi vì trong ngôn ngữ Go, package chỉ được import tĩnh, chúng ta không thể thêm plugin mới vào plugin đã có sẵn là `protoc-gen-go`. Chúng ta sẽ `re-clone` lại hàm main để build lại `protoc-gen-go`.

```go
package main

import (
    "io/ioutil"
    "os"
    // import các package cần thiết
    "github.com/golang/protobuf/proto"
    "github.com/golang/protobuf/protoc-gen-go/generator"
)
// bắt đầu hàm main
func main() {
    // sinh ra một đối tượng plugin mới
    g := generator.New()
    // đọc lệnh từ console vào biến data
    data, err := ioutil.ReadAll(os.Stdin)
    // in ra lỗi nếu có
    if err != nil {
        g.Error(err, "reading input")
    }
    // unmarsal data thành cấu trúc Request
    if err := proto.Unmarshal(data, g.Request); 
    // in ra lỗi nếu có
    err != nil {
        g.Error(err, "parsing input proto")
    }
    // kiểm tra xem tên file có hợp lệ không
    if len(g.Request.FileToGenerate) == 0 {
        g.Fail("no files to generate")
    }
    // đăng ký các tham số
    g.CommandLineParameters(g.Request.GetParameter())
    g.WrapTypes()
    // thiết lập tên package
    g.SetPackageNames()
    g.BuildTypeNameMap()

    // sinh ra các file mã nguồn
    g.GenerateAllFiles()

    // Trả về kết quả
    data, err = proto.Marshal(g.Response)
    if err != nil {
        g.Error(err, "failed to marshal output proto")
    }
    // ghi kết quả ra màn hình
    _, err = os.Stdout.Write(data)
    // in ra lỗi nếu có
    if err != nil {
        g.Error(err, "failed to write output proto")
    }
}
```

Để tránh việc trùng tên với protoc-gen-go plugin, chúng ta sẽ đặt tên cho plugin mới là `protoc-gen-go-netrpc`, và dự định biên dịch lại `hello.proto` với lệnh sau:

```sh
$ protoc --go-netrpc_out=plugins=netrpc:. hello.proto
```

Tham số `--go-netrpc_out` sẽ nói cho bộ biên dịch protoc biết là nó phải tải một plugin với tên gọi là protoc-gen-go-netrpc. Bây giờ, tiếp tục phát triển netrpcPlugin plugin với mục tiêu cuối cùng là sinh ra lớp Interface RPC. Đầu tiên chúng ta sẽ phải xây dựng genImportCode:

```go
func (p *netrpcPlugin) genImportCode(file *generator.FileDescriptor) {
    p.P(`import "net/rpc"`)
}
```

Chúng ta sẽ định nghĩa kiểu ServiceSpec được mô tả như là thông tin thêm vào của service.

```go
type ServiceSpec struct {
    // Tên của service
    ServiceName string
    // Danh sách cách Service method
    MethodList  []ServiceMethodSpec
}

type ServiceMethodSpec struct {
    MethodName     string
    InputTypeName  string
    OutputTypeName string
}
```

Chúng ta sẽ tạo ra một phương thức `buildServiceSpec`, nó sẽ parse thông tin thêm vào service được định nghĩa trong ServiceSpec cho mỗi service.

```go
// phương thức buildServiceSpec
func (p *netrpcPlugin) buildServiceSpec(
    // tham số truyền vào thuộc kiểu ServiceDescriptorProto 
    // mô tả thông tin về service
    svc *descriptor.ServiceDescriptorProto,
) *ServiceSpec {
    // khởi tạo đối tượng
    spec := &ServiceSpec{
        // svc.GetName(): lấy tên service được định nghĩa ở Protobuf file
        // sau đó chuyển đổi chúng về style CamelCase
        ServiceName: generator.CamelCase(svc.GetName()),
    }
    // mới mỗi phương thức RPC, ta thêm một cấu trúc tương ứng vào danh sách
    for _, m := range svc.Method {
        spec.MethodList = append(spec.MethodList, ServiceMethodSpec{
            // m.GetName(): lấy tên phương thức
            MethodName:     generator.CamelCase(m.GetName()),
            // m.GetInputType(): lấy kiểu dữ liệu tham số đầu vào
            InputTypeName:  p.TypeName(p.ObjectNamed(m.GetInputType())),
            OutputTypeName: p.TypeName(p.ObjectNamed(m.GetOutputType())),
        })
    }
    // trả về cấu trúc trên
    return spec
}
```

Sau đó chúng ta sẽ sinh ra mã nguồn của service dựa trên thông tin mô tả đó, được xây dựng bởi phương thức `buildServiceSpec` :

```go
func (p *netrpcPlugin) genServiceCode(svc *descriptor.ServiceDescriptorProto) {
    // gọi hàm được định nghĩa ở trên
    spec := p.buildServiceSpec(svc)
    // buf là biến chứa dữ liệu
    var buf bytes.Buffer
    // dùng tmplService cho việc sinh mã nguồn
    t := template.Must(template.New("").Parse(tmplService))
    // thực thi việc sinh mã nguồn
    err := t.Execute(&buf, spec)
    // in ra lỗi nếu có
    if err != nil {
        log.Fatal(err)
    }
    // ghi buf.String() vào file
    p.P(buf.String())
}
```

Chúng ta mong đợi vào mã nguồn cuối cùng được sinh ra như sau:

```go
type HelloServiceInterface interface {
    Hello(in String, out *String) error
}

func RegisterHelloService(srv *rpc.Server, x HelloService) error {
    if err := srv.RegisterName("HelloService", x); err != nil {
        return err
    }
    return nil
}

type HelloServiceClient struct {
    *rpc.Client
}

var _ HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (*HelloServiceClient, error) {
    c, err := rpc.Dial(network, address)
    if err != nil {
        return nil, err
    }
    return &HelloServiceClient{Client: c}, nil
}

func (p *HelloServiceClient) Hello(in String, out *String) error {
    return p.Client.Call("HelloService.Hello", in, out)
}
```

Để làm được như vậy, template của chúng ta được viết như sau:

```go
const tmplService = `
{{$root := .}}

type {{.ServiceName}}Interface interface {
    {{- range $_, $m := .MethodList}}
    {{$m.MethodName}}(*{{$m.InputTypeName}}, *{{$m.OutputTypeName}}) error
    {{- end}}
}

func Register{{.ServiceName}}(
    srv *rpc.Server, x {{.ServiceName}}Interface,
) error {
    if err := srv.RegisterName("{{.ServiceName}}", x); err != nil {
        return err
    }
    return nil
}

type {{.ServiceName}}Client struct {
    *rpc.Client
}

var _ {{.ServiceName}}Interface = (*{{.ServiceName}}Client)(nil)

func Dial{{.ServiceName}}(network, address string) (
    *{{.ServiceName}}Client, error,
) {
    c, err := rpc.Dial(network, address)
    if err != nil {
        return nil, err
    }
    return &{{.ServiceName}}Client{Client: c}, nil
}

{{range $_, $m := .MethodList}}
func (p *{{$root.ServiceName}}Client) {{$m.MethodName}}(
    in *{{$m.InputTypeName}}, out *{{$m.OutputTypeName}},
) error {
    return p.Client.Call("{{$root.ServiceName}}.{{$m.MethodName}}", in, out)
}
{{end}}`
```

Khi plugin mới của protoc được hoàn thành, mã nguồn có thể được sinh ra mỗi khi RPC service thay đổi trong `hello.proto` file. Chúng ta có thể điều chỉnh hoặc thêm nội dung của mã nguồn được sinh ra bằng việc cập nhật template plugin.

<div style="display: flex; justify-content: space-around;">
<span> <a href="ch3-01-rpc-go.md">&lt Phần 3.1</a>
</span>
<span><a href="../SUMMARY.md"> Mục lục</a>  </span>
<span> <a href="ch3-03-grpc.md">Phần 3.3 &gt</a> </span>
</div>
