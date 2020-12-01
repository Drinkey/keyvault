package namespace

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/certio"
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
	certOU, tlsEnabled := certio.ParseClientCertOU(c.Request)
	if tlsEnabled && certOU != ns_data.Name {
		log.Printf("OU=%s and Namespace=%s should be the same", certOU, ns_data.Name)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprintf("Not a authorized client to create namespace=%s. Cert OU and Namespace must be the same", ns_data.Name),
		})
		return
	}
	ns_data.MasterKey = internal.GenerateMasterKey()
	var ns_model model.Namespace
	err := ns_model.Create(ns_data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newNs, err := ns_model.Get(ns_data.Name)
	newNs.MasterKey = internal.KeyMask
	c.JSON(http.StatusCreated, newNs)
}
