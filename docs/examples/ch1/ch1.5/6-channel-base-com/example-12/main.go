package main

var done = make(chan bool)
var msg string

func aGoroutine() {
	msg = "Hello World"
	close(done)
}

func main() {
	go aGoroutine()
	<-done
	println(msg)
}
