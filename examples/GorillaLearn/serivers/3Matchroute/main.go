package main

import (
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
)

//路由参数
func main() {
	r := mux.NewRouter()
	// Only matches if domain is "www.example.com".
	//r.Host("localhost:8080")
	// Matches a dynamic subdomain.
	//	sub:=r.Host("{subdomain:[a-z]+}.example.com").Subrouter()
	// only when the method is 'GET', the handkle function will be called
	// r.HandleFunc("/products/{title:[a-z]+}", TitleHandler).Methods("GET")

	// There are several other matchers that can be added. To match path prefixes:
	// the pathfrefix must be the part of path in the hannlder

	//pathSub := r.PathPrefix("/wxl").Subrouter()
	//pathSub.HandleFunc("/products/{title:[a-z]+}", TitleHandler)

	//r.HandleFunc("/products/{title:[a-z]+}", TitleHandler).PathPrefix("/products")
	// how to get the value and key ?
	// http://localhost:8080/products/wxl?key=value
	//r.HandleFunc("/products/{title:[a-z]+}", TitleHandler).Queries("key", "value")
	queryRouter := r.Queries("key", "value").Subrouter()
	queryRouter.HandleFunc("/products/{title:[a-z]+}", TitleHandler)

	// combine the domain method and schema
	//r.HandleFunc("/products", ProductsHandler).
	//Host("www.example.com").
	//Methods("GET").
	//Schemes("http")

	// use a custom matcher function
	// only when the request has the key 'X-Requested-With', the value should be 'wxldonna', the handller function can be called
	// Accept: */*
	// User-Agent: Thunder Client (https://www.thunderclient.io)
	// X-Requested-With: wxldonna
	/*
		r.HandleFunc("/products/{title:[a-z]+}", TitleHandler).Methods("GET").MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			if r.Header.Get("X-Requested-With") == "wxldonna" {
				return true
			}
			return false
		})
	*/
	http.ListenAndServe("localhost:8080", r)
}

//https://github.com/gorilla/mux#examples
func TitleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // 获取参数
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "title: %v\n", vars["title"])
}
