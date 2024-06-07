package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Service struct {
	Logger *slog.Logger
	Router *chi.Mux
	// DB     database.Service
}

func New(logger *slog.Logger) *Service {
	s := &Service{
		Logger: logger,
		Router: chi.NewRouter(),
	}

	// INFO: Define all middlewares before mounting the handlers
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.Router.Use(middleware.CleanPath)

	// INFO: Mount all the handlers
	s.Router.Mount("/public", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))

	s.Router.Get("/", handleErrors(s.Home))
	s.MountHandlers()

	return s
}

func (s *Service) MountHandlers() {
	// s.Router.Get("/", handleErrors(s.Home))
	//
	// s.Router.Post("/login", handleErrors(s.Login))
	// s.Router.Post("/register", handleErrors(s.CreateCompany))
	// s.Router.Get("/register", handleErrors(s.Register))
	//
	// s.Router.Group(func(r chi.Router) {
	// 	sr := &BCAService{Service: *s, Router: r}
	//
	// 	sr.Router.Use(sr.AuthMiddleware)
	//
	// 	sr.Router.Get("/bca", handleErrors(sr.BCAHome))
	// 	sr.Router.Get("/logout", handleErrors(sr.Logout))
	// })
}
