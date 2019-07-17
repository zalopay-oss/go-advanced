package main

import (
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

func main() {
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	l := zk.NewLock(c, "/lock", zk.WorldACL(zk.PermAll))
	err = l.Lock()
	if err != nil {
		panic(err)
	}
	println("lock succ, do your business logic")

	time.Sleep(time.Second * 10)

	// do some thing
	l.Unlock()
	println("unlock succ, finish business logic")
}