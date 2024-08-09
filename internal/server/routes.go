package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"bca-go-final/internal/views"
)

func (s *Server) RegisterRoutes(r chi.Router) http.Handler {
	r.HandleFunc("/api/login", s.Login)
	r.HandleFunc("/api/register", s.Register)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	component := views.WelcomeView()
	component.Render(r.Context(), w)
}
