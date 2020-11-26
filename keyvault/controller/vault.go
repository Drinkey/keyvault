package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Secrets struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

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
	c.JSON(http.StatusOK, gin.H{
		"key":   key,
		"value": "some_secrets",
	})
}

func CreateSecret(c *gin.Context) {
	var secret Secrets
	if err := c.ShouldBindJSON(&secret); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := secret.Name
	value := secret.Value
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("secret %s with value %s created success", name, value),
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
