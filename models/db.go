package models

import (
	"log"
	"os"

	"github.com/Drinkey/keyvault/internal"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DBPATH string = os.Getenv("DB_PATH")

var db *gorm.DB

func init() {
	log.SetPrefix("model - ")

	dsn := DBPATH

	log.Printf("initializing database: %s", DBPATH)

	if dsn == "" {
		log.Print("Specify the DB Path in environment variable DB_PATH")
		log.Print("Without DB file, run in-memory")
		dsn = "file::memory:?cache=shared"
	}

	initDbRequired := true

	if internal.FileExist(DBPATH) {
		log.Printf("database file %s already exist", DBPATH)
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
