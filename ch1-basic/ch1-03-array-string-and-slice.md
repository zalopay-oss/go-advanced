
## 1.3 Array, strings và slices

`Arrays` và một số cấu trúc dữ liệu liên quan khác được sử dụng thường xuyên trong các ngôn ngữ lập trình. Chỉ khi chúng không đáp ứng được yêu cầu chúng ta mới cân nhắc sử dụng `linked lists` (danh sách liên kết) và `hash tables` (bảng băm) hoặc nhiều cấu trúc dữ liệu tự định nghĩa phức tạp khác.

`Arrays`, `strings` và `slices` trong ngôn ngữ Go là các cấu trúc dữ liệu liên quan mật thiết với nhau. Ba kiểu dữ liệu đó có cùng cấu trúc vùng nhớ lưu trữ bên dưới, và chỉ có những hành vi thể hiện ra bên ngoài khác nhau tùy thuộc vào ràng buộc ngữ nghĩa. Đầu tiên, trong ngôn ngữ Go, `array` là một kiểu giá trị. Mặc dù những phần tử của array có thể được chỉnh sửa, phép gán của array hoặc khi truyền array như là một tham số của hàm thì chúng sẽ được xử lý toàn bộ, có thể hiểu là khi đó chúng được sao chép lại toàn bộ thành một bản sao rồi mới xử lý trên bản sao đó - khác với kiểu truyền tham khảo. Bên dưới dữ liệu của ngôn ngữ Go, `string` cũng là một array của các `byte` dữ liệu, nhưng khác với array những phần tử của string không được phép chỉnh sửa. Phép gán string chỉ đơn giản là sao chép hai thành phần đó là con trỏ tới vùng nhớ của `string` và độ dài `string`, mà không phải sao chép toàn bộ string. `Slices` thì phức tạp hơn, cấu trúc của chúng cũng như `string`, tuy nhiên việc giới hạn chỉ-đọc như string được lược bỏ, mỗi slice có thêm hai thông tin là `len` (độ dài) và `capacity` (sức chứa). Phép gán của slice và khi truyền slice như tham số của hàm thì thông tin về header của slice sẽ được xử lý theo giá trị. Bởi vì slice header chứa con trỏ đến dữ liệu bên dưới, phép gán sẽ không gây ra việc sao chép toàn bộ dữ liệu. Trong thực tế, phép gán trong Go và quy luật truyền tham số hàm trong Go rất đơn giản. Ngoại trừ hàm `closure` có tham khảo tới biến toàn cục bên ngoài, thì hầu hết những phép gán và truyền tham số khác được truyền bằng giá trị. Để hiểu được ba cách để xử lý arrays, strings và slices cần phải hiểu chi tiết tầng lưu trữ bên dưới của chúng.


### 1.3.1 Array

Một array là một chuỗi độ dài cố định của các phần tử có kiểu dữ liệu nào đó, một array có thể bao gồm không hoặc nhiều phần tử. Độ dài của array là một phần thông tin được chứa trong nó, các array có độ dài khác nhau hoặc kiểu phần tử bên trong khác nhau được xem là các kiểu dữ liệu khác nhau, và không được phép gán cho nhau, vì thế array hiếm khi được sử dụng trong Go. Một kiểu dữ liệu tương ứng với array là slice, một slice cũng là một chuỗi nhưng có thể tăng giảm kích thước một cách động, và các hàm hỗ trợ kiểu slice thì rất linh hoạt, nhưng để hiểu slice hoạt động thế nào, chúng ta phải hiểu array.

Đầu tiên, xem cách định nghĩa một array.

```go
var a [3]int  // Định nghĩa một mảng kiểu int độ dài 3, các phần tử đều bằng 0
var b = [...]int{1, 2, 3} // Định nghĩa một mảng có ba phần tử 1, 2, 3, do đó độ dài là 3
var c = [...]int{2: 3, 1: 2} // Mảng này có 3 phần tử theo thứ tự là 0, 2, 3
var d = [...]int{1, 2, 4: 5, 6} // Mảng này chứa dãy các phần tử là 1, 2, 0 , 0, 5, 6
```


[>> mã nguồn](../examples/chapter1/ch1.3/1-arrays/example-1/main.go)

Cách đầu tiên là cách cơ bản nhất để định nghĩa một array. Độ dài của array sẽ được ràng buộc trước, và mỗi phần tử trong array sẽ được khởi tạo với giá trị ban đầu là 0.

Cách thứ hai cũng được dùng để định nghĩa một array, chúng ta có thể đặc tả những giá trị được khởi tạo cho nó theo thứ tự được định nghĩa. Chiều dài của array được tự động tính toán theo số lượng phần tử được khởi tạo.

Cách thứ ba khi khởi tạo array sẽ kèm theo chỉ mục của từng phần tử, theo đó thứ tự xuất hiện của các phần tử  trong array ứng với thứ tự chỉ mục của nó. Cách khởi tạo này tương tự với kiểu `map[int]Type`. Độ dài của array sẽ dựa trên chỉ mục lớn nhất từng xuất hiện, và những phần tử ứng với chỉ mục không được khai báo sẽ có giá trị bằng 0.

Cách thứ tư, là pha trộn giữa cách thứ hai và thứ ba, hai phần tử đầu tiên trong ví dụ sẽ được khởi tạo tuần tự, phần tử thứ ba và thứ tư sẽ được khởi tạo với giá trị 0, phần tử thứ năm được khởi tạo theo chỉ mục đã cho, và phần tử cuối cùng theo sau phần tử thứ 5 cũng sẽ được khởi tạo theo thứ tự này.

Cấu trúc vùng nhớ của array thì rất đơn giản. Ví dụ cho một array `[4]int{2,3,5,7}` thì cấu trúc bên dưới sẽ như sau:

<p align="center" width="600">
<img src="../images/ch1-7-array-4int.ditaa.png">
<br/>
<span>Hình 1-7 Array layout</span>
</p>

Array trong ngôn ngữ Go mang ngữ nghĩa giá trị. Biến thể hiện array được xem như là toàn bộ array. Nó không phải là một con trỏ ngầm định tới phần tử đầu tiên (như trong ngôn ngữ C), mà hoàn toàn là một giá trị. Khi biến array được gán hoặc truyền, thì toàn bộ array sẽ được sao chép. Nếu kích thước của array lớn, thì phép gán array sẽ chịu tổn phí lớn. Để tránh việc `overhead` (tổn phí) trong việc sao chép array, bạn có thể truyền con trỏ tới array, lưu ý con trỏ array thì không phải là một array.

```go
var a = [...]int{1, 2, 3} // a là một array
var b = &a                // b là một con trỏ tới array a

fmt.Println(a[0], a[1])   // in ra hai phần tử đầu tiên của array a
fmt.Println(b[0], b[1])   // truy xuất các phần tử của con trỏ array cũng giống như truy xuất các phần tử của array

for i, v := range b {     // duyệt qua các phần tử trong con trỏ array, giống như duyệt qua array
    fmt.Println(i, v)
}
```

[>> mã nguồn](../examples/chapter1/ch1.3/1-arrays/example-2/main.go)

Trong khi `b` là một con trỏ tới array `a`, nhưng khi làm việc với `b` cũng giống như `a`. Thì hoàn toàn có thể lặp qua (dùng `for range`) đối với con trỏ array, khi chúng ta gán hoặc truyền vào hàm một con trỏ array thì chỉ có giá trị con trỏ array được sao chép. Tuy nhiên con trỏ array cũng không đủ linh hoạt, bởi vì thông tin về chiều dài của array là một phần của array, do đó nếu hai con trỏ tới hai array có độ dài khác nhau thì hai con trỏ đó cũng thuộc kiểu khác nhau.

Bạn có thể nghĩ về array là một cấu trúc dữ liệu đặc biệt. Các trường trong cấu trúc sẽ tương ứng với chỉ mục của array, và số lượng phần tử của cấu trúc đó cũng được cố định.  Hàm `len` được dựng sẵn có thể dùng để lấy thông tin về độ dài của array, và hàm `cap` sẽ tính toán sức chứa của array. Nhưng trong kiểu array cả hai hàm này sẽ cùng trả về một giá trị giống nhau (điều này khác với slice).

Chúng ta có thể  dùng vòng lặp `for` dể duyệt qua các phần tử của array. Sau đây là những cách thường dùng để duyệt qua một array

```go
for i := range a {
    fmt.Printf("a[%d]: %d\n", i, a[i])
}
for i, v := range b {
    fmt.Printf("b[%d]: %d\n", i, v)
}
for i := 0; i < len(c); i++ {
    fmt.Printf("c[%d]: %d\n", i, c[i])
}
```


[>> mã nguồn](../examples/chapter1/ch1.3/1-arrays/example-3/main.go)

`for range` là cách tốt nhất để duyệt qua các phần tử trong array, bởi vì cách này sẽ đảm các việc truy xuất sẽ không vượt quá giới hạn của array.

Khi dùng `for range`, chúng ta có thể phớt lờ đi các tham số  đi kèm bằng cách sau

```go
var times [5][0]int
for range times {
    fmt.Println("hello")
}
```


[>> mã nguồn](../examples/chapter1/ch1.3/1-arrays/example-4/main.go)

Biến `times` sẽ tương ứng với kiểu array `[5][0]int`, mặc dù chiều thứ nhất của array có độ dài là 5, nhưng độ dài của array `[0]int` là 0, do đó kích thước của toàn bộ `array` là 0. Bỏ qua chi phí cho việc khởi tạo vùng nhớ chúng ta sẽ thực hiện 5 vòng lặp nhanh chóng.

Các phần tử của array không nhất thiết là kiểu số học, nên cũng có thể là string, struct, function, interface, và pipe, v,v..

```go
// Mảng string
var s1 = [2]string{"hello", "world"}
var s2 = [...]string{"Hello!", "World"}
var s3 = [...]string{1: "Hello", 0: "World", }

// Mảng struct
var line1 [2]image.Point
var line2 = [...]image.Point{image.Point{X: 0, Y: 0}, image.Point{X: 1, Y: 1}}
var line3 = [...]image.Point{{0, 0}, {1, 1}}

// Mảng decoder của hình ảnh
var decoder1 [2]func(io.Reader) (image.Image, error)
var decoder2 = [...]func(io.Reader) (image.Image, error){
    png.Decode,
    jpeg.Decode,
}

// Mảng interface{}
var unknown1 [2]interface{}
var unknown2 = [...]interface{}{123, "Hello!"}

// Mảng pipe
var chanList = [2]chan int{}
```

[>> mã nguồn](../examples/chapter1/ch1.3/1-arrays/example-5/main.go)


Chúng ta cũng có thể định nghĩa một array rỗng

```go
var d [0]int       // Định nghĩa một array chiều dài 0
var e = [0]int{}   // Tương tự trên
var f = [...]int{} // Tương tự như trên
```


[>> mã nguồn](../examples/chapter1/ch1.3/1-arrays/example-6/main.go)

Một array có chiều dài 0 thì không chiếm không gian lưu trữ. Một mảng rỗng hiếm khi được sử dụng trực tiếp, có có ích trong trường hợp như sau để đồng bộ luồng thực thi, khi mà việc phát sinh thêm vùng nhớ là không thực sự cần thiết

```go
c1 := make(chan [0]int)
go func() {
    fmt.Println("c1")
    c1 <- [0]int{}
}()
<-c1
```


[>> mã nguồn](../examples/chapter1/ch1.3/1-arrays/example-7/main.go)

Ở đây, chúng ta không quan tâm về kiểu thực sự được truyền vào pipeline, trong khi thực thi lệnh nhận hoặc gửi chỉ nhằm mục đích đồng bộ thông điệp. Trong ngữ cảnh đó, chúng ta có thể sử dụng mảng rỗng trong pipe để hạn chế phí tổn của phép gán pipe. Dĩ nhiên, nó thích hợp hơn khi thay thế bằng một kiểu struct vô danh.

```go
c2 := make(chan struct{})
go func() {
    fmt.Println("c2")
    c2 <- struct{}{}
}()
<-c2
```


[>> mã nguồn](../examples/chapter1/ch1.3/1-arrays/example-8/main.go)

Chúng ta có thể sử dụng hàm `fmt.Printf`, chúng cho phép in ra kiểu cũng như chi tiết của array thông qua các chỉ thị `%T` hoặc `%#v`

```go
fmt.Printf("b: %T\n", b)  // b: [3]int
fmt.Printf("b: %#v\n", b) // b: [3]int{1, 2, 3}
```


[>> mã nguồn](../examples/chapter1/ch1.3/1-arrays/example-9/main.go)

Trong Go, kiểu array là một kiểu cơ bản như là slice và strings. Nhiều ví dụ về array phía trên có thể được áp dụng trực tiếp cho strings hoặc slices


### 1.3.2 String


Một string là một chuỗi các giá trị `byte` không được thay đổi, và string thường được dùng để biểu diễn giá trị con người có thể đọc được. Không giống như array, những phần tử trong string sẽ không được thay đổi, và chỉ có thể đọc. Chiều dài của mỗi string sẽ được cố định, những thông tin chiều dài đó không là một phần của kiểu string. Do mã nguồn của Go được yêu cầu là kiểu `UTF8`. Nội dung của string trong mã nguồn với kiểu Unicode sẽ được chuyển thành UTF8. Bởi vì mỗi phần tử của string cũng thực chất được lưu trữ thành những byte chỉ-đọc, một string có thể chứa những dữ liệu tùy ý, có thể toàn những byte zero (không). Chúng ta có thể dùng string để biểu diễn kiểu không phải là UTF8 bằng cách mã hóa chúng như là GBK, nhưng cơ bản không nên làm như vậy bởi vì hàm mệnh đề `for range` trong Go không hỗ trợ duyệt string mang kí tự không phải kiểu UTF8.


Bên dưới cấu trúc string của ngôn ngữ Go `reflect.StringHeader` được định nghĩa với 

```go
type StringHeader struct {
    Data uintptr
    Len  int
}
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-1/main.go)

Cấu trúc của string chứa hai phần thông tin: đầu tiên là con trỏ array tới địa chỉ chứa string, thứ hai là chiều dài của string. Một string thực sự là một cấu trúc, do đó phép gán string thực chất là việc sao chép cấu trúc `reflect.StringHeader`, và không gây ra việc sao chép bên dưới phần dữ liệu. `[2]string`, cấu trúc bên dưới string được đề cập ở chương trước là `[2]reflect.StringHeader` cũng giống với cấu trúc dưới đây. 

Chúng ta có thể thấy cấu trúc vùng nhớ tương ứng với dòng string "Hello, world" là 

<p align="center" width="600">
<img src="../images/ch1-8-string-1.ditaa.png">
<br/>
<span>Hình 1-8 String layout</span>
</p>

Phân tích ra chúng ta có thể thấy rằng bên dưới dòng chữ "Hello, world" trong string chính xác là một array như sau

```go
var  data = [...] byte {
    'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd',
}
```

[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-1/main.go)


Mặc dù string không phải là slice nhưng nó cũng hỗ trợ thao tác (slicing) cắt. Một vài phần của vùng nhớ cũng được truy cập bên dưới slice tại một số nơi khác nhau.

```go
s := "hello, world"
hello := s[:5]
world := s[7:]
s1 := "hello, world"[:5]
s2 := "hello, world"[7:]
```

[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-1/main.go)


Tương tự như array, String cũng có một hàm dựng sẵn là `len` dùng để trả về chiều dài của string, ngoài ra bạn có thể  dùng `reflect.StringHeader` để truy xuất chiều dài của string theo cách như sau

```go
fmt.Println("len(s): ", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len)
fmt.Println("len(s1): ", (*reflect.StringHeader)(unsafe.Pointer(&s1)).Len)
fmt.Println("len(s2): ", (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len)
```

[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-1/main.go)


Theo như mô tả của ngôn ngữ Go, mã nguồn của ngôn ngữ được encoded (mã hóa) dưới dạng UTF8. Do đó, hằng string cũng được mã hóa dưới dạng UTF8. Khi đề cập tới Go string, chúng ta thường giả định rằng string là tương ứng với một chuỗi kí tự UTF8 hợp lệ. Bạn có thể dùng hàm dựng sẵn là `print` hoặc `fmt.Print` để in trực tiếp nó, hoặc có thể dùng vòng lặp `for range` qua chuỗi UTF8 một cách trực tiếp.

Chuỗi "Hello, 世界" chứa kí tự tiếng Trung có thể được in ra

```go
fmt.Printf("%#v\n", []byte("Hello, 世界"))

// Kết quả là
[]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c}
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-2/main.go)


Phân tích ra chúng ta có thể nhận thấy rằng các số hexa `0xe4, 0xb8, 0x96` ứng với từ "World" trong tiếng Trung, và `0xe7, 0x95, 0x8c` ứng với "Hello"

```go
fmt.Println("\xe4\xb8\x96")
fmt.Println("\xe7\x95\x8c")
// Kết quả là
世
界
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-3/main.go)


<p align="center" width="600">
<img src="../images/ch1-9-string-2.ditaa.png">
<br/>
<span>Hình 1-9 String layout</span>
</p>


Vì phần tử của string có thể  là những byte nhị phân, nên có thể bắt gặp một số trường hợp các kí tự UTF8 sẽ không được mã hóa chuẩn xác. Nếu bạn phát hiện được trường hợp nào mà UTF8 không encoded (mã hóa) đúng, một kí tự Unicode đặt biệt sẽ được in ra là 'uFFFD'. Kí tự này sẽ trông khác nhau ở những phầm mềm khác nhau. Thường thì kí tự này là một hình tứ giác hoặc kim cương màu đen, ở giữa chứa dấu hỏi.

Trong chuỗi sau, chúng ta sẽ cố tình làm hỏng byte thứ hai và thứ ba của ký tự đầu tiên, do đó, ký tự đầu tiên sẽ được in là "�", byte thứ hai và thứ ba sẽ bị phớt lờ, tiếp theo là "abc" Vẫn có thể giải mã in bình thường (mã hóa lỗi không lan truyền ngược là một trong những tính năng tuyệt vời của mã hóa UTF8)

```go
fmt.Println("\xe4\x00\x00\xe7\x95\x8cabc") // �界abc
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-4/main.go)


Tuy nhiên, khi mà `for range` trên những chuỗi UTF8 bị hỏng như trên, các byte thứ hai và thứ ba của kí tự đầu tiên vẫn sẽ được lặp lại một cách độc lập, nhưng giá trị của lần lặp này là 0 sau khi bị gặp lỗi.

```go
for i, c := range "\xe4\x00\x00\xe7\x95\x8cabc" {
	fmt.Println(i, c)
}
// 0 65533  // \uFFFD, 对应 �
// 1 0      // 空字符
// 2 0      // 空字符
// 3 30028  // 界
// 6 97     // a
// 7 98     // b
// 8 99     // c
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-5/main.go)


Nếu bạn không muốn decode (giải mã) chuỗi UTF8 và muốn duyệt trực tiếp qua nó, bạn có thể bắt string có thể chuyển qua chuỗi `[]byte` sau đó sẽ duyệt (sự chuyển đổi này sẽ không gây ra phí tổn khi chạy chương trình)  

```go
for i, c := range []byte("世界abc") {
    fmt.Println(i, c)
}
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-6/main.go)


Hoặc bạn có thể duyệt một dãy các byte của string như sau


```go
const s = "\xe4\x00\x00\xe7\x95\x8cabc"
for i := 0; i < len(s); i++ {
    fmt.Printf("%d %x\n", i, s[i])
}
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-7/main.go)


Hơn nữa, `for range` sẽ nhờ vào cú pháp UTF8 mà Go có thể hỗ trợ kiểu đặc biệt `[]rune` để chuyển từ kiểu string sang kiểu khác.

```go
fmt.Printf("%#v\n", []rune("世界"))      // []int32{19990, 30028}
fmt.Printf("%#v\n", string([]rune{'世', '界'})) // 世界
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-8/main.go)


Từ kết quả của đoạn mã nguồn trên, chúng ta có thể thấy `[]rune` thực sự là kiểu `[]int32`, từ đây, `rune` là một tên gọi khác của `int32`, `rune` được dùng để biểu diễn mỗi điểm unicode, hiện tại thì chỉ 21 bits được sử dụng.

Ép kiểu cho string sẽ liên quan đến hai kiểu dữ liệu `[]byte` và `[]rune`. Mỗi cách chuyển đổi sẽ kèm theo chi phí để cấp phát lại vùng nhớ, và trong trường hợp tệ hơn, thời gian tính toán sẽ xấp xỉ `O(n)`. Tuy nhiên, `[]rune` string thì đặc biệt hơn, bởi vì thông thường chỉ cần cast hai biến, thì cấu trúc vùng nhớ bên dưới phải đồng nhất hết sức có thể. Hiển nhiên, bên dưới kiểu `[]byte` và `[]int32` có thể hoàn toàn khác so với lớp trung gian, do đó sự chuyển đổi kiểu thực chất là chuyển đổi vùng nhớ.

Theo ví dụ sau, mã giả có thể được sử dụng để minh họa cơ bản một số tác vụ được xây dựng sẵn trong Go cho chuỗi, do đó độ phức tạp về thời gian và không gian của mỗi tác vụ sẽ dễ hiểu hơn.

**`for range` trong khi duyệt giả lập trong string**

```go
func forOnString(s string, forBody func(i int, r rune)) {
	for i := 0; len(s) > 0; {
		r, size := utf8.DecodeRuneInString(s)
		forBody(i, r)
		s = s[size:]
		i += size
	}
}
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-9/main.go)


`for range` khi lặp qua một string, mỗi lần chúng ta decode một ký tự Unicode và sau đó nhập vào thân vòng lặp for khi bắt gặp một kí tự broken code sẽ không gây dừng vòng lặp.

**`[]byte(s)` hiện thực mô phỏng chuyển đổi**

```go
func str2bytes(s string) []byte {
	p := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		p[i] = c
	}
	return p
}
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-10/main.go)


Một slice mới sẽ được tạo ra trong mô phỏng và sau đó một array của string sẽ được sap chép thành một slice theo từng phần tử, theo thứ tựu để đảm bảo ngữ nghĩa của string là chỉ đọc, Dĩ nhiên, khi chúng ta chuyển một string sang một array các byte `[]byte`, nếu trong quá trình chuyển đổi không thay đổi dữ liệu, thì bộ biên dịch sẽ trả về dữ liệu trực tiếp trỏ tới chuỗi gốc.

**`string(bytes)` mô phỏng hiện thực chuyển đổi kiểu**

```go
func bytes2str(s []byte) (p string) {
	data := make([]byte, len(s))
	for i, c := range s {
		data[i] = c
	}

	hdr := (*reflect.StringHeader)(unsafe.Pointer(&p))
	hdr.Data = uintptr(unsafe.Pointer(&data[0]))
	hdr.Len = len(s)

	return p
}
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-11/main.go)


Bởi vì string trong ngôn ngữ Go là chỉ-đọc, hoàn toàn không thể cấu trúc bên dưới một mảng kiểu byte để sinh ra một string. Để mô phỏng cách hiện thực, `unsafe` - một cấu trúc dữ liệu bên dưới của string sẽ được chứa trong một package, và sau đó một dữ liệu slice sẽ được sao chép thành chuỗi string tuần tự, nó giúp đảm bảo rằng ngữ nghĩa của string là chỉ được đọc không bị ảnh hưởng bởi slice. Nếu trong khi chuyển đổi, chuỗi byte không bị thay đổi trong suốt thời gian tồn tại của biến gốc, trình biên dịch sẽ xây dựng một mảng các `[]byte` để tạo thành string một cách trực tiếp dựa vào dữ liệu bên dưới.

**`[]rune(s)` Hiện thực mô phỏng chuyển đổi**

```go
func str2runes(s []byte) []rune {
	var p []int32
	for len(s) > 0 {
		r, size := utf8.DecodeRune(s)
		p = append(p, int32(r))
		s = s[size:]
	}
	return []rune(p)
}
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-12/main.go)


Bởi vì sự khác nhau bên dưới cấu trúc dữ liệu bên dưới, một string được chuyển đổi sang `[]rune` sẽ không thể không cấp phát lại vùng nhớ, và sau đó một chuỗi được decode và sao chép tuần tự tương tứng với chuỗi Unicode. Sự ép kiểu đó sẽ không có một sự tối ưu về string và bytes như được đề cập từ trước

**`string(runes)` mô phỏng chuyển đổi kiểu**

```go
func runes2string(s []int32) string {
	var p []byte
	buf := make([]byte, 3)
	for _, r := range s {
		n := utf8.EncodeRune(buf, r)
		p = append(p, buf[:n]...)
	}
	return string(p)
}
```

[>> mã nguồn](../examples/chapter1/ch1.3/2-strings/example-13/main.go)



Cũng bởi vì một sự khác nhau bên dưới cấu trúc lưu trữ, `[]rune`, việc chuyển đổi thành một chuỗi chắc chắn sẽ dẫn đến việc xây dựng lại chuỗi. Cách này không không có tối ưu hóa như mô tả ở trên.

### 1.3.3 Slice

Đơn giản mà nói, slice là một phiên bản đơn giản của mảng động. Bởi vì chiều dài của một mảng động không được cố định, chiều dài của slice thông thường không là một phần của kiểu dữ liệu. Array có nơi mà nó được áp dụng, nhưng kiểu array và những tác vụ trên nó sẽ không đủ linh hoạt, do đó array không được sử dụng nhiều trong ngôn ngữ Go. Slice thường được dùng một cách phổ biến hơn, và hiểu được ý nghĩa cũng như nguyên tắc sử dụng slice sẽ đòi hỏi phải có nhiều kĩ năng của người lập trình viên Go.

Hãy nhìn vào việc định nghĩa cấu trúc bên dưới slice `reflect.SliceHeader`

```go
type  SliceHeader  struct {
	Data  uintptr 
	Len   int 
	Cap   int 
}
```


[>> mã nguồn](../examples/chapter1/ch1.3/2-slices/example-1/main.go)


Có thể nhìn thấy rằng khởi đầu một slice là giống như Go String, nhưng slice có thêm thuộc tính `Cap` chỉ ra kích thước tối đa mà vùng nhớ trỏ tới slice được cấp phát. Hình bên dưới sẽ mô phỏng với `x := []int{2,3,5,7,11}` và `y := x[1:3]` cấu trúc vùng nhớ tương ứng với chương thứ hai.


<p align="center" width="600">
<img src="../images/ch1-10-slice-1.ditaa.png">
<br/>
<span>Hình 1-10 Slice layout</span>
</p>

Hãy nhìn vào định nghiã slices bên dưới:

```go
var (
	a [] int                // nil slice, equal to nil, generally used to represent a non-existent slice 
	b = []int{}            // empty slice, not equal to nil, generally used to represent an empty set 
	c = [] int {1,2,3}     // There are 3 elements of the slice, both len and cap are 3 
	d = c[:2]              // There are 2 elements of the slice, len is 2, cap is 3 
	e = c[0:2:cap(c)]      // There are 2 elements of the slice, len is 2, cap is 3 
	f = c[:0]              // There are 0 elements of the slice, len is 0, cap is 3 
	g = make ([]int,3)     // There are 3 elements of the slice, len and cap are 3 
	h = make ([]int,2,3) // there are 2 elements of the slice, len is 2, cap is 3 
	i = make ([]int,0,3) // There are 0 elements of the slice, len is 0, cap is 3 
)
```


[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-1/main.go)



Giống như array, hàm dựng sẵn `Len` sẽ trả về chiều dài của những phần tử hợp lệ trong khoảng `Cap` của slice. Hàm `Cap` được dựng sẵn sẽ trả về kích thước của slice, nó sẽ lớn hơn hay bằng với chiều dài của slice. Có thể dùng `reflect.SliceHeader` để truy cập vào thông tin trong slice thông qua cấu trúc đó (không khuyến khích). Slice khi `nil` có thể được so sánh với slice khác, và bản thân slice `nil` chỉ khi con trỏ của slice trỏ tới vùng nhớ rỗng. Khi đó, chiều dài của slice và sức chứa của nó sẽ không hợp lệ. Nếu con trỏ dữ liệu bên dưới slice là rỗng nhưng chiều dài hoặc sức chứa của nó khác 0, thì bản thân slice đó có thể bị corrupted (sai sót) (ví dụ slice bị sửa đổi sai bởi gói `reflect.SliceHeader` hoặc trực tiếp là `unsafe`).

Duyệt qua slice thì tương tự như duyệt qua một arrays

```go
for i := range a {
	fmt.Printf("a[%d]: %d\n", i, a[i])
}
for i, v := range b {
	fmt.Printf("b[%d]: %d\n", i, v)
}
for i := 0; i < len(c); i++ {
	fmt.Printf("c[%d]: %d\n", i, c[i])
}
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-2/main.go)

Trên thực tế, phép duyệt sẽ thông qua con trỏ dữ liệu bên dưới, chiều dài và sức chứa của slice sẽ không bị thay đổi, phép duyệt slice sẽ đọc và thay đổi phần tử như là array. Khi gán một giá trị hoặc truyền vào một tham số cho bản thân slice, nó hoạt động giống như array các con trỏ chỉ sao chép phần thông tin header của slice (`reflect.SliceHeader`). Ở các kiểu đó, điểm khác biệt lớn nhất đối với array và slice chính là thông tin về chiều dài, bên cạnh đó, slice có cùng kiểu dữ liệu sẽ có cùng kiểu slice.

Như đã đề cập ở trên, slicing là một phiên bản đơn giản của mảng động, nó là cấu trúc của kiểu slice. Bên cạnh đó khi xây dựng một slice hoặc duyệt qua slice, thêm phần tử vào slice hoặc xóa phần tử ra khỏi slice là những tác vụ thường gặp trên slice.

**Thêm phần tử vào slice**

Hàm dựng sẵn `append` có thể thêm phần tử thứ `N` vào cuối cùng của slice:

```go
var a []int
a = append(a, 1)               // nối thêm phần tử 1
a = append(a, 1, 2, 3)         // nối thêm phần tử 1, 2, 3
a = append(a, []int{1,2,3}...) // nối thêm các phần tử 1, 2, 3 bằng cách truyền vào một mảng
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-3/main.go)

Tuy nhiên, chú ý rằng trong trường hợp không đủ sức chứa, hàm `append` sẽ gây ra kết quả là vùng nhớ sẽ được phân bố lại, nó dẫn đến chi phí của việc phân bố và sao chép là rất lớn. Mặc dù khi sức chứa không đủ, bạn sẽ cần hàm `append` để cập nhật lại bản thân slice và là giá trị được trả về bởi hàm, bởi vì chiều dài của slice mới đã bị thay đổi.

Bên cạnh thêm phần tử vào cuối slice, chúng ta cũng có thể thêm phần tử vào đầu slice như sau

```go
var a = []int{1,2,3}
a = append([]int{0}, a...)        // thêm phần tử 0 vào đầu slice a
a = append([]int{-3,-2,-1}, a...) // thêm các phần tử -3, -2, -1 vào đầu slice a
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-4/main.go)

Đầu tiên, việc thêm phần tử vào đầu slice sẽ gây ra việc tổ chức lại vùng nhớ, nó cũng sẽ làm những phần tử đang tồn tại trong slice sẽ được sao chép một lần nữa. Do đó, hiệu suất của việc thêm phần tử  vào đầu slice sẽ tệ hơn là thêm phần tử vào cuối slice.

Do hàm `append` sẽ trả về một slice mới, nó sẽ hỗ trợ một dãy các tác vụ. Chúng ta có thể kết hợp nhiều hàm `append` để chèn một vài phần tử vào giữa slice.

```go
var a []int
a = append(a[:i], append([]int{x}, a[i:]...)...)     // chèn x ở vị trí thứ i
a = append(a[:i], append([]int{1,2,3}, a[i:]...)...) // chèn một slice con vào slice ở vị trí thứ i
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-5/main.go)

Cách `append` thứ hai sẽ gây ra việc tạo một slice tạm thời, slice `a[i:]` sẽ sao chép nội dung vào slice mới được tạo, và thêm slice tạm thời này vào `a[:i]`
Bạn cũng có thể sử dụng hàm `copy` và `append` kết hợp với nhau để tránh việc khởi tạo những slice tạm thời như vậy, cũng như có thể hoàn thành việc thêm phần tử vào một vị trí bất kỳ trong slice như sau

```go
a = append(a, 0)
copy(a[i+1:], a[i:]) // lùi những phần tử từ i trở về sau của a
a[i] = x             // gán vị trí thứ i bằng x
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-6/main.go)

Dòng đầu tiên dùng `append` để mở rộng kích thước của slice và tạo không gian cho phần tử mới được thêm vào. Ở dòng thứ hai sẽ sao chép các phần tử trong slice dời về sau kể từ vị trí thứ i. Dòng cuối cùng sẽ gán giá trị mới vào vị trí thứ i. Mặc dù cách làm trên sẽ dài dòng, tuy nhiên chúng ta có thể lượt bỏ việc phải sao chép một slice tạm thời khi so sánh với cách làm trước.

Chúng ta cũng có thể chèn nhiều phần tử vào vị trí chính giữa bằng việc kết hợp hàm `copy` và `append` như sau


```go
a = append(a, x...)       // mở rộng không gian của slice a với array x
copy(a[i+len(x):], a[i:]) // sao chép len(x) phần tử lùi về sau
copy(a[i:], x)            // sao chép array x vào giữa
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-7/main.go)

**Xóa những phần tử trong slice** 

Có ba trường hợp phụ thuộc vào nơi mà chúng ta muốn xóa phần tử, từ đầu, từ cuối hoặc từ chính giữa, trong đó xóa phần tử từ cuối là nhanh nhất.

```go
a = []int{1, 2, 3}
a = a[:len(a)-1]   // xóa một phần tử ở cuối 
a = a[:len(a)-N]   // xóa N phần tử ở cuối
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-8/main.go)

Xóa phần tử ở đầu thì thực chất là di chuyển con trỏ dữ liệu về sau

```go
a = []int{1, 2, 3}
a = a[1:] // xóa phần tử đầu tiên
a = a[N:] // xóa N phần tử đầu tiên
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-9/main.go)

Bạn cũng có thể xóa bỏ con trỏ dữ liệu mà không di chuyển phần còn lại về phía sau, nhưng sẽ di chuyển chúng tới nơi bắt đầu sẽ có thể thực hiện bởi hàm `append`, chúng không làm thay đổi cấu trúc không gian vùng nhớ

```go
a = []int{1, 2, 3}
a = append(a[:0], a[1:]...) // xóa phần tử đầu tiên
a = append(a[:0], a[N:]...) // xóa N phần tử đầu tiên
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-10/main.go)

Bạn cũng có thể dùng hàm `copy` để hoàn thành nhiệm vụ xóa

```go
a = []int{1, 2, 3}
a = a[:copy(a, a[1:])] // xóa phần tử đầu tiên
a = a[:copy(a, a[N:])] // xóa N phần tử đầu tiên
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-11/main.go)

Khi xóa phần tử ở giữa, bạn cần dịch chuyển những phần tử ở phía sau lên trước, điều đó có thể được thực hiện như sau

```go
a = []int{1, 2, 3, ...}

a = append(a[:i], a[i+1:]...) //  xóa phần tử ở vị trí i
a = append(a[:i], a[i+N:]...) //  xóa N phần tử từ vị trí i

a = a[:i+copy(a[i:], a[i+1:])]  // xóa phần tử ở vị trí i
a = a[:i+copy(a[i:], a[i+N:])]  // xáo N phần từ từ vị trí i
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-12/main.go)

Xóa phần tử đầu hoặc phần tử cuối, có thể được xem là những trường hợp đặc biệt của xóa nhũng phần tử ở giữa.

**Kĩ thuật quản lý vùng nhớ trong slice**

Ở chương array bắt đầu bằng việc đề cập sự tương tự của mảng rỗng `[0]int`, và mảng rỗng thường hiếm khi được sử dụng. Nhưng trong slice, khi chiều dài `len` là 0 và `cap` thì lớn hơn 0, thì slice có nhiều tính chất thú vị. Dĩ nhiên, nếu `len` và `cap` bằng 0, khi đó, slice sẽ trở thành slice rỗng, mặc dùng giá trị của slice đó khác `nil`. Để biết slice có trống hay không, ta sẽ dùng thuộc tính chiều dài của slice, và `nil` thì rất hiếm để so sánh trực tiếp slice và một giá trị.

Ví dụ như sau, hàm `TrimSpace` sau sẽ xóa một vùng không gian `[]byte`. Hiện thực hàm này với độ phức tạp O(n) thể đạt được sự hiệu quả và đơn giản.

```go
func TrimSpace(s []byte) []byte {
    b := s[:0]
    for _, x := range s {
        if x != ' ' {
            b = append(b, x)
        }
    }
    return b
}
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-13/main.go)

Trong thực thế những giải thuật tương tự để xóa những phần tử trong slice thỏa một điều kiện nào đó, có thể được xử lý theo cách trên, (bởi vì không có chi phí vùng nhớ phụ cho tác vụ xóa).

```go
func Filter(s []byte, fn func(x byte) bool) []byte {
    b := s[:0]
    for _, x := range s {
        if !fn(x) {
            b = append(b, x)
        }
    }
    return b
}
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-14/main.go)

Điểm chính của những tác vụ được coi là hiệu quả trên slice là hạn chế việc phải phân bố lại vùng nhớ, cố gắng để hàm `append` sẽ không đạt tới `cap` sức chứa của slice, là giảm số lần cấp phát vùng nhớ và giảm kích thước vùng nhớ cấp phát tại mọi thời điểm.

**Tránh gây ra memory leak trên slice**

Như đã đề cập ở trên, những tác vụ trên slice sẽ không sao chép vùng nhớ bên dưới slice. Bên dưới array sẽ được lưu trữ trên vùng nhớ cho đến khi nào không còn được tham chiếu nữa. Nhưng thỉnh thoảng toàn bộ dữ liệu bên dưới array sẽ có thể mang một trạng thái được sử dụng, bởi vì có một phần nhỏ vùng nhớ được tham chiếu, nó sẽ trì hoãn quá trình tự động thu gom vùng nhớ để đòi lại vùng nhớ đã cấp phát.

Ví dụ là hàm `FindPhoneNumber` sẽ tải toàn bộ file vào bộ nhớ, và sau đó tìm kiếm số điện thoại đầu tiên xuất hiện trong file, và kết quả cuối cùng sẽ trả về một array.

```go
func FindPhoneNumber(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    return regexp.MustCompile("[0-9]+").Find(b)
}
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-15/main.go)

Mã nguồn này sẽ trả về  một mảng các `byte` trỏ tới toàn bộ file. Bởi vì slice tham khảo tới toàn bộ array gốc, cơ chế tự động thu gom rác không thể giải phóng không gian bên dưới array trong thời gian đó. Một yêu cầu kết quả nhỏ, những phải lưu trữ toàn bộ dữ liệu trong một thời gian dài. Mặc dù nó không phải là `memory leak` trong ngữ cảnh truyền thống, nó có thể làm chậm hiệu suất của toàn hệ thống.

Để khắc phục vấn đề này, bạn có thể sao chép dữ liệu cần thiết thành một slice mới (kiểu giá trị của dữ liệu là một triết lý khi lập trình Go, mặc dù có được giá trị đó sẽ có một cái giá phải trả, nhưng lợi ích của việc tách biệt từ dữ liệu gốc như bên dưới)

```go
func FindPhoneNumber(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    b = regexp.MustCompile("[0-9]+").Find(b)
    return append([]byte{}, b...)
}
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-16/main.go)

Vấn đề tương tự có thể gặp phải khi xóa những phần tử trong slice. Giả sử rằng con trỏ đối tượng được lưu trữ trong cấu trúc của slice, sau khi xóa đi phần tử cuối, thì phần tử được xóa có thể còn được tham khảo bên dưới mảng slice, vùng nhớ có thể được giải phóng tự động trong thời gian đó (nó phụ thuộc vào cách hiện thực cơ chế thu hồi vùng nhớ)

```go
var a []*int{ ... }
a = a[:len(a)-1]    // phần tử cuối cùng dù được xóa nhưng vẫn được tham chiếu, do đó cơ chế thu gom rác tự động không thu hồi nó
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-17/main.go)

Phương pháp đảm bảo là đầu tiên thiết lập phần tử cần thu hồi về `nil` để đảm bảo quá trị thu gom tự động có thể tìm thấy chúng, sau đó xóa slices đó.

```go
var a []*int{ ... }
a[len(a)-1] = nil // phần tử cuối cùng sẽ được gán giá trị nil
a = a[:len(a)-1]  // xóa phần tử cuối cùng ra khỏi slice
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-18/main.go)

Dĩ nhiên, nếu ở cách làm trước đối với slice có kích thước nhỏ, bạn sẽ không gặp phải vấn đề về  tham chiếu treo. Bởi vì nếu bản thân slice có thể được giải phóng bởi GC (Garbage collector), mỗi phần tử ứng với slice có thể được thu gom tự nhiên.


**Ép kiểu slice**

Vì lý do an toàn, khi hai kiểu slice là `T[]` và `[]Y`, bên dưới phần dữ liệu thô sẽ khác nhau, ngôn ngữ Go sẽ không trực tiếp chuyển đổi kiểu. Tuy nhiên, tính an toàn đi kèm với một chi phí. Thỉnh thoảng phép chuyển đổi đi kèm theo giá trị- có thể đơn giản code hoặc cải thiện code đó. Ví dụ, trên hệ điều hành 64 bít, bạn cần phải sắp xếp một mảng `[]float64` với tốc độ cao. Chúng ta có thể ép chúng về kiểu `[]int` slice và sắp xếp chúng (bởi vì `float64` là chuẩn dấu chấm động `IEEE754` được sử dụng, số nguyên ứng với nó cũng sẽ theo thứ tự đó) điều đó không có gì xa lạ.

Đoạn mã bên dưới sẽ chuyển slice `[]float64` đến slice `[]int` bằng hai cách.

```go
// +build amd64 arm64

import "sort"

var a = []float64{4, 2, 5, 7, 2, 1, 88, 1}

func SortFloat64FastV1(a []float64) {
    var b []int = ((*[1 << 20]int)(unsafe.Pointer(&a[0])))[:len(a):cap(a)]
    sort.Ints(b)
}

func SortFloat64FastV2(a []float64) {
    var c []int
    aHdr := (*reflect.SliceHeader)(unsafe.Pointer(&a))
    cHdr := (*reflect.SliceHeader)(unsafe.Pointer(&c))
    *cHdr = *aHdr

    sort.Ints(c)
}
```

[>> mã nguồn](../examples/chapter1/ch1.3/3-slices/example-19/main.go)

Cách ép kiểu đầu tiên ban đầu sẽ chuyển địa chỉ bắt đầu của slice thành con trỏ đến mảng lớn hơn, sau đó sẽ `re-slice` array tương ứng với con trỏ array. Ở giữa `unsafe.Pointer` cần phải kết nối tới kiểu dữ liệu khác của pointer để truyền. Nên chú ý rằng, kiểu array none-zero sẽ tối đa 2GB chiều dài, do đó chúng ta có thể tính toán chiều dài tối đa của array cho kiểu array đó (kiểu `[]uint8` có kích thước tối đa 2GB, kiểu `[]uint16` tối đa 1GB, nhưng kiểu `[]struct{}` kích thước tối đa 2GB).

Cách chuyển đổi thứ hai sẽ chứa hai kiểu dữ liệu về thông tin header của slice và bên dưới cấu trúc `reflect.SliceHeader` của thông tin header sẽ ứng với cấu trúc slice, sau đó thông tin sẽ được cập nhật cấu trúc, sau đó hiện thực biến `a` ứng với kiểu `[]float64` bởi cấu trúc `[]int`. Đây là phép chuyển đổi kiểu của slice.

Thông qua việc benmarking, chúng ta có thể thấy rằng hiệu suất của việc sắp xếp `sort.Ints` của kiểu `[]int` sẽ tốt hơn là `sort.Float64s`. Tuy nhiên, bạn phải chú ý rằng, tiền đề của phương pháp đó, sẽ đảm bảo rằng `[]float64` sẽ không có dấu phẩy động chính tắc như NaN và Inf (vì NaN không thể sắp xếp theo số đấu chấm động, dương 0 và âm 0 bằng nhau, nhưng không có trường hợp nào như vậy trong số nguyên).