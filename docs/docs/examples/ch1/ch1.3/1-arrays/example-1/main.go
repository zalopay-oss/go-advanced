package main

import "fmt"

func main() {

	var a [3]int                    // Định nghĩa một mảng kiểu int độ dài 3, các phần tử đều bằng 0
	var b = [...]int{1, 2, 3}       // Định nghĩa một mảng có ba phần tử 1, 2, 3, do đó độ dài là 3
	var c = [...]int{2: 3, 1: 2}    // Mảng này có 3 phần tử theo thứ tự là 0, 2, 3
	var d = [...]int{1, 2, 4: 5, 6} // Mảng này chứa dãy các phần tử là 1, 2, 0 , 0, 5, 6

	fmt.Println(a, b, c, d)
}
