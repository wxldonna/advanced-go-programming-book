package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	findRandomNumber(rand.Intn(100))
}

func findRandomNumber(randomNumber int) {
	count := 1
	numberfounded := false
	for {
		number := rand.Intn(10000)
		if number == randomNumber {
			numberfounded = true
			break
		}
		count++
	}
	if numberfounded {
		log.Printf("number %v found after %v attemp(s)", randomNumber, count)
	}
}
