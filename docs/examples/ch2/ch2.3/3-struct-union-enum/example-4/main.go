package main

/*
struct A {
    int   size: 10; // Trường bit không thể truy cập
    float arr[];    // Mạng có độ dài bằng 0 cũng không thể truy cập được
};
*/
import "C"
import "fmt"

func main() {
	var a C.struct_A
	fmt.Println(a.size) // Lỗi không thể truy cập trường bit
	fmt.Println(a.arr)  // Lỗi mảng có độ dài bằng 0
}
