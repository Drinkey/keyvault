package main

import (
	"net/http"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/routers"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func getRouter() *gin.Engine {
	router := gin.Default()

	return router
}

func main() {
	r := routers.InitRouter()

	tlsConfig := certio.BuildTLSConfig()

	httpServer := &http.Server{
		Addr:      ":443",
		Handler:   r,
		TLSConfig: tlsConfig,
	}
	httpServer.ListenAndServeTLS(certio.CertFiles.ServerCertPath, certio.CertFiles.ServerPrivKeyPath)
}
