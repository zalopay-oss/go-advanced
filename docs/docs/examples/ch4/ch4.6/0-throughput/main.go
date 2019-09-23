package main

import (
	"io"
	"log"
	"net/http"
)

func sayhello(wr http.ResponseWriter, r *http.Request) {
	wr.WriteHeader(200)
	io.WriteString(wr, "hello world")
}

func main() {
	http.HandleFunc("/", sayhello)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
