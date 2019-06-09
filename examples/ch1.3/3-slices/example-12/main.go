package main

import "fmt"

func main() {

	i := 2
	N := 2
	var a = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	a = append(a[:i], a[i+1:]...) //  xóa phần tử ở vị trí i
	a = append(a[:i], a[i+N:]...) //  xóa N phần tử từ vị trí i

	a = a[:i+copy(a[i:], a[i+1:])] // xóa phần tử ở vị trí i
	a = a[:i+copy(a[i:], a[i+N:])] // xóa N phần từ từ vị trí i

	fmt.Println(a)
}
