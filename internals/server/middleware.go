package server

import (
	"context"
	"log/slog"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"github.com/alcb1310/bca/internals/utils"
)

func (s *BCAService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := utils.ReadCookie(r)
		if err != nil {
			slog.Error("AuthMiddleware: Unable to read cookie", "error", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		user, err := utils.ValidateJWTToken(w, r, cookie.Value)
		if err != nil {
			slog.Error("AuthMiddleware: Unable to validate token", "error", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.Clone(ctx)
		http.SetCookie(w, utils.GenerateCookie(cookie.Value))

		next.ServeHTTP(w, r)
	})
}
