package main

import "fmt"

func main() {
	var a = []int{1, 2, 3, 4, 5, 6}
	N := 2
	a = a[1:] // xóa phần tử đầu tiên
	a = a[N:] // xóa N phần tử đầu tiên
	fmt.Println(a)
}
