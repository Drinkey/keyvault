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

func (n Namespace) Create() error {
	return models.CreateNamespace(
		n.Name,
		crypt.EncodeByte(crypt.GenerateMasterKey()),
		crypt.EncodeByte(crypt.GenerateNonce()),
	)
}

func (n Namespace) Get() (*models.Namespace, error) {
	newNs, err := models.GetNamespace(n.Name)
	if err != nil {
		return nil, err
	}
	newNs.MasterKey = crypt.KeyMask
	newNs.Nonce = crypt.KeyMask
	return &newNs, nil
}
