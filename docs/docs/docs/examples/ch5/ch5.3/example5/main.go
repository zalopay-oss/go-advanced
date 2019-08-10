package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var logg = log.New(os.Stdout, "INFO ", 3)

type middleware func(http.Handler) http.Handler

type Router struct {
	middlewareChain []middleware
	mux             map[string]http.Handler
}

func NewRouter() *Router {
	return &Router{nil, make(map[string]http.Handler)}
}

func (r *Router) Use(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *Router) Add(route string, h http.Handler) {
	var mergedHandler = h

	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}

	r.mux[route] = mergedHandler
}

// Implement the ServeHTTP method on our Router type
func (r *Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if handler := r.mux[req.URL.Path]; handler != nil {
		handler.ServeHTTP(resp, req)
	} else {
		http.Error(resp, "Bad Request", 400) // Or Redirect?
	}
}

// Implement our logger middleware
func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(wr, r)
		logg.Println("call middleware logger")
	})
}

// Implement our timeout middleware
func timeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		next.ServeHTTP(wr, r)
		timeElapsed := time.Since(timeStart)
		logg.Println("call middleare timeout ", timeElapsed)
	})
}

// Implement our helloHandler
func helloHandler(wr http.ResponseWriter, r *http.Request) {
	wr.Write([]byte("call helloHandler"))
}

func main() {
	r := NewRouter()
	r.Use(logger)
	r.Use(timeout)
	r.Add("/", http.HandlerFunc(helloHandler))
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}
