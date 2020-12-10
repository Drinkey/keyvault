package routers

import (
	v1 "github.com/Drinkey/keyvault/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/ping", v1.Ping)
		// only list namespace of vault
		apiv1.GET("/namespace", v1.ListNamespace)
		apiv1.POST("/namespace", v1.CreateNamespace)
		// respond to URL like /vault/gitlab/?q=k8s_password
		apiv1.GET("/vault/:namespace", v1.GetSecrets)
		apiv1.POST("/vault/:namespace", v1.CreateSecrets)
		// TODO:
		// v1.DELETE("/vault/:namespace", v1.secret.Delete)
		// TODO:
		// v1.PUT("/vault/:namespace", v1.secret.Update)
		// another user story
		// TODO:
		// v1.POST("/cert/", v1.controller.SignCSR)
	}
	return r
}
