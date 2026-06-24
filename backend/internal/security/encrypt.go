package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"parser/internal/config"
)

func deriveKey(password string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hasher.Sum(nil)
}

func EncryptPayload(plaintext string) (string, error) {
	secretKey := cfg.Security.SecretKey
	key := deriveKey(secretKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM block: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to read random nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	result := make([]byte, len(nonce)+len(ciphertext))
	copy(result, nonce)
	copy(result[len(nonce):], ciphertext)

	return hex.EncodeToString(result), nil
}
