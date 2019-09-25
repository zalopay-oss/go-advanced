package main

import (
	"fmt"
	"time"
)

func main() {
	// sử dụng từ khoá go để tạo goroutine
	go fmt.Println("Hello from another goroutine")
	fmt.Println("Hello from main goroutine")

	// chờ 1 giây để có thể chạy được goroutine
	// trước khi hàm main kết thúc
	// sau khi chương trình kết thúc tất cả gorotine bị huỷ
	time.Sleep(time.Second)
}
