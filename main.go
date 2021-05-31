package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Env defines routes, as well as objects that should be shared across requests
// (database connections, config, etc.)
//
// Defining routes like this allows for easy dependency injection in tests.
type Env struct {
	db Database
}

func (env *Env) ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (env *Env) serveIndex(c *gin.Context) {
	now := time.Now()

	csrfToken, err := GenerateToken()
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	session, err := env.db.CreatePreAuthSession(csrfToken, now)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Content-Type", "text/html")
	sessionCookie := http.Cookie{
		Name:     "session_token",
		Value:    session.SessionToken,
		Secure:   true,
		HttpOnly: true,
		// Same-site is lax by default
	}

	// Gin has a slightly weirder interface for cookies for some reason
	http.SetCookie(c.Writer, &sessionCookie)
	c.HTML(http.StatusOK, "index.html", gin.H{"csrf": csrfToken})
}

func (env *Env) login(c *gin.Context) {
	now := time.Now()
	csrfToken := c.Request.Header.Get("CSRF")

	sessionToken, err := c.Cookie("session_token")
	// should look into deduplicating this
	// should also look into a logging library w/ levels (debug, info, warning, etc)
	if err != nil {
		log.Println("Missing session token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	session, err := env.db.FetchSession(sessionToken)
	if err != nil {
		log.Printf("Error fetching session: %v\n", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if session.Expired(now) {
		log.Println("Session expired")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !IsSessionOwner(session, csrfToken) {
		log.Println("Not owner of session")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if session.Authenticated() {
		log.Println("Already logged in")
		c.AbortWithStatus(http.StatusOK)
		return
	}

	email, password, ok := c.Request.BasicAuth()
	if !ok {
		log.Println("Missing basic auth")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	account, err := env.db.FetchAccount(email)
	if err != nil {
		log.Println("Missing account")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !IsCorrectPassword(account, password) {
		log.Println("Bad password")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// both session tuples should only be accessed by this user, so we won't bother
	// with a transaction
	authenticatedSession, err := env.db.CreateSession(account.AccountId, csrfToken, now)
	if err != nil {
		log.Printf("Error creating session: %v\n", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = env.db.DeleteSession(sessionToken)
	if err != nil {
		log.Printf("Error deleting old session: %v\n", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sessionCookie := http.Cookie{
		Name:     "session_token",
		Value:    authenticatedSession.SessionToken,
		Secure:   true,
		HttpOnly: true,
		// Same-site is lax by default
	}

	http.SetCookie(c.Writer, &sessionCookie)
	c.AbortWithStatus(http.StatusOK)
}

func (env *Env) logout(c *gin.Context) {
	// TODO: move logged in check to its own function
	now := time.Now()
	csrfToken := c.Request.Header.Get("CSRF")

	sessionToken, err := c.Cookie("session_token")
	// should look into deduplicating this
	// should also look into a logging library w/ levels (debug, info, warning, etc)
	if err != nil {
		log.Println("Missing session token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	session, err := env.db.FetchSession(sessionToken)
	if err != nil {
		log.Printf("Error fetching session: %v\n", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if session.Expired(now) {
		log.Println("Session expired")
		c.AbortWithStatus(http.StatusOK)
		return
	}

	if !IsSessionOwner(session, csrfToken) {
		log.Println("Not owner of session")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = env.db.DeleteSession(sessionToken)
	if err != nil {
		log.Printf("Error deleting session: %v\n", err)
	}

	c.AbortWithStatus(http.StatusOK)
}

func setupRouter(env *Env) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("web/index.html")
	router.GET("/ping", env.ping)
	router.GET("/index.html", env.serveIndex)
	router.GET("/", env.serveIndex)
	router.POST("/api/login", env.login)
	router.POST("/api/logout", env.logout)

	router.Static("/static", "./web/static")
	rootLevelFiles := []string{"asset-manifest.json", "favicon.ico", "logo192.png", "logo512.png", "manifest.json", "robots.txt"}
	for _, file := range rootLevelFiles {
		router.StaticFile(fmt.Sprintf("/%s", file), fmt.Sprintf("./web/%s", file))
	}

	return router
}

func getEnvironment() (*Env, error) {
	dbUsername, err := os.ReadFile("/run/secrets/mysql_user")
	if err != nil {
		return nil, err
	}

	dbPassword, err := os.ReadFile("/run/secrets/mysql_pw")
	if err != nil {
		return nil, err
	}

	dbName, err := os.ReadFile("/run/secrets/mysql_db")
	if err != nil {
		return nil, err
	}

	db, err := NewMySqlDatabase(strings.TrimSpace(string(dbUsername)), strings.TrimSpace(string(dbPassword)), strings.TrimSpace(string(dbName)))
	if err != nil {
		return nil, err
	}
	return &Env{db}, nil
}

func main() {
	env, err := getEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		tick := time.NewTicker(time.Minute)
		for {
			now := <-tick.C
			log.Println("Clearing expired sessions")
			err := env.db.DeleteExpiredSessions(now)
			if err != nil {
				log.Println("Error while cleaning up expired sessions: %v\n", err)
			}
		}
	}()

	router := setupRouter(env)
	log.Fatal(router.RunTLS(
		":8080", "/run/secrets/server-cert.pem", "/run/secrets/server-key.pem",
	))
}
