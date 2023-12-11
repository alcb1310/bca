package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"bca-go-final/internal/database"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	DB   database.Service
}

type contextPayload struct {
	Id         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	CompanyId  uuid.UUID `json:"company_id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	IsLoggedIn bool      `json:"is_logged_in"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		DB:   database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func getMyPaload(r *http.Request) (contextPayload, error) {
	ctx := r.Context()
	val := ctx.Value("token")

	x, ok := val.([]byte)
	if !ok {
		return contextPayload{}, errors.New("Unable to load context")
	}
	var p contextPayload
	err := json.Unmarshal(x, &p)
	if err != nil {
		return contextPayload{}, errors.New("Unable to parse context")
	}
	return p, nil
}
