package main

import (
    "fmt"
)

var limit = make(chan int, 3)

func work ()(s string){
    fmt.Print("aaaaaa")
}

func main() {
    for _, w := range work {
        go func() {
            limit <- 1
            w()
            <-limit
        }()
    }
    select{}
}