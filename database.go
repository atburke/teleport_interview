package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// tokens and password hash could be []byte objects, but they're strings in the
// db and in the api, so leaving them as hex strings will save time

// Account represents a user's account and login information.
type Account struct {
	AccountId, Email, PasswordHash, Salt string
}

// Session represent's a user's session.
type Session struct {
	AccountId, SessionToken, CSRFToken string
	ExpireIdle, ExpireAbs              time.Time
}

// Expired checks if the session has expired.
func (session *Session) Expired(t time.Time) bool {
	return t.After(session.ExpireIdle) || t.After(session.ExpireAbs)
}

// Authenticated returns true if the session is associated with an authenticated
// user, and false if the session is pre-auth.
func (session *Session) Authenticated() bool {
	return session.AccountId != ""
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

	// DeleteExpiredSessions deletes sessions that have expired (after either their
	// idle or absolute timeout has passed)
	DeleteExpiredSessions(t time.Time) error

	// Close closes any underlying resources.
	Close()
}

type MySqlDatabase struct {
	driver *sql.DB
}

// let's not bother with reusable prepared statements; the performance impact
// should be minimal for us

func NewMySqlDatabase(username, password, databaseName string) (*MySqlDatabase, error) {
	// hardcode non-sensitive info to speed us along
	driver, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?parseTime=true", username, password, databaseName),
	)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database: %w", err)
	}
	return &MySqlDatabase{driver}, nil
}

func (db *MySqlDatabase) CreatePreAuthSession(csrfToken string, initTime time.Time) (*Session, error) {
	// an empty account ID is functionally equivalent to NULL. Consider doing this
	// instead of having account_id nullable?
	return db.CreateSession("", csrfToken, initTime)
}

func (db *MySqlDatabase) CreateSession(accountId, csrfToken string, initTime time.Time) (*Session, error) {
	sessionToken, err := GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("Cound not generate token: %w", err)
	}
	stmt := "INSERT INTO Sessions(account_id, session_token, csrf_token, expire_idle, expire_abs) VALUES (?, ?, ?, ?)"

	// TODO: consider explicitly handling case where sessionToken happens to be a duplicate?
	_, err = db.driver.Exec(stmt, accountId, sessionToken, csrfToken, initTime, initTime)
	if err != nil {
		return nil, fmt.Errorf("Error inserting new session: %w", err)
	}

	return &Session{accountId, sessionToken, csrfToken, initTime, initTime}, nil
}

func (db *MySqlDatabase) FetchSession(sessionToken string) (*Session, error) {
	var accountIdRaw sql.NullString
	var accountId, csrfToken string
	var expireIdle, expireAbs time.Time

	stmt := "SELECT account_id, csrf_token, expire_idle, expire_abs FROM Sessions WHERE session_token = ?"

	err := db.driver.QueryRow(
		stmt, sessionToken,
	).Scan(&accountIdRaw, &csrfToken, &expireIdle, &expireAbs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Session not found")
		}

		return nil, fmt.Errorf("SQL Error: %w", err)
	}

	if accountIdRaw.Valid {
		accountId = accountIdRaw.String
	}

	return &Session{accountId, sessionToken, csrfToken, expireIdle, expireAbs}, nil
}

func (db *MySqlDatabase) DeleteSession(sessionToken string) error {
	stmt := "DELETE FROM Sessions WHERE session_token = ?"
	_, err := db.driver.Exec(stmt, sessionToken)
	if err != nil {
		return fmt.Errorf("Could not delete session: %w", err)
	}
	return nil
}

func (db *MySqlDatabase) FetchAccount(email string) (*Account, error) {
	var accountId, passwordHash, salt string
	stmt := "SELECT account_id, password_hash, salt FROM Accounts WHERE email = ?"
	err := db.driver.QueryRow(
		stmt, email,
	).Scan(&accountId, &passwordHash, &salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("No matching account found")
		}
		return nil, fmt.Errorf("SQL error: %w", err)
	}

	return &Account{accountId, email, passwordHash, salt}, nil
}

func (db *MySqlDatabase) DeleteExpiredSessions(t time.Time) error {
	stmt := "DELETE FROM Sessions WHERE expire_idle < ? OR expire_abs < ?"
	_, err := db.driver.Exec(stmt, t, t)
	if err != nil {
		return fmt.Errorf("SQL Error: %w", err)
	}

	return nil
}

func (db *MySqlDatabase) Close() {
	db.driver.Close()
}
