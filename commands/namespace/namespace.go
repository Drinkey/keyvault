package namespace

import (
	"log"
)

type NamespaceCommand struct {
	Action    string // which action to take
	Name      string // will be used for OU property
	ConfigDir string // location to store generated certs, tokens
}

func (c NamespaceCommand) Print() {
	log.Printf("Client (OU): %s", c.Name)
	log.Printf("Configuration DIR: %s", c.ConfigDir)
}

func (c *NamespaceCommand) Run() {
	c.Print()
	api := NamespaceAPI{Config: c}
	switch c.Action {
	case "get":
		err := api.Read(c.Name)
		if err != nil {
			log.Panic("API: get certificate error")
		}
	case "create":
		// send csr to keyvault service for signing
		api.Create(c.Name)
	}
}

type NamespaceAPI struct {
	Config *NamespaceCommand
}

func (n NamespaceAPI) Create(name string) {
	return
}

func (n NamespaceAPI) Read(name string) error {
	return nil
}
