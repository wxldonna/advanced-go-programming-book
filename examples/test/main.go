package main

import (
	"fmt"
	"strings"
)

func main() {
	/*
		t, err := time.ParseInLocation(time.RFC3339, "2016-01-02T15:04:05+09:00", time.UTC) //time.Parse(time.RFC3339, "2016-01-02T15:04:05+07:00")
		//t, err := time.Parse(time.RFC3339, "t, _ = time.Parse(time.RFC3339, \"2006-01-02T15:04:05+07:00\")")
		if err != nil {
			log.Printf("%v", err)
		}
		log.Printf("original time is %v \n ", t)
		utcNewT := t.UTC()
		log.Printf("old time is %v \n ", utcNewT)

		newTime := utcNewT.Format("2006-01-02T15:04:05.999999999")
		//newTime := t.Format(time.RFC3339Nano)
		//ti := fmt.Sprintf("%v", t)
		log.Printf("time is %s ", newTime)

		test := map[string]string{
			"test1": "qw",
			"test2": "qwer",
		}
		c := test["test1"]
		c = "hello"
		log.Printf("result is %v \n", c)

		log.Printf("test is %v \n", test)


	*/
	fmt.Print(strings.Trim("/Hello, Gophers/qww/", "/"))

}
