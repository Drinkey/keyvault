package certio

import (
	"crypto/rsa"
	"crypto/x509"
	"os"
)

var CERT_CONF_FILE = os.Getenv("CERT_CONF_FILE")
var CERT_DIR = os.Getenv("CERT_DIR")
var CertFiles = GetCertFiles(CERT_DIR)

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
