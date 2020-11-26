package controller

import (
	"fmt"
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
	var secret model.Secrets
	if err := c.ShouldBindJSON(&secret); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// var db model.Secrets
	// newSecret := db.Create(secret)
	// c.JSON(http.StatusCreated, gin.H{
	// 	"message": fmt.Sprintf("secret %s with value %s created success", name, value),
	// })
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
