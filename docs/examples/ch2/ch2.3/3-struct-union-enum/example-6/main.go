package main

/*
#include <stdint.h>

union B {
    int i;
    float f;
};
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	var b C.union_B
	fmt.Println("b.i:", *(*C.int)(unsafe.Pointer(&b)))
	fmt.Println("b.f:", *(*C.float)(unsafe.Pointer(&b)))
}
