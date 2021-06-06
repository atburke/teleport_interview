package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"github.com/atburke/teleport_interview/internal/types"
	"golang.org/x/crypto/argon2"
)

// TODO: use bytes everywhere, rather than constantly switching between strings
// and bytes

// GenHash generates a password hash from a password and salt. The salt and
// hash are hex-encoded strings.
func GenerateHash(password, salt string) (string, error) {
	passwordIn := []byte(password)
	saltIn, err := hex.DecodeString(salt)
	if err != nil {
		return "", errors.New("Salt is not a hex-encoded string")
	}
	hash := argon2.IDKey(passwordIn, saltIn, 2, 15*1024, 1, 16)
	return hex.EncodeToString(hash), nil
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
func IsSessionOwner(session *types.Session, csrfToken string) bool {
	// non-hex-encoded tokens are never valid
	expectedCSRFToken, err := hex.DecodeString(session.CSRFToken)
	if err != nil {
		return false
	}
	givenCSRFToken, err := hex.DecodeString(csrfToken)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(expectedCSRFToken, givenCSRFToken) == 1
}

// IsCorrectPassword checks if a client's provided password matches the password
// for the account, in constant time.
func IsCorrectPassword(account *types.Account, password string) bool {
	expectedHash, err := hex.DecodeString(account.PasswordHash)
	// non-hex-encoded hash is never valid
	if err != nil {
		return false
	}

	hashBytes, err := GenerateHash(password, account.Salt)
	if err != nil {
		return false
	}

	// guaranteed to not error, since hashBytes is the result of EncodeToString
	hash, _ := hex.DecodeString(hashBytes)
	return subtle.ConstantTimeCompare(expectedHash, hash) == 1
}
