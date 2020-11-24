package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

const CONF_DIR = "etc"
const CA_CONF = "ca.json"
const CA_FILE = "ca.crt"
const CERT_CONF = "cert.json"
const CERT_FILE = "cert.pem"

type SubjectConfig struct {
	Organization string `json: "organization"`
	Country      string `json:"country"`
	Province     string `json:"province"`
	Locality     string `json:"locality"`
	Address      string `json:"address"`
	PostalCode   string `json:"postal_code"`
}

type CAConfig struct {
	SerialNumber int64         `json:"serial_number"`
	Subject      SubjectConfig `json:"subject"`
	Valid        int           `json:"valid"`
	KeyLength    int           `json:"key_length"`
}

type CertConfig struct {
	SerialNumber int64         `json:"serial_number"`
	Subject      SubjectConfig `json:"subject"`
	IPv4Address  string        `json:"ip4addr"`
	Valid        int           `json:"valid"`
	KeyLength    int           `json:"key_length"`
}

type CertFiles struct {
	CaCert         string
	CaPrivKey      string
	ServerCert     string
	ServerPrivKey  string
	CaCertConf     string
	ServerCertConf string
}

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
	certs := GetCertFiles(CONF_DIR)
	if !FileExist(certs.CaCert) {
		log.Print("CA Cert is not exist, try to create new CA")
		if !FileExist(certs.CaCertConf) {
			log.Panic("Unable to create new CA because no configuration for CA found")
		}
		ca, err := InitCACertificate(certs)
		if err != nil {
			log.Fatal("creating CA failed: ", err)
		}
		if err := CreateCertificate(ca, certs); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	} else if !FileExist(certs.ServerPrivKey) {
		log.Print("Certificate is not exist, try to create new certificate")
		if !FileExist(certs.ServerCertConf) {
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
