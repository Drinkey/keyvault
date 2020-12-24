package main

import (
	"flag"
	"os"

	"github.com/Drinkey/keyvault/commands/client"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	clientCommand := flag.NewFlagSet("client", flag.ExitOnError)
	configInit := clientCommand.Bool("init", false, "Whether initialize the client")
	configDir := clientCommand.String("dir", "", "Specify the DIR to store certificates and configuration files")
	clientName := clientCommand.String("name", "keyvault", "The client name, which is same as namespace name")
	clientConfigFile := clientCommand.String("config", "cert.json", "The client certificate config file in JSON")

	adminCommand := flag.NewFlagSet("admin", flag.ExitOnError)

	namespaceCommand := flag.NewFlagSet("namespace", flag.ExitOnError)
	secretCommand := flag.NewFlagSet("secret", flag.ExitOnError)

	switch os.Args[1] {
	case "client":
		clientCommand.Parse(os.Args[2:])
	case "admin":
		adminCommand.Parse(os.Args[2:])
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
			Init:       *configInit,
			ConfigDir:  *configDir,
			ClientName: *clientName,
		}
		c.Configuration.Path = *clientConfigFile
		c.Run()
	}
}
