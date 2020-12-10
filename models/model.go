package models

import (
	"log"

	"github.com/Drinkey/keyvault/internal"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// type Secret struct {
// 	ID          uint `gorm:"primaryKey;autoIncrement"`
// 	Key         string
// 	Value       string
// 	NamespaceID uint
// 	Namespace   Namespace `gorm:"foreignKey:NamespaceID"`
// }

// type Namspace struct {
// 	ID        uint   `gorm:"primaryKey;autoIncrement"`
// 	Name      string `gorm:"unique;not null"`
// 	MasterKey string
// 	Nonce     string
// }

// type Certificate struct {
// 	ID          uint   `gorm:"primaryKey;autoIncrement"`
// 	Name        string `gorm:"unique;not null"`
// 	SignRequest string
// 	Certificate string
// 	Token       string `gorm:"unique;not null"`
// }

const SECRET_DB_SCHEMA = `
CREATE TABLE secrets (
    id INTEGER PRIMARY KEY,
    key TEXT,
    value TEXT,
    namespace_id INTEGER NOT NULL,
	FOREIGN KEY(namespace_id) REFERENCES namespace(namespace_id),
	UNIQUE (namespace_id, key)
    );
`
const NS_DB_SCHEMA = `
CREATE TABLE namespace (
    namespace_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
	master_key TEXT,
	nonce TEXT,
	CONSTRAINT uniqueName UNIQUE(name)
    );
`

const CERT_DB_SCHEMA = `
CREATE TABLE certificate (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	namespace_id INTEGER NOT NULL,
	csr TEXT,
	cert TEXT,
	token TEXT NOT NULL,
)
`

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
	// gorm.Open(sqlite.Open(DBPATH), &gorm.Config{})

	connDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	db = connDB

	// if err = db.Ping(); err != nil {
	// 	log.Fatal("database unreachable:", err)
	// }

	if initDbRequired {
		log.Println("first install, initializing database schema")
		db.AutoMigrate(&Secret{})
		db.AutoMigrate(&Namespace{})
		db.AutoMigrate(&Certificate{})
		// tables := map[string]interface{
		// 	"namespace":   Namespace{},
		// 	"certificate": Certificate{},
		// 	"secrets":     Secret{},
		// }
		// for table, obj := range tables {
		// 	db.AutoMigrate(&obj)
		// 	// _, err = conn.Exec(sql)
		// 	// if err != nil {
		// 	// 	log.Fatalf("Create table [%s] failed: %q: %s\n", table, err, sql)
		// 	// }
		// 	log.Printf("Create table %s successful", table)
		// }
	}
}
