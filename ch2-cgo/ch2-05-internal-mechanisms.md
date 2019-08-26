# 2.5. Cơ chế bên trong CGO

Để hiểu rõ hơn cách CGO vận hành, trong phần này chúng ta sẽ phân tích luồng hoạt động của các hàm ngôn ngữ Go và C từ code được tạo.

## 2.5.1. Các file trung gian được CGO tạo ra

Để hiểu cơ chế  của  CGO, trước tiên ta cần biết các file trung gian mà CGO tạo ra. Ta có thể thêm một thư mục `-work` chứa file trung gian output khi build package cgo và giữ file trung gian khi quá trình build hoàn tất. Nếu là một đoạn code cgo đơn giản, chúng ta cũng có thể trực tiếp xem file trung gian được tạo bằng cách gọi lệnh `go tool cgo`.

Trong file nguồn Go, nếu một lệnh `import "C"` thực thi thì lệnh cgo sẽ được gọi để tạo ra file trung gian tương ứng. Dưới đây là sơ đồ đơn giản mô tả các file trung gian được cgo tạo ra:

<div align="center">
	<img src="../images/ch2-4-cgo-generated-files.dot.png">
	<br/>
	<span align="center">
		<i>Các file trung gian được CGO tạo ra</i>
	</span>
    <br/>
</div>
  
Có 4 file Go trong package, trong đó các file nocgo chứa `import "C"` và hai file còn lại chứa code cgo. Lệnh cgo tạo ra hai file trung gian cho mỗi file chứa mã cgo. Ví dụ: main.go tạo ra hai file trung gian là `main.cgo1.go` và `main.cgo2.c`.

Kế đó file `_cgo_gotypes.go` được tạo cho toàn bộ package chứa một phần code hỗ trợ của Go. Đồng thời quá trình này cũng  tạo ra các file  `_cgo_export.h`  và `_cgo_export.c`, để export các kiểu và hàm trong Go tới kiểu và hàm tương ứng trong C.

## 2.5.2. Go gọi hàm của C

Go gọi các hàm trong  C là  trường hợp ứng dụng phổ biến nhất của CGO. Chúng ta sẽ bắt đầu với ví dụ đơn giản nhất để phân tích chi tiết luồng hoạt động của quá trình này.

Đoạn code cụ thể như sau:

```go
//main.go
package main

//int sum(int a, int b) { return a+b; }
import "C"

func main() {
    println(C.sum(1, 1))
}
```

Tiếp theo ta tạo một file trung gian trong thư mục _obj thông qua command line cgo:

```sh
$ go tool cgo main.go
```

Vào thư mục _obj  để xem các file trung gian:

```go
$ ls _obj
_cgo_.o
_cgo_export.c
_cgo_export.h
_cgo_flags
_cgo_gotypes.go
_cgo_main.c
main.cgo1.go
main.cgo2.c
```

Trong đó `_cgo_.o`, `_cgo_flags` và `_cgo_main.c` có code không liên quan logic trực tiếp với nhau, bạn có thể bỏ qua.

Trước tiên chúng ta hãy xem file [`main.cgo1.go`](../examples/ch2/ch2.5/2-go-call-c/example-1/_obj/main.cgo1.go) chứa code Go sau khi file `main.go` expand các hàm và biến số liên quan trong package C ảo:

```go
package main

//int sum(int a, int b) { return a+b; }
import _ "unsafe"

func main() {
    println((_Cfunc_sum)(1, 1))
}
```

Lời gọi `C.sum(1, 1)` được thay thế thành `(_Cfunc_sum)(1, 1)`. Mỗi dạng `C.xxx` của hàm được thay thế bằng hàm Go thuần túy dạng `_Cfunc_xxx`, trong đó tiền tố `_Cfunc_` chỉ ra rằng đây là hàm C.

Hàm `_Cfunc_sum` được định nghĩa trong file [`_cgo_gotypes.go`](../examples/ch2/ch2.5/2-go-call-c/example-1/_obj/_cgo_gotypes.go) do CGO tạo ra như sau:

```go
//go:cgo_unsafe_args
func _Cfunc_sum(p0 _Ctype_int, p1 _Ctype_int) (r1 _Ctype_int) {
    _cgo_runtime_cgocall(_cgo_506f45f9fa85_Cfunc_sum, uintptr(unsafe.Pointer(&p0)))
    if _Cgo_always_false {
        _Cgo_use(p0)
        _Cgo_use(p1)
    }
    return
}
```

Tham số của hàm `_Cfunc_sum` và kiểu `_Ctype_int` của giá trị trả về tương ứng với kiểu `C.int`, các quy tắc đặt tên `_Cfunc_xxx` là tương tự nhau và các tiền tố khác nhau được sử dụng để phân biệt giữa các hàm và kiểu.

Hàm `_cgo_runtime_cgocall` tương ứng với `runtime.cgocall`, khai báo của hàm như sau:

```go
func runtime.cgocall(fn, arg unsafe.Pointer) int32
```

Tham số đầu tiên là địa chỉ của hàm ngôn ngữ C và tham số thứ hai là địa chỉ của struct tham số tương ứng với hàm ngôn ngữ C.

Trong ví dụ này, hàm được truyền vào: `_cgo_506f45f9fa85_Cfunc_sum` cũng là một hàm trung gian được CGO tạo ra. Hàm  được định nghĩa trong `main.cgo2.c1`:

```c
void _cgo_506f45f9fa85_Cfunc_sum(void *v) {
    struct {
        int p0;
        int p1;
        int r;
        char __pad12[4];
    } __attribute__((__packed__)) *a = v;
    char *stktop = _cgo_topofstack();
    __typeof__(a->r) r;
    _cgo_tsan_acquire();
    r = sum(a->p0, a->p1);
    _cgo_tsan_release();
    a = (void*)((char*)a + (_cgo_topofstack() - stktop));
    a->r = r;
}
```

Tham số hàm này chỉ có một con trỏ trỏ tới kiểu void và hàm không có giá trị trả về.

Struct được trỏ  tới bởi con trỏ hàm  `_cgo_506f45f9fa85_Cfunc_sum` là:

```go
struct {
    int p0;
    int p1;
    int r;
    char __pad12[4];
} __attribute__((__packed__)) *a = v;
```

Thành phần p0 tương ứng với tham số đầu tiên của `sum`, thành phần p1 tương ứng với tham số thứ hai  và thành phần `__pad12` được sử dụng để điền vào struct cho mục đích  đảm bảo alignment của CPU.

Sau khi có được các tham số (trỏ tới struct), hàm `sum` của phiên bản ngôn ngữ C được gọi và giá trị trả về được lưu vào thành phần tương ứng trong thân struct.

Toàn bộ biểu đồ luồng hoạt động của cuộc gọi `C.sum` như sau:

<div align="center">
	<img src="../images/ch2-5-call-c-sum-v1.uml.png">
	<br/>
	<span align="center">
		<i>Gọi hàm C</i>
	</span>
    <br/>
</div>

Trong đó hàm  `runtime.cgocall` là chìa khóa để thực hiện cuộc gọi  vượt ranh giới của hàm ngôn ngữ Go sang hàm ngôn ngữ C. Thông tin chi tiết có thể tham khảo <https://golang.org/src/cmd/cgo/doc.go>.

## 2.5.3. C gọi hàm của Go

Bây giờ chúng ta sẽ phân tích luồng của cuộc gọi ngược lại: C gọi đến hàm Go. Tương tự, ta cũng khởi tạo một hàm Go, tên file là sum.go:

```go
package main

//int sum(int a, int b);
import "C"

//export sum
func sum(a, b C.int) C.int {
    return a + b
}

func main() {}
```

Các chi tiết về cú pháp của CGO không được mô tả ở đây. Để sử dụng hàm `sum` trong C, chúng ta cần biên dịch mã Go vào thư viện tĩnh của C:

```sh
$ go build -buildmode=c-archive -o sum.a sum.go
```

Nếu không có lỗi, lệnh biên dịch ở trên sẽ tạo ra một thư viện tĩnh `sum.a` và file header `sum.h`. File này sẽ chứa khai báo của hàm sum và thư viện tĩnh sẽ chứa hiện thực của hàm.

Để phân tích luồng hoạt động của cuộc gọi hàm từ phiên bản ngôn ngữ C ta cũng cần phải phân tích các file trung gian do CGO tạo ra:

```sh
$ go tool cgo sum.go
```

Thư mục _obj vẫn chứa các file trung gian được tạo tương tự như phần trước. Để ngắn gọn, ta bỏ qua một vài file không liên quan:

```sh
$ ls _obj
_cgo_export.c
_cgo_export.h
_cgo_gotypes.go
main.cgo1.go
main.cgo2.c
```

Trong đó nội dung của file `_cgo_export.h` và file do C tạo ra khi nó tạo thư viện tĩnh `sum.h` là giống nhau, đều khai báo hàm sum.

Vì ngôn ngữ C là bên gọi, chúng ta cần bắt đầu với việc hiện thực phiên bản ngôn ngữ C của hàm sum. Phiên bản này nằm trong file `_cgo_export.c` (file chứa phần hiện thực hàm của C tương ứng với hàm export  của Go):

```c
int sum(int p0, int p1)
{
    __SIZE_TYPE__ _cgo_ctxt = _cgo_wait_runtime_init_done();
    struct {
        int p0;
        int p1;
        int r0;
        char __pad0[4];
    } __attribute__((__packed__)) a;
    a.p0 = p0;
    a.p1 = p1;
    _cgo_tsan_release();
    crosscall2(_cgoexp_8313eaf44386_sum, &a, 16, _cgo_ctxt);
    _cgo_tsan_acquire();
    _cgo_release_context(_cgo_ctxt);
    return a.r0;
}
```

Hàm sum sử dụng một kỹ thuật tương tự như phần trước trình bày để đóng gói các tham số và trả về các giá trị của hàm thành một  struct, sau đó truyền struct `runtime/cgo.crosscall2` vào hàm thực thi thông qua hàm `_cgoexp_8313eaf44386_sum`.

Hàm `runtime/cgo.crosscall2` được hiện thực bằng hợp ngữ và khai báo hàm của nó như sau:

```go
func runtime/cgo.crosscall2(
    fn func(a unsafe.Pointer, n int32, ctxt uintptr),
    a unsafe.Pointer, n int32,
    ctxt uintptr,
)
```

Điểm cần chú ý ở đây là `fn` và `a`, `fn` là con trỏ tới hàm trung gian (proxy) và `a` là con trỏ tới struct tương ứng với đối số truyền đi khi gọi (và cũng chứa luôn giá trị trả về).

Hàm trung gian `_cgoexp_8313eaf44386_sum` có trong file [`_cgo_gotypes.go`](../examples/ch2/ch2.5/3-c-call-go/example-1/_obj/_cgo_gotypes.go):

```go
func _cgoexp_8313eaf44386_sum(a unsafe.Pointer, n int32, ctxt uintptr) {
    fn := _cgoexpwrap_8313eaf44386_sum
    _cgo_runtime_cgocallback(**(**unsafe.Pointer)(unsafe.Pointer(&fn)), a, uintptr(n), ctxt);
}

func _cgoexpwrap_8313eaf44386_sum(p0 _Ctype_int, p1 _Ctype_int) (r0 _Ctype_int) {
    return sum(p0, p1)
}
```

Hàm bao ngoài `_cgoexpwrap_8313eaf44386_sum` của `sum`  được sử dụng như một con trỏ hàm. Hàm `_cgo_runtime_cgocallback` tương ứng với hàm `runtime.cgocallback`:

```go
func runtime.cgocallback(fn, frame unsafe.Pointer, framesize, ctxt uintptr)
```

Các tham số là con trỏ hàm, tham số hàm và giá trị trả về tương ứng với con trỏ của struct, kích thước frame của lời gọi hàm và tham số ngữ cảnh.

Toàn bộ biểu đồ luồng cuộc gọi như sau:

<div align="center">

<img src="../images//ch2-6-call-c-sum-v2.uml.png">
<br/>
<span align="center"><i>Gọi hàm Go đã export</i></span>
    <br/>

</div>

Trong đó, hàm `runtime.cgocallback` là chìa khóa để thực hiện cuộc gọi vượt ranh giới từ ngôn ngữ C sang Go. Chi tiết  có thể được tìm thấy trong hiện thực [runtime.cgocallback.go](https://github.com/golang/go/blob/master/src/runtime/cgocallback.go)

[Tiếp theo](ch2-06-qsort.md)