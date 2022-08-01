package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	addr := flag.String("addr", "localhost:4040", "dial address")
	certFile := flag.String("certfile", "/home/vagrant/advanced-go-programming-book/examples/TLS/ExaminServerCertificates/cert.pem", "trusted CA certificate")

	flag.Parse()

	cert, err := os.ReadFile(*certFile)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("unable to parse cert from %s", *certFile)
	}

	cfg := tls.Config{RootCAs: certPool}
	conn, err := tls.Dial("tcp", *addr, &cfg)
	if err != nil {
		log.Fatal("TLS connection failed: " + err.Error())
	}
	defer conn.Close()

	certChain := conn.ConnectionState().PeerCertificates
	for i, cert := range certChain {
		fmt.Println(i)
		fmt.Println("Issuer:", cert.Issuer)
		fmt.Println("Subject:", cert.Subject)
		fmt.Println("Version:", cert.Version)
		fmt.Println("NotAfter:", cert.NotAfter)
		fmt.Println("DNS names:", cert.DNSNames)
		fmt.Println("")
	}
}
