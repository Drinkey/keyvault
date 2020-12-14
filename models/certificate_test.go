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

func TestCertificateUpdateSignedCertificateByName(t *testing.T) {
	certStr := "---some-request---xxxfffwww"
	token := "some_token__rand2"
	err := CreateCertificateRequest("DB_USER_3", certStr, token)
	if err != nil {
		t.Log("fail to create certificate record.")
		t.Fail()
	}
	signedCertPem := "Some_Signed_Cert_PEM"
	u, err := UpdateSignedCertificateByName("DB_USER_3", signedCertPem)
	if err != nil {
		t.Log("fail to update certificate record.")
		t.Fail()
	}
	t.Log(u)
	if u.Token != token || u.SignRequest != certStr || u.Certificate != signedCertPem {
		t.Fail()
	}
}