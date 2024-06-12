package utils

import (
	"errors"
	"log/slog"
	"net/http"
	"time"
)

func ReadCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("bca")
	if err != nil {
		slog.Error("AuthMiddleware", "error", err)
		return nil, err
	}

	if cookie.Value == "" {
		return nil, errors.New("No cookie value set")
	}

	return cookie, nil
}

func GenerateCookie(tokenString string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "bca",
		Value:    tokenString,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   86400,
		Raw:      "",
		HttpOnly: true,
		Secure:   true,
	}

	return cookie
}
