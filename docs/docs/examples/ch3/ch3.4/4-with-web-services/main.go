package main

import (
	"crypto/tls"
	fmt "fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var port = ":5000"

type myGrpcServer struct{}

func (s *myGrpcServer) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "[rpc] Hello " + in.Name}, nil
}

func main() {
	go startServer()
	time.Sleep(time.Second)

	doClientWork()
	doClient2()
}

func startServer() {
	creds, err := credentials.NewServerTLSFromFile("tls-config/server.crt", "tls-config/server.key")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	RegisterGreeterServer(grpcServer, new(myGrpcServer))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "hello")
	})

	http.ListenAndServeTLS(port, "tls-config/server.crt", "tls-config/server.key", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	}))
}

func doClient2() {
	log.SetFlags(log.Lshortfile)
	cert, err := tls.LoadX509KeyPair("tls-config/server.crt", "server.grpc.io")

	conf := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	conn, err := tls.Dial("tcp", "localhost:5000", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}

	println(string(buf[:n]))
}

func doClientWork() {
	creds, err := credentials.NewClientTLSFromFile("tls-config/server.crt", "server.grpc.io")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := NewGreeterClient(conn)

	r, err := c.SayHello(context.Background(), &HelloRequest{Name: "gopher"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("[client]: %s", r.Message)
}
