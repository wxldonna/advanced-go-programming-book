package main

import (
	"encoding/json"
	"log"
	"time"
)

func main() {
	type stu struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		Sex  string `json:"sex"`
	}

	st1 := stu{
		Name: "wxl",
		Age:  10,
		Sex:  "male",
	}
	stB, err := json.Marshal(&st1)
	if err != nil {
		log.Panic(err)
	}
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			timer.Reset(time.Second)
			log.Printf("st1 :%s", string(stB))

		default:
			//log.Printf("st1 :%s", string(stB))
		}

	}

}
