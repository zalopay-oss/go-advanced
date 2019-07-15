package main

import "fmt"

func main() {

	var a [3]int
	var b = [...]int{1, 2, 3}
	var c = [...]int{2: 3, 1: 2}

	for i := range a {
		fmt.Printf("a[%d]: %d\n", i, a[i])
	}
	for i, v := range b {
		fmt.Printf("b[%d]: %d\n", i, v)
	}
	for i := 0; i < len(c); i++ {
		fmt.Printf("c[%d]: %d\n", i, c[i])
	}
}
