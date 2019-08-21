package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func FindPhoneNumber(filename string) []byte {
	b, _ := ioutil.ReadFile(filename)
	b = regexp.MustCompile("[0-9]+").Find(b)
	fmt.Println(b)
	return append([]byte{}, b...)
}

func main() {
	b := FindPhoneNumber("dial.txt")
	fmt.Println(string(b))
}
