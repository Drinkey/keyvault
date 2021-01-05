package request

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"

	"github.com/Drinkey/keyvault/certio"
)

func send(client *http.Client, method, url string, header http.Header, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type Requests struct {
	Client       *http.Client
	Header       http.Header
	Certificates certio.CertFilePaths
}

func (r *Requests) InitClient(insecure bool) {
	var tlsConfig *tls.Config
	if insecure {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		caCert, cert := certio.LoadCertificates(r.Certificates)
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Setup HTTPS client
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	r.Client = &http.Client{Transport: transport}
}

// func (r Requests)addHeader() http.Header{
// 	var h http.Header
// 	if len(r.header) != 0 {
// 		for k, v := range r.header {
// 			h.Add(k, v)
// 		}
// 	}
// 	return h
// }

func (r Requests) Get(url string) (*http.Response, error) {
	return send(r.Client, "GET", url, r.Header, nil)
}

func (r Requests) Post(url string, body io.Reader) (*http.Response, error) {
	return send(r.Client, "POST", url, r.Header, body)
}

func (r Requests) Put(url string, body io.Reader) (*http.Response, error) {
	return send(r.Client, "PUT", url, r.Header, body)
}

func (r Requests) Delete(url string) (*http.Response, error) {
	return send(r.Client, "DELETE", url, r.Header, nil)
}
