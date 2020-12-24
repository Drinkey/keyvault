package main

import (
	"log"

	"github.com/Drinkey/keyvault/pkg/server"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

// @title Keyvault API Document
// @version 1.0
// @description Keyvault API Document

// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

func main() {
	g.Go(server.CreateSecretServer)
	g.Go(server.CreateDefaultServer)
	g.Go(server.CreateHTTPServer)
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
