package client

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/commands/configuration"
	"github.com/Drinkey/keyvault/commands/restclient"
	"github.com/Drinkey/keyvault/pkg/settings"
)

func getSubjectName(subject settings.Subject, namespace string) pkix.Name {
	return pkix.Name{
		Organization:       []string{subject.Organization},
		OrganizationalUnit: []string{namespace},
		Country:            []string{subject.Country},
		Province:           []string{subject.Province},
		Locality:           []string{subject.Locality},
		StreetAddress:      []string{subject.Address},
		PostalCode:         []string{subject.PostalCode},
		CommonName:         subject.CommonName,
	}
}

type ClientCommand struct {
	io            certio.CertIO
	Action        string // which action to take
	ClientName    string // will be used for OU property
	ConfigDir     string // location to store generated certs, tokens
	Configuration configuration.Configuration
}

func (c ClientCommand) Print() {
	log.Printf("Client (OU): %s", c.ClientName)
	log.Printf("Configuration DIR: %s", c.ConfigDir)
}

func (c *ClientCommand) Run() {
	c.Print()
	c.Configuration.Read()
	c.Configuration.Directory = c.ConfigDir
	api := CertificateAPI{Config: c}
	switch c.Action {
	case "get":
		r, err := api.Read(c.ClientName)
		if err != nil {
			log.Panic("API: get certificate error")
		}
		//save cert to local
		if err = c.Configuration.SaveCertificate([]byte(r.Certificate)); err != nil {
			fmt.Printf("Save certificate failed")
			fmt.Println(err)
		}
		//get CA cert
		r, err = api.Read("ca")
		if err != nil {
			log.Panic("API: get CA certificate error")
		}
		//save cert to local
		if err = c.Configuration.SaveCA([]byte(r.CA)); err != nil {
			fmt.Printf("Save CA certificate failed")
			fmt.Println(err)
		}
	case "init":
		log.Printf("Initializing client...")
		ci := CertificateInitiator{Config: c}
		// create private key pair and csr, save private key
		csr := ci.Init()
		log.Print("CSR and Private initialized successful.")
		// send csr to keyvault service for signing
		api.Create(c.ClientName, csr)
	}
}

// CertificateInitiator initiates client certificates
type CertificateInitiator struct {
	Config *ClientCommand
}

func (c *CertificateInitiator) Init() []byte {
	csrConfig := c.Config.Configuration.Certificate
	// generate private key and save to destinated location
	pk := c.initClientPrivateKey(csrConfig)
	return c.initClientCSR(csrConfig, pk)
}

func (c *CertificateInitiator) initClientPrivateKey(csr settings.WebCertConfig) *rsa.PrivateKey {
	clientPrivKeyPath := fmt.Sprintf("%s/client.key", c.Config.ConfigDir)
	var web certio.WebCertificate

	webPrivkey, err := web.PrivKey.Generate(csr.KeyLength)
	if err != nil {
		log.Fatal(err)
	}
	err = web.PrivKey.Save(clientPrivKeyPath, webPrivkey)
	if err != nil {
		log.Print("save web cert private key failed")
		log.Fatal(err)
	}
	return webPrivkey
}

func (c *CertificateInitiator) initClientCSR(csr settings.WebCertConfig, pk *rsa.PrivateKey) []byte {
	//Create cert request
	template := &x509.CertificateRequest{
		Subject:            getSubjectName(csr.Subject, c.Config.ClientName),
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, template, pk)
	if err != nil {
		log.Print("Create certificate request error")
		log.Fatal(err)
	}
	pemBuffer := new(bytes.Buffer)
	pem.Encode(pemBuffer, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})
	return pemBuffer.Bytes()
}

type CertReq struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	SignRequest string `json:"req"`
	Certificate string `json:"certificate"`
	Token       string `json:"token"`
	CA          string `json:"ca"`
}

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	CertReq `json:"data"`
}

// CertificateAPI interacts keyvault service /api/v*/certificate
type CertificateAPI struct {
	Config *ClientCommand
	api    restclient.RESTFulClienter
}

func (c *CertificateAPI) initAPI() {
	if c.api == nil {
		rclient := restclient.RESTFulClient{
			Host:     c.Config.Configuration.Server.Host,
			Port:     c.Config.Configuration.Server.TLSMaintenancePort,
			Insecure: true,
		}
		c.api = restclient.Certificate{rclient}
	}
}

func (c *CertificateAPI) Create(name string, csr []byte) error {
	// create certificate item in keyvault service
	c.initAPI()
	csrString := string(csr)
	encode := ""
	for _, line := range strings.Split(csrString, "\n") {
		encode = fmt.Sprintf("%s%s", encode, line)
	}
	r := CertReq{Name: name, SignRequest: csrString}

	resp, err := c.api.Create("", r)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		panic(fmt.Sprintf("Create CSR in keyvault service failed, status code %d", resp.StatusCode))
	}
	log.Print("Create CSR in keyvault service success, wait for admin to sign the request")
	respBytes, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%q", respBytes)
	return nil
}

func (c *CertificateAPI) Read(name string) (Resp, error) {
	c.initAPI()
	q := map[string]string{
		"q": name,
	}
	resp, err := c.api.Read("", q)
	if err != nil {
		log.Fatalf("failed to GET certificate %s: %s", name, err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to GET certificate %s: service returns error", name)
		log.Fatalf("status=%s", resp.Status)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	var r Resp
	err = json.Unmarshal(respBytes, &r)
	if err != nil {
		log.Printf("failed to decode response body")
		return Resp{}, err
	}
	return r, nil
}
