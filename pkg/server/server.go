package server

import (
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/routers"
	"github.com/gin-gonic/gin"
)

var certPaths = certio.Cfg.Paths

func createHTTPSServer(level, port string) error {
	var r *gin.Engine
	switch level {
	case "SECRET":
		r = routers.InitSecretRouter()
	case "MAINTENANCE":
		r = routers.InitDefaultRouter()
	}

	tlsConfig := certio.BuildTLSConfig(certPaths, level)
	secretServer := &http.Server{
		Addr:      port,
		Handler:   r,
		TLSConfig: tlsConfig,
	}
	err := secretServer.ListenAndServeTLS(certPaths.WebCertPath, certPaths.WebPrivKeyPath)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	return err
}

func CreateSecretServer() error {
	return createHTTPSServer("SECRET", ":443")
}

func CreateDefaultServer() error {
	return createHTTPSServer("MAINTENANCE", ":1443")
}

func CreateHTTPServer() error {
	apiServerHandler := routers.InitAPIRouter()
	apiServer := &http.Server{
		Addr:    ":8080",
		Handler: apiServerHandler,
	}
	err := apiServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	return err
}
