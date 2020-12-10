package v1

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/internal"
	"github.com/gin-gonic/gin"
)

func CreateCertificate(c *gin.Context) {

}

func createCertificateTemplate(csr *x509.CertificateRequest) *x509.Certificate {
	NotBefore, NotAfter := internal.TimeRange(5)
	return &x509.Certificate{
		SerialNumber:       big.NewInt(210201),
		Subject:            csr.Subject,
		NotBefore:          NotBefore,
		NotAfter:           NotAfter,
		Signature:          csr.Signature,
		SignatureAlgorithm: csr.SignatureAlgorithm,
		PublicKeyAlgorithm: csr.PublicKeyAlgorithm,
		PublicKey:          csr.PublicKey,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:           x509.KeyUsageDigitalSignature,
	}
}

func IssueCertificate(c *gin.Context) {
	// 1. Parse CSR
	// 2. Load CA
	// 3. Sign the CSR
	// 4. Construct response payload and return
	var csrPEM certio.CertificateSigningRequest
	if err := c.ShouldBindJSON(&csrPEM); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(csrPEM)
	block, _ := pem.Decode([]byte(csrPEM.Request))

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		e := fmt.Sprintf("Failed to parse Certificate Request in CSR: %s", err.Error())
		log.Printf(e)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": e,
		})
		return
	}

	if err = csr.CheckSignature(); err != nil {
		e := fmt.Sprintf("check csr signature failed: %s", err.Error())
		log.Printf(e)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": e,
		})
		return
	}
	log.Printf("check CSR signature passed")

	certTemplate := createCertificateTemplate(csr)

	ca, err := certio.LoadCACertificate(certio.CertFiles)
	if err != nil {
		log.Fatal("CA cert exists but failed to load it. Please delete the CA cert and re-generate it")
	}

	certBytes, err := certio.Issue(certTemplate, ca, csr.PublicKey)
	if err != nil {
		e := fmt.Sprintf("Failed to issue certificate: %s", err.Error())
		log.Printf(e)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": e,
		})
		return
	}
	log.Print("Certificate issued, creating response")

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca.CaCertBytes,
	})

	var response certio.CertificateResponse
	response.Certificate = certPEM.String()
	response.CA = caPEM.String()

	c.JSON(http.StatusCreated, response)
}

func GetCertificate(c *gin.Context) {

}
