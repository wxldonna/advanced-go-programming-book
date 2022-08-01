package trc

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
)

var offset uint64

//go:noescape
func id(offset uint64) int64

func GoroutineID() int64 {
	return id(offset)
}

func init() {
	idFromStack := func() int64 {
		b := make([]byte, 32)
		b = b[:runtime.Stack(b, false)]
		b = bytes.TrimPrefix(b, []byte("goroutine "))
		b = b[:bytes.IndexByte(b, ' ')]
		n, _ := strconv.ParseInt(string(b), 10, 64)
		return n
	}
	thisID := idFromStack()
	var i uint64
Loop:
	for i = 0; i < 100*8; i += 8 {
		cand := id(i)
		if cand == thisID {
			attempts := 10
			var wg sync.WaitGroup
			wg.Add(attempts)
			c := make(chan bool, attempts)
			for j := 0; j < attempts; j++ {
				go func() {
					gID := idFromStack()
					c <- gID == id(i)
					wg.Done()
				}()
			}
			wg.Wait()
			close(c)
			for b := range c {
				if !b {
					continue Loop
				}
			}
			offset = i
			return
		}
	}
}
