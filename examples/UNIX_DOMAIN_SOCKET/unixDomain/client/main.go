package main

import (
	"log"
	"net"
)

func main() {
	c, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		panic(err)
	}
	hello := "hello server"
	msgRe := make([]byte, 1024)

	for i := 0; i < 10; i++ {
		n, err := c.Write([]byte(hello))
		if err != nil {
			log.Printf("dial failed with error %v", err)
		}
		log.Printf("write the lenth %d", n)

		serverN, err := c.Read(msgRe)
		if err != nil {
			log.Printf("dial failed with error %v", err)
		}
		log.Printf("wread the lenth %d \n", serverN)
		log.Printf("wread the msg %s \n", string(msgRe))
	}

}
