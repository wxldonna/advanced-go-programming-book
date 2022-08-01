package main

import (
	"log"
	"math/rand"
)

type linkedList struct {
	value   int
	pointer *linkedList
}

func main() {
	list := createlist(10, nil)

	for ; list != nil; list = list.pointer {
		log.Printf("%v\n", list.value)
	}
}

func createlist(num int, temp *linkedList) *linkedList {
	if temp == nil {
		temp = &linkedList{rand.Intn(100), nil}
		num--
	}
	tempList := temp
	for i := 0; i < num; i++ {
		t := &linkedList{rand.Intn(100), nil}
		tempList.pointer = t
		tempList = tempList.pointer
	}

	return temp
}
