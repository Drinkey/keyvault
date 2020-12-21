package e

// Msg defines the error code to text message mapping
var Msg = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "Invalid Parameters",
	NOT_AUTHORIED:  "Client Not Authorized",
	NOT_FOUND:      "Resource Not Found",

	DECODING_ERROR: "Decoding Error",
}

// GetMsg returns the message to corresponding error code
func GetMsg(code int) (m string) {
	m, ok := Msg[code]
	if !ok {
		return Msg[ERROR]
	}
	return m
}
