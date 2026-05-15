// Package crypto provides encryption and decryption utilities for envcrypt.
// It uses AES-256-GCM for authenticated encryption of environment files.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// KeySize is the required AES-256 key size in bytes.
	KeySize = 32
	// SaltSize is the size of the random salt used in key derivation.
	SaltSize = 16
	// NonceSize is the GCM nonce size in bytes.
	NonceSize = 12
	// PBKDF2Iterations is the number of iterations for key derivation.
	PBKDF2Iterations = 100_000
)

// ErrInvalidKey is returned when the provided key has an invalid length.
var ErrInvalidKey = errors.New("invalid key: must be 32 bytes for AES-256")

// ErrDecryptionFailed is returned when decryption or authentication fails.
var ErrDecryptionFailed = errors.New("decryption failed: invalid key or corrupted ciphertext")

// GenerateKey creates a cryptographically secure random 32-byte key
// suitable for use with AES-256-GCM encryption.
func GenerateKey() ([]byte, error) {
	key := make([]byte, KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

// DeriveKey derives a 32-byte AES key from a passphrase and salt using PBKDF2-SHA256.
// If salt is nil, a new random salt is generated and returned.
func DeriveKey(passphrase string, salt []byte) (key []byte, usedSalt []byte, err error) {
	if salt == nil {
		salt = make([]byte, SaltSize)
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			return nil, nil, fmt.Errorf("failed to generate salt: %w", err)
		}
	}
	key = pbkdf2.Key([]byte(passphrase), salt, PBKDF2Iterations, KeySize, sha256.New)
	return key, salt, nil
}

// Encrypt encrypts plaintext using AES-256-GCM with the provided key.
// The returned ciphertext includes the random nonce prepended to the encrypted data.
// The key must be exactly 32 bytes.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext produced by Encrypt using AES-256-GCM.
// It expects the nonce to be prepended to the ciphertext.
// The key must be exactly 32 bytes.
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(ciphertext) < NonceSize+gcm.Overhead() {
		return nil, ErrDecryptionFailed
	}

	nonce, ciphertext := ciphertext[:NonceSize], ciphertext[NonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

// EncodeKey encodes a raw key as a base64 URL-safe string for storage.
func EncodeKey(key []byte) string {
	return base64.URLEncoding.EncodeToString(key)
}

// DecodeKey decodes a base64 URL-safe encoded key string into raw bytes.
func DecodeKey(encoded string) ([]byte, error) {
	key, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}
	if len(key) != KeySize {
		return nil, ErrInvalidKey
	}
	return key, nil
}
