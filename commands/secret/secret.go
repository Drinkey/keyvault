package namespace

import (
	"log"
)

type SecretCommand struct {
	Action    string // which action to take
	Namespace string // will be used for OU property
	Key       string // key to operate
	Value     string // key's value
	ConfigDir string // location to store generated certs, tokens
}

func (c SecretCommand) Print() {
	log.Printf("Client (OU): %s", c.Namespace)
	log.Printf("Configuration DIR: %s", c.ConfigDir)
}

func (c *SecretCommand) Run() {
	c.Print()
	api := SecretAPI{Config: c}
	switch c.Action {
	case "get":
		r, err := api.Read(c.Namespace, c.Key)
		if err != nil {
			log.Panic("API: get certificate error")
		}
	case "create":
		// send csr to keyvault service for signing
		api.Create(c.Namespace, c.Key, c.Value)
	}
}
