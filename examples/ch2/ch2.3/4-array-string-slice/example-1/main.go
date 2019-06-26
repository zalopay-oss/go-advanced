package main

/*
static char arr[10];
static char *s = "Hello";
*/
import "C"
import (
	"reflect"
	"unsafe"
	"fmt"
)

func main() {
	// chuyển đổi bằng reflect.SliceHeader
	var arr0 []byte
	var arr0Hdr = (*reflect.SliceHeader)(unsafe.Pointer(&arr0))
	arr0Hdr.Data = uintptr(unsafe.Pointer(&C.arr[0]))
	arr0Hdr.Len = 10
	arr0Hdr.Cap = 10

	// chuyển đổi slice
	arr1 := (*[31]byte)(unsafe.Pointer(&C.arr[0]))[:10:10]

	var s0 string
	var s0Hdr = (*reflect.StringHeader)(unsafe.Pointer(&s0))
	s0Hdr.Data = uintptr(unsafe.Pointer(C.s))
	s0Hdr.Len = int(C.strlen(C.s))

	sLen := int(C.strlen(C.s))
	s1 := string((*[31]byte)(unsafe.Pointer(&C.s[0]))[:sLen:sLen])

	fmt.Println(arr1, s1)
}
