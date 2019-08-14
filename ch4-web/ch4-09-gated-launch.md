# 4.9 Grayscale Publishing và kiểm định A/B

Những công ty có tầm cỡ trung bình thường cung cấp dịch vụ cho hàng triệu users, trong khi hệ thống của các công ty lớn sẽ phải phục vụ cho chục triệu, thậm chí hàng tỉ users. Đầu vào của những requests từ các hệ thống lớn thường bất tận, và bất cứ thay đổi nào cũng sẽ được cảm nhận từ người dùng cuối. Ví dụ, nếu hệ thống của bạn từ chối một số upstream requests trên đường xử lý, nguyên nhân đến từ hệ thống của bạn và nó không có tính chịu lỗi, sau đó lỗi sẽ được ném ra đến người dùng cuối. Tạo thành một thiệt hại thực sự đến user, loại của thiệt hại này sẽ hiển thị một pops up chứa những thông điệp lạ lẫm trên làm màn hình của user app. Do users không thể biết nó là gì, users có thể quên nó đi, bằng việc refreshing lại trang. Nhưng nó cũng làm users mất đi cơ hội mua một món hàng nào đó vì có hàng chục hoặc hàng ngàn người mua khác cùng một thời điểm, bởi vì những vấn đề nhỏ như vậy trong mã nguồn, làm mất đi lợi thế ban đầu, và mất đi cơ hội mua được món hàng yêu thích mà họ đã chờ đợi trong vài tháng. Mức độ thiệt hại mà users phải gánh chịu sẽ phụ thuộc vào tầm quan trọng của hệ thống của bạn đối với users.

Trong trường hợp này, tính chất fault tolerance (tính chịu lỗi) cực kỳ quan trọng trong hệ thống lớn. Mặc dù ngày nay, những công ty làm về internet thường nói rằng họ đã kiểm tra nghiêm ngặt và triệt để trước khi đưa vào khai thác, cho dù họ làm vậy chăng nữa, code bugs vẫn không thể tránh khỏi. Ngay cả khi mã nguồn không có bugs, sự giao tiếp giữa các services trong hệ thống phân tán cũng có thể  gây ra lỗi.

Vào thời điểm này, "Grayscale release" cực kì quan trọng. Grayscale release cũng thường được được gọi là "Canary release" (canary: chim hoàng yến). Vào thế kỉ 17, những người công nhân hầm mỏ ở Anh đã phát hiện ra rằng chim hoàng yến rất nhạy cảm với khí gas. Khi khí gas đạt tới một nồng độ nào đó, thì chim hoàng yến sẽ chết, nhưng ở nồng độ gas làm cho chim hoàng yến chết lại không gây hại cho con người, do đó chim hoàng yến được dùng để làm công cụ phát hiện khí gas. Grayscale publishing của một hệ thống internet thông thường đạt được thông qua hai cách:

1. Hiện thực Grayscale publishing thông qua batch deployment.
2. Grayscale publishing thông qua business rules.

Phương pháp đầu tiên được sử dụng nhiều trong các hàm cũ của hệ thống. Khi mà một hàm mới được đưa vào hoạt động, thì phương pháp thứ hai sẽ được dùng nhiều hơn. Dĩ nhiên, khi gây ra một số thay đổi chính đến những hàm cũ mà chúng quan trọng, thì thông thường sẽ tốt hơn nếu publish chúng theo business rules, bởi vì độ rủi ro khi mở tất cả các hàm cho người dùng là khá lớn.

## 4.9.1 Hiện thực grayscale publishing bằng cách deployment theo nhóm

Nếu service được deploy trên 15 instanses (có thể là physical machines hoặc containers), chúng ta chia 15 instances thành nhóm theo thứ tự độ ưu tiên, sẽ có 1-2-4-8 machines, mỗi thời điểm. Khi mở rộng ra, số lượng tăng gấp đôi.

<div align="center">
	<img src="../images/ch5-online-group.png">
	<br/>
	<span align="center">
		<i>Group deployment</i>
	</span>
</div>
<br/>

Tại sao lại gấp đôi? Nó sẽ đảm bảo rằng chúng ta không chia nhỏ group quá nhiều, không quan trọng bao nhiêu machines mà chúng ta có. Cho ví dụ, 1024 machines, trong thực tế, chỉ cần 1-2-4-8-16-32-64-128-256-512 là 10 lần deployment để toàn bộ được deployed.

Theo cách đó, những users đầu tiên bị ảnh hưởng bởi sự thay đổi sẽ chiếm một phần nhỏ trong tổng số users, như là service của 1000 machines. Nếu không có vấn đề nào sau khi chúng ta deployed toàn bộ và đi vào hoạt động, nó sẽ chỉ ảnh hưởng tới 1/1000 users. Nếu 10 groups hoàn toàn được chia ra đều nhau, thì sẽ ảnh hưởng tới 1/10 users ngay lặp tức, và 1/10 của business problems, nó sẽ là một tai họa không thể khắc phục của công ty.

Khi đi vào hoạt động, cách hiệu quả nhất để quan sát là nhìn vào `error log` của chương trình. Nếu không có nhiều lỗi logic, thì chúng ta sẽ scroll nhanh error log khi xem. Những errors có thể được báo cáo trong hệ thống monitoring của công ty như là `metrics`. Do đó, trong suốt quá trình đi vào hoạt động, có thể thấy rằng bất cứ lỗi bất thường nào xảy ra sẽ được monitoring.

Nếu có một trường hợp bất thường, việc làm đầu tiên là roll back.

## 4.9.2 Grayscale publishing thông qua business rules

Có nhiều chiến lược Grayscale phổ biến. Ví dụ, chiến lược của chúng ta là publish trong hàng ngàn points. Sau đó chúng ta có thể dùng user id, mobile phone number, user device information, v,v để sinh ra một giá trị hash.

```go
// pass 3/1000
func passed() bool {
    key := hashFunctions(userID) % 1000
    if key <= 2 {
        return true
    }

    return false
}
```

## 4.9.2.1 Các rules tuỳ chọn

Một số hệ thống Grayscale publishing phổ biến sẽ có một số rules chọn từ:

```
1. Published by city
2. Publish by probability
3. Published by percentage
4. Publish by whitelist
5. Published by line of business
6. Publish by UA (APP, Web, PC)
7. Publish by distribution channel
```

Publishing bởi whitelist thì tương đối đơn giản. Khi tính năng đưa vào hoạt động, chúng ta hy vọng rằng chỉ những người nhân viên là testers trong công ty có thể truy cập các tính năng đó. Họ sẽ trực tiếp cho các accounts và mailboxs vào whitelist và từ chối truy cập đến các accounts khác.

Publishing theo xác suất được hiện thực bởi một function đơn giản:

```go
func isTrue() bool {
    //return true/false according to the rate provided by user
}
```

Có thể thấy kết quả trả về theo xác suất được ghi nhận bởi users là `true` hoặc `false`. Dĩ nhiên xác suất của `true` và `false` trong mã nguồn trên có thể là `100% true` và `0% false` hoặc `0% true` và `100% false`. Function này không yêu cầu bất cứ input nào.

Publishing theo phần trăm nghĩa là sẽ hiện thực một function như sau:

```go
func isTrue(phone string) bool {
    if hash of phone matches {
        return true
    }

    return false
}
```

Trường hợp này có thể trả về kết quả `true` hoặc `false` theo tỷ lệ phần trăm được đặc tả trước, và ở trên là sự khác biệt đơn giản theo xác suất mà chúng ta cần người gọi cung cấp một tham số input parameter. Chúng ta sử dụng input parameter như là một thông số để tính toán giá trị hash, sau đó trả về kết quả là một model. Điều này đảm bảo rằng user sẽ trả về cùng một kết quả qua nhiều lần gọi, trong ngữ cảnh sau, thuật toán sẽ phân đoạn được kết quả mong đợi

<div align="center">
	<img src="../images/ch5-set-time-line.png">
	<br/>
	<span align="center">
		<i>First set and then get immediately</i>
	</span>
</div>
<br/>

Nếu bạn dùng chiến lược random, bạn sẽ gặp một vấn đề như hình 5-22

<div align="center">
	<img src="../images/ch5-set-time-line_2.png">
	<br/>
	<span align="center">
		<i>First set and then get immediately</i>
	</span>
</div>
<br/>

## 4.9.3 Làm thế nào để hiện thực hệ thống Grayscale publishing

### 4.9.3.1 Business-related simple grayscale

Công ty thông thường sẽ có một bảng ánh xạ giữa tên của thành phố và `ids`. Nếu business chỉ trong vòng một quốc gia, số thành phố sẽ không thực sự lớn, và `ids` có thể trong vòng `10,000`. Chúng ta chỉ cần một mảng `bool` có kích thước khoảng `10,000` đạt được nhu cầu.

```go
var cityID2Open = [12000]bool{}

func init() {
    readConfig()
    for i:=0;i<len(cityID2Open);i++ {
        if city i is opened in configs {
            cityID2Open[i] = true
        }
    }
}

func isPassed(cityID int) bool {
    return cityID2Open[cityID]
}
```

Nếu công ty sử dụng giá trị lớn lơn cho cityID, thì chúng ta có thể cân nhắc sử dụng map để lưu trữ giá trị. Câu truy vấn map sẽ thường chậm hơn array, nhưng việc mở rộng sẽ linh hoạt hơn

```go
var cityID2Open = map[int]struct{}{}

func init() {
    readConfig()
    for _, city := range openCities {
        cityID2Open[city] = struct{}{}
    }
}

func isPassed(cityID int) bool {
    if _, ok := cityID2Open[cityID]; ok {
        return true
    }

    return false
}
```

Publishing bởi probability (xác suất) đặc biệt hơn tí, nhưng rất dễ để hiện thực mà không cần tới input.

```go
func init() {
    rand.Seed(time.Now().UnixNano())
}

// rate từ 0 tới 100
func isPassed(rate int) bool {
    if rate >= 100 {
        return true
    }

    if rate > 0 && rand.Int(100) > rate {
        return true
    }

    return false
}
```

Chú ý tới khởi tạo `seed`.

### 4.9.3.2 Thuật toán Hash

Có nhiều thuật thoán hash như là `md5`, `crc32`, `sha1`, v,v,.. nhưng mục đích mà chúng ta hướng đến là ánh xạ những data tới key tương ứng, và ta không muốn sử dụng quá nhiều CPU cho việc tính toán hash. Đa số các thuật toán đều `murmurhash`, sau đây là kết quả benchmark cho những thuật toán hash phổ biến đó.

Sau khi dùng thư viện chuẩn `md5`, `sha1` và opensource hiện thực `murmur3` cho việc so sánh

```go
package main

import (
    "crypto/md5"
    "crypto/sha1"

    "github.com/spaolacci/murmur3"
)

var str = "hello world"

func md5Hash() [16]byte {
    return md5.Sum([]byte(str))
}

func sha1Hash() [20]byte {
    return sha1.Sum([]byte(str))
}

func murmur32() uint32 {
    return murmur3.Sum32([]byte(str))
}

func murmur64() uint64 {
    return murmur3.Sum64([]byte(str))
}
```

Viết benchmark test cho các thuật toán đó:

```go
package main

import "testing"

func BenchmarkMD5(b *testing.B) {
    for i := 0; i < b.N; i++ {
        md5Hash()
    }
}

func BenchmarkSHA1(b *testing.B) {
    for i := 0; i < b.N; i++ {
        sha1Hash()
    }
}

func BenchmarkMurmurHash32(b *testing.B) {
    for i := 0; i < b.N; i++ {
        murmur32()
    }
}

func BenchmarkMurmurHash64(b *testing.B) {
    for i := 0; i < b.N; i++ {
        murmur64()
    }
}
```

Sau đó xem kết quả chạy như sau:

```sh
$ go test -bench=.
goos: darwin
goarch: amd64
BenchmarkMD5-4          10000000 180 ns/op
BenchmarkSHA1-4         10000000 211 ns/op
BenchmarkMurmurHash32-4 50000000  25.7 ns/op
BenchmarkMurmurHash64-4 20000000  66.2 ns/op
PASS
ok _/Users/caochunhui/test/go/hash_bench 7.050s
```

Có thể thấy rằng **murmurhash** có hiệu suất cao gấp ba so với các hàm thuật toán hash khác. Hiển nhiên, để thực hiện việc `load balancing` (cân bằng tải), murmurhash sẽ tốt hơn `md5` và `sha1`. Thực tế, đó là thuật toán hash hiệu quả trong cộng đồng vài năm vừa qua, người đọc có thể tự nghiên cứu.

### 4.9.3.3 Liệu có một mô hình phân phối chuẩn không ?

Cho một thuật toán hash, xem xét vấn đề hiệu suất, sẽ thực sự cần thiết khi xem xét khi nào giá trị hash sẽ theo phân phối chuẩn. Nếu giá trị sau khi hash không theo phân phối chuẩn, thì sẽ không đạt được hiệu ứng `uniform gray`.

Xét `murmur3` là một ví dụ, hãy bắt đầu với `15810000000`, sinh ra mười triệu con số di động tương tự nhau, sau đó tính toán giá trị hash và chia vào mười `buckets` và quan sát khi nào giá trị đều nhau

```go
package main

import (
    "fmt"

    "github.com/spaolacci/murmur3"
)

var bucketSize = 10

func main() {
    var bucketMap = map[uint64]int{}
    for i := 15000000000; i < 15000000000+10000000; i++ {
        hashInt := murmur64(fmt.Sprint(i)) % uint64(bucketSize)
        bucketMap[hashInt]++
    }
    fmt.Println(bucketMap)
}

func murmur64(p string) uint64 {
    return murmur3.Sum64([]byte(p))
}
```


Hãy xem kết quả thực thi

```go
map[7:999475 5:1000359 1:999945 6:1000200 3:1000193 9:1000765 2:1000044 \
4:1000343 8:1000823 0:997853]
```

Độ sai lệch trong vòng 1/100 và có thể chấp nhận được. Khi độc giả đối chiếu với các thuật toán khác và đánh giá cái nào sẽ được dùng cho Grayscale publishing, có thể xem xét tới hiệu suất và tính cân bằng trong chương này.
