package main

func main() {
	done := make(chan int)

	go func() {
		println("Hello World")
		done <- 1
	}()

	<-done
}
