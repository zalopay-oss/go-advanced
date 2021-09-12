package main

/*
struct A {
    int   type;  // type là một từ khóa trong Go
    float _type; // chặn CGO truy cập type trên kia
};
*/
import "C"
import "fmt"

func main() {
	var a C.struct_A
	fmt.Println(a._type) // _type tương ứng với _type
}
