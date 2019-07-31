package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var logger = log.New(os.Stdout, "", 0)

func helloHandler(wr http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	wr.Write([]byte("call helloHandler"))
	timeElapsed := time.Since(timeStart)
	logger.Println(timeElapsed)
}

func showInfoHandler(wr http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	wr.Write([]byte("call showInfoHandler"))
	timeElapsed := time.Since(timeStart)
	logger.Println(timeElapsed)
}

func showEmailHandler(wr http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	wr.Write([]byte("call showEmailHandler"))
	timeElapsed := time.Since(timeStart)
	logger.Println(timeElapsed)
}

func showFriendsHandler(wr http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	wr.Write([]byte("your friends is tom and alex"))
	timeElapsed := time.Since(timeStart)
	logger.Println(timeElapsed)
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/info/show", showInfoHandler)
	http.HandleFunc("/email/show", showEmailHandler)
	http.HandleFunc("/friends/show", showFriendsHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
