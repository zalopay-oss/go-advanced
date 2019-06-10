package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

func main() {
	// Mảng string
	var s1 = [2]string{"hello", "world"}
	var s2 = [...]string{"Hello!", "World"}
	var s3 = [...]string{1: "Hello", 0: "World"}

	fmt.Println(s1, s2, s3)
	// Mảng struct
	var line1 [2]image.Point
	var line2 = [...]image.Point{image.Point{X: 0, Y: 0}, image.Point{X: 1, Y: 1}}
	var line3 = [...]image.Point{{0, 0}, {1, 1}}
	fmt.Println(line1, line2, line3)
	// Mảng decoder của hình ảnh
	var decoder1 [2]func(io.Reader) (image.Image, error)
	var decoder2 = [...]func(io.Reader) (image.Image, error){
		png.Decode,
		jpeg.Decode,
	}
	fmt.Println(decoder1, decoder2)
	// Mảng interface{}
	var unknown1 [2]interface{}
	var unknown2 = [...]interface{}{123, "Hello!"}
	fmt.Println(unknown1, unknown2)
	// Mảng pipe
	var chanList = [2]chan int{}
	fmt.Println(chanList)
}
