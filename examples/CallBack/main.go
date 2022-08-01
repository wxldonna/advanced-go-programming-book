package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sum := 0

	go func() {
		for i := 0; i < 5; i++ {
			sum += i
			//	log.Printf("index %d sum := %d", i, sum)
			fmt.Printf("index %d sum := %d \n", i, sum)
		}
	}()

	go func() {
		for i := 5; i < 10; i++ {
			sum += i
			fmt.Printf("index %d sum := %d \n", i, sum)
		}

	}()
	ctrlc := make(chan os.Signal)
	signal.Notify(ctrlc, os.Interrupt, syscall.SIGTERM)

	<-ctrlc
	fmt.Printf("Ctrl+C pressed")

}
