package main

import "fmt"

func main() {
	b := [3]int{1, 2, 3}
	fmt.Printf("b: %T\n", b)  // b: [3]int
	fmt.Printf("b: %#v\n", b) // b: [3]int{1, 2, 3}
}
