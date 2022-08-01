package main

import (
	"encoding/json"
	"log"
	"sort"
)

func main() {
	/*
		type QueryID struct {
			ID string `json:"id"`
		}
		id := QueryID{
			"1234567890string",
		}

	*/
	type Message struct {
		Type string `json:"type"`
		Text string `json:"text" ,json:"message"`
	}
	type Result struct {
		Columns  []string        `json:"columns"`
		Result   [][]interface{} `json:"result"`
		Messages []Message       `json:"messages"`
	}
	/*
		rvResult := Result{
			Columns: []string{
				"ke1",
				"key2",
				"field3",
				"field4",
			},
			Result: [][]interface{}{
				[]interface{}{
					"800",
					"!KI",
					123,
					456,
				},
				[]interface{}{
					"800",
					"!KI",
					123,
					456,
				},
			},
			Messages: []Message{
				{
					Type: "E",
					Text: "XXXXX",
				},
			},
		}
	*/

	type FederationVersion struct {
		SupportVersion []map[string]interface{} `json:"supportedVersions"`
	}

	//data := " {\"columns\":[\"ke1\",\"key2\",\"field3\",\"field4\"],\"result\":[[\"800\",\"!KI\",123,456],[\"800\",\"!KI\",123,456]],\"messages\":[{\"type\":\"E\",\"message\":\"XXXXX\"}]}\n"
	data := "{\n  \"supportedVersions\": [\n    {\n      \"version\": \"v1\",\n      \"feature1\": true,\n      \"feature2\": false\n    },\n    {\n      \"version\": \"v2\",\n      \"feature1\": true,\n      \"feature2\": false\n    }\n  ],\n  \"messages\": [\n    {\n      \"type\": \"E\",\n      \"text\": \"xxxx\"\n    }\n  ]\n}"
	/*
		idByte, _ := json.Marshal(rvResult)
		log.Printf("%s", fmt.Sprintf("%s", idByte))

	*/
	var result FederationVersion
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Panic(err)
	}
	//log.Printf("result %v", result)
	sort.Slice(result.SupportVersion, func(i, j int) bool {
		versionI, _ := result.SupportVersion[i]["version"].(string)
		versionJ, _ := result.SupportVersion[j]["version"].(string)
		return versionI > versionJ
		//return result.SupportVersions[i]["version"] > result.SupportVersions[j]["version"]
	})
	log.Printf("%v", result.SupportVersion[0])
}
