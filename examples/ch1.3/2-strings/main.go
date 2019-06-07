package main

import (
	"fmt"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

type StringHeader struct {
	Data uintptr
	Len  int
}

func forOnString(s string, forBody func(i int, r rune)) {
	for i := 0; len(s) > 0; {
		r, size := utf8.DecodeRuneInString(s)
		forBody(i, r)
		s = s[size:]
		i += size
	}
}

func str2bytes(s string) []byte {
	p := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		p[i] = c
	}
	return p
}

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

func str2runes(s []byte) []rune {
	var p []int32
	for len(s) > 0 {
		r, size := utf8.DecodeRune(s)
		p = append(p, int32(r))
		s = s[size:]
	}
	return []rune(p)
}

func runes2string(s []int32) string {
	var p []byte
	buf := make([]byte, 3)
	for _, r := range s {
		n := utf8.EncodeRune(buf, r)
		p = append(p, buf[:n]...)
	}
	return string(p)
}

func main() {
	///////////////////////////////////////////////
	var data = []byte{
		'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd',
	}
	fmt.Println(string(data))
	///////////////////////////////////////////////
	s := "hello, world"
	hello := s[:5]
	world := s[7:]
	fmt.Println(hello, world)
	s1 := "hello, world"[:5]
	s2 := "hello, world"[7:]
	///////////////////////////////////////////////
	fmt.Println("len(s): ", (*reflect.StringHeader)(unsafe.Pointer(&s)).Len)
	fmt.Println("len(s1): ", (*reflect.StringHeader)(unsafe.Pointer(&s1)).Len)
	fmt.Println("len(s2): ", (*reflect.StringHeader)(unsafe.Pointer(&s2)).Len)
	///////////////////////////////////////////////
	fmt.Printf("%#v\n", []byte("Hello, 世界"))
	// Kết quả là
	// []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0xe4, 0xb8, 0x96, 0xe7, 0x95, 0x8c}
	///////////////////////////////////////////////
	fmt.Println("\xe4\xb8\x96")
	fmt.Println("\xe7\x95\x8c")
	///////////////////////////////////////////////
	fmt.Println("\xe4\x00\x00\xe7\x95\x8cabc")
	///////////////////////////////////////////////
	for i, c := range "\xe4\x00\x00\xe7\x95\x8cabc" {
		fmt.Println(i, c)
	}
	///////////////////////////////////////////////
	for i, c := range []byte("世界abc") {
		fmt.Println(i, c)
	}
	///////////////////////////////////////////////
	const l = "\xe4\x00\x00\xe7\x95\x8cabc"
	for i := 0; i < len(l); i++ {
		fmt.Printf("%d %x\n", i, s[i])
	}
	///////////////////////////////////////////////
	fmt.Printf("%#v\n", []rune("世界"))             // []int32{19990, 30028}
	fmt.Printf("%#v\n", string([]rune{'世', '界'})) // 世界
	///////////////////////////////////////////////
}
