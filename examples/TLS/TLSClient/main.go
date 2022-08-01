package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", "localhost:4000", "HTTPS server address")
	certFile := flag.String("certfile", "/home/vagrant/advanced-go-programming-book/examples/TLS/TLSServer/cert.pem", "trusted CA certificate") // server certificate

	// if user need provide the certificate , need client certifiacte and private key
	clientCertFile := flag.String("clientcert", "/home/vagrant/advanced-go-programming-book/examples/TLS/TLSClient/clientcert.pem", "certificate PEM for client")
	clientKeyFile := flag.String("clientkey", "/home/vagrant/advanced-go-programming-book/examples/TLS/TLSClient/clientkey.pem", "key PEM for client")

	flag.Parse()

	// Load our client certificate and key.
	clientCert, err := tls.LoadX509KeyPair(*clientCertFile, *clientKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	cert, err := os.ReadFile(*certFile)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("unable to parse cert from %s", *certFile)
	}

	conf := &tls.Config{
		// server certificate
		RootCAs: certPool,
		// client certificate
		Certificates: []tls.Certificate{clientCert},
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// server certificate
				RootCAs: certPool,
				// client certificate
				Certificates: []tls.Certificate{clientCert},
			},

			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				c, err := tls.Dial("tcp", addr, conf)
				//c, err := net.Dial("tcp", addr) //localhost:4000:https
				if err != nil {
					panic(err)
				}
				log.Printf("DialTLSContext network and addr is %s,%s", network, addr)
				return c, err
			},
		},
	}

	r, err := client.Get("https://" + *addr)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	html, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", r.Status)
	fmt.Printf(string(html))
}
