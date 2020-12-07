package certio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/Drinkey/keyvault/internal"
)

func CertConfigParser(filename string, schema *CertConfigSchema) {
	file, _ := ioutil.ReadFile(filename)
	_ = json.Unmarshal([]byte(file), &schema)
}

func GetCertFiles(dir string) CertFilePath {
	return CertFilePath{
		CaCertPath:        fmt.Sprintf("%s/%s", dir, CA_FILE),
		CaPrivKeyPath:     fmt.Sprintf("%s/ca_priv.key", dir),
		ServerCertPath:    fmt.Sprintf("%s/%s", dir, CERT_FILE),
		ServerPrivKeyPath: fmt.Sprintf("%s/cert_priv.key", dir),
	}
}

func init() {
	log.SetPrefix("certio: ")
	log.Printf("initialize certificates under %s", CERT_DIR)
	if !internal.FileExist(CertFiles.CaCertPath) {
		log.Print("CA Cert is not exist, try to create new CA")
		if !internal.FileExist(CERT_CONF_FILE) {
			log.Panic("Unable to create new CA because no configuration for CA was found")
		}
		ca, err := InitCACertificate(CertFiles, CERT_CONF_FILE)
		if err != nil {
			log.Fatal("creating CA failed: ", err)
		}
		if err := CreateCertificate(ca, CertFiles, CERT_CONF_FILE); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	} else if !internal.FileExist(CertFiles.ServerPrivKeyPath) {
		log.Print("Certificate private key is not exist, try to create new certificate with new key")
		if !internal.FileExist(CERT_CONF_FILE) {
			log.Panic("Unable to create new certificate because no configuration for certificate found")
		}
		ca, err := LoadCACertificate(CertFiles)
		if err != nil {
			log.Fatal("CA cert exists but failed to load it. Please delete the CA cert and re-generate it")
		}
		if err := CreateCertificate(ca, CertFiles, CERT_CONF_FILE); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	}

	log.Print("All required certificates are in place")
}
