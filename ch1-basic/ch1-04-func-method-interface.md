# 1.4 Functions, Methods và Interfaces

Hàm (function) là  thành phần cơ bản của chương trình. Các hàm trong ngôn ngữ Go có thể có tên hoặc ẩn danh (anonymous function): hàm được đặt tên thường tương ứng với hàm cấp package (package function). Đây là trường hợp đặc biệt của hàm ẩn danh. Khi một hàm ẩn danh tham chiếu một biến trong phạm vi bên ngoài, nó sẽ trở thành hàm đóng. Các package function là cốt lõi của một ngôn ngữ lập trình hàm (functional programming).

Phương thức (Method) được liên kết với một hàm đặc biệt của một kiểu cụ thể. Các phương thức trong ngôn ngữ Go phụ thuộc vào kiểu và phải được ràng buộc tĩnh tại thời gian biên dịch.

Một Interface xác định một tập hợp các phương thức phụ thuộc vào đối tượng Interface trong thời gian thực thi, vì vậy các phương thức tương ứng với Interface được ràng buộc động khi thực thi. Ngôn ngữ Go hiện thực mô hình hướng đối tượng thông qua cơ chế Interface ngầm định.

Việc khởi tạo và thực thi chương trình Go luôn bắt đầu từ hàm `main.main`. Nếu package `main` có import  các package khác, chúng sẽ được thêm vào package `main` theo thứ tự khai báo.

- Nếu một package được import nhiều lần, sẽ chỉ được tính là một khi thực thi.
- Khi một package được import mà nó lại import các package khác, trước tiên Go sẽ import các package khác đó trước, sau đó  khởi tạo các hằng và biến của package, rồi gọi hàm `init` trong từng package.
- Nếu một package có nhiều hàm `init` và thứ tự gọi không được xác định cụ thể (phần implement có thể được gọi theo thứ tự tên file), thì chúng sẽ được gọi theo thứ tự xuất hiện (`init` không phải là hàm thông thường, nó có thể có nhiều định nghĩa, và các hàm khác không thể sử dụng nó). Cuối cùng, khi `main` đã có đủ tất cả hằng và biến ở cấp package, chúng sẽ được khởi tạo bằng cách thực thi hàm `init`, tiếp theo chương trình đi vào hàm `main.main` và  bắt đầu thực thi. Hình dưới đây là sơ đồ nguyên lý  một chuỗi bắt đầu của chương trình hàm trong Go:

<p align="center">

<img src="../images/ch1-11-init.ditaa.png">
<span align="center">Hình 1-11. Tiến trình khởi tạo package</span>

</p>

Cần lưu ý rằng trong `main.main` tất cả các mã lệnh đều chạy trong cùng một Goroutine trước khi hàm được thực thi, đây là thread chính của chương trình. Do đó, nếu một hàm `init` khởi chạy từ hàm `main` trong một Goroutine mới với từ khóa go, thì Goroutine đó chỉ có `main.main` có thể được thực thi sau khi vào hàm.

Cần lưu ý rằng trước khi hàm `main.main` được thực thi thì tất cả code đều chạy trong cùng một Goroutine, đây là thread chính của chương trình. Do đó, nếu một hàm `init` khởi động bên trong một Goroutine mới với từ khóa go, Goroutine đó chỉ có thể được thực thi sau khi vào hàm `main.main`.

## 1.4.1 Function

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

Khi đối số có thể thay đổi là một kiểu interface null,  việc người gọi có phân giải (unpack) đối số đó hay không sẽ dẫn đến những kết quả khác nhau:

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

Với ý tưởng là tạo ra một biến duy nhất cho mỗi hàm `defer` trong mỗi lần lặp. Có hai cách để làm điều này: ([source1,](../examples/ch1/ch1.4/1-function/example-2/main.go)[(source2](../examples/ch1/ch1.4/1-function/example-3/main.go)

[>> mã nguồn 1,](../examples/ch1/ch1.4/1-function/example-2/main.go)[>> mã nguồn 2](../examples/ch1/ch1.4/1-function/example-3/main.go)

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

Mỗi Goroutine sẽ  được phân bổ một stack nhỏ (4 hoặc 8KB, tùy thuộc vào implement) ngay sau khi khởi động. Kích thước stack có thể được điều chỉnh động khi cần. Stack có thể đạt đến mức GB (tùy theo cách implement, trong phiên bản hiện tại là 32 bit) Kiến trúc là 250MB và kiến ​​trúc 64 bit là 1GB).

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
- Hàm thứ hai, mặc dù lời gọi `new` tạo một đối tượng con trỏ kiểu `*int`, nhưng vẫn không biết nó sẽ được lưu ở đâu. Một điều nói riêng với những lập trình viên có kinh nghiệm với C/C ++ là trình biên dịch và thực thi (runtime) của Go sẽ giúp chúng ta không phải lo lắng về stack và heap của hàm. Do đó đừng cho rằng địa chỉ của biến trong bộ nhớ là cố định vì con trỏ có thể thay đổi bất cứ lúc nào, đặc biệt là những khi chúng ta không mong đợi nó thay đổi nhất.

## 1.4.2 Method

Phương thức (Method) là một tính năng của lập trình hướng đối tượng (OOP). Trong ngôn ngữ C++, phương thức  tương ứng với một hàm thành viên của một đối tượng lớp, được liên kết với một bảng ảo trên một đối tượng cụ thể. Tuy nhiên, phương thức trong ngôn ngữ Go được liên kết với kiểu, do đó liên kết tĩnh của phương thức có thể được tạo thành trong giai đoạn biên dịch.

Một chương trình hướng đối tượng sử dụng các phương thức để thể hiện nhưng thao tác trên thuộc tính (properties) của nó, qua đó người dùng có thể sử dụng đối tượng mà không cần phải thao tác trực tiếp với đối tượng mà là thông qua các phương thức. C++ thường được xem là nơi mà lập trình hướng đối tượng bắt đầu phát triển mạnh. C++ hỗ trợ các tính năng hướng đối tượng (như lớp) dựa trên cơ sở ngôn ngữ C. Sau đó đến Java được gọi là ngôn ngữ hướng đối tượng thuần túy  vì các hàm của nó không thể tồn tại độc lập mà phải thuộc về một lớp nhất định.

Lập trình hướng đối tượng là một ý tưởng. Nhiều ngôn ngữ tuyên bố hỗ trợ lập trình hướng đối tượng chỉ đơn giản là kết hợp các tính năng thường được sử dụng vào ngôn ngữ. Mặc dù ngôn ngữ C tổ tiên của ngôn ngữ Go không phải là ngôn ngữ hướng đối tượng, các hàm liên quan đến file trong thư viện chuẩn ngôn ngữ C cũng sử dụng ý tưởng lập trình hướng đối tượng. Dưới đây là hiện thực một tập hợp các hàm làm việc với file theo kiểu ngôn ngữ C:

[>> mã nguồn](../examples/ch1/ch1.4/2-method/example-1/main.go)

```go
// đối tượng File
type File struct {
    fd int
}

// mở file
func OpenFile(name string) (f *File, err error) {
    // ...
}

// đóng file
func CloseFile(f *File) error {
    // ...
}

// đọc dữ liệu từ file
func ReadFile(f *File, offset int64, data []byte) int {
    // ...
}
```

Hàm `OpenFile` xây dựng như constructor để mở một đối tượng kiểu file, `CloseFile` tương tự như destructor dùng để đóng lại đối tượng, `ReadFile` là một hàm thành viên, ba hàm này đều là các hàm thông thường. Với `CloseFile` Và `ReadFile` ta cần chiếm tài nguyên tên trong không gian cấp độ package. Tuy nhiên `CloseFile` hay `ReadFile` chỉ là các hàm thao tác trên đối tượng kiểu `File`. Tại thời điểm này, ta muốn các hàm đó được gắn chặt với các kiểu đối tượng hoạt động.

Ngôn ngữ Go thực hiện `CloseFile` và `ReadFile` bằng cách chuyển tham số đầu tiên lên đầu của tên hàm:

```go
// đóng file
func (f *File) CloseFile() error {
    // ...
}

// đọc dữ liệu từ file
func (f *File) ReadFile (offset int64, data []byte) int {
    // ...
}
```

Trong trường hợp này, hàm `CloseFile` và `ReadFile` trở thành  phương thức duy nhất của kiểu  `File`(thay vì phương thức đối tượng `File`). Chúng cũng không còn chiếm tài nguyên tên trong không gian cấp độ package và kiểu `File` đã làm rõ các thao tác trên đối tượng của chúng, vì vậy tên phương thức thường được đơn giản hóa thành `Close` và `Read`:

```go
// đóng file
func (f *File) Close() error {
    // ...
}

// đọc dữ liệu từ file
func (f *File) Read(offset int64, data []byte) int {
    // ...
}
```

Việc di chuyển tham số đầu tiên của hàm lên phía đầu của tên hàm chỉ là một thay đổi nhỏ trong code, nhưng từ quan điểm triết lý lập trình, ngôn ngữ Go đã đứng trong hàng ngũ của các ngôn ngữ hướng đối tượng. Ta có thể thêm một hoặc nhiều phương thức cho bất kỳ kiểu tùy chỉnh nào (custom type). Phương thức cho mỗi kiểu phải nằm trong cùng một package với định nghĩa kiểu, do đó không thể thêm phương thức vào các kiểu dựng sẵn đó (vì định nghĩa của phương thức và định nghĩa của kiểu không nằm trong package). Đối với một kiểu nhất định, tên của mỗi phương thức phải là duy nhất và các phương thức cũng như hàm đều không hỗ trợ overload.

Phương thức được bắt nguồn từ hàm, chỉ là di chuyển tham số đối tượng đầu tiên của hàm lên phía trước tên hàm. Vì vậy, chúng ta vẫn có thể sử dụng phương thức theo tư duy thủ tục (procedure). Ta có thể biến một phương thức thành một loại hàm thông thường bằng cách gọi các thuộc tính trong biểu thức của nó:

[>> mã nguồn](../examples/ch1/ch1.4/2-method/example-2/main.go)

```go
// không phụ thuộc vào đối tượng file cụ thể
// func CloseFile(f *File) error
var CloseFile = (*File).Close

// không phụ thuộc vào đối tượng file cụ thể
// func ReadFile(f *File, offset int64, data []byte) int
var ReadFile = (*File).Read

// xử lý file
f, _ := OpenFile("foo.dat")
ReadFile(f, 0, data)
CloseFile(f)
```

Trong một số tình huống, ta quan tâm nhiều hơn đến một chuỗi thao tác ví dụ  như `Read` đọc một số mảng và sau đó gọi `Close` để đóng, trong ngữ cảnh này, người dùng không quan tâm đến kiểu của đối tượng, miễn là nó có thể đáp ứng được các thao tác của `Read` và `Close`. Tuy nhiên trong các biểu thức phương thức của `ReadFile`, `CloseFile` có chỉ rõ kiểu `File` trong tham số kiểu sẽ khiến chúng không bị phụ thuộc vào đối tượng nào cụ thể. Việc này có thể khắc phục bằng cách sử dụng thuộc tính bao đóng (closure property):

[>> mã nguồn](../examples/ch1/ch1.4/2-method/example-3/main.go)

```go
// mở đối tượng file
f, _ := OpenFile("foo.dat")

// liên kết với đối tượng f
// func Close() error
var Close = func() error {
    return (*File).Close(f)
}

// liên kết với đối tượng f
// func Read (offset int64, data []byte) int
var Read = func(offset int64, data []byte) int {
    return (*File).Read(f, offset, data)
}

// xử lý file
Read(0, data)
Close()
```

Đây chính là vấn đề mà giá trị phương thức cần giải quyết. Chúng ta có thể đơn giản hóa việc  hiện thực với các tính năng:

[>> mã nguồn](../examples/ch1/ch1.4/2-method/example-4/main.go)

```go
// mở đối tượng file
f, _ := OpenFile("foo.dat")

// giá trị phương thức: ràng buộc với đối tượng f
// func Close() error
var Close = f.Close

// giá trị phương thức: ràng buộc với đối tượng f
// func Read (offset int64, data []byte) int
var Read = f.Read

// xử lý file
Read(0, data)
Close()
```

Go không hỗ trợ tính năng kế thừa như các ngôn ngữ hướng đối tượng truyền thống nhưng sẽ hỗ trợ việc kế thừa phương thức theo sự kết hợp độc đáo của riêng mình. Với ngôn ngữ Go, tính kế thừa đạt được bằng cách xây dựng các thành phần ẩn danh trong structure:

```go
import "image/color"

type Point struct{ X, Y float64 }

type ColoredPoint struct {
    Point
    Color color.RGBA
}
```

Chúng ta có thể định nghĩa `ColoredPoint` như một struct có 3 trường, nhưng ở đây chúng ta sẽ dùng struct `Point` chứa `X` và `Y` để thay thế.

```go
var cp ColoredPoint
cp.X = 1
fmt.Println(cp.Point.X) // "1"
cp.Point.Y = 2
fmt.Println(cp.Y)       // "2"
```

Bằng cách sử dụng các thành phần ẩn danh, chúng ta có thể kế thừa không chỉ các thành phần nội bộ (`X` và `Y`), mà cả các phương thức tương ứng với các kiểu của chúng. Ta thường nghĩ rằng `Point` là một lớp cơ sở và `ColoredPoint` là lớp kế thừa hoặc lớp con của nó. Tuy nhiên, phương thức được kế thừa theo cách này không thể hiện tính đa hình của  hàm ảo trong C++. Tham số chỗ  hàm nhận của tất cả các phương thức được kế thừa vẫn là thành phần ẩn danh, không phải là biến hiện tại.

```go
type Cache struct {
    m map[string]string
    sync.Mutex
}

func (p *Cache) Lookup(key string) string {
    p.Lock()
    defer p.Unlock()

    return p.m[key]
}
```

Cấu trúc `Cache` nhúng một kiểu ẩn danh `sync.Mutex` để kế thừa phương thức  `Lock` và `Unlock` từ đó, các lời gọi `p.Lock()` và `p.Unlock()` với `p`là đối tượng nhận của phương thức,  chúng sẽ được triển khai thành `p.Mutex.Lock()` và `p.Mutex.Unlock()`. Sự mở rộng này được hoàn thành lúc biên dịch và không mất chi phí runtime.

Đối với tính kế thừa trong  ngôn ngữ hướng đối tượng truyền thống (như C ++ hoặc Java), phương thức ở lớp con được liên kết động với đối tượng khi chạy, do đó một số phương thức hiện thực lớp cơ sở `this` có thể không tương ứng với kiểu của lớp cơ sở. Những đối tượng khác nhau gây ra sự không chắc chắn trong hoạt động của phương thức lớp cơ sở. Phương thức của lớp cơ sở trong ngôn ngữ Go "kế thừa" bằng cách nhúng thêm các thành phần ẩn danh `this` là đối tượng hiện thực kiểu của phương thức. Phương thức trong ngôn ngữ Go bị ràng buộc tĩnh tại thời gian biên dịch.

Nếu cần tính chất đa hình ở các hàm ảo, chúng ta cần triển khai nó với Interface.

## 1.4.3 Interface

Rob Pike, cha đẻ của ngôn ngữ Go, đã từng nói một câu nói nổi tiếng:

> Languages ​​that try to disallow idiocy become themselves idiotic 
> (Các ngôn ngữ cố gắng tránh các hành vi ngu ngốc cuối cùng trở thành ngôn ngữ ngu ngốc).

Các ngôn ngữ lập trình tĩnh nói chung có các hệ thống kiểu nghiêm ngặt, cho phép trình biên dịch đi sâu vào xem liệu lập trình viên có thực hiện bất kỳ động thái bất thường nào không. Tuy nhiên, một hệ thống kiểu quá nghiêm ngặt có thể làm cho việc lập trình trở nên quá cồng kềnh và khiến  lập trình viên lãng phí rất nhiều thời gian tuổi trẻ trong công cuộc đấu tranh với trình biên dịch.

Ngôn ngữ Go  vì thế cố gắng cung cấp sự cân bằng giữa lập trình an toàn và lập trình linh hoạt. Nó  hỗ trợ  duck-typing thông qua interface đồng thời cũng có  kiểm tra kiểu nghiêm ngặt, giúp việc lập trình tương đối nhẹ nhàng hơn.

Interface type của Go là một sự trừu tượng hóa và khái quát hóa các loại hành vi khác, bởi vì kiểu interface không gắn với các chi tiết implement cụ thể, chúng ta có thể làm cho đối tượng linh hoạt hơn và dễ dùng hơn thông qua sự trừu tượng hóa này.

Nhiều ngôn ngữ hướng đối tượng có các khái niệm interface tương tự, nhưng interface trong Go là duy nhất ở chỗ nó là duck-typing thỏa mãn việc implement ngầm định. Duck-type nói rằng: *Miễn là nó đi như vịt và kêu như vịt, bạn có thể sử dụng nó như một con vịt*.

Nếu một đối tượng trông giống như phần  hiện thực của một interface, thì nó có thể được sử dụng như thể nó thuộc kiểu interface đó. Thiết kế này cho phép chúng ta tạo ra một interface mới thỏa mãn kiểu hiện có mà không phải  hủy đi định nghĩa ban đầu của chúng, thiết kế này đặc biệt linh hoạt và hữu ích khi các kiểu mà ta sử dụng đến từ những package không thuộc quyền kiểm soát của ta. Interface của ngôn ngữ Go là loại liên kết trễ (delay binding), có thể hiện thực các chức năng đa hình như các  hàm ảo.

Các  interface có mặt khắp nơi trong ngôn ngữ Go. Trong ví dụ "Hello world", `fmt.Printf` là hàm có thiết kế hoàn toàn dựa trên  interface và chức năng thực sự của nó được `fmt.Fprintf` thực hiện bởi các hàm. Kiểu `error` được sử dụng để chỉ ra lỗi là  một kiểu  interface tích hợp. Trong C, `printf` chỉ cho phép một số lượng hạn chế các kiểu dữ liệu cơ bản có thể được in vào các đối tượng file. Tuy nhiên, nhờ tính năng  interface linh hoạt của Go mà `fmt.Fprintf` có thể in ra bất kỳ đối tượng output stream tùy chỉnh nào, in ra file hoặc output tiêu chuẩn, in ra mạng hoặc thậm chí in ra file nén. Đồng thời, dữ liệu in không bị giới hạn. Đối với các kiểu cơ bản được tích hợp vào ngôn ngữ, bất kỳđối tượng  `fmt.Stringer` nào hoàn toàn thỏa mãn  interface đều có thể được in. Nếu  interface của `fmt.Stringer` không được thỏa mãn , nó vẫn có thể được in bằng kỹ thuật reflection. Protorype của hàm `fmt.Fprintf`  như sau:

```go
func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)
```

Trong đó `io.Writer` là interface output, `error` là built-in interface làm việc với lỗi được định nghĩa như sau:

```go
type io.Writer interface {
    Write(p []byte) (n int, err error)
}

type error interface {
    Error() string
}
```

Chúng ta có thể output từng kí tự thành kí tự in hoa bằng cách tùy chỉnh lại đối tượng output của nó: [(source)](../examples/ch1/ch1.4/3-interface/example-1/main.go)

[>> mã nguồn](../examples/ch1/ch1.4/3-interface/example-1/main.go)

```go
type UpperWriter struct {
    io.Writer
}

func (p *UpperWriter) Write(data []byte) (n int, err error) {
    return p.Writer.Write(bytes.ToUpper(data))
}

func main() {
    fmt.Fprintln(&UpperWriter{os.Stdout}, "hello, world")
}
```

Tất nhiên ta cũng có thể định nghĩa định dạng in riêng để đạt được hiệu quả tương tự. Với mỗi đối tượng được in ra, nếu interface `fmt.Stringer` được thỏa mãn, kết quả kiểu `String` được trả về bởi phương thức của đối tượng được in mặc định:

[>> mã nguồn](../examples/ch1/ch1.4/3-interface/example-2/main.go)

```go
type UpperString string

func (s UpperString) String() string {
    return strings.ToUpper(string(s))
}

type fmt.Stringer interface {
    String() string
}

func main() {
    fmt.Fprintln(os.Stdout, UpperString("hello, world"))
}
```

Trong ngôn ngữ Go, chuyển đổi ngầm định không được hỗ trợ với các kiểu cơ bản (kiểu không có interface). Chúng ta không thể gán giá trị  của một kiểu `int` trực tiếp cho một biến  kiểu `int64`, chúng ta cũng không thể gán giá trị của kiểu `int` cho kiểu được đặt tên mới của kiểu cơ sở.

Các yêu cầu về tính nhất quán của ngôn ngữ Go đối với kiểu cơ bản là rất nghiêm ngặt, nhưng Go rất linh hoạt để chuyển đổi kiểu interface. Chuyển đổi giữa các đối tượng và interface, chuyển đổi giữa các interface và interface đều có thể là chuyển đổi ngầm định. Bạn có thể xem ví dụ sau:

```go
var (
    a io.ReadCloser = (*os.File)(f) // chuyển đổi ngầm định, *os.File thỏa  interface io.ReadCloser
    b io.Reader     = a             // chuyển đổi ngầm định, io.ReadCloser thỏa interface io.Reader
    c io.Closer     = a             // chuyển đổi ngầm định, io.ReadCloser thỏa interface io.Closer
    d io.Reader     = c.(io.Reader) // chuyển đổi tường minh, io.Closer 不thỏa interface io.Reader
)
```

Đôi khi đối tượng và interface quá linh hoạt dẫn đến việc chúng ta bị hạn chế vào việc bắt buộc phải sử dụng chúng. Một ví dụ phổ biến là định nghĩa một phương thức đặc biệt để phân biệt các interface. Ví dụ: interface `runtime` trong package `Error` xác định một phương thức duy nhất `RuntimeError` để chặn các kiểu khác vô tình điều chỉnh interface:

```go
type runtime.Error interface {
    error

    // RuntimeError is a no-op function but
    // serves to distinguish types that are run time
    // errors from ordinary errors: a type is a
    // run time error if it has a RuntimeError method.
    RuntimeError()
}
```

Trong protobuf, interface `Message`  cũng áp dụng một phương thức tương tự: định nghĩa một phương thức duy nhất `ProtoMessage` để ngăn các kiểu dữ liệu khác vô tình thỏa mãn interface:

```go
type proto.Message interface {
    Reset()
    String() string
    ProtoMessage()
}
```

`proto.Message` rất dễ  bị ai đó cố tình giả mạo interface. Một cách tiếp cận chặt chẽ hơn là xác định một phương thức riêng cho  interface. Chỉ các đối tượng thỏa mãn phương thức riêng  này mới có thể thỏa mãn interface đó và tên của phương thức riêng chứa tên đường dẫn tuyệt đối của package, vì vậy phương thức riêng này chỉ có thể được hiện thực bên trong package để đáp ứng  interface này. `testing.TB` interface trong gói thử nghiệm sử dụng một kỹ thuật tương tự:

```go
type testing.TB interface {
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    ...

    // A private method to prevent users implementing the
    // interface and so future additions to it will not
    // violate Go 1 compatibility.
    private()
}
```

Tuy nhiên, phương pháp chặn  các đối tượng bên ngoài thực hiện interface thông qua các phương thức private phải lưu ý:

- Thứ nhất, interface này chỉ có thể được sử dụng bên trong gói và các gói bên ngoài thường không thể tạo ra các đối tượng thỏa mãn interface,
- Thứ hai, cơ chế bảo vệ này cũng không phải tuyệt đối, người dùng nếu cố tình vẫn có thể bỏ qua được.

Như  đã đề cập trong phần Method, ta có thể kế thừa các phương thức của  kiểu ẩn danh bằng cách nhúng các thành phần thuộc  kiểu đó vào struct. Trong thực tế, thành phần ẩn danh này không nhất thiết phải là một kiểu bình thường, mà có thể một kiểu interface cũng được. Chúng ta có thể làm giả  phương thức private `testing.TB` bằng cách nhúng vào các interface ẩn danh, bởi vì các phương thức trong interface thuộc loại lazy binding và không thành vấn đề nếu phương thức `private` thực sự tồn tại ở compile-time. [(source)](../examples/ch1/ch1.4/3-interface/example-3/main.go)

[>> mã nguồn](../examples/ch1/ch1.4/3-interface/example-3/main.go)

```go
package main

import (
    "fmt"
    "testing"
)

type TB struct {
    testing.TB
}

func (p *TB) Fatal(args ...interface{}) {
    fmt.Println("TB.Fatal disabled!")
}

func main() {
    var tb testing.TB = new(TB)
    tb.Fatal("Hello, playground")
}
```

Kế thừa  bằng cách nhúng vào interface ẩn danh hoặc nhúng vào đối tượng con trỏ ẩn danh thực sự implement là một thừa kế ảo thuần túy. Ta chỉ kế thừa đặc tả được chỉ định bởi interface và phần hiện thực chỉ thực sự được đưa vào trong thời gian thực thi. Ví dụ: chúng ta có thể mô phỏng một plugin thực hiện gRPC:

```go
type grpcPlugin struct {
    *generator.Generator
}

func (p *grpcPlugin) Name() string { return "grpc" }

func (p *grpcPlugin) Init(g *generator.Generator) {
    p.Generator = g
}

func (p *grpcPlugin) GenerateImports(file *generator.FileDescriptor) {
    if len(file.Service) == 0 {
        return
    }

    p.P(`import "google.golang.org/grpc"`)
    // ...
}
```

Đối tượng kiểu `grpcPlugin`  được xây dựng phải thỏa mãn  interface `generate.Plugin` (trong package "github.com/golang/protobuf/protoc-gen-go/generator"):

```go
type Plugin interface {
    // Name identifies the plugin.
    Name() string
    // Init is called once after data structures are built but before
    // code generation begins.
    Init(g *Generator)
    // Generate produces the code generated by the plugin for this file,
    // except for the imports, by calling the generator's methods
    // P, In, and Out.
    Generate(file *FileDescriptor)
    // GenerateImports produces the import declarations for this file.
    // It is called after Generate.
    GenerateImports(file *FileDescriptor)
}
```

Hàm `GenerateImports` được sử dụng trong phương thức của kiểu `generate.Plugin` tương ứng với interface `p.P(...)` được hiện thực bởi `Init` đối tượng `generator.Generator`. `generator.Generator` này tương ứng với một kiểu cụ thể, nhưng nếu nó là một kiểu interface, chúng ta  có thể vượt truyền nó thẳng vào trong phần hiện thực.

Ngôn ngữ Go dễ dàng hiện thực các tính năng nâng cao như hướng đối tượng với duck-typing và kế thừa ảo thông qua sự kết hợp của một số tính năng đơn giản, điều này thực sự đáng kinh ngạc.
