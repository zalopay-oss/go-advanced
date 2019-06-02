# 6.1 Distributed id generator

Đôi khi chúng ta cần có tạo ra một id tương tự như ID tăng tự động của MySQL và nó không được lặp lại. Để hỗ trợ các ngữ cảnh đồng thời nghiêm ngặc trong kinh doanh. Điển hình, khi có chương trình khuyến mãi trong thương mại điện tử, một số lượng lớn đơn đặt hàng sẽ tràn vào hệ thống trong một khoảng thời gian ngắn, chẳng hạn như 10w+ mỗi giây. Khi ngôi sao điện ảnh gặp vấn đề, sẽ có rất nhiều người hâm mộ nhiệt tình gửi microblog để bày tỏ ý kiến của họ, và tạo ra rất nhiều tin tức trong một thời gian ngắn.

Trước khi chèn vào cơ sở dữ liệu, chúng ta cần cung cấp cho các thông tin đó một ID trước, sau đó chèn nó vào cơ sở dữ liệu. Yêu cầu đối với id này là nó có thông tin về thời gian. Do đó, cơ sở dữ liệu phụ của chúng ta vẫn có thể sắp xếp các tin nhắn đó theo thứ tự thời gian.

Thuật toán snowflake của Twitter là giải pháp điển hình trong ngữ cảnh ngày. Hãy nhìn vào hình 6-1:

![ch6-1-snowflake](../images/ch6-1-snowflake.png)

*Hình 6-1 Phân phối Bit trong snowflake*

Đầu tiên, ta xác định rằng giá trị là 64 bit, loại int64, được chia thành bốn phần:

- Không sử dụng bit đầu tiên vì bit này là bit dấu. 
- `timestamp`: sử dụng 41 bit để biểu thị timestamp khi nhận được yêu cầu, tính bằng mili giây
- `datacenter_id`: 5 chữ số để biểu thị id của trung tâm dữ liệu
- `worker_id`: 5 chữ số để biểu thị id của máy host
- `sequence_id`: cuối cùng là id tăng vòng lặp 12 bit ( tăng từ 0 đến 111111111111 rồi trở về 0).

Cơ chế này có thể hỗ trợ chúng ta tạo ra `2 ^ 12 = 4096` tin nhắn trong cùng một mili giây của cùng một máy. Tổng cộng có 4.096 triệu tin nhắn trong một giây. Nó là hoàn toàn đủ để sử dụng.

Trung tâm dữ liệu cộng với id đối tượng có tổng cộng 10 chữ số, hỗ trợ 32 máy cho mỗi trung tâm dữ liệu và 1024 trường hợp trong tất cả các trung tâm dữ liệu.

`timestamp` gồm 41 chữ số nên có thể hỗ trợ chúng ta trong 69 năm. Tất nhiên, thời gian tính bằng mili giây của chúng ta sẽ không thực sự bắt đầu từ năm 1970, vì nếu như vậy thì hệ thống sẽ chỉ hoạt động đến `2039/9/7 23: 47: 35` , do đó, chúng chỉ nên tăng so với một mốc thời gian nhất định. Ví dụ: hệ thống của chúng ta là bắt đầu chạy vào 2018-08-01, thì chúng ta có thể coi mốc thời gian này là `2018-08-01 00: 00: 00.000` và tăng lên từ đó.

## 6.1.1 phân công worker_id

`timestamp`, `datacenter_id`, `worker_id` và `sequence_id` là bốn trường, riêng timestamp và sequence_id được tạo bởi chương trình khi chạy. Còn datacenter_id và worker_id cần lấy trong giai đoạn triển khai và một khi chương trình được khởi động, nó không thể thay đổi (nếu bạn có thể thay đổi nó theo ý muốn, nó có thể vô tình bị sửa đổi, gây ra xung đột cho giá trị id cuối cùng).

Generally, machines in different data centers will provide corresponding APIs for obtaining data center ids, so datacenter_idwe can easily obtain them during the deployment phase. And worker_id is an id that we logically assign to the machine. What should we do? The simpler idea is supported by tools that provide this self-increasing id feature, such as MySQL:

Nhìn chung, các máy trong các trung tâm dữ liệu khác nhau sẽ cung cấp các API tương ứng để lấy id của trung tâm dữ liệu, vì vậy datacenter_id có thể dễ dàng lấy trong giai đoạn triển khai. Còn worker_id là một id mà chúng ta gán một cách logically cho từng máy. Chúng ta nên xử lý thế nào? Ý tưởng đơn giản được hỗ trợ bởi các công cụ là cung cấp tính năng id tự tăng, chẳng hạn như trong MySQL:

mysql> insert into a (ip) values("10.1.2.101");
Query OK, 1 row affected (0.00 sec)

```sh
mysql> select last_insert_id();
+------------------+
| last_insert_id() |
+------------------+
|                2 |
+------------------+
1 row in set (0.00 sec)
```

Khi bạn đã nhận được nó từ MySQL, worker_id sẽ được duy trì cục bộ, tránh được trường hợp làm mới mỗi khi bạn truy cập worker_id. Để các worker_id duy nhất luôn giữ nguyên.

Tất nhiên, sử dụng MySQL đồng nghĩa với việc thêm một phụ thuộc bên ngoài vào dịch vụ tạo id của chúng ta. Càng thêm nhiều phụ thuộc, khả năng phục vụ của dịch vụ càng tệ.

Considering that even if there is a single id generation service instance hanging in the cluster, it is a part of the id that is lost for a while, so we can also be more violent and worker_idwrite it directly in the worker configuration. When going online, the deployment script completes the worker_idfield. replace.

Xem xét đến việc có một dịch vụ tạo id bị lỗi trong cluster, một phần của id bị mất trong một khoảng thời gian, thì chúng ta cần một cách bạo lực là viết trực tiếp worker_id vào cấu hình worker. Khi lên mạng, tập lệnh triển khai sẽ hoàn thành trường worker_id.

## 6.1.2 Các nguồn mở

### 6.1.2.1 Standard snowflake implementation Hiện thực theo snowflake chuẩn

github.com/bwmarrin/snowflake It's a fairly lightweight implementation of Snowflake's Go. Its document definition you use, see Figure 6-2 below.

Đây là một hiện thực Snowflake's Go khá nhẹ. Bạn có thể sử dụng tài liệu định nghĩa của nó, xem Hình 6-2 dưới đây.

![](../images/ch6-2-snowflake-easy.png)

Hình 6-2 thư viện snowflake

It is exactly the same as the standard snowflake. It is relatively simple to use:

Nó giống hoàn toàn như snowflake tiêu chuẩn. Và tương đối đơn giản để sử dụng:

```go

package main

import (
    "fmt"
    "os"

    "github.com/bwmarrin/snowflake"
)

func main() {
    n, err := snowflake.NewNode(1)
    if err != nil {
        println(err)
        os.Exit(1)
    }

    for i := 0; i < 3; i++ {
        id := n.Generate()
        fmt.Println("id", id)
        fmt.Println(
            "node: ", id.Node(),
            "step: ", id.Step(),
            "time: ", id.Time(),
            "\n",
        )
    }
}

```

Of course, this library also left us with a customized backend, which reserved some customizable fields:

Dĩ nhiên, thư viện này cũng cho phép chúng ta tùy chỉnh, một số trường có thể tùy chỉnh như:

```go
    // Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC
    // You may customize this to set a different epoch for your application.
    Epoch int64 = 1288834974657

    // Số bit để sử dụng cho Node
    // Nhớ rằng, bạn có tổng 22 bits chia sẻ giữa Node/Step
    NodeBits uint8 = 10

    // Số bit để sử dụng cho Step
    // Nhớ rằng, bạn có tổng 22 bits chia sẻ giữa Node/Step
    StepBits uint8 = 12
```

`Epoch` It is the start time at the beginning of this section, NodeBitswhich refers to the bit length of the machine number, StepBitswhich refers to the bit length of the self-incrementing sequence.

`Epoch` là thời gian bắt đầu, `NodeBits` dùng để chỉ độ dài bit của máy chủ đang chaỵ và `StepBits` đề cập đến độ dài bit của chuỗi tự tăng.

### 6.1.2.2 sonyflake

Sonyflake is an open source project of Sony. The basic idea is similar to snowflake, but the bit allocation is slightly different. See *Figure 6-3*:

Sonyflower là một dự án nguồn mở của Sony. Ý tưởng cơ bản tương tự như snowflake, nhưng phân bổ bit hơi khác. Xem * Hình 6-3 *:

![sonyflake](../images/ch6-snoyflake.png)

*Hình 6-3 sonyflake*

The time here only uses 39 bits, but the unit of time becomes 10ms, so theoretically it is longer than the time indicated by 41 bits (174 years).

Thời gian ở đây chỉ sử dụng 39 bit, nhưng đơn vị thời gian trở thành 10ms. Vì vậy, về mặt lý thuyết, nó dài hơn thời gian được chỉ định bởi 41 bit (174 năm).

`Sequence ID` is consistent with the previous definition, `Machine ID` is actually the node id. The difference between `sonyflake` is its configuration parameters during the startup phase:

`Sequence ID` giống với với định nghĩa trước đó, `Machine ID` thực sự là id của Node. Sự khác biệt giữa `sonyflower` là các tham số cấu hình của nó trong giai đoạn khởi động:

```go
Func NewSonyflake(st Settings) *Sonyflake
```

Cấu trúc của `Settings` như sau:

```go

Type Settings struct {
    StartTime time.Time
    MachineID func() (uint16, error)
    CheckMachineID func(uint16) bool
}

```

Tùy chọn `StartTime` tương tự như `Epoch` trước đây của chúng ta. Nếu không được đặt, mặc định là bắt đầu vào `2014-09-01 00:00:00 +0000 UTC`.

`MachineID` can be a user-defined function. If the user does not define it, the lower 16 bits of the native IP will be used as the `machine id` by default.

`MachineID` có thể là hàm do người dùng định nghĩa. Nếu người dùng không định nghĩa nó, 16 bit cuối của IP gốc sẽ mặc định được sử dụng làm `machine id`.

`CheckMachineID` is a function provided by the user to check if `MachineID` conflicts. The design here is quite clever. If there is another centralized storage and support for checking the duplicate storage, then we can customize the logic to check if the `MachineID` conflicts according to our own ideas. If the company has a ready-made Redis cluster, then we can easily check for conflicts with Redis' collection types.

`CheckMachineID` là một chức năng được cung cấp bởi người dùng để kiểm tra xem `MachineID` có trùng hay không. Đây là thiết kế khá thông minh. Nếu có một bộ lưu trữ tập trung và hỗ trợ kiểm tra sự trùng lặp, thì chúng ta có thể tùy chỉnh logic để kiểm tra xem `MachineID` có trùng không. Nếu công ty có cụm Redis được tạo sẵn, thì chúnv
```shell
Redis 127.0.0.1:6379> SADD base64_encoding_of_last16bits MzI0Mgo=
(integer) 1
Redis 127.0.0.1:6379> SADD base64_encoding_of_last16bits MzI0Mgo=
(integer) 0
```

Nó cũng khá đơn giản để sử dụng:
<!-- và một số hàm với logic đơn giản bị bỏ qua: -->

```go

Package main

Import (
   "fmt"
   "os"
   "time"

   "github.com/sony/sonyflake"
)

Func getMachineID() (uint16, error) {
   Var machineID uint16
   Var err error
   machineID = readMachineIDFromLocalFile()
   If machineID == 0 {
      machineID, err = generateMachineID()
      If err != nil {
         Return 0, err
      }
   }

   Return machineID, nil
}

Func checkMachineID(machineID uint16) bool {
   saddResult, err := saddMachineIDToRedisSet()
   If err != nil || saddResult == 0 {
      Return true
   }

   Err := saveMachineIDToLocalFile(machineID)
   If err != nil {
      Return true
   }

   Return false
}

Func main() {
   t, _ := time.Parse("2006-01-02", "2018-01-01")
   Settings := sonyflake.Settings{
      StartTime: t,
      MachineID: getMachineID,
      CheckMachineID: checkMachineID,
   }

   Sf := sonyflake.NewSonyflake(settings)
   Id, err := sf.NextID()
   If err != nil {
      fmt.Println(err)
      os.Exit(1)
   }

   fmt.Println(id)
}

```
