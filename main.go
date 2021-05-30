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
	indexData, err := os.ReadFile("./web/index.html")
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

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

	// could consider templating this properly
	newIndexData := strings.Replace(
		string(indexData),
		"window.csrfToken = '';",
		fmt.Sprintf("window.csrfToken = '%s';", csrfToken),
		1,
	)

	c.Header("Content-Type", "text/html")
	sessionCookie := http.Cookie{
		Name:     "session_token",
		Value:    session.SessionToken,
		Secure:   true,
		HttpOnly: true,
		// Same-site is lax by default
	}

	// Gin has a slightly weirder interface for this for some reason
	http.SetCookie(c.Writer, &sessionCookie)
	c.String(http.StatusOK, newIndexData)
}

func setupRouter(env *Env) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", env.ping)
	router.GET("/index.html", env.serveIndex)
	router.GET("/", env.serveIndex)

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

	db, err := NewMySqlDatabase(string(dbUsername), string(dbPassword), string(dbName))
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

	router := setupRouter(env)
	log.Fatal(router.RunTLS(
		":8080", "/run/secrets/server-cert.pem", "/run/secrets/server-key.pem",
	))
}
