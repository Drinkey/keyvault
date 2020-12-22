package certificate_service

import (
	"log"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/models"
	"github.com/Drinkey/keyvault/pkg/crypt"
)

type Certificate struct {
}

func (c Certificate) GetCA() string {
	return certio.CaContainer.String
}
func (c Certificate) Get(name string) (*models.Certificate, error) {
	cert, err := models.GetCertificate(name)

	if err != nil {
		return nil, err
	}
	if cert.IsEmpty() {
		return nil, err
	}
	return &cert, nil
}

func (c Certificate) Create(name, request string) (*models.Certificate, error) {
	err := models.CreateCertificateRequest(
		name,
		request,
		crypt.EncodeByte(crypt.GenerateRandomKey(20)),
	)
	if err != nil {
		return nil, err
	}
	newCert, err := models.GetCertificate(name)
	log.Print("got from db:")
	log.Println(newCert)
	if err != nil {
		return nil, err
	}
	newCert.SignRequest = crypt.KeyMask
	return &newCert, nil
}
