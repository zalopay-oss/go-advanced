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
            // thực hiện bình thường
        case <-cannel:
            // thoát
        }
    }
}

func main() {
    cannel := make(chan bool)
    go worker(cannel)

    time.Sleep(time.Second)
    cannel <- true
}