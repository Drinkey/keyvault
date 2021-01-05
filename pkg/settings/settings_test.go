package settings

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Drinkey/keyvault/pkg/utils"
)

func setup() (dir, config string) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Println(pwd)
	projectRoot := utils.DirUpLevel(pwd, -2)
	log.Println(projectRoot)
	testDir := "/tmp/certs"
	testConfig := fmt.Sprintf("%s/keyvaultd-config.json", projectRoot)
	os.Mkdir(testDir, 0777)

	os.Setenv("KV_CERT_DIR", testDir)
	os.Setenv("KV_CONFIG_FILE", testConfig)
	return testDir, testConfig
}

func teardown(dir string) {
	os.RemoveAll(dir)
	os.Unsetenv("KV_CERT_DIR")
	os.Unsetenv("KV_CONFIG_FILE")
}

func TestSettingsConfigurationParsing(t *testing.T) {
	testDir, testConfig := setup()
	Settings.Parse()
	if Settings.CertificateDir != testDir {
		t.Fail()
	}
	if Settings.ConfigFile != testConfig {
		t.Fail()
	}

	if Settings.Certificate.CA.Subject.CommonName != "keyvault.org" {
		t.Fail()
	}
	if Settings.Account.Username != "keyvault" {
		t.Fail()
	}
	if Settings.Service.APIPort != 8080 {
		t.Fail()
	}
	teardown(testDir)
}
