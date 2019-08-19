package main

import (
  "sync"
)

// biến toàn cục
var counter int

func main() {
  var wg sync.WaitGroup
  var l sync.Mutex
  for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
      defer wg.Done()
      l.Lock()
      counter++
      l.Unlock()
    }()
  }

  wg.Wait()
  println(counter)
}