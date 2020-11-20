package internal

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
)

func CreateCACertificate() error {
	log.Printf("Creating CA Certificate %s with %s", CA_FILE, CA_CONF)
	configFile := fmt.Sprintf("%s/%s", CONF_DIR, CA_CONF)

	config := CAConfig{}
	CaConfigParser(configFile, &config)
	fmt.Println(config)

	NotBefore, NotAfter := TimeRange(config.Valid)

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(config.SerialNumber),
		Subject: pkix.Name{
			Organization:  []string{config.Subject.Organization},
			Country:       []string{config.Subject.Country},
			Province:      []string{config.Subject.Province},
			Locality:      []string{config.Subject.Locality},
			StreetAddress: []string{config.Subject.Address},
			PostalCode:    []string{config.Subject.PostalCode},
		},
		NotBefore:             NotBefore,
		NotAfter:              NotAfter,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, config.KeyLength)
	if err != nil {
		return err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return err
	}

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	caFile := fmt.Sprintf("%s/%s", CONF_DIR, CA_FILE)
	err = ioutil.WriteFile(caFile, caPEM.Bytes(), 0600)
	if err != nil {
		return err
	}

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	caPrivKeyFile := fmt.Sprintf("%s/ca_priv.key", CONF_DIR)
	err = ioutil.WriteFile(caPrivKeyFile, caPrivKeyPEM.Bytes(), 0600)
	if err != nil {
		return err
	}
	return nil
}

func CreateCertificate() error {
	log.Printf("Creating Certificate %s with %s", CERT_FILE, CERT_CONF)
	filename := fmt.Sprintf("%s/%s", CONF_DIR, CERT_CONF)
	data := CertConfig{}
	CertConfigParser(filename, &data)
	fmt.Println(data)
	return nil
}
