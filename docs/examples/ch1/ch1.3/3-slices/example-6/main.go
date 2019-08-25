package main

import "fmt"

func main() {
	var (
		a = []int{1, 2, 3}
		i = 2
		x = 3
	)
	a = append(a, 0)
	copy(a[i+1:], a[i:]) // lùi những phần tử từ i trở về sau của a
	a[i] = x             // gán vị trí thứ i bằng x
	fmt.Println(a)
}
