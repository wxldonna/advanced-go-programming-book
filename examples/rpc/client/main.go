package main

import (
	"fmt"
	"log"
	"net/rpc"

	"chai2010.cn/gobook/examples/rpc/server"
)

const serverAddress = "localhost"

func main() {
	client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
	if err != nil {
		log.Fatalf("dialing: %v", err)
	}
	/*
		// Synchronous call
		args := &server.Args{7, 8}
		var reply int
		err = client.Call("Arith.Multiply", args, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		fmt.Printf("Arith: %d*%d=%d \n", args.A, args.B, reply)

		dividResult := server.Quotient{}
		err = client.Call("Arith.Divide", args, &dividResult)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		fmt.Printf("Divide: %d/%d=%v \n", args.A, args.B, dividResult)

	*/
	call, quan := AsynchronCall(client)
	fmt.Printf("AsynchronCall call is finished %+v\n", <-call.Done)
	fmt.Printf("result :%v\n", quan)

}

func AsynchronCall(client *rpc.Client) (*rpc.Call, *server.Quotient) {
	// Asynchronous call
	quotient := new(server.Quotient)
	args := &server.Args{7, 8}
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	//replyCall := <-divCall.Done	// will be equal to divCall
	return divCall, quotient
	// check errors, print, etc
}
