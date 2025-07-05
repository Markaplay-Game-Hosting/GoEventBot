package crypter

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type Crypt struct {
	AesGCM cipher.AEAD
}

func New(secret []byte) (*Crypt, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return &Crypt{AesGCM: aesGCM}, nil
}

func (c Crypt) Encrypt(plainText string) string {
	return base64.StdEncoding.EncodeToString(c.AesGCM.Seal(make([]byte, c.AesGCM.NonceSize()), make([]byte, c.AesGCM.NonceSize()), []byte(plainText), nil))
}

func (c Crypt) Decrypt(cipherText string) (string, error) {
	// Decode from base64 first
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	dataString := string(data)

	nonceSize := c.AesGCM.NonceSize()
	if len(dataString) < nonceSize {
		return "", errors.New("cipherText too short")
	}

	nonce, cipherTextBytes := dataString[:nonceSize], dataString[nonceSize:]

	plaintext, err := c.AesGCM.Open(nil, []byte(nonce), []byte(cipherTextBytes), nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
