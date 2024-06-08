package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type DatabaseService interface{}

type service struct {
	DB *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func CreateConnection() DatabaseService {
	db := &service{}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	d, err := sql.Open("pgx", connStr)
	if err != nil {
		slog.Error("Error connecting to the database", "error", err)
		os.Exit(1)
	}
	db.DB = d

	if err := db.DB.Ping(); err != nil {
		slog.Error("Error connecting to the database", "error", err)
		os.Exit(1)
	}

	slog.Debug("Connected to database", "name", database)

	return &db
}
