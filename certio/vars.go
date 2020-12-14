/*
Package certio provides all operations against certificate.
*/
package certio

import (
	"crypto/rsa"
	"crypto/x509"
)

const caFile = "ca.crt"
const certFile = "cert.pem"

// CertFilePath Stores certificate related file path
type CertFilePath struct {
	CaCertPath        string
	CaPrivKeyPath     string
	ServerCertPath    string
	ServerPrivKeyPath string
}

// CertificateAuthority stores CA cert and CA private key
type CertificateAuthority struct {
	CaCert      *x509.Certificate
	CaCertBytes []byte
	CaPrivKey   *rsa.PrivateKey
}
