package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func EncryptString(src string, key []byte) (encryptedString []byte, err error) {
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

	return gcm.Seal(nonce, nonce, []byte(src), nil), nil
}

func DecryptString(src string, key []byte) (decryptedString []byte, err error) {
	cphr, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, fmt.Errorf("aes new chipher err: %v", err)
	}

	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		return []byte{}, fmt.Errorf("aes chipher err: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len([]byte(src)) < nonceSize {
		return []byte{}, fmt.Errorf("nonceSize does not match src length")
	}

	nonce, encryptedMessage := []byte(src)[:nonceSize], []byte(src)[nonceSize:]
	return gcm.Open(nil, nonce, encryptedMessage, nil)
}
