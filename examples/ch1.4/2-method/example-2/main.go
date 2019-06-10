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
func (f *File) Close() error {
	fmt.Println("Close file")
	return nil
}

// đọc dữ liệu từ file
func (f *File) Read(offset int64, data []byte) int {
	fmt.Println("Read file")
	return 0
}

func main() {
	var data []byte

	// không phụ thuộc vào đối tượng file cụ thể
	// func CloseFile(f *File) error
	var CloseFile = (*File).Close

	// không phụ thuộc vào đối tượng file cụ thể
	// func ReadFile(f *File, offset int64, data []byte) int
	var ReadFile = (*File).Read

	// xử lý file
	f, _ := OpenFile("foo.dat")
	ReadFile(f, 0, data)
	CloseFile(f)
}
