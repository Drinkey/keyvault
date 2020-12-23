package settings

type Subject struct {
	Organization string `json:"organization"`
	Country      string `json:"country"`
	Province     string `json:"province"`
	Locality     string `json:"locality"`
	Address      string `json:"address"`
	PostalCode   string `json:"postal_code"`
	CommonName   string `json:"common_name"`
}

type CaCertConfig struct {
	SerialNumber int64   `json:"serial_number"`
	Subject      Subject `json:"subject"`
	Valid        int     `json:"valid_year"`
	KeyLength    int     `json:"key_length"`
}

type WebCertConfig struct {
	SerialNumber int64   `json:"serial_number"`
	Subject      Subject `json:"subject"`
	DNSName      string  `json:"dns_name"`
	Valid        int     `json:"valid_year"`
	KeyLength    int     `json:"key_length"`
}

type CertJSON struct {
	CA  CaCertConfig  `json:"ca"`
	Web WebCertConfig `json:"web"`
}
