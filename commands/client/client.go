package client

import (
	"fmt"
	"log"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/commands/configuration"
)

type ClientCommand struct {
	Init          bool
	ClientName    string
	ConfigDir     string
	Configuration configuration.Configuration
}

func (c ClientCommand) Print() {
	log.Printf("Client: %s", c.ClientName)
	log.Printf("Configuration DIR: %s", c.ConfigDir)
	// log.Printf("Certificate Config: %s", c.CertConfig)
}

func (c *ClientCommand) Run() {
	c.Print()
	c.Configuration.Read()
	log.Printf("Parsed config")
	log.Println(c.Configuration)
	if c.Init {
		log.Printf("Initializing client...")
		c.initCertificate()
	}

}

func (c *ClientCommand) initCertificate() {
	// clientCertPath := fmt.Sprintf("%s/client.crt", c.ConfigDir)
	clientPrivKeyPath := fmt.Sprintf("%s/client.key", c.ConfigDir)
	// caCertPath := fmt.Sprintf("%s/ca.crt", c.ConfigDir)

	var web certio.WebCertificate

	webPrivkey, err := web.PrivKey.Generate(c.Configuration.Certificate.KeyLength)
	if err != nil {
		log.Fatal(err)
	}
	err = web.PrivKey.Save(clientPrivKeyPath, webPrivkey)
	if err != nil {
		log.Print("save web cert private key failed")
		log.Fatal(err)
	}

	//Create cert request
	// webtempl := web.CreateTemplate(c.Configuration.Certificate)
}

// func certificateInit(){
// 	parse config
// 	if cert file exist:
// 		return
// 	else
// 		create pkey

// 	if pkey file exist:
// 		create new cert
// 		return
// }

// func createNewCertificate(){
// 	parse cert config
// 	set OU same as namespace
// 	create csr PEM
// 	parse csr PEM and replace \n with \\n (convert the file in one line)
// 	post to service
// 	get the signed cert and save to [namespace].crt
// }

// func createHTTPSClient() *http.Client {
// 	// Load client cert
// 	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Load CA cert
// 	caCert, err := ioutil.ReadFile(*caFile)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	caCertPool := x509.NewCertPool()
// 	caCertPool.AppendCertsFromPEM(caCert)

// 	// Setup HTTPS client
// 	tlsConfig := &tls.Config{
// 		Certificates: []tls.Certificate{cert},
// 		RootCAs:      caCertPool,
// 	}
// 	tlsConfig.BuildNameToCertificate()
// 	transport := &http.Transport{TLSClientConfig: tlsConfig}
// 	return &http.Client{Transport: transport}
// }
