package certio

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Subject struct {
	Organization string `json:"organization"`
	Country      string `json:"country"`
	Province     string `json:"province"`
	Locality     string `json:"locality"`
	Address      string `json:"address"`
	PostalCode   string `json:"postal_code"`
	CommonName   string `json:"common_name"`
}

type CaCertConfig struct {
	SerialNumber int64   `json:"serial_number"`
	Subject      Subject `json:"subject"`
	Valid        int     `json:"valid_year"`
	KeyLength    int     `json:"key_length"`
}

type WebCertConfig struct {
	SerialNumber int64   `json:"serial_number"`
	Subject      Subject `json:"subject"`
	DNSName      string  `json:"dns_name"`
	Valid        int     `json:"valid_year"`
	KeyLength    int     `json:"key_length"`
}

type CertJSON struct {
	CA  CaCertConfig  `json:"ca"`
	Web WebCertConfig `json:"web"`
}

type CertificateAuthority struct {
	Certificate *x509.Certificate
	String      string

	certFlag   bool // true if already set
	PrivateKey *rsa.PrivateKey
	pkeyFlag   bool // true if already set
}

func (c *CertificateAuthority) Cache(cert *x509.Certificate, s string, p *rsa.PrivateKey) {
	c.SetCertificate(cert)
	c.SetPrivateKey(p)
	c.String = s
}

func (c *CertificateAuthority) SetCertificate(cert *x509.Certificate) {
	c.Certificate = cert
	c.certFlag = true
}

func (c *CertificateAuthority) SetPrivateKey(p *rsa.PrivateKey) {
	c.PrivateKey = p
	c.pkeyFlag = true
}

func (c CertificateAuthority) IsSet() bool {
	return c.certFlag && c.pkeyFlag
}

// CertFilePath Stores certificate related file path
type CertFilePaths struct {
	CaCertPath     string
	CaPrivKeyPath  string
	WebCertPath    string
	WebPrivKeyPath string
}

type CertificateConfiguration struct {
	Paths  CertFilePaths
	dir    string
	file   string // the JSON config file
	config *CertJSON
}

func (config *CertificateConfiguration) Parse() {
	config.getKvCertDir()
	config.getKvCertConfig()
	config.getCertPaths()
	config.parseJSON()
}

func (config *CertificateConfiguration) parseJSON() {
	var schema *CertJSON
	log.Printf("reading config file %s", config.file)
	contentBytes, _ := ioutil.ReadFile(config.file)
	_ = json.Unmarshal(contentBytes, &schema)
	config.config = schema
	log.Println(config.config)
}

func (config *CertificateConfiguration) getKvCertDir() {
	log.Print("reading env cert dir")
	config.dir = os.Getenv("KV_CERT_DIR")
	log.Printf("got config.dir = %s", config.dir)
}

func (config *CertificateConfiguration) getKvCertConfig() {
	log.Print("reading env cert config")
	config.file = os.Getenv("KV_CERT_CONF")
	log.Printf("got config.file = %s", config.file)
}

func (config *CertificateConfiguration) getCertPaths() {
	// var dir = getKvCertDir()
	log.Printf("got cert dir %s", config.dir)
	const caFile = "ca.crt"
	const certFile = "cert.pem"
	config.Paths = CertFilePaths{
		CaCertPath:     fmt.Sprintf("%s/%s", config.dir, caFile),
		CaPrivKeyPath:  fmt.Sprintf("%s/ca_priv.key", config.dir),
		WebCertPath:    fmt.Sprintf("%s/%s", config.dir, certFile),
		WebPrivKeyPath: fmt.Sprintf("%s/cert_priv.key", config.dir),
	}
}
