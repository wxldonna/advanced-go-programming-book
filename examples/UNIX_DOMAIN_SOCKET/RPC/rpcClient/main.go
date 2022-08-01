package main

import (
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("unix", "/tmp/rpc.sock")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	/*
		name := "Joe"
		var reply string
		err = client.Call("Greeter.Greet", &name, &reply)
		if err != nil {
			log.Fatal("greeter error:", err)
		}
		fmt.Printf("Got '%s'\n", reply)

		err = client.Call("Greeter.Nice", name, &reply)
		if err != nil {
			log.Fatal("Nice error:", err)
		}
		fmt.Printf("Got '%s'\n", reply)

	*/
	AsynchronCall(client)
}

func AsynchronCall(client *rpc.Client) {
	// Asynchronous call
	name := "Joe"
	var reply string
	replyCall := client.Go("Greeter.Greet", &name, &reply, nil)
	// check errors, print, etc
	log.Printf("async call finished with %v", <-replyCall.Done)
	log.Printf("asysn call result %s", reply)

}
