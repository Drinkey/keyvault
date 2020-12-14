/*
Package certio provides all operations against certificate.
*/
package certio

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"os"
)

func Issue(c *x509.Certificate, ca CertificateAuthority, k interface{}) ([]byte, error) {
	certBytes, err := x509.CreateCertificate(rand.Reader, c, ca.CaCert, k, ca.CaPrivKey)
	if err != nil {
		return []byte{}, err
	}
	return certBytes, nil
}

func SavePemFile(certype string, content []byte, filename string, perm os.FileMode) error {
	PEM := new(bytes.Buffer)
	pem.Encode(PEM, &pem.Block{
		Type:  certype,
		Bytes: content,
	})

	return ioutil.WriteFile(filename, PEM.Bytes(), 0600)
}

func getSubjectName(sc SubjectSchema) pkix.Name {
	return pkix.Name{
		Organization:  []string{sc.Organization},
		Country:       []string{sc.Country},
		Province:      []string{sc.Province},
		Locality:      []string{sc.Locality},
		StreetAddress: []string{sc.Address},
		PostalCode:    []string{sc.PostalCode},
		CommonName:    sc.CommonName,
	}
}
