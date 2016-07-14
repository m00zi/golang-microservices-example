package main

import (
	"crypto/aes"
	"io"
	"crypto/cipher"
	"errors"
	"crypto/rand"
	"log"
)

const (
	// 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes
	cipherKeyBytesCount = 32
)

func aesEncrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize + len(string(text)))

	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], text)
	return ciphertext, nil
}

func aesDecrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	plaintext := ciphertext
	return plaintext, nil
}

func generateAesKey() []byte {
	key := make([]byte, cipherKeyBytesCount)

	_, err := rand.Read(key)
	if err != nil {
		log.Panicf("Error during AES key generation: %+v", key)
	}

	return key
}


