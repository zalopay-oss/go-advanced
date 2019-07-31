package main

import (
    "fmt"
)

func main() {
    done := make(chan int, 1) // pipeline cache

    go func(){
        fmt.Println("Hello World")
        done <- 1
    }()

    <-done
}