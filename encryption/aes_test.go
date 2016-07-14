package main

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	value := "my-value"

	aesKey := generateAesKey()
	encrypted, err := aesEncrypt(aesKey, []byte(value))
	if err != nil {
		t.Fatal("Test failed during encrypting using AES.")
	}

	decrypted, err := aesDecrypt(aesKey, encrypted)
	if err != nil {
		t.Fatal("Test failed during decrypting using AES.")
	}

	if value != string(decrypted) {
		t.Fatalf("Encrypted and decrypted value is not the same: actual '%s', expected '%s'", decrypted, value)
	}
}
