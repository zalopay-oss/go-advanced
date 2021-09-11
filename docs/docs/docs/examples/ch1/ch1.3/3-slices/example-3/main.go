package main

import "fmt"

func main() {
	var a = []int{-1, -2, -3}
	a = append(a, 1)                 // nối thêm phần tử 1
	a = append(a, 1, 2, 3)           // nối thêm phần tử 1, 2, 3
	a = append(a, []int{1, 2, 3}...) // nối thêm các phần tử 1, 2, 3 bằng cách truyền vào một mảng

	fmt.Print(a)
}
