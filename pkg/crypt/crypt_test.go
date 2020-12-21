package crypt

import (
	"encoding/hex"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCryptEncryptionProcess(t *testing.T) {
	Convey("Scenario Encrypt Process", t, func() {
		Convey("Given nonce hex, plain text, encryption key, and expected encoded ciphertext", func() {
			expectEncrypted := "E7sm0JlmEVPFvfwuiFjXIYmjNnbz71DoVpxKoOpNZQ=="
			nonceHex := "3660f5f97eb4180a767731ab"
			text := "试一下中文"
			key := []byte("passphrasewhichneedstobe32bytes!")
			Convey("When decode nonce hex to byte", func() {
				nonce, err := hex.DecodeString(nonceHex)
				Convey("The decoding should have no error", func() {
					So(err, ShouldEqual, nil)
				})
				Convey("When Encrypting the text with nonce and key to ciphertext", func() {
					s := Encrypt(text, key, nonce)
					Convey("And encode ciphertext with base64", func() {
						b64CipherText := EncodeByte(s)
						Convey("The encoded ciphertext should equal to expected string", func() {
							So(b64CipherText, ShouldEqual, expectEncrypted)
						})
					})
				})
			})
		})
	})

}

func TestCryptDecryptionProcess(t *testing.T) {
	Convey("Scenario Decryption Process", t, func() {
		Convey("Given nonce hex string, base64 encoded cipher text and key", func() {
			nonceHex := "3660f5f97eb4180a767731ab"
			b64CipherText := "E7sm0JlmEVPFvfwuiFjXIYmjNnbz71DoVpxKoOpNZQ=="
			key := []byte("passphrasewhichneedstobe32bytes!")
			text := "试一下中文"
			Convey("When decode nonce hex string", func() {
				nonce, err := hex.DecodeString(nonceHex)
				Convey("The nonce hex decoding should have no error", func() {
					So(err, ShouldEqual, nil)
				})
				Convey("When Decode base64 encoded cipher text to byte", func() {
					cipherText, err := DecodeString(b64CipherText)
					Convey("The nonce hex decoding should have no error", func() {
						So(err, ShouldEqual, nil)
					})
					Convey("When Decrypt cipher text to plaintext", func() {
						plaintext := Decrypt(cipherText, key, nonce)
						Convey("The plaintext should equal to expect plaintext", func() {
							So(plaintext, ShouldEqual, text)
						})
					})
				})
			})
		})
	})
}

func TestCryptGenerateRandomKey(t *testing.T) {
	k := GenerateRandomKey(20)
	if len(k) != 20 {
		t.Logf("expected key len=20, actual key len=%d", len(k))
		t.Fail()
	}
}

func TestCryptGenerateMasterKey(t *testing.T) {
	k := GenerateMasterKey()
	if len(k) != 32 {
		t.Logf("expected key len=32, actual key len=%d", len(k))
		t.Fail()
	}
}

func TestCryptGenerateNonce(t *testing.T) {
	k := GenerateNonce()
	if len(k) != 12 {
		t.Logf("expected key len=12, actual key len=%d", len(k))
		t.Fail()
	}
}
