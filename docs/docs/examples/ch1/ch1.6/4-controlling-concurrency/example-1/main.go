package main

import (
    "fmt"
)

var limit = make(chan int, 3)

<<<<<<< HEAD
func work ()(s string){
    fmt.Print("aaaaaa")
}

=======
>>>>>>> 039d41a5ffac593cb424dd3bee29b440339ea376
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