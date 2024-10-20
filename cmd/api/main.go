package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	timezone     = os.Getenv("TIMEZONE")
)

var intTimezone int

func init() {
	environmentValidation()
	loggerSetup(env)
}

func main() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, databasePort, databaseName)
	db := database.New(connStr)
	server := server.NewServer(db, secretKey, intTimezone)

	slog.Info("Listening on port", "port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), server.Router); err != nil {
		panic(err)
	}
}

func loggerSetup(env string) {
	handlerOptions := &slog.HandlerOptions{}

	switch strings.ToLower(env) {
	case "debug":
		handlerOptions.Level = slog.LevelDebug
	case "info":
		handlerOptions.Level = slog.LevelInfo
	case "warn":
		handlerOptions.Level = slog.LevelWarn
	default:
		handlerOptions.Level = slog.LevelError
	}

	loggerHandler := slog.NewJSONHandler(os.Stdout, handlerOptions)
	slog.SetDefault(slog.New(loggerHandler))
}

func environmentValidation() {
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

	if timezone == "" {
		panic("TIMEZONE must be set")
	}

	var err error
	intTimezone, err = strconv.Atoi(timezone)
	if err != nil {
		panic("TIMEZONE must be an integer")
	}
	if intTimezone < -12 || intTimezone > 12 {
		panic("TIMEZONE must be between -12 and 12")
	}
}
