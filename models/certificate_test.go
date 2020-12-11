package models

import "testing"

func TestCertificateCreate(t *testing.T) {
	err := CreateCertificateRequest("DB_USER", "---some-request---xxxfffwww", "tttoken")
	if err != nil {
		t.Log("fail to create certificate record.")
		t.Fail()
	}
}

func TestCertificateGet(t *testing.T) {
	certStr := "---some-request---xxxfffwww"
	token := "some_token"
	err := CreateCertificateRequest("DB_USER_2", certStr, token)
	if err != nil {
		t.Log("fail to create certificate record.")
		t.Fail()
	}
	cert, err := GetCertificate("DB_USER_2")
	if err != nil {
		t.Log("failed to get certificate")
		t.Log(err)
		t.Fail()
	}
	if cert.Name != "DB_USER_2" || cert.SignRequest != certStr || cert.Token != token {
		t.Log("got data from db is not match what saved")
		t.Log(cert)
		t.Fail()
	}
}
