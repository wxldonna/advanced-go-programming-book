package main

import (
	"bytes"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

/*
func ExampleMarshal() {
	type Item struct {
		Foo string
	}

	b, err := msgpack.Marshal(&Item{Foo: "bar"})
	if err != nil {
		panic(err)
	}

	var item Item
	err = msgpack.Unmarshal(b, &item)
	if err != nil {
		panic(err)
	}
	fmt.Println(item.Foo)
	// Output: bar
}
*/

type Item struct {
	//_msgpack struct{} `msgpack:",as_array"`
	Foo string
	Bar string
}

func encoderStructure(buf *bytes.Buffer) {
	// create writer using buf
	enc := msgpack.NewEncoder(buf)

	enc.EncodeInt(1) //when you encode some value , you need decode them as the correct sequence
	// for the reader , you need decoder.DecodeInt()
	// write content to buf
	err := enc.Encode(&Item{Foo: "foo", Bar: "bar"})
	if err != nil {
		panic(err)
	}
}
func decodeStructureBuffer(buf *bytes.Buffer) {
	fmt.Printf("byte %b \n", buf)

	// create reader using the buf
	dec := msgpack.NewDecoder(buf)
	//if you add the int value to the encoder,
	// decode it
	len, _ := dec.DecodeInt()
	fmt.Printf("the length is %d \n", len)
	deItem := Item{}
	dec.Decode(&deItem)
	fmt.Printf("decode item %v \n", deItem)
}

func decodeStructure(buf []byte) {
	fmt.Printf("byte %b \n", buf)
	// create reader using the buf
	dec := msgpack.NewDecoder(bytes.NewReader(buf))
	//if you add the int value to the encoder,
	// decode it
	len, _ := dec.DecodeInt()
	fmt.Printf("the length is %d \n", len)
	deItem := Item{}
	dec.Decode(&deItem)
	fmt.Printf("decode item %v \n", deItem)
}

type VScalar = interface{}

type VStruct = []VScalar

type VTable = []VStruct

func createTable() VTable {
	result := make(VTable, 0)
	for i := 0; i < 10; i++ {
		reStr := make(VStruct, 0)
		for i := 0; i < 2; i++ {
			reStr = append(reStr, "10")
		}
		result = append(result, reStr)
	}
	return result
}

func encoderTable(writer *bytes.Buffer, table VTable) {
	// create writer using bytes Buffer
	enc := msgpack.NewEncoder(writer)
	enc.EncodeArrayLen(len(table))
	// write table to the writer Buffer and encode it
	if err := enc.Encode(table); err != nil {
		panic(err)
	}

}

func decodeTable(reader *bytes.Buffer) {

	dec := msgpack.NewDecoder(reader)
	if len, err := dec.DecodeArrayLen(); err != nil {
		panic(err)
	} else {
		fmt.Printf("array length is %d \n", len)
	}

	deTable := VTable{}
	dec.Decode(&deTable)
	/*
		if IterTable, err := dec.DecodeInterface(); err != nil {
			fmt.Printf("decode error is %v", err)
		} else {
			fmt.Printf("interfance table is \n %v ", IterTable)
		}
	*/
	fmt.Printf("decode table %v \n", deTable)
}
func main() {
	//ExampleMarshal()

	var buf bytes.Buffer
	encoderStructure(&buf)
	fmt.Printf("buf %b", buf.Bytes())
	//decodeStructureBuffer(&buf)
	//decodeStructure(buf.Bytes())

	/*
		var buf bytes.Buffer
		result := createTable()
		encoderTable(&buf, result)
		decodeTable(&buf)
	*/
}
