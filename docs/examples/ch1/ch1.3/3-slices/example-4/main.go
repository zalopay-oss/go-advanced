package main

import "fmt"

func main() {
	var a = []int{1, 2, 3}
	a = append([]int{0}, a...)          // thêm phần tử 0 vào đầu slice a
	a = append([]int{-3, -2, -1}, a...) // thêm các phần tử -3, -2, -1 vào đầu slice a
	fmt.Print(a)
}
