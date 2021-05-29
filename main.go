package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", ping)
	router.Static("/static", "./web/static")
	rootLevelFiles := []string{"asset-manifest.json", "favicon.ico", "index.html", "logo192.png", "logo512.png", "manifest.json", "robots.txt"}
	for _, file := range rootLevelFiles {
		router.StaticFile(fmt.Sprintf("/%s", file), fmt.Sprintf("./web/%s", file))
	}
	router.StaticFile("/", "./web/index.html")
	return router
}

func main() {
	router := setupRouter()
	router.Run()
}
