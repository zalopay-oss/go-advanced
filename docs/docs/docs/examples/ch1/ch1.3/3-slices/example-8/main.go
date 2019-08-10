package main

import "fmt"

func main() {
	var a = []int{1, 2, 3, 4, 5, 6}
	N := 2
	a = a[:len(a)-1] // xóa một phần tử ở cuối
	a = a[:len(a)-N] // xóa N phần tử ở cuối
	fmt.Println(a)
}
