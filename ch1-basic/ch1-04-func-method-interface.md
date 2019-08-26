# 1.4. Functions, Methods và Interfaces

Trong phần này chúng ta sẽ tìm hiểu cụ thể về các khái niệm cơ bản trong Golang: Function, Method và Interface.

## 1.4.1. Function

Hàm (function) là  thành phần cơ bản của chương trình. Các hàm trong ngôn ngữ Go có thể có tên hoặc ẩn danh (anonymous function):

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

Một hàm trong ngôn ngữ Go có thể có nhiều tham số và nhiều giá trị trả về. Cả tham số và giá trị trả về trao đổi dữ liệu  với hàm theo cách truyền vào giá trị (pass by value). Về mặt cú pháp, hàm cũng hỗ trợ số lượng tham số thay đổi, biến số lượng tham số phải là tham số cuối cùng và biến này phải là kiểu slice.

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

    // tương đương với lời gọi trực tiếp `Print(123, "abc")`
    Print(a...) // 123 abc

    // tương đương với lời gọi `Print([]interface{}{123, "abc"})`
    Print(a)    // [123 abc]
}

func Print(a ...interface{}) {
    fmt.Println(a...)
}
```

Cả tham số truyền vào và các giá trị trả về đều có thể được đặt tên:

```go
func Find(m map[int]int, key int) (value int, ok bool) {
    value, ok = m[key]
    return
}
```

### Defer trong Function

Lệnh `defer` trì hoãn việc thực thi hàm cho tới khi hàm bao ngoài nó return. Các đối số trong lời gọi defer được đánh giá ngay lặp tức nhưng lời gọi không được thực thi cho tới khi hàm bao ngoài nó return.

```go
func main() {
    defer fmt.Println("world")

    fmt.Println("hello")
}
// kết quả: hello world
```

Mỗi lời gọi `defer` được push vào stack và thực thi theo thứ tự ngược lại khi hàm bao ngoài nó kết thúc.

Ta thường sử dụng `defer` cho việc đóng hoặc giải phóng tài nguyên:

- Đóng file giống như `try-finally`:

    ```go
    func main() {
        f, err := os.Create("file")
        if err != nil {
            panic("cannot create file")
        }

        // chắc chắn file sẽ được close dù hàm có bị panic hay return
        defer f.Close()
        fmt.Fprintf(f,"hello")
    }
    ```

- Đóng file và xử lý panic giống như `try-catch-finally`:

    ```go
    func main() {
        defer func() {
            msg := recover()
            fmt.Println(msg)
        }()

        // . là folder hiện tại
        f, err := os.Create(".")
        if err != nil {
            panic("cannot create file")
        }
        defer f.Close()

        // không quan trọng chuyện gì xảy ra thì file cũng sẽ được close
        // để đơn giản nên ở đây bỏ qua bước kiểm ra close result
        fmt.Fprintf(f,"hello")
    }
    ```

- Cũng giống như block `finally` thì lời gọi defer cũng có thể làm cho kết quả trả về thay đổi:

    ```go
    func yes() (text string) {
        defer func() {
            text = "no"
        }()
        return "yes"
    }

    func main() {
        fmt.Println(yes())
    }
    ```

### Slice trong Function

<div align="center">
	<img src="../images/slice_1.png" width="300">
	<br/>
	<span align="center">
		<i>Minh hoạ slice</i>
	</span>
</div>

Mọi thứ trong Go đều được truyền theo kiểu pass by value, slice cũng thế. Nhưng vì giá trị của slice là một *header* (chứa con trỏ tới dữ liệu array bên dưới) nên khi truyền slice vào hàm, quá trình copy sẽ bao gồm luôn địa chỉ tới array chứa dữ liệu thực sự.

Ví dụ sau cho thấy ý nghĩa của việc truyền tham số kiểu slice vào hàm thay vì array:

```go
// truyền vào array sẽ giúp
// nội dung của biến x không bị thay đổi
func once(x [3]int) {
    for i := range x {
        x[i] *= 2
    }
}

// truyền vào con trỏ ngầm định (slice)
// khiến nội dung của biến x bị thay đổi
func twice(x []int) {
    for i := range x {
        x[i] *= 2
    }
}

func main() {
    data := [3]int{8,9,0}

    once(data)
    fmt.Println(data)

    twice(data[0:])
    fmt.Println(data)

    // kết quả:
    // [8 9 0]
    // [16 18 0]
}
```

### Tham số trả về được đặt tên

Cũng như tham số nhận vào, giá trị trả về cũng có thể được đặt tên, nhờ đó có thể đơn giản hoá lệnh return:

```go
func ReadFull(r Reader, buf []byte) (n int, err error) {
    for len(buf) > 0 && err == nil {
        var nr int
        nr, err = r.Read(buf)
        n += nr
        buf = buf[nr:]
    }

    // hàm trả về n mà không cần phải chỉ rõ
    return
}
```

## 1.4.2. Method

Go không có class, tuy nhiên chúng ta có thể định nghĩa các phương thức (Method) cho *type* (kiểu).

Phương thức là một hàm với đối số (argument) đặc biệt gọi là *receiver*.

```go
type Vertex struct {
    X, Y float64
}

// method Abs() với receiver 'v'
func (v Vertex) Abs() float64 {
    return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func main() {
    v := Vertex{3, 4}
    fmt.Println(v.Abs())

    // kết quả:
    // 5
}
```

Phương thức (Method) là một tính năng của lập trình hướng đối tượng (OOP). Trong ngôn ngữ C++, phương thức  tương ứng với một hàm thành viên của một class, liên kết với một  đối tượng cụ thể. Tuy nhiên, phương thức trong ngôn ngữ Go được liên kết với kiểu, do đó liên kết tĩnh của phương thức có thể được tạo thành trong giai đoạn biên dịch.

Một chương trình hướng đối tượng sử dụng các phương thức để thể hiện những thao tác trên thuộc tính (properties) của nó, qua đó người dùng có thể sử dụng đối tượng mà không cần phải thao tác trực tiếp với đối tượng mà là thông qua các phương thức. C++ thường được xem là một dấu mốc mà lập trình hướng đối tượng bắt đầu phát triển mạnh mẽ, nó hỗ trợ các tính năng hướng đối tượng (như class) dựa trên cơ sở ngôn ngữ C. Kế đến là Java, ngôn ngữ được gọi là hướng đối tượng thuần túy  vì các hàm của nó không thể tồn tại độc lập mà phải thuộc về một class nhất định.

Đối với một kiểu nhất định, tên của mỗi phương thức phải là duy nhất và các phương thức cũng như hàm đều không hỗ trợ overload.

Dưới đây là hiện thực các phương thức làm việc với File theo kiểu ngôn ngữ C:

```go
type File struct {
    fd int
}

// mở file
func OpenFile(name string) (f *File, err error) {
    fmt.Println("Opening file ", name)
    return nil, nil
}

// đóng file
func (f *File) Close() error {
    fmt.Println("Close file")
    return nil
}

// đọc dữ liệu từ file
func (f *File) Read(offset int64, data []byte) int {
    fmt.Println("Read file")
    return 0
}
```

Ta sử dụng các phương thức này như sau:

```go
func main() {
    var data []byte

    // khởi tạo một đối tượng File
    f, _ := OpenFile("foo.dat")

    f.Read(0, data)
    f.Close()
}
```

Trong một số tình huống, ta quan tâm nhiều hơn đến một chuỗi thao tác ví dụ  như `Read` đọc một số mảng và sau đó gọi `Close` để đóng, trong ngữ cảnh này, người dùng không quan tâm đến kiểu của đối tượng, miễn là nó có thể đáp ứng được các thao tác của `Read` và `Close`. Tuy nhiên trong các biểu thức phương thức của `ReadFile`, `CloseFile` có chỉ rõ kiểu `File` trong tham số kiểu sẽ khiến chúng không bị phụ thuộc vào đối tượng nào cụ thể. Việc này có thể khắc phục bằng cách sử dụng thuộc tính closure (closure property):

```go
func main() {
    var data []byte

    // khởi tạo một đối tượng File
    f, _ := OpenFile("foo.dat")

    // một hàm closure có thể gọi tới đối tượng f ngoài hàm
    // sẽ liên kết với đối tượng f
    var Close = func() error {
        return (*File).Close(f)
    }

    // tương tự với hàm Close
    var Read = func (offset int64, data []byte) int {
        return (*File).Read(f, offset, data)
    }

    // xử lý file
    Read(0, data)
    Close()
}
```

Chúng ta có thể đơn giản hóa thành như sau:

```go
func main() {
    var data []byte

    // mở đối tượng file
    f, _ := OpenFile("foo.dat")

    // ràng buộc với đối tượng f
    var Close = f.Close

    // ràng buộc với đối tượng f
    var Read = f.Read

    // khi gọi không cần chỉ rõ đối tượng nữa
    // vì đã được ràng buộc trước đó
    Read(0, data)
    Close()
}
```

### Kế thừa phương thức

Go không hỗ trợ tính năng kế thừa như các ngôn ngữ hướng đối tượng truyền thống mà có cách của riêng mình. Tính kế thừa đạt được bằng cách xây dựng các thuộc tính ẩn danh trong struct:

```go
type Point struct{ X, Y float64 }

type ColoredPoint struct {
    // thuộc tính ẩn danh
    Point

    // thuộc tính bình thường
    Color color.RGBA
}
```

Chúng ta có thể định nghĩa `ColoredPoint` như một struct có 3 trường, nhưng ở đây chúng ta sẽ dùng struct `Point` chứa `X` và `Y` để thay thế.

```go
// khai báo một đối tượng thuộc struct
var cp ColoredPoint

// có thể gán thẳng vào thuộc tính X
// không cần phải thông qua Point
cp.X = 1

// có thể truy cập X bằng cách này
fmt.Println(cp.Point.X)
// "1"

// hoặc gán vào Y thông qua Point
cp.Point.Y = 2

// và truy cập Y bằng cách này
fmt.Println(cp.Y)
// "2"
```

Có thể đạt được kết quả tương tự ngay cả với phương thức.

```go
// lấy ví dụ với struct Mutex có sẵn
type Mutex struct {}
func (m *Mutex) Lock()
func (m *Mutex) Unlock()

// struct Cache kế thừa Mutex bằng cách
// khai báo một thuộc tính ẩn danh là sync.Mutex
type Cache struct {
    m map[string]string
    sync.Mutex
}


// Lookup tìm trên Cache với dữ liệu key và trả về value tương ứng
func (p *Cache) Lookup(key string) string {
    // p có thể gọi thẳng tới phương thức Lock và Unlock
    // nhờ kế thừa từ sync.Mutex
    p.Lock()
    defer p.Unlock()

    return p.m[key]
}
```

Khả năng liên kết trực tiếp tới kiểu được kế thừa này được hoàn thành lúc biên dịch và không mất chi phí runtime.

Ví dụ trên có thể làm ta nghĩ rằng `sync.Mutex` là một lớp cơ sở và `Cache` là lớp kế thừa hoặc lớp con của nó. Tuy nhiên, phương thức được kế thừa theo cách này không thể hiện được tính đa hình bởi vì cái mà đối tượng `p` gọi tới là phương thức gốc mà không phải của nó (của struct Cache).

Nếu cần tính chất đa hình như các ngôn ngữ OOP khác, chúng ta cần triển khai nó với Interface.

## 1.4.3. Interface

Các interface trong Go cung cấp một cách để xác định hành vi của một đối tượng: nếu đối tượng đó có thể làm những việc *như thế này*, thì nó có thể được sử dụng *ở đây*.

Ngôn ngữ Go hiện thực mô hình hướng đối tượng thông qua cơ chế Interface.

Rob Pike, cha đẻ của ngôn ngữ Go, đã từng nói một câu nói nổi tiếng:

> Languages ​​that try to disallow idiocy become themselves idiotic

Các ngôn ngữ lập trình tĩnh nói chung có các hệ thống kiểu nghiêm ngặt, cho phép trình biên dịch đi sâu vào xem liệu lập trình viên có thực hiện bất kỳ động thái bất thường nào không. Tuy nhiên, một hệ thống kiểu quá nghiêm ngặt có thể làm cho việc lập trình trở nên quá cồng kềnh và khiến chúng ta phải mất nhiều thời gian cho nó.

Ngôn ngữ Go  vì thế cố gắng cung cấp sự cân bằng giữa sự linh hoạt và tính an toàn: có cơ chế `duck-typing` thông qua interface nhưng đồng thời cũng  kiểm tra kiểu nghiêm ngặt.

### Duck typing

Duck-typing với ý tưởng đơn giản:

> If something looks like a duck, swims like a duck and quacks like a duck then it’s probably a duck.

<div align="center">
	<img src="../images/duck-typing.png" width="300">
	<br/>
	<span align="center">
	</span>
</div>

Ví dụ có một interface con vịt, xác định khả năng `Quacks`:

```go
type Duck interface {
   Quacks()
}
```

Và cách ta áp dụng *duck-typing*:

```go
// một struct động vật bất kì
type Animal struct {
}

// con này có khả năng `Quacks` như vịt
func (a Animal) Quacks() {
   fmt.Println("The animal quacks");
}

// hàm dành cho vịt
func Scream(duck Duck) {
   duck.Quacks()
}

func main() {
    // a là một một vật thuộc struct Animal
   a := Animal{}

   // vì a có khẳng năng `Quacks` như vịt nên
   // ta có thể sử dụng nó như một con vịt trong hàm này
   Scream(a)
}
```

Thiết kế này cho phép chúng ta tạo ra một interface mới thỏa mãn kiểu hiện có mà không phải  hủy đi định nghĩa ban đầu của chúng, điều này đặc biệt linh hoạt và hữu ích khi các kiểu mà ta sử dụng đến từ những package không thuộc quyền kiểm soát của mình.

### Chuyển đổi kiểu trong Go

Trong Golang, chuyển đổi kiểu ngầm định không được hỗ trợ với các kiểu cơ bản (kiểu không có interface): không thể gán giá trị  của một biến kiểu `int` trực tiếp cho một biến  kiểu `int64`.

Các yêu cầu về tính nhất quán của ngôn ngữ Go đối với kiểu cơ bản nghiêm ngặt là thế, nhưng nó lại khá linh hoạt để chuyển đổi kiểu giữa các interface: Chuyển đổi giữa đối tượng - interface hoặc chuyển đổi giữa interface - interface đều có thể là chuyển đổi ngầm định. Bạn có thể xem ví dụ sau:

```go
var (
    // chuyển đổi ngầm định khi *os.File thỏa  interface io.ReadCloser
    a io.ReadCloser = (*os.File)(f)

    // chuyển đổi ngầm định khi io.ReadCloser thỏa interface io.Reader
    b io.Reader  = a

    // chuyển đổi ngầm định khi io.ReadCloser thỏa interface io.Closer
    c io.Closer  = a

    // chuyển đổi tường minh khi io.Closer thỏa interface io.Reader
    d io.Reader     = c.(io.Reader)
)
```

#### Một số sai lầm khi sử dụng Interface

Đôi khi đối tượng và interface quá linh hoạt dẫn đến việc chúng ta có thể mắc sai lầm khi struct khác vô tình điều chỉnh interface. Để khắc phục ta định nghĩa một phương thức đặc biệt để phân biệt các interface:

```go
type runtime.Error interface {
    error

    // RuntimeError là một hàm rỗng được dùng chỉ với mục đích là
    // phân biệt lỗi runtime  với các lỗi khác nhờ tính chất:
    // một type là runtime error chỉ khi nào nó có method RuntimeError
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

Interface `proto.Message` rất dễ  bị "giả mạo", để tránh điều đó ta nên định nghĩa một phương thức riêng cho nó. Chỉ các đối tượng thỏa mãn phương thức riêng  này mới có thể thỏa mãn interface đó và tên của phương thức riêng chứa tên đường dẫn tuyệt đối của package, vì vậy phương thức này chỉ có thể được hiện thực bên trong package để đáp ứng  interface. `testing.TB` là interface trong package `test` sử dụng kỹ thuật này:

```go
type testing.TB interface {
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    ...

    // Phương thức private ngăn user khác implement interface
    private()
}
```

#### Khả năng bị "làm giả" phương thức thuộc interface

Như  đã đề cập trong phần Method, ta có thể kế thừa các phương thức của  kiểu ẩn danh bằng cách thêm các thuộc tính ẩn danh thuộc  kiểu đó vào struct. Vậy điều gì xảy ra nếu thuộc tính ẩn danh này không phải là một kiểu bình thường  mà là một kiểu interface?

Chúng ta có thể làm giả  phương thức `private` của `testing.TB` bằng cách nhúng vào struct `TB` interface ẩn danh:

```go
package main

import (
    "fmt"
    "testing"
)

// TB có thể kế thừa phương thức `private` từ interface `testing.TB`
type TB struct {
    testing.TB
}

// phương thức thuộc struct TB
func (p *TB) Fatal(args ...interface{}) {
    fmt.Println("TB.Fatal disabled!")
}

func main() {
    // khởi tạo một đối tượng thuộc interface testing.TB
    var tb testing.TB = new(TB)

    // lúc này nó có thể sử dụng phương thức Fatal mà TB đã hiện thực
    tb.Fatal("Hello, playground")
}
```

## 1.4.4. Luồng thực thi của một chương trình Go

Việc khởi tạo và thực thi chương trình Go luôn bắt đầu từ hàm `main.main`. Nếu package `main` có import  các package khác, chúng sẽ được thêm vào package `main` theo thứ tự khai báo.

`init` không phải là hàm thông thường, nó có thể có nhiều định nghĩa, và các hàm khác không thể sử dụng nó.

- Nếu một package được import nhiều lần, sẽ chỉ được tính là một khi thực thi.
- Khi một package được import mà nó lại import các package khác, trước tiên Go sẽ import các package khác đó trước, sau đó  khởi tạo các hằng và biến của package, rồi gọi hàm `init` trong từng package.
- Nếu một package có nhiều hàm `init` và thứ tự gọi không được xác định cụ thể thì chúng sẽ được gọi theo thứ tự xuất hiện. Cuối cùng, khi `main` đã có đủ tất cả hằng và biến ở package-level thì nó sẽ được khởi tạo bằng cách thực thi hàm `init`, tiếp theo chương trình đi vào hàm `main.main` và  bắt đầu thực thi. Hình dưới đây là sơ đồ nguyên lý  một chuỗi bắt đầu của chương trình hàm trong Go:

<div align="center">
<img src="../images/ch1-11-init.ditaa.png">
<br/>
<span align="center"><i>Tiến trình khởi tạo package</i></span>
</div>

Cần lưu ý rằng trước khi hàm nào khác được thực thi thì tất cả code đều chạy trong cùng một Goroutine `main.main`, đây là thread chính của chương trình. Do đó, nếu một Goroutine khởi chạy trong hàm `main.main` thì nó chỉ có thể được thực thi sau khi vào chương trình đã thực thi xong `init`.

[Tiếp theo](ch1-05-concurrency-parallelism.md)