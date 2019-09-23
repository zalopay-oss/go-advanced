package main

import (
	"fmt"
	"time"
)

func MyPrintln(id int, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		fmt.Println("Xin chào, tôi là goroutine: ", id)
	}()
}

func main() {
	for i := 0; i < 100; i++ {
		MyPrintln(i, 1*time.Second)
	}

	time.Sleep(10 * time.Second)
	fmt.Println("Chương trình kết thúc")
}
