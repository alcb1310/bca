package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"

	"bca-go-final/internal/database"
	"bca-go-final/internal/server"
)

var (
	databaseName = os.Getenv("DB_DATABASE")
	password     = os.Getenv("DB_PASSWORD")
	username     = os.Getenv("DB_USERNAME")
	dbPort       = os.Getenv("DB_PORT")
	host         = os.Getenv("DB_HOST")
	port         = os.Getenv("PORT")
)

func init() {
	if databaseName == "" {
		panic("DB_DATABASE is not set")
	}

	if password == "" {
		panic("DB_PASSWORD is not set")
	}

	if username == "" {
		panic("DB_USERNAME is not set")
	}

	if dbPort == "" {
		panic("DB_PORT is not set")
	}

	if host == "" {
		panic("DB_HOST is not set")
	}

	if port == "" {
		panic("PORT is not set")
	}

	if os.Getenv("SECRET") == "" {
		panic("SECRET is not set")
	}
}

func main() {
	db := database.New(username, password, host, dbPort, databaseName)
	server, _ := server.NewServer(db)

	err := server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
