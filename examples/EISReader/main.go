package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

type Message struct {
	Id              string    `json:"id"`
	Source          string    `json:"source"`
	Specversion     string    `json:"specversion"`
	Type            string    `json:"type"`
	Subject         string    `json:"subject"`
	Time            time.Time `json:"time"`
	Datacontenttype string    `json:"datacontenttype"`
	Data            struct {
		InitialLoadId      string `json:"initial-load-id"`
		Type               string `json:"type"`
		DestinationPrefix  string `json:"destination-prefix"`
		DestinationAuthUri string `json:"destination-auth-uri"`
		RowCount           int    `json:"row-count"`
		Status             string `json:"status"`
		StatusMessage      string `json:"status-message"`
		RandomContent      string `json:"random_content"`
	} `json:"data"`
}

/*
{
	"id": "ile007_12345",
	"source": "/us-west-2/concur/pdw3324rfwwt",
	"specversion": "1.0",
	"type": "sap.ism.InitialLoad.Created.v1",
	"subject": "asdf12345",
	"time": "2018-04-05T17:31:00Z",
	"datacontenttype": "application/json",
	"data": {
		"initial-load-id": "1111",
		"type": "concur.ExpenseReport.Created.v1",
		"destination-prefix": "/lob=concur/bobj=Employee/tenant=us-west-2.concur.pdw3324rfwwt",
		"destination-auth-uri": "https://api.business.network.sap/ism/v0/auth",
		"row-count": 100240,
		"status": "Submitted",
		"status-message": ""
	}
}
*/
func main() {

	rawMsg := GenerateMessage(1)
	log.Printf("message is %s \n", string(rawMsg))
	dec := gob.NewDecoder(bytes.NewReader(rawMsg))

	msg := Message{}

	if err := dec.Decode(&msg); err != nil {
		panic(err)
	}
	log.Printf("msg %v \n", msg)

	type EISMessage []Message
	msgs := EISMessage{}
	msgs = append(msgs, msg)
	msgByte, err := json.Marshal(&msgs)
	if err != nil {
		panic(err)
	}
	log.Printf("jsonMessage %s \n", string(msgByte))
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GenerateMessage(i int) []byte {
	msg := Message{
		Id:              "ile007_12345",
		Source:          "/us-west-2/concur/pdw3324rfwwt",
		Specversion:     "1.0",
		Type:            "sap.ism.InitialLoad.Created.v1",
		Subject:         "asdf12345",
		Time:            time.Now(),
		Datacontenttype: "application/json",
		Data: struct {
			InitialLoadId      string `json:"initial-load-id"`
			Type               string `json:"type"`
			DestinationPrefix  string `json:"destination-prefix"`
			DestinationAuthUri string `json:"destination-auth-uri"`
			RowCount           int    `json:"row-count"`
			Status             string `json:"status"`
			StatusMessage      string `json:"status-message"`
			RandomContent      string `json:"random_content"`
		}{
			InitialLoadId:      GetRandomString(i),
			Type:               "sap.ism.InitialLoad.Created.v1",
			DestinationPrefix:  "/lob=concur/bobj=Employee/tenant=us-west-2.concur.pdw3324rfwwt",
			DestinationAuthUri: "https://api.business.network.sap/ism/v0/auth",
			RowCount:           100240,
			Status:             "Submitted",
			RandomContent:      GetRandomString(i),
		},
	}

	buf := new(bytes.Buffer)
	end := gob.NewEncoder(buf)
	if err := end.Encode(&msg); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
