# 5.5 Cân bằng tải (Loadbalancer)
<div align="center">
	<img src="../images/ch6-loadbalancer.png" width="400">
	<br/>
	<span align="center">
		<i>Loadbalancer</i>
	</span>
</div>

Phần này sẽ thảo luận về các phương pháp phổ biến trong cân bằng tải hệ thống phân tán.

## 5.5.1 Ý tưởng cân bằng tải

Khi có n node cùng cung cấp service và chúng ta cần chọn một trong số đó để thực hiện quy trình business. Có một số ý tưởng:

1. Chọn theo thứ tự: lần gần nhất bạn chọn cái đầu tiên, thì lần này bạn chọn cái thứ hai, rồi cứ thế với cái tiếp theo. Nếu bạn đã đạt đến cái cuối cùng, thì cái tiếp theo bắt đầu từ cái đầu tiên. Trong trường hợp này, chúng ta có thể lưu trữ thông tin node dịch vụ trong một mảng. Sau khi mỗi yêu cầu được hoàn thành xuôi dòng, chúng ta di chuyển chỉ mục đi tiếp. Di chuyển trở lại đầu của mảng khi bạn di chuyển đến cuối.

2. Chọn ngẫu nhiên: Chọn node một cách ngẫu nhiên. Giả sử rằng máy thứ x được chọn, thì x có thể được chọn từ hàm `rand.Intn()%n`.

3. Sắp xếp các node theo một trọng lượng nhất định và chọn một node có trọng lượng lớn nhất hoặc nhỏ nhất.

Nếu yêu cầu không thành công, chúng ta vẫn cần cơ chế để thử lại. Đối với thuật toán ngẫu nhiên, có khả năng bạn sẽ chọn node lỗi lần nữa.

## 5.5.2 Cân bằng tải dựa trên thuật toán xáo trộn

Giả sử chúng ta cần chọn ngẫu nhiên node gửi yêu cầu và thử lại các node khác khi có lỗi trả về. Vì vậy, chúng ta thiết kế một mảng chỉ mục với kích thước bằng số node. Mỗi lần chúng ta có một yêu cầu mới, chúng ta xáo trộn mảng chỉ mục, sau đó lấy phần tử đầu tiên làm node dịch vụ. Nếu yêu cầu thất bại, ta chọn node tiếp theo. Cứ thử lại và tiếp tục, ...

```go
var endpoints = []string {
    "100.69.62.1:3232",
    "100.69.62.32:3232",
    "100.69.62.42:3232",
    "100.69.62.81:3232",
    "100.69.62.11:3232",
    "100.69.62.113:3232",
    "100.69.62.101:3232",
}

// shuffle hàm xáo trộn chỉ mục
func shuffle(slice []int) {
    for i := 0; i < len(slice); i++ {
        a := rand.Intn(len(slice))
        b := rand.Intn(len(slice))
        slice[a], slice[b] = slice[b], slice[a]
    }
}

func request(params map[string]interface{}) error {
    var indexes = []int {0,1,2,3,4,5,6}
    var err error

    // gọi shuffle để xáo trộn các index
    shuffle(indexes)

    // số lần thử lại là 3
    maxRetryTimes := 3

    idx := 0
    for i := 0; i < maxRetryTimes; i++ {
        err = apiRequest(params, indexes[idx])
        if err == nil {
            break
        }
        idx++
    }

    if err != nil {
        // logging
        return err
    }

    return nil
}
```

Chúng ta duyệt qua các chỉ mục và hoán đổi chúng, tương tự như phương pháp xáo trộn mà chúng ta thường sử dụng khi chơi bài.

### 5.5.2.1 Tải không cân bằng gây ra do xáo trộn không chính xác

Thực sự không có vấn đề? Trong thực tế, vẫn còn vấn đề. Có hai cạm bẫy tiềm ẩn trong chương trình trên là:

1. Không có random seed. Khi không có `random seed`, trình tự của các lần random  `rand.Intn()` là cố định.

2. Xáo trộn không đều, điều này sẽ khiến node đầu tiên của toàn bộ mảng có xác suất được chọn cao và phân phối tải giữa các node không cân bằng.

Điểm đầu tiên tương đối đơn giản nên chúng tôi không nêu ví dụ cụ thể. Về điểm thứ hai, chúng ta có thể sử dụng kiến ​​thức về xác suất để chứng minh điều đó. Giả sử rằng mỗi lựa chọn là thực sự ngẫu nhiên, xác suất mà node ở vị trí đầu tiên không được chọn trong trao đổi `len(slice)` là `((6/7)*(6/7))^7 ≈ 0.34`. Trong trường hợp phân phối đồng đều, chúng ta chắc chắn muốn xác suất phần tử đầu tiên được phân phối tại bất kỳ vị trí nào bằng nhau, do đó xác suất được chọn ngẫu nhiên phải xấp xỉ bằng `1/7≈0.14`.

Rõ ràng, thuật toán xáo trộn được đưa ra ở đây có xác suất 30% không hoán đổi các yếu tố cho bất kỳ vị trí nào. Vì vậy, tất cả các yếu tố có xu hướng ở lại vị trí ban đầu của chúng. Bởi vì mỗi lần chúng ta nhập cùng một chuỗi cho mảng `shuffle`, phần tử đầu tiên có xác suất được chọn cao hơn. Trong trường hợp cân bằng tải, có nghĩa là tải máy đầu tiên trong mảng node sẽ cao hơn nhiều so với các máy khác (ít nhất gấp 3 lần).

### 5.5.2.2 Sửa thuật toán xáo trộn

Thuật toán [fishing-yates](https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle) đã chứng minh tính đúng đắn về mặt toán học. Ý tưởng chính của nó là chọn một giá trị ngẫu nhiên rồi đặt ở cuối mảng, và cứ thế tiếp tục. Ví dụ:

```go
func shuffle(indexes []int) {
    for i:=len(indexes); i>0; i-- {
        lastIdx := i - 1
        idx := rand.Int(i)
        indexes[lastIdx], indexes[idx] = indexes[idx], indexes[lastIdx]
    }
}
```

Thuật toán đã được hiện thực trong thư viện chuẩn [ math/rand](https://golang.org/pkg/math/rand/) của Go:

```go
func shuffle(n int) []int {
    b := rand.Perm(n)
    return b
}
```

Hiện tại, chúng ta có thể sử dụng `rand.Perm` để lấy mảng chỉ mục mà chúng ta muốn.

## 5.5.3 Vấn đề chọn node ngẫu nhiên cho cụm ZooKeeper

Giả sử, ta cần chọn một node từ N node để gửi yêu cầu. Sau khi yêu cầu ban đầu kết thúc, các yêu cầu tiếp theo sẽ xáo trộn lại mảng, do đó không có mối quan hệ nào giữa hai yêu cầu. Ví thế, thuật toán ở trên sẽ không cần khởi tạo bất kì random seed nào.

Tuy nhiên, trong một số trường hợp đặc biệt, chẳng hạn như khi sử dụng ZooKeeper, khi máy khách khởi tạo việc lựa chọn node từ nhiều node dịch vụ, một kết nối được thiết lập cho node. Yêu cầu máy khách sau đó được gửi đến node. Node tiếp theo trong danh sách được chọn cho đến khi không còn node nào có sẵn. Lúc này, việc lựa chọn node kết nối ban đầu là "đúng chuẩn" ngẫu nhiên. Tuy nhiên, tất cả các máy khách sẽ kết nối với cùng một ZooKeeper khi chúng khởi động cùng lúc, lúc này sẽ không có tải cân bằng. Nếu doanh nghiệp của bạn cần phát triển tính năng hàng ngày, thì bạn phải xem xét liệu có một tình huống tương tự như trên xảy ra không. Cách đặt random seed cho thư viện rand:

```go
rand.Seed(time.Now().UnixNano())
```

Lý do cho những kết luận này là phiên bản trước của thư viện Open source ZooKeeper được sử dụng rộng rãi đã mắc phải những lỗi trên và mãi đến đầu năm 2016, vấn đề mới được khắc phục.

## 5.5.4 Kiểm tra lại ảnh hưởng của thuật toán cân bằng tải

Chúng ta không xét trường hợp cân bằng tải có trọng số ở đây. Bây giờ, điều quan trọng nhất là sự cân bằng. Chúng ta chỉ đơn giản so sánh thuật toán xáo trộn trong phần mở đầu với kết quả của thuật toán fisher yates:

***main.go***
```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func shuffle1(slice []int) {
    for i := 0; i < len(slice); i++ {
        a := rand.Intn(len(slice))
        b := rand.Intn(len(slice))
        slice[a], slice[b] = slice[b], slice[a]
    }
}

func shuffle2(indexes []int) {
    for i := len(indexes); i > 0; i-- {
        lastIdx := i - 1
        idx := rand.Intn(i)
        indexes[lastIdx], indexes[idx] = indexes[idx], indexes[lastIdx]
    }
}

func main() {
    var cnt1 = map[int]int{}
    for i := 0; i < 1000000; i++ {
        var sl = []int{0, 1, 2, 3, 4, 5, 6}
        shuffle1(sl)
        cnt1[sl[0]]++
    }

    var cnt2 = map[int]int{}
    for i := 0; i < 1000000; i++ {
        var sl = []int{0, 1, 2, 3, 4, 5, 6}
        shuffle2(sl)
        cnt2[sl[0]]++
    }

    fmt.Println(cnt1, "\n", cnt2)
}
```

Kết quả:

```shell
map[0:224436 1:128780 5:129310 6:129194 2:129643 3:129384 4:129253]
map[6:143275 5:143054 3:143584 2:143031 1:141898 0:142631 4:142527]
```

Kết quả trên phù hợp với kết luận đã đưa ra.

[Tiếp theo](ch5-06-config.md)