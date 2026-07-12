package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

func DecryptPayload(ciphertextHex string, secretKey string) (string, error) {
	if len(ciphertextHex) < 24 {
		return ciphertextHex, nil
	}

	rawBytes, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return ciphertextHex, nil
	}

	key := deriveKey(secretKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM block: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(rawBytes) < nonceSize {
		return ciphertextHex, nil
	}

	nonce := rawBytes[:nonceSize]
	ciphertext := rawBytes[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return ciphertextHex, nil
	}

	return string(plaintext), nil
}
