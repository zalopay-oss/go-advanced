package main

import "fmt"

func main() {
	N := 2
	var a = []int{1, 2, 3}
	a = append(a[:0], a[1:]...) // xóa phần tử đầu tiên
	a = append(a[:0], a[N:]...) // xóa N phần tử đầu tiên

	fmt.Println(a)
}
