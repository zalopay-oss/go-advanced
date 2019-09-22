package main

import "fmt"

func main() {
	x1 := 1
	x2 := 2
	x3 := 3
	var a = []*int{&x1, &x2, &x3}
	a[len(a)-1] = nil // phần tử cuối cùng sẽ được gán giá trị nil
	a = a[:len(a)-1]  // xóa phần tử cuối cùng ra khỏi slice
	fmt.Println(x3)
}
