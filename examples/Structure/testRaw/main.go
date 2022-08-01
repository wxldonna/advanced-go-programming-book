package main

import (
	"encoding/json"
	"log"
)

/*
type RequestBody struct {
	Metadata []interface{} `json:"metadata"`
	Query    interface{}   `json:"query"`
}

*/
type RequestBody struct {
	Metadata json.RawMessage `json:"metadata"`
	Query    json.RawMessage `json:"query"`
}

func main() {
	str := `{
  "metadata": [
    {
      "name": "carrid",
      "type": "NVARCHAR",
      "length": 5
    }
  ],
  "query": "mockQuery"
}`
	request := RequestBody{}
	err := json.Unmarshal([]byte(str), &request)
	if err != nil {
		panic(err)
	}
	log.Printf("hello %s", request)

	/*
		requestByte, err := json.Marshal(&request)
		if err != nil {
			panic(err)
		}
		log.Printf("byte %s", string(requestByte))

	*/
}
