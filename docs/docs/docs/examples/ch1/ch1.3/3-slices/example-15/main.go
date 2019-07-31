package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func FindPhoneNumber(filename string) []byte {
	b, _ := ioutil.ReadFile(filename)
	return regexp.MustCompile("[0-9]+").Find(b)
}

func main() {
	b := FindPhoneNumber("dial.txt")
	fmt.Println(string(b))
}
