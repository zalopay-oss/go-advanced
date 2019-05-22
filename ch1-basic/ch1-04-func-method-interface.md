## 1.4 Functions, Methods và Interfaces

Hàm (Function) tương ứng với chuỗi các thao tác và là thành phần cơ bản của chương trình. Các hàm trong ngôn ngữ Go có thể có tên hoặc ẩn danh (anonymous function): hàm được đặt tên thường tương ứng với hàm cấp package (package function). Đây là trường hợp đặc biệt của hàm ẩn danh. Khi một hàm ẩn danh tham chiếu một biến trong phạm vi bên ngoài, nó sẽ trở thành hàm đóng. Các package function là cốt lõi của một ngôn ngữ lập trình hàm (functional programming).

Phương thức (Method) được liên kết với một hàm đặc biệt của một kiểu cụ thể. Các phương thức trong ngôn ngữ Go phụ thuộc vào kiểu và phải được ràng buộc tĩnh tại thời gian biên dịch.

Một Interface xác định một tập hợp các phương thức phụ thuộc vào đối tượng Interface trong thời gian thực thi, vì vậy các phương thức tương ứng với Interface được ràng buộc động khi thực thi. Ngôn ngữ Go hiện thực mô hình hướng đối tượng thông qua cơ chế Interface ngầm định.

Việc khởi tạo và thực thi chương trình Go luôn bắt đầu từ hàm `main.main`. Nếu package `main` có import  các package khác, chúng sẽ được thêm vào package `main` theo thứ tự khai báo.

- Nếu một package được import nhiều lần, sẽ chỉ được tính là một khi thực thi.
- Khi một package được import mà nó lại import các package khác, trước tiên Go sẽ import các package khác đó trước, sau đó  khởi tạo các hằng và biến của package, rồi gọi hàm `init` trong từng package.
- Nếu một package có nhiều hàm `init` và thứ tự gọi không được xác định cụ thể(phần implement có thể được gọi theo thứ tự tên file), thì chúng sẽ được gọi theo thứ tự xuất hiện (`init` không phải là hàm thông thường, nó có thể có nhiều định nghĩa, và các hàm khác không thể sử dụng nó). Cuối cùng, khi `main` đã có đủ tất cả hằng và biến ở cấp package, chúng sẽ được khởi tạo bằng cách thực thi hàm `init`, tiếp theo chương trình đi vào hàm `main.main` và  bắt đầu thực thi. Hình dưới đây là sơ đồ nguyên lý của một chuỗi bắt đầu của chương trình hàm trong Go:

<p align="center">

<img src="../images/ch1-11-init.ditaa.png">
<span align="center">Hình 1-11. Tiến trình khởi tạo package</span>

</p>

Cần lưu ý rằng trong `main.main` tất cả các mã lệnh đều chạy trong cùng một goroutine trước khi hàm được thực thi, đây là thread chính của chương trình. Do đó, nếu một hàm `init` khởi chạy từ hàm `main` trong một goroutine mới với từ khóa go, thì goroutine đó chỉ có `main.main` có thể được thực thi sau khi vào hàm.

Cần lưu ý rằng trước khi hàm `main.main` được thực thi thì tất cả code đều chạy trong cùng một goroutine, đây là thread chính của chương trình. Do đó, nếu một hàm `init` khởi động bên trong một goroutine mới với từ khóa go, goroutine đó chỉ có thể được thực thi sau khi vào hàm `main.main`.

### 1.4.1 Function

Trong Go, hàm là kiểu đầu tiên của đối tượng  và chúng ta có thể giữ hàm trong một biến. Hàm có thể được đặt tên hoặc ẩn danh (anonymous). Các hàm cấp độ package thường là các hàm được đặt tên. Hàm được đặt tên là một trường hợp đặc biệt của hàm ẩn danh. Tất nhiên, mỗi kiểu trong ngôn ngữ Go cũng có thể có các phương thức riêng, và đó có thể là là một hàm:

```go
// hàm được đặt tên
func Add(a, b int) int {
    return a+b
}

// hàm ẩn danh
var Add = func(a, b int) int {
    return a+b
}

```

Một hàm trong ngôn ngữ Go có thể có nhiều tham số và nhiều giá trị trả về. Cả tham số và giá trị trả về trao đổi dữ liệu  với hàm được gọi theo cách truyền vào giá trị (pass by value). Về mặt cú pháp, hàm cũng hỗ trợ số lượng tham số thay đổi, biến số lượng tham số phải là tham số cuối cùng và biến này phải là kiểu slice.

```go
// Nhiều tham số và nhiều giá trị trả về
func Swap(a, b int) (int, int) {
    return b, a
}

// Biến số lượng tham số 'more'
// Tương ứng với kiểu [] int, là một slice
func Sum(a int, more ...int) int {
    for _, v := range more {
        a += v
    }
    return a
}
```

Khi đối số có thể thay đổi là một kiểu interface null,  việc người gọi có giải nén (unpack) đối số đó hay không sẽ dẫn đến những kết quả khác nhau:

```go
func main() {
    var a = []interface{}{123, "abc"}

    Print(a...) // 123 abc
    Print(a)    // [123 abc]
}

func Print(a ...interface{}) {
    fmt.Println(a...)
}
```

Lời gọi `Print` đầu tiên  truyền vào  `a...` tương đương với lời gọi trực tiếp `Print(123, "abc")`. Lời gọi `Print` thứ hai truyền vào `a`, tương đương với lời gọi `Print([]interface{}{123, "abc"})`.

Cả tham số truyền vào và các giá trị trả về đều có thể được đặt tên:

```go
func Find(m map[int]int, key int) (value int, ok bool) {
    value, ok = m[key]
    return
}
```

Nếu giá trị trả về được đặt tên, nó có thể sửa đổi  bằng tên hoặc có thể sửa đổi bằng lệnh `defer` sẽ thực thi sau lệnh `return`

```go
func Inc() (v int) {
    defer func(){ v++ } ()
    return 42
}
// giá trị v cuối cùng là 43
```

Câu lệnh `defer` sẽ trì hoãn việc thực thi của hàm ẩn danh (trong ví dụ trên) vì hàm này lấy biến cục bộ `v` của hàm bên ngoài. Hàm này được gọi là bao đóng. Bao đóng không truy cập tới biến bên ngoài (như `v`) theo kiểu giá trị (value-by-value) mà truy cập bằng tham chiếu (reference).

Hành vi truy cập các biến bên ngoài bằng tham chiếu này đến các bao đóng có thể dẫn đến một số vấn đề tiềm ẩn:

```go
func main() {
    for i := 0; i < 3; i++ {
        defer func(){ println(i) } ()
    }
}
// Output:
// 3
// 3
// 3
```

Bởi vì nó là một bao đóng (hàm trong câu lệnh lặp for), mỗi câu lệnh `defer` trì hoãn việc thực hiện tham chiếu hàm tới cùng một biến lặp i, giá trị của biến này sau khi kết thúc vòng lặp là 3, do đó đầu ra cuối cùng là 3.

Với ý tưởng là tạo ra một biến duy nhất cho mỗi hàm `defer` trong mỗi lần lặp. Có hai cách để làm điều này:

```go
func main() {
    for i := 0; i < 3; i++ {
        i := i // Xác định một biến cục bộ i trong vòng lặp
        defer func(){ println(i) } ()
    }
}

func main() {
    for i := 0; i < 3; i++ {
        // truyền i vào hàm (pass by value)
        // câu lệnh defer sẽ lấy các tham số từ lời gọi
        defer func(i int){ println(i) } (i)
    }
}
```

- Phương pháp đầu tiên là xác định một biến cục bộ bên trong thân vòng lặp, để hàm bao đóng của câu lệnh `defer` lấy các biến khác nhau cho mỗi lần lặp. Các giá trị của các biến này tương ứng với các giá trị tại thời điểm lặp.
- Cách thứ hai là truyền biến lặp iterator thông qua các tham số của hàm bao đóng và câu lệnh `defer` sẽ ngay lập tức lấy các tham số từ lời gọi (trường hợp này là lấy `i`).

Cả hai phương pháp đều hoạt động. Tuy nhiên, đây không phải là cách thực hành tốt để thực thi câu lệnh `defer` bên trong vòng lặp for. Đây chỉ là ví dụ và không được khuyến khích.

Trong ngôn ngữ Go, nếu một hàm được gọi với tham số là kiểu slice thì một tham số ảo sẽ được truyền vào bởi vì phần tử của slice có thể được sửa đổi bên trong hàm được gọi.

Trong thực tế, trường hợp mà một tham số ở lời gọi hàm bị sửa đổi bởi thao tác trong  hàm là bởi vì nó là con trỏ được truyền tường minh hoặc ngầm định vào tham số hàm. Đặc tả tham số hàm chỉ đề cập đến phần cố định của cấu trúc dữ liệu, chẳng hạn như cấu trúc con trỏ hoặc độ dài chuỗi (trong cấu trúc chuỗi) hoặc slice tương ứng, nhưng không chứa nội dung trỏ tới bởi con trỏ gián tiếp.

Việc thay thế tham số của kiểu slice với cấu trúc tương tự là `reflect.SliceHeader` là một ví dụ để hiểu ý nghĩa của việc truyền vào giá trị kiểu slice (pass by value):

```go
// truyền vào con trỏ ngầm định khiến
// nội dung của biến x bị thay đổi
func twice(x []int) {
    for i := range x {
        x[i] *= 2
    }
}

type IntSliceHeader struct {
    Data []int
    Len  int
    Cap  int
}

func twice(x IntSliceHeader) {
    for i := 0; i < x.Len; i++ {
        x.Data[i] *= 2
    }
}
```

Vì phần array bên dưới của kiểu slice được truyền bởi con trỏ ngầm định (chính con trỏ vẫn được truyền, nhưng con trỏ trỏ đến cùng một dữ liệu), nên hàm được gọi có thể sửa đổi dữ liệu trong slice thông qua con trỏ.  cấu trúc `IntSliceHeader` chứa không chỉ dữ liệu mà còn có thông tin về độ dài và dung lượng slice, hai thành phần này cũng được truyền theo giá trị. Nếu có hàm nào điều chỉnh `Len` hoặc `Cap` được gọi, nó sẽ không thể hiện sự thay đổi đó trong biến slice của tham số hàm được. Lúc này, ta nên cập nhật slice trước bằng cách trả về slice đã sửa đổi. Đây cũng là lý do tại sao hàm `append` (built-in) phải trả về một slice.

Trong ngôn ngữ Go, các hàm cũng có thể tự gọi chính nó trực tiếp hoặc gián tiếp (gọi đệ quy). Không có giới hạn về độ sâu của lệnh gọi đệ quy trong Go. Stack của lệnh gọi hàm không có lỗi tràn, vì trong thời gian thực thi Go tự động điều chỉnh kích thước của stack hàm khi cần.

Mỗi goroutine sẽ  được phân bổ một stack nhỏ (4 hoặc 8KB, tùy thuộc vào implement) ngay sau khi khởi động. Kích thước stack có thể được điều chỉnh động khi cần. Stack có thể đạt đến mức GB (tùy theo cách implement, trong phiên bản hiện tại là 32 bit) Kiến trúc là 250MB và kiến ​​trúc 64 bit là 1GB).

Trước phiên bản 1.4, Go sử dụng stack động phân đoạn (Segmented dynamic stack). Về cơ bản thì một danh sách liên kết (linked list) được sử dụng để hiện thực các stack động. Địa chỉ bộ nhớ của các node trong mỗi danh sách liên kết là không thay đổi. Tuy nhiên, các stack động này có ảnh hưởng lớn đến hiệu suất của một số lời gọi ở những thời điểm quan trọng. Nguyên nhân là bởi  vì các node  trong danh sách liên kết dù có liền kề cũng sẽ không liền kề trong địa chỉ bộ nhớ, làm tăng khả năng xảy ra lỗi bộ nhớ cache của CPU (cache hit failure).

Để giải quyết vấn đề về tỉ lệ trúng CPU cache (hit rate) nói trên, Go 1.4 sử dụng hiện thực stack động liên tục (Continuous dynamic stack), nghĩa là dùng một cấu trúc tương tự như mảng động để biểu diễn stack. Tuy nhiên, stack động liên tục cũng mang đến một vấn đề mới: khi stack tăng kích thước động, nó cần di chuyển dữ liệu trước đó sang không gian bộ nhớ mới, điều này sẽ khiến địa chỉ của tất cả các biến trong stack trước đó thay đổi.

Mặc dù trong thời điểm thực thi Go tự động cập nhật các con trỏ để lưu trữ (vào stack) các biến tham chiếu tới địa chỉ mới, nhưng quan trọng  là các con trỏ trong Go không còn cố định nữa(vì vậy ta không thể giữ con trỏ trong các biến theo ý muốn, địa chỉ trong Go không thể được lưu vào môi trường không được kiểm soát bởi GC, đó là lý do địa chỉ của đối tượng Go không thể được giữ bằng ngôn ngữ C trong một thời gian dài khi sử dụng CGO).

Vì stack của các hàm trong Go sẽ tự động thay đổi kích thước, lập trình viên hiếm khi cần quan tâm đến cơ chế hoạt động của stack. Trong đặc tả ngôn ngữ, ngay cả khái niệm stack và heap cũng không được đề cập một cách có chủ ý. Chúng ta không thể biết được một tham số hàm hoặc một biến cục bộ sẽ lưu trữ trên stack hay trên heap. Chúng ta chỉ cần biết rằng chúng hoạt động tốt là được. Hãy xem ví dụ sau:

```go
func f(x int) *int {
    return &x
}

func g() int {
    x = new(int)
    return *x
}
```

- Hàm đầu tiên trả về trực tiếp địa chỉ của biến tham số hàm (biến `x`) - điều này có vẻ là không khả thi bởi vì nếu biến tham số nằm trên stack sẽ trở thành không hợp lệ sau khi hàm trả về và địa chỉ được trả về dĩ nhiên bị lỗi. Nhưng trình biên dịch của  Go thông minh hơn khi đảm bảo rằng các biến được trỏ bởi con trỏ sẽ ở đúng vị trí.
- Hàm thứ hai, mặc dù lời gọi `new` tạo một đối tượng con trỏ kiểu `*int`, nhưng vẫn không biết nó sẽ được lưu ở đâu. Với những lập trình viên có kinh nghiệm với C/C ++, điều quan trọng phải nhấn mạnh rằng trình biên dịch và thực thi (runtime) sẽ giúp chúng ta không phải lo lắng về stack và heap của hàm trong Go. Ngoài ra cũng đừng cho rằng vị trí của biến trong bộ nhớ là cố định. Con trỏ có thể thay đổi bất cứ lúc nào, đặc biệt là những khi chúng ta không mong đợi nó thay đổi nhất.
