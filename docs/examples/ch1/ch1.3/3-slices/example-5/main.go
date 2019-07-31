package main

import "fmt"

func main() {
	var (
		a = []int{-1, -2, -3}
		x = 3
		i = 2
	)
	a = append(a[:i], append([]int{x}, a[i:]...)...)       // chèn x ở vị trí thứ i
	a = append(a[:i], append([]int{1, 2, 3}, a[i:]...)...) // chèn một slice con vào slice ở vị trí thứ i
	fmt.Println(a)
}
