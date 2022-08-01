package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

type transport struct {
	base http.RoundTripper
	cnt  int
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Printf("request counter %d", t.cnt)
	log.Printf("request header is %v", r.Header)
	t.cnt = t.cnt + 1
	return t.base.RoundTrip(r)
}

func NewTransport() *transport {
	return &transport{
		cnt: 0,
		//base: http.DefaultTransport,
		base: &http.Transport{
			Proxy: func(r *http.Request) (*url.URL, error) {
				log.Printf("request URL %v", r.URL)
				return r.URL, nil
			},
			DialContext: func(ctx context.Context, network, addr string) (c net.Conn, err error) {

				c, err = net.Dial("tcp", addr)
				if err != nil {
					panic(err)
				}

				log.Printf("network and addr is %s,%s", network, addr)
				return c, err
			},
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				c, err := net.Dial("tcp", addr)
				if err != nil {
					panic(err)
				}
				log.Printf("DialTLSContext network and addr is %s,%s", network, addr)
				return c, err
			},
		},
	}
}

func main() {
	// create http client
	client := http.Client{
		Transport: NewTransport(),
	}

	// create request
	ctx := context.TODO()
	ctx = context.WithValue(ctx, "header1", "test")
	var buf bytes.Buffer
	buf.WriteString("hello xiaoliang")
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/test", &buf)
	for i := 0; i < 10; i++ {
		res, err := client.Do(r)
		if err != nil {
			panic(err)
		}

		result, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		log.Printf("response is %s", string(result))
	}

}
