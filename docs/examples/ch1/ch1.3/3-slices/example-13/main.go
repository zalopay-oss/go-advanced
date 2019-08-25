package main

import "fmt"

func TrimSpace(s []byte) []byte {
	b := s[:0]
	for _, x := range s {
		if x != ' ' {
			b = append(b, x)
		}
	}
	return b
}

func main() {
	var s = []byte{' ', '1', '2', ' ', '3', '4', ' ', ' '}
	b := TrimSpace(s)

	fmt.Println(string(b))
}
