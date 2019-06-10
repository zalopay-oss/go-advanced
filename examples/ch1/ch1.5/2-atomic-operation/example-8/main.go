package main

func main() {
	done := make(chan int)

	go func() {
		println("你好, 世界")
		done <- 1
	}()

	<-done
}
