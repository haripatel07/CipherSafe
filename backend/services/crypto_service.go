package services

import (
	"ciphersafe/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Encrypt encrypts plaintext using AES-GCM with a random nonce.
// The output is hex-encoded "nonce||ciphertext".
func Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(config.AppConfig.MasterEncryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a new nonce for every encryption
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Seal encrypts and authenticates the plaintext.
	// We pass nil for additionalData.
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Prepend the nonce to the ciphertext
	fullCiphertext := append(nonce, ciphertext...)

	return hex.EncodeToString(fullCiphertext), nil
}

// Decrypt decrypts hex-encoded "nonce||ciphertext" string.
func Decrypt(hexCiphertext string) (string, error) {
	ciphertext, err := hex.DecodeString(hexCiphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(config.AppConfig.MasterEncryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Split the nonce and the actual ciphertext
	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Open decrypts and authenticates the ciphertext
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}
