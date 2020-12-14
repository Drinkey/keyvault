/*
Package certio provides all operations against certificate.
*/
package certio

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/Drinkey/keyvault/internal"
)

func LoadCACertificate(c CertFilePath) (CertificateAuthority, error) {
	caTxt, err := ioutil.ReadFile(c.CaCertPath)
	if err != nil {
		log.Fatal("read certificate file failed", err)
	}
	caBlock, _ := pem.Decode([]byte(caTxt))
	caCert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		log.Fatal("parse certificate content failed", err)
	}

	caPrivKeyTxt, err := ioutil.ReadFile(c.CaPrivKeyPath)
	if err != nil {
		log.Fatal("read private key file failed", err)
	}
	caPrivKeyBlock, _ := pem.Decode([]byte(caPrivKeyTxt))
	caPrivKey, err := x509.ParsePKCS1PrivateKey(caPrivKeyBlock.Bytes)
	if err != nil {
		log.Fatal("parse private key content failed", err)
	}

	return CertificateAuthority{CaCert: caCert, CaCertBytes: caBlock.Bytes, CaPrivKey: caPrivKey}, nil
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
