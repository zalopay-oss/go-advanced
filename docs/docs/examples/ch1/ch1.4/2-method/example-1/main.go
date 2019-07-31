package main

import (
	"fmt"
)

// đối tượng File
type File struct {
	fd int
}

// mở file
func OpenFile(name string) (f *File, err error) {
	fmt.Println("Opening file ", name)
	return nil, nil
}

// đóng file
func CloseFile(f *File) error {
	fmt.Println("Close file")
	return nil
}

// đọc dữ liệu từ file
func ReadFile(f *File, offset int64, data []byte) int {
	fmt.Println("Read file")
	return 0
}

func main() {
	myFile := File{fd: 4}
	var data []byte

	OpenFile("file xx")
	
	ReadFile(&myFile, 0, data)
	
	CloseFile(&myFile)
}
