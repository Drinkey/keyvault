package configuration

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Drinkey/keyvault/pkg/utils"
)

func TestConfigurationParsing(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Println(pwd)
	projectRoot := utils.DirUpLevel(pwd, -2)
	log.Println(projectRoot)
	cfg := Configuration{Path: fmt.Sprintf("%s/keyvault-client.json", projectRoot)}
	cfg.Read()
	if cfg.Server.HTTPPort != 8080 ||
		cfg.Server.Host != "keyvault.org" ||
		cfg.Certificate.SerialNumber != 49001 ||
		cfg.Certificate.Subject.Organization != "KeyVault, ORG. Client" {
		t.Fail()
	}
}
