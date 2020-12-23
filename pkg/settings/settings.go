package settings

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Service struct {
	TLSPort int `json:"tls_port"`
	APIPort int `json:"api_port"`
}

// Settings stores all values of configuration, make it singleton
var Settings Setting

type Setting struct {
	Account        Account  `json:"account"`
	Service        Service  `json:"service"`
	Certificate    CertJSON `json:"certificates"`
	ConfigFile     string   `json:"-"`
	CertificateDir string   `json:"-"`
	DatabasePath   string   `json:"-"`
	// CertificateConfig CertificateConfiguration `json:"-"`
}

func (s *Setting) Parse() {
	s.getKvCertDir()
	s.getKvConfig()
	s.getKvDatabasePath()
	s.parseJSON()
}

func (s *Setting) parseJSON() {
	// var schema *setting.CertJSON
	log.Printf("Reading config file %s", s.ConfigFile)
	contentBytes, _ := ioutil.ReadFile(s.ConfigFile)
	_ = json.Unmarshal(contentBytes, &s)
	log.Println(s)
}

func (s *Setting) getKvCertDir() {
	log.Print("reading env cert dir")
	s.CertificateDir = os.Getenv("KV_CERT_DIR")
	log.Printf("got settings.CertificateDir = %s", s.CertificateDir)
}

func (s *Setting) getKvConfig() {
	log.Print("reading env config")
	s.ConfigFile = os.Getenv("KV_CONFIG_FILE")
	log.Printf("got settings.ConfigFile = %s", s.ConfigFile)
}

func (s *Setting) getKvDatabasePath() {
	log.Print("reading env config")
	s.DatabasePath = os.Getenv("KV_DB_PATH")
	log.Printf("got settings.DatabasePath = %s", s.DatabasePath)
}

func init() {
	log.SetPrefix("settings: ")
	Settings.Parse()
}
