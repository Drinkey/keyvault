package model

import (
	"database/sql"
	"log"
	"os"
)

const KEY_DB_SCHEMA = `
CREATE TABLE vault (
    id integer not null primary key,
    namespace text not null,
    key text,
    value text
    );
`
const MASTER_KEY_SCHEMA = `
CREATE TABLE masterkey (
    id integer not null primary key,
    namespace text not null,
    master_key text
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
		_, err = conn.Exec(KEY_DB_SCHEMA)
		if err != nil {
			log.Fatalf("db %q: %s\n", err, KEY_DB_SCHEMA)
		}
		log.Printf("Table %s created", DBPATH)

		_, err = conn.Exec(MASTER_KEY_SCHEMA)
		if err != nil {
			log.Fatalf("db %q: %s\n", err, MASTER_KEY_SCHEMA)
		}
		log.Printf("Table masterkey created")
	}

}
