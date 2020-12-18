/*
Package certio provides all operations against certificate.
*/
package certio

import (
	"log"

	"github.com/Drinkey/keyvault/internal"
)

var Cfg CertificateConfiguration
var CaContainer CertificateAuthority

func initCACertificate(cfg CertificateConfiguration) error {
	log.Printf("Creating CA Certificate %s with %s", Cfg.Paths.CaCertPath, Cfg.file)
	var ca CA
	caprivkey, err := ca.privateKey.Generate(Cfg.config.CA.KeyLength)
	if err != nil {
		log.Fatal(err)
	}

	if err := ca.privateKey.Save(Cfg.Paths.CaPrivKeyPath, caprivkey); err != nil {
		log.Print("save private key failed")
		log.Fatal(err)
	}

	caTempl := ca.CreateTemplate(Cfg.config.CA)
	caCert, err := ca.Issue(caTempl, caTempl, &caprivkey.PublicKey, caprivkey)
	if err != nil {
		log.Fatal(err)
	}

	if err = ca.Save(Cfg.Paths.CaCertPath, caCert); err != nil {
		log.Fatal(err)
	}
	return nil
}

func createWebCertificate(cfg CertificateConfiguration) error {
	log.Printf("Creating Web Certificate %s with %s", Cfg.Paths.WebCertPath, Cfg.file)
	var ca CA
	caCert, caprivkey := ca.Load(Cfg.Paths.CaCertPath, Cfg.Paths.CaPrivKeyPath)

	var web WebCertificate

	webPrivkey, err := web.privateKey.Generate(Cfg.config.Web.KeyLength)
	if err != nil {
		log.Fatal(err)
	}
	err = web.privateKey.Save(Cfg.Paths.WebPrivKeyPath, webPrivkey)
	if err != nil {
		log.Print("save web cert private key failed")
		log.Fatal(err)
	}

	//Create cert request
	webtempl := web.CreateTemplate(Cfg.config.Web)
	// CA sign the request
	webCert, err := ca.Issue(caCert, webtempl, &webPrivkey.PublicKey, caprivkey)
	if err != nil {
		log.Fatal(err)
	}
	if err = web.Save(Cfg.Paths.WebCertPath, webCert); err != nil {
		log.Fatal(err)
	}

	initCaContainer()

	return nil
}

func initCaContainer() {

	if CaContainer.IsSet() {
		log.Print("CA Cache already set, no need to init")
	}
	log.Printf("First time load CA, Storing CA in memory when creating web certificate")
	var ca CA
	caCert, caprivkey := ca.Load(Cfg.Paths.CaCertPath, Cfg.Paths.CaPrivKeyPath)
	CaContainer.Cache(caCert, ca.String, caprivkey)
}

func init() {
	log.SetPrefix("certio: ")

	Cfg.Parse()

	var certConfigFile = Cfg.file
	log.Printf("got cert config path %s", certConfigFile)
	var certDirectory = Cfg.dir
	var CertFiles = Cfg.Paths

	log.Printf("initialize certificates under %s", certDirectory)
	if !internal.FileExist(CertFiles.CaCertPath) {
		log.Printf("CA Cert is not exist, try to create new CA with config file %s", certConfigFile)
		if !internal.FileExist(certConfigFile) {
			log.Panic("Unable to create new CA because no configuration for CA was found")
		}
		if err := initCACertificate(Cfg); err != nil {
			log.Fatal("creating CA failed: ", err)
		}
		if err := createWebCertificate(Cfg); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	} else if !internal.FileExist(Cfg.Paths.WebCertPath) {
		log.Print("Certificate private key is not exist, try to create new certificate with new key")
		if !internal.FileExist(Cfg.Paths.CaCertPath) {
			log.Panic("Unable to create new certificate because no configuration for certificate found")
		}
		if err := createWebCertificate(Cfg); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	}
	initCaContainer()
	log.Print("All required certificates are in place")
}
