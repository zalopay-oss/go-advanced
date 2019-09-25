package main

import (
	"log"
	"net"
	"net/rpc"

	pb "../hello_pb"
)

// đối tượng RPC HelloService
type HelloService struct{}

// hiện thực lời gọi RPC
func (p *HelloService) Hello(request *pb.String, reply *pb.String) error {
	*reply = pb.String{Value: "Hello" + request.Value}
	return nil
}

// hàm main phía server
func main() {
	pb.RegisterHelloService(&rpc.Server{}, new(HelloService))
	// lắng nghe kết nối từ phía client
	listener, err := net.Listen("tcp", ":1234")
	// log ra lỗi nếu có (vd: trùng port, v,v..)
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	// vòng lặp tiếp nhận nhiều kết nối client
	for {
		// chấp nhận kết nối từ một client nào đó
		conn, err := listener.Accept()
		// in ra lỗi nếu có
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		// phục vụ kết nối trên một goroutine khác
		// để main thread tiếp tục vòng lặp accept client khác
		go rpc.ServeConn(conn)
	}
}
