# 2.6. Tạo ra package `qsort`

Hàm quick sort (`qsort`) là một hàm bậc cao ([higher-order function](https://en.wikipedia.org/wiki/Higher-order_function)) của ngôn ngữ C. Nó sử dụng các hàm so sánh để sắp xếp có thể tùy chỉnh và có thể sắp xếp bất kỳ kiểu mảng nào. Trong phần này, chúng ta sẽ thử tạo một package phiên bản ngôn ngữ Go của hàm `qsort` dựa trên hàm `qsort` của ngôn ngữ C.

## 2.6.1. Tìm hiểu về hàm `qsort`

Hàm `qsort` được cung cấp bởi thư viện chuẩn <stdlib.h>:

```c
void qsort(
    // `base` là địa chỉ phần tử đầu tiên của mảng
    // `num` là số phần tử
    // `size` là kích thước của mỗi phần tử
    void* base, size_t num, size_t size,

    // hàm so sánh sử dụng để sắp xếp hai phần tử bất kỳ
    int (*cmp)(const void*, const void*)
    // hai tham số con trỏ trong hàm là địa chỉ của
    // hai phần tử được so sánh. Nếu phần tử tương ứng
    // của tham số thứ nhất lớn hơn phần tử tương ứng
    // của tham số thứ hai thì kết quả lớn hơn 0
);
```

Ví dụ sau sắp xếp một mảng kiểu int với `qsort` trong C:

```c
#include <stdio.h>
#include <stdlib.h>

// macro sử dụng để tính các phần tử trong mảng
#define DIM(x) (sizeof(x)/sizeof((x)[0]))

// cmp là hàm callback so sánh kích thước của hai phần tử.
// Để tránh làm lộn xộn global namespace  chúng ta
// xác định hàm `cmp` là một hàm static chỉ có thể
// truy cập trong file hiện tại.
static int cmp(const void* a, const void* b) {
    const int* pa = (int*)a;
    const int* pb = (int*)b;
    return *pa - *pb;
}

int main() {
    int values[] = { 42, 8, 109, 97, 23, 25 };
    int i;

    qsort(values, DIM(values), sizeof(values[0]), cmp);

    for(i = 0; i < DIM(values); i++) {
        printf ("%d ",values[i]);
    }
    return 0;
}

// kết quả 8 23 25 42 97 109
```

## 2.6.2. Hiện thực `qsort` bằng Go

Để tạo điều kiện cho người dùng không phải CGO của ngôn ngữ Go sử dụng hàm `qsort`, chúng ta cần phải bọc hàm `qsort` của ngôn ngữ C dưới dạng hàm Go có thể truy cập được từ bên ngoài.

Tạo package cho hàm `qsort`:

```go
package qsort

//typedef int (*qsort_cmp_func_t)(const void* a, const void* b);
import "C"
import "unsafe"

func Sort(base unsafe.Pointer, num, size C.size_t, cmp C.qsort_cmp_func_t) {
    C.qsort(base, num, size, cmp)
}
```

Kiểu của  hàm so sánh được xác định là `qsort_cmp_func_t`  trong  ngôn ngữ C.

Mặc dù hàm `Sort` đã được export, nhưng hàm này không available cho người dùng ở bên ngoài package qsort. Các tham số của hàm `Sort` cũng chứa các  kiểu  được cung cấp bởi package C ảo. Như chúng tôi đã đề cập trong phần trước ([chương 2.5](./ch2-05-internal-mechanisms.md)), bất kỳ tên nào trong package C ảo sẽ thực sự được ánh xạ thành một tên riêng trong package. Ví dụ, `C.size_t` sẽ được mở rộng thành `_Ctype_size_t`,  kiểu `C.qsort_cmp_func_t` sẽ mở rộng thành `_Ctype_qsort_cmp_func_t`.

Hàm `Sort` có kiểu dữ liệu đã được CGO xử lý như sau:

```go
func Sort(
    base unsafe.Pointer, num, size _Ctype_size_t,
    cmp _Ctype_qsort_cmp_func_t,
)
```

Điều này sẽ khiến package không thể được sử dụng từ bên ngoài do các tham số không thể khởi tạo từ kiểu  `_Ctype_size_t` và `_Ctype_qsort_cmp_func_t` vì vậy mà hàm `Sort` cũng không thể được sử dụng. Các tham số và giá trị trả về của hàm `Sort` được export cần phải  tránh phụ thuộc vào package C ảo.

Điều chỉnh lại kiểu của tham số và triển khai hàm `Sort` như sau:

```go
/*
#include <stdlib.h>

typedef int (*qsort_cmp_func_t)(const void* a, const void* b);
*/
import "C"
import "unsafe"

type CompareFunc C.qsort_cmp_func_t

func Sort(base unsafe.Pointer, num, size int, cmp CompareFunc) {
    C.qsort(base, C.size_t(num), C.size_t(size), C.qsort_cmp_func_t(cmp))
}
```

Chúng ta thay thế kiểu trong package C ảo bằng kiểu của ngôn ngữ Go và chuyển đổi lại thành kiểu được yêu cầu bởi hàm C khi gọi hàm. Do đó, người dùng bên ngoài sẽ không còn phụ thuộc vào package C ảo trong package qsort.

Đoạn mã sau cho biết cách sử dụng hàm `Sort`:

```go
package main

//extern int go_qsort_compare(void* a, void* b);
import "C"

import (
    "fmt"
    "unsafe"

    qsort "./qsort"
)

//export go_qsort_compare
func go_qsort_compare(a, b unsafe.Pointer) C.int {
    pa, pb := (*C.int)(a), (*C.int)(b)
    return C.int(*pa - *pb)
}

func main() {
    values := []int32{42, 9, 101, 95, 27, 25}

    qsort.Sort(unsafe.Pointer(&values[0]),
        len(values), int(unsafe.Sizeof(values[0])),
        qsort.CompareFunc(C.go_qsort_compare),
    )
    fmt.Println(values)
}
```

Để sử dụng hàm `Sort`, chúng ta cần lấy thông tin của địa chỉ phần tử đầu tiên, số lượng phần tử, kích thước của phần tử trong ngôn ngữ Go làm tham số cho hàm gọi và đồng thời cung cấp hàm so sánh của đặc tả ngôn ngữ C. Trong đó `go_qsort_compare` được hiện thực bằng ngôn ngữ Go và được export sang hàm  C.

Việc đóng gói package ban đầu của qsort cho ngôn ngữ C đã được hiện thực và có thể được người dùng khác sử dụng thông qua package đó.

Tuy nhiên, hàm `qsort.Sort` có rất nhiều bất tiện vì người dùng cần cung cấp hàm so sánh trong C. Cho nên tiếp theo sau đây chúng ta sẽ tiếp tục cải tiến hàm wrapper của hàm qsort, cố gắng thay thế hàm so sánh trong C bằng hàm closure. Từ đó hướng đến bỏ đi sự phụ thuộc  của người dùng vào code CGO.

### Cải tiến 1: Loại bỏ hàm so sánh

Trước khi đi vào chi tiết, chúng ta sẽ xem xét interface của hàm `Slice` đi kèm với [package sort](https://godoc.org/github.com/golang/go/src/sort#Slice) trong Go:

```go
func Slice(slice interface{}, less func(i, j int) bool)
```

`Sort.Slice` của thư viện chuẩn rất đơn giản để sắp xếp các slice vì nó hỗ trợ chức năng so sánh được chỉ định bởi hàm closure:

```go
import (
    "fmt"
    "sort"
)
func main() {
    values := []int32{42, 9, 101, 95, 27, 25}

    sort.Slice(values, func(i, j int) bool {
        return values[i] < values[j]
    })

    fmt.Println(values)
}
```

Chúng ta cũng sẽ bọc hàm qsort của ngôn ngữ C dưới dạng hàm ngôn ngữ Go theo định dạng sau:

```go
package qsort

func Sort(base unsafe.Pointer, num, size int, cmp func(a, b unsafe.Pointer) int)
```

Hàm closure không thể được export dưới dạng hàm ngôn ngữ C, vì vậy hàm closure không thể được truyền trực tiếp sang hàm qsort của ngôn ngữ C. Để làm điều này, chúng ta có thể khởi tạo một hàm proxy có thể được export sang C bằng cách sử dụng Go và tạm thời lưu hàm so sánh closure hiện tại vào một biến toàn cục. Cụ thể như sau:

```go
var go_qsort_compare_info struct {
    fn func(a, b unsafe.Pointer) int
    sync.Mutex
}

//export _cgo_qsort_compare
func _cgo_qsort_compare(a, b unsafe.Pointer) C.int {
    return C.int(go_qsort_compare_info.fn(a, b))
}
```

Hàm ngôn ngữ C đã export `_cgo_qsort_compare` là hàm so sánh qsort được public, bên trong `go_qsort_compare_info.fn` gọi hàm so sánh closure hiện tại.

Hàm `Sort` mới được hiện thực như sau:

```go
/*
#include <stdlib.h>

typedef int (*qsort_cmp_func_t)(const void* a, const void* b);
extern int _cgo_qsort_compare(void* a, void* b);
*/
import "C"

func Sort(base unsafe.Pointer, num, size int, cmp func(a, b unsafe.Pointer) int) {
    go_qsort_compare_info.Lock()
    defer go_qsort_compare_info.Unlock()

    go_qsort_compare_info.fn = cmp

    C.qsort(base, C.size_t(num), C.size_t(size),
        C.qsort_cmp_func_t(C._cgo_qsort_compare),
    )
}
```

Trước mỗi lần sắp xếp, lock biến toàn cục `go_qsort_compare_info`, lưu hàm closure hiện tại vào biến toàn cục và sau đó gọi hàm qsort của ngôn ngữ C.

Dựa trên hàm mới được bọc, chúng ta có thể đơn giản hóa code sắp xếp trước đó:

```go
func main() {
    values := []int32{42, 9, 101, 95, 27, 25}

    qsort.Sort(unsafe.Pointer(&values[0]), len(values), int(unsafe.Sizeof(values[0])),
        func(a, b unsafe.Pointer) int {
            pa, pb := (*int32)(a), (*int32)(b)
            return int(*pa - *pb)
        },
    )

    fmt.Println(values)
}
```

Bây giờ việc sắp xếp không còn cần phải hiện thực phiên bản ngôn ngữ C của hàm so sánh thông qua CGO, bạn có thể chuyển hàm closure của ngôn ngữ Go làm hàm so sánh. Nhưng hàm `Sort` được import vẫn dựa vào package `unsafe`.

### Cải tiến 2: Loại bỏ sự phụ thuộc vào package unsafe

Phần này chúng ta sẽ thử loại bỏ sự phụ thuộc vào package `unsafe` và hiện thực hàm `Sort` tương tự như `sort.Slice` trong thư viện chuẩn.

Hàm wrapper mới được khai báo như sau:

```go
package qsort

func Slice(slice interface{}, less func(a, b int) bool)
```

Đầu tiên, chúng ta truyền slice dưới dạng tham số kiểu interface để có thể tương thích với các kiểu slice khác nhau. Sau đó, địa chỉ, số lượng phần tử và kích thước của phần tử đầu tiên của slice có thể được lấy từ slice bằng package reflection.

Struct để chuyển tham số thay đổi thành:

```go
var go_qsort_compare_info struct {
    base     unsafe.Pointer
    elemnum  int
    elemsize int
    less     func(a, b int) bool
    sync.Mutex
}
```

Hàm so sánh  cần tính toán chỉ số mảng của các phần tử tương ứng theo con trỏ tới phần tử, địa chỉ bắt đầu của mảng được sắp xếp và kích thước của phần tử, sau đó trả về kết quả so sánh theo định dạng giống với kết quả trả về của hàm `less`:

```go
//export _cgo_qsort_compare
func _cgo_qsort_compare(a, b unsafe.Pointer) C.int {
    var (
        // array memory is locked
        base     = uintptr(go_qsort_compare_info.base)
        elemsize = uintptr(go_qsort_compare_info.elemsize)
    )

    i := int((uintptr(a) - base) / elemsize)
    j := int((uintptr(b) - base) / elemsize)

    switch {
    case go_qsort_compare_info.less(i, j): // v[i] < v[j]
        return -1
    case go_qsort_compare_info.less(j, i): // v[i] > v[j]
        return +1
    default:
        return 0
    }
}
```

Hiện thực hàm Slice mới như sau:

```go
func Slice(slice interface{}, less func(a, b int) bool) {
    sv := reflect.ValueOf(slice)
    if sv.Kind() != reflect.Slice {
        panic(fmt.Sprintf("qsort called with non-slice value of type %T", slice))
    }
    if sv.Len() == 0 {
        return
    }

    // Để tránh thông tin ngữ cảnh của mảng được sắp xếp là `go_qsort_compare_info`
    // bị sửa đổi trong quá trình sắp xếp, chúng tôi đã thực hiện lock global.
    go_qsort_compare_info.Lock()
    defer go_qsort_compare_info.Unlock()

    defer func() {
        go_qsort_compare_info.base = nil
        go_qsort_compare_info.elemnum = 0
        go_qsort_compare_info.elemsize = 0
        go_qsort_compare_info.less = nil
    }()

    // baseMem = unsafe.Pointer(sv.Index(0).Addr().Pointer())
    // baseMem maybe moved, so must saved after call C.fn
    go_qsort_compare_info.base = unsafe.Pointer(sv.Index(0).Addr().Pointer())
    go_qsort_compare_info.elemnum = sv.Len()
    go_qsort_compare_info.elemsize = int(sv.Type().Elem().Size())
    go_qsort_compare_info.less = less

    C.qsort(
        go_qsort_compare_info.base,
        C.size_t(go_qsort_compare_info.elemnum),
        C.size_t(go_qsort_compare_info.elemsize),
        C.qsort_cmp_func_t(C._cgo_qsort_compare),
    )
}
```

Interface được truyền vào phải là kiểu slice. Sau đó lấy thông tin slice cần thiết của hàm qsort thông qua reflection và gọi hàm qsort của ngôn ngữ C.

Dựa trên hàm mới được bọc, chúng ta có thể sắp xếp slice theo cách tương tự với thư viện chuẩn:

```go
import (
    "fmt"

    qsort "."
)

func main() {
    values := []int64{42, 9, 101, 95, 27, 25}

    qsort.Slice(values, func(i, j int) bool {
        return values[i] < values[j]
    })

    fmt.Println(values)
}
```

Để tránh thông tin ngữ cảnh của mảng được sắp xếp là `go_qsort_compare_info` bị sửa đổi trong quá trình sắp xếp, chúng tôi đã thực hiện lock global. Do đó, phiên bản hiện tại của hàm `qsort.Slice` không thể được thực thi đồng thời. Bạn đọc có thể thử cải tiến giới hạn này.

[Tiếp theo](ch2-07-cgo-mem.md)