package server

import (
	"bca-go-final/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, r)
	})
}

func (s *Server) authVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "/api/v1") {
			next.ServeHTTP(w, r)
			return
		}
		cookie, _ := r.Cookie("bca")

		token := strings.Split(cookie.String(), "=")
		if len(token) != 2 {
			// http.Redirect(w, r, "/login", http.StatusSeeOther)
			header := r.Header.Get("x-access-token")
			if header == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token = append(token, header)
		}

		secretKey := os.Getenv("SECRET")
		maker, err := utils.NewJWTMaker(secretKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenData, err := maker.VerifyToken(token[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		marshalStr, _ := json.Marshal(tokenData)
		ctx := r.Context()
		ctx = context.WithValue(r.Context(), "token", marshalStr)
		r = r.Clone(ctx)

		if !s.DB.IsLoggedIn(token[1], tokenData.ID) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
