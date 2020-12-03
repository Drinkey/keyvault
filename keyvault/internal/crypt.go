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

func aesGcmCipher(key string) cipher.AEAD {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	return aesgcm
}

// TODO: implement the master key generation
func GenerateMasterKey() string {
	// keyLen := 24
	return "passphrasewhichneedstobe32bytes!"
}

func Sha256Sum(input string) string {
	shaSum := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", shaSum)
}

func GenerateNonce() []byte {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	return nonce
}

func EncodeByte(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeString(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// encrypt the plaintext using key and nonce
// Return encrypted in byte
func Encrypt(plaintext string, key string, nonce []byte) []byte {
	aesgcm := aesGcmCipher(key)
	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	return ciphertext
}

// Decrypt the ciphertext in byte using key and nonce
// Return the plain text
func Decrypt(ciphertext []byte, key string, nonce []byte) string {
	aesgcm := aesGcmCipher(key)
	plaintext, err := aesgcm.Open(nil, nonce, []byte(ciphertext), nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
