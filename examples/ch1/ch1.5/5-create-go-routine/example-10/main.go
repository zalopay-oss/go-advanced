package main

var a string

func f() {
	print(a)
}

func hello() {
	a = "hello world"
	go f()
}

func main() {

}
