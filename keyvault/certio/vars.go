package certio

import (
	"crypto/rsa"
	"crypto/x509"
	"os"
)

var CONF_DIR = os.Getenv("CERT_CONF_DIR")

const CA_CONF = "ca.json"
const CERT_CONF = "cert.json"

var CERT_DIR = os.Getenv("CERT_DIR")

const CA_FILE = "ca.crt"
const CERT_FILE = "cert.pem"

type SubjectConfig struct {
	Organization string `json: "organization"`
	Country      string `json:"country"`
	Province     string `json:"province"`
	Locality     string `json:"locality"`
	Address      string `json:"address"`
	PostalCode   string `json:"postal_code"`
	CommonName   string `json:"common_name"`
}

type CAConfig struct {
	SerialNumber int64         `json:"serial_number"`
	Subject      SubjectConfig `json:"subject"`
	Valid        int           `json:"valid_year"`
	KeyLength    int           `json:"key_length"`
}

type CertConfig struct {
	SerialNumber int64         `json:"serial_number"`
	Subject      SubjectConfig `json:"subject"`
	DNSName      string        `json:"dns_name"`
	Valid        int           `json:"valid_year"`
	KeyLength    int           `json:"key_length"`
}

// Store certificate related file path
type CertFiles struct {
	CaCert        string
	CaPrivKey     string
	ServerCert    string
	ServerPrivKey string
}

type CertConfFiles struct {
	CaCertConf     string
	ServerCertConf string
}

// store CA cert and CA private key
type CertificateAuthority struct {
	CaCert    *x509.Certificate
	CaPrivKey *rsa.PrivateKey
}
