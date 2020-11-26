package model

import (
	"database/sql"
	"log"
	"os"

	"github.com/Drinkey/keyvault/internal"
	_ "github.com/mattn/go-sqlite3"
)

const SECRET_DB_SCHEMA = `
CREATE TABLE secrets (
    id INTEGER PRIMARY KEY,
    key TEXT,
    value TEXT,
    namespace_id INTEGER NOT NULL,
    FOREIGN KEY(namespace_id) REFERENCES namespace(namespace_id)
    );
`
const NS_DB_SCHEMA = `
CREATE TABLE namespace (
    namespace_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    master_key TEXT
    );
`

var DBPATH string = os.Getenv("DB_PATH")

var conn *sql.DB

func init() {
	log.SetPrefix("model - init():")
	log.Println("initializing database")

	if DBPATH == "" {
		log.Panic("Specify the DB Path in environment variable DB_PATH")
	}

	initDbRequired := true

	if internal.FileExist(DBPATH) {
		log.Println("database file already exist")
		initDbRequired = false
	}

	connDB, err := sql.Open("sqlite3", DBPATH)
	if err != nil {
		log.Panic(err)
	}
	conn = connDB

	if err = conn.Ping(); err != nil {
		log.Fatal("database unreachable:", err)
	}

	if initDbRequired {
		log.Println("first install, initializing database schema")
		_, err = conn.Exec(NS_DB_SCHEMA)
		if err != nil {
			log.Fatalf("db %q: %s\n", err, NS_DB_SCHEMA)
		}
		log.Printf("Table namespace created")

		_, err = conn.Exec(SECRET_DB_SCHEMA)
		if err != nil {
			log.Fatalf("db %q: %s\n", err, SECRET_DB_SCHEMA)
		}
		log.Printf("Table secret created")
	}

}
