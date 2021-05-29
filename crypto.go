package main

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateToken generates a cryptographically-random token.
//
// Since all of our tokens will be the same length, the output length is hard-coded.
func GenerateToken() (string, error) {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}
