package app

import (
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
