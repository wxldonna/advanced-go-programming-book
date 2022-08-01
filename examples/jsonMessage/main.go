package main

import (
	"encoding/json"
	"log"
)

type Str1 struct {
	Com1 map[string]json.RawMessage `json:"com_1"`
}

func main() {
	st1 := &Str1{
		Com1: map[string]json.RawMessage{
			"abc": json.RawMessage("123"),
		},
	}

	str1, err := json.Marshal(st1)
	if err != nil {
		panic(err)
	}
	log.Printf("sting ===%s", str1)
}
