# 2.6 Đóng package một hàm `qsort`

Hàm quick sort (`qsort`) là một hàm bậc cao của ngôn ngữ C. Nó sử dụng các hàm so sánh sắp xếp tùy chỉnh và có thể sắp xếp bất kỳ kiểu mảng nào. Trong phần này, chúng tôi sẽ cố gắng package gọn một phiên bản ngôn ngữ Go của hàm `qsort` dựa trên hàm `qsort` của ngôn ngữ C.

## 2.6.1 Tìm hiểu về hàm `qsort`

Hàm `qsort` được cung cấp bởi thư viện chuẩn <stdlib.h>. Khai báo hàm như sau:

```c
void qsort(
    void* base, size_t num, size_t size,
    int (*cmp)(const void*, const void*)
);
```

Tham số `base` là địa chỉ phần tử đầu tiên của mảng được sắp xếp, `num` là số phần tử trong mảng và `size` là kích thước của mỗi phần tử trong mảng. Điểm đáng lưu ý là hàm so sánh `cmp`, được sử dụng để sắp xếp bất kỳ hai phần tử nào trong mảng. Hai tham số con trỏ của hàm sắp xếp `cmp` là địa chỉ của hai phần tử được so sánh. Nếu phần tử tương ứng của tham số thứ nhất lớn hơn phần tử tương ứng của tham số thứ hai, kết quả lớn hơn 0. Nếu hai phần tử bằng nhau thì trả về 0. Trường hợp còn lại thì trả về kết quả bé hơn 0.

Ví dụ sau sắp xếp một mảng kiểu int với `qsort` trong C:

```c
#include <stdio.h>
#include <stdlib.h>

#define DIM(x) (sizeof(x)/sizeof((x)[0]))

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
```

Trong đó macro `DIM(values)` được sử dụng để tính các phần tử mảng, `sizeof(values[0])` để tính kích thước của phần tử mảng. `cmp` là một hàm callback so sánh kích thước của hai phần tử khi sắp xếp. Để tránh làm lộn xộn global namespace  chúng ta xác định hàm `cmp` là một hàm static chỉ có thể truy cập trong file hiện tại.

## 2.6.2 Export hàm `qsort` từ Go package

Để tạo điều kiện cho người dùng không phải CGO của ngôn ngữ Go sử dụng hàm `qsort`, chúng ta cần phải bọc hàm `qsort` của ngôn ngữ C dưới dạng hàm Go có thể truy cập được từ bên ngoài.

Đóng package hàm `qsort` như một hàm `qsort.Sort` trong Go :

```go
package qsort

//typedef int (*qsort_cmp_func_t)(const void* a, const void* b);
import "C"
import "unsafe"

func Sort(base unsafe.Pointer, num, size C.size_t, cmp C.qsort_cmp_func_t) {
    C.qsort(base, num, size, cmp)
}
```

Kiểu của  hàm so sánh được xác định là `qsort_cmp_func_t`  trong không gian ngôn ngữ C.

Mặc dù hàm `Sort` đã được export, nhưng hàm này không available cho người dùng ở bên ngoài package qsort. Các tham số của hàm `Sort` cũng chứa các  kiểu  được cung cấp bởi package C ảo. Như chúng tôi đã đề cập trong phần trước ([chương 2.5](./ch2-05-internal-mechanisms.md)), bất kỳ tên nào trong package C ảo sẽ thực sự được ánh xạ thành một tên riêng trong package. Ví dụ, `C.size_t` sẽ được mở rộng thành `_Ctype_size_t`,  kiểu `C.qsort_cmp_func_t` sẽ mở rộng thành `_Ctype_qsort_cmp_func_t`.

Hàm `Sort` có kiểu dữ liệu đã được xử lý bởi CGO như sau:

```go
func Sort(
    base unsafe.Pointer, num, size _Ctype_size_t,
    cmp _Ctype_qsort_cmp_func_t,
)
```

Điều này sẽ khiến package không thể được sử dụng từ bên ngoài do các tham số không thể khởi tạo từ kiểu  `_Ctype_size_t` và `_Ctype_qsort_cmp_func_t` vì vậy mà hàm `Sort` cũng không thể được sử dụng. Vì thế các tham số và giá trị trả về của hàm `Sort` được export cần phải  tránh phụ thuộc vào package C ảo.

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

Chúng ta thay thế kiểu trong package C ảo bằng kiểu ngôn ngữ Go và chuyển đổi lại thành kiểu được yêu cầu bởi hàm C khi gọi hàm. Do đó, người dùng bên ngoài sẽ không còn phụ thuộc vào các package C ảo trong package qsort.

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

Việc đóng gói package ban đầu của qsort cho ngôn ngữ C đã được hiện thực và có thể được người dùng khác sử dụng thông qua package đó. Tuy nhiên, hàm `qsort.Sort` có rất nhiều bất tiện vì người dùng cần cung cấp hàm so sánh trong C, đây là một việc không dễ đối với nhiều người dùng ngôn ngữ Go. Cho nên tiếp theo sau đây chúng ta sẽ tiếp tục cải tiến hàm wrapper của hàm qsort, cố gắng thay thế hàm so sánh trong C bằng hàm closure. Từ đó hướng đến bỏ đi sự phụ thuộc  của người dùng vào code CGO.

## 2.6.3 Cải tiến 1: hàm so sánh closure

Trước khi đi vào chi tiết, chúng ta sẽ xem xét interface của hàm `Sort` đi kèm với package sort trong Go:

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

Bây giờ việc sắp xếp không còn cần phải hiện thực phiên bản ngôn ngữ C của hàm so sánh thông qua CGO, bạn có thể chuyển hàm closure của ngôn ngữ Go làm hàm so sánh. Nhưng hàm `Sort` được import vẫn dựa vào package `unsafe`, điều này đi ngược lại thói quen lập trình ngôn ngữ Go.

## 2.6.4 Cải tiến 2: bỏ đi sự phụ thuộc vào package unsafe

Phiên bản trước của hàm wrapper `qsort.Sort` dễ sử dụng hơn nhiều so với phiên bản ngôn ngữ C gốc của qsort, nhưng vẫn giữ lại nhiều chi tiết về cấu trúc dữ liệu cơ bản của ngôn ngữ C. Bây giờ chúng tôi sẽ tiếp tục cải tiến hàm này, cố gắng loại bỏ sự phụ thuộc vào package `unsafe` và hiện thực hàm `Sort` tương tự như `sort.Slice` trong thư viện chuẩn.

Hàm wrapper mới được khai báo như sau:

```go
package qsort

func Slice(slice interface{}, less func(a, b int) bool)
```

Đầu tiên, chúng ta truyền slice dưới dạng tham số kiểu interface để có thể thích ứng với các kiểu slice khác nhau. Sau đó, địa chỉ, số lượng phần tử và kích thước của phần tử đầu tiên của slice có thể được lấy từ slice bằng package reflection.

Để lưu thông tin ngữ cảnh sắp xếp cần thiết, chúng ta cần tăng  số lượng phần tử và kích thước của phần tử trong biến package global. Hàm so sánh được thay đổi thành như sau:

```go
var go_qsort_compare_info struct {
    base     unsafe.Pointer
    elemnum  int
    elemsize int
    less     func(a, b int) bool
    sync.Mutex
}
```

Hàm so sánh  cần tính toán chỉ số mảng của các phần tử tương ứng theo con trỏ tới phần tử, địa chỉ bắt đầu của mảng được sắp xếp và kích thước của phần tử, sau đó trả về kết quả so sánh theo định dạng giống với kết quả trả về của hàm `less`

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

Trước tiên, cần xác định rằng kiểu interface được truyền vào phải là kiểu slice. Sau đó lấy thông tin slice cần thiết của hàm qsort thông qua reflection và gọi hàm qsort của ngôn ngữ C.

Dựa trên hàm mới được bọc, chúng ta có thể sắp xếp các slice theo cách tương tự với thư viện chuẩn:

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

Để tránh thông tin ngữ cảnh của mảng được sắp xếp là `go_qsort_compare_info` bị sửa đổi trong quá trình sắp xếp, chúng tôi đã thực hiện lock global. Do đó, phiên bản hiện tại của hàm `qsort.Slice` không thể được thực thi đồng thời. Bạn đọc có thể thử tự cải tiến giới hạn này.
