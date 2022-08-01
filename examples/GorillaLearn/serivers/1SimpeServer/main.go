package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"log"

	"github.com/gorilla/mux"

 	_ "net/http/pprof"
)

const serveraddress = "localhost:8080"

var ctx = context.Background()

type server struct {
	*http.Server
	wg *sync.WaitGroup
}
type MiddlewareFunc func(http.Handler) http.Handler

func HomeHandller(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//version := ctx.Value("version").(string)
	fmt.Fprintf(w, "version %v \n", ctx.Value("version"))
	w.Write([]byte("hello home package"))
}
func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//	w.Write([]byte("product is called \n"))
	fmt.Fprintf(w, "version %v \n", ctx.Value("version"))
	fmt.Fprintf(w, "param is %v \n", vars["title"])

}

// logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// add the context
		ctx = context.WithValue(ctx, "version", "1.0.0")
		vars := mux.Vars(r) // 获取参数
		log.Printf("1variables is %v\n", vars["title"])
		//r = r.Clone(ctx)
		log.Printf("2variables is %v\n", vars["title"])
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// subrouter which use the middleware
func subMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("subrouter is called ")
		next.ServeHTTP(w, r)
	})
}

// subrouter which use the parameters is called
func subMiddlewarePara(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("subMiddlewarePara is called ")
		next.ServeHTTP(w, r)
	})
}

func TitleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // 获取参数
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "title: %v\n", vars["title"])
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/debug/pprof").Handler(http.DefaultServeMux)
	r.HandleFunc("/", HomeHandller)
	r.HandleFunc("/products/{title}", ProductsHandler)
	r.HandleFunc("/articles/{title}", TitleHandler)
	RegisterHandlers(r)
	r.Use(loggingMiddleware)

	//PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	/*
		//create server using http.server
		ctx := context.Background()
		ctx = context.WithValue(ctx, "version", "1.0.0")
		startServer(r, ctx)
	*/
	log.Printf("server start")
	/*
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	 */
	if err := http.ListenAndServe(serveraddress, r); err != nil {
		log.Fatalf("server start failed %v", err)
	}

}

func Addresshandller(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "address is called ")
}

func Beijing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.Path)
}

//handle the parameters
func MyHome(w http.ResponseWriter, r *http.Request) {
	// get the varaiables
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Area: %v\n", vars["area"])
	//	fmt.Fprintf(w, "ID is %v \n", vars["id1"])
}

func RegisterHandlers(base *mux.Router) {

	subrouter := base.NewRoute().Subrouter()
	subrouter.HandleFunc("/address", Addresshandller)
	subPrefix := subrouter.PathPrefix("/city").Subrouter()
	subPrefix.HandleFunc("/beijing", Beijing)

	subPrefix.Use(subMiddleware)
	// r.HandleFunc("/articles/{category}/", ArticlesCategoryHandler)
	// r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
	//subrouterwithpara := subPrefix.NewRoute().Subrouter()
	// /{id:[0:9]+}
	base.HandleFunc("/haidian/{area}", MyHome)
	base.Use(subMiddlewarePara)
}

func startServer(r *mux.Router, ctx context.Context) *server {
	srv := &server{
		Server: &http.Server{
			Addr:    serveraddress,
			Handler: r,
			BaseContext: func(net.Listener) context.Context {
				return ctx
			},
		},
		wg: &sync.WaitGroup{},
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Printf("server failed with error %v", err)
	}
	log.Printf("server is start")
	srv.wg.Add(1)

	go func() {
		defer srv.wg.Done()
		if err := srv.Serve(ln); err != nil {
			log.Printf("server failed with error %v", err)
			return
		}

		log.Printf("seriver is stopped ")
	}()
	srv.wg.Wait()
	return srv
}
