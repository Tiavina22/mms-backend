package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"mms-backend/config"
)

// getEncryptionKey derives a 32-byte key from the config
func getEncryptionKey() []byte {
	key := config.AppConfig.Security.EncryptionKey
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

// Encrypt encrypts plain text using AES-256-GCM
func Encrypt(plainText string) (string, error) {
	if plainText == "" {
		return "", errors.New("plain text cannot be empty")
	}

	key := getEncryptionKey()
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	
	// Encode to base64 for storage
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts cipher text using AES-256-GCM
func Decrypt(cipherText string) (string, error) {
	if cipherText == "" {
		return "", errors.New("cipher text cannot be empty")
	}

	key := getEncryptionKey()
	
	// Decode from base64
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("cipher text too short")
	}

	nonce, cipherTextBytes := data[:nonceSize], data[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

