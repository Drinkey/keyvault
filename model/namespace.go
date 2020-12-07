package model

import (
	"fmt"
	"log"

	"github.com/Drinkey/keyvault/internal"
)

type Namespace struct {
	ID        int    `json:"namespace_id"`
	Name      string `json:"name"`
	MasterKey string `json:"master_key"`
	Nonce     string `json:"nonce"`
}

func (ns Namespace) IsEmpty() bool {
	return ns.ID == 0
}

func (ns Namespace) IsNameExist(n string) bool {
	return ns.Name == n
}

func (ns Namespace) Get(name string) (Namespace, error) {
	sqlTemplate := `
	SELECT namespace_id, name, master_key, nonce
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
		err = rows.Scan(&r.ID, &r.Name, &r.MasterKey, &r.Nonce)
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
		var _ns Namespace
		err = rows.Scan(&_ns.ID, &_ns.Name)
		if err != nil {
			log.Fatal(err)
		}
		_ns.MasterKey = internal.KeyMask
		namespaces = append(namespaces, _ns)
	}
	return namespaces
}

func (ns Namespace) Create(n Namespace) error {
	sqlTemplate := `
		INSERT INTO namespace
		(name, master_key, nonce)
		VALUES
		("%s", "%s", "%s");
	`

	sqlStmt := fmt.Sprintf(sqlTemplate, n.Name, n.MasterKey, n.Nonce)
	_, err := conn.Exec(sqlStmt)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
