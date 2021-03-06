package models

import (
	"log"

	"github.com/Drinkey/keyvault/pkg/settings"
	"github.com/Drinkey/keyvault/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	log.SetPrefix("model - ")
	dbPath := settings.Settings.DatabasePath

	dsn := dbPath

	log.Printf("initializing database: %s", dbPath)

	if dsn == "" {
		log.Print("Specify the DB Path in environment variable DB_PATH")
		log.Print("Without DB file, run in-memory")
		dsn = "file::memory:?cache=shared"
	}

	initDbRequired := true

	if utils.FileExist(dbPath) {
		log.Printf("database file %s already exist", dbPath)
		initDbRequired = false
	}

	connDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	db = connDB

	if initDbRequired {
		log.Println("first install, initializing database schema")
		db.AutoMigrate(&Secret{})
		db.AutoMigrate(&Namespace{})
		db.AutoMigrate(&Certificate{})
	}
}
