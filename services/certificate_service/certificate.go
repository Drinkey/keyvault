package certificate_service

import (
	"github.com/Drinkey/keyvault/models"
)

type Certificate struct {
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
