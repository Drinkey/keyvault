/*
Package certio provides all operations against certificate.
*/
package certio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Drinkey/keyvault/internal"
)

func CertConfigParser(filename string, schema *CertConfigSchema) {
	file, _ := ioutil.ReadFile(filename)
	_ = json.Unmarshal([]byte(file), &schema)
}

func getKvCertDir() string {
	log.Print("reading env cert dir")
	return os.Getenv("KV_CERT_DIR")
}

func getKvCertConfig() string {
	log.Print("reading env cert config")
	return os.Getenv("KV_CERT_CONF")
}

func GetCertFiles() CertFilePath {
	var dir = getKvCertDir()
	log.Printf("got cert dir %s", dir)
	return CertFilePath{
		CaCertPath:        fmt.Sprintf("%s/%s", dir, caFile),
		CaPrivKeyPath:     fmt.Sprintf("%s/ca_priv.key", dir),
		ServerCertPath:    fmt.Sprintf("%s/%s", dir, certFile),
		ServerPrivKeyPath: fmt.Sprintf("%s/cert_priv.key", dir),
	}
}

func init() {
	log.SetPrefix("certio: ")
	var certConfigFile = getKvCertConfig()
	log.Printf("got cert config path %s", certConfigFile)
	var certDirectory = getKvCertDir()
	var CertFiles = GetCertFiles()

	log.Printf("initialize certificates under %s", certDirectory)
	if !internal.FileExist(CertFiles.CaCertPath) {
		log.Printf("CA Cert is not exist, try to create new CA with config file %s", certConfigFile)
		if !internal.FileExist(certConfigFile) {
			log.Panic("Unable to create new CA because no configuration for CA was found")
		}
		ca, err := InitCACertificate(CertFiles, certConfigFile)
		if err != nil {
			log.Fatal("creating CA failed: ", err)
		}
		if err := CreateCertificate(ca, CertFiles, certConfigFile); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	} else if !internal.FileExist(CertFiles.ServerPrivKeyPath) {
		log.Print("Certificate private key is not exist, try to create new certificate with new key")
		if !internal.FileExist(certConfigFile) {
			log.Panic("Unable to create new certificate because no configuration for certificate found")
		}
		ca, err := LoadCACertificate(CertFiles)
		if err != nil {
			log.Fatal("CA cert exists but failed to load it. Please delete the CA cert and re-generate it")
		}
		if err := CreateCertificate(ca, CertFiles, certConfigFile); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	}

	log.Print("All required certificates are in place")
}
