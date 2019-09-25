package main

//extern int go_qsort_compare(void* a, void* b);
import "C"

import (
    qsort "./qsort"
    "bufio"
	"log"
    "os"
    "unsafe"
    //"math/rand"
	"strconv"
)


// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// var lines []string
	var res []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// lines = append(lines, scanner.Text())
		x, _ := strconv.Atoi(scanner.Text())
		res = append(res, x)
	}
	return res, scanner.Err()
}

//export go_qsort_compare
func go_qsort_compare(a, b unsafe.Pointer) C.int {
    pa, pb := (*C.int)(a), (*C.int)(b)
    return C.int(*pa - *pb)
}

func main() {
    
	values, err := readLines("foo.out.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

    qsort.Sort(unsafe.Pointer(&values[0]),
        len(values), int(unsafe.Sizeof(values[0])),
        qsort.CompareFunc(C.go_qsort_compare),
    )
    //fmt.Println(values)
}