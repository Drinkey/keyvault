package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/models"
	"github.com/Drinkey/keyvault/pkg/crypt"
	"github.com/gin-gonic/gin"
)

type Certificate struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	SignRequest string `json:"req"`
	Certificate string `json:"certificate"`
	Token       string `json:"token"`
}

type CACertificate struct {
	Certificate string `json:"ca"`
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
		crypt.EncodeByte(crypt.GenerateRandomKey(20)),
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
	newCert.SignRequest = crypt.KeyMask
	c.JSON(http.StatusCreated, newCert)
}

func IssueCertificate(c *gin.Context) {
	// 1. Parse CSR
	// 2. Load CA
	// 3. Sign the CSR
	// 4. Construct response payload and return
	var csrPEM certio.CertificateRequest
	if err := c.ShouldBindJSON(&csrPEM); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(csrPEM)
	// // pem string to certificate request
	// cert, ca := certio.IssueCertificate(csrPEM.Request)

	var r certio.CertificateResponse
	r.Certificate, r.CA = certio.IssueCertificate(csrPEM.Request)

	c.JSON(http.StatusCreated, r)
}

func GetCACertificate(c *gin.Context) {
	log.Print("reading CA certificate")
	r := CACertificate{Certificate: certio.CaContainer.String}
	c.JSON(http.StatusOK, r)
}

func GetCertificate(c *gin.Context) {
	q := c.Query("q")
	log.Printf("got query for %s", q)
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
