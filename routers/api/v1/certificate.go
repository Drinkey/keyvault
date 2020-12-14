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
	"github.com/Drinkey/keyvault/models"
	"github.com/gin-gonic/gin"
)

type Certificate struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	SignRequest string `json:"req"`
	Certificate string `json:"certificate"`
	Token       string `json:"token"`
}

func CreateCertificateRequest(c *gin.Context) {
	var req Certificate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := models.CreateCertificateRequest(
		req.Name,
		req.SignRequest,
		internal.EncodeByte(internal.GenerateRandomKey(20)),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newCert, err := models.GetCertificate(req.Name)
	log.Print("got from db:")
	log.Println(newCert)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newCert.SignRequest = internal.KeyMask
	c.JSON(http.StatusCreated, newCert)
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

	ca, err := certio.LoadCACertificate(certio.GetCertFiles())
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
	q := c.Query("q")
	cert, err := models.GetCertificate(q)
	if err != nil {
		log.Printf("get certificate %s error", q)
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("et certificate error: %s", err.Error()),
		})
		return
	}
	if cert.IsEmpty() {
		log.Printf("failed to get certificate %s, 0 result found", q)
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("failed to get certificate %s, 0 result found", q),
		})
		return
	}
	c.JSON(http.StatusOK, cert)
}
