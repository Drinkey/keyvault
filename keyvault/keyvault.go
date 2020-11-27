package main

import (
	"github.com/Drinkey/keyvault/controller"
	"github.com/Drinkey/keyvault/controller/namespace"
	"github.com/Drinkey/keyvault/controller/secret"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/ping", controller.Ping)
		// only list namespace of vault
		v1.GET("/vault", namespace.List)
		v1.POST("/vault", namespace.Create)
		// respond to URL like /vault/gitlab/?q=k8s_password
		v1.GET("/vault/:namespace", secret.Query)
		v1.POST("/vault/:namespace", secret.Create)
		// TODO:
		v1.DELETE("/vault/:namespace", secret.Delete)
		// TODO:
		v1.PUT("/vault/:namespace", secret.Update)
		// another user story
		// TODO:
		v1.POST("/certs/sign", controller.SignCSR)
	}
	return router
}

func main() {
	// fmt.Println(certio.DoNothing())
	// log.SetPrefix("main: ")
	router := setupRouter()
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
