package main

import (
	"context"
	"log"
	"net/http"

	"fmt"

	"github.com/gorilla/mux"
)

//路由参数
func main() {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	// r.HandleFunc("/articles/{title}", TitleHandler)

	r.HandleFunc("/articles/{title:[a-z]+}", TitleHandler)

	http.ListenAndServe(":8080", r)
}

//https://github.com/gorilla/mux#examples
func TitleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // 获取参数
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "title: %v\n", vars["title"])
	ctx := r.Context()
	fmt.Fprintf(w, "version is %s", ctx.Value("version"))
}

// logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// add the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "version", "1.0.0")

		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
