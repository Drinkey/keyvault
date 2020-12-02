package controller

import (
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/certio"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// Description:
// Validate if namespace and OU in client certificate are the same.
// Authorized client can only access the namespace OU specified
//
// For example, a client with OU=GITLAB in certificate trying to access namespace
// GITLAB, then the client is a authorized client. This is just a protocol.
//
// Return true if namespace and OU are exactly the same.
func IsClientAuthorized(r *http.Request, namespace string) bool {
	certOU, tlsEnabled := certio.ParseClientCertOU(r)
	if tlsEnabled && certOU != namespace {
		log.Printf("OU=%s and Namespace=%s should be the same", certOU, namespace)
		return false
	}
	return true
}
