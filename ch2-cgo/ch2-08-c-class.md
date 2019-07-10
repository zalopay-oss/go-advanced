# 2.8 C++ Class Packaging

CGO là một cầu nối giữa C và Go. Về nguyên tắc, class trong C++ không được hỗ trợ trực tiếp. Nguyên nhân đến từ việc CGO thiếu hỗ trợ cho cú pháp C++. C++ không có sẵn một Binary Interface Specification (ABI) nào cả. Làm thế nào class C++ contructor có thể sinh ra link symbol names khi biên dịch ra các object files, các phương thức khác nhau từ các platform khác nhau của các phiên bản C++ khác nhau? Nhưng C++ cũng tương thích với ngôn ngữ C, do đó chúng ta có thể thêm một tập các hàm C interface như là cầu nối giữa C++ class và CGO, do đó, việc giao tiếp giữa C++ và Go có thể được nhận dạng một cách trực tiếp. Dĩ nhiên, bởi vì CGO chỉ hỗ trợ kiểu dữ liệu của ngôn ngữ C, chúng ta không thể trực tiếp dùng C++ reference parameters và các tính năng khác.

## 2.8.1 Từ Class trong C++ đến Object trong Go

Việc hiện thực packaging C++ class thành Object trong Go yêu cầu một số bước. Đầu tiên, C++ class được bọc bởi một interface C thuần, tiếp theo, hàm interface C thuần sẽ map với hàm của Go bằng CGO, cuối cùng, đối tượng Go wrapper được tạo ra. Hiện thực class C++ thành các phương thức sử dụng Go objects.

### 2.8.1.1 Chuẩn bị một C++ class

Để minh họa đơn giản, chúng ta sẽ dựa trên `str::string` để làm một class đơn giản `MyBuffer`. Thêm vào các hàm constructor và destructor, chỉ hai phương thức được trả về kiểu con trỏ và kích thước cache. Bởi vì đó là binarry cache, chúng ta có thể thay thế các thông tin tùy ý cho nó.

```c++
// my_buffer.h
#include <string>

struct MyBuffer {
    std::string* s_;

    MyBuffer(int size) {
        this->s_ = new std::string(size, char('\0'));
    }
    ~MyBuffer() {
        delete this->s_;
    }

    int Size() const {
        return this->s_->size();
    }
    char* Data() {
        return (char*)this->s_->data();
    }
};
```

Chúng ta đặc tả kích thước của cache và cấp phát không gian cho constructor, và release vùng nhớ cục bộ thông qua destructor sau khi dùng. Đây là cách chúng ta dùng nó:

```c++
int main() {
    auto pBuf = new MyBuffer(1024);

    auto data = pBuf->Data();
    auto size = pBuf->Size();

    delete pBuf;
}
```

Để thuận tiện cho việc chuyển giao giữa interface ngôn ngữ C, từ đây ta sẽ không định nghĩa bản sao C++ constructor. ta phải cấp phát và giải phóng cache objects với lệnh `new` và `delete`, không theo kiểu giá trị.

### 2.8.1.2 Đóng gói class trong C++ với Interface C

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

Đặc tả file header `my_buffer_capi.h` như sau:

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

extern "C" {
    #include "./my_buffer_capi.h"
}

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

Bởi vì header file `my_buffer_capi.h` dành cho CGO, nó phải có một quy luật đặt tên được dùng trong đặc tả của ngôn ngữ C. Mệnh đề `extern "C"` được yêu cầu khi các file mã nguồn C++ được included. Thêm vào đó, việc hiện thực `MyBuffer_T` chỉ là một class kế thừa từ `MyBuffer`, nó đơn giản là hiện thực việc wrapper code. Giờ đây, khi giao tiếp với CGO `MyBuffer_T`, chúng ta phải truyền qua pointer. Chúng ta không thể thể hiện việc đặc tả cho việc hiện thực CGO cụ thể bởi vì việc hiện thực sẽ chứa cú pháp `C++`, và `CGO` không thể nhận diện C++ feature.

Sau khi wrapping C++ class như là một C interface thuần, bước tiếp theo là chuyển đổi hàm C đến hàm Go.

### 2.8.1.3 Chuyển đổi pure C interface function sang hàm Go

Quá trình bọc hàm C thuần thành một hàm Go là tương đối đơn giản. Chú ý rằng bởi vì package của chúng ta chứa cú pháp C++11, chúng ta cần flag `#cgo CXXFLAGS: -std=c++11` để mở tùy chọn `C++11`.

```go
// my_buffer_capi.go

package main

/*
#cgo CXXFLAGS: -std=c++11

#include "my_buffer_capi.h"
*/
import "C"

type cgo_MyBuffer_T C.MyBuffer_T

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

Để phân biệt, chúng ta thêm vào một tiền tố `cgo_` cho mỗi hàm được đặt tên trong Go. Ví dụ, `cgo_MyBuffer_T` là một kiểu `MyBuffer_T` trong C.

Để đơn giản, khi đóng gói một hàm C thuần thành một hàm Go, bằng việc thêm vào kiểu `cgo_MyBuffer_T`, chúng ta vẫn dùng kiểu ngôn ngữ C bên dưới kiểu dữ liệu cho tham số đầu vào và cho kết quả trả về.

### 2.8.1.4 Wrapper là một đối tượng của Go

Sau khi bọc interface C thuần thành một hàm Go, chúng ta có thể dễ dàng xây dựng một đối tượng Go dựa trên hàm wrapped Go. Bởi vì `cgo_MyBuffer_T` là một kiểu được imported trong không gian ngôn ngữ C, không thể định nghĩa những phương thức cho riêng chúng, do đó chúng ta phải xây dựng một kiểu `MyBuffer` sẽ giữ đối tượng cache của ngôn ngữ C được trỏ tới bởi `cgo_MyBuffer_T`.

```go
// my_buffer.go

package main

import "unsafe"

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

func (p *MyBuffer) Data() []byte {
    data := cgo_MyBuffer_Data(p.cptr)
    size := cgo_MyBuffer_Size(p.cptr)
    return ((*[1 << 31]byte)(unsafe.Pointer(data)))[0:int(size):int(size)]
}
```

Bởi vì bản thân ngôn ngữ Go có kiểu slice chứa thông tin về chiều dài, chúng ta có thể kết hợp hai hàm `cgo_MyBuffer_Data` và `cgo_MyBuffer_Size` vào trong phương thức `MyBuffer.Data`, chúng sẽ trả về một slice ứng với cache space trong ngôn ngữ C.

Bây giờ chúng ta có thể dễ dàng sử dụng wrapped cache object trong ngôn ngữ Go (bên dưới là phần hiện thực `std::string` C++)

```go
package main

//#include <stdio.h>
import "C"
import "unsafe"

func main() {
    buf := NewMyBuffer(1024)
    defer buf.Delete()

    copy(buf.Data(), []byte("hello\x00"))
    C.puts((*C.char)(unsafe.Pointer(&(buf.Data()[0]))))
}
```

Trong ví dụ chúng ta tạo ra 1024-byte cache và sau đó phân bổ string bằng hàm copy. Để thuận tiện cho việc hiện thực các hàm C, chúng ta sẽ mặc định đặt cuối mỗi string kí tự `\0`. Cuối cùng, chúng ta sẽ trực tiếp lấy ra thông tin con trỏ của cache và in nội dung của bộ đệm bằng hàm `put` của C.

## 2.8.2 Chuyển đổi đối tượng Go sang class C++

Để hiện thực việc đóng gói các đối tượng ngôn ngữ Go vào các class C++, cần có các bước như sau. Trước tiên, ánh xạ đối tượng Go sang một id. Sau đó export hàm interface C tương ứng dựa trên id. Cuối cùng đóng gói đối tượng C++ dựa trên hàm interface C.

### 2.8.2.1 Xây dựng một đối tượng Go

Để cho dễ theo dõi, chúng tôi đã xây dựng một đối tượng `Person` trong Go, mỗi đối tượng có thông tin về tên và tuổi:

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

Nếu đối tượng Person muốn được truy cập trong C/C++, thì nó cần được truy cập thông qua interface CGO export C.

### 2.8.2.2 export interface C

Chúng tôi đã mô hình hóa đối tượng C++ theo interface C và trừu tượng hóa một tập các interface C để mô tả đối tượng Person. Tạo một `person_capi.h` file tương ứng với  file đặc tả interface C:

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

Sau đó, các hàm C này được thực hiện bằng ngôn ngữ Go.

Cần lưu ý rằng khi export các hàm C thông qua CGO, cả hai kiểu của tham số đầu vào và kiểu của giá trị trả về đều không hỗ trợ sửa đổi hằng số *const* và cũng không hỗ trợ các hàm có tham số biến. Đồng thời, như được mô tả trong phần trước ([chương 2.7](./ch2-07-cgo-mem.md)), chúng ta không thể truy cập trực tiếp các đối tượng bộ nhớ Go trong C/C++ trong một thời gian dài. Vì vậy, chúng tôi đã sử dụng kỹ thuật được mô tả trong phần trước để ánh xạ đối tượng Go thành một id số nguyên.

Sau đây là file `person_capi.go` hiện thực các hàm trong  interface C:

```go
// person_capi.go
package main

//#include "./person_capi.h"
import "C"
import "unsafe"

//export person_new
func person_new(name *C.char, age C.int) C.person_handle_t {
    id := NewObjectId(NewPerson(C.GoString(name), int(age)))
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

Sau khi tạo đối tượng Go ta ánh xạ tới id thông qua `NewObjectId`. Sau đó, buộc id phải được exit dưới dạng `person_handle_t`. Các hàm interface khác dựa trên id được thể hiện bởi `person_handle_t`, nhờ đó đối tượng Go tương ứng được parse theo id.

### 2.8.2.3 Đóng gói các đối tượng C++

Đóng gói các đối tượng C++ với interface C tương đối đơn giản. Một cách thực hiện phổ biến là tạo một class Person mới, chứa một thành viên của kiểu `person_handle_t` tương ứng với đối tượng Go, sau đó tạo một đối tượng Go thông qua interface C trong hàm tạo (constructor) của class Person và giải phóng đối tượng Go qua interface C trong hàm hủy (destructor). Đây là một hiện thực sử dụng kỹ thuật này:

```c++
extern "C" {
    #include "./person_capi.h"
}

struct Person {
    person_handle_t goobj_;

    Person(const char* name, int age) {
        this->goobj_ = person_new((char*)name, age);
    }
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

### 2.8.2.4 Cải tiến đóng gói đối tượng C++

Trong lần hiện thực đóng gói các đối tượng C++ trước đây, mỗi lần tạo một instance Person mới, ta cần thực hiện hai lần cấp phát bộ nhớ: một lần cho phiên bản Person của C++ và một lần nữa cho phiên bản Person của ngôn ngữ Go. Trong thực tế, phiên bản C++ của Person chỉ có một id thuộc kiểu `person_handle_t`, được sử dụng để ánh xạ các đối tượng Go. Chúng ta có thể sử dụng `person_handle_t` trực tiếp trong đối tượng C++.

Các phương pháp đóng gói được cải tiến như sau đây:

```c++
extern "C" {
    #include "./person_capi.h"
}

struct Person {
    static Person* New(const char* name, int age) {
        return (Person*)person_new((char*)name, age);
    }
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

Ta thêm một hàm thành viên static mới vào class Person để tạo một cá thể Person mới. Trong hàm `new`, instance Person được tạo bằng cách gọi `person_new`, trả về kiểu `person_handle_tid` và chúng ta sử dụng nó làm con trỏ kiểu `Person*`. Trong các hàm thành viên khác, ta đảo ngược việc chuyển đổi con trỏ này thành một kiểu `person_handle_t` và sau đó gọi hàm tương ứng thông qua interface C.

Ở thời điểm này, ta đã đạt được mục tiêu export đối tượng Go dưới dạng interface C và sau đó đóng gói lại dưới dạng đối tượng C++ dựa trên interface C.

## 2.8.3 giải phóng hoàn toàn con trỏ `this` của C++

Qua việc quen thuộc với sử dụng ngôn ngữ Go sẽ cho thấy rằng các phương thức trong ngôn ngữ Go bị ràng buộc kiểu. Ví dụ nếu chúng ta xác định kiểu `Int` mới dựa trên int, chúng ta có thể có phương thức riêng:

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

Class `Int` mới được thêm vào thêm phương thức `Twice` nhưng mất quyền  chuyển về kiểu int. Tại thời điểm này, không chỉ `printf` không thể tự export giá trị của `Int` mà còn mất tất cả các tính năng của operation kiểu int. `this` là chủ ý của ngôn ngữ C++: để đổi lợi ích từ việc sử dụng class lấy cái giá là mất tất cả các tính năng ban đầu của nó.

Nguyên nhân gốc rễ của vấn đề này là do kiểu con trỏ được cố định vào class trong C++. Hãy xem xét lại bản chất của `this` trong ngôn ngữ Go:

```go
func (this Int) Twice() int
func Int_Twice(this Int) int
```

Trong Go, kiểu của tham số receiver có hàm tương tự như `this` chỉ là một tham số hàm bình thường. Chúng ta có thể tự do chọn giá trị hoặc kiểu con trỏ.

Nếu bạn nghĩ theo thuật ngữ C, `this` chỉ là một con trỏ `void*` tới một kiểu bình thường và chúng ta có thể tự do chuyển đổi nó thành các kiểu khác.

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
