package main

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/argon2"
	"testing"
)

func genHash(password, salt string) string {
	passwordIn := []byte(password)
	saltIn, _ := hex.DecodeString(salt)
	hash := argon2.IDKey(passwordIn, saltIn, 2, 15*1024, 1, 16)
	return hex.EncodeToString(hash)
}

func TestMatchingCSRFTokens(t *testing.T) {
	token := "4233af9dc30344cd"
	session := Session{CSRFToken: token}
	assert.True(t, IsSessionOwner(&session, token))
}

func TestNotMatchingCSRFTokens(t *testing.T) {
	token1 := "4233af9dc30344cd"
	token2 := "2433af9dc30344cd"
	session := Session{CSRFToken: token1}
	assert.False(t, IsSessionOwner(&session, token2))
}

func TestCorrectPassword(t *testing.T) {
	password := "mypassword"
	salt := "d7c7dd775f746f67f76ded1cedc7b57f"
	expectedHash := genHash(password, salt)
	account := Account{PasswordHash: expectedHash, Salt: salt}
	assert.True(t, IsCorrectPassword(&account, password))
}

func TestNotCorrectPassword(t *testing.T) {
	password := "mypassword"
	salt := "d7c7dd775f746f67f76ded1cedc7b57f"
	expectedHash := genHash(password, salt)
	account := Account{PasswordHash: expectedHash, Salt: salt}
	assert.False(t, IsCorrectPassword(&account, "mypasswordd"))
}
