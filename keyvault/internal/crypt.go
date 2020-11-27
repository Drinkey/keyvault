package internal

import "fmt"

const KeyMask = "******"

// TODO: implement the master key generation
func GenerateMasterKey() string {
	// keyLen := 24
	return "xDeifu-fkeI19-vs313dR"
}

// TODO: implement the encryption
// encrypt the text using key. Encode with base64 and return
func Encrypt(text string, key string) string {
	return fmt.Sprintf("base64.Encode(Encrypt(%s, %s))", text, key)
}

// TODO: implement the decryption
// Decode with base64, decrypt the cipher using key and return the plain text
func Decrypt(cipher string, key string) string {
	return fmt.Sprintf("decrypt(base64.Decode(%s), %s)", cipher, key)
}
