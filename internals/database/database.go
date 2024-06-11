package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"

	"github.com/alcb1310/bca/internals/types"
)

type DatabaseService interface {
	LoadScript()

	// Creates a new company with its admin user, this function assumes that
	// all data is validated and transformed, will return an error if either the
	// company or the user couldn't be created
	CreateCompany(c *types.Company, u *types.CreateUser) error

	// Logs in to the application, it receives a validated email address as the
	// user's name and a un hashed password, then it will validate that the
	// email - password convination is correct and on success it will login and
	// on failure it will return an "Invalid Credentials" message
	Login(email, password string) (types.User, error)
}

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

	return db
}

func (s *service) LoadScript() {
	data, err := os.OpenFile("./scripts/tables.sql", os.O_RDONLY, 0644)
	if err != nil {
		slog.Error("Unable to open scripts file", "err", err)
		os.Exit(1)
	}
	defer data.Close()

	info, _ := data.Stat()
	bs := make([]byte, info.Size())
	if _, err := bufio.NewReader(data).Read(bs); err != nil {
		slog.Error("Unable to read file", "err", err)
		os.Exit(1)
	}

	tx, err := s.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		slog.Error("Unable to create transaction", "err", err)
		os.Exit(1)
	}
	queries := strings.Split(string(bs), ";")

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			slog.Error("Unable to create tables", "err", err)
			tx.Rollback()
			os.Exit(1)
		}
	}

	q := "select count(*) from role"
	res := 0

	if err := tx.QueryRow(q).Scan(&res); err != nil {
		slog.Error("Could not get roles info", "err", err)
		tx.Rollback()
		os.Exit(1)
	}

	slog.Debug("Total", "res", res)

	if res == 0 {
		q = "insert into role (id, name) values ('a', 'admin')"
		if _, err := tx.Exec(q); err != nil {
			slog.Error("Could not insert the roles", "err", err)
			tx.Rollback()
			os.Exit(1)
		}
	}

	tx.Commit()
	slog.Info("Tables created successfully")

}
