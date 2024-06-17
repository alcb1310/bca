package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/alcb1310/bca/internals/database"
)

type Service struct {
	Logger *slog.Logger
	Router *chi.Mux
	DB     database.DatabaseService
}

type BCAService struct {
	Service

	Router chi.Router
}

func New(
	logger *slog.Logger,
	db database.DatabaseService,
) *Service {
	s := &Service{
		Logger: logger,
		Router: chi.NewRouter(),
		DB:     db,
	}

	// INFO: Define all middlewares before mounting the handlers
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.Router.Use(middleware.CleanPath)

	// INFO: Mount all the handlers
	s.Router.Mount("/public", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))

	s.Router.Get("/", handleErrors(s.Home))
	s.Router.Get("/register", handleErrors(s.Register))
	s.Router.Post("/register", handleErrors(s.RegisterForm))
	s.Router.Post("/login", handleErrors(s.Login))
	s.MountHandlers()

	return s
}

func (s *Service) MountHandlers() {
	s.Router.Group(func(r chi.Router) {
		sr := &BCAService{Service: *s, Router: r}
		sr.Router.Use(sr.AuthMiddleware)
		sr.Router.Get("/bca", handleErrors(sr.BCAHome))
		sr.Router.Get("/bca/logout", handleErrors(sr.Logout))

		sr.Router.Get("/bca/proyectos", handleErrors(sr.ProjectsPage))
		sr.Router.Get("/bca/proyectos/table", handleErrors(sr.ProjectsTable))
		sr.Router.Get("/bca/proyectos/form", handleErrors(sr.ProjectsForm))
        sr.Router.Post("/bca/proyectos", handleErrors(sr.CreateProject))
	})
}
