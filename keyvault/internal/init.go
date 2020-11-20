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

func CaConfigParser(filename string, schema *CAConfig) {
	file, _ := ioutil.ReadFile(filename)
	_ = json.Unmarshal([]byte(file), &schema)
}

func CertConfigParser(filename string, schema *CertConfig) {
	file, _ := ioutil.ReadFile(filename)
	_ = json.Unmarshal([]byte(file), &schema)
}

func init() {
	if !FileExist(fmt.Sprintf("%s/%s", CONF_DIR, CA_FILE)) {
		log.Print("CA Cert is not exist, try to create new CA")
		if !FileExist(fmt.Sprintf("%s/%s", CONF_DIR, CA_CONF)) {
			log.Panic("Unable to create new CA because no configuration for CA found")
		}
		if err := CreateCACertificate(); err != nil {
			log.Fatal("creating CA failed: ", err)
		}
		if err := CreateCertificate(); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	} else if !FileExist(fmt.Sprintf("%s/%s", CONF_DIR, CERT_FILE)) {
		log.Print("Certificate is not exist, try to create new certificate")
		if !FileExist(fmt.Sprintf("%s/%s", CONF_DIR, CERT_CONF)) {
			log.Panic("Unable to create new certificate because no configuration for certificate found")
		}
		CreateCertificate()
	}

	log.Print("All required certificates are in place")
}

func DoNothing() string {
	log.Println("hoho")
	return "init"
}
