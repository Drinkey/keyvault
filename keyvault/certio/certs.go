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
	"net"

	"github.com/Drinkey/keyvault/internal"
)

func LoadCACertificate(certs CertFiles) (CertificateAuthority, error) {
	ca_txt, err := ioutil.ReadFile(certs.CaCert)
	if err != nil {
		log.Fatal("read certificate file failed", err)
	}
	ca_block, _ := pem.Decode([]byte(ca_txt))
	cacert, err := x509.ParseCertificate(ca_block.Bytes)
	if err != nil {
		log.Fatal("parse certificate content failed", err)
	}

	ca_pkey_txt, err := ioutil.ReadFile(certs.CaPrivKey)
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

func getSubjectName(sc SubjectConfig) pkix.Name {
	return pkix.Name{
		Organization:  []string{sc.Organization},
		Country:       []string{sc.Country},
		Province:      []string{sc.Province},
		Locality:      []string{sc.Locality},
		StreetAddress: []string{sc.Address},
		PostalCode:    []string{sc.PostalCode},
	}
}

func InitCACertificate(f CertFiles) (CertificateAuthority, error) {
	log.Printf("Creating CA Certificate %s with %s", f.CaCert, f.CaCertConf)

	config := CAConfig{}
	CaConfigParser(f.CaCertConf, &config)
	fmt.Println(config)

	NotBefore, NotAfter := internal.TimeRange(config.Valid)

	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(config.SerialNumber),
		Subject:               getSubjectName(config.Subject),
		NotBefore:             NotBefore,
		NotAfter:              NotAfter,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, config.KeyLength)
	if err != nil {
		return CertificateAuthority{}, err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return CertificateAuthority{}, err
	}

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	err = ioutil.WriteFile(f.CaCert, caPEM.Bytes(), 0600)
	if err != nil {
		return CertificateAuthority{}, err
	}

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	err = ioutil.WriteFile(f.CaPrivKey, caPrivKeyPEM.Bytes(), 0600)
	if err != nil {
		return CertificateAuthority{}, err
	}

	return CertificateAuthority{CaCert: ca, CaPrivKey: caPrivKey}, nil
}

func CreateCertificate(ca CertificateAuthority, f CertFiles) error {
	log.Printf("Creating Server Certificate %s with %s", f.ServerCert, f.ServerCertConf)

	config := CertConfig{}
	CertConfigParser(f.ServerCertConf, &config)
	fmt.Println(config)

	NotBefore, NotAfter := internal.TimeRange(config.Valid)

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(config.SerialNumber),
		Subject:      getSubjectName(config.Subject),
		IPAddresses:  []net.IP{net.ParseIP(config.IPv4Address), net.IPv6loopback},
		NotBefore:    NotBefore,
		NotAfter:     NotAfter,
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certPrivKey, err := rsa.GenerateKey(rand.Reader, config.KeyLength)
	if err != nil {
		return err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.CaCert, &certPrivKey.PublicKey, ca.CaPrivKey)
	if err != nil {
		return err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	err = ioutil.WriteFile(f.ServerCert, certPEM.Bytes(), 0600)
	if err != nil {
		return err
	}

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	err = ioutil.WriteFile(f.ServerPrivKey, certPrivKeyPEM.Bytes(), 0600)
	if err != nil {
		return err
	}
	return nil
}
