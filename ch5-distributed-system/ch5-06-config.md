# 5.6 Quản lý cấu hình trong hệ thống phân tán

Trong hệ thống phân tán, thường có những vấn đề gây phiền cho chúng ta. Mặc dù, hiện tại có một số cách để thay đổi cấu hình nhưng vẫn bị hạn chế bởi cách hoạt động nội bộ của hệ thống và lúc này sẽ không có cách nào để thay đổi cấu hình một cách thuận tiện nhất. Ví dụ: để giới hạn dòng chảy xuống `downstream`, chúng ta có thể tích luỹ dữ liệu lại và sau khi lượng tích luỹ đến một ngưỡng thời gian hay tổng số tiền thì ta bắt đầu gửi đi, điều này tránh được việc gửi quá nhiều cho `downstream`. Trong trường hợp này, ta lại rất khó để thay đổi cấu hình.

Do đó, mục tiêu của chúng ta là tránh áp dụng hoặc bỏ qua phương pháp trực tuyến và thực hiện một số sửa đổi cho chương trình trực tuyến. Một sửa đổi điển hình là mục cấu hình của chương trình.

## 5.6.1 Thảo luận các ví dụ

### 5.6.1.1 Hệ thống báo cáo (Reporting system)

Trong các hệ thống [OLAP](https://en.wikipedia.org/wiki/Online_analytical_processing) hoặc một số nền tảng dữ liệu ngoại tuyến, sau một thời gian dài phát triển, các chức năng của toàn bộ hệ thống đã dần ổn định. Các dữ liệu đã có sẵn và hầu hết các thay đổi để hiển thị chỉ liên quan tới việc thay đổi câu truy vấn SQL. Lúc này, ta nghĩ tới việc có thể cấu hình được các câu truy vấn SQL mà không cần phải sửa đổi code.

Khi doanh nghiệp đưa ra các yêu cầu mới, việc chúng ta cần làm là cấu hình lại câu SQL cho hệ thống. Những thay đổi này có thể được thực hiện trực tiếp mà không cần khởi động lại.

### 5.6.1.2 Cấu hình mang tính doanh nghiệp

Nền tảng (Platform) của một công ty lớn luôn phục vụ cho nhiều business khác nhau và mỗi business được gán một ID duy nhất. Nền tảng này được tạo thành từ nhiều module và cùng chia sẻ một business. Khi công ty mở một dây chuyền sản phẩm mới, nó cần phải được thông qua bởi tất cả các hệ thống trong nền tảng. Lúc này, chắc chắn là sẽ tốn rất nhiều thời gian để nó có thể chạy được. Ngoài ra, các loại cấu hình toàn cục cần phải được quản lý theo cách thống nhất, các logic cộng và trừ cũng phải được quản lý theo cách thống nhất. Khi cấu hình này được thay đổi, hệ thống cần phải tự động thông báo cho toàn bộ hệ thống của nền tảng mà không cần sự can thiệp của con người (hoặc chỉ can thiệp rất đơn giản, chẳng hạn như kiểm toán nhấp chuột một phát).

Ngoài quản lý trong lĩnh vực kinh doanh, nhiều công ty Internet còn phải kinh doanh theo quy định của thành phố. Khi doanh nghiệp được mở ở một thành phố, ID thành phố mới sẽ tự động được thêm vào danh sách trong hệ thống. Bằng cách này, quá trình kinh doanh có thể chạy tự động.

Một ví dụ khác, có nhiều loại hoạt động trong hệ điều hành của một công ty. Một số hoạt động có thể gặp những sự kiện bất ngờ (như khủng hoảng quan hệ công chúng), và hệ thống cần tắt chức năng liên quan lĩnh vực đó đi. Lúc này, một số công tắc sẽ được sử dụng để tắt nhanh các chức năng tương ứng. Hoặc nhanh chóng xóa ID của hoạt động mà bạn muốn khỏi danh sách chứa. Trong chương Web, chúng ta biết rằng đôi khi cần phải có một hệ thống để đo được lưu lượng truy cập vào các chức năng. Chúng ta có thể chủ động lấy thông tin này kết hợp với cấu hình hệ thống để tắt một tính năng trong trường hợp có lưu lượng lớn bất thường.

## 5.6.2 Sử dụng etcd để thực hiện cập nhật cấu hình

Chúng ta sẽ sử dụng etcd để thực hiện đọc cấu hình và cập nhật tự động để hiểu về một quy trình cập nhật cấu hình trực tuyến.

### 5.6.2.1 Định nghĩa cấu hình

Cấu hình đơn giản, bạn có thể lưu trữ nội dung hoàn toàn trong etcd. Ví dụ:

```shell
etcdctl get /configs/remote_config.json
{
  "addr" : "127.0.0.1:1080",
  "aes_key" : "01B345B7A9ABC00F0123456789ABCDAF",
  "https" : false,
  "secret" : "",
  "private_key_path" : "",
  "cert_file_path" : ""
}
```

### 5.6.2.2 Tạo ứng dụng khách etcd
Cấu trúc khởi tạo kết nối bằng package etcd cho người dùng.

```go
cfg := client.Config{
  Endpoints:               []string{"http://127.0.0.1:2379"},
  Transport:               client.DefaultTransport,
  HeaderTimeoutPerRequest: time.Second,
}
```

### 5.6.2.3 Lấy cấu hình

```go
resp, err = kapi.Get(context.Background(), "/path/to/your/config", nil)
if err != nil {
  log.Fatal(err)
} else {
  log.Printf("Get is done. Metadata is %q\n", resp)
  log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
}
```

Dùng phương thức `Get()` của KeysAPI trong etcd tương đối đơn giản. Các bạn có thể tham khảo thêm các API khác [ở đây](https://godoc.org/github.com/coreos/etcd/client).

### 5.6.2.4 Đăng ký tự động cập nhật cấu hình

```go
kapi := client.NewKeysAPI(c)
w := kapi.Watcher("/path/to/your/config", nil)
go func() {
  for {
    resp, err := w.Next(context.Background())
    log.Println(resp, err)
    log.Println("new values is ", resp.Node.Value)
  }
}()
```

Bằng cách theo dõi những thay đổi sự kiện của đường dẫn cấu hình, khi có nội dung thay đổi trong đường dẫn, chúng ta có thể nhận được thông báo thay đổi cùng với giá trị đã thay đổi.

### 5.6.2.5 Chương trình hoàn chỉnh

```go
package main

import (
  "log"
  "time"

  "golang.org/x/net/context"
  "github.com/coreos/etcd/client"
)

var configPath =  `/configs/remote_config.json`
var kapi client.KeysAPI

type ConfigStruct struct {
  Addr           string `json:"addr"`
  AesKey         string `json:"aes_key"`
  HTTPS          bool   `json:"https"`
  Secret         string `json:"secret"`
  PrivateKeyPath string `json:"private_key_path"`
  CertFilePath   string `json:"cert_file_path"`
}

var appConfig ConfigStruct

func init() {
  cfg := client.Config{
    Endpoints:               []string{"http://127.0.0.1:2379"},
    Transport:               client.DefaultTransport,
    HeaderTimeoutPerRequest: time.Second,
  }

  c, err := client.New(cfg)
  if err != nil {
    log.Fatal(err)
  }
  kapi = client.NewKeysAPI(c)
  initConfig()
}

func watchAndUpdate() {
  w := kapi.Watcher(configPath, nil)
  go func() {
    for {
      resp, err := w.Next(context.Background())
      if err != nil {
        log.Fatal(err)
      }
      log.Println("new values is ", resp.Node.Value)

      err = json.Unmarshal([]byte(resp.Node.Value), &appConfig)
      if err != nil {
        log.Fatal(err)
      }
    }
  }()
}

func initConfig() {
  resp, err = kapi.Get(context.Background(), configPath, nil)
  if err != nil {
    log.Fatal(err)
  }

  err := json.Unmarshal(resp.Node.Value, &appConfig)
  if err != nil {
    log.Fatal(err)
  }
}

func getConfig() ConfigStruct {
  return appConfig
}

func main() {
  // khởi tạo ứng dụng của bạn
}
```

Nếu là doanh nghiệp nhỏ, có thể sử dụng luôn ví dụ trên để hiện thực chức năng mà bạn cần.

Có một vài lưu ý ở đây, chúng ta sẽ làm rất nhiều thứ khi cập nhật cấu hình: phản hồi watch, phân tích json, và các hoạt động này không phải là `atomic`. Khi cấu hình bị thay đổi nhiều lần trong một quy trình dịch vụ, có thể xuất hiện sự không thống nhất logic giữa các yêu cầu xảy ra trước và sau khi cấu hình thay đổi. Do đó, khi bạn sử dụng cách tiếp cận trên để cập nhật cấu hình của mình, bạn cần sử dụng cùng một cấu hình trong suốt vòng đời của một yêu cầu. Cách thức thực hiện cụ thể nên lấy cấu hình một lần khi yêu cầu bắt đầu, và sau đó được truyền tiếp đi cho đến hết vòng đời của yêu cầu.

## 5.6.3 Sự phình to của cấu hình

Khi doanh nghiệp phát triển, áp lực lên hệ thống cấu hình có thể ngày càng lớn hơn và số lượng tệp cấu hình có thể là hàng chục nghìn. Máy khách cũng có hàng chục nghìn và việc lưu trữ nội dung cấu hình bên trong etcd không còn phù hợp nữa. Khi số lượng tệp cấu hình mở rộng, ngoài các vấn đề về thông lượng của hệ thống lưu trữ, còn có các vấn đề về quản lý đối với thông tin cấu hình. Chúng ta cần quản lý các quyền của cấu hình tương ứng và chúng ta cần cấu hình cụm lưu trữ theo lưu lượng truy cập. Nếu có quá nhiều máy khách, khiến hệ thống lưu trữ cấu hình không thể chịu được lượng lớn QPS, thì có thể cần phải thực hiện tối ưu hóa cache ở phía máy khách, ....

Đó là lý do tại sao các công ty lớn luôn phải phát triển một hệ thống cấu hình phức tạp cho doanh nghiệp của họ.

## 5.6.4 Quản lý phiên bản của cấu hình

Trong quy trình quản lý cấu hình, không thể tránh khỏi việc người điều hành thực hiện sai. Ví dụ: khi cập nhật cấu hình, một cấu hình không thể đọc được. Lúc này, chúng ta giải quyết bằng cách kiểm tra toàn bộ.

Đôi khi việc sai cấu hình có thể không phải là do vấn đề với định dạng, mà là vấn đề logic. Ví dụ: khi chúng ta viết SQL, chúng ta chọn ít trường hơn. Khi chúng ta cập nhật cấu hình, chúng ta vô tình làm mất một trường trong chuỗi json và khiến chương trình hiểu cấu hình mới và thực hiện một logic mới. Cách nhanh nhất và hiệu quả nhất để ngăn chặn những lỗi lầm này nhanh chóng là quản lý phiên bản và hỗ trợ khôi phục theo phiên bản.

Khi cấu hình được cập nhật, chúng ta sẽ chỉ định số phiên bản cho từng nội dung của cấu hình, và luôn ghi lại nội dung và số phiên bản trước mỗi lần thay đổi, thực hiện quay ngược bản trước khi phát hiện sự cố với cấu hình mới.

Một cách phổ biến trong thực tế là sử dụng MySQL để lưu trữ các phiên bản khác nhau của tệp cấu hình hoặc chuỗi cấu hình. Khi bạn cần quay lại, chỉ cần thực hiện một truy vấn đơn giản.

## 5.6.5 Khả năng chịu lỗi ở máy khách

Sau khi cấu hình của hệ thống kinh doanh được chuyển đến trung tâm cấu hình, điều đó không có nghĩa là hệ thống của chúng ta đã hoàn thành nhiệm vụ. Khi trung tâm cấu hình ngừng hoạt động, chúng ta cũng cần một vài cơ chế chịu lỗi, ít nhất là để đảm bảo rằng doanh nghiệp vẫn có thể hoạt động trong thời gian này. Điều này đòi hỏi hệ thống phải lấy đủ thông tin cấu hình cần thiết trước khi trung tâm cấu hình ngừng hoạt động. Ngay cả khi thông tin này không đủ mới.

Cụ thể, khi cung cấp SDK đọc cấu hình cho một dịch vụ, tốt nhất là lưu cache cấu hình thu được trên đĩa của máy nghiệp vụ. Khi trung tâm cấu hình không hoạt động, bạn có thể trực tiếp sử dụng nội dung của đĩa cứng. Khi kết nối lại được với trung tâm cấu hình, các nội dung sẽ được cập nhật.

Hãy xem xét kỹ vấn đề thống nhất dữ liệu khi thêm cache. Các máy kinh doanh có thể không thống nhất về cấu hình do lỗi mạng, chúng ta có thể biết được nó đang diễn ra bằng hệ thống giám sát.

Chúng ta sử dụng một cách để giải quyết các vấn đề của việc cập nhật cấu hình, nhưng đồng thời chúng ta lại mang đến những vấn đề mới bằng việc sử dụng cách đó. Trong thực tế, chúng ta phải suy nghĩ rất nhiều về từng quyết định để chúng ta không bị thiệt hại quá nhiều khi vấn đề xảy ra.

[Tiếp theo](ch5-07-crawler.md)