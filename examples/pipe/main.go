package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	//pprof.StartCPUProfile(os.Stdout)
	go func() {
		log.Printf("pprof server is start")
		log.Println(http.ListenAndServe("localhost:8081", nil))
	}()
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	//defer cancel()
	log.Printf("main function start")
	defer log.Printf("main function finished")
	test(ctx)
	//pprof.StopCPUProfile()
	time.Sleep(1000 * time.Hour)
}

func test(ctx context.Context) {

	log.Printf("sub function is called ")
	defer log.Printf("sub function is finished ")
	pr, pw := io.Pipe()

	//br := bufio.NewReader(pr)
	//bw := bufio.NewWriter(pw)

	//sig := make(chan struct{})
	go func(ctx context.Context) {
		log.Printf("slow function is called")
		defer log.Printf("slow function is finished")
		//defer pw.Close()
		defer func() {
			//sig <- struct{}{}
			pw.Close()
			//		bw.Flush()

			//	bw.Close()
			/*
				select {
				case <-ctx.Done():
					return
				default:
					//bw.Flush()
					bw.Flush()
					sig <- struct{}{}
				}

			*/
		}()
		for {
			time.Sleep(1 * time.Second)
			//fmt.Fprintf(pw,"API return")
			n, err := pw.Write([]byte("API return\n"))
			//pw.Close()
			if err != nil {
				panic(err)
			}

			log.Printf("write lengthn===%d", n)
		}

		//log.Printf("Flush Error %v", bw.Flush())
	}(ctx)

	/*
		n, err := io.Copy(os.Stdout, pr)
		if err != nil {
			panic(err)
		}
	*/

	// Creating a buffer

	buffer := new(bytes.Buffer)
	sig := make(chan struct{})
	go func() {
		for {
			time.Sleep(1 * time.Second)
			io.Copy(buffer, pr)

			log.Printf("read  %b", buffer.Bytes())
		}

		sig <- struct{}{}
	}()
	<-sig
	//log.Printf("read  %b", buffer.Bytes())
	/*
		select {
		case <-ctx.Done():
			log.Printf("%v", ctx.Err())
			return
		case <-sig:
			io.Copy(buffer, br)
			log.Printf("read  %s", buffer.String())
		}
	*/
	//buffer.ReadFrom(pr)

}
