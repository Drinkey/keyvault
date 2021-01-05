package request

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/Drinkey/keyvault/certio"
	"github.com/Drinkey/keyvault/pkg/utils"
)

func getProjectConfigFile() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	projectRoot := utils.DirUpLevel(pwd, -1)
	return fmt.Sprintf("%s/keyvault-client.json", projectRoot)
}

func TestInsecureSendRequestNoTLS(t *testing.T) {
	req := Requests{Certificates: certio.CertFilePaths{}}
	req.InitClient(true)
	resp, err := req.Get("http://keyvault.org:8080/ping")
	if err != nil {
		t.Fail()
	}
	fmt.Println(resp)
	t.Logf("status: %s", resp.Status)
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer resp.Body.Close()
	fmt.Printf("%s", respBytes)
}

// func TestInsecureSendRequestTLS(t *testing.T) {
// 	config := getProjectConfigFile()
// 	cfg := configuration.Configuration{Path: config}
// 	cfg.Read()
// 	req := Request{Certificates: cfg}
// 	req.InitClient(true)
// 	resp, err := req.Get("http://keyvault.org:8080/ping")
// 	if err != nil {
// 		t.Fail()
// 	}
// 	fmt.Println(resp)
// 	t.Logf("status: %s", resp.Status)
// 	respBytes, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		t.Log(err)
// 		t.Fail()
// 	}
// 	defer resp.Body.Close()
// 	fmt.Printf("%s", respBytes)
// }
