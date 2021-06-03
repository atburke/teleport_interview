package crypto

import (
	"fmt"
	"github.com/atburke/teleport_interview/internal/types"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMatchingCSRFTokens(t *testing.T) {
	token := "4233af9dc30344cd"
	session := types.Session{CSRFToken: token}
	assert.True(t, IsSessionOwner(&session, token))
}

func TestNotMatchingCSRFTokens(t *testing.T) {
	token1 := "4233af9dc30344cd"
	token2 := "2433af9dc30344cd"
	session := types.Session{CSRFToken: token1}
	assert.False(t, IsSessionOwner(&session, token2))
}

func TestCorrectPassword(t *testing.T) {
	password := "mypassword"
	salt := "d7c7dd775f746f67f76ded1cedc7b57f"
	expectedHash := GenHash(password, salt)
	account := types.Account{PasswordHash: expectedHash, Salt: salt}
	assert.True(t, IsCorrectPassword(&account, password))
}

func TestNotCorrectPassword(t *testing.T) {
	password := "mypassword"
	salt := "d7c7dd775f746f67f76ded1cedc7b57f"
	expectedHash := GenHash(password, salt)
	account := types.Account{PasswordHash: expectedHash, Salt: salt}
	assert.False(t, IsCorrectPassword(&account, "mypasswordd"))
}

func TestMain(m *testing.M) {
	// this is a more reliable way to get our test hash than with the web version
	// I was using
	fmt.Println("admin@example.com:sneakyadminpassword")
	fmt.Printf("Hash is %s\n", GenHash("sneakyadminpassword", "8bc78e90a114942e38ee62a89b2f22cf"))
	os.Exit(m.Run())
}
