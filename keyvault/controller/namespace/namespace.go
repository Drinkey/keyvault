package namespace

import (
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/internal"
	"github.com/Drinkey/keyvault/model"
	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	var ns model.Namespace
	r := ns.List()
	c.JSON(http.StatusOK, gin.H{"namespace": r})
}
