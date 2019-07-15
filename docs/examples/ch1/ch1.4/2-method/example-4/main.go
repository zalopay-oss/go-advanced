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

	// giá trị phương thức: ràng buộc với đối tượng f
	// func Close() error
	var Close = f.Close

	// giá trị phương thức: ràng buộc với đối tượng f
	// func Read (offset int64, data []byte) int
	var Read = f.Read

	// xử lý file
	Read(0, data)
	Close()
}
