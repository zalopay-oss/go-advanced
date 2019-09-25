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

	// mở đối tượng file
	f, _ := OpenFile("foo.dat")

	// liên kết với đối tượng f
	// func Close() error
	var Close = func() error {
		return (*File).Close(f)
	}

	// liên kết với đối tượng f
	// func Read (offset int64, data []byte) int
	var Read = func (offset int64, data []byte) int {
		return (*File).Read(f, offset, data)
	}

	// xử lý file
	Read(0, data)
	Close()
}
