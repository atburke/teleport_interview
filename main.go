package main

import (
	"github.com/atburke/teleport_interview/internal/database"
	"github.com/atburke/teleport_interview/internal/server"
	"log"
	"os"
	"strings"
)

func getEnvironment() (*server.Env, error) {
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

	db, err := database.NewMySqlDatabase(
		strings.TrimSpace(string(dbUsername)),
		strings.TrimSpace(string(dbPassword)),
		strings.TrimSpace(string(dbName)),
	)
	if err != nil {
		return nil, err
	}
	return server.NewEnv(db, "./internal/server/web/"), nil
}

func main() {
	env, err := getEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Run(env))
}
