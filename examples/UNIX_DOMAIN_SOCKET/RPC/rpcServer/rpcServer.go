package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

const SockAddr = "/tmp/rpc.sock"

type Greeter struct {
}

func (g Greeter) Greet(name *string, reply *string) error {
	*reply = "Hello, " + *name
	return nil
}

func (g Greeter) Nice(name string, reply *string) error {
	*reply = "nice, " + name
	return nil
}

func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	greeter := new(Greeter)
	rpc.Register(greeter)
	rpc.HandleHTTP()
	l, e := net.Listen("unix", SockAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	fmt.Println("Serving...")
	http.Serve(l, nil)
}
