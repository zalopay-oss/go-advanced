package main

import "fmt"

func main() {
	var (
		a []int               // nil slice, equal to nil, generally used to represent a non-existent slice
		b = []int{}           // empty slice, not equal to nil, generally used to represent an empty set
		c = []int{1, 2, 3}    // There are 3 elements of the slice, both len and cap are 3
		d = c[:2]             // There are 2 elements of the slice, len is 2, cap is 3
		e = c[0:2:cap(c)]     // There are 2 elements of the slice, len is 2, cap is 3
		f = c[:0]             // There are 0 elements of the slice, len is 0, cap is 3
		g = make([]int, 3)    // There are 3 elements of the slice, len and cap are 3
		h = make([]int, 2, 3) // there are 2 elements of the slice, len is 2, cap is 3
		i = make([]int, 0, 3) // There are 0 elements of the slice, len is 0, cap is 3
	)

	fmt.Println(a, b, c, d, e, f, g, h, i)
}
