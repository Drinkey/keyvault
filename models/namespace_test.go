package models

import "testing"

func TestNamespaceCreateSuccess(t *testing.T) {
	err := CreateNamespace("name1", "key1", "nonce1")
	if err != nil {
		t.Fail()
	}
}

func TestNamespaceCreateDuplicateFailure(t *testing.T) {
	err := CreateNamespace("name1", "key1", "nonce1")
	if err == nil {
		t.Log("create duplicated namespace expect failure but success")
		t.Log(err.Error())
		t.Fail()
	}
}

func TestNamespaceGetSuccess(t *testing.T) {
	err := CreateNamespace("name2", "key2", "nonce2")
	if err != nil {
		t.Log("create namespace failed")
		t.Fail()
	}
	n, err := GetNamespace("name2")
	if err != nil {
		t.Logf("get namespace failed %s", err.Error())
		t.Fail()
	}
	if n.IsEmpty() {
		t.Log("got empty namespace is not expected")
		t.Fail()
	}
	if n.MasterKey != "key2" {
		t.Log("get namespace with wrong value")
		t.Log(n)
		t.Fail()
	}
}

func TestNamespaceListSuccess(t *testing.T) {
	var err error
	if err = CreateNamespace("TestNamespaceListSuccess1", "KKK1", "NNN1"); err != nil{
		t.Fail()
	}
	if err = CreateNamespace("TestNamespaceListSuccess2", "KKK2", "NNN2");  err != nil{
		t.Fail()
	}
	if err = CreateNamespace("TestNamespaceListSuccess3", "KKK3", "NNN3");  err != nil{
		t.Fail()
	}
	ns,err := ListNamespace()
	if err != nil {
		t.Logf("list namespace failed: %s", err.Error())
		t.Fail()
	}
	
	if len(ns) < 3 {
		t.Log(ns)
		t.Logf("total %d", len(ns))
		t.Fail()
	}
}