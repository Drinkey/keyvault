package namespace_service

import (
	"github.com/Drinkey/keyvault/models"
	"github.com/Drinkey/keyvault/pkg/crypt"
)

type Namespace struct {
	ID        uint   `json:"namespace_id"`
	Name      string `json:"name"`
	MasterKey string `json:"master_key"`
	Nonce     string `json:"nonce"`
}

func (n Namespace) Create(name string) error {
	return models.CreateNamespace(
		name,
		crypt.EncodeByte(crypt.GenerateMasterKey()),
		crypt.EncodeByte(crypt.GenerateNonce()),
	)
}

func (n Namespace) Get(name string) (*models.Namespace, error) {
	newNs, err := models.GetNamespace(name)
	if err != nil {
		return nil, err
	}
	newNs.MasterKey = crypt.KeyMask
	newNs.Nonce = crypt.KeyMask
	return &newNs, nil
}
