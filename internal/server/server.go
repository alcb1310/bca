package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"bca-go-final/internal/database"
)

type Server struct {
	DB     database.Service
	Router *chi.Mux
}

func NewServer(db database.Service) *Server {
	NewServer := &Server{
		DB:     db,
		Router: chi.NewRouter(),
	}
  NewServer.Router.Use(middleware.Logger)

  NewServer.Router.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))

	return NewServer
}
