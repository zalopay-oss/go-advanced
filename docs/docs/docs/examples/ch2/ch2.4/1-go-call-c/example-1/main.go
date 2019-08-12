package main

/*
static int add(int a, int b) {
    return a+b;
}
*/
import "C"

func main() {
	C.add(1, 1)
}
