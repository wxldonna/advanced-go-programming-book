package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	/*
		a := [1 << 3]byte{1, 2, 3, 4}
		fmt.Printf("a length is %d, the cap is %d \n", len(a), cap(a))

		b := (*[1 << 31]byte)(unsafe.Pointer(&a[0]))[:len(a):len(a)]
		fmt.Printf("array a is %v \n", a)
		fmt.Printf("slice a is %v \n", b)

		c := []byte{1, 2}
		fmt.Printf("[]byte c is %v \n", c)

		d := (*[1 << 2]byte)(unsafe.Pointer(&c[0]))
		fmt.Printf("[]byte c is %v \n", d)
	*/

	/*
			type RemoteObject struct {
				connection       interface{} `json:"connection"`
				name             string      `json:"name"`
				remoteObjectType string      `json:"remoteObjectType"`
				qualifiedName    string      `json:"qualifiedName"`
				description      string      `json:"description"`
			}

			str := `{
		  "connection": {
		    "id": "",
		    "type": "ABAP"
		  },
		  "description": "ZXL_ALL_TYPES",
		  "name": "ZXL_ALL_TYPES01",
		  "qualifiedName": "/TABLES/TMP/ZXL_ALL_TYPES01",
		  "remoteObjectType": ""
		}`
			obj := RemoteObject{}
			err := json.Unmarshal([]byte(str), &obj)
			if err != nil {
				fmt.Printf("err: %v", err)
			}
			fmt.Printf("onj:%v", obj)

			obj.description = "ZXL_ALL_TYPES"
			obj.name = "ZXL_ALL_TYPES01"
			obj.qualifiedName = "/TABLES/TMP/ZXL_ALL_TYPES01"
			obj.remoteObjectType = ""
			objB, _ := json.Marshal(&obj)
			fmt.Printf("objB:%s", string(objB))
			obj2 := RemoteObject{}
			json.Unmarshal(objB, &obj2)
			fmt.Printf("obj2:%v", obj2)
	*/
	type response2 struct {
		Page   int      `json:"page"`
		Fruits []string `json:"fruits"`
	}
	type RemoteObject struct {
		//connection       interface{} `json:"connection"`
		Name string `json:"name"`
		//remoteObjectType string `json:"remoteObjectType"`
		//qualifiedName    string `json:"qualifiedName"`
		//description      string `json:"description"`
	}

	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	res := response2{}
	json.Unmarshal([]byte(str), &res)
	fmt.Printf("res:%v \n", res)

	str2 := `{"name": "ZXL_ALL_TYPES01"}`
	res2 := RemoteObject{}
	json.Unmarshal([]byte(str2), &res2)
	fmt.Printf("res2:%v", res2)

	var obj interface{}

	obj = " {\n    \"connection\": {\n      \"id\": \"\",\n      \"type\": \"ABAP\"\n    },\n    \"name\": \"ZXL_ALL_TYPES01\",\n    \"remoteObjectType\": \"\",\n    \"qualifiedName\": \"/TABLES/TMP/ZXL_ALL_TYPES01\",\n    \"description\": \"ZXL_ALL_TYPES\"\n  }"
	type Person struct {
		FirstName string `json:"FirstName"`
		LastName  string `json:"LastName"`
	}
	objB, err := json.Marshal(&obj)
	p := Person{}
	err = json.Unmarshal([]byte(objB), &p)
	fmt.Printf("name:%v \n", p)
	fmt.Printf("err:%v \n", err)
	/*
		p := Person{
			"xiao", "Wang",
		}


		obj = p
	*/
	/*
		stu := obj.(Person)
		stuB, err := json.Marshal(&stu)
		fmt.Printf("name:%s \n", stu.FirstName)
		fmt.Printf("stuB:%s and err %v", string(stuB), err)

	*/

}
