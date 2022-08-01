package main

import (
	"counterfeiterTest/foo/foofakes"
	"errors"
	"testing"

	_ "github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"
)

func main() {
	errHello := errors.New("hello")
	fake1 := foofakes.FakeMySpecialInterface1{}
	// set the return value
	fake1.DoThings1Returns(100, errHello)
	num, err := fake1.DoThings1("hello", 10)
	if !errors.Is(err, errHello) {
		panic(err)
	}
	t := &testing.T{}
	assert.Equal(t, 100, num)

	//Fakes record the arguments they were called with:

	fake2 := foofakes.FakeMySpecialInterface2{}
	fake2.DoThings2("hello", 200)
	str, num1 := fake2.DoThings2ArgsForCall(0)
	assert.Equal(t, "hello", str)
	assert.Equal(t, 200, num1)

}
