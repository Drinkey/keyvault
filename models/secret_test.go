package models

import "testing"

func TestSecretCreate(t *testing.T) {
	err := CreateNamespace("secret_namespace1", "key1", "nonce1")
	if err != nil {
		t.Log("create name space failed")
		t.Fail()
	}
	err = CreateSecret("sec_key1", "sec_value", "secret_namespace1")
	if err != nil {
		t.Log("create secret failed")
		t.Fail()
	}
}

func TestSecretGet(t *testing.T) {
	err := CreateNamespace("secret_namespace2", "key2", "nonce2")
	if err != nil {
		t.Log("create namespace failed")
		t.Fail()
	}
	err = CreateSecret("sec_key2", "sec_value2", "secret_namespace2")
	if err != nil {
		t.Log("create secret failed")
		t.Fail()
	}
	s, err := GetSecret("sec_key2", "secret_namespace2")
	t.Log(s)
	if err != nil {
		t.Log("get secret failed")
		t.Fail()
	}
	if s.Key != "sec_key2" || s.Value != "sec_value2" {
		t.Log("validate retrieved secret data failed")
		t.Logf("got key=%s, value=%s", s.Key, s.Value)
		t.Fail()
	}
	if s.Namespace.MasterKey != "key2" || s.Namespace.Nonce != "nonce2" {
		t.Log("validate retrieved secret data Namespace field failed")
		t.Log("got namespace: ")
		t.Log(s.Namespace)
		t.Fail()
	}
}
