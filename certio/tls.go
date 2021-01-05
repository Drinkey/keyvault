/*
Package certio provides all operations against certificate.
*/
package certio

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func LoadCertificates(c CertFilePaths) ([]byte, tls.Certificate) {
	// 1. load ca to cert pool
	// 2. load cert and key
	// 3. construct tls config

	caPEM, err := ioutil.ReadFile(c.CaCertPath)
	if err != nil {
		log.Printf("load CA certificate [%s] failed", c.CaCertPath)
		log.Panic("panic because no CA found")
	}

	cert, err := tls.LoadX509KeyPair(c.WebCertPath, c.WebPrivKeyPath)
	if err != nil {
		log.Printf("cert=%s, pkey=%s", c.WebCertPath, c.WebPrivKeyPath)
		log.Panic("load server certificate failed")
	}

	return caPEM, cert
}

func BuildTLSConfig(certs CertFilePaths, level string) *tls.Config {
	clientAuthMap := map[string]tls.ClientAuthType{
		"SECRET":      tls.RequireAndVerifyClientCert,
		"MAINTENANCE": tls.VerifyClientCertIfGiven,
	}
	auth, ok := clientAuthMap[level]
	if !ok {
		log.Panic("unable to create client auth profile because level is unknown")
	}
	log.Printf("Build config for %s, type: %d", level, auth)

	caPEM, webcert := LoadCertificates(certs)

	roots := x509.NewCertPool()
	ok = roots.AppendCertsFromPEM(caPEM)
	if !ok {
		log.Panic("append CA certificate to CertPool failed")
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{webcert},
		ClientAuth:   auth,
		ClientCAs:    roots,
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
