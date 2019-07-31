package main

/*
#include <stdint.h>

union B1 {
    int i;
    float f;
};

union B2 {
    int8_t i8;
    int64_t i64;
};
*/
import "C"
import "fmt"

func main() {
	var b1 C.union_B1
	fmt.Printf("%T\n", b1) // [4]uint8

	var b2 C.union_B2
	fmt.Printf("%T\n", b2) // [8]uint8
}
