package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/alcb1310/bca/internal/database"
	"github.com/alcb1310/bca/internal/server"
)

var (
	databaseName = os.Getenv("DB_DATABASE")
	password     = os.Getenv("DB_PASSWORD")
	username     = os.Getenv("DB_USERNAME")
	databasePort = os.Getenv("DB_PORT")
	host         = os.Getenv("DB_HOST")
	port         = os.Getenv("PORT")
	secretKey    = os.Getenv("SECRET")
	env          = os.Getenv("APP_ENV")
)

func init() {
  environmentValidation()
  loggerSetup(env)
}

func main() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, databasePort, databaseName)
	db := database.New(connStr)
	server := server.NewServer(db, secretKey)

	slog.Info("Listening on port", "port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), server.Router); err != nil {
		panic(err)
	}
}
