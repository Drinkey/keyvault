package certificate_service

import (
	"log"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/pkg/settings"
	"github.com/Drinkey/keyvault/pkg/utils"
)

func init() {
	log.SetPrefix("certio: ")

	settings.Settings.Parse()

	certio.Cfg.Parse()

	var certConfigFile = certio.Cfg.File
	log.Printf("got cert config path %s", certConfigFile)
	var certDirectory = certio.Cfg.Dir
	var CertFiles = certio.Cfg.Paths

	log.Printf("initialize certificates under %s", certDirectory)
	if !utils.FileExist(CertFiles.CaCertPath) {
		log.Printf("CA Cert is not exist, try to create new CA with config file %s", certConfigFile)
		if !utils.FileExist(certConfigFile) {
			log.Panic("Unable to create new CA because no configuration for CA was found")
		}
		if err := certio.InitCACertificate(certio.Cfg); err != nil {
			log.Fatal("creating CA failed: ", err)
		}
		if err := certio.CreateWebCertificate(certio.Cfg); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	} else if !utils.FileExist(certio.Cfg.Paths.WebCertPath) {
		log.Print("Certificate private key is not exist, try to create new certificate with new key")
		if !utils.FileExist(certio.Cfg.Paths.CaCertPath) {
			log.Panic("Unable to create new certificate because no configuration for certificate found")
		}
		if err := certio.CreateWebCertificate(certio.Cfg); err != nil {
			log.Fatal("creating certificate failed: ", err)
		}
	}
	certio.InitCaContainer()
	log.Print("All required certificates are in place")
}
