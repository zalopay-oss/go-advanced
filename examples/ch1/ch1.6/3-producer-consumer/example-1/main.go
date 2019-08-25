package main

import (
    "fmt"
    "time"
)

// Producer: tạo ra một chuỗi số nguyên dựa trên bội số factor
func Producer(factor int, out chan<- int) {
    for i := 0; ; i++ {
        out <- i*factor
    }
}

// Consumer
func Consumer(in <-chan int) {
    for v := range in {
        fmt.Println(v)
    }
}
func main() {
    ch := make(chan int, 64) // hàng đợi kết quả

    go Producer(3, ch) // Tạo một chuỗi số với bội số 3
    go Producer(5, ch) // Tạo một chuỗi số với bội số 5
    go Consumer(ch)    // Tạo consumer

    // Thoát ra sau khi chạy trong một khoảng thời gian nhất định
    time.Sleep(5 * time.Second)
}