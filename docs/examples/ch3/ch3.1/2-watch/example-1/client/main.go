package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

func doClientWork(client *rpc.Client) {
	go func() {
		var keyChanged string
		err := client.Call("KVStoreService.Watch", 30, &keyChanged)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("watch:", keyChanged)
	}()

	err := client.Call(
		"KVStoreService.Set", [2]string{"abc", "value 1"},
		new(struct{}),
	)
	err = client.Call(
		"KVStoreService.Set", [2]string{"abc", "another value"},
		new(struct{}),
	)

	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	doClientWork(client)
}
