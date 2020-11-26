package model

import (
	"fmt"
	"log"
)

type Secrets struct {
	ID        int    `json:"id"`
	NameSpace string `json:"namespace"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

// type SecretsDB interface{}

func (s Secrets) IsEmpty() bool {
	return s.ID == 0
}

func (s Secrets) Get(q string, ns string) Secrets {
	sqlTemplate := `
	SELECT s.id AS id, s.key AS key, s.value AS value, ns.name AS namespace
	FROM secrets AS s, namespace AS ns
	WHERE s.namespace_id = ns.namespace_id
	AND ns.name="%s"
	AND s.key="%s"`
	sqlStmt := fmt.Sprintf(sqlTemplate, ns, q)
	rows, err := conn.Query(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	var r Secrets
	for rows.Next() {
		err = rows.Scan(&s.ID, &s.Key, &s.Value, &s.NameSpace)
		if err != nil {
			log.Fatal(err)
		}
	}
	return r
}

func (d Secrets) Create(s Secrets) {

}

func Delete() {

}

func Update() {

}
