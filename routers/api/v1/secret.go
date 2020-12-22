package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/pkg/app"
	"github.com/Drinkey/keyvault/pkg/e"
	"github.com/Drinkey/keyvault/services/secret_service"
	"github.com/gin-gonic/gin"
)

type Secret struct {
	ID          uint      `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	NamespaceID uint      `json:"namespace_id"`
	Namespace   Namespace `json:"namespace"`
}

func GetSecret(c *gin.Context) {
	app := app.KvContext{Context: c}

	namespace := c.Param("namespace")

	if !app.IsClientAuthorized(namespace) {
		msg := fmt.Sprintf(`Client not authorized to access the namespace=%s`, namespace)
		app.Response(http.StatusUnauthorized, e.NOT_AUTHORIED, msg)
		return
	}

	key := c.Query("q")
	log.Printf("Query secret [%s] under namespace %s", key, namespace)
	var ss secret_service.Secret
	data, err := ss.Get(key, namespace)
	if err != nil {
		msg := fmt.Sprintf("error when finding secret: NameSpace=%s, Key=%s", namespace, key)
		app.Response(http.StatusNotFound, e.NOT_FOUND, msg)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, data)
}

// CreateSecret create a secret record in database, sensitive fields of info
// are encrypted or hashed.
func CreateSecret(c *gin.Context) {
	app := app.KvContext{Context: c}

	namespace := c.Param("namespace")

	if !app.IsClientAuthorized(namespace) {
		msg := fmt.Sprintf("Client not authorized to create new key under the namespace=%s", namespace)
		app.Response(http.StatusUnauthorized, e.NOT_AUTHORIED, msg)
		return
	}

	log.Printf("Creating new secret under %s", namespace)

	var req Secret
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, err.Error())
		return
	}
	var ss secret_service.Secret
	err := ss.Create(namespace, req.Key, req.Value)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR,
			fmt.Sprintf("Error when creating secret: %s", err.Error()))
		return
	}
	app.Response(http.StatusCreated, e.SUCCESS,
		fmt.Sprintf("Secret namespace=%s, key=%s created success", namespace, req.Key))
}
