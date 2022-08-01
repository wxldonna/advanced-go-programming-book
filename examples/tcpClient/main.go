package main

import (
	"bytes"
	"fmt"
	"net"
	"os"

	"github.com/vmihailenco/msgpack/v5"
)

type Student struct {
	Name string
	Sex  string
	Age  int
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "please specfify the server host \n ")
	}
	hostIP := args[1]
	hostPort := args[2]
	conn, err := net.Dial("tcp4", hostIP+":"+hostPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "the error is %w", err)
	}
	//conn.Write([]byte("hello client \n"))
	// send the structure
	student1 := Student{
		Name: "xiaoliangwang",
		Sex:  "male",
		Age:  20,
	}
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	err = enc.Encode(&student1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode error is %w", err)
	}
	structMessage := buf.Bytes()
	fmt.Fprintf(os.Stdout, "struct message is %b", structMessage)
	conn.Write(structMessage)
	go func() {
		clientMessage := make([]byte, 100)
		len, err := conn.Read(clientMessage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read message is %w", err)
		}
		fmt.Fprintf(os.Stdout, "receive message %s and the length is %d", string(clientMessage), len)
	}()
	go func() {
		//conn.Write([]byte("hello world \n"))
	}()
	stop := make(chan chan struct{})
	<-stop
}
