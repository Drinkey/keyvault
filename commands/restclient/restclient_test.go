package restclient

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Drinkey/keyvault/commands/configuration"
	"github.com/Drinkey/keyvault/pkg/utils"
)

func getProjectConfigFile() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	projectRoot := utils.DirUpLevel(pwd, -2)
	return fmt.Sprintf("%s/keyvault-client.json", projectRoot)
}
func TestNamespaceUrlWithoutQueryString(t *testing.T) {
	config := getProjectConfigFile()
	cfg := configuration.Configuration{Path: config}
	cfg.Read()

	// var client RESTFulClient
	rclient := RESTFulClient{Host: cfg.Server.Host, Port: cfg.Server.TLSMaintenancePort}
	client := Certificate{rclient}
	url := client.Url("")
	fmt.Println(url)
	if url != "https://keyvault.org:1443/api/v1/certificate" {
		t.Logf("Actual URL: %s", url)
		t.Fail()
	}
}

func TestNamespaceGetWithQueryString(t *testing.T) {
	config := getProjectConfigFile()
	cfg := configuration.Configuration{Path: config}
	cfg.Read()

	rclient := RESTFulClient{Host: cfg.Server.Host, Port: cfg.Server.TLSMaintenancePort}
	client := Certificate{rclient}
	query := map[string]string{
		"q": "ca",
	}
	_, _ = client.Read("", query)
	// fmt.Println(r)
	url := client.Url("", query)

	if url != "https://keyvault.org:1443/api/v1/certificate?q=ca" {
		t.Logf("Actual URL: %s", url)
		t.Fail()
	}
}

func TestGetWithPathAndQueryString(t *testing.T) {
	config := getProjectConfigFile()
	cfg := configuration.Configuration{Path: config}
	cfg.Read()

	rclient := RESTFulClient{Host: cfg.Server.Host, Port: cfg.Server.TLSMaintenancePort}
	client := Certificate{rclient}
	query := map[string]string{
		"q": "ca",
	}
	_, _ = client.Read("keytest", query)
	// fmt.Println(r)
	url := client.Url("keytest", query)

	if url != "https://keyvault.org:1443/api/v1/certificate/keytest?q=ca" {
		t.Logf("Actual URL: %s", url)
		t.Fail()
	}
}
