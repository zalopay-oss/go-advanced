package main

import (
    "fmt"
)

func main() {
    done := make(chan int)

    go func(){
        fmt.Println("Hello World")
        <-done
    }()

    done <- 1
}