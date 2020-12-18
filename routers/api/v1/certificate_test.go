package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

const certReq = `-----BEGIN CERTIFICATE REQUEST-----\nMIIEvTCCAqUCAQAweDELMAkGA1UEBhMCQ04xCzAJBgNVBAgMAlNDMQswCQYDVQQH\nDAJDRDEYMBYGA1UECgwPS2V5VmF1bHQgQ2xpZW50MRMwEQYDVQQLDApLVUJFUk5F\nVEVTMSAwHgYDVQQDDBdLVUJFUk5FVEVTLmtleXZhdWx0Lm9yZzCCAiIwDQYJKoZI\nhvcNAQEBBQADggIPADCCAgoCggIBANAr+7KoWS1T70Z9TBEdgGABvJ5jA+Ezd/sb\ntBILZDqn58rxFfgbyMKL+vXRPIUbXocs0uQ2uwySKmJUewWi6oNu2EmCb2wHSm0T\nwmR+68p04AtE/zqR3KaDFXZWNko9Bl0yxqH03T5WscDPcH2fWAoPQxov0qw9Mw/W\n/jNzdiGbvNMfW1QKjYC97y3w9vXJoEpgh+onpGdj0dUgTff01D48FQEXxPwXH940\n7/Cp02bN+e7UAS5daEnHFloTwwyl/UhbhFzeJBco6YEfW3zMQd2iqruEfvDUfsSw\nKVF1tOFKIHb4VKK2l8lObbIwaEHq/flTLBTp21eJWrPNOo0hn8ChBwvBptVC16xZ\nkXJYHF+hTu/5P/0MFcYEF0jchebUEQ4er9Y13oeWNa5vNmo4kBlmMHZDQn1WUKYM\n2OHrx+6DKtIYFF8EAnhJ+3xMqNgYW6VNn7Br44S3yubr2OohG61nqGOjOI6DuufB\nNZx1pFjsF36WZgkJz6Z3b/rskBOl5aPw8pb6hYpaUwM5sIzru8ZOmY1HxZUChjd+\ne7Hlr9Gh+4p+mDv0eoAdk8QdEsHKGykmCCKBSIN+Wh1lWiwz38gavMnBj6hB1yE6\nZTWadp/NFXjrzoN+W+8S4pVh9e8MYrqAtCds/NG3rzr/zQ+1re4WqkUuLwm3zOW/\n0lGIsU/rAgMBAAGgADANBgkqhkiG9w0BAQsFAAOCAgEATm7CJrq0bDFj0YxK201h\nPnM+1QSEZoHgoueauOxTpSMKSld2vKy33eN3v5ak4aQLsIefZiAKiVJnPdhVB5Up\nMWnSewobU4ubbaPGjcCF1kc5NNPDmWrspjwO3b+C4wUW3txPfszz5/ZrhY6BIgQc\nxlz3JAMq1f81dbH3mG6iPzuncmeb4AuBakh6b+hsPlIXHepBaGWBkafrR7XwUjWy\nXCA0HBoO5KdWdmOtm5faVEdGAM0RTUfajmPUgCHe6aqLNqeGoN/ZuDTmi5HW4+Hw\nvzE9Akw8xLENie32gfkr+ADIWSF3kBbuGvP+5ez7slJOKd+cQZUyovVB2E215t8T\nASjv9LNy0bhgVwB6XK0MgOkdl1zNXWMGLrAM/99wmRPrhrZGa2gsZ8OK+rNdZKWy\nRiLQGG46VHExDKOharFhRBu3fZvvsTyc5tOS/3BBXtPkGXC8VGq5E1Ew7fzXl/vr\ncklNaVT5xcesb1gpWlOYOaLOxq1fnUQQHByMfnbkPJspy3t0Yimc6Sq5Mtj+hsl8\nlLTft+xa+e/R72+tqiqHPExrQ/xKCw49FoFSS73pIFk//4cpythyNfwP0ygFeAcD\nhv/HVQPS4t+C1U3v+G94Ihg8yXOg/0zlBkbggequ9RQj1tTH0iXvWzNAI9wFvpm4\n9G3WWhl+rekaNomQPITBVkU=\n-----END CERTIFICATE REQUEST-----`

func TestCertificateCreate(t *testing.T) {
	r := gin.Default()
	w := httptest.NewRecorder()
	uri := "/api/v1/cert/req"
	r.POST(uri, CreateCertificateRequest)

	reqJSON, _ := json.Marshal(map[string]string{
		"name": "DATABASE",
		"req":  certReq,
	})
	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(reqJSON))
	r.ServeHTTP(w, req)

	var resp Certificate
	t.Log(w.Body.String())
	json.Unmarshal(w.Body.Bytes(), &resp)

	if w.Code != http.StatusCreated {
		t.Logf("response code validation failed: code=%d, expected=%d", w.Code, http.StatusCreated)
		t.Fail()
	}
	if resp.Name != "DATABASE" {
		t.Log("response body validation failed:")
		t.Logf("got response %s, expect: DATABASE", resp.Name)
		t.Logf("got response %s, expect: %s", resp.SignRequest, certReq)
		t.Log(resp)
		t.Fail()
	}
}

func TestCertificateGet(t *testing.T) {
	r := gin.Default()
	w := httptest.NewRecorder()
	uri := "/api/v1/cert"
	r.POST(fmt.Sprintf("%s/req", uri), CreateCertificateRequest)
	r.GET(fmt.Sprintf("%s/", uri), GetCertificate)

	reqJSON, _ := json.Marshal(map[string]string{
		"name": "DATABASE_MYSQL",
		"req":  certReq,
	})
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/req", uri), bytes.NewBuffer(reqJSON))
	r.ServeHTTP(w, req)

	var resp Certificate
	t.Log(w.Body.String())
	json.Unmarshal(w.Body.Bytes(), &resp)

	if w.Code != http.StatusCreated {
		t.Logf("response code validation failed: code=%d, expected=%d", w.Code, http.StatusCreated)
		t.Fail()
	}
	if resp.Name != "DATABASE_MYSQL" {
		t.Log("response body validation failed:")
		t.Logf("got response %s, expect: DATABASE", resp.Name)
		t.Logf("got response %s, expect: %s", resp.SignRequest, certReq)
		t.Log(resp)
		t.Fail()
	}

	t.Log("Getting the cert just added")

	req, _ = http.NewRequest("GET", fmt.Sprintf("%s/?q=DATABASE_MYSQL", uri), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Log(w.Body.String())
	json.Unmarshal(w.Body.Bytes(), &resp)
	if w.Code != http.StatusOK {
		t.Logf("response code validation failed: code=%d, expected=%d", w.Code, http.StatusCreated)
		t.Fail()
	}
	if resp.Name != "DATABASE_MYSQL" {
		t.Log("response body validation failed:")
		t.Logf("got response %s, expect: DATABASE", resp.Name)
		t.Logf("got response %s, expect: %s", resp.SignRequest, certReq)
		t.Log(resp)
		t.Fail()
	}
}

func TestCertificateCACertGet(t *testing.T) {
	r := gin.Default()
	w := httptest.NewRecorder()
	uri := "/api/v1/cert"
	r.GET(fmt.Sprintf("%s/ca", uri), GetCACertificate)
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/ca", uri), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var resp CACertificate
	json.Unmarshal(w.Body.Bytes(), &resp)
	if !strings.Contains(resp.Certificate, "CERTIFICATE") {
		t.Logf("got actual response: %s", resp.Certificate)
		t.Fail()
	}
}
