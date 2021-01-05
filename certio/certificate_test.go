package certio

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/Drinkey/keyvault/pkg/settings"
	"github.com/Drinkey/keyvault/pkg/utils"
)

func setup() (dir, config string) {
	log.Print("Test Setup")
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	log.Println(pwd)
	projectRoot := utils.DirUpLevel(pwd, -1)
	log.Println(projectRoot)
	testDir := "/tmp/certs"
	testConfig := fmt.Sprintf("%s/keyvaultd-config.json", projectRoot)
	os.Mkdir(testDir, 0777)

	// os.Setenv("KV_CERT_DIR", testDir)
	// os.Setenv("KV_CONFIG_FILE", testConfig)
	// os.Setenv("KV_DB_PATH", "")
	settings.Settings.Parse()
	return testDir, testConfig
}

func TestCertificateCASelfSigned(t *testing.T) {
	setup()
	var cfg CertificateConfiguration
	cfg.Parse()
	t.Log("Current setting")
	t.Log(settings.Settings)
	t.Log("Current config")
	t.Log(cfg)

	var ca CA
	capkey, err := ca.privateKey.Generate(cfg.config.CA.KeyLength)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = ca.privateKey.Save(cfg.Paths.CaPrivKeyPath, capkey)
	if err != nil {
		t.Logf("save private key failed: %s", err.Error())
		t.Log(err)
		t.Fail()
	}
	if !utils.FileExist(cfg.Paths.CaPrivKeyPath) {
		t.Logf("Private Key file does not exist after save: %s", cfg.Paths.CaPrivKeyPath)
		t.Fail()
	}

	catempl := ca.CreateTemplate(cfg.config.CA)
	caCert, err := ca.Issue(catempl, catempl, &capkey.PublicKey, capkey)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	ca.Save(cfg.Paths.CaCertPath, caCert)
	if !utils.FileExist(cfg.Paths.CaCertPath) {
		t.Logf("CA Cert file does not exist after save: %s", cfg.Paths.CaCertPath)
		t.Fail()
	}
	// teardown(testDir)
}

func TestCertificateSignCertificate(t *testing.T) {
	setup()
	var cfg CertificateConfiguration
	cfg.Parse()
	t.Log("Current setting")
	t.Log(settings.Settings)
	t.Log("Current config")
	t.Log(cfg)

	var ca CA
	cacert, caprivkey := ca.Load(cfg.Paths.CaCertPath, cfg.Paths.CaPrivKeyPath)

	var web WebCertificate
	webpkey, err := web.PrivKey.Generate(cfg.config.Web.KeyLength)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = web.PrivKey.Save(cfg.Paths.WebPrivKeyPath, webpkey)
	if err != nil {
		t.Log("save private key failed")
		t.Log(err)
		t.Fail()
	}
	if !utils.FileExist(cfg.Paths.WebPrivKeyPath) {
		t.Logf("Private Key file does not exist after save: %s", cfg.Paths.WebPrivKeyPath)
		t.Fail()
	}

	//Create cert request
	webtempl := web.CreateTemplate(cfg.config.Web)
	// CA sign the request
	webCert, err := ca.Issue(cacert, webtempl, &webpkey.PublicKey, caprivkey)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	web.Save(cfg.Paths.WebCertPath, webCert)
	if !utils.FileExist(cfg.Paths.WebCertPath) {
		t.Logf("Web Cert file does not exist after save: %s", cfg.Paths.WebCertPath)
		t.Fail()
	}
}

func TestCertificateCACertString(t *testing.T) {
	setup()
	var cfg CertificateConfiguration
	cfg.Parse()
	t.Log("Current setting")
	t.Log(settings.Settings)
	t.Log("Current config")
	t.Log(cfg)

	var ca CA
	pem := ca.Read(cfg.Paths.CaCertPath)
	if !strings.Contains(pem, "CERTIFICATE") {
		t.Fail()
	}
}

func TestCertificateAuthorityCache(t *testing.T) {
	setup()
	var cache CertificateAuthority
	var ca CA
	var cfg CertificateConfiguration
	cfg.Parse()

	caCert, caPrivKey := ca.Load(cfg.Paths.CaCertPath, cfg.Paths.CaPrivKeyPath)
	cache.Cache(caCert, ca.String, caPrivKey)
	if !strings.Contains(ca.String, "CERTIFICATE") {
		t.Log("ca String doesn't contain keyword CERTIFICATE")
		t.Log(ca.String)
		t.Fail()
	}
	if !strings.Contains(cache.String, "CERTIFICATE") {
		t.Log("ca String doesn't contain keyword CERTIFICATE")
		t.Log(cache.String)
		t.Fail()
	}
}
