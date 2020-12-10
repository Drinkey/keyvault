package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/internal"
	"github.com/Drinkey/keyvault/models"
	"github.com/gin-gonic/gin"
)

type Secret struct {
	ID          uint      `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	NamespaceID uint      `json:"namespace_id"`
	Namespace   Namespace `json:"namespace"`
}

func GetSecrets(c *gin.Context) {

	namespace := c.Param("namespace")

	if !IsClientAuthorized(c.Request, namespace) {
		msg := fmt.Sprintf(`Client not authorized to access the namespace=%s`, namespace)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": msg,
		})
		return
	}

	key := c.Query("q")
	log.Printf("Query secret [%s] under namespace %s", key, namespace)

	s, err := models.GetSecret(internal.Sha256Sum(key), namespace)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("error when finding secret: NameSpace=%s, Key=%s", namespace, key),
		})
		return
	}

	if s.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Record Not Found: NameSpace=%s, Key=%s", namespace, key),
		})
		return
	}

	cipherTextBytes, err := internal.DecodeString(s.Value)
	if err != nil {
		log.Printf("failed to decode string %s", s.Value)
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("failed to decode secret value string. NameSpace=%s, Key=%s", namespace, key),
		})
		return
	}
	nonceByte, err := internal.DecodeString(s.Namespace.Nonce)
	masterKeyByte, err := internal.DecodeString(s.Namespace.MasterKey)
	// var resp Secret
	s.Value = internal.Decrypt(cipherTextBytes, masterKeyByte, nonceByte)
	s.Namespace.MasterKey = internal.KeyMask
	s.Namespace.Nonce = internal.KeyMask

	c.JSON(http.StatusOK, s)
}

func CreateSecret(c *gin.Context) {

	namespace := c.Param("namespace")

	if !IsClientAuthorized(c.Request, namespace) {
		msg := fmt.Sprintf("Client not authorized to create new key under the namespace=%s", namespace)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": msg,
		})
		return
	}

	log.Printf("Creating new secret under %s", namespace)

	var req Secret
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 1. Query namespace database for id, masterkey
	ns, err := models.GetNamespace(namespace)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ns.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("namespace %s does not exist, create it first", namespace),
		})
		return
	}

	// 2. encrypt secret value with master key
	nonceByte, err := internal.DecodeString(ns.Nonce)
	masterKeyByte, err := internal.DecodeString(ns.MasterKey)

	encryptDataBase64 := internal.EncodeByte(
		internal.Encrypt(req.Value, masterKeyByte, nonceByte),
	)
	hashKey := internal.Sha256Sum(req.Key)

	// 3. save to database
	err = models.CreateSecret(hashKey, encryptDataBase64, namespace)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Secret namespace=%s, key=%s created success", namespace, req.Key),
	})
}
