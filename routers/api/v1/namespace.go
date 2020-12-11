package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/internal"
	"github.com/Drinkey/keyvault/models"
	"github.com/gin-gonic/gin"
)

type Namespace struct {
	ID        uint   `json:"namespace_id"`
	Name      string `json:"name"`
	MasterKey string `json:"master_key"`
	Nonce     string `json:"nonce"`
}

func GetNamespaces(c *gin.Context) {

}

func ListNamespaces(c *gin.Context) {
	// var ns models.Namespace
	r, err := models.ListNamespace()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"namespace": r})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r})
}

func CreateNamespace(c *gin.Context) {
	log.Printf("Creating new namespace")

	var req Namespace
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !IsClientAuthorized(c.Request, req.Name) {
		msg := fmt.Sprintf(`Client not authorized to create namespace=%s.
		Cert OU and Namespace must be the same`, req.Name)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": msg,
		})
		return
	}

	req.MasterKey = internal.EncodeByte(internal.GenerateMasterKey())
	req.Nonce = internal.EncodeByte(internal.GenerateNonce())

	// var data models.Namespace
	err := models.CreateNamespace(req.Name, req.MasterKey, req.Nonce)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get the record just saved and mask sensitive data
	newNs, err := models.GetNamespace(req.Name)
	newNs.MasterKey = internal.KeyMask
	newNs.Nonce = internal.KeyMask
	c.JSON(http.StatusCreated, newNs)
}
