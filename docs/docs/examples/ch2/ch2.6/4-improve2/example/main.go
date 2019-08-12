package main 

import (
	"fmt"
    qsort "./qsort"
    "bufio"
	"log"
    "os"
    "math/rand"
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

// writeLines writes the lines to the given file.
func writeLines(lines []int, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, token := range lines {
		fmt.Fprintln(w, token)
	}
	return w.Flush()
}

func main() {
	arr := make([]int, 100000)
	for i := 0; i< 100000; i++ {
		arr[i] = rand.Intn(1000000)
	}
	

	if err := writeLines(arr, "foo.out.txt"); err != nil {
		log.Fatalf("writeLines: %s", err)
	}

	values, err := readLines("foo.out.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	// for i, line := range lines {
	// 	fmt.Println(i, line)
    // }

    qsort.Slice(values, func(i, j int) bool {
        return values[i] < values[j]
    })

    //fmt.Println(values)
}