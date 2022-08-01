package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	//pprof.StartCPUProfile(os.Stdout)
	go func() {
		log.Printf("pprof server is start")
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	log.Printf("main function start")
	defer log.Printf("main function finished")
	test(ctx)
	//pprof.StopCPUProfile()

	for {

	}
}

func test(ctx context.Context) {
	log.Printf("sub function is called ")
	defer log.Printf("sub function is finished ")
	msg := make(chan string)
	go func(ctx context.Context) {
		log.Printf("slow function is called")
		defer log.Printf("slow function is finished")
		time.Sleep(20 * time.Second)
		select {
		case msg <- "API return":
			//case <-ctx.Done():
			//	close(msg)
			//	return
		}
	}(ctx)

	select {
	case <-msg:
		log.Printf("get data from slow function ")
	case <-ctx.Done():
		//	close(msg)
		log.Printf("timeout error %v", ctx.Err())
	}
}
