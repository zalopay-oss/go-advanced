package main

import (
	"fmt"
	"log"

	pb "../hello_pb"
)

func main() {
	client, err := pb.DialHelloService("tcp", "localhost:1234")
	// log ra lỗi nếu có
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// biến lưu kết quả từ lời gọi RPC
	var reply = &pb.String{}
	// thực thi lệnh gọi RPC
	err = client.Hello(&pb.String{Value: "World"}, reply)
	// log ra lỗi nếu có
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.Value)
}
