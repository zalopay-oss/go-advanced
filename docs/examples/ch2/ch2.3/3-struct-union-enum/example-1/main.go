package main

/*
struct A {
    int i;
    float f;
};
*/
import "C"
import "fmt"

func main() {
	var a C.struct_A
	fmt.Println(a.i)
	fmt.Println(a.f)
}
