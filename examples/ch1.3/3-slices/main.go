package main

import "fmt"

type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

func main() {
	var (
		a []int               // nil slice, equal to nil, generally used to represent a non-existent slice
		b = []int{}           // empty slice, not equal to nil, generally used to represent an empty set
		c = []int{1, 2, 3}    // There are 3 elements of the slice, both len and cap are 3
		d = c[:2]             // There are 2 elements of the slice, len is 2, cap is 3
		e = c[0:2:cap(c)]     // There are 2 elements of the slice, len is 2, cap is 3
		f = c[:0]             // There are 0 elements of the slice, len is 0, cap is 3
		g = make([]int, 3)    // There are 3 elements of the slice, len and cap are 3
		h = make([]int, 2, 3) // there are 2 elements of the slice, len is 2, cap is 3
		i = make([]int, 0, 3) // There are 0 elements of the slice, len is 0, cap is 3
	)
	////////////////////////////////////////////////////////
	for i := range a {
		fmt.Printf("a[%d]: %d\n", i, a[i])
	}
	for i, v := range b {
		fmt.Printf("b[%d]: %d\n", i, v)
	}
	for i := 0; i < len(c); i++ {
		fmt.Printf("c[%d]: %d\n", i, c[i])
	}
	/////////////////////////////////////////////////////////
	a = append(a, 1)                 // nối thêm phần tử 1
	a = append(a, 1, 2, 3)           // nối thêm phần tử 1, 2, 3
	a = append(a, []int{1, 2, 3}...) // nối thêm các phần tử 1, 2, 3 bằng cách truyền vào một mảng
	/////////////////////////////////////////////////////////

	var a = []int{1, 2, 3}
	a = append([]int{0}, a...)          // thêm phần tử 0 vào đầu slice a
	a = append([]int{-3, -2, -1}, a...) // thêm các phần tử -3, -2, -1 vào đầu slice a
	/////////////////////////////////////////////////////////
	var a []int
	a = append(a[:i], append([]int{x}, a[i:]...)...)       // chèn x ở vị trí thứ i
	a = append(a[:i], append([]int{1, 2, 3}, a[i:]...)...) // chèn một slice con vào slice ở vị trí thứ i
	//////////////////////////////////////////////////////////
	a = append(a, 0)
	copy(a[i+1:], a[i:]) // lùi những phần tử từ i trở về sau của a
	a[i] = x             // gán vị trí thứ i bằng x
	//////////////////////////////////////////////////////////
	a = append(a, x...)       // mở rộng không gian của slice a với array x
	copy(a[i+len(x):], a[i:]) // sao chép len(x) phần tử lùi về sau
	copy(a[i:], x)            // sao chép array x vào giữa
	/////////////////////////////////////////////////////////
	a = []int{1, 2, 3}
	a = a[:len(a)-1] // xóa một phần tử ở cuối
	a = a[:len(a)-N] // xóa N phần tử ở cuối
	/////////////////////////////////////////////////////////
	a = []int{1, 2, 3}
	a = a[1:] // xóa phần tử đầu tiên
	a = a[N:] // xóa N phần tử đầu tiên
	////////////////////////////////////////////////////////
	a = []int{1, 2, 3}
	a = append(a[:0], a[1:]...) // xóa phần tử đầu tiên
	a = append(a[:0], a[N:]...) // xóa N phần tử đầu tiên
	////////////////////////////////////////////////////////
	a = []int{1, 2, 3}
	a = a[:copy(a, a[1:])] // xóa phần tử đầu tiên
	a = a[:copy(a, a[N:])] // xóa N phần tử đầu tiên
	///////////////////////////////////////////////////////
	a = []int{1, 2, 3, ...}

	a = append(a[:i], a[i+1:]...) //  xóa phần tử ở vị trí i
	a = append(a[:i], a[i+N:]...) //  xóa N phần tử từ vị trí i

	a = a[:i+copy(a[i:], a[i+1:])]  // xóa phần tử ở vị trí i
	a = a[:i+copy(a[i:], a[i+N:])]  // xáo N phần từ từ vị trí i
	////////////////////////////////////////////////////////
	func TrimSpace(s []byte) []byte {
		b := s[:0]
		for _, x := range s {
			if x != ' ' {
				b = append(b, x)
			}
		}
		return b
	}
	////////////////////////////////////////////////////////
	func Filter(s []byte, fn func(x byte) bool) []byte {
		b := s[:0]
		for _, x := range s {
			if !fn(x) {
				b = append(b, x)
			}
		}
		return b
	}
	///////////////////////////////////////////////////////////
	func FindPhoneNumber(filename string) []byte {
		b, _ := ioutil.ReadFile(filename)
		return regexp.MustCompile("[0-9]+").Find(b)
	}
	/////////////////////////////////////////////////////
	func FindPhoneNumber(filename string) []byte {
		b, _ := ioutil.ReadFile(filename)
		b = regexp.MustCompile("[0-9]+").Find(b)
		return append([]byte{}, b...)
	}
	////////////////////////////////////////////////////////
	var a []*int{ ... }
	a = a[:len(a)-1]    // phần tử cuối cùng dù được xóa nhưng vẫn được tham chiếu, do đó cơ chế thu gom rác tự động không thu hồi nó
	//////////////////////////////////////////////////////
	var a []*int{ ... }
	a[len(a)-1] = nil // phần tử cuối cùng sẽ được gán giá trị nil
	a = a[:len(a)-1]  // xóa phần tử cuối cùng ra khỏi slice
	////////////////////////////////////////////////////////
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
	/////////////////////////////////////////////////////
	
}
