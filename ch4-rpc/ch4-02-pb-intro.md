# 4.2 Protobuf

Protobuf là tên gọi rút gọn của Protocols Buffers, một ngôn ngữ mô tả được phát triển bởi Google và được công bố vào năm 2008. Có thể hình dung Protobuf cũng như XML, JSON hay những ngôn ngữ mô tả khác. Chúng ta có thể sinh ra code từ nó và cung cấp cơ chế để serialization cấu trúc dữ liệu thông qua một số công cụ đi kèm. Nhưng chúng ta thường quan tâm nhiều hơn về Protobuf như là một ngôn ngữ mô tả dữ liệu để đặc tả lớp giao diện, nó có thể được xem như là một phương tiện cơ bản cho việc thiết kế những lớp giao diện RPC được bảo mật và an toàn hơn.

## 4.2.1 Bắt đầu với Protobuf

Cho những ai chưa từng làm quen với protobuf, chúng tôi khuyên hãy hiểu cách sử dụng cơ bản chúng từ trang chủ [protobuf](https://developers.google.com/protocol-buffers/). Ở đây chúng ta thử kết hợp protobuf và RPC, và cuối cùng sẽ đảm bảo rằng giao diện được đặc tả và tính bảo mật của lệnh gọi RPC thông qua protobuf. Đơn vị cơ bản của dữ liệu Protobuf là một message, chúng cũng tương tự như cấu trúc của ngôn ngữ Go. Các thành phần của message hoặc những message từ những kiểu dữ liêu bên dưới khác có thể được lồng vào mesage đó.

Đầu tiên chúng ta tạo file `hello.proto` nó sẽ chứa kiểu string được dùng cho HelloService service:

```go
syntax = "proto3";

package main;

message String {
    string value = 1;
}
```

Cú pháp của statement trên bắt đầu bằng việc định nghĩa trường "syntax" là "proto3" - Phiên bản ngôn ngữ protobuf thứ ba, và tất cả những thành phần được khởi tạo với giá trị 0 giống như Go (không có hỗ trợ custom), sau đó những thành phần của message sẽ không cần những thuộc tính được yêu cầu. Sau đó package được chỉ thị là "main package" (nó nên đồng nhất với tên package trong Go, đơn giản như code ví dụ), dĩ nhiên, user cũng có thể tùy chỉnh đường dẫn package tương ứng khác cho những ngôn ngữ khác nhau. Cuối cùng, từ khóa "message" sẽ định nghĩa một kiểu dữ liệu string mới là String, nó sẽ ứng với cấu trúc string ở mã nguồn cuối cùng được sinh ra từ chúng trong Go, và những thành phần sẽ được định nghĩa với một số tên gọi.

Trong ngôn ngữ mô tả như là XML hay JSON, kiểu dữ liệu tương ứng thông thường sẽ được bao bọc bởi tên các thành viên. Tuy nhiên, Protobuf encoding sẽ kết hợp những dữ liệu bằng một số duy nhất của dữ liệu đó, do đó dung lượng của protobuf sẽ dữ liệu protobuf được encoded sẽ nhỏ, nhưng nó không dễ dàng để con người có thể đọc được. Chúng ta sẽ không quan tâm đến công nghệ encode dữ liệu của protobuf. Két quả của cấu trúc Go có thể được encode bằng JSON hay không, do đó chúng ta có thể tạm thời phớt lờ đi việc Protobuf encode dữ liệu trong tài liệu này.

Phần lõi của Protobuf được phát triển dựa trên ngôn ngữ C++, và chúng không dùng ngôn ngữ Go trong bộ biên dịch protoc. Để sinh ra mã nguồn Go tương ứng với file hello.go ở trên, chúng ta sẽ phải cần cài đặt một số plugin khác. Đầu tiên là cài đặt bộ biên dịch protoc, chúng có thể được tải về tại https://github.com/google/protobuf/releases. Sau đó là cài đặt một plugin cho Go, chúng ta có thể cài đặt thông qua `go get github.com/golang/protobuf/protoc-gen-go`.

Sau đó chúng ta sẽ sinh ra mã nguồn Go bằng lệnh sau:

```
$ protoc --go_out=. hello.proto
```

Tham số `go_out` chỉ ra cho protoc sẽ tải công cụ `protoc-gen-go`, sau đó sinh ra code thông qua công cụ đó và sinh code trong cùng một đường dẫn với file `hello.proto`. Cuối cùng, sau đây là danh sách chuỗi các file protobuf cho tiến trình đó.

Chỉ có một file `hello.pb.go` được sinh ra, cấu trúc String sẽ như sau

```go
type String struct {
    Value string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
}

func (m *String) Reset()         { *m = String{} }
func (m *String) String() string { return proto.CompactTextString(m) }
func (*String) ProtoMessage()    {}
func (*String) Descriptor() ([]byte, []int) {
    return fileDescriptor_hello_069698f99dd8f029, []int{0}
}

func (m *String) GetValue() string {
    if m != nil {
        return m.Value
    }
    return ""
}
```

Cấu trúc được sinh ra của chứa một số hàm với tiền tố `XXX_`, chúng ta có thể ẩn đi những thành phần đó. Cùng một thời điểm, kiểu String cũng có thể được tự động sinh ra một tập hợp các phương thức, trong số đó ProtoMessage chỉ ra rằng đó là một hàm được hiện thực giao diện proto.Message. Thêm vào đó, Protobuf sẽ sinh ra những phương thức Get cho mỗi thành phần, trong nó sẽ kiểm tra dữ liệu null và trả về chuỗi rỗng.

Dựa trên kiểu String mới, chúng ta có thể hiện thực lại HelloService service

```go
type HelloService struct{}

func (p *HelloService) Hello(request *String, reply *String) error {
    reply.Value = "hello:" + request.GetValue()
    return nil
}
```

Tham số đầu vào và tham số đầu ra của phương thức Hello được thể hiện bởi kiểu String được định nghĩa bởi protobuf. Bởi vì tham số đầu vào mới là một kiểu cấu trúc, kiểu dữ liệu con trỏ được sử dụng như là tham số đầu vào, và một mã nguồn bên trong của hàm sẽ cũng được hiệu chỉnh cho phù hợp.

Chúng ta đầu tiên sẽ nhận ra sự kết hợp giữa Protobuf và RPC. Khi chúng ta bắt đầu một RPC service, chúng ta có thể vẫn chọn một kiểu mặc định hoặc định nghĩa lại với Json, và sau đó sẽ hiện thực lại plugin dựa trên mã nguồn protobuf. Mặc dù chúng ta có thể làm rất nhiều công việc, nó dường như chúng ta không thể đạt được gì đáng kể.

Nhìn lại giao diện RPC khá bảo mật của chương 1, chúng ta đã rất nỗ lực để đảm bảo bảo mật cho dịch vụ RPC. Kết quả của mã nguồn RPC trên sẽ an toàn hơn và rất là tuyệt vời để  bảo trì thủ công, chúng ta có thể bảo mật hóa những mã nguồn liên quan mà nó chỉ sẵn có ở môi trường ngôn ngữ Go. Do đó đầu vào và đầu ra của tham số được định nghĩa bởi Protobuf được dùng, có thể giao diện RPC được định nghĩa bởi protobuf. Việc áp dụng protobuf được định nghĩa ở mức độc lập ngôn ngữ dịch vụ RPC và giao diện của chúng ở giá trị thực tế.

Cập nhật hello.proto file bên dưới dễ định nghĩa dịch vụ RPC HelloService service thông qua protobuf.


```go
service HelloService {
    rpc Hello (String) returns (String);
}
```

Nhưng khi sinh lại mã nguồn Go, chúng cũng không thay đổi. Đó là bởi vì có hàng triệu hiện thực RPC trên thế giới, và bộ biên dịch protoc sẽ không thể biết sinh ra mã nguồn của HelloService như thế nào.

Tuy nhiên `grpc`, một plugin khác đã được tích hợp bên trong `protoc-gen-go` để sinh ra mã nguồn cho gRPC.

```
$ protoc --go_out=plugins=grpc:. hello.proto
```

Trong mã nguồn được sinh ra, sẽ có một số kiểu mới như là HelloServiceServer và HelloServiceClient. Những kiểu đó được dùng bởi gRPC và không cần trong phần hiện thực RPC của chúng ta.

Tuy nhiên gRPC plugin sẽ cung cấp cho chúng ta những ý tưởng mới. Bên dưới chúng ta sẽ khám phá ra làm thế nào để sinh ra mã nguồn bảo mật trong RPC.

## 4.2.2 Tùy chỉnh mã nguồn được sinh ra bởi plugin.

Bộ biên dịch protoc của Protobuf sẽ hiện thực để hỗ trợ những ngôn ngữ khác nhau thông qua cơ chế plugin. cho ví dụ, nếu lệnh protoc có tham số  được định dạng `--xxx_out`, thì sau đó proto sẽ yêu cầu plugin được xây dựng dựa trên ngôn ngữ `xxx`. Sau đó plugin sẽ sinh ra mã nguồn, ví dụ `protoc-gen-go` sẽ sinh ra mã nguồn Go bằng tham số `--go_out=plugins=grpc` sẽ sinh ra mã nguồn cho gRPC, nếu không chúng sẽ chỉ sinh ra mã nguồn liên quan đến message đó.

Đề cập đến mã nguồn từ gRPC plugin, chúng ta có thể thấy rằng hàm `generator.RegisterPlugin` sẽ có thể được dùng để đăng kí `plugin` đó. Dưới đây sẽ sinh ra giao diện plugin.

```go
// A Plugin provides functionality to add to the output during
// Go code generation, such as to produce RPC stubs.
type Plugin interface {
    // Name identifies the plugin.
    Name() string
    // Init is called once after data structures are built but before
    // code generation begins.
    Init(g *Generator)
    // Generate produces the code generated by the plugin for this file,
    // except for the imports, by calling the generator's methods P, In,
    // and Out.
    Generate(file *FileDescriptor)
    // GenerateImports produces the import declarations for this file.
    // It is called after Generate.
    GenerateImports(file *FileDescriptor)
}
```

Phương thức Name sẽ trả về tên của plugin. Đó là một plugin ở hệ thống cho việc hiện thực Protobuf của ngôn ngữ Go. Sẽ không có gì để làm với tên của protoc plugin. Sau đó hàm Init sẽ khởi tạo plugin với tham số `g`, nó sẽ chứa toàn bộ thông tin về Proto file. Cuối cùng phương thức Generate và GenerateImports sẽ được dùng để sinh ra phần thân của mã nguồn tương ứng với package được import.

Do đó chúng ta có thể hiện thực lại hàm `netrpcPlugin` để sinh ra mã nguồn cho thư viện RPC chuẩn của Go.


```go
import (
    "github.com/golang/protobuf/protoc-gen-go/generator"
)

type netrpcPlugin struct{ *generator.Generator }

func (p *netrpcPlugin) Name() string                { return "netrpc" }
func (p *netrpcPlugin) Init(g *generator.Generator) { p.Generator = g }

func (p *netrpcPlugin) GenerateImports(file *generator.FileDescriptor) {
    if len(file.Service) > 0 {
        p.genImportCode(file)
    }
}

func (p *netrpcPlugin) Generate(file *generator.FileDescriptor) {
    for _, svc := range file.Service {
        p.genServiceCode(svc)
    }
}
```

Đầu tiên phương thức Name sẽ trả về tên của plugin. `netrpcPlugin` có một hàm dựng sẵn `*generator.Generator`, và được khởi tạo với tham số `g` khi hàm Init được khởi tạo, do đó, plugin sẽ kế thừa tất cả những phương thức được public  từ tham số g này. Phương thức ` GenerateImports` sẽ gọi hàm `genImportCode` được định nghĩa để sinh ra mã nguồn import. Và Phương thức `Generate` sẽ gọi phương thức được custom `genServiceCode` để sinh ra mã nguồn cho mỗi service.

Hiện tại, phương thức `genImportCode` và `genServiceCode` chỉ đơn giản là có một dòng comment đơn giản

```go
func (p *netrpcPlugin) genImportCode(file *generator.FileDescriptor) {
    p.P("// TODO: import code")
}

func (p *netrpcPlugin) genServiceCode(svc *descriptor.ServiceDescriptorProto) {
    p.P("// TODO: service code, Name = " + svc.GetName())
}
```

Để sử dụng plugin, chúng ta cần phải đăng kí plugin đó với hàm `generator.RegisterPlugin`, chúng có thể được hoàn tất nhờ vào hàm init.

```go
func init() {
    generator.RegisterPlugin(new(netrpcPlugin))
}
```

Bởi vì trong ngôn ngữ Go, package chỉ được import một cách tĩnh, chúng ta không thể thêm plugin mới vào plugin đã có sẵn là `protoc-gen-go`. Chúng ta sẽ `re-clone` lại hàm main để tương ứng với `protoc-gen-go`

```go
package main

import (
    "io/ioutil"
    "os"

    "github.com/golang/protobuf/proto"
    "github.com/golang/protobuf/protoc-gen-go/generator"
)

func main() {
    g := generator.New()

    data, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        g.Error(err, "reading input")
    }

    if err := proto.Unmarshal(data, g.Request); err != nil {
        g.Error(err, "parsing input proto")
    }

    if len(g.Request.FileToGenerate) == 0 {
        g.Fail("no files to generate")
    }

    g.CommandLineParameters(g.Request.GetParameter())

    // Create a wrapped version of the Descriptors and EnumDescriptors that
    // point to the file that defines them.
    g.WrapTypes()

    g.SetPackageNames()
    g.BuildTypeNameMap()

    g.GenerateAllFiles()

    // Send back the results.
    data, err = proto.Marshal(g.Response)
    if err != nil {
        g.Error(err, "failed to marshal output proto")
    }
    _, err = os.Stdout.Write(data)
    if err != nil {
        g.Error(err, "failed to write output proto")
    }
}
```

Để tránh việc trùng tên với protoc-gen-go plugin, chúng ta sẽ đặt tên cho chương trình thực thi trên là `protoc-gen-go-netrpc`, điều đó có ý nghĩa là `netrpc` plugin đã được bao gồm. Sau đó chúng ta sẽ biên dịch lại hello.proto file với lệnh sau

```
$ protoc --go-netrpc_out=plugins=netrpc:. hello.proto
```

Tham số `--go-netrpc_out` sẽ nói cho bộ biên dịch protoc biết là nó phải tải một plugin với tên gọi là protoc-gen-go-netrpc, với `plugins=netrpc` chỉ ra rằng những plugin nội bộ `netrpcPlugin` được kích hoạt. Chúng ta thêm vào chú thích cho mã nguồn trên, chúng sẽ được included trong file `hello.pb.go` mới sinh ra.

Tại thời điểm này, plugin được chúng ta tạo ra cuối cùng đã hoạt động.

## 4.2.3 Tự động sinh ra toàn bộ mã nguồn RPC

Trong ví dụ trước chúng ta đã xây dựng một plugin nho nhỏ là `netrpcPlugin` và tạo ra một plugin mới cho `protoc-gen-go-netrpc` bởi việc sao chép lại chương trình chính của protoc-gen-go. Bây giờ tiếp tục phát triển netrpcPlugin plugin với mục tiêu cuối cùng là sinh ra lớp giao diện RPC bảo mật.
Đầu tiên chúng ta sẽ phải sinh ra mã nguồn của package được import bằng phương thức đã được định nghĩa genImportCode.


```go
func (p *netrpcPlugin) genImportCode(file *generator.FileDescriptor) {
    p.P(`import "net/rpc"`)
}
```

Sau đó sinh ra những mã nguồn liên quan cho mỗi service của phương thức genServiceCode được tạo ra. Chúng ta có thể phân tích thấy rằng thứ quan trọng nhất của mỗi service là tên của service, và sau đó mỗi service sẽ có một tập hợp các phương thức. Việc định nghiã phương thức có thành phần quan trọng nhất là tên của service cũng như là tham số đầu vào và tham số đầu ra.

Chúng ta sẽ định nghĩa kiểu ServiceSpec được mô tả như là thông tin thêm vào của service.


```go
type ServiceSpec struct {
    ServiceName string
    MethodList  []ServiceMethodSpec
}

type ServiceMethodSpec struct {
    MethodName     string
    InputTypeName  string
    OutputTypeName string
}
```

Chúng ta sẽ tạo ra một phương thức mới là `buildServiceSpec` nó sẽ parse thông tin thêm vào service được định nghĩa trong ServiceSpec cho mỗi service.

```go
func (p *netrpcPlugin) buildServiceSpec(
    svc *descriptor.ServiceDescriptorProto,
) *ServiceSpec {
    spec := &ServiceSpec{
        ServiceName: generator.CamelCase(svc.GetName()),
    }

    for _, m := range svc.Method {
        spec.MethodList = append(spec.MethodList, ServiceMethodSpec{
            MethodName:     generator.CamelCase(m.GetName()),
            InputTypeName:  p.TypeName(p.ObjectNamed(m.GetInputType())),
            OutputTypeName: p.TypeName(p.ObjectNamed(m.GetOutputType())),
        })
    }

    return spec
}
```

Kiểu tham số đầu vào là `*descriptor.ServiceDescriptorProto` nó sẽ hoàn toàn mô tả tất cả những thông tin về service. Sau đó `svc.GetName()` sẽ lấy tên của service được định nghĩa ở protobuf file. Sau khi tên của Protobuf file được chuyển thành tên trong ngôn ngữ Go , `generator.CamelCase` một sự chuyển đổi được yêu cầu trong hàm đó. Tương tự, trong vòng lặp chúng ta sẽ `m.GetName()` để lấy ra tên của phương thức và sau đó thay đổi chúng để ứng với tên trong ngôn ngữ Go. Vơi sự phức tạp của việc phân tích tên tham số đầu vào và đầu ra: đầu tiên chúng ta cần có `m.GetInputType()` để lấy kiểu dữ liệu tham số đầu vào, sau đó là `p.ObjectNamed` để đạt được thông tin về class của đối tượng tương ứng với kiểu đó.

Sau đó chúng ta sẽ sinh ra mã nguồn của service dựa trên thông tin mô tả đó, được xây dựng bởi phương thức `buildServiceSpec` :

```go
func (p *netrpcPlugin) genServiceCode(svc *descriptor.ServiceDescriptorProto) {
    spec := p.buildServiceSpec(svc)

    var buf bytes.Buffer
    t := template.Must(template.New("").Parse(tmplService))
    err := t.Execute(&buf, spec)
    if err != nil {
        log.Fatal(err)
    }

    p.P(buf.String())
}
```

Để dễ dàng bảo trì, chúng ta sẽ phải sinh ra service code dựa trên template của ngôn ngữ Go, khi đó `tmplService` là template của service.

Bởi việc viết template, hãy nhìn  những gì chúng ta mong đợi vào mã nguồn cuối cùng được sinh ra

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

HelloService là tên của service, bên cạnh mỗi chuỗi những tên service liên quan khác.

Template sau đây sẽ có thể được xây dựng với sự reference tới mã nguồn cuối cùng được sinh ra.

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
{{end}}
`
```

Khi plugin của protobuf được tùy chỉnh hoàn thành, mã nguồn có thể được từ động sinh ra mỗi thời điểm mà RPC service thay đổi trong hello.proto file. Chúng ta có thể điều chỉnh hoặc tăng nội dung của mã nguồn được sinh ra bởi việc cập nhật template được plugin. Sau khi chúng ta đã có thể xây dựng một plugin riêng, chúng ta sẽ hoàn thành toàn bộ công nghệ đó.

