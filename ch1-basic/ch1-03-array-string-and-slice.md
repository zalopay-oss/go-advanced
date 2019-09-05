# 1.3. Array, strings và slices

`Arrays` và một số cấu trúc dữ liệu liên quan khác được sử dụng thường xuyên trong các ngôn ngữ lập trình. Chỉ khi chúng không đáp ứng được yêu cầu chúng ta mới cân nhắc sử dụng `linked lists` (danh sách liên kết) và `hash tables` (bảng băm) hoặc nhiều cấu trúc dữ liệu tự định nghĩa phức tạp khác.

`Arrays`, `strings` và `slices` trong ngôn ngữ Go là các cấu trúc dữ liệu liên quan mật thiết với nhau. Ba kiểu dữ liệu đó có cùng cấu trúc vùng nhớ lưu trữ bên dưới, và chỉ có những hành vi thể hiện ra bên ngoài khác nhau tùy thuộc vào ràng buộc ngữ nghĩa.

## 1.3.1. Array

<div align="center">
	<img src="../images/ch1-1-array-and-array-index-representation.png" width="600">
	<br/>
	<span align="center">
		<i>Array</i>
	</span>
</div>

Trong ngôn ngữ Go, `array` là một kiểu giá trị. Mặc dù những phần tử của array có thể được chỉnh sửa, phép gán của array hoặc khi truyền array như là một tham số của hàm thì chúng sẽ được xử lý toàn bộ, có thể hiểu là khi đó chúng được sao chép lại toàn bộ thành một bản sao rồi mới xử lý trên bản sao đó (khác với kiểu truyền tham khảo).

Một array là một chuỗi độ dài cố định của các phần tử có kiểu dữ liệu nào đó, một array có thể bao gồm không hoặc nhiều phần tử. Độ dài của array là một phần thông tin được chứa trong nó, các array có độ dài khác nhau hoặc kiểu phần tử bên trong khác nhau được xem là các kiểu dữ liệu khác nhau, và không được phép gán cho nhau, vì thế array hiếm khi được sử dụng trong Go.

Cách định nghĩa một array:

```go
// Định nghĩa một mảng kiểu int độ dài 3, các phần tử đều bằng 0
var a [3]int
// Định nghĩa một mảng có ba phần tử 1, 2, 3, do đó độ dài là 3
var b = [...]int{1, 2, 3} 
// Mảng này có 3 phần tử theo thứ tự là 0, 2, 3
var c = [...]int{2: 3, 1: 2} 
// Mảng này chứa dãy các phần tử là 1, 2, 0 , 0, 5, 6
var d = [...]int{1, 2, 4: 5, 6}
```

Cấu trúc vùng nhớ của array thì rất đơn giản. Ví dụ cho một array `[4]int{2,3,5,7}` thì cấu trúc bên dưới sẽ như sau:

<div align="center">
	<img src="../images/ch1-7-array-4int.ditaa.png">
	<br/>
	<span align="center">
		<i>Array layout</i>
	</span>
</div>

Khi biến array được gán hoặc truyền, thì toàn bộ array sẽ được sao chép. Nếu kích thước của array lớn, thì phép gán array sẽ chịu tổn phí lớn. Để tránh việc `overhead` (tổn phí) trong việc sao chép array, bạn có thể truyền con trỏ của array.

**Lưu ý:** con trỏ array thì không phải là một array.

```go
// a là một array
var a = [...]int{1, 2, 3}
// b là một con trỏ tới array a
var b = &a
// in ra hai phần tử đầu tiên của array a
fmt.Println(a[0], a[1])
// truy xuất các phần tử của con trỏ array cũng giống như truy xuất các phần tử của array
fmt.Println(b[0], b[1])
// duyệt qua các phần tử trong con trỏ array, giống như duyệt qua array
for index, value := range b {
// thay đổi từng phần tử trong b
    b[index] += 1
    fmt.Println(index, value)
}
// giá trị của các phần tử trong a bị thay đổi vì b
for index, value := range a {
    fmt.Println(index, value)
}
```

Kết quả là :

```sh
$ go run main.go
1 2
1 2
0 1
1 2
2 3
0 2
1 3
2 4
```

Hàm `len` có thể dùng để lấy thông tin về độ dài của array, và hàm `cap` sẽ tính toán độ dài tối đa của array. Nhưng trong kiểu array cả hai hàm này sẽ cùng trả về một giá trị giống nhau (điều này khác với slice).

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

`for range` là cách tốt nhất để duyệt qua các phần tử trong array, bởi vì cách này sẽ đảm các việc truy xuất sẽ không vượt quá giới hạn của array.

Các phần tử của array không nhất thiết là kiểu số học, nên cũng có thể là string, struct, function, interface, và channel, v,v..

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

// Mảng channel
var chanList = [2]chan int{}
```

Chúng ta cũng có thể định nghĩa một array rỗng.

```go
// Định nghĩa một array chiều dài 0
var d [0]int
// Tương tự trên
var e = [0]int{}
// Tương tự như trên
var f = [...]int{}
```

Một array có chiều dài 0 thì không chiếm không gian lưu trữ.

## 1.3.2. String

<div align="center">
	<img src="../images/ch1-string.png" width="400">
</div>
<br/>

`string` cũng là một array của các `byte` dữ liệu, nhưng khác với array những phần tử của string là [immutable](https://en.wikipedia.org/wiki/Immutable_object).

Cấu trúc [reflect.StringHeader](https://golang.org/src/reflect/value.go?s=56526:56578#L1873) được dùng để biểu diễn string :

```go
type StringHeader struct {
    // con trỏ địa chỉ vùng nhớ string
    Data uintptr
    // chiều dài của string
    Len  int
}
```

Một string là một cấu trúc, do đó phép gán string thực chất là việc sao chép cấu trúc `reflect.StringHeader`, và không gây ra việc sao chép bên dưới phần dữ liệu.

Cấu trúc vùng nhớ tương ứng với string "Hello World" là:

<div align="center">
	<img src="../images/ch1-8-string-1.ditaa.png"width="600">
	<br/>
	<span align="center">
		<i>String layout</i>
	</span>
</div>

Có thể thấy rằng bên dưới string "Hello World" là một array như sau:

```go
var  data = [...] byte {
    'H', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd',
}
```

Mặc dù string không phải là slice nhưng nó cũng hỗ trợ thao tác (slicing) cắt. Một vài phần của vùng nhớ cũng được truy cập bên dưới slice tại một số nơi khác nhau:

```go
// s là biến string
var s = "hello world"
// lấy phần tử từ index 0 tới 4
hello := s[:5]
// lấy phần tử từ index 6 tới hết
world := s[6:]
// có thể thao tác trực tiếp
s1 := "hello world"[:5]
s2 := "hello world"[6:]
```

Tương tự như array, String cũng có một hàm built-in là `len` dùng để trả về chiều dài của string, ngoài ra bạn có thể  dùng `reflect.StringHeader` để truy xuất chiều dài của string theo cách như sau

```go
fmt.Println("len(s): ", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len)
```

## 1.3.3. Slice

<div align="center">
	<img src="../images/1-3-golang-slices-length-capacity.jpg"width="500">
	<br/>
	<span align="center">
		<i>Cấu trúc Slice</i>
	</span>
</div>

`Slices` thì phức tạp hơn, cấu trúc của chúng cũng như `string`, tuy nhiên việc giới hạn chỉ-đọc như string được lược bỏ.

Cấu trúc của slice là `reflect.SliceHeader`

```go
type  SliceHeader  struct {
	Data  uintptr 
	Len   int 
	Cap   int 
}
```
Ngoài `Data` và `Len`, slice có thêm thuộc tính `Cap` chỉ ra kích thước tối đa mà vùng nhớ trỏ tới slice được cấp phát. 

Hình bên dưới sẽ miêu tả slice `x := []int{2,3,5,7,11}` và slice `y := x[1:3]`:

<div align="center">
	<img src="../images/ch1-10-slice-1.ditaa.png">
	<br/>
	<span align="center">
		<i>Slice layout</i>
	</span>
</div>

Các cách định nghĩa slice:

```go
var (
	// nil slice
	a = []int
	// empty slice, khác với nil
	b = []int{}
	// có 3 phần tử trong slice, cả len và cap đều bằng 3
	c = []int{1,2,3}
	// có 2 phần tử trong slice, len bằng 2 và cap bằng 3
	d = c[:2]
	// có 2 phần tử trong slice, len bằng 2 và cap bằng 3
	e = c[0:2:cap(c)]
	// có 0 phần tử trong slice, len bằng 0 và cap bằng 3
	f = c[:0]
	// có 3 phần tử trong slice, len và cap bằng 3
	g = make ([]int,3)
	// có 2 phần tử trong slice, len bằng 2, cap bằng 3
	h = make ([]int,2,3)
	// có 0 phần tử trong slice, len bằng 0, cap bằng 3
	i = make ([]int,0,3)
)
```

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

Các tác vụ cơ bản trong slice bao gồm:
  * Thêm phần tử vào slice
  * Xóa phần tử trong slcie
  * Duyệt qua các phần tử của slice

#### Thêm phần tử vào slice

Hàm  `append` có thể thêm phần tử thứ `N` vào cuối cùng của slice:

```go
var a []int
// nối thêm phần tử 1
a = append(a, 1)
// nối thêm phần tử 1, 2, 3
a = append(a, 1, 2, 3)
// nối thêm các phần tử 1, 2, 3 bằng cách truyền vào một mảng
a = append(a, []int{1,2,3}...)
```

Trong trường hợp slice ban đầu không đủ sức chứa khi thêm vào phần tử, hàm `append` sẽ hiện thực cấp phát lại vùng nhớ, chi phí của việc cấp phát và sao chép là tương đối đáng kể.

Bên cạnh thêm phần tử vào cuối slice, chúng ta cũng có thể thêm phần tử vào đầu slice như sau

```go
var a = []int{1,2,3}
// thêm phần tử 0 vào đầu slice a
a = append([]int{0}, a...)
// thêm các phần tử -3, -2, -1 vào đầu slice a
a = append([]int{-3,-2,-1}, a...)
```

Thêm phần tử vào đầu slice sẽ gây ra việc cấp phát lại vùng nhớ và làm những phần tử đang tồn tại trong slice sẽ được sao chép một lần nữa. Do đó, hiệu suất của việc thêm phần tử  vào đầu slice sẽ thấp hơn thêm phần tử vào cuối slice.

Do hàm `append` sẽ trả về một slice mới, chúng ta có thể kết hợp nhiều hàm `append` để chèn một vài phần tử vào giữa slice.

```go
// khai báo slice a
var a []int
// chèn x ở vị trí thứ i
a = append(
    a[:i],
    // tạo ra một slice tạm thời để nối với a[:i]
    append(
        []int{x},a[i:]...
    )...
)
// chèn một slice con vào slice ở vị trí thứ i
a = append(
    a[:i],
    append(
        []int{1,2,3},
        a[i:]...
    )...
)
```

Bạn cũng có thể sử dụng hàm `copy` và `append` kết hợp với nhau để tránh việc khởi tạo những slice tạm thời như vậy:

```go
// thêm phần tử 0 vào cuối slice a
a = append(a, 0)
// lùi những phần tử từ i trở về sau
copy(a[i+1:], a[i:])
// gán vị trí thứ i bằng x
a[i] = x
```


Chúng ta cũng có thể chèn nhiều phần tử vào vị trí chính giữa bằng việc kết hợp hàm `copy` và `append` như sau

```go
// mở rộng không gian của slice a với array x
a = append(a, x...)
// sao chép len(x) phần tử lùi về sau
copy(a[i+len(x):], a[i:])
// sao chép array x vào giữa
copy(a[i:], x)
```

#### Xóa những phần tử trong slice

Có ba trường hợp xóa các phần tử:
  * Ở đầu
  * Ở giữa
  * Ở cuối


Trong đó xóa phần tử ở cuối là nhanh nhất

```go
a = []int{1, 2, 3}
// xóa một phần tử ở cuối
a = a[:len(a)-1]
// xóa N phần tử ở cuối
a = a[:len(a)-N]
```

Xóa phần tử ở đầu thì thực chất là di chuyển con trỏ dữ liệu về sau

```go
a = []int{1, 2, 3}
// xóa phần tử đầu tiên
a = a[1:]
// xóa N phần tử đầu tiên
a = a[N:]
```

Khi xóa phần tử ở giữa, bạn cần dịch chuyển những phần tử ở phía sau lên trước, điều đó có thể được thực hiện như sau

```go
a = []int{1, 2, 3, ...}
// xóa phần tử ở vị trí i
a = append(a[:i], a[i+1:]...)
// xóa N phần tử từ vị trí i
a = append(a[:i], a[i+N:]...)
// xóa phần tử ở vị trí i
a = a[:i+copy(a[i:], a[i+1:])]
// xoá N phần từ từ vị trí i
a = a[:i+copy(a[i:], a[i+N:])]
```

#### Kỹ thuật quản lý vùng nhớ trong slice

Hàm `TrimSpace` sau sẽ xóa đi các khoảng trắng. Hiện thực hàm này với độ phức tạp O(n) để đạt được sự hiệu quả và đơn giản.

```go
func TrimSpace(s []byte) []byte {
    b := s[:0]
    // duyệt qua slice s để tìm phần tử thỏa điều kiện
    for _, x := range s {
        // kiểm tra điều kiện
        if x != ' ' {
        // tạo ra slice mới từ slice ban đầu thêm vào phần tử x
            b = append(b, x)
        }
    }
    return b
}
```

Thực tế, những giải thuật tương tự để xóa những phần tử trong slice thỏa một điều kiện nào đó, có thể được xử lý theo cách trên.

**Lưu ý:**  hàm `append()` không cấp phát lại vùng nhớ khi chưa đạt tới sức chứa tối đa `cap`.

Hàm `Filter` sẽ lọc các phần tử thỏa điều kiện trong slice.

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

Điểm chính của những tác vụ làm việc hiệu quả trên slice là hạn chế việc phải cấp lại vùng nhớ, cố gắng để hàm `append` sẽ không đạt tới `cap` sức chứa của slice, là giảm số lần cấp phát vùng nhớ và giảm kích thước vùng nhớ cấp phát tại mọi thời điểm.

####  Tránh gây ra memory leak trên slice

Những tác vụ trên slice sẽ không thay đổi vùng nhớ bên dưới slice, mà thực chất là thay đổi các tham số như `Data`, `Len`. Vùng nhớ bên dưới vẫn còn cho đến khi nào không còn được tham chiếu nữa.

Thỉnh thoảng toàn bộ dữ liệu bên dưới slice sẽ có thể mang một trạng thái đang được sử dụng vì chỉ có một phần nhỏ vùng nhớ được tham chiếu tới, nó sẽ trì hoãn quá trình tự động thu gom vùng nhớ để đòi lại vùng nhớ đã cấp phát.

Ví dụ là hàm `FindPhoneNumber` sẽ tải toàn bộ file vào bộ nhớ, và sau đó tìm kiếm số điện thoại đầu tiên xuất hiện trong file, và kết quả cuối cùng sẽ trả về một array.

```go
func FindPhoneNumber(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    return regexp.MustCompile("[0-9]+").Find(b)
}
```

Mã nguồn này sẽ trả về  một mảng các `byte` trỏ tới toàn bộ file. Bởi vì slice tham khảo tới toàn bộ array gốc, cơ chế tự động thu gom rác không thể giải phóng không gian bên dưới array trong thời gian đó. Một yêu cầu kết quả nhỏ, những phải lưu trữ toàn bộ dữ liệu trong một thời gian dài. Mặc dù nó không phải là `memory leak` trong ngữ cảnh truyền thống, nó có thể làm chậm hiệu suất của toàn hệ thống.

Để khắc phục vấn đề này, bạn có thể sao chép dữ liệu cần thiết thành một slice mới.

```go
func FindPhoneNumber(filename string) []byte {
    b, _ := ioutil.ReadFile(filename)
    b = regexp.MustCompile("[0-9]+").Find(b)
    return append([]byte{}, b...)
}
```

Vấn đề tương tự có thể gặp phải khi xóa những phần tử trong slice. Giả sử rằng con trỏ đối tượng được lưu trữ trong cấu trúc của slice, sau khi xóa đi phần tử cuối, thì phần tử được xóa có thể còn được tham khảo bên dưới mảng slice, vùng nhớ có thể được giải phóng tự động trong thời gian đó (nó phụ thuộc vào cách hiện thực cơ chế thu hồi vùng nhớ)

```go
var a []*int{ ... }
// phần tử cuối cùng dù được xóa nhưng vẫn được tham chiếu,
// do đó cơ chế thu gom rác tự động không thu hồi nó
a = a[:len(a)-1]
```

Cách giải quyết là đầu tiên thiết lập phần tử cần thu hồi về `nil` để đảm bảo giá trị thu gom tự động có thể tìm thấy chúng, sau đó xóa slices đó.

```go
var a []*int{ ... }
// phần tử cuối cùng sẽ được gán giá trị nil
a[len(a)-1] = nil
// xóa phần tử cuối cùng ra khỏi slice
a = a[:len(a)-1]
```

[Tiếp theo](ch1-04-func-method-interface.md)