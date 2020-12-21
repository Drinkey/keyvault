package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/models"
	"github.com/Drinkey/keyvault/pkg/crypt"
	"github.com/Drinkey/keyvault/pkg/e"
	"github.com/gin-gonic/gin"
)

type Namespace struct {
	ID        uint   `json:"namespace_id"`
	Name      string `json:"name"`
	MasterKey string `json:"master_key"`
	Nonce     string `json:"nonce"`
}

func ListNamespaces(c *gin.Context) {
	// var ns models.Namespace
	r, err := models.ListNamespace()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  e.GetMsg(e.ERROR),
			"data": fmt.Sprintf("Error when retrieving namespace: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": r,
	})
}

func CreateNamespace(c *gin.Context) {
	log.Printf("Creating new namespace")

	var req Namespace
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  e.GetMsg(e.INVALID_PARAMS),
			"data": fmt.Sprintf("The POST payload is invalid: %s", err.Error()),
		})
		return
	}

	if !IsClientAuthorized(c.Request, req.Name) {
		msg := fmt.Sprintf(`Client not authorized to create namespace=%s.
		Cert OU and Namespace must be the same`, req.Name)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": e.NOT_AUTHORIED,
			"msg":  e.GetMsg(e.NOT_AUTHORIED),
			"data": msg,
		})
		return
	}

	req.MasterKey = crypt.EncodeByte(crypt.GenerateMasterKey())
	req.Nonce = crypt.EncodeByte(crypt.GenerateNonce())

	// var data models.Namespace
	err := models.CreateNamespace(req.Name, req.MasterKey, req.Nonce)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": e.ERROR,
			"msg":  e.GetMsg(e.ERROR),
			"data": fmt.Sprintf("Error when creating new namespace %s: %s", req.Name, err.Error()),
		})
		return
	}
	// Get the record just saved and mask sensitive data
	newNs, err := models.GetNamespace(req.Name)
	newNs.MasterKey = crypt.KeyMask
	newNs.Nonce = crypt.KeyMask
	c.JSON(http.StatusCreated, MakeResponse(e.SUCCESS, newNs))
}
