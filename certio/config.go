package certio

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"log"

	"github.com/Drinkey/keyvault/pkg/settings"
)

type CertificateAuthority struct {
	Certificate *x509.Certificate
	String      string
	certFlag    bool // true if already set
	PrivateKey  *rsa.PrivateKey
	pkeyFlag    bool // true if already set
}

// Cache saves CA cert and CA private key in memory, and keeps CA cert in PEM encoded
// string in CertificateAuthority.String
func (c *CertificateAuthority) Cache(cert *x509.Certificate, s string, p *rsa.PrivateKey) {
	c.setCertificate(cert)
	c.setPrivateKey(p)
	c.String = s
}

func (c *CertificateAuthority) setCertificate(cert *x509.Certificate) {
	c.Certificate = cert
	c.certFlag = true
}

func (c *CertificateAuthority) setPrivateKey(p *rsa.PrivateKey) {
	c.PrivateKey = p
	c.pkeyFlag = true
}

// IsSet returns true if certificate and private key are already cached
func (c CertificateAuthority) IsSet() bool {
	return c.certFlag && c.pkeyFlag
}

// CertFilePaths is a collection of certificate related file paths
type CertFilePaths struct {
	CaCertPath     string
	CaPrivKeyPath  string
	WebCertPath    string
	WebPrivKeyPath string
}

// CertificateConfiguration has all parameters of certio configuration
type CertificateConfiguration struct {
	Paths  CertFilePaths
	dir    string
	file   string // the JSON config file
	config *settings.CertJSON
}

// Parse initializes the parameters from settings.Settings
func (config *CertificateConfiguration) Parse() {
	config.config = &settings.Settings.Certificate
	config.dir = settings.Settings.CertificateDir
	config.file = settings.Settings.ConfigFile
	config.getCertPaths()
}

func (config *CertificateConfiguration) getCertPaths() {
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
