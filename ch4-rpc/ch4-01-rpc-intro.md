# 4.1 Bắt đầu với RPC

RPC viết tắt của remote procedure call (lời gọi hàm từ xa) là một cách thức giao tiếp giữa các node của distributed system (hệ thống phân tán). Trong lịch sử của internet, RPC đã trở thành một cơ sở hạ tầng không thể thiếu giống như là IPC (inter process communication- giao tiếp giữa các tiến trình). Do đó, thư viện chuẩn của Go đã hỗ trợ phiên bản hiện thực RPC đơn giản, và chúng ta sẽ dùng chúng như là một đối tượng để học RPC.

## 4.1.1 RPC phiên bản "Hello World"

Package RPC trong ngôn ngữ Go là `net/rpc`, được đặt trong thư mục package `net`. Do đó chúng ta có thể đoán được rằng package RPC được hiện thực dựa trên package `net`. Tại phần cuối của chương 1 phần "Cuộc cách mạng Hello World", chúng ta đã hiện thực việc in ra một ví dụ mẫu dựa trên `http`. Bên dưới chúng ta sẽ thử hiện thực tương tự dựa trên `rpc`.

```go
type HelloService struct {}

func (p *HelloService) Hello(request string, reply *string) error {
    *reply = "hello:" + request
    return nil
}
```

Hàm Hello sẽ phải thỏa mãn những quy tắt của RPC trong ngôn ngữ Go: phương thức chỉ có thể có hai tham số để serialize, tham số thứ hai là kiểu con trỏ, giá trị trả về là kiểu error và nó phải là một phương thức public.

Sau đó chúng ta có thể đăng kí đối tượng thuộc kiểu HelloService là một RPC Service.

```go
func main() {
    rpc.RegisterName("HelloService", new(HelloService))

    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("ListenTCP error:", err)
    }

    conn, err := listener.Accept()
    if err != nil {
        log.Fatal("Accept error:", err)
    }

    rpc.ServeConn(conn)
}
```


Hàm `rpc.Register` sẽ đăng kí những đối tượng thỏa mãn quy tắt RPC như là RPC functions, và tất cả những phương thức bên dưới không gian "HelloService" service. Sau đó chúng ta sẽ tạo ra một liên kết TCP duy nhất và cung cấp service RPC đến các thành phần khác qua liên kết TCP được hỗ trợ bởi hàm `rpc.ServeConn`.

Dưới đây là mã nguồn client để yêu cầu service Hello:

```go
func main() {
    client, err := rpc.Dial("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }

    var reply string
    err = client.Call("HelloService.Hello", "hello", &reply)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(reply)
}
```

Lựa chọn đầu tiên là kết nối tới service RPC thông qua `rpc.Dial` và sau đó gọi một phương thức RPC cụ thể thông qua `client.Call`. Khi gọi `client.Call`, tham số đầu tiên là tên RPC Service và tên phương thức qua dấu ".", và sau đó tham số thứ hai và thứ ba sẽ tương ứng với hai tham số định nghĩa ở phương thức RPC.
Từ ví dụ trên, có thể thấy rằng chúng ta dùng RPC thật sự rất đơn giản.

## 4.1.2 RPC Interface an toàn

Trong ứng dụng gọi RPC, sẽ thường có ít nhất ba nhà phát triển: người thứ nhất là nhà phát triển sẽ hiện thực phương thức RPC ở bên phía server, người thứ hai là người gọi RPC bên phía client, và người cuối cùng là người cực kì quan trọng, họ sẽ phát triển interface giữa server và client RPC. Trong ví dụ trước, chúng ta đặt tất cả những vai trò trên lại với nhau cho đơn giản, mặc dùng nó dường như là cách đơn giản để hiện thực, nhưng nó không thuận lợi cho việc bảo trì và phân công công việc về sau.

Nếu bạn muốn refactor lại service HelloService, bước đầu tiên là phân định rạch ròi giữa tên và inteface của service;

```go
const HelloServiceName = "path/to/pkg.HelloService"

type HelloServiceInterface = interface {
    Hello(request string, reply *string) error
}

func RegisterHelloService(svc HelloServiceInterface) error {
    return rpc.RegisterName(HelloServiceName, svc)
}
```

Chúng ta chia đặc tả interface của service RPC thành ba phần: đầu tiên là tên của service, sau đó là danh sách chi tiết của những phương thức cần hiện thực của service và cuối cùng là function đăng ký service. Để tránh xung đột tên, chúng ta thêm tiền tố của package vào tên của service RPC (nó là đường dẫn đến package của lớp service trừu tượng RPC, không phải là đường dẫn tới package của ngôn ngữ Go). Chúng ta sẽ đăng kí phương thức `RegisterHelloService` đến service và bộ biên dịch sẽ yêu cầu đối tượng tới để thỏa mãn interface `HelloServiceInterface`.

Sau khi định nghĩa lớp interface của service RPC, client có thể viết mã nguồn để gọi lệnh RPC theo đặc tả như sau:

```go
func main() {
    client, err := rpc.Dial("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }

    var reply string
    err = client.Call(HelloServiceName+".Hello", "hello", &reply)
    if err != nil {
        log.Fatal(err)
    }
}
```

Sự thay đổi duy nhất là đối số đầu tiên của hàm client.Call thay thế `HelloService.Hello` với `HelloServiceName+".Hello"`. Tuy nhiên, gọi phương thức RPC thông qua hàm `client.Call` vẫn rất cồng kềnh, và kiểu của tham số không thể có tính an toàn do trình biên dịch cung cấp.

Để đơn giản lời gọi từ Client tới hàm RPC, chúng ta thêm vào hàm wrapper ở phía client trong interface đặc tả như sau:

```go
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

func (p *HelloServiceClient) Hello(request string, reply *string) error {
    return p.Client.Call(HelloServiceName+".Hello", request, reply)
}
```

Chúng ta thêm một kiểu mới là `HelloServiceClient` bên phía client trong đặc tả. Kiểu này cũng phải thõa mãn interface `HelloServiceInterface`, do đó client cần phải trực tiếp gọi phương thức RPC thông qua hàm tương ứng của interface đó. Đồng thời, phương thức DialHelloService được cung cấp trực tiếp để gọi service HelloService.

Dựa trên interface client mới, chúng ta sẽ đơn giản hóa mã nguồn bên phía  client như sau:

```go
func main() {
    client, err := DialHelloService("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }

    var reply string
    err = client.Hello("hello", &reply)
    if err != nil {
        log.Fatal(err)
    }
}
```


Giờ đây, client không còn phải lo lắng về low-level errors như là tên phương thức RPC hoặc kiểu dữ liệu không trùng khớp.

Cuối cùng, mã nguồn server thực sự sẽ được viết dựa trên interface đặc tả RPC

```go
type HelloService struct {}

func (p *HelloService) Hello(request string, reply *string) error {
    *reply = "hello:" + request
    return nil
}

func main() {
    RegisterHelloService(new(HelloService))

    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("ListenTCP error:", err)
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal("Accept error:", err)
        }

        go rpc.ServeConn(conn)
    }
}
```


Ở phiên bản hiện thực RPC mới, chúng ta sử dụng hàm `RegisterHelloService` để đăng kí, nó không chỉ tránh công việc đặt tên cho service với những tên gọi của service, mà còn đảm bảo rằng những đối tượng của service mang đến sẽ thỏa mãn định nghĩa của interface RPC. Cuối cùng, service mới của chúng ta sẽ hỗ trợ nhiều liên kết TCP và do đó sẽ cung cấp service RPC cho mỗi đường dẫn TCP.


## 4.1.3 Cross-language RPC (đa ngôn ngữ trên RPC)

Thư viện chuẩn của RPC sẽ mặc định đóng gói dữ liệu theo đặc tả của Go encoding, do đó sẽ khó hơn nhiều để gọi service RPC từ những ngôn ngữ khác. Trong những micro-service trên môi trường mạng, mỗi RPC và người dùng dịch vụ có thể sử dụng những ngôn ngữ lập trình khác nhau, do đó để cross-language (vượt qua rào cản ngôn ngữ) là điều kiện chính cho sự tồn tại của RPC trên môi trường internet.

Framework RPC của ngôn ngữ Go có nhiều hơn hai thiết kế đặc biệt: một là cho phép chúng ta có thể thay đổi quá trình encoding và decoding trong quá trình kết nối khi gói dữ liệu được đóng gói; và hai là interface RPC được xây dựng dựa trên interface `io.ReadWriteClose`, chúng ta có thể  xây dựng RPC trên những protocol giao tiếp khác nhau. Từ đây chúng ta có thể hiện thực việc cross-language thông qua phần mở rộng của `net/rpc/jsonrpc`

Đầu tiên chúng ta có thể hiện thực lại RPC service dựa trên json encoding như sau:

```go
func main() {
    rpc.RegisterName("HelloService", new(HelloService))

    listener, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("ListenTCP error:", err)
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal("Accept error:", err)
        }

        go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
}
```


Sự thay đổi lớn nhất trong mã nguồn là thay thế hàm `rpc.ServeConn` với `rpc.ServeCodec`. Tham số được truyền ở trong là json codec cho server.

Sau đó, client sẽ hiện thực phiên bản json như sau:

```go
func main() {
    conn, err := net.Dial("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("net.Dial:", err)
    }

    client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

    var reply string
    err = client.Call("HelloService.Hello", "hello", &reply)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(reply)
}
```

Đầu tiên chúng ta sẽ gọi hàm `net.Dial` để thiết lập kết nối TCP, sau đó là xây dựng json codec cho client dựa trên liên kết đó.

Sau khi đảm bảo rằng client có thể gọi RPC service một cách bình thường, chúng ta sẽ thay thế phiên bản ngôn ngữ Go với TCP service bình thường, do đó chúng ta có thể thấy dữ liệu được định dạng được gửi đến client. Ví dụ `nc -l 1234`, bắt đầu một TCP service trong cùng một port sử dụng lệnh `nc`. Sau đó thực thi lời gọi RPC một lần nữa để thấy rằng kết quả của `nc` sẽ là thông tin sau

```
{"method":"HelloService.Hello","params":["hello"],"id":0}
```

Đây không phải là dữ liệu json-encoded, khi một phần của phương thức ứng với tên của rpc service và tên hàm được gọi, phần tử đầu tiên của "params" là tham số, và số id được đảm bảo phải duy nhất bởi phía gọi.

Đối tượng dữ liệu json sẽ tương ứng với hai cấu trúc sau: bên phía client là clientRequest và bên phía server là serverRequest. Nội dung của cấu trúc clientRequest và serverRequest về cơ bản sẽ giống nhau:

```go
type clientRequest struct {
    Method string         `json:"method"`
    Params [1]interface{} `json:"params"`
    Id     uint64         `json:"id"`
}

type serverRequest struct {
    Method string           `json:"method"`
    Params *json.RawMessage `json:"params"`
    Id     *json.RawMessage `json:"id"`
}
```

Sau khi định nghĩa kiểu dữ liệu json để gọi RPC, chúng ta có thể gửi dữ liệu json để mô phỏng lệnh gọi RPC một cách trực tiếp đến RPC server mà xây dựng RPC service

```
$ echo -e '{"method":"HelloService.Hello","params":["hello"],"id":1}' | nc localhost 1234
```

Kết quả trả về cũng là chuỗi dữ liệu json được định dạng như sau

```
{"id":1,"result":"hello:hello","error":null}
```

Trong khi id tương ứng với tham số input id, kết quả là giá trị của "result" và phần "error" sẽ chỉ ra thông điệp lỗi khi có vấn đề xảy ra. Cho chuỗi các lệnh tuần tự, id không được yêu cầu phải có. Tuy nhiên, framework RPC của ngôn ngữ Go sẽ hỗ trợ lệnh gọi bất đồng bộ. Khi thứ tự của kết quả trả về không tương ứng với thứ tự của các lần gọi, lệnh gọi tương ứng sẽ được nhận dạng bởi id.

Kết quả dữ liệu json được trả về sẽ tương ứng với hai thành phần bên trong, đối với phía client là clientResponse, và phía server là serverResponse. Nội dung của hai cấu trúc trên cũng sẽ tương tự nhau

```
type clientResponse struct {
    Id     uint64           `json:"id"`
    Result *json.RawMessage `json:"result"`
    Error  interface{}      `json:"error"`
}

type serverResponse struct {
    Id     *json.RawMessage `json:"id"`
    Result interface{}      `json:"result"`
    Error  interface{}      `json:"error"`
}
```

Do đó không có vấn đề gì về rào cản ngôn ngữ, chỉ theo định dạng của kiểu dữ liệu json trên, chúng ta có thể giao tiếp với nhiều RPC service được viết bởi Go hay những ngôn ngữ khác. Do đó chúng ta hoàn toàn có thể hiện thực việc cross-language trong RPC.

## 4.1.4 RPC trên HTTP

RPC framework sẽ thừa hưởng từ ngôn ngữ Go đã hỗ trợ sẵn dịch vụ RPC trên giao thức HTTP. Tuy nhiên, frameword http service cũng có giao thức được xây dựng sẵn, và nó không cung cấp interface để sử dụng cho những protocol khác. Trong ví dụ trước, chúng ta sẽ hiện thực jsonrpc service dựa trên giao thức TCP, và đã hiện thực thành công lời gọi RPC thông qua lệnh `nc`. Bây giờ chúng ta sẽ thử cung cấp service rpcjson trên giao thức HTTP.

RPC Service  mới sẽ thực sự tuân thủ theo chuẩn interface REST, do đó chúng sẽ nhận yêu cầu và xử lý chúng theo quá trình bên dưới

```go
func main() {
    rpc.RegisterName("HelloService", new(HelloService))

    http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
        var conn io.ReadWriteCloser = struct {
            io.Writer
            io.ReadCloser
        }{
            ReadCloser: r.Body,
            Writer:     w,
        }

        rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
    })

    http.ListenAndServe(":1234", nil)
}
```


RPC Service sẽ thiết lập đường dẫn `/jsonrpc` và kênh `conn` thuộc kiểu `io.ReadWriteCloser` được xây dựng dựa trên tham số thuộc kiểu `http.ResponseWriter` và `http.Request`. Một json codec cho server sẽ được xây dựng dựa trên `conn`. Cuối cùng phương thức gọi RPC được xử lý một lần cho mỗi request thông qua hàm `rpc.ServeRequest`.

Quá trình để mô phỏng lệnh gọi RPC để gửi chuỗi json đến kết nối đó như sau:

```
$ curl localhost:1234/jsonrpc -X POST \
    --data '{"method":"HelloService.Hello","params":["hello"],"id":0}'
```

Kết quả vẫn là một chuỗi json như sau

```
{"id":0,"result":"hello:hello","error":null}
```

Điều đó làm việc gọi service RPC từ những ngôn ngữ khác dễ dàng hơn.
