package v1

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Drinkey/keyvault/pkg/app"
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

// CreateCertificateRequest create a new CSR records in database
// @Summary create a new CSR records in database
// @Description create a new CSR records in database
// @Produce json
// @Tags Certificate
// @Param body body Certificate true "body"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "err"
// @Failure 401 {string} string "err"
// @Failure 500 {string} string "err"
// @Router /api/v1/cert [post]
func CreateCertificateRequest(c *gin.Context) {

	var (
		app = app.KvContext{Context: c}
		req Certificate
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		data := fmt.Sprintf("The POST payload is invalid: %s", err.Error())
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, data)
		return
	}

	var cs certificate_service.Certificate

	data, err := cs.Create(req.Name, req.SignRequest)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	app.Response(http.StatusCreated, e.SUCCESS, data)
}

// IssueCertificate issues certificate requests
// TODO:
// func IssueCertificate(c *gin.Context) {
// 	app := app.KvContext{Context: c}
// 	// 1. Parse CSR
// 	// 2. Load CA
// 	// 3. Sign the CSR
// 	// 4. Construct response payload and return
// 	var csrPEM certio.CertificateRequest
// 	if err := c.ShouldBindJSON(&csrPEM); err != nil {
// 		data := fmt.Sprintf("The POST payload is invalid: %s", err.Error())
// 		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, data)
// 		return
// 	}
// 	log.Println(csrPEM)
// 	// // pem string to certificate request
// 	// cert, ca := certio.IssueCertificate(csrPEM.Request)

// 	var r certio.CertificateResponse
// 	r.Certificate, r.CA = certio.IssueCertificate(csrPEM.Request)

// 	app.Response(http.StatusCreated, e.SUCCESS, r)
// }

func GetCertificate(c *gin.Context) {
	app := app.KvContext{Context: c}
	q := c.Query("q")
	log.Printf("got query for %s,validating query string", q)

	if strings.ToLower(q) == "ca" {
		log.Print("reading CA certificate")
		var cs certificate_service.Certificate
		r := CACertificate{Certificate: cs.GetCA()}
		app.Response(http.StatusOK, e.SUCCESS, r)
		return
	}

	var cs certificate_service.Certificate
	cert, err := cs.Get(q)
	if err != nil {
		data := fmt.Sprintf("failed to get certificate %s", q)
		app.Response(http.StatusNotFound, e.NOT_FOUND, data)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, cert)
}
