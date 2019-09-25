package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
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
	//lines1 := []int{1, 2, 3, 4, 5, 6, 6}
	arr := make([]int, 10000)
	for i := 0; i< 10000; i++ {
		arr[i] = rand.Intn(100000)
	}
	

	if err := writeLines(arr, "foo.out.txt"); err != nil {
		log.Fatalf("writeLines: %s", err)
	}

	//lines, err := readLines("foo.out.txt")
	// if err != nil {
	// 	log.Fatalf("readLines: %s", err)
	// }
	// for i, line := range lines {
	// 	fmt.Println(i, line)
	// }

	
}
