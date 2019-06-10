# 4.1 Bắt đầu với RPC

RPC là viết tắt của remote procedure call (lời gọi hàm từ xa) và nó là một cách thức giao tiếp giữa các nốt của distributed system (hệ phân bố). Trong lịch sử của internet, RPC đã trở thành một cơ sở hạ tầng không thể thiếu cũng như là IPC (inter process communication- giao tiếp giữa các tiến trình). Do đó, thư viện chuẩn của Go đã hỗ trợ phiên bản hiện thực RPC đơn giản, và chúng ta sẽ dùng chúng như là một đối tượng để học RPC.

## 4.1.1 RPC phiên bản "Hello, World"

Một nhánh hiện thực RPC của ngôn ngữ Go là `net/rpc`, nó sẽ nằm dưới đường dẫn của gói `net`. Do đó chúng ta có thể đoán được rằng gói RPC được hiện thực dựa trên gói `net`. Tại phần cuối của chương "Cuộc cách mạng Hello, World", chúng ta đã hiện thực việc in ra một ví dụ mẫu dựa trên `http`. Bên dưới chúng ta sẽ thử hiện thực tương tự dựa trên `rpc`.

```go
type HelloService struct {}

func (p *HelloService) Hello(request string, reply *string) error {
    *reply = "hello:" + request
    return nil
}
```

Hàm Hello sẽ phải thỏa mãn những quy luật của RPC trong ngôn ngữ Go: phương thức chỉ có thể có hai tham số để serialize, tham số thứ hai là kiểu con trỏ, và giá trị trả về là kiểu error, và nó phải là một phương thức public.

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

Hàm rpc.Register sẽ đăng kí những đối tượng thỏa mãn quy luật RPC như là RPC functions, và tất cả những phương thức bên dưới không gian  "HelloService" service. Sau đó chúng ta sẽ tạo ra một đường dẫn TCP duy nhất và cung cấp dịch vụ RPC đến bên khác trong qua đường truyền RPC được hỗ trợ bởi hàm `rpc.ServeConn`.

Dưới đây là mã nguồn client để yêu cầu dịch vụ Hello:

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

Lựa chọn đầu tiên là kết nối tới dịch vụ RPC thông qua rpc.Dial và sau đó gọi một phương thức RPC cụ thể thông qua Client.Call. Khi gọi Client.Call, tham số đầu tiên là tên RPC Service và tên phương thức qua dấu ".", và sau đó tham số thứ hai và thứ ba sẽ tương ứng với hai tham số định nghĩa ở phương thức RPC.
Từ ví dụ trên, có thể thấy rằng chúng ta dùng RPC thật sự rất đơn giản.

## 4.1.2 RPC Interface an toàn

Trong ứng dụng gọi RPC, sẽ thường có ít nhất ba vai trò như là developers: người thứ nhất là nhà phát triển sẽ hiện thực phương thức RPC ở bên phía server, người thứ hai là người gọi RPC bên phía client, và người cuối cùng là người cực kì quan trọng, họ sẽ phát triển server và client RPC - người đặc tả giao diện. Trong ví dụ trước, chúng ta đặt tất cả những vai trò trên lại với nhau cho đơn giản, mặc dùng nó dường như là đơn giản để hiện thực, nhưng nó không thuận lợi cho việc bảo trì và phân công công việc về sau.

Nếu bạn muốn refactor lại dịch vụ HelloService, bước đầu tiên là phân định rạch ròi giữa tên và giao diện của service;

```go
const HelloServiceName = "path/to/pkg.HelloService"

type HelloServiceInterface = interface {
    Hello(request string, reply *string) error
}

func RegisterHelloService(svc HelloServiceInterface) error {
    return rpc.RegisterName(HelloServiceName, svc)
}
```

Chúng ta chia đặc tả giao diện của dịch vụ RPC thành ba phần: đầu tiên, tên của dịch vụ, sau đó là danh sách chi tiết của những phương thức cần hiện thực của dịch vụ. Để tránh xung đột tên, chúng ta thêm vào tên của gói thành tiền tố của tên dịch vụ RPC (nó là đường dẫn đến gói của lớp dịch vụ trừu tượng RPC, không phải là đường dẫn tới gói của ngôn ngữ Go). Chúng ta sẽ đăng kí phương thức RegisterHelloService đến dịch vụ và bộ biên dịch sẽ hỏi những đối tượng được mang tới để thỏa mãn giao diện HelloServiceInterface.

Sau khi định nghĩa lớp giao diện của dịch vụ RPC, client có thể ghi mã nguồn để gọi lệnh RPC theo đặc tả như sau:

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

Sự thay đổi nhỏ được nhìn thấy ở đối số đầu tiên của hàm client.Call thay thế `HelloService.Hello` với 
`HelloServiceName+".Hello"`. Tuy nhiên, gọi phương thức RPC thông qua hàm client.Call vẫn rất cồng kềnh, và kiểu của tham số không thể có tính an toàn của trình biên dịch.

Để đơn giản lời gọi từ Client tới hàm RPC, chúng ta thêm vào hàm wrapper ở phía client trong giao diện đặc tả như sau:

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

