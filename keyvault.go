package main

import (
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/routers"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"
)

var g errgroup.Group
var certPaths certio.CertFilePaths

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

func createSecretServer() error {
	return createHTTPSServer("SECRET", ":443")
}

func createDefaultServer() error {
	return createHTTPSServer("MAINTENANCE", ":1443")
}

func createHTTPServer() error {
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

// @title Keyvault API Document
// @version 1.0
// @description Keyvault API Document

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

func main() {
	certPaths = certio.Cfg.Paths

	g.Go(createSecretServer)
	g.Go(createDefaultServer)
	g.Go(createHTTPServer)

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
