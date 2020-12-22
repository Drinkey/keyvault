package v1

import (
	"net/http"

	"github.com/Drinkey/keyvault/pkg/app"
	"github.com/Drinkey/keyvault/pkg/e"
	"github.com/gin-gonic/gin"
)

type KvResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Ping responses pong to the client. It can be used for service probing
func Ping(c *gin.Context) {
	app := app.KvContext{Context: c}
	app.Response(http.StatusOK, e.SUCCESS, "PONG")
}
