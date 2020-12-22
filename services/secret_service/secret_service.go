package secret_service

import (
	"fmt"
	"log"

	"github.com/Drinkey/keyvault/models"
	"github.com/Drinkey/keyvault/pkg/crypt"
)

type Secret struct {
}

func (s Secret) Get(key, namespace string) (*models.Secret, error) {
	secret, err := models.GetSecret(crypt.Sha256Sum(key), namespace)

	if err != nil {
		return nil, err
	}

	if secret.IsEmpty() {
		return nil, fmt.Errorf("Record Not Found: NameSpace=%s, Key=%s", namespace, key)
	}

	cipherTextBytes, err := crypt.DecodeString(secret.Value)
	if err != nil {
		log.Printf("failed to decode string %s", secret.Value)
		return nil, fmt.Errorf("failed to decode secret value string. NameSpace=%s, Key=%s", namespace, key)
	}
	nonceByte, err := crypt.DecodeString(secret.Namespace.Nonce)
	masterKeyByte, err := crypt.DecodeString(secret.Namespace.MasterKey)

	secret.Value = crypt.Decrypt(cipherTextBytes, masterKeyByte, nonceByte)
	secret.Namespace.MasterKey = crypt.KeyMask
	secret.Namespace.Nonce = crypt.KeyMask

	return &secret, nil
}

func (s Secret) Create(namespace, key, value string) error {
	// 1. Query namespace database for id, masterkey
	ns, err := models.GetNamespace(namespace)

	if err != nil {
		return fmt.Errorf("Error when finding namespace %s: %s", namespace, err.Error())
	}
	if ns.IsEmpty() {
		return fmt.Errorf("namespace %s does not exist, create it first", namespace)
	}

	// 2. encrypt secret value with master key
	nonceByte, err := crypt.DecodeString(ns.Nonce)
	masterKeyByte, err := crypt.DecodeString(ns.MasterKey)

	encryptDataBase64 := crypt.EncodeByte(
		crypt.Encrypt(value, masterKeyByte, nonceByte),
	)
	hashKey := crypt.Sha256Sum(key)

	// 3. save to database
	err = models.CreateSecret(hashKey, encryptDataBase64, namespace)
	if err != nil {
		return fmt.Errorf("Error when creating secret: %s", err.Error())
	}
	return nil
}
