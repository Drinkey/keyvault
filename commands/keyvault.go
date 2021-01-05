package main

import (
	"flag"
	"os"

	"github.com/Drinkey/keyvault/commands/client"
	"github.com/Drinkey/keyvault/commands/namespace"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	clientCommand := flag.NewFlagSet("client", flag.ExitOnError)
	clientAction := clientCommand.String("action", "get", "Client action to take")
	configDir := clientCommand.String("dir", "", "Specify the DIR to store certificates and configuration files")
	clientName := clientCommand.String("name", "keyvault", "The client name, which is same as namespace name")
	clientConfigFile := clientCommand.String("config", "cert.json", "The client certificate config file in JSON")

	namespaceCommand := flag.NewFlagSet("namespace", flag.ExitOnError)
	nsAction := namespaceCommand.String("action", "", "Which action to take")
	nsConfigDir := namespaceCommand.String("dir", "", "Directory to store credentials")
	nsName := namespaceCommand.String("name", "", "namespace to operate")

	secretCommand := flag.NewFlagSet("secret", flag.ExitOnError)
	// sAction := secretCommand.String("action", "", "Which action to take")
	// sConfigDir := secretCommand.String("dir", "", "Directory to store credentials")
	// sNamespace := secretCommand.String("namespace", "", "namespace to operate")
	// sKey := secretCommand.String("key", "", "key to operate")

	switch os.Args[1] {
	case "client":
		clientCommand.Parse(os.Args[2:])
	case "namespace":
		namespaceCommand.Parse(os.Args[2:])
	case "secret":
		secretCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if clientCommand.Parsed() {
		c := client.ClientCommand{
			Action:     *clientAction,
			ConfigDir:  *configDir,
			ClientName: *clientName,
		}
		c.Configuration.Path = *clientConfigFile
		c.Run()
	}

	if namespaceCommand.Parsed() {
		c := namespace.NamespaceCommand{
			Action:    *nsAction,
			ConfigDir: *nsConfigDir,
			Name:      *nsName,
		}
		c.Run()
	}

	if secretCommand.Parsed() {

	}
}
