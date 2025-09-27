package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	data := append(plaintext, key...)

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func Decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]

	data, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, err
	}

	if len(data) < len(key) {
		return nil, fmt.Errorf("decrypted data too short")
	}
	plaintext, embeddedKey := data[:len(data)-len(key)], data[len(data)-len(key):]

	if !bytes.Equal(key, embeddedKey) {
		return nil, fmt.Errorf("key mismatch: wrong key used for decryption")
	}

	return plaintext, nil
}

func SHA256(data string) (hash []byte) {
	h := sha256.New()
	h.Write([]byte(data))
	return h.Sum(nil)
}
