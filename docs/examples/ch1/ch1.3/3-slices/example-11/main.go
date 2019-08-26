package main

import "fmt"

func main() {

	N := 2
	var a = []int{1, 2, 3, 4, 5, 6}
	a = a[:copy(a, a[1:])] // xóa phần tử đầu tiên
	a = a[:copy(a, a[N:])] // xóa N phần tử đầu tiên

	fmt.Println(a)
}
