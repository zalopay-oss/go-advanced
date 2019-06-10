package main

import "fmt"

func main() {
	x1 := 1
	x2 := 2
	x3 := 3
	var a = []*int{&x1, &x2, &x3}
	a = a[:len(a)-1]
	// phần tử cuối cùng dù được xóa nhưng vẫn được tham chiếu,
	// do đó cơ chế thu gom rác tự động không thu hồi nó
	fmt.Println(a)
}
