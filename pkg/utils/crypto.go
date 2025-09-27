package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
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

func SHA256(data string) [32]byte {
	return sha256.Sum256([]byte(data))
}
func t410fToHex(address string) (string, error) {
	if !(strings.HasPrefix(address, "t410f") || strings.HasPrefix(address, "f410f")) {
		return "", errors.New("address must start with t410f or f410f")
	}

	// Strip prefix (t410f / f410f)
	encoded := strings.ToUpper(address[5:]) // base32 expects uppercase

	// RFC4648 base32 decoder, no padding
	decoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	decoded, err := decoder.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	if len(decoded) <= 4 {
		return "", errors.New("decoded payload too short")
	}

	// Drop last 4 bytes (checksum)
	payload := decoded[:len(decoded)-4]
	if len(payload) != 20 {
		return "", fmt.Errorf("unexpected payload length %d (want 20)", len(payload))
	}

	return "0x" + hex.EncodeToString(payload), nil
}
