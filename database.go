package main

import (
	"time"
)

/*
 * tokens and password hash could be []byte objects, but they're strings in the
 * db and in the api, so leaving them as hex strings will save time
 */

// Account represents a user's account and login information.
type Account struct {
	AccountID, Email, PasswordHash, Salt string
}

// Session represent's a user's session.
type Session struct {
	AccountId, SessionToken, CSRFToken string
	ExpireIdle, ExpireAbs              time.Time
}

func (session *Session) Expired(t time.Time) bool {
	return t.After(session.ExpireIdle) || t.After(session.ExpireAbs)
}

// Database represents a database connection.
type Database interface {

	// CreatePreAuthSession creates a new session for a user that has not yet
	// authenticated.
	CreatePreAuthSession(csrfToken string, initTime time.Time) (*Session, error)

	// CreateSession creates a new session for an authenticated user.
	CreateSession(accountId, csrfToken string, initTime time.Time) (*Session, error)

	// FetchSession fetches the session identified by the provided token.
	FetchSession(sessionToken string) (*Session, error)

	// DeleteSession deletes a session for a user.
	DeleteSession(sessionToken string) error

	// FetchAccount fetches the account for a user.
	FetchAccount(email string) (*Account, error)
}
