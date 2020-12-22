package app

import (
	"log"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/pkg/e"
	"github.com/gin-gonic/gin"
)

type KvContext struct {
	Context *gin.Context `json:"-"`
	Code    int          `json:"code"`
	Msg     string       `json:"msg"`
	Data    interface{}  `json:"data"`
}

func (c KvContext) Response(status_code, err_code int, data interface{}) {
	c.Context.JSON(status_code, KvContext{
		Code: err_code,
		Msg:  e.GetMsg(err_code),
		Data: data,
	})
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
func (c KvContext) IsClientAuthorized(namespace string) bool {
	certOU, tlsEnabled := certio.ParseClientCertOU(c.Context.Request)
	if tlsEnabled && certOU != namespace {
		log.Printf("OU=%s and Namespace=%s should be the same", certOU, namespace)
		return false
	}
	return true
}
