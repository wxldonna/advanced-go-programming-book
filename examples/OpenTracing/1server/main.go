package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.wdf.sap.corp/velocity/trc"

	abc "1server/openTracing"
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
	trc.Application = "1Server"

	abc.Const(tracer)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //
	title := vars["title"]
	tracer.Infof("get products from parameter title %s", string(title))

	ctx := r.Context()

	tracer.Debugf("version is %s", ctx.Value("version"))
	fmt.Fprintf(w, "version is %s", ctx.Value("version"))
	// create request and call 2Server
	newrequest, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8090/list/apple", nil)
	if err != nil {
		log.Fatal(err)
	}
	// attach span to the outbround request
	//AttachToRequest(ctx, newrequest)
	abc.AttachToRequest(ctx, newrequest)
	client := &http.Client{}
	res, err := client.Do(newrequest)
	if err != nil {
		log.Panicf("create request failed with error %v", err)
		os.Exit(2)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("ReadAll failed with err %v", err)
		os.Exit(3)
	}
	tracer.Debugf("the response data is %s ", string(data))
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Printf("handle the list failed with status %s \n", res.Status)
		os.Exit(1)
	}

}

func main() {
	flag.Parse()
	log.Printf("trcFlag is %v", *trcFlag)

	log.Printf("trcFlag is %v", *trcFlag)
	if err := trc.ReconfigFromString(*trcFlag); err != nil {
		log.Fatal(err)
	}
	abc.Initopentracing("1Server")
	r := mux.NewRouter()
	r.Use(loggingMiddleware, abc.Middleware)
	r.HandleFunc("/products/{title}", GetProducts)

	http.ListenAndServe(":8080", r)
}
