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

func getRouter() *gin.Engine {
	router := gin.Default()

	return router
}

// @title Keyvault API Document
// @version 1.0
// @description Keyvault API Document

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

func main() {
	r := routers.InitRouter()

	certPaths := certio.Cfg.Paths

	tlsConfig := certio.BuildTLSConfig(certPaths)

	httpServer := &http.Server{
		Addr:      ":443",
		Handler:   r,
		TLSConfig: tlsConfig,
	}
	g.Go(func() error {
		err := httpServer.ListenAndServeTLS(certPaths.WebCertPath, certPaths.WebPrivKeyPath)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	apiServerHandler := routers.InitAPIRouter()
	apiServer := &http.Server{
		Addr:    ":8080",
		Handler: apiServerHandler,
	}
	g.Go(func() error {
		err := apiServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
