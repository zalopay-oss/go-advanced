package main 

import (
    "fmt"
    "time"
)

func worker(cannel chan bool) {
    for {
        select {
        default:
            fmt.Println("hello")
            // hoạt động bình thường
        case <-cannel:
            // thoát
        }
    }
}

func main() {
    cancel := make(chan bool)

    for i := 0; i < 10; i++ {
        go worker(cancel)
    }

    time.Sleep(time.Second)
    close(cancel)
}