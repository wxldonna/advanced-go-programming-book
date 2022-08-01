package main

import (
	"fmt"
	"log"
	"runtime/debug"
)

/*
func my_div(a, b int) int {
	return a / b
}

func handle_panic() {
	if r := recover(); r != nil {
		fmt.Printf("failed: %v\n", r)
	}
}

func example() {
	defer handle_panic()

	fmt.Printf("result: %v\n", my_div(0, 0))
}

func main() {
	go example()
	time.Sleep(1 * time.Second)
}

*/
func my_div(a, b int) (result int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic: %v\n", r)
			log.Printf("panic: %v\n", r)
			//log.Printf("panic stack: %s\n", debug.Stack())
			debug.PrintStack()
			result = 0
		}
	}()

	return a / b
}
func main() {
	my_div(1, 0)
}
