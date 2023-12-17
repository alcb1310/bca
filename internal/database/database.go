package database

import (
	"bca-go-final/internal/types"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
	CreateCompany(company *types.CompanyCreate) error
	Login(l *types.Login) (string, error)
	IsLoggedIn(token string, user uuid.UUID) bool

	// database/users.go
	GetAllUsers(companyId uuid.UUID) ([]types.User, error)
	CreateUser(u *types.UserCreate) (types.User, error)
	GetUser(id, companyId uuid.UUID) (types.User, error)
	UpdateUser(u types.User, id, companyId uuid.UUID) (types.User, error)
	UpdatePassword(pass string, id, companyId uuid.UUID) (types.User, error)
	DeleteUser(id, companyId uuid.UUID) error

	// database/projects.go
	GetAllProjects(companyId uuid.UUID) ([]types.Project, error)
	CreateProject(p types.Project) (types.Project, error)
	GetProject(id, companyId uuid.UUID) (types.Project, error)
	UpdateProject(p types.Project, id, companyId uuid.UUID) error

	// database/suppliers.go
	GetAllSuppliers(companyId uuid.UUID) ([]types.Supplier, error)
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	s := &service{db: db}

	if err := createTables(db); err != nil {
		log.Fatalf(fmt.Sprintf("error creating tables. Err: %v", err))
	}

	if err := loadRoles(db); err != nil {
		log.Fatalf(fmt.Sprintf("error loading roles. Err: %v", err))
	}

	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
