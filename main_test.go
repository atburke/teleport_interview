package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"
)

// TODO: see if there's a better home for MockDatabase
type MockDatabase struct {
	session *Session
	account *Account
	error   error
}

func (m *MockDatabase) CreatePreAuthSession(csrfToken string, initTime time.Time) (*Session, error) {
	return m.session, m.error
}

func (m *MockDatabase) CreateSession(accountId, csrfToken string, initTime time.Time) (*Session, error) {
	return m.session, m.error
}

func (m *MockDatabase) FetchSession(sessionToken string) (*Session, error) {
	return m.session, m.error
}

func (m *MockDatabase) DeleteSession(sessionToken string) error {
	return m.error
}

func (m *MockDatabase) FetchAccount(email string) (*Account, error) {
	return m.account, m.error
}

func (m *MockDatabase) DeleteExpiredSessions(t time.Time) error {
	return m.error
}

func (m *MockDatabase) Close() {}

func TestPing(t *testing.T) {
	router := setupRouter(&Env{&MockDatabase{}})
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestIndex(t *testing.T) {
	fakeToken := "00112233445566778899aabbccddeeff"
	db := MockDatabase{session: &Session{SessionToken: fakeToken}}
	env := &Env{&db}
	router := setupRouter(env)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/index.html", nil)
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Regexp(t, regexp.MustCompile(`window\.csrfToken = '\w{32}';`), w.Body.String())

	cookies := w.Result().Cookies()
	require.Equal(t, 1, len(cookies))
	sessionCookie := cookies[0]
	assert.Equal(t, "session_token", sessionCookie.Name)
	assert.Equal(t, fakeToken, sessionCookie.Value)
}
