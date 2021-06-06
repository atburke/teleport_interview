package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/atburke/teleport_interview/internal/crypto"
	"github.com/atburke/teleport_interview/internal/types"
	"github.com/benbjohnson/clock"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// tokens and password hash could be []byte objects, but they're strings in the
// db and in the api, so leaving them as hex strings will save time

// Database represents a database connection.
type Database interface {

	// CreatePreAuthSession creates a new session for a user that has not yet
	// authenticated.
	CreatePreAuthSession(csrfToken string) (*types.Session, error)

	// CreateSession creates a new session for an authenticated user.
	CreateSession(accountId, csrfToken string) (*types.Session, error)

	// FetchSession fetches the session identified by the provided token.
	FetchSession(sessionToken string) (*types.Session, error)

	// DeleteSession deletes a session for a user.
	DeleteSession(sessionToken string) error

	// FetchAccount fetches the account for a user.
	FetchAccount(email string) (*types.Account, error)

	// DeleteExpiredSessions deletes sessions that have expired (after their
	// absolute timeout has passed)
	DeleteExpiredSessions() error

	// Close closes any underlying resources.
	Close()
}

type MySqlDatabase struct {
	driver *sql.DB
	clock  clock.Clock
}

// let's not bother with reusable prepared statements; the performance impact
// should be minimal for us

func NewMySqlDatabase(username, password, databaseName string, clock clock.Clock) (*MySqlDatabase, error) {
	// hardcode non-sensitive info to speed us along
	driver, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(db:3306)/%s?parseTime=true", username, password, databaseName),
	)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database: %w", err)
	}
	return &MySqlDatabase{driver, clock}, nil
}

func (db *MySqlDatabase) CreatePreAuthSession(csrfToken string) (*types.Session, error) {
	// an empty account ID is functionally equivalent to NULL. Consider doing this
	// instead of having account_id nullable?
	return db.CreateSession("", csrfToken)
}

func (db *MySqlDatabase) CreateSession(accountId, csrfToken string) (*types.Session, error) {
	sessionToken, err := crypto.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("Cound not generate token: %w", err)
	}

	// hooray more hardcoding!
	const expireAbsTime = 8 * time.Hour

	stmt := "INSERT INTO Sessions(account_id, session_token, csrf_token, expire_abs) VALUES (?, ?, ?, ?)"

	// TODO: consider explicitly handling case where sessionToken happens to be a duplicate?
	initTime := db.clock.Now()
	_, err = db.driver.Exec(stmt, accountId, sessionToken, csrfToken, initTime.Add(expireAbsTime))
	if err != nil {
		return nil, fmt.Errorf("Error inserting new session: %w", err)
	}

	return &types.Session{accountId, sessionToken, csrfToken, initTime}, nil
}

func (db *MySqlDatabase) FetchSession(sessionToken string) (*types.Session, error) {
	var accountIdRaw sql.NullString
	var accountId, csrfToken string
	var expireAbs time.Time

	stmt := "SELECT account_id, csrf_token, expire_abs FROM Sessions WHERE session_token = ?"

	err := db.driver.QueryRow(
		stmt, sessionToken,
	).Scan(&accountIdRaw, &csrfToken, &expireAbs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Session not found")
		}

		return nil, fmt.Errorf("SQL Error: %w", err)
	}

	if accountIdRaw.Valid {
		accountId = accountIdRaw.String
	}

	return &types.Session{accountId, sessionToken, csrfToken, expireAbs}, nil
}

func (db *MySqlDatabase) DeleteSession(sessionToken string) error {
	stmt := "DELETE FROM Sessions WHERE session_token = ?"
	_, err := db.driver.Exec(stmt, sessionToken)
	if err != nil {
		return fmt.Errorf("Could not delete session: %w", err)
	}
	return nil
}

func (db *MySqlDatabase) FetchAccount(email string) (*types.Account, error) {
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

	return &types.Account{accountId, email, passwordHash, salt}, nil
}

func (db *MySqlDatabase) DeleteExpiredSessions() error {
	stmt := "DELETE FROM Sessions WHERE expire_abs < ?"
	t := db.clock.Now()
	_, err := db.driver.Exec(stmt, t, t)
	if err != nil {
		return fmt.Errorf("SQL Error: %w", err)
	}

	return nil
}

func (db *MySqlDatabase) Close() {
	db.driver.Close()
}
