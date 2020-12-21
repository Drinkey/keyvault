package v1

import (
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/pkg/e"
	"github.com/gin-gonic/gin"
)

type KvResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func MakeResponse(code int, data interface{}) (r KvResponse) {
	return KvResponse{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: data,
	}
}

// Ping responses pong to the client. It can be used for service probing
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, MakeResponse(e.SUCCESS, "PONG"))
}

/*
Description:
Validate if namespace and OU in client certificate are the same.
Authorized client can only access the namespace OU specified

For example, a client with OU=GITLAB in certificate trying to access namespace
GITLAB, then the client is a authorized client. This is just a protocol.

Return true if namespace and OU are exactly the same.
*/

// IsClientAuthorized validates whether the current request authorized.
func IsClientAuthorized(r *http.Request, namespace string) bool {
	certOU, tlsEnabled := certio.ParseClientCertOU(r)
	if tlsEnabled && certOU != namespace {
		log.Printf("OU=%s and Namespace=%s should be the same", certOU, namespace)
		return false
	}
	return true
}
