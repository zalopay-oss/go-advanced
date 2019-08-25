package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type StringHeader struct {
	Data uintptr
	Len  int
}

func main() {
	var data = []byte{
		'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd',
	}

	fmt.Println(string(data))

	s := "hello world"
	hello := s[:5]
	world := s[6:]

	fmt.Println(hello + world)
	s1 := "hello world"[:5]
	s2 := "hello world"[6:]

	fmt.Println("len(s): ", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len)
	fmt.Println("len(s1): ", (*reflect.StringHeader)(unsafe.Pointer(&s1)).Len)
	fmt.Println("len(s2): ", (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len)
}
