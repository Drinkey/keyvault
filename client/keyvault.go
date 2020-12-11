package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	url      = flag.String("url", "someUrl", "the URL to get")
	method   = flag.String("method", "GET", "Specify the HTTP Method")
	certFile = flag.String("cert", "client.crt", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "client.pkey", "A PEM encoded private key file.")
	caFile   = flag.String("CA", "ca.crt", "A PEM eoncoded CA's certificate file.")
)

// func certificateInit(){
// 	parse config
// 	if cert file exist:
// 		return
// 	else
// 		create pkey

// 	if pkey file exist:
// 		create new cert
// 		return
// }

// func createNewCertificate(){
// 	parse cert config
// 	set OU same as namespace
// 	create csr PEM
// 	parse csr PEM and replace \n with \\n (convert the file in one line)
// 	post to service
// 	get the signed cert and save to [namespace].crt
// }

func createHTTPSClient() *http.Client {
	// Load client cert
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: transport}
}

func main() {
	flag.Parse()

	client := createHTTPSClient()

	// Do GET something
	resp, err := client.Get(*url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
}
