package main

import (
	"errors"
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
	session, authSession       *Session
	account                    *Account
	sessionError, accountError error
}

func (m *MockDatabase) CreatePreAuthSession(csrfToken string, initTime time.Time) (*Session, error) {
	return m.session, m.sessionError
}

func (m *MockDatabase) CreateSession(accountId, csrfToken string, initTime time.Time) (*Session, error) {
	return m.authSession, m.sessionError
}

func (m *MockDatabase) FetchSession(sessionToken string) (*Session, error) {
	return m.session, m.sessionError
}

func (m *MockDatabase) DeleteSession(sessionToken string) error {
	return m.sessionError
}

func (m *MockDatabase) FetchAccount(email string) (*Account, error) {
	return m.account, m.accountError
}

func (m *MockDatabase) DeleteExpiredSessions(t time.Time) error {
	return m.sessionError
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

// reuse these in login tests
const sessionToken = "00112233445566778899aabbccddeeff"
const csrfToken = "ffeeddccbbaa99887766554433221100"
const accountId = "a28f1766-15be-46a7-8132-b267344faa9c"
const email = "test@example.com"
const password = "sneakytime"
const salt = "d7c7dd775f746f67f76ded1cedc7b57f"

func TestLogin(t *testing.T) {
	theFuture := time.Now().AddDate(0, 0, 1)
	token2 := "abcdabcdabcdabcdabcdabcdabcdabcd"
	session := Session{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		ExpireIdle:   theFuture,
		ExpireAbs:    theFuture,
	}
	authSession := Session{
		SessionToken: token2,
		CSRFToken:    csrfToken,
		ExpireIdle:   theFuture,
		ExpireAbs:    theFuture,
	}
	account := Account{
		AccountId:    accountId,
		Email:        email,
		PasswordHash: GenHash(password, salt),
		Salt:         salt,
	}
	db := MockDatabase{session: &session, account: &account, authSession: &authSession}
	env := &Env{&db}
	router := setupRouter(env)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/login", nil)

	sessionCookie := http.Cookie{Name: "session_token", Value: sessionToken}
	request.AddCookie(&sessionCookie)
	request.Header.Set("CSRF", csrfToken)
	request.SetBasicAuth(email, password)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)
	cookies := w.Result().Cookies()
	require.Equal(t, 1, len(cookies))
	cookie := cookies[0]
	assert.Equal(t, "session_token", cookie.Name)
	assert.NotEqual(t, "", cookie.Value)
	// should have generated a new session
	assert.NotEqual(t, sessionToken, cookie.Value)
}

func TestLoginNoAccount(t *testing.T) {
	theFuture := time.Now().AddDate(0, 0, 1)
	session := Session{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		ExpireIdle:   theFuture,
		ExpireAbs:    theFuture,
	}
	db := MockDatabase{session: &session, accountError: errors.New("Account not found")}
	env := &Env{&db}
	router := setupRouter(env)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/login", nil)

	sessionCookie := http.Cookie{Name: "session_token", Value: sessionToken}
	request.AddCookie(&sessionCookie)
	request.Header.Set("CSRF", csrfToken)
	request.SetBasicAuth(email, password)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// should have similar tests for no session/csrf token

func TestLoginBadSessionToken(t *testing.T) {
	theFuture := time.Now().AddDate(0, 0, 1)
	session := Session{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		ExpireIdle:   theFuture,
		ExpireAbs:    theFuture,
	}
	account := Account{
		AccountId:    accountId,
		Email:        email,
		PasswordHash: GenHash(password, salt),
		Salt:         salt,
	}
	db := MockDatabase{
		session:      &session,
		account:      &account,
		sessionError: errors.New("Bad session token"),
	}

	env := &Env{&db}
	router := setupRouter(env)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/login", nil)

	badToken := "55443344554433445544334455443344"

	sessionCookie := http.Cookie{Name: "session_token", Value: badToken}
	request.AddCookie(&sessionCookie)
	request.Header.Set("CSRF", csrfToken)
	request.SetBasicAuth(email, password)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginBadCSRF(t *testing.T) {
	theFuture := time.Now().AddDate(0, 0, 1)
	session := Session{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		ExpireIdle:   theFuture,
		ExpireAbs:    theFuture,
	}
	account := Account{
		AccountId:    accountId,
		Email:        email,
		PasswordHash: GenHash(password, salt),
		Salt:         salt,
	}
	db := MockDatabase{session: &session, account: &account}
	env := &Env{&db}
	router := setupRouter(env)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/login", nil)

	badToken := "55443344554433445544334455443344"

	sessionCookie := http.Cookie{Name: "session_token", Value: sessionToken}
	request.AddCookie(&sessionCookie)
	request.Header.Set("CSRF", badToken)
	request.SetBasicAuth(email, password)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginBadPassword(t *testing.T) {
	theFuture := time.Now().AddDate(0, 0, 1)
	session := Session{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		ExpireIdle:   theFuture,
		ExpireAbs:    theFuture,
	}
	account := Account{
		AccountId:    accountId,
		Email:        email,
		PasswordHash: GenHash(password, salt),
		Salt:         salt,
	}
	db := MockDatabase{session: &session, account: &account}
	env := &Env{&db}
	router := setupRouter(env)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/login", nil)

	badPassword := "hackslol"

	sessionCookie := http.Cookie{Name: "session_token", Value: sessionToken}
	request.AddCookie(&sessionCookie)
	request.Header.Set("CSRF", csrfToken)
	request.SetBasicAuth(email, badPassword)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogout(t *testing.T) {
	theFuture := time.Now().AddDate(0, 0, 1)
	session := Session{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		ExpireIdle:   theFuture,
		ExpireAbs:    theFuture,
	}
	db := MockDatabase{session: &session}
	env := &Env{&db}
	router := setupRouter(env)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/logout", nil)

	sessionCookie := http.Cookie{Name: "session_token", Value: sessionToken}
	request.AddCookie(&sessionCookie)
	request.Header.Set("CSRF", csrfToken)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogoutNotLoggedIn(t *testing.T) {
	db := MockDatabase{sessionError: errors.New("No session")}
	env := &Env{&db}
	router := setupRouter(env)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/logout", nil)

	sessionCookie := http.Cookie{Name: "session_token", Value: sessionToken}
	request.AddCookie(&sessionCookie)
	request.Header.Set("CSRF", csrfToken)

	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogoutBadCSRF(t *testing.T) {
	theFuture := time.Now().AddDate(0, 0, 1)
	session := Session{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		ExpireIdle:   theFuture,
		ExpireAbs:    theFuture,
	}
	db := MockDatabase{session: &session}
	env := &Env{&db}
	router := setupRouter(env)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/logout", nil)

	badToken := "00558833775588449933775566448855"

	sessionCookie := http.Cookie{Name: "session_token", Value: sessionToken}
	request.AddCookie(&sessionCookie)
	request.Header.Set("CSRF", badToken)

	router.ServeHTTP(w, request)

	// TODO: when we have better mocking, check that session wasn't deleted
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
