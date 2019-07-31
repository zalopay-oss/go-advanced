package main

import "fmt"

func main() {
	for i, c := range "\xe4\x00\x00\xe7\x95\x8cabc" {
		fmt.Println(i, c)
	}
	// 0 65533  // \uFFFD, 对应 �
	// 1 0      // 空字符
	// 2 0      // 空字符
	// 3 30028  // 界
	// 6 97     // a
	// 7 98     // b
	// 8 99     // c
}
