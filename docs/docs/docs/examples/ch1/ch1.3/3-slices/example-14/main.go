package main

import "fmt"

func Filter(s []byte, fn func(x byte) bool) []byte {
	b := s[:0]
	for _, x := range s {
		if !fn(x) {
			b = append(b, x)
		}
	}
	return b
}

func main() {
	var s = "thoainguyen"

	b := Filter([]byte(s), func(x byte) bool {
		if x == 'o' {
			return true
		}
		return false
	})

	fmt.Println(string(b))
}
