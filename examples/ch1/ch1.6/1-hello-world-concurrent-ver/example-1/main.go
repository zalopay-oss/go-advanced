package main

import (
    "fmt"
    "sync"
)

func main() {
    var mu sync.Mutex

    go func(){
        fmt.Println("Hello World")
        mu.Lock()
    }()

    mu.Unlock()
}