package certio

type SubjectSchema struct {
	Organization string `json: "organization"`
	Country      string `json:"country"`
	Province     string `json:"province"`
	Locality     string `json:"locality"`
	Address      string `json:"address"`
	PostalCode   string `json:"postal_code"`
	CommonName   string `json:"common_name"`
}

type CASchema struct {
	SerialNumber int64         `json:"serial_number"`
	Subject      SubjectSchema `json:"subject"`
	Valid        int           `json:"valid_year"`
	KeyLength    int           `json:"key_length"`
}

type CertSchema struct {
	SerialNumber int64         `json:"serial_number"`
	Subject      SubjectSchema `json:"subject"`
	DNSName      string        `json:"dns_name"`
	Valid        int           `json:"valid_year"`
	KeyLength    int           `json:"key_length"`
}

type CertConfigSchema struct {
	CA          CASchema   `json:"ca"`
	Certificate CertSchema `json:"cert"`
}

type CertificateSigningRequest struct {
	Request string `json:"csr"`
}

type CertificateResponse struct {
	Certificate string `json:"certificate"`
	CA          string `json:"ca"`
}
