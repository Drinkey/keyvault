package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/internal"
	"github.com/Drinkey/keyvault/model"
	"github.com/gin-gonic/gin"
)

type Vault struct {
	Name string `json:"name"`
}

func ListVault(c *gin.Context) {
	vaults := []Vault{Vault{Name: "vault1"}, Vault{Name: "vault2"}}
	c.JSON(http.StatusOK, gin.H{
		"vaults": vaults,
	})
}

func QuerySecret(c *gin.Context) {
	namespace := c.Param("namespace")
	key := c.Query("q")

	var db model.Secrets
	secret := db.Get(key, namespace)

	if secret.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Record Not Found: NameSpace=%s, Key=%s", namespace, key),
		})
		return
	}

	secret.Value = internal.Decrypt(secret.Value, secret.NameSpace.MasterKey)
	secret.NameSpace.MasterKey = "******"

	c.JSON(http.StatusOK, secret)
}

func CreateSecret(c *gin.Context) {

	namespace := c.Param("namespace")
	log.Printf("Creating new secret under %s", namespace)

	var secret model.Secrets
	if err := c.ShouldBindJSON(&secret); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 1. Query namespace database for id, masterkey
	var ns model.Namespace
	secret.NameSpace = ns.Get(namespace)
	// 2. encrypt secret value with master key
	secret.Value = internal.Encrypt(secret.Value, secret.NameSpace.MasterKey)
	fmt.Println(secret)
	// 3. save to database
	var secret_model model.Secrets
	err := secret_model.Create(secret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Secret namespace=%s, key=%s created success", namespace, secret.Key),
	})
}

func DeleteSecret(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("delete"),
	})
}

func UpdateSecret(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("update"),
	})
}
