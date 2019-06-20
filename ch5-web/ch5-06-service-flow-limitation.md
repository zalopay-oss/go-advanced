# 5.6 Ratelimit Service Flow Limit

Chương trình máy tính có thể được phân loại theo kiểu thắt cổ chai disk IO:

- Thắt cổ chai do CPU tính toán
- Thắt cổ chai do băng thông mạng
- Và đôi khi hệ thống bên ngoài gây ra tình trạng thắt cổ chai trong hệ thống phân tán.

Phần quan trọng nhất của hệ thống Web là mạng. Cho dù đó là tiếp nhận, phân tích request của người dùng, truy cập bộ nhớ hay trả về dữ liệu response đều cần phải truy cập trực tuyến. Trước IO multiplexing interface `epoll/kqueue` do hệ thống cung cấp đã từng có một sự cố C10k trong máy tính hiện đại đa lõi. Vấn đề C10k có thể khiến máy tính không thể sử dụng toàn bộ CPU để xử lý nhiều kết nối người dùng. Do đó cần phải chú ý tối ưu hóa chương trình để tăng mức sử dụng CPU.

Kể từ khi trên Linux có `epoll`, FreeBSD hiện thực `kqueue`, chúng ta có thể dễ dàng giải quyết vấn đề C10k với API do kernel cung cấp. Điều đó có nghĩa là nếu chương trình của chúng ta chủ yếu xử lý qua mạng, thì thắt cổ chai phải nằm phía người dùng, chứ không nằm ở kernel của hệ điều hành.

Ngày nay, việc phát triển ở lớp application hầu như không thể thấy trong chương trình. Hầu hết, chúng ta chỉ cần tập trung vào logic nghiệp vụ. Thư viện `net` của Go đóng gói các syscall API khác nhau cho các nền tảng khác nhau. Thư viện `http` được xây dựng trên nền của thư viện `net`, vì vậy trong Golang, chúng ta có thể viết các service `http` hiệu suất cao với sự trợ giúp của thư viện chuẩn. Đây là một service `hello world` đơn giản:

```go
package main

import (
    "io"
    "log"
    "net/http"
)

func sayhello(wr http.ResponseWriter, r *http.Request) {
    wr.WriteHeader(200)
    io.WriteString(wr, "hello world")
}

func main() {
    http.HandleFunc("/", sayhello)
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}
```

Chúng ta cần đo throughput của Web service này, đồng thời đó cũng là QPS của interface. Sử dụng [wrk](https://github.com/wg/wrk) trên máy tính cá nhân có cấu hình như sau:

```sh
ThinkPad-T470
-------------------------------
OS: Ubuntu 18.04.2 LTS x86_64
Host: 20HES39900 ThinkPad T470
Kernel: 4.15.0-52-generic

Resolution: 1920x1080, 1920x1080

CPU: Intel i5-7200U (4) @ 3.100GHz
GPU: Intel HD Graphics 620
Memory: 2977MiB / 15708MiB
```

Kết quả test:

```sh
➜ wrk -c 10 -d 10s -t10 http://localhost:9090
Running 10s test @ http://localhost:9090
  10 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   360.57us  769.72us  14.27ms   91.45%
    Req/Sec     7.07k     0.91k    9.25k    76.10%
  704529 requests in 10.01s, 86.00MB read
Requests/sec:  70363.89
Transfer/sec:      8.59MB

~ took 10s
➜ wrk -c 10 -d 10s -t10 http://localhost:9090
Running 10s test @ http://localhost:9090
  10 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   395.85us    0.90ms  18.93ms   91.88%
    Req/Sec     6.92k     0.95k    9.34k    75.30%
  688941 requests in 10.02s, 84.10MB read
Requests/sec:  68783.96
Transfer/sec:      8.40MB

~ took 10s
➜ wrk -c 10 -d 10s -t10 http://localhost:9090
Running 10s test @ http://localhost:9090
  10 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   374.92us  814.34us  13.76ms   91.47%
    Req/Sec     7.09k     1.03k   21.18k    82.82%
  706968 requests in 10.09s, 86.30MB read
Requests/sec:  70040.20
Transfer/sec:      8.55MB
```

Kết quả của nhiều thử nghiệm là khoảng 70.000 QPS và thời gian phản hồi là khoảng 15ms. Đối với một ứng dụng web, đây đã là một kết quả rất tốt. Đây chỉ là một máy tính cá nhân nên với server có cấu hình cao hơn sẽ đạt được kết quả còn tốt hơn nữa.

Một số hệ thống thường bị thắt cổ chai do mạng, chẳng hạn như dịch vụ CDN và dịch vụ Proxy. Một số chương trình thường thắt cổ chai do CPU / GPU, như các dịch vụ xác minh đăng nhập và dịch vụ xử lý hình ảnh. Một số lại do truy cập disk như hệ thống lưu trữ, cơ sở dữ liệu. Các nút thắt chương trình khác nhau được phản ánh ở những nơi khác nhau và các dịch vụ đơn chức năng được đề cập ở trên tương đối dễ phân tích. Nếu bạn gặp một mô-đun có số lượng logic nghiệp vụ lớn và số lượng code lớn, thì vấn đề thắt cổ chai có thể phải trả qua nhiều lần stress test mới có thể phát hiện ra.

Đối với lớp nút cổ chai IO / Network, hiệu suất là IO / đĩa IO sẽ đầy trước CPU. Trong trường hợp này, ngay cả khi CPU được tối ưu hóa, thông lượng của toàn bộ hệ thống có thể được cải thiện và tốc độ đọc và ghi của đĩa có thể tăng lên. Kích thước bộ nhớ làm tăng băng thông của NIC để cải thiện hiệu suất tổng thể. Chương trình thắt cổ chai CPU là mức sử dụng CPU đạt 100% trước khi bộ nhớ và card mạng không đầy. CPU bận rộn với nhiều tác vụ điện toán khác nhau và các thiết bị IO tương đối nhàn rỗi.



