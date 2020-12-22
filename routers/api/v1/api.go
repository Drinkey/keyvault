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
// @Summary Response to service probing
// @Description probing
// @Produce json
// @Success 200 {string} string "ok" "PONG"
// @Router /api/v1/ping [get]
func Ping(c *gin.Context) {
	app := app.KvContext{Context: c}
	app.Response(http.StatusOK, e.SUCCESS, "PONG")
}
