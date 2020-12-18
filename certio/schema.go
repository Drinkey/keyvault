/*
Package certio provides all operations against certificate.
*/
package certio

type CertificateSigningRequest struct {
	Request string `json:"csr"`
}

type CertificateResponse struct {
	Certificate string `json:"signed"`
	CA          string `json:"ca"`
}
