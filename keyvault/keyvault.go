package main

import (
	"github.com/Drinkey/keyvault/controller"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// fmt.Println(certio.DoNothing())
	// log.SetPrefix("main: ")
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/ping", controller.Ping)
		// only list namespace of vault
		v1.GET("/vault", controller.ListVault)
		// respond to URL like /vault/gitlab/?q=k8s_password
		v1.GET("/vault/:namespace", controller.QuerySecret)
		v1.POST("/vault/:namespace", controller.CreateSecret)
		v1.DELETE("/vault/:namespace", controller.DeleteSecret)
		v1.PUT("/vault/:namespace", controller.UpdateSecret)
		// another user story
		v1.POST("/certs/sign", controller.SignCSR)
	}
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
