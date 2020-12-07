package certio

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"github.com/Drinkey/keyvault/internal"
)

func LoadCACertificate(c CertFilePath) (CertificateAuthority, error) {
	ca_txt, err := ioutil.ReadFile(c.CaCertPath)
	if err != nil {
		log.Fatal("read certificate file failed", err)
	}
	ca_block, _ := pem.Decode([]byte(ca_txt))
	cacert, err := x509.ParseCertificate(ca_block.Bytes)
	if err != nil {
		log.Fatal("parse certificate content failed", err)
	}

	ca_pkey_txt, err := ioutil.ReadFile(c.CaPrivKeyPath)
	if err != nil {
		log.Fatal("read private key file failed", err)
	}
	ca_pkey_block, _ := pem.Decode([]byte(ca_pkey_txt))
	ca_pkey, err := x509.ParsePKCS1PrivateKey(ca_pkey_block.Bytes)
	if err != nil {
		log.Fatal("parse private key content failed", err)
	}

	return CertificateAuthority{CaCert: cacert, CaPrivKey: ca_pkey}, nil
}

func Issue(c *x509.Certificate, ca CertificateAuthority, k interface{}) ([]byte, error) {
	certBytes, err := x509.CreateCertificate(rand.Reader, c, ca.CaCert, &k, ca.CaPrivKey)
	if err != nil {
		return []byte{}, err
	}
	return certBytes, nil
}

func getSubjectName(sc SubjectSchema) pkix.Name {
	return pkix.Name{
		Organization:  []string{sc.Organization},
		Country:       []string{sc.Country},
		Province:      []string{sc.Province},
		Locality:      []string{sc.Locality},
		StreetAddress: []string{sc.Address},
		PostalCode:    []string{sc.PostalCode},
		CommonName:    sc.CommonName,
	}
}

func SavePemFile(certype string, content []byte, filename string, perm os.FileMode) error {
	PEM := new(bytes.Buffer)
	pem.Encode(PEM, &pem.Block{
		Type:  certype,
		Bytes: content,
	})

	return ioutil.WriteFile(filename, PEM.Bytes(), 0600)
}

func InitCACertificate(f CertFilePath, confFiles string) (CertificateAuthority, error) {
	log.Printf("Creating CA Certificate %s with %s", f.CaCertPath, confFiles)

	config := CertConfigSchema{}
	CertConfigParser(confFiles, &config)
	fmt.Println(config)

	NotBefore, NotAfter := internal.TimeRange(config.CA.Valid)

	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(config.CA.SerialNumber),
		Subject:               getSubjectName(config.CA.Subject),
		NotBefore:             NotBefore,
		NotAfter:              NotAfter,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, config.CA.KeyLength)
	if err != nil {
		return CertificateAuthority{}, err
	}
	err = SavePemFile("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(caPrivKey), f.CaPrivKeyPath, 0600)
	if err != nil {
		return CertificateAuthority{}, err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return CertificateAuthority{}, err
	}
	err = SavePemFile("CERTIFICATE", caBytes, f.CaCertPath, 0600)
	if err != nil {
		return CertificateAuthority{}, err
	}

	return CertificateAuthority{CaCert: ca, CaPrivKey: caPrivKey}, nil
}

func CreateCertificate(ca CertificateAuthority, f CertFilePath, confFiles string) error {
	log.Printf("Creating Server Certificate %s with %s", f.ServerCertPath, confFiles)

	config := CertConfigSchema{}
	CertConfigParser(confFiles, &config)
	fmt.Println(config)

	NotBefore, NotAfter := internal.TimeRange(config.Certificate.Valid)

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(config.Certificate.SerialNumber),
		Subject:      getSubjectName(config.Certificate.Subject),
		DNSNames:     []string{config.Certificate.DNSName, "localhost"},
		NotBefore:    NotBefore,
		NotAfter:     NotAfter,
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certPrivKey, err := rsa.GenerateKey(rand.Reader, config.Certificate.KeyLength)
	if err != nil {
		return err
	}
	err = SavePemFile("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(certPrivKey), f.ServerPrivKeyPath, 0600)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.CaCert, &certPrivKey.PublicKey, ca.CaPrivKey)
	if err != nil {
		return err
	}
	err = SavePemFile("CERTIFICATE", certBytes, f.ServerCertPath, 0600)
	if err != nil {
		return err
	}

	return nil
}
