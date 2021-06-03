package crypto

import (
	"fmt"
	"github.com/atburke/teleport_interview/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	expectedHash, err := GenerateHash(password, salt)
	require.Nil(t, err)
	account := types.Account{PasswordHash: expectedHash, Salt: salt}
	isCorrectPassword, err := IsCorrectPassword(&account, password)
	assert.Nil(t, err)
	assert.True(t, isCorrectPassword)
}

func TestNotCorrectPassword(t *testing.T) {
	password := "mypassword"
	salt := "d7c7dd775f746f67f76ded1cedc7b57f"
	expectedHash, err := GenerateHash(password, salt)
	require.Nil(t, err)
	account := types.Account{PasswordHash: expectedHash, Salt: salt}
	isCorrectPassword, err := IsCorrectPassword(&account, "mypasswordd")
	assert.Nil(t, err)
	assert.False(t, isCorrectPassword)
}

func TestBadSalt(t *testing.T) {
	_, err := GenerateHash("anything", "spaghettitime")
	require.NotNil(t, err)
}

func TestMain(m *testing.M) {
	// this is a more reliable way to get our test hash than with the web version
	// I was using
	fmt.Println("admin@example.com:sneakyadminpassword")
	hash, _ := GenerateHash("sneakyadminpassword", "8bc78e90a114942e38ee62a89b2f22cf")
	fmt.Printf("Hash is %s\n", hash)
	os.Exit(m.Run())
}
