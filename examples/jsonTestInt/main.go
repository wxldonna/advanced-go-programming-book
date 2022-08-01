package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

func main() {
	type Parameters map[string]interface{}
	pars := map[string]interface{}{}
	rfcParamters := map[string]interface{}{
		"IV_QUERY_JSON":    "quert",
		"IV_METADATA_JSON": "metadata",
		"size":             100,
	}
	var network bytes.Buffer
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.
	enc.Encode(rfcParamters)
	dec.Decode(&pars)
	log.Println(pars)
	/*
		raw, err := json.Marshal(rfcParamters)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(raw, &pars)


	*/

}
