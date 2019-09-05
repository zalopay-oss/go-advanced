package main

import (
	"fmt"
)

func searchByBaidu(inp string)(out string){
	return (inp+ " from Baidu")
}

func searchByBing(inp string)(out string){
	return (inp+ " from Bing")
	 
}

func searchByGoogle(inp string)(out string){
	return (inp+ " from Google")
}

func main() {
    ch := make(chan string, 32)
	
	go func() {
		ch <- searchByBaidu("golang")
	}()
    go func() {
        ch <- searchByBing("golang")
    }()
    go func() {
        ch <- searchByGoogle("golang")
    }()

    fmt.Println(<-ch)
}