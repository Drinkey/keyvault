package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

const KeyMask = "******"

func aesGcmCipher(key []byte) cipher.AEAD {
	// keyBytes := []byte(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	return aesgcm
}

func GenerateRandomKey(l int) (k []byte) {
	k = make([]byte, l)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		panic(err.Error())
	}
	return
}

func GenerateMasterKey() []byte {
	return GenerateRandomKey(32)
}

func Sha256Sum(input string) string {
	shaSum := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", shaSum)
}

func GenerateNonce() []byte {
	return GenerateRandomKey(12)
}

func EncodeByte(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// encrypt the plaintext using key and nonce
// Return encrypted in byte
func Encrypt(plaintext string, key []byte, nonce []byte) []byte {
	aesgcm := aesGcmCipher(key)
	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	return ciphertext
}

// Decrypt the ciphertext in byte using key and nonce
// Return the plain text
func Decrypt(ciphertext []byte, key []byte, nonce []byte) string {
	aesgcm := aesGcmCipher(key)
	plaintext, err := aesgcm.Open(nil, nonce, []byte(ciphertext), nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
