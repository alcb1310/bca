package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alcb1310/bca/internals/types"
)

// Handler is a type alias for a function that handles an HTTP request that returns an error.
type Handler func(w http.ResponseWriter, r *http.Request) error

// handleErrors returns a handler that calls the given handler and handles errors.
func handleErrors(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if e, ok := err.(*pgconn.PgError); ok {
				if e.Code == "23505" {
					if strings.Contains(e.Message, "company") {
						w.WriteHeader(http.StatusConflict)
						w.Write([]byte(fmt.Sprintf("<p>Company already exists: %s</p>", err.Error())))
						return
					}

					if strings.Contains(e.Message, "user") {
						w.WriteHeader(http.StatusConflict)
						w.Write([]byte(fmt.Sprintf("<p>User already exists: %s</p>", err.Error())))
						return
					}
				}
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("<p>Unknown error %s</p>", err.Error())))
			return
		}
	}
}

func renderPage(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	return c.Render(r.Context(), w)
}

func getUserFromContext(r *http.Request) (types.User, error) {
	ctx := r.Context()
	ctxValue := ctx.Value("user")

	user, ok := ctxValue.(types.User)
	if !ok {
		return user, errors.New("Invalid context")
	}

	return user, nil
}
