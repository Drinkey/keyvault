package certio

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
)

type CertIO struct {
	Type         string
	ContentBytes []byte
}

func (c CertIO) decoder(b []byte) []byte {
	block, _ := pem.Decode(b)
	return block.Bytes
}

func (c CertIO) reader(f string) []byte {
	text, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal("read certificate file failed", err)
		return nil
	}
	return text
}

// loader read from file f and use pem to decode the data to bytes
func (c CertIO) loader(f string) []byte {
	text := c.reader(f)
	return c.decoder(text)
}

func (c CertIO) encoder() *bytes.Buffer {
	pemBuffer := new(bytes.Buffer)
	pem.Encode(pemBuffer, &pem.Block{
		Type:  c.Type,
		Bytes: c.ContentBytes,
	})
	return pemBuffer
}

// readPemToString reads the file f and encode as PEM, and returns the PEM string
func (c *CertIO) readPemToString(f string) string {
	c.ContentBytes = c.reader(f)
	c.Type = CertificateType
	s, e := c.toPEMString()
	if e != nil {
		log.Print("got empty PEM string")
		return ""
	}
	return s
}

// loadFileToCert reads from file f and parse the content to x509.Certificate
func (c *CertIO) loadFileToCert(f string) (*x509.Certificate, error) {
	c.ContentBytes = c.loader(f)
	return x509.ParseCertificate(c.ContentBytes)
}

// loadFileToPrivKey reads from file f and parse the content to rsa.PrivateKey
func (c *CertIO) loadFileToPrivKey(f string) (*rsa.PrivateKey, error) {
	c.ContentBytes = c.loader(f)
	return x509.ParsePKCS1PrivateKey(c.ContentBytes)
}

// toPEMString returns the c.ContentBytes encoded by PEM
func (c CertIO) toPEMString() (string, error) {
	r := c.encoder()
	return r.String(), nil
}

// Save writes ContentBytes of CertIO encoded by PEM as Type type to file f
func (c CertIO) save(f string, perm os.FileMode) error {
	r := c.encoder()
	return ioutil.WriteFile(f, r.Bytes(), perm)
}
