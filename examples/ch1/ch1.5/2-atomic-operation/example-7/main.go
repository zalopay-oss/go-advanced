package main

func main() {
	var a string
	var done bool
	
	func setup() {
		a = "hello world"
		done = true
	}
	
	func main() {
		go setup()
		for !done {}
		print(a)
	}
}
