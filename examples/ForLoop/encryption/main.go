package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	stringChan := make(chan string)
	tower1Chan := make(chan string)
	tower2Chan := make(chan string)

	offset := 3
	go tower1(stringChan, tower1Chan, offset)
	go tower2(stringChan, tower2Chan, offset)
	go func() {
		for {
			select {
			case msg := <-tower1Chan:
				fmt.Printf("\nControl Tower:Message from Tower 1 - %v", msg)
			case msg := <-tower2Chan:
				fmt.Printf("\nControl Tower Message from Tower 2 - %v", msg)
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()
	ctrlc := make(chan os.Signal)
	signal.Notify(ctrlc, os.Interrupt, syscall.SIGTERM)
	<-ctrlc
	fmt.Println("server is shut-down")

}
func tower1(s chan string, t1 chan string, offset int) {
	inputStream := bufio.NewReader(os.Stdin)
	fmt.Printf("Tower1:enter your message for tower2 \n")
	userInput, _ := inputStream.ReadString('\n')

	fmt.Printf("\nTower1: Original String :%s", userInput)

	var secretString string
	for _, c := range userInput {
		secretString += string(c + int32(offset))
	}
	fmt.Printf("\nTower 1: Encrypted String: %s", secretString)
	s <- secretString
	t1 <- "Msg send to Tower 2"
}

func tower2(s chan string, t2 chan string, offset int) {
	secretString := <-s
	var orgString string

	for _, c := range secretString {
		orgString += string(c - int32(offset))
	}
	fmt.Printf("\nTower2: Decrypted Message :%s", orgString)
	t2 <- "Message received from Tower 1"
}
