package configuration

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/Drinkey/keyvault/pkg/settings"
)

type Configuration struct {
	Path        string                 `json:"-"`
	ServerAddr  string                 `json:"server_host"`
	Certificate settings.WebCertConfig `json:"certificate"`
}

func (c *Configuration) Read() {
	log.Printf("client configuration reading file: %s", c.Path)
	contentBytes, _ := ioutil.ReadFile(c.Path)
	if err := json.Unmarshal(contentBytes, &c); err != nil {
		log.Print("failed to read config file")
		log.Println(err)
	}
	log.Println(c)
}
