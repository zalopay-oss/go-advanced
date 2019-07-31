package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup

    // Mở N thread
    for i := 0; i < 10; i++ {
        wg.Add(1)

        go func() {
            fmt.Println("Hello World")
            wg.Done()
        }()
    }

    // Đợi N thread hoàn thành
    wg.Wait()
}