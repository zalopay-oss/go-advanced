package main

//#include <stdio.h>
import "C"
import (
	"fmt"
	"time"
)



func main() {
	C.puts(
		C.CString("Hello World\n"),
	)
	fmt.Println("Test")
	time.Sleep(2 * time.Second)

}
