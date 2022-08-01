package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.wdf.sap.corp/velocity/trc"

	abc "github.wdf.sap.corp/abap/openTracing"
)

var tracer = trc.InitTraceTopic("Initlization", "executable package")

var trcFlag = flag.String("trc", "", "e.g. -trc=debug,main:warning")

// logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		tracer.Debugf("Request URL %s", r.RequestURI)
		// get the context from request
		ctx := r.Context()
		// add the context
		ctx = context.WithValue(ctx, "version", "1.0.0")

		r = r.WithContext(ctx)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func init() {
	trc.Application = "2Server"
	abc.Const(tracer)
}

func GetTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //
	title := vars["title"]
	tracer.Infof("get products from parameter title %s", string(title))

	ctx := r.Context()

	tracer.Debugf("version is %s", ctx.Value("version"))
	fmt.Fprintf(w, "version is %s", ctx.Value("version"))
	// write response
	fmt.Fprintf(w, title)
}

func main() {
	flag.Parse()
	log.Printf("trcFlag is %v", *trcFlag)
	if err := trc.ReconfigFromString(*trcFlag); err != nil {
		log.Fatal(err)
	}
	//Initopentracing("2Server")
	abc.Initopentracing("2Server")
	r := mux.NewRouter()
	r.Use(loggingMiddleware, abc.Middleware)
	r.HandleFunc("/list/{title}", GetTitle)

	if err := http.ListenAndServe(":8090", r); err != nil {
		log.Fatal(err)
	}
}
