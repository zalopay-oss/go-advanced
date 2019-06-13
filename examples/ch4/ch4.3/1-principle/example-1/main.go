package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func (client *Client) Call(
	serviceMethod string, args interface{},
	reply interface{},
) error {
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}

func doClientWork(client *rpc.Client) {
	helloCall := client.Go("HelloService.Hello", "hello", new(string), nil)

	// do some thing

	helloCall = <-helloCall.Done
	if err := helloCall.Error; err != nil {
		log.Fatal(err)
	}

	args := helloCall.Args.(string)
	reply := helloCall.Reply.(string)
	fmt.Println(args, reply)
}

func (client *Client) Go(
	serviceMethod string, args interface{},
	reply interface{},
	done chan *Call,
) *Call {
	call := new(Call)
	call.ServiceMethod = serviceMethod
	call.Args = args
	call.Reply = reply
	call.Done = make(chan *Call, 10) // buffered.

	client.send(call)
	return call
}
