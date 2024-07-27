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
	databasePort = os.Getenv("DB_PORT")
	host         = os.Getenv("DB_HOST")
	port         = os.Getenv("PORT")
	secretKey    = os.Getenv("SECRET")
)

func init() {
	if databaseName == "" {
		panic("DB_DATABASE must be set")
	}
	if password == "" {
		panic("DB_PASSWORD must be set")
	}
	if username == "" {
		panic("DB_USERNAME must be set")
	}
	if databasePort == "" {
		panic("DB_PORT must be set")
	}
	if host == "" {
		panic("DB_HOST must be set")
	}
	if port == "" {
		panic("PORT must be set")
	}
	if secretKey == "" || len(secretKey) < 8 {
		panic("SECRET must be set and of at least 8 characters")
	}
}

func main() {
	db := database.New(databaseName, username, password, host, databasePort)
	server := server.NewServer(db)

	err := server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
