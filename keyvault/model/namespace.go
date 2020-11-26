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

func (ns Namespace) Get(name string) (Namespace, error) {
	sqlTemplate := `
	SELECT namespace_id, name, master_key
	FROM namespace
	WHERE name="%s"
	`
	sqlStmt := fmt.Sprintf(sqlTemplate, name)

	rows, err := conn.Query(sqlStmt)
	if err != nil {
		log.Println(err)
		return ns, err
	}

	var r Namespace
	for rows.Next() {
		err = rows.Scan(&r.ID, &r.Name, &r.MasterKey)
		if err != nil {
			log.Println(err)
			return ns, err
		}
	}
	return r, nil
}

func (ns Namespace) List() []Namespace {
	sqlStmt := `
	SELECT namespace_id, name
	FROM namespace
	`
	rows, err := conn.Query(sqlStmt)
	if err != nil {
		log.Println(err)
	}

	var namespaces []Namespace
	for rows.Next() {
		var ns Namespace
		err = rows.Scan(&ns.ID, &ns.Name)
		if err != nil {
			log.Fatal(err)
		}
		ns.MasterKey = "******"
		namespaces = append(namespaces, ns)
	}
	return namespaces
}
