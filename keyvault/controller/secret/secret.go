package secret

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/controller"
	"github.com/Drinkey/keyvault/internal"
	"github.com/Drinkey/keyvault/model"
	"github.com/gin-gonic/gin"
)

func Query(c *gin.Context) {

	namespace := c.Param("namespace")

	if !controller.IsClientAuthorized(c.Request, namespace) {
		msg := fmt.Sprintf(`Client not authorized to access the namespace=%s`, namespace)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": msg,
		})
		return
	}

	key := c.Query("q")
	log.Printf("Query secret [%s] under namespace %s", key, namespace)

	var secret_model model.Secrets
	secret := secret_model.Get(internal.Sha256Sum(key), namespace)

	if secret.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Record Not Found: NameSpace=%s, Key=%s", namespace, key),
		})
		return
	}

	cipherTextBytes, err := internal.DecodeString(secret.Value)
	if err != nil {
		log.Printf("failed to decode string %s", secret.Value)
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("failed to decode secret value string. NameSpace=%s, Key=%s", namespace, key),
		})
		return
	}
	nonceByte, err := internal.DecodeString(secret.NameSpace.Nonce)

	secret.Value = internal.Decrypt(cipherTextBytes, secret.NameSpace.MasterKey, nonceByte)
	secret.NameSpace.MasterKey = internal.KeyMask
	secret.NameSpace.Nonce = internal.KeyMask

	c.JSON(http.StatusOK, secret)
}

func Create(c *gin.Context) {

	namespace := c.Param("namespace")

	if !controller.IsClientAuthorized(c.Request, namespace) {
		msg := fmt.Sprintf("Client not authorized to create new key under the namespace=%s", namespace)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": msg,
		})
		return
	}

	log.Printf("Creating new secret under %s", namespace)

	var secret_data model.Secrets
	if err := c.ShouldBindJSON(&secret_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 1. Query namespace database for id, masterkey
	var ns_model model.Namespace
	ns, err := ns_model.Get(namespace)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ns.IsEmpty() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("namespace %s does not exist, create it first", namespace),
		})
		return
	}
	secret_data.NameSpace = ns

	// 2. encrypt secret value with master key
	nonceByte, err := internal.DecodeString(ns.Nonce)
	secret_data.Value = internal.EncodeByte(internal.Encrypt(secret_data.Value, ns.MasterKey, nonceByte))
	secret_data.Key = internal.Sha256Sum(secret_data.Key)

	// 3. save to database
	var secret_model model.Secrets
	err = secret_model.Create(secret_data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Secret namespace=%s, key=%s created success", namespace, secret_data.Key),
	})
}

func Delete(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("delete"),
	})
}

func Update(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("update"),
	})
}
