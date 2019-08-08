# 4.1 Bắt đầu với RPC

[Remote Procedure Call](https://en.wikipedia.org/wiki/Remote_procedure_call) (viết tắt: RPC) là phương pháp gọi hàm từ một máy tính từ xa để lấy về kết quả. Trong lịch sử phát triển của internet, RPC đã trở thành một cơ sở hạ tầng không thể thiếu cũng giống như là IPC (inter process communication) ngoài việc chúng dùng để giao tiếp giữa các máy tính chứ không những là giữa các tiến trình, ngoài ra RPC còn hay được sử dụng trong các hệ thống phân tán.

<div align="center">
	<img src="../images/ch4-1-rpc-arch.png" width="500">
    <br/>
    <span align="center">
		<i>Mô hình giao tiếp client/server trong RPC</i>
	</span>
</div>

## 4.1.1 Chương trình "Hello, World" bằng RPC

Thư viện chuẩn của Go chứa gói [net/rpc](https://golang.org/pkg/net/rpc/) dùng để hiện thực chương trình RPC, chương trình RPC đầu tiên của chúng ta sẽ in ra chuỗi "Hello, World" được tạo ra và trả về từ máy khác:

***service/hello.go***: định nghĩa service Hello

```go
package service

// định nghĩa struct register service
type HelloService struct{}

// định nghĩa hàm service Hello, quy tắc:
// 1. Hàm service phải public (viết hoa)
// 2. Có hai tham số trong hàm
// 3. Tham số thứ hai phải kiểu con trỏ
// 4. Phải trả về kiểu error

func (p *HelloService) Hello(request string, reply *string) error {
    *reply = "Hello, " + request
    // trả về error = nil nếu thành công
    return nil
}

```

***server/main.go:*** chương trình phía server

```go
package main

import (
    "log"
    "net"
    "net/rpc"

    // import rpc service
    "../service"
)

func main() {
    // đăng kí tên service với đối tượng rpc service
    rpc.RegisterName("HelloService", new(service.HelloService))
    // chạy rpc server trên port 1234
    listener, err := net.Listen("tcp", ":1234")
    // nếu có lỗi xảy ra thì in ra
    if err != nil {
    log.Fatal("ListenTCP error:", err)
    }
    // vòng lặp để xử lý nhiều kết nối client
    for {
        // chấp nhận một connection đến
        conn, err := listener.Accept()
        // in ra nếu bị lỗi khi Accept
        if err != nil {
            log.Fatal("Accept error:", err)
        }
        // phục vụ RPC cho client trên một goroutine khác
        // để giải phóng main thread tiếp tục kết nối client khác
        go rpc.ServeConn(conn)
    }
}
```


***client/main.go:*** mã nguồn client để gọi service Hello

```go
package main

import (
    "fmt"
    "log"
    "net/rpc"
)

func main() {
    // kết nối đến rpc server
    client, err := rpc.Dial("tcp", "localhost:1234")
    // in ra lỗi nếu có
    if err != nil {
        log.Fatal("dialing:", err)
    }
    // biến chứa giá trị trả về sau lời gọi rpc
    var reply string
    // gọi rpc với tên service đã register, tham số và biến
    err = client.Call("HelloService.Hello", "World", &reply)
    if err != nil {
        log.Fatal(err)
    }
    // in ra kết quả
    fmt.Println(reply)
}
```

Kết quả khi chạy Hello Service :

```sh
$ go run server/main.go
```

Ở một terminal khác chạy client:

```sh
$ go run client/main.go
Hello, World
```

Từ ví dụ trên, có thể thấy rằng chúng ta dùng RPC trong Go thật sự đơn giản.

## 4.1.2 Tạo interface cho RPC

Ứng dụng gọi RPC sẽ có ít nhất ba thành phần: thứ nhất là chương trình hiện thực phương thức RPC ở bên phía server, thứ hai là chương trình gọi RPC bên phía client, và cuối cùng là thành phần cực kì quan trọng: service đóng vai trò là interface giữa server và client.

Trong ví dụ trước, chúng ta đã đặt tất cả những thành phần trên trong ba thư mục **server**, **client**, **service**, nếu bạn muốn refactor lại mã nguồn HelloService, đầu tiên hãy tạo ra một inteface như sau:

***Interface của RPC service:***

```go
// tên của service, chứa tiền tố pkg để tránh xung đột tên về sau
const HelloServiceName = "path/to/pkg.HelloService"
// interface RPC của HelloService
type HelloServiceInterface = interface {
    // định nghĩa danh sách các function trong service
    Hello(request string, reply *string) error
}
// hàm đăng kí service
func RegisterHelloService(svc HelloServiceInterface) error {
    // gọi hàm register của gói net/rpc
    return rpc.RegisterName(HelloServiceName, svc)
}
```

Sau khi định nghĩa lớp interface của RPC service, client có thể viết mã nguồn để gọi lệnh RPC :

***Hàm main phía client:***

```go
// hàm main bên phía client
func main() {
    // kết nối rpc server qua port 1234
    client, err := rpc.Dial("tcp", "localhost:1234")
    // log ra lỗi nếu có
    if err != nil {
        log.Fatal("dialing:", err)
    }
    // biến chứa kết quả sau khi gọi RPC
    var reply string
    // gọi hàm RPC được định nghĩa phía server
    err = client.Call(service.HelloServiceName+".Hello", "hello", &reply)
    // log ra chi tiết lỗi nếu có
    if err != nil {
        log.Fatal(err)
    }
}
```

Tuy nhiên, gọi phương thức RPC thông qua hàm `client.Call` vẫn rất cồng kềnh, để đơn giản chúng ta nên wrapper biến connection vào trong:

***Wrapper các đối tượng:***

```go
// struct chứa đối tượng
type HelloServiceClient struct {
    // wrapper server connection
    *rpc.Client
}
// tạo hàm wrapper lời gọi Dial tới server
func DialHelloService(network, address string) (*HelloServiceClient, error) {
    // gọi Dial tới server bên trong
    c, err := rpc.Dial(network, address)
    // trả về rỗng và lỗi nếu có
    if err != nil {
        return nil, err
    }
    // trả về rpc struct và error=nil nếu thành công
    return &HelloServiceClient{Client: c}, nil
}
// wrapper lại lời gọi hàm Hello phía client
func (p *HelloServiceClient) Hello(request string, reply *string) error {
    return p.Client.Call(HelloServiceName+".Hello", request, reply)
}
```

Dựa trên các hàm wrapper trên, chúng ta sẽ viết lại mã nguồn phía client:

***Hàm main phía client sau khi refactor:***

```go
func main() {
    // kết nối RPC server bằng hàm wrapper
    client, err := DialHelloService("tcp", "localhost:1234")
    // log ra lỗi nếu có
    if err != nil {
        log.Fatal("dialing:", err)
    }
    // biến lưu kết quả từ lời gọi RPC
    var reply string
    // thực thi lệnh gọi RPC
    err = client.Hello("World", &reply)
    // log ra lỗi nếu có
    if err != nil {
        log.Fatal(err)
    }
}
```

Cuối cùng, mã nguồn server thực sự sẽ được viết lại như sau:

***Chương trình bên phía server:***

```go
// đối tượng RPC HelloService
type HelloService struct {}
// hiện thực lời gọi RPC
func (p *HelloService) Hello(request string, reply *string) error {
    *reply = "Hello, " + request
    return nil
}
// hàm main phía server
func main() {
    // gọi wrapper đăng ký đối tượng HelloService 
    RegisterHelloService(new(HelloService))
    // lắng nghe kết nối từ phía client
    listener, err := net.Listen("tcp", ":1234")
    // log ra lỗi nếu có (vd: trùng port, v,v..)
    if err != nil {
        log.Fatal("ListenTCP error:", err)
    }
    // vòng lặp tiếp nhận nhiều kết nối client
    for {
        // chấp nhận kết nối từ một client nào đó
        conn, err := listener.Accept()
        // in ra lỗi nếu có
        if err != nil {
            log.Fatal("Accept error:", err)
        }
        // phục vụ kết nối trên một goroutine khác
        // để main thread tiếp tục vòng lặp accept client khác
        go rpc.ServeConn(conn)
    }
}
```

Ở phiên bản refactor, chúng ta sử dụng hàm `RegisterHelloService` để đăng ký RPC service, nó tránh việc trực tiếp đặt tên cho service, và đảm bảo bất cứ đối tượng nào hiện thực các hàm trong interface của RPC service cũng đều có thể phục vụ lời gọi RPC từ phía client.

## 4.1.3 Vấn đề gọi RPC trên các ngôn ngữ khác nhau:

Trong hệ thống microservice, mỗi dịch vụ có thể viết bằng các ngôn ngữ lập trình khác nhau, do đó để **cross-language** (vượt qua rào cản ngôn ngữ) là điều kiện thiết yếu cho sự tồn tại của RPC trong môi trường internet.

Thư viện chuẩn RPC của Go mặc định đóng gói dữ liệu theo đặc tả của [Go encoding](https://golang.org/pkg/encoding/), do đó sẽ rất khó để gọi RPC service từ những ngôn ngữ khác.

May mắn là thư viện `net/rpc` của Go có ít nhất hai thiết kế đặc biệt:
   * Một là cho phép chúng ta có thể thay đổi quá trình encoding và decoding gói tin RPC.
   * Hai là interface RPC được xây dựng dựa trên interface `io.ReadWriteClose`, chúng ta có thể  xây dựng RPC trên những protocol giao tiếp khác nhau.

Từ đây chúng ta có thể hiện thực việc cross-language thông qua gói `net/rpc/jsonrpc` :

***Hàm main mới phía server:***

```go
package main

import (
    "log"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
)

// định nghĩa struct register service
type HelloService struct{}

func (p *HelloService) Hello(request string, reply *string) error {
    *reply = "Hello, " + request
    // trả về error = nil nếu thành công
    return nil
}

func main() {
    // đăng kí HelloService (dùng cách cũ cho đơn giản)
    rpc.RegisterName("HelloService", new(HelloService))
    // lắng nghe connection từ phía client
    listener, err := net.Listen("tcp", ":1234")
    // in ra lỗi (vd: trùng port,..) nếu có
    if err != nil {
        log.Fatal("ListenTCP error:", err)
    }
    // thực hiện vòng lặp phục vụ nhiều RPC client
    for {
        // chấp nhận kết nối từ RPC client
        conn, err := listener.Accept()
        // in ra lỗi nếu có
        if err != nil {
            log.Fatal("Accept error:", err)
        }
        // phục vụ client trên một goroutine khác, lúc này:
        // 1. rpc.ServeConn được thay thế bằng rpc.ServeCodec
        // 2. dùng jsonrpc.NewServerCodec để bao đối tượng conn
        go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
}
```

Sau đó, client sẽ hiện thực phiên bản json :

***Hàm main bên phía client:***

```go
package main

import (
    "fmt"
    "log"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
)

func main() {
    // kết nối đến RPC server
    conn, err := net.Dial("tcp", "localhost:1234")
    // in ra lỗi nếu có
    if err != nil {
        log.Fatal("net.Dial:", err)
    }
    // gọi dịch vụ RPC Server được encoding bằng json Codec
    client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
    // biến lưu giá trị sau lời gọi hàm rpc
    var reply string
    // gọi dịch vụ RPC
    err = client.Call("HelloService.Hello", "World", &reply)
    // in ra lỗi nếu có
    if err != nil {
        log.Fatal(err)
    }
    // in ra kết quả của lệnh gọi RPC
    fmt.Println(reply)
}
```

***Kết quả:***
  * Chạy server:
```sh
$ go run server/main.go
```
  * Chạy client:
```sh
$ go run client/main.go
Hello, World
```

Để thấy dữ liệu được client gửi cho server, đầu tiên tắt chương trình server và gọi lệnh [nc](http://www.tutorialspoint.com/unix_commands/nc.htm) :


```sh
$ go run server/main.go
// Ctrl+C
$ nc -l 1234
```

Sau đó gọi chương trình client: `$ go run client/main.go` một lần nữa, ta sẽ thấy kết quả:

```sh
$ nc -l 1234
{"method":"HelloService.Hello","params":["World"],"id":0}
```

Ngược lại, nếu muốn thấy thông điệp mà phía server gửi cho client, chạy RPC service phía server: `$ go run server/main.go` và ở một terminal khác chạy lệnh:

```
$ echo -e '{"method":"HelloService.Hello","params":["World"],"id":1}' | nc localhost 1234
// kết quả mà RPC server trả về 
{"id":1,"result":"Hello, World","error":null}
// Trong đó:
// - "id" : để nhận dạng kết quả ứng với yêu cầu vì việc thực thi lời gọi RPC là bất đồng bộ
// - "result" : kết quả trả về của lời gọi hàm
// - "error" : chứa thông điệp lỗi nếu có
```

Dữ liệu json ở hai ví dụ trên sẽ tương ứng với hai cấu trúc sau: 


```go
// cấu trúc json phía client
type clientRequest struct {
    // tên phương thức RPC
    Method string         `json:"method"`
    // parameter truyền vào
    Params [1]interface{} `json:"params"`
    // id của lời gọi RPC
    Id     uint64         `json:"id"`
}

// cấu trúc json phía server
type serverRequest struct {
    // tên phương thức RPC
    Method string           `json:"method"`
    // tham số truyền vào
    Params *json.RawMessage `json:"params"`
    // thông tin về index
    Id     *json.RawMessage `json:"id"`
}
```

Ta có thể thấy rằng, chỉ cần theo định dạng json như trên là có thể giao tiếp với RPC service được viết bởi Go hay bất kỳ ngôn ngữ nào khác, nói cách khác chúng ta có thể hiện thực việc cross-language trong RPC.

## 4.1.4 Go RPC qua giao thức HTTP

Trong ví dụ trước, chúng ta đã gọi RPC thông qua lệnh `nc`, bây giờ chúng ta sẽ thử cung cấp RPC service trên giao thức HTTP. RPC Service mới sẽ tuân thủ theo chuẩn [REST](https://restfulapi.net/), chúng sẽ nhận yêu cầu và xử lý chúng như dưới đây:

***Chương trình phía server:***

```go
func main() {
    // đăng ký tên của RPC service
    rpc.RegisterName("HelloService", new(HelloService))
    // routing uri /jsonrpc đến hàm xử lý tương ứng
    http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
        // conn là một biến thuộc kiểu io.ReadWriteCloser
        var conn io.ReadWriteCloser = struct {
            // là struct gồm hai biến io đọc và ghi
            io.Writer
            io.ReadCloser
        }{  // được khởi tạo với nội dung:
            // ReadCloser là nội dung nhận được
            ReadCloser: r.Body,
            // Writer là đối tượng dùng ghi kết quả
            Writer: w,
        }
        // truyền dịch vụ RPC với biến conn
        rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
    })
    // lắng nghe kết nối từ client trên port 1234
    http.ListenAndServe(":1234", nil)
}

```

Lệnh gọi RPC để gửi chuỗi json đến kết nối đó :

``` 
$ curl localhost:1234/jsonrpc -X POST \
    --data '{"method":"HelloService.Hello","params":["hello"],"id":0}'
```

Kết quả vẫn là một chuỗi json :

```
{"id":0,"result":"hello:hello","error":null}
```

Điều đó làm việc gọi RPC service từ những ngôn ngữ khác dễ dàng hơn.