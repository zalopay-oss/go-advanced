package main

import (
	"fmt"
)

func main() {
	// sử dụng từ khoá go để tạo goroutine
	go fmt.Println("Hello from another goroutine")
	fmt.Println("Hello from main goroutine")
}
