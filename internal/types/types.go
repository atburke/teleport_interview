package types

import (
	"time"
)

// Account represents a user's account and login information.
type Account struct {
	AccountId, Email, PasswordHash, Salt string
}

// Session represent's a user's session.
type Session struct {
	AccountId, SessionToken, CSRFToken string
	ExpireAbs                          time.Time
}

// Expired checks if the session has expired.
func (session *Session) Expired(t time.Time) bool {
	return t.After(session.ExpireAbs)
}

// Authenticated returns true if the session is associated with an authenticated
// user, and false if the session is pre-auth.
func (session *Session) Authenticated() bool {
	return session.AccountId != ""
}
