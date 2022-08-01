package main

import (
	"fmt"
	"strings"
)

func main() {
	name := "concursap.cred..CreditCardTransaction"
	pos := strings.LastIndex(name, ".")
	fmt.Println(name[pos+1:])
}
