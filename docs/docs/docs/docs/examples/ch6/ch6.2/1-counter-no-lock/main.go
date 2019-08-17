package main

import (
  "sync"
)

// biến toàn cục
var counter int

func main() {
  var wg sync.WaitGroup
  for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
    defer wg.Done()
      counter++
    }()
  }

  wg.Wait()
  println(counter)
}