package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/models"
	"github.com/Drinkey/keyvault/pkg/app"
	"github.com/Drinkey/keyvault/pkg/crypt"
	"github.com/Drinkey/keyvault/pkg/e"
	"github.com/Drinkey/keyvault/services/certificate_service"
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
		data := fmt.Sprintf("The POST payload is invalid: %s", err.Error())
		c.JSON(http.StatusBadRequest, MakeResponse(e.INVALID_PARAMS, data))
		return
	}
	err := models.CreateCertificateRequest(
		req.Name,
		req.SignRequest,
		crypt.EncodeByte(crypt.GenerateRandomKey(20)),
	)
	if err != nil {
		data := fmt.Sprintf("Error when creating record: %s", err.Error())
		c.JSON(http.StatusInternalServerError, MakeResponse(e.ERROR, data))
		return
	}
	newCert, err := models.GetCertificate(req.Name)
	log.Print("got from db:")
	log.Println(newCert)
	if err != nil {
		data := fmt.Sprintf("Error when retrieving record: %s", err.Error())
		c.JSON(http.StatusInternalServerError, MakeResponse(e.ERROR, data))
		return
	}
	newCert.SignRequest = crypt.KeyMask
	c.JSON(http.StatusCreated, MakeResponse(e.SUCCESS, newCert))
}

func IssueCertificate(c *gin.Context) {
	// 1. Parse CSR
	// 2. Load CA
	// 3. Sign the CSR
	// 4. Construct response payload and return
	var csrPEM certio.CertificateRequest
	if err := c.ShouldBindJSON(&csrPEM); err != nil {
		data := fmt.Sprintf("The POST payload is invalid: %s", err.Error())
		c.JSON(http.StatusBadRequest, MakeResponse(e.INVALID_PARAMS, data))
		return
	}
	log.Println(csrPEM)
	// // pem string to certificate request
	// cert, ca := certio.IssueCertificate(csrPEM.Request)

	var r certio.CertificateResponse
	r.Certificate, r.CA = certio.IssueCertificate(csrPEM.Request)

	c.JSON(http.StatusCreated, MakeResponse(e.SUCCESS, r))
}

func GetCACertificate(c *gin.Context) {
	log.Print("reading CA certificate")
	r := CACertificate{Certificate: certio.CaContainer.String}
	c.JSON(http.StatusOK, MakeResponse(e.SUCCESS, r))
}

func GetCertificate(c *gin.Context) {
	app := app.KvContext{Context: c}
	q := c.Query("q")
	log.Printf("got query for %s", q)
	var cs certificate_service.Certificate
	cert, err := cs.Get(q)
	if err != nil {
		// data := fmt.Sprintf("failed to get certificate %s", q)
		app.Response(http.StatusNotFound, e.NOT_FOUND, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, cert)
}
