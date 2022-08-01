package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"sync"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		wg:=sync.WaitGroup{}
		wg.Add(10000)
		for i:=0;i<10000;i++{

			go func() {
				log.Printf("number is %d",i)
				wg.Done()
			}()

		}
		wg.Wait()
		defer pprof.StopCPUProfile()
	}
}