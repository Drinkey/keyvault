package certio

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func CreateTLSConfig(c CertFiles) (*tls.Config, error) {
	// 1. load ca to cert pool
	// 2. load cert and key
	// 3. construct tls config
	caPEM, err := ioutil.ReadFile(c.CaCert)
	if err != nil {
		log.Panic("load CA certificate failed, unable to recover")
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caPEM)
	if !ok {
		log.Panic("append CA certificate to CertPool failed")
	}

	cert, err := tls.LoadX509KeyPair(c.ServerCert, c.ServerPrivKey)
	if err != nil {
		log.Panic("load server certificate failed")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    roots,
	}, nil
}

func BuildTLSConfig() (*tls.Config, CertFiles) {
	f := GetCertFiles(CONF_DIR)
	tlsConfig, err := CreateTLSConfig(f)
	if err != nil {
		log.Panic("create TLS server config failed")
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig, f
}

func ParseClientCertOU(r *http.Request) (string, bool) {
	mode := os.Getenv("KVMODE")
	// not panic only when mode is TEST if connection is TLS.
	// getting peer certificate will panic if connection is not TLS
	if r.TLS == nil && mode == "TEST" {
		// later "false" means TLS is not enabled
		return "", false
	}
	// The peer certificate must be the first one
	peerCert := r.TLS.PeerCertificates[0]
	return peerCert.Subject.OrganizationalUnit[0], true
}
