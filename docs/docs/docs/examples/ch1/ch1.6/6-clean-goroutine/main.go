package main

import (
	"fmt"
	"time"
)

func worker(queue chan int, worknumber int, done, ks chan bool) {
	for true {
		// dùng select để chờ cùng lúc trên cả 2 channel
		select {
		// xử lý job trong channel queue
		case k := <-queue:
			fmt.Println("doing work!", k, "worknumber", worknumber)
			done <- true

		// nếu nhận được kill signal thì return
		case <-ks:
			fmt.Println("worker halted, number", worknumber)
			return
		}
	}
}

func main() {
	// channel để terminate các worker
	killsignal := make(chan bool)

	// queue các jobs
	q := make(chan int)
	// done channel nhận vào kết quả của các job
	done := make(chan bool)

	// số lượng worker trong pool
	numberOfWorkers := 4
	for i := 0; i < numberOfWorkers; i++ {
		go worker(q, i, done, killsignal)
	}

	// đưa job vào queue
	numberOfJobs := 17
	for j := 0; j < numberOfJobs; j++ {
		go func(j int) {
			q <- j
		}(j)
	}

	// chờ để nhận đủ kết quả
	for c := 0; c < numberOfJobs; c++ {
		<-done
	}

	// dọn dẹp các worker
	close(killsignal)
	time.Sleep(2 * time.Second)
}
