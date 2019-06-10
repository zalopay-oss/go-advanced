package main

import "sync"

func main() {
	var mu sync.Mutex

	mu.Lock()
	go func() {
		println("你好, 世界")
		mu.Unlock()
	}()

	mu.Lock()
}
