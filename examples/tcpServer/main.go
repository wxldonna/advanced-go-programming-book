package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/vmihailenco/msgpack/v5"
)

type Student struct {
	Name string
	Sex  string
	Age  int
}

func main() {
	// Listen for incoming connections.
	//	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	l, err := net.Listen("tcp4", ":0") // #nosec G102
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	address := l.Addr().(*net.TCPAddr)

	fmt.Printf("Listening on %s \n ", address.IP.String()+":"+strconv.Itoa(address.Port))
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
		//go handleRequestRead(conn)
	}
}

// Handles incoming requests.
// Client send structure data and Unmarshal it back in server
func handleRequest(conn net.Conn) {

	//client echo -n "test client" | nc localhost 35215    the content is 11, so make([]byte,11),
	//if you make shorter, you can only receove part of message and couldn't write message to client
	buf := make([]byte, 11)

	if err := binary.Read(conn, binary.LittleEndian, &buf); err != nil {
		fmt.Printf("read connection data is %v", err)
	} else {
		fmt.Printf("received message %s \n", string(buf))
	}

	//if err := msgpack.Unmarshal(headerBytes, &msg.Header); err != nil {

	// Send a response back to person contacting us.
	conn.Write([]byte("Message received. \n"))
	// Close the connection when you're done with it.
	conn.Close()
	fmt.Printf("write finished")
}

// can we use Read to unmarshal the byte to get the correct structure?
func handleRequestRead(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	//
	targetByte := buf[:reqLen]
	dec := msgpack.NewDecoder(bytes.NewReader(targetByte))
	// encode and decode must use the global type
	student := Student{}
	if err := dec.Decode(&student); err != nil {
		fmt.Println(err)
	}
	//
	fmt.Printf("received %v and the length is %d \n", student, reqLen)
	// Send a response back to person contacting us.
	conn.Write([]byte("Message received. \n"))
	// Close the connection when you're done with it.
	conn.Close()
	fmt.Printf("write finished %d \n", reqLen)
}
