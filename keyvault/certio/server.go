package certio

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func CreateTLSConfig(c CertFilePath) (*tls.Config, error) {
	// 1. load ca to cert pool
	// 2. load cert and key
	// 3. construct tls config
	caPEM, err := ioutil.ReadFile(c.CaCertPath)
	if err != nil {
		log.Panic("load CA certificate failed, unable to recover")
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caPEM)
	if !ok {
		log.Panic("append CA certificate to CertPool failed")
	}

	cert, err := tls.LoadX509KeyPair(c.ServerCertPath, c.ServerPrivKeyPath)
	if err != nil {
		log.Panic("load server certificate failed")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    roots,
	}, nil
}

func BuildTLSConfig() *tls.Config {
	// f := GetCertFiles(CERT_DIR)
	tlsConfig, err := CreateTLSConfig(CertFiles)
	if err != nil {
		log.Panic("create TLS server config failed")
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig
}

func ParseClientCertOU(r *http.Request) (string, bool) {
	kvMode := os.Getenv("KV_MODE")
	ginMode := os.Getenv("GIN_MODE")
	// not panic only when mode is TEST and gin in debug mode if connection is TLS.
	// getting peer certificate will panic if connection is not TLS
	if r.TLS == nil && kvMode == "TEST" && ginMode == "debug" {
		// later "false" means TLS is not enabled
		log.Printf("Allow non-TLS connection in TEST mode, gin is debug mode")
		return "", false
	}
	// The peer certificate must be the first one
	peerCert := r.TLS.PeerCertificates[0]
	return peerCert.Subject.OrganizationalUnit[0], true
}
