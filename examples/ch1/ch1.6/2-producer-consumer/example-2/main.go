package main

import (
    "fmt"
    "os"
    "os/signal"
	"syscall"
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

    // Ctrl+C để thoát
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    fmt.Printf("quit (%v)\n", <-sig)
}