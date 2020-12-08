package namespace

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/controller"
	"github.com/Drinkey/keyvault/internal"
	"github.com/Drinkey/keyvault/model"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	var ns model.Namespace
	r := ns.List()
	c.JSON(http.StatusOK, gin.H{"namespace": r})
}

func Create(c *gin.Context) {
	log.Printf("Creating new namespace")

	var ns_data model.Namespace
	if err := c.ShouldBindJSON(&ns_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !controller.IsClientAuthorized(c.Request, ns_data.Name) {
		msg := fmt.Sprintf(`Client not authorized to create namespace=%s.
		Cert OU and Namespace must be the same`, ns_data.Name)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": msg,
		})
		return
	}

	ns_data.MasterKey = internal.EncodeByte(internal.GenerateMasterKey())
	ns_data.Nonce = internal.EncodeByte(internal.GenerateNonce())

	var ns_model model.Namespace
	err := ns_model.Create(ns_data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newNs, err := ns_model.Get(ns_data.Name)
	newNs.MasterKey = internal.KeyMask
	newNs.Nonce = internal.KeyMask
	c.JSON(http.StatusCreated, newNs)
}
