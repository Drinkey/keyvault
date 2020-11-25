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
		v1.GET("/vault", controller.ListVault)
		v1.GET("/vault/:id", controller.QuerySecret)
		v1.POST("/vault/:id", controller.CreateSecret)
		v1.DELETE("/vault/:id", controller.DeleteSecret)
		v1.PUT("/vault/:id", controller.UpdateSecret)
		v1.POST("/certs/sign", controller.SignCSR)
	}
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
