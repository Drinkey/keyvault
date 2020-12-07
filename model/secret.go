package model

import (
	"fmt"
	"log"
)

type Secrets struct {
	ID        int       `json:"id"`
	NameSpace Namespace `json:"namespace"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
}

// type SecretsDB interface{}

func (s Secrets) IsEmpty() bool {
	return s.ID == 0
}

func (s Secrets) Get(q string, ns string) Secrets {
	sqlTemplate := `
	SELECT 
	s.id AS id,
	s.key AS key,
	s.value AS value,
	s.namespace_id as namespace_id,
	ns.name AS namespace,
	ns.master_key as encryption_key,
	ns.nonce as nonce
	FROM secrets AS s, namespace AS ns
	WHERE s.namespace_id = ns.namespace_id
	AND ns.name="%s"
	AND s.key="%s"`
	sqlStmt := fmt.Sprintf(sqlTemplate, ns, q)

	rows, err := conn.Query(sqlStmt)
	if err != nil {
		log.Println(err)
	}

	var r Secrets
	for rows.Next() {
		err = rows.Scan(&r.ID, &r.Key, &r.Value, &r.NameSpace.ID,
			&r.NameSpace.Name, &r.NameSpace.MasterKey, &r.NameSpace.Nonce)
		if err != nil {
			log.Println(err)
		}
	}
	return r
}

func (s Secrets) Create(secret Secrets) error {
	sqlTemplate := `
	INSERT INTO secrets
	(key, value, namespace_id)
	VALUES 
	("%s", "%s", (SELECT namespace_id FROM namespace WHERE name="%s"));
	`

	sqlStmt := fmt.Sprintf(sqlTemplate, secret.Key, secret.Value, secret.NameSpace.Name)
	_, err := conn.Exec(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Delete() {

}

func Update() {

}
