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

Trong ngôn ngữ Go, nếu một hàm được gọi với một slice làm tham số, đôi khi một ảo ảnh truyền tham chiếu được đưa ra cho tham số: bởi vì phần tử của slice đến có thể được sửa đổi bên trong hàm được gọi. Trong thực tế, bất kỳ tình huống trong đó một tham số cuộc gọi có thể được sửa đổi bởi một tham số chức năng là bởi vì tham số con trỏ được truyền rõ ràng hoặc ngầm định trong tham số chức năng. Đặc tả của giá trị tham số hàm là chính xác hơn. Nó chỉ đề cập đến phần cố định của cấu trúc dữ liệu, chẳng hạn như cấu trúc con trỏ hoặc độ dài chuỗi trong cấu trúc chuỗi hoặc lát tương ứng, nhưng không chứa nội dung trỏ con trỏ gián tiếp. . Thay thế các tham số của loại lát bằng một cái gì đó giống như cấu trúc phản chiếu.SliceHeader là một cách tốt để hiểu ý nghĩa của giá trị lát:

```go
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