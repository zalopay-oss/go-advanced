package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// KVStoreService store key-value db
type KVStoreService struct {
	m      map[string]string
	filter map[string]func(key string)
	mu     sync.Mutex
}

// NewKVStoreService construct new obj
func NewKVStoreService() *KVStoreService {
	return &KVStoreService{
		m:      make(map[string]string),
		filter: make(map[string]func(key string)),
	}
}

// Get api
func (p *KVStoreService) Get(key string, value *string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if v, ok := p.m[key]; ok {
		*value = v
		return nil
	}

	return fmt.Errorf("not found")
}

// Set api
func (p *KVStoreService) Set(kv [2]string, reply *struct{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	key, value := kv[0], kv[1]

	if p.m != nil {
		
		if oldValue := p.m[key]; oldValue != value {
			for _, fn := range p.filter {
				fn(key)
			}
		}
	}
	
	p.m[key] = value
	fmt.Println("m[key]", p.m)
	return nil
}

// Watch to monitor
func (p *KVStoreService) Watch(timeoutSecond int, keyChanged *string) error {
	id := fmt.Sprintf("watch-%s-%03d", time.Now(), rand.Intn(999))
	ch := make(chan string, 10)

	p.mu.Lock()
	p.filter[id] = func(key string) { ch <- key }
	p.mu.Unlock()

	select {
	case <-time.After(time.Duration(timeoutSecond) * time.Second):
		return fmt.Errorf("timeout")
	case key := <-ch:
		*keyChanged = key
		return nil
	}
}

// ServeKVStoreService create conn
func ServeKVStoreService(conn net.Conn) {
	p := rpc.NewServer()
	p.Register(NewKVStoreService())
	p.ServeConn(conn)
}

func main() {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	fmt.Println("Server is running ", "localhost:1234")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go ServeKVStoreService(conn)
	}
}
