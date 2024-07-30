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
	s := &Server{
		DB:     db,
		Router: chi.NewRouter(),
	}
	s.Router.Use(middleware.Logger)
	s.Router.Use(s.authVerify)

	s.Router.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))

	s.Router.Get("/", s.HelloWorldHandler)

	s.Router.Get("/login", s.DisplayLogin)
	s.Router.Post("/login", s.LoginView) // TODO: test form validation
	s.RegisterRoutes()
  // TODO: properly migrate to chi router

	return s
}
