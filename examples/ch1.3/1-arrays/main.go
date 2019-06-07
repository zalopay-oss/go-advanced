package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

func main() {

	var a [3]int                    // Định nghĩa một mảng kiểu int độ dài 3, các phần tử đều bằng 0
	var b = [...]int{1, 2, 3}       // Định nghĩa một mảng có ba phần tử 1, 2, 3, do đó độ dài là 3
	var c = [...]int{2: 3, 1: 2}    // Mảng này có 3 phần tử theo thứ tự là 0, 2, 3
	var d = [...]int{1, 2, 4: 5, 6} // Mảng này chứa dãy các phần tử là 1, 2, 0 , 0, 5, 6

	fmt.Print(a, b, c, d)

	////////////////////////////////////////////////////

	fmt.Println(a[0], a[1]) // in ra hai phần tử đầu tiên của array a
	fmt.Println(b[0], b[1]) // truy xuất các phần tử của con trỏ array cũng giống như truy xuất các phần tử của array

	for i, v := range b { // duyệt qua các phần tử trong con trỏ array, giống như duyệt qua array
		fmt.Println(i, v)
	}

	/////////////////////////////////////////////////////

	for i := range a {
		fmt.Printf("a[%d]: %d\n", i, a[i])
	}
	for i, v := range b {
		fmt.Printf("b[%d]: %d\n", i, v)
	}
	for i := 0; i < len(c); i++ {
		fmt.Printf("c[%d]: %d\n", i, c[i])
	}
	/////////////////////////////////////////////////////

	var times [5][0]int
	for range times {
		fmt.Println("hello")
	}

	////////////////////////////////////////////////////

	// Mảng string
	var s1 = [2]string{"hello", "world"}
	var s2 = [...]string{"Hello!", "World"}
	var s3 = [...]string{1: "Hello", 0: "World"}

	// Mảng struct
	var line1 [2]image.Point
	var line2 = [...]image.Point{image.Point{X: 0, Y: 0}, image.Point{X: 1, Y: 1}}
	var line3 = [...]image.Point{{0, 0}, {1, 1}}

	// Mảng decoder của hình ảnh
	var decoder1 [2]func(io.Reader) (image.Image, error)
	var decoder2 = [...]func(io.Reader) (image.Image, error){
		png.Decode,
		jpeg.Decode,
	}

	// Mảng interface{}
	var unknown1 [2]interface{}
	var unknown2 = [...]interface{}{123, "Hello!"}

	// Mảng pipe
	var chanList = [2]chan int{}

	///////////////////////////////////////////////////////
}
