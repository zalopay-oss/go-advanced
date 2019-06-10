package main

import "fmt"

func main() {
	for i, c := range []byte("世界abc") {
		fmt.Println(i, c)
	}
}
