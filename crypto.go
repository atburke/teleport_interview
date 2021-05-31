package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"golang.org/x/crypto/argon2"
)

// TODO: use bytes everywhere, rather than constantly switching between strings
// and bytes

// GenHash generates a password hash from a password and salt. The salt and
// hash are hex-encoded strings.
func GenHash(password, salt string) string {
	passwordIn := []byte(password)
	saltIn, _ := hex.DecodeString(salt)
	hash := argon2.IDKey(passwordIn, saltIn, 2, 15*1024, 1, 16)
	return hex.EncodeToString(hash)
}

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

// IsSessionOwner checks if a client's csrf token matches the one in its session
// in constant time.
func IsSessionOwner(session *Session, csrfToken string) bool {
	// If we have malformed tokens, this function should return false, so we can
	// ignore decode errors.
	expectedCSRFToken, _ := hex.DecodeString(session.CSRFToken)
	givenCSRFToken, _ := hex.DecodeString(csrfToken)
	return subtle.ConstantTimeCompare(expectedCSRFToken, givenCSRFToken) == 1
}

// IsCorrectPassword checks if a client's provided password matches the password
// for the account, in constant time.
func IsCorrectPassword(account *Account, password string) bool {
	expectedHash, _ := hex.DecodeString(account.PasswordHash)
	hash, _ := hex.DecodeString(GenHash(password, account.Salt))
	return subtle.ConstantTimeCompare(expectedHash, hash) == 1
}
