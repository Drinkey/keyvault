package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/models"
	"github.com/Drinkey/keyvault/pkg/app"
	"github.com/Drinkey/keyvault/pkg/e"
	"github.com/Drinkey/keyvault/services/namespace_service"
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
	var (
		app = app.KvContext{Context: c}
		req Namespace
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS,
			fmt.Sprintf("The POST payload is invalid: %s", err.Error()))
		return
	}

	if !IsClientAuthorized(c.Request, req.Name) {
		msg := fmt.Sprintf(`Mismatched OU in Certificate and Namespace`)
		app.Response(http.StatusUnauthorized, e.NOT_AUTHORIED, msg)
		return
	}

	namespaceService := namespace_service.Namespace{
		Name: req.Name,
	}

	err := namespaceService.Create()
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR,
			fmt.Sprintf("Error when creating new namespace %s: %s", req.Name, err.Error()),
		)
		return
	}

	newNs, err := namespaceService.Get()
	if err != nil {
		app.Response(http.StatusNotFound, e.NOT_FOUND,
			fmt.Sprintf("Error when getting new namespace %s just created: %s", req.Name, err.Error()),
		)
		return
	}
	app.Response(http.StatusCreated, e.SUCCESS, newNs)
}
