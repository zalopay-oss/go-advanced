package main

import "fmt"

func main() {

	var a = [...]int{1, 2, 3}
	var b = &a

	fmt.Println(a[0], a[1]) // in ra hai phần tử đầu tiên của array a
	fmt.Println(b[0], b[1]) // truy xuất các phần tử của con trỏ array cũng giống như truy xuất các phần tử của array

	for i, v := range b { // duyệt qua các phần tử trong con trỏ array, giống như duyệt qua array
		fmt.Println(i, v)
	}

}
