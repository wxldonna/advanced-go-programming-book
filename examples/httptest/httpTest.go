package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

func main() {
	TestHeader3D()
}

func TestHeader3D() {
	resp := httptest.NewRecorder()

	uri := "/3D/header/?"
	path := "/home/test"
	unlno := "997225821"

	param := make(url.Values)
	param["param1"] = []string{path}
	param["param2"] = []string{unlno}

	req, err := http.NewRequest("GET", uri+param.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}

	http.DefaultServeMux.ServeHTTP(resp, req)
	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		log.Fatal(err)
	} else {
		if strings.Contains(string(p), "Error") {
			log.Fatal("header response shouldn't return error: %s", p)
		} else if !strings.Contains(string(p), `expected result`) {
			log.Fatal("header response doen't match:\n%s", string(p))
		}
	}
}
