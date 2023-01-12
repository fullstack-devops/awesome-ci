package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func EncryptString(src string, key []byte) (encryptedString string, err error) {
	cphr, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	token, err := gcm.Seal(nonce, nonce, []byte(src), nil), nil

	return string(token), err
}

func DecryptString(src string, key []byte) (decryptedString string, err error) {
	cphr, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes new chipher err: %v", err)
	}

	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return "", fmt.Errorf("aes chipher err: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len([]byte(src)) < nonceSize {
		return "", fmt.Errorf("nonceSize does not match src length")
	}

	nonce, encryptedMessage := []byte(src)[:nonceSize], []byte(src)[nonceSize:]
	token, err := gcm.Open(nil, nonce, encryptedMessage, nil)

	return string(token), err
}
