package services

import (
	"ciphersafe/config"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up test environment variables with a 32-byte key
	os.Setenv("MASTER_ENCRYPTION_KEY", "12345678901234567890123456789012") // Exactly 32 bytes
	os.Setenv("JWT_SECRET_KEY", "test-jwt-secret")
	os.Setenv("DATABASE_URL", "host=localhost user=test password=test dbname=test port=5432 sslmode=disable")

	// Load config
	config.LoadConfig()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestEncryptDecrypt(t *testing.T) {
	plaintext := "my-secret-password"

	// Encrypt
	ciphertext, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	if ciphertext == "" {
		t.Fatal("Ciphertext is empty")
	}

	if ciphertext == plaintext {
		t.Fatal("Ciphertext should not equal plaintext")
	}

	// Decrypt
	decrypted, err := Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf("Decrypted text '%s' does not match original '%s'", decrypted, plaintext)
	}
}

func TestDecryptInvalidCiphertext(t *testing.T) {
	invalidCiphertext := "invalid-base64"

	_, err := Decrypt(invalidCiphertext)
	if err == nil {
		t.Fatal("Expected error for invalid ciphertext")
	}
}

func TestDecryptTamperedCiphertext(t *testing.T) {
	plaintext := "test-secret"

	ciphertext, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	// Tamper with the ciphertext
	tampered := ciphertext[:len(ciphertext)-1] + "x"

	_, err = Decrypt(tampered)
	if err == nil {
		t.Fatal("Expected error for tampered ciphertext")
	}
}
