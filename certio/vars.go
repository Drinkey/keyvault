package certio

import (
	"crypto/rsa"
	"crypto/x509"
)

const CA_FILE = "ca.crt"
const CERT_FILE = "cert.pem"

// Store certificate related file path
type CertFilePath struct {
	CaCertPath        string
	CaPrivKeyPath     string
	ServerCertPath    string
	ServerPrivKeyPath string
}

// store CA cert and CA private key
type CertificateAuthority struct {
	CaCert      *x509.Certificate
	CaCertBytes []byte
	CaPrivKey   *rsa.PrivateKey
}
