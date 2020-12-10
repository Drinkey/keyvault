package models

import (
	"os"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

var DBPATH string = os.Getenv("DB_PATH")

var db *gorm.DB
