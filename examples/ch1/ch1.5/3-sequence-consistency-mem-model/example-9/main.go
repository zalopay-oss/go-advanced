package main

import "sync"

func main() {
	var mu sync.Mutex

	mu.Lock()
	go func() {
		println("Hello World")
		mu.Unlock()
	}()

	mu.Lock()
}
