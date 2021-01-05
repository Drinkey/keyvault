package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Drinkey/keyvault/pkg/settings"
)

type ServerConfig struct {
	Host               string `json:"host"`
	HTTPPort           int    `json:"http_port"`
	TLSSecretPort      int    `json:"tls_secret_port"`
	TLSMaintenancePort int    `json:"tls_maintenance_port"`
}

type Configuration struct {
	Directory   string                 `json:"-"`
	Path        string                 `json:"-"`
	Server      ServerConfig           `json:"server"`
	Certificate settings.WebCertConfig `json:"certificate"`
}

func (c *Configuration) Read() {
	log.Printf("client configuration file: %s", c.Path)
	contentBytes, _ := ioutil.ReadFile(c.Path)
	if err := json.Unmarshal(contentBytes, &c); err != nil {
		log.Print("failed to read config file")
		log.Println(err)
	}
}

func (c *Configuration) Save(filename string, content []byte, perm os.FileMode) error {
	f := fmt.Sprintf("%s/%s", c.Directory, filename)
	return ioutil.WriteFile(f, content, perm)
}

func (c Configuration) SaveCertificate(content []byte) error {
	return c.Save("client.crt", content, 0600)
}

func (c Configuration) SaveCA(content []byte) error {
	return c.Save("ca.crt", content, 0600)
}
