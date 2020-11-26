package model

import (
	"fmt"
	"log"
)

type Namespace struct {
	ID        int    `json:"namespace_id"`
	Name      string `json:"name"`
	MasterKey string `json:"master_key"`
}

func (ns Namespace) IsEmpty() bool {
	return ns.ID == 0
}

func (ns Namespace) IsNameExist(n string) bool {
	return ns.Name == n
}

func (ns Namespace) Get(name string) Namespace {
	sqlTemplate := `
	SELECT namespace_id, name, master_key
	FROM namespace
	WHERE name="%s"
	`
	sqlStmt := fmt.Sprintf(sqlTemplate, name)

	rows, err := conn.Query(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	var r Namespace
	for rows.Next() {
		err = rows.Scan(&r.ID, &r.Name, &r.MasterKey)
		if err != nil {
			log.Fatal(err)
		}
	}
	return r
}
