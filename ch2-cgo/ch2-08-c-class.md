# 2.8. C++ Class Packaging

CGO không hỗ trợ trực tiếp tính năng về class trong C++. Nguyên nhân đến từ việc CGO không có hỗ trợ cho cú pháp C++, ngoài ra C++ cũng không có sẵn một Application Binary Interface ([ABI](https://stackoverflow.com/questions/2171177/what-is-an-application-binary-interface-abi)) nào cả.

Nhưng do C++ tương thích với C, do đó chúng ta có thể thêm một tập các hàm C cho interface như là cầu nối giữa C++ class và CGO. Dĩ nhiên, bởi vì CGO chỉ hỗ trợ kiểu dữ liệu của ngôn ngữ C, chúng ta không thể trực tiếp dùng C++ reference parameters và các tính năng khác.

## 2.8.1. Từ Class của C++ đến Object của Go

Việc hiện thực packaging C++ class thành Object trong Go yêu cầu một số bước:

- Đầu tiên, C++ class được bọc bởi một interface C thuần,
- Tiếp theo hàm trong interface sẽ map với hàm của Go bằng CGO
- Cuối cùng là tạo ra đối tượng Go wrapper. Lúc này ta có thể hiện thực class C++ thành các phương thức sử dụng Go objects.

### Bước 1: Chuẩn bị một C++ class

Để minh họa, chúng ta sẽ dựa trên `str::string` để làm một class đơn giản `MyBuffer`:

```c++
// my_buffer.h
#include <string>

struct MyBuffer {
    std::string* s_;

    // thêm vào constructor
    MyBuffer(int size) {
        this->s_ = new std::string(size, char('\0'));
    }

    // và destructor
    ~MyBuffer() {
        delete this->s_;
    }

    // trả về kích thước buffer
    int Size() const {
        return this->s_->size();
    }

    // trả về con trỏ tới data
    char* Data() {
        return (char*)this->s_->data();
    }
};
```

Tiếp theo là cách chúng ta sử dụng:

```c++
int main() {
    auto pBuf = new MyBuffer(1024);

    auto data = pBuf->Data();
    auto size = pBuf->Size();

    delete pBuf;
}
```

### Bước 2: Đóng gói Class C++ với Interface C

Chúng ta có thể ánh xạ từ khóa `new` và `delete` C++ sang C và tương tự các phương thức của đối tượng tới hàm của ngôn ngữ C.

Trong C chúng ta mong đợi class `MyBuffer` được dùng như sau:

```c++
int main() {
    MyBuffer* pBuf = NewMyBuffer(1024);

    char* data = MyBuffer_Data(pBuf);
    auto size = MyBuffer_Size(pBuf);

    DeleteMyBuffer(pBuf);
}
```

Đặc tả file header `my_buffer_capi.h`:

```c++
// my_buffer_capi.h
typedef struct MyBuffer_T MyBuffer_T;

MyBuffer_T* NewMyBuffer(int size);
void DeleteMyBuffer(MyBuffer_T* p);

char* MyBuffer_Data(MyBuffer_T* p);
int MyBuffer_Size(MyBuffer_T* p);
```

Sau đó chúng ta có thể định nghĩa các hàm wrapper dựa trên class C++ là `MyBuffer`. File `my_buffer_capi.cc` tương ứng như sau:

```c++
// my_buffer_capi.cc

#include "./my_buffer.h"

// chỉ ra các file C++ được include
// mục đích để hàm C gọi được C++
extern "C" {
    #include "./my_buffer_capi.h"
}

// một class kế thừa từ `MyBuffer`, thực ra
// là hiện thực việc wrapper code
struct MyBuffer_T: MyBuffer {
    MyBuffer_T(int size): MyBuffer(size) {}
    ~MyBuffer_T() {}
};

MyBuffer_T* NewMyBuffer(int size) {
    auto p = new MyBuffer_T(size);
    return p;
}
void DeleteMyBuffer(MyBuffer_T* p) {
    delete p;
}

char* MyBuffer_Data(MyBuffer_T* p) {
    return p->Data();
}
int MyBuffer_Size(MyBuffer_T* p) {
    return p->Size();
}
```

Lúc này trở đi, `MyBuffer_T` giao tiếp với CGO có thể thông qua việc truyền pointer.

Sau khi wrapping C++ class thành một C interface thuần, bước tiếp theo là chuyển đổi hàm C sang hàm Go.

### Bước 3: Chuyển đổi hàm trong C interface sang Go

Quá trình bọc hàm C thành một hàm Go là tương đối đơn giản. Chú ý rằng bởi vì package của chúng ta chứa cú pháp C++11, chúng ta cần flag `#cgo CXXFLAGS: -std=c++11` để mở tùy chọn `C++11`.

```go
// my_buffer_capi.go

package main

/*
#cgo CXXFLAGS: -std=c++11

#include "my_buffer_capi.h"
*/
import "C"

type cgo_MyBuffer_T C.MyBuffer_T

// Để phân biệt, chúng ta thêm vào một tiền tố `cgo_`
// cho mỗi hàm được đặt tên trong Go.
// `cgo_MyBuffer_T` là một kiểu `MyBuffer_T` trong C
func cgo_NewMyBuffer(size int) *cgo_MyBuffer_T {
    p := C.NewMyBuffer(C.int(size))
    return (*cgo_MyBuffer_T)(p)
}

func cgo_DeleteMyBuffer(p *cgo_MyBuffer_T) {
    C.DeleteMyBuffer((*C.MyBuffer_T)(p))
}

func cgo_MyBuffer_Data(p *cgo_MyBuffer_T) *C.char {
    return C.MyBuffer_Data((*C.MyBuffer_T)(p))
}

func cgo_MyBuffer_Size(p *cgo_MyBuffer_T) C.int {
    return C.MyBuffer_Size((*C.MyBuffer_T)(p))
}
```

Khi đóng gói một hàm C thành một hàm Go, bằng việc thêm vào kiểu `cgo_MyBuffer_T`, chúng ta vẫn dùng kiểu của C bên dưới kiểu dữ liệu cho tham số đầu vào và cho giá trị trả về.

### Bước 4: Tạo đối tượng Go wrapper

Sau khi bọc interface C thuần thành một hàm Go, chúng ta sẽ xây dựng một đối tượng Go wrapper.

Bởi vì `cgo_MyBuffer_T` là một kiểu được import trong không gian ngôn ngữ C, không thể định nghĩa những phương thức cho riêng chúng, do đó chúng ta phải xây dựng một kiểu `MyBuffer` sẽ .

```go
// my_buffer.go

package main

import "unsafe"

// giữ buffer của ngôn ngữ C được trỏ tới bởi `cgo_MyBuffer_T`
type MyBuffer struct {
    cptr *cgo_MyBuffer_T
}

func NewMyBuffer(size int) *MyBuffer {
    return &MyBuffer{
        cptr: cgo_NewMyBuffer(size),
    }
}

func (p *MyBuffer) Delete() {
    cgo_DeleteMyBuffer(p.cptr)
}

// vì kiểu slice của Go có chứa cả thông tin chiều dài,
// chúng ta có thể kết hợp hai hàm `cgo_MyBuffer_Data`
// và `cgo_MyBuffer_Size` vào trong phương thức `MyBuffer.Data`
func (p *MyBuffer) Data() []byte {
    data := cgo_MyBuffer_Data(p.cptr)
    size := cgo_MyBuffer_Size(p.cptr)


    // trả về một slice ứng với buffer trong C
    return ((*[1 << 31]byte)(unsafe.Pointer(data)))[0:int(size):int(size)]
}
```

Bây giờ chúng ta có thể dễ dàng sử dụng wrapped buffer object trong ngôn ngữ Go (ngầm bên trong là phần hiện thực `std::string` C++)

```go
package main

//#include <stdio.h>
import "C"
import "unsafe"

func main() {
    // tạo ra 1024-byte buffer
    buf := NewMyBuffer(1024)
    defer buf.Delete()

    // cấp phát string bằng copy
    copy(buf.Data(), []byte("hello\x00"))

    //trực tiếp lấy ra thông tin con trỏ của buffer
    // và in nội dung của buffer bằng hàm `put` của C
    C.puts((*C.char)(unsafe.Pointer(&(buf.Data()[0]))))
}
```

## 2.8.2. Chuyển đổi Object Go sang Class C++

Để hiện thực việc đóng gói các đối tượng ngôn ngữ Go vào các class C++, cần có các bước như sau:

- Trước tiên ánh xạ đối tượng Go sang một id
- Sau đó export hàm interface C tương ứng dựa trên id.
- Cuối cùng đóng gói đối tượng C++ dựa trên hàm interface C.

### Bước 1: Xây dựng một đối tượng Go

Để cho dễ theo dõi, tôi đã xây dựng một đối tượng `Person` trong Go, mỗi đối tượng có thông tin về tên và tuổi:

```go
package main

type Person struct {
    name string
    age  int
}

func NewPerson(name string, age int) *Person {
    return &Person{
        name: name,
        age:  age,
    }
}

func (p *Person) Set(name string, age int) {
    p.name = name
    p.age = age
}

func (p *Person) Get() (name string, age int) {
    return p.name, p.age
}
```

Nếu đối tượng Person muốn được truy cập trong C/C++, thì nó cần được truy cập thông qua interface C.

### Bước 2: Export object Go sang interface C

Tạo một file tương ứng với  file đặc tả interface C:

```c++
// person_capi.h
#include <stdint.h>

typedef uintptr_t person_handle_t;

person_handle_t person_new(char* name, int age);
void person_delete(person_handle_t p);

void person_set(person_handle_t p, char* name, int age);
char* person_get_name(person_handle_t p, char* buf, int size);
int person_get_age(person_handle_t p);
```

Sau đó, các hàm C này được hiện thực bằng ngôn ngữ Go.

Cần lưu ý rằng khi export ra các hàm C thông qua CGO, cả kiểu của tham số đầu vào và kiểu của giá trị trả về đều không hỗ trợ sửa đổi hằng số *const* và cũng không hỗ trợ các hàm có tham số biến. Đồng thời như đã mô tả trong phần trước ([chương 2.7](./ch2-07-cgo-mem.md)), chúng ta không thể truy cập trực tiếp các đối tượng bộ nhớ Go trong C/C++ trong một thời gian dài. Vì vậy, chúng ta cần ánh xạ đối tượng Go thành một id số nguyên.

Sau đây là file `person_capi.go` hiện thực các hàm trong  interface C:

```go
// person_capi.go
package main

//#include "./person_capi.h"
import "C"
import "unsafe"

//export person_new
func person_new(name *C.char, age C.int) C.person_handle_t {
    // ánh xạ tới id thông qua `NewObjectId`
    id := NewObjectId(NewPerson(C.GoString(name), int(age)))

    //  buộc id phải được trả về dưới dạng `person_handle_t`
    return C.person_handle_t(id)
}

//export person_delete
func person_delete(h C.person_handle_t) {
    ObjectId(h).Free()
}

//export person_set
func person_set(h C.person_handle_t, name *C.char, age C.int) {
    p := ObjectId(h).Get().(*Person)
    p.Set(C.GoString(name), int(age))
}

//export person_get_name
func person_get_name(h C.person_handle_t, buf *C.char, size C.int) *C.char {
    p := ObjectId(h).Get().(*Person)
    name, _ := p.Get()

    n := int(size) - 1
    bufSlice := ((*[1 << 31]byte)(unsafe.Pointer(buf)))[0:n:n]
    n = copy(bufSlice, []byte(name))
    bufSlice[n] = 0

    return buf
}

//export person_get_age
func person_get_age(h C.person_handle_t) C.int {
    p := ObjectId(h).Get().(*Person)
    _, age := p.Get()
    return C.int(age)
}
```

Các hàm interface khác dựa trên id được thể hiện bởi `person_handle_t`, nhờ đó đối tượng Go tương ứng được parse theo id.

### Bước 3: Đóng gói các đối tượng C++

Một cách thực hiện phổ biến là:

```c++
extern "C" {
    #include "./person_capi.h"
}

// tạo một class Person mới
struct Person {
    // chứa một thành viên thuộc kiểu `person_handle_t`
    // tương ứng với đối tượng Go
    person_handle_t goobj_;

    // tạo một đối tượng Go thông qua interface
    // C trong hàm constructor của class Person
    Person(const char* name, int age) {
        this->goobj_ = person_new((char*)name, age);
    }

    // giải phóng đối tượng Go qua interface C trong hàm destructor
    ~Person() {
        person_delete(this->goobj_);
    }

    void Set(char* name, int age) {
        person_set(this->goobj_, name, age);
    }
    char* GetName(char* buf, int size) {
        return person_get_name(this->goobj_ buf, size);
    }
    int GetAge() {
        return person_get_age(this->goobj_);
    }
}
```

Sau khi đóng gói, chúng ta có thể sử dụng nó như một class C++ bình thường:

```c++
#include "person.h"

#include <stdio.h>

int main() {
    auto p = new Person("gopher", 10);

    char buf[64];
    char* name = p->GetName(buf, sizeof(buf)-1);
    int age = p->GetAge();

    printf("%s, %d years old.\n", name, age);
    delete p;

    return 0;
}
```

### Bước 4: Cải tiến đóng gói đối tượng C++

Trong lần hiện thực đóng gói các đối tượng C++ trước đây, mỗi lần tạo một instance Person mới, ta cần thực hiện hai lần cấp phát bộ nhớ: một lần cho phiên bản Person của C++ và một lần nữa cho phiên bản Person của ngôn ngữ Go.

Trong thực tế, phiên bản C++ của Person chỉ có một id thuộc kiểu `person_handle_t`, được sử dụng để ánh xạ các đối tượng Go. Chúng ta có thể sử dụng `person_handle_t` trực tiếp trong đối tượng C++.

Các phương pháp đóng gói được cải tiến như sau đây:

```c++
extern "C" {
    #include "./person_capi.h"
}

struct Person {
    // thêm một hàm thành viên static mới vào class Person
    // để tạo một instance Person mới
    static Person* New(const char* name, int age) {
        // instance Person được tạo bằng cách gọi
        // person_new, trả về kiểu `person_handle_tid`
        // và chúng ta sử dụng nó làm con trỏ kiểu `Person*`
        return (Person*)person_new((char*)name, age);
    }

    // trong các hàm thành viên khác, ta chuyển đổi con trỏ this
    // thành một kiểu `person_handle_t` và sau đó gọi hàm
    // tương ứng thông qua interface C.
    void Delete() {
        person_delete(person_handle_t(this));
    }

    void Set(char* name, int age) {
        person_set(person_handle_t(this), name, age);
    }
    char* GetName(char* buf, int size) {
        return person_get_name(person_handle_t(this), buf, size);
    }
    int GetAge() {
        return person_get_age(person_handle_t(this));
    }
};
```

Ở thời điểm này, ta đã đạt được mục tiêu export đối tượng Go dưới dạng interface C và sau đó đóng gói lại thành đối tượng C++ dựa trên interface C.

## 2.8.3. Con trỏ `this` của C++

Các phương thức trong ngôn ngữ Go đều bị ràng buộc kiểu. Ví dụ nếu chúng ta xác định kiểu `Int` mới dựa trên int, chúng ta có thể có phương thức riêng:

```go
type Int int

func (p Int) Twice() int {
    return int(p)*2
}

func main() {
    var x = Int(42)
    fmt.Println(int(x))
    fmt.Println(x.Twice())
}
```

`this` cho phép bạn tự do chuyển đổi các kiểu int và `Int` để sử dụng các biến mà không thay đổi cấu trúc bộ nhớ cơ bản của dữ liệu gốc.

Để đạt được các tính năng tương tự trong C++, các cách hiện thực sau thường được sử dụng:

```c++
class Int {
    int v_;

    Int(v int) { this.v_ = v; }
    int Twice() const{ return this.v_*2; }
};

int main() {
    Int v(42);

    printf("%d\n", v); // error
    printf("%d\n", v.Twice());
}
```

Class `Int` mới được thêm vào thêm phương thức `Twice` nhưng mất quyền  chuyển về kiểu int. Tại thời điểm này, không chỉ `printf` không thể tự export giá trị của `Int` mà còn mất tất cả các tính năng của operation kiểu int. `this` là chủ ý của ngôn ngữ C++: đổi lợi ích từ việc sử dụng class lấy cái giá là mất tất cả các tính năng ban đầu của nó.

Nguyên nhân gốc rễ của vấn đề này là do kiểu con trỏ được cố định vào class trong C++. Hãy xem xét lại bản chất của `this` trong ngôn ngữ Go:

```go
func (this Int) Twice() int
func Int_Twice(this Int) int
```

Trong Go, kiểu của tham số receiver có chức năng tương tự như `this` chỉ là một tham số hàm bình thường. Chúng ta có thể tự do chọn giá trị hoặc kiểu con trỏ.

Nếu nghĩ theo thuật ngữ C thì `this` chỉ là một con trỏ `void*` tới một kiểu bình thường và chúng ta có thể tự do chuyển đổi nó thành các kiểu khác.

```c
struct Int {
    int Twice() {
        const int* p = (int*)(this);
        return (*p) * 2;
    }
};
int main() {
    int x = 42;
    printf("%d\n", x);
    printf("%d\n", ((Int*)(&x))->Twice());
    return 0;
}
```

Bằng cách này, chúng ta có thể xây dựng một đối tượng `Int` bằng cách buộc con trỏ kiểu int thành con trỏ kiểu `Int` thay vì hàm tạo mặc định (default constructor). Bên trong hàm `Twice`, bằng cách chuyển con trỏ `this` trở lại con trỏ int trong thao tác ngược lại, giá trị kiểu int ban đầu đã có thể được parse. Tại thời điểm này, kiểu `Int` chỉ là một lớp vỏ trong thời gian biên dịch và không chiếm thêm bộ nhớ khi chạy.

Do đó, phương thức C++ cũng có thể được sử dụng cho các kiểu không phải class. C++ cho các hàm thành phần thông thường cũng có thể được liên kết với các kiểu. Chỉ có các phương thức ảo thuần túy được ràng buộc với đối tượng và đó là interface.

[Tiếp theo](ch2-09-static-shared-lib.md)