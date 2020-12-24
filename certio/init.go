/*
Package certio provides all operations against certificate.
*/
package certio

import (
	"log"
)

var Cfg CertificateConfiguration
var CaContainer CertificateAuthority

func InitCACertificate(cfg CertificateConfiguration) error {
	log.Printf("Creating CA Certificate %s with %s", Cfg.Paths.CaCertPath, Cfg.File)
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

func CreateWebCertificate(cfg CertificateConfiguration) error {
	log.Printf("Creating Web Certificate %s with %s", Cfg.Paths.WebCertPath, Cfg.File)
	var ca CA
	caCert, caprivkey := ca.Load(Cfg.Paths.CaCertPath, Cfg.Paths.CaPrivKeyPath)

	var web WebCertificate

	webPrivkey, err := web.PrivKey.Generate(Cfg.config.Web.KeyLength)
	if err != nil {
		log.Fatal(err)
	}
	err = web.PrivKey.Save(Cfg.Paths.WebPrivKeyPath, webPrivkey)
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

	InitCaContainer()

	return nil
}

func InitCaContainer() {

	if CaContainer.IsSet() {
		log.Print("CA Cache already set, no need to init")
		return
	}
	log.Printf("First time load CA, Storing CA in memory when creating web certificate")
	var ca CA
	caCert, caprivkey := ca.Load(Cfg.Paths.CaCertPath, Cfg.Paths.CaPrivKeyPath)
	CaContainer.Cache(caCert, ca.String, caprivkey)
}
