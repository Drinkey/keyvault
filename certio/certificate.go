package certio

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"

	"github.com/Drinkey/keyvault/pkg/settings"
	"github.com/Drinkey/keyvault/pkg/utils"
)

const (
	// PrivateKeyType represents the "RSA Private key" String
	PrivateKeyType = "RSA PRIVATE KEY"
	// CertificateType represents the "Certificate" String
	CertificateType = "CERTIFICATE"
)

func getSubjectName(subject settings.Subject) pkix.Name {
	return pkix.Name{
		Organization:  []string{subject.Organization},
		Country:       []string{subject.Country},
		Province:      []string{subject.Province},
		Locality:      []string{subject.Locality},
		StreetAddress: []string{subject.Address},
		PostalCode:    []string{subject.PostalCode},
		CommonName:    subject.CommonName,
	}
}

type WebCertificate struct {
	io         CertIO
	PrivateKey PrivateKey
}

func (c WebCertificate) CreateTemplate(config settings.WebCertConfig) *x509.Certificate {
	NotBefore, NotAfter := utils.TimeRange(config.Valid)
	return &x509.Certificate{
		SerialNumber: big.NewInt(config.SerialNumber),
		Subject:      getSubjectName(config.Subject),
		DNSNames:     []string{config.DNSName, "localhost"},
		NotBefore:    NotBefore,
		NotAfter:     NotAfter,
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
}
func (c WebCertificate) Save(f string, content []byte) error {
	c.io.ContentBytes = content
	c.io.Type = CertificateType
	return c.io.save(f, 0600)
}

func (c WebCertificate) Read(f string) string {
	return c.io.readPemToString(f)
}

type CA struct {
	io         CertIO
	privateKey PrivateKey
	Bytes      []byte // CA cert in []byte
	String     string //CA Cert in string
}

func (ca CA) CreateTemplate(config settings.CaCertConfig) *x509.Certificate {
	NotBefore, NotAfter := utils.TimeRange(config.Valid)

	return &x509.Certificate{
		SerialNumber:          big.NewInt(config.SerialNumber),
		Subject:               getSubjectName(config.Subject),
		NotBefore:             NotBefore,
		NotAfter:              NotAfter,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

func (ca CA) Issue(cacert *x509.Certificate, cert *x509.Certificate, pub interface{}, priv interface{}) ([]byte, error) {
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cacert, pub, priv)
	if err != nil {
		return nil, err
	}
	return certBytes, nil
}

func (ca CA) Save(f string, content []byte) error {
	ca.io.ContentBytes = content
	ca.io.Type = CertificateType
	return ca.io.save(f, 0600)
}

// Load reads certificate and private key from file specified in certPath and privkeyPath,
// then returns *x509.Certificate and *rsa.PrivateKey of CA
func (ca *CA) Load(certPath, privkeyPath string) (*x509.Certificate, *rsa.PrivateKey) {
	cert, err := ca.io.loadFileToCert(certPath)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	// ca.Bytes = ca.io.ContentBytes
	ca.io.Type = CertificateType
	ca.String, err = ca.io.toPEMString()
	if err != nil {
		log.Print("unable to save PEM string")
		log.Println(err)
		log.Println(ca.io.ContentBytes)
		return nil, nil
	}

	pkey, err := ca.privateKey.io.loadFileToPrivKey(privkeyPath)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	return cert, pkey
}

// Read returns the PEM encoded string in file f
func (ca *CA) Read(f string) string {
	return ca.io.readPemToString(f)
}

type PrivateKey struct {
	io CertIO
}

func (p PrivateKey) Generate(keyLen int) (*rsa.PrivateKey, error) {
	pk, err := rsa.GenerateKey(rand.Reader, keyLen)
	if err != nil {
		return nil, err
	}
	return pk, nil
}

func (p PrivateKey) Save(f string, key *rsa.PrivateKey) error {
	p.io.ContentBytes = x509.MarshalPKCS1PrivateKey(key)
	p.io.Type = PrivateKeyType
	return p.io.save(f, 0600)
}

type CertificateRequest struct {
	io      CertIO
	Request string `json:"csr"`
}

func (c CertificateRequest) CreateTemplate(csr *x509.CertificateRequest) *x509.Certificate {
	notBefore, notAfter := utils.TimeRange(5)
	return &x509.Certificate{
		SerialNumber:       big.NewInt(210201),
		Subject:            csr.Subject,
		NotBefore:          notBefore,
		NotAfter:           notAfter,
		Signature:          csr.Signature,
		SignatureAlgorithm: csr.SignatureAlgorithm,
		PublicKeyAlgorithm: csr.PublicKeyAlgorithm,
		PublicKey:          csr.PublicKey,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:           x509.KeyUsageDigitalSignature,
	}
}

func (c CertificateRequest) ParsePEMString(pem string) (*x509.CertificateRequest, error) {
	b := c.io.decoder([]byte(pem))
	csr, err := x509.ParseCertificateRequest(b)
	if err != nil {
		log.Printf("Failed to parse Certificate Request in CSR: %s", err.Error())
		return nil, err
	}

	if err = csr.CheckSignature(); err != nil {
		log.Printf("check csr signature failed: %s", err.Error())
		return nil, err
	}
	return csr, nil
}

func IssueCertificate(csrString string) (certPem, caPem string) {
	if !CaContainer.IsSet() {
		log.Fatal("Something is wrong, unable to find CA in memory")
		return
	}
	var req CertificateRequest

	csr, err := req.ParsePEMString(csrString)
	if err != nil {
		log.Fatal("Failed to parse certificate request")
	}
	reqTemplate := req.CreateTemplate(csr)

	var ca CA
	reqCert, err := ca.Issue(CaContainer.Certificate, reqTemplate, &csr.PublicKey, CaContainer.PrivateKey)
	if err != nil {
		log.Fatal("Issue certificate failed")
		return
	}

	var certIo CertIO
	certIo.ContentBytes = reqCert
	certIo.Type = CertificateType
	certPem, err = certIo.toPEMString()

	caPem = CaContainer.String
	return
}
