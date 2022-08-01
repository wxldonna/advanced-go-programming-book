package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const SockAddr = "/tmp/echo.sock"

func echoServer(c net.Conn) {
	log.Printf("Client connected [%s]", c.RemoteAddr().Network())
	// get the source content from src reader
	var b []byte
	b = make([]byte, 1024)
	for {
		n, err := c.Read(b)

		if err != nil {
			if err == io.EOF {

				log.Printf("read finish \n")
				continue
			}
			break
		}
		log.Printf("read %d content from source", n)
		t := fmt.Sprintf("echo from server %s", string(b))
		c.Write([]byte(t))

	}

	c.Close()
}

func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("unix", SockAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()

	for {
		// Accept new connections, dispatching them to echoServer
		// in a goroutine.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		go echoServer(conn)
	}
}
