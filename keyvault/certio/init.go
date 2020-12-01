package certio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Drinkey/keyvault/internal"
)

func CaConfigParser(filename string, schema *CAConfig) {
	file, _ := ioutil.ReadFile(filename)
	_ = json.Unmarshal([]byte(file), &schema)
}

func CertConfigParser(filename string, schema *CertConfig) {
	file, _ := ioutil.ReadFile(filename)
	_ = json.Unmarshal([]byte(file), &schema)
}

func GetCertFiles(dir string) CertFiles {
	return CertFiles{
		CaCert:         fmt.Sprintf("%s/%s", dir, CA_FILE),
		CaPrivKey:      fmt.Sprintf("%s/ca_priv.key", dir),
		ServerCert:     fmt.Sprintf("%s/%s", dir, CERT_FILE),
		ServerPrivKey:  fmt.Sprintf("%s/cert_priv.key", dir),
		CaCertConf:     fmt.Sprintf("%s/%s", dir, CA_CONF),
		ServerCertConf: fmt.Sprintf("%s/%s", dir, CERT_CONF),
	}
}

func init() {
	log.SetPrefix("certio: ")
	log.Print("init()")
	certs := GetCertFiles(CONF_DIR)
	if !internal.FileExist(certs.CaCert) {
		log.Print("CA Cert is not exist, try to create new CA")
		if !internal.FileExist(certs.CaCertConf) {
			log.Panic("Unable to create new CA because no configuration for CA found")
		}
		ca, err := InitCACertificate(certs)
		if err != nil {
			log.Fatal("creating CA failed: ", err)
		}
		if err := CreateCertificate(ca, certs); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	} else if !internal.FileExist(certs.ServerPrivKey) {
		log.Print("Certificate private key is not exist, try to create new certificate with new key")
		if !internal.FileExist(certs.ServerCertConf) {
			log.Panic("Unable to create new certificate because no configuration for certificate found")
		}
		ca, err := LoadCACertificate(certs)
		if err != nil {
			log.Fatal("CA cert exists but failed to load it. Please delete the CA cert and re-generate it")
		}
		if err := CreateCertificate(ca, certs); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	}

	log.Print("All required certificates are in place")
}

func DoNothing() string {
	log.Println("hoho")
	return "init"
}
