package main

import (
    "fmt"
)

func main() {
    done := make(chan int, 10)

    // Mở ra N thread
    for i := 0; i < cap(done); i++ {
        go func(){
            fmt.Println("Hello World")
            done <- 1
        }()
    }

    // Đợi cả 10 thread hoàn thành
    for i := 0; i < cap(done); i++ {
        <-done
    }
}