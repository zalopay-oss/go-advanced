package main

import "fmt"

func worker(queue chan int, worknumber int, done chan bool) {
	for j := range queue {
		fmt.Println("worker", worknumber, "finished job", j)
		done <- true
	}
}

func main() {

	// queue of jobs
	q := make(chan int)

	// done channel lấy ra kết quả của jobs
	done := make(chan bool)

	// số lượng worker trong pool
	numberOfWorkers := 4
	for i := 0; i < numberOfWorkers; i++ {
		go worker(q, i, done)
	}

	// đưa job vào queue
	numberOfJobs := 17
	for j := 0; j < numberOfJobs; j++ {
		go func(j int) {
			q <- j
		}(j)
	}

	// chờ nhận đủ kết quả
	for c := 0; c < numberOfJobs; c++ {
		<-done
	}
}
