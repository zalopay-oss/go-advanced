package main

import "fmt"

func main() {
	var (
		a = []int{1, 2, 3}
		i = 2
		x = []int{4, 5}
	)
	a = append(a, x...)       // mở rộng không gian của slice a với array x
	copy(a[i+len(x):], a[i:]) // sao chép len(x) phần tử lùi về sau
	copy(a[i:], x)            // sao chép array x vào giữa
	fmt.Println(a)
}
