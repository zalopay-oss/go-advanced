# 4.5 gRPC Nâng cao

Các framework RPC cơ bản thường gặp phải nhiều vấn đề về bảo mật và khả năng mở rộng. Phần này sẽ mô tả ngắn gọn một số cách xác thực an toàn bằng gRPC. Sau đó giới thiệu tính năng interceptor trên gRPC và cách triển khai cơ chế xác thực Token một cách tốt nhất, theo dõi các lời gọi RPC và bắt các Panic thông qua interceptor. Cuối cùng là cách gRPC service kết hợp với Web service khác như thế nào.

## 4.5.1 Xác thực qua chứng chỉ

gRPC được xây dựng dựa trên giao thức HTTP/2 và hỗ trợ TLS rất tốt. gRPC service trong chương trước chúng tôi không cung cấp hỗ trợ chứng chỉ, vì vậy client `grpc.WithInsecure()` có thể  thông qua tùy chọn mà bỏ qua việc xác thực chứng chỉ   trong server được kết nối. gRPC service không có chứng chỉ được kích hoạt sẽ phải giao tiếp hoàn toàn bằng plain-text với client và có nguy cơ cao bị giám sát bởi một bên thứ ba khác. Để đảm bảo rằng giao tiếp gRPC không bị giả mạo hoặc giả mạo bởi các bên thứ ba, chúng ta có thể kích hoạt mã hóa TLS trên server.

Bạn có thể tạo private key và certificate (chứng chỉ) cho server và client riêng biệt bằng các lệnh sau:

```sh
$ openssl genrsa -out server.key 2048
$ openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=grpc-server/CN=server.grpc.io" \
    -key server.key -out server.crt

$ openssl genrsa -out client.key 2048
$ openssl req -new -x509 -days 3650 \
    -subj "/C=GB/L=China/O=grpc-client/CN=client.grpc.io" \
    -key client.key -out client.crt
```

Lệnh trên sẽ tạo ra 4 file: *server.key*, *server.crt*, *client.key* và *client.crt*. File private key có phần mở rộng *.key* và được  cần được giữ bảo mật an toàn. File certificate có phần mở rộng *.crt* được hiểu như public key và không cần giữ bí mật.

Với certificate đấy ta có thể truyền nó vào tham số để bắt đầu một gRPC service:

```go
func main() {
    creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
    if err != nil {
        log.Fatal(err)
    }

    server := grpc.NewServer(grpc.Creds(creds))

    ...
}
```

Hàm `credentials.NewServerTLSFromFile` khởi tạo đối tượng certificate từ file cho server, sau đó bọc certificate dưới dạng tùy chọn thông qua hàm `grpc.Creds(creds)` và truyền nó dưới dạng tham số cho hàm `grpc.NewServer`.

Server có thể được xác thực ở client dựa trên chứng chỉ của server và tên của nó:

```go
func main() {
    creds, err := credentials.NewClientTLSFromFile(
        "server.crt", "server.grpc.io",
    )
    if err != nil {
        log.Fatal(err)
    }

    conn, err := grpc.Dial("localhost:5000",
        grpc.WithTransportCredentials(creds),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    ...
}
```

Trong code client mới, ta không còn phụ thuộc trực tiếp vào các file chứng chỉ phía server. Trong lệnh gọi `credentials.NewTLS`, client xác thực server bằng cách đưa vào chứng chỉ CA root và tên của server. Khi client liên kết với server, trước tiên, nó sẽ yêu cầu chứng chỉ của server, sau đó sử dụng chứng chỉ CA root để xác minh chứng chỉ phía server mà nó nhận được.

Nếu chứng chỉ của client cũng được ký bởi chứng chỉ CA root, server cũng có thể thực hiện xác thực chứng chỉ trên client. Ở đây ta sử dụng chứng chỉ CA root để ký chứng chỉ client:

```sh
$ openssl req -new \
    -subj "/C=GB/L=China/O=client/CN=client.io" \
    -key client.key \
    -out client.csr
$ openssl x509 -req -sha256 \
    -CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
    -in client.csr \
    -out client.crt
```

Xem [Makefile](../examples/ch4/ch4.5/1-tls-certificate/tls-config/Makefile)

Chứng chỉ root được cấu hình lúc khởi động server:

```go
func main() {
    certificate, err := tls.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        log.Fatal(err)
    }

    certPool := x509.NewCertPool()
    ca, err := ioutil.ReadFile("ca.crt")
    if err != nil {
        log.Fatal(err)
    }
    if ok := certPool.AppendCertsFromPEM(ca); !ok {
        log.Fatal("failed to append certs")
    }

    creds := credentials.NewTLS(&tls.Config{
        Certificates: []tls.Certificate{certificate},
        ClientAuth:   tls.RequireAndVerifyClientCert, // NOTE: this is optional!
        ClientCAs:    certPool,
    })

    server := grpc.NewServer(grpc.Creds(creds))
    ...
}
```

Server cũng sử dụng hàm `credentials.NewTLS` để tạo chứng chỉ, chọn chứng chỉ CA root thông qua ClientCA và cho phép Client được xác thực bằng tùy chọn `ClientAuth`.

Như vậy chúng ta đã xây dựng được một hệ thống gRPC đáng tin cậy để kết nối giữa Client và Server thông qua xác thực chứng chỉ từ cả 2 chiều.

## 4.5.2 Xác thực token

Xác thực dựa trên chứng chỉ được mô tả ở trên là dành cho từng kết nối gRPC. Ngoài ra gRPC cũng  hỗ trợ xác thực cho mỗi lệnh gọi   gRPC, để việc quản lý quyền có thể thực hiện trên các kết nối khác nhau dựa trên user token.

Để hiện thực cơ chế xác thực cho từng phương thức gRPC, ta cần triển khai interface `grpc.PerRPCCredentials`:

```go
type PerRPCCredentials interface {
    // GetRequestMetadata gets the current request metadata, refreshing
    // tokens if required. This should be called by the transport layer on
    // each request, and the data should be populated in headers or other
    // context. If a status code is returned, it will be used as the status
    // for the RPC. uri is the URI of the entry point for the request.
    // When supported by the underlying implementation, ctx can be used for
    // timeout and cancellation.
    // TODO(zhaoq): Define the set of the qualified keys instead of leaving
    // it as an arbitrary string.
    GetRequestMetadata(ctx context.Context, uri ...string) (
        map[string]string,    error,
    )
    // RequireTransportSecurity indicates whether the credentials requires
    // transport security.
    RequireTransportSecurity() bool
}
```

Trả về thông tin cần thiết để xác thực trong phương thức `GetRequestMetadata`. Phương thức `RequireTransportSecurity` cho biết kết nối bảo mật ở tầng transport có cần không. Trong thực tế, nên yêu cầu các kết nối có hỗ trợ tầng bảo mật cơ bản này để thông tin xác thực không có nguy cơ bị xâm phạm và giả mạo.

Ta có thể tạo ra kiểu `Authentication` để xác thực username và password:

```go
type Authentication struct {
    User     string
    Password string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (
    map[string]string, error,
) {
    return map[string]string{"user":a.User, "password": a.Password}, nil
}
func (a *Authentication) RequireTransportSecurity() bool {
    return false
}
```

Trong đó phương thức `GetRequestMetadata` trả về thông tin xác thực cục bộ gói cả thông tin đăng nhập và password. Để code được đơn giản hơn nên `RequireTransportSecurity` không cần thiết.

Thông tin token có thể được truyền vào như tham số cho mỗi gRPC service được yêu cầu:

```go
func main() {
    auth := Authentication{
        Login:    "gopher",
        Password: "password",
    }

    conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    ...
}
```

Đối tượng `Authentication` được chuyển đổi thành tham số của `grpc.Dial` bằng hàm `grpc.WithPerRPCCredentials`. Vì secure link không được kích hoạt nên ta cần phải truyền vào `grpc.WithInsecure()` để bỏ qua bước xác thực chứng chỉ bảo mật.

Kế đó trong mỗi phương thức của gRPC server, danh tính người dùng được xác thực bởi phương thức `Authentication` của `Auth`:

```go
type grpcServer struct { auth *Authentication }

func (p *grpcServer) SomeMethod(
    ctx context.Context, in *HelloRequest,
) (*HelloReply, error) {
    if err := p.auth.Auth(ctx); err != nil {
        return nil, err
    }

    return &HelloReply{Message: "Hello " + in.Name}, nil
}

func (a *Authentication) Auth(ctx context.Context) error {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return fmt.Errorf("missing credentials")
    }

    var appid string
    var appkey string

    if val, ok := md["login"]; ok { appid = val[0] }
    if val, ok := md["password"]; ok { appkey = val[0] }

    if appid != a.Login || appkey != a.Password {
        return grpc.Errorf(codes.Unauthenticated, "invalid token")
    }

    return nil
}
```

Công việc xác thực chi tiết chủ yếu được thực hiện trong phương thức `Authentication.Auth`. Đầu tiên, thông tin mô tả (meta infomation) được lấy từ biến ngữ cảnh  `ctx` thông qua `metadata.FromIncomeContext` và sau đó thông tin xác thực tương ứng được lấy ra để xác thực. Nếu xác thực thất bại, nó sẽ trả về lỗi thuộc kiểu `code.Unauthenticated`.

## 4.5.3 Interceptor

`Grpc.UnaryInterceptor` và `grpc.StreamInterceptor` trong gRPC cung cấp hỗ trợ interceptor cho các phương thức thông thường và phương thức stream. Ở đây chúng ta hãy tìm hiểu về việc sử dụng  interceptor cho phương thức thông thường.

Để hiện thực một interceptor như vậy, bạn cần phải hiện thực hàm cho tham số của `grpc.UnaryInterceptor`:

```go
func filter(ctx context.Context,
    req interface{}, info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (resp interface{}, err error) {
    log.Println("fileter:", info)
    return handler(ctx, req)
}
```

Trong đó `ctx` và `req`   là  tham số  của   phương thức RPC bình thường. Tham số `info`  chỉ ra phương thức gRPC tương ứng hiện đang được sử dụng và tham số `handler` tương ứng với hàm phương thức gRPC hiện tại. Dòng đầu của hàm để ghi ra log   tham số `info` sau đó gọi tới phương thức gRPC gắn với `handler`.

Để sử dụng hàm filter interceptor, chỉ cần truyền nó vào lời gọi hàm khi bắt đầu một gRPC service:

```go
server := grpc.NewServer(grpc.UnaryInterceptor(filter))
```

Sau đó server sẽ ghi ra log trước khi nhận lời gọi gRPC, rồi mới gọi tới phương thức được yêu cầu.

Nếu hàm interceptor trả về lỗi thì lệnh gọi phương thức gRPC sẽ được coi là failure. Do đó, chúng ta có thể thực hiện một số công việc xác thực đơn giản trên các tham số đầu và cả kết quả trả về của Interceptor. Interceptor là tính năng rất phù hợp cho việc chứng thực Token đã giới thiệu ở phần trước.

Sau đây là một interceptor có chức năng thêm một exception cho phương thức gRPC:

```go
func filter(
    ctx context.Context, req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (resp interface{}, err error) {
    log.Println("fileter:", info)

    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()

    return handler(ctx, req)
}
```

Tuy nhiên, chỉ một interceptor có thể được gắn cho một service trong gRPC framework, cho nên tất cả chức năng interceptor chỉ có thể thực hiện trong một hàm. Package go-grpc-middleware trong dự án opensource grpc-ecosystem có hiện thực cơ chế hỗ trợ cho một chuỗi interceptor dựa trên gRPC.

Một ví dụ về cách sử dụng chuỗi interceptor trong package go-grpc-middleware:

```go
import "github.com/grpc-ecosystem/go-grpc-middleware"

myServer := grpc.NewServer(
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        filter1, filter2, ...
    )),
    grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
        filter1, filter2, ...
    )),
)
```

Xem chi tiết: [go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)

## 4.5.4 gRPC kết hợp với Web service

gRPC được xây dựng bên trên giao thức HTTP/2 nên chúng ta có thể đặt gRPC service vào các port giống  như một web service bình thường.

Với các service không sử dụng giao thức TLS thì cần phải thực hiện một số tùy chỉnh trong chức năng của HTTP/2:

```go
func main() {
    mux := http.NewServeMux()

    h2Handler := h2c.NewHandler(mux, &http2.Server{})
    server = &http.Server{Addr: ":3999", Handler: h2Handler}
    server.ListenAndServe()
}
```

Cho phép kích hoạt một server https thông thường thì rất đơn giản:

```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        fmt.Fprintln(w, "hello")
    })

    http.ListenAndServeTLS(port, "server.crt", "server.key",
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            mux.ServeHTTP(w, r)
            return
        }),
    )
}
```

Tương tự như vậy kích hoạt một gRPC service với các chứng chỉ riêng cũng rất đơn giản:

```go
func main() {
    creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
    if err != nil {
        log.Fatal(err)
    }

    grpcServer := grpc.NewServer(grpc.Creds(creds))

    ...
}
```

Vì gRPC service đã hiện thực phương thức `ServeHTTP` trước đó nên nó có thể được sử dụng  làm đối tượng xử lý định tuyến Web (routing). Nếu  đặt gRPC và Web service lại với nhau, sẽ dẫn đến xung đột giữa đường dẫn gRPC và Web. Chúng ta cần phân biệt giữa hai loại service này khi xử lý.

Việc tạo ra các handler xử lý việc routing hỗ trợ cả Web và gRPC có thể thực hiện như sau:

```go
func main() {
    ...

    http.ListenAndServeTLS(port, "server.crt", "server.key",
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if r.ProtoMajor != 2 {
                mux.ServeHTTP(w, r)
                return
            }
            if strings.Contains(
                r.Header.Get("Content-Type"), "application/grpc",
            ) {
                grpcServer.ServeHTTP(w, r) // gRPC Server
                return
            }

            mux.ServeHTTP(w, r)
            return
        }),
    )
}
```

Hàm `if` đầu tiên để đảm bảo nếu HTTP không phải là phiên bản  HTTP/2 thì sẽ không hỗ trợ gRPC. Kế tiếp để xét nếu `Content-Type` ở header của request là  *"application/grpc"* thì thực thi lời gọi gRPC tương ứng (ở đây là `ServeHTTP`).

Theo cách này chúng ta có thể cung cấp cả web serive và gRPC chung port cùng một lúc.
