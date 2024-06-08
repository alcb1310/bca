package server

import (
	"net/http"

	"github.com/a-h/templ"
)

// Handler is a type alias for a function that handles an HTTP request that returns an error.
type Handler func(w http.ResponseWriter, r *http.Request) error

// handleErrors returns a handler that calls the given handler and handles errors.
func handleErrors(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func renderPage(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	return c.Render(r.Context(), w)
}
