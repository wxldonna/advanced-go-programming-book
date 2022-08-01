package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"log"
)

func main() {
	/*
			// store to byte array
			strs := []string{"foo", "bar"}
			buf := &bytes.Buffer{}
			gob.NewEncoder(buf).Encode(strs)
			bs := buf.Bytes()
			fmt.Printf("%q", bs)

			// Decode it back
			strs2 := []string{}
			gob.NewDecoder(buf).Decode(&strs2)
			fmt.Printf("%v", strs2)
		}

	*/

	var students []string
	students = make([]string, 0)
	students = append(students, "xiaoliang")
	students = append(students, "xiaoming")
	writer := bytes.Buffer{}
	encoder := gob.NewEncoder(&writer)
	encoder.Encode(students)

	s := []string{}

	scb := writer.Bytes()
	remainingInput := scb
	for len(remainingInput) > 0 {
		nextStart, _, err := bufio.ScanLines(remainingInput, true)
		if err != nil {
			panic(err)
		}

		//log.Printf("line %s \n", string(line))
		reader := bytes.Buffer{}
		reader.Write(remainingInput)
		gob.NewDecoder(&reader).Decode(&s)
		log.Printf("s %s \n", s)
		log.Printf("nextStart is %d \n", nextStart)
		remainingInput = remainingInput[nextStart:]
	}
}
