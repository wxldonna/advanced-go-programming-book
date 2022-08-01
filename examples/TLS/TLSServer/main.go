package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTPS network address")
	certFile := flag.String("certfile", "/home/vagrant/advanced-go-programming-book/examples/TLS/TLSServer/cert.pem", "certificate PEM file")
	keyFile := flag.String("keyfile", "/home/vagrant/advanced-go-programming-book/examples/TLS/TLSServer/key.pem", "key PEM file")
	// Load our client certificate and key.
	serverCert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}
	// if we need verify the client certificate , need this
	clientCertFile := flag.String("clientcert", "/home/vagrant/advanced-go-programming-book/examples/TLS/TLSClient/clientcert.pem", "certificate PEM for client authentication")
	// Trusted client certificate.
	clientCert, err := os.ReadFile(*clientCertFile)
	if err != nil {
		log.Fatal(err)
	}
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCert)

	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "Proudly served with Go and HTTPS!")
	})

	srv := &http.Server{
		Addr:    *addr,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
			//server certificates
			Certificates: []tls.Certificate{serverCert},
			// the two lines are for client certificate
			ClientCAs:  clientCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	log.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS(*certFile, *keyFile)
	log.Fatal(err)
}
