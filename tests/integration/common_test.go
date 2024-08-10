package integration

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"bca-go-final/internal/database"
	"bca-go-final/internal/server"
)

func login(t *testing.T, s *server.Server) ([]*http.Cookie, error) {
	form := url.Values{
		"email":    {"test@test.com"},
		"password": {"password123"},
	}
	req, err := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	assert.NoError(t, err)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()
	s.Router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusSeeOther, res.Code)
	if res.Code != http.StatusSeeOther {
		return nil, errors.New("failed to login")
	}

	cookies := res.Result().Cookies()
	assert.Equal(t, 1, len(cookies))
	if len(cookies) != 1 {
		return nil, errors.New("failed to get cookies")
	}

	assert.Equal(t, "jwt", cookies[0].Name)
	assert.NotEmpty(t, cookies[0].Value)

	return cookies, nil
}

func createServer(t *testing.T, ctx context.Context, pgContainer *postgres.PostgresContainer) (*server.Server, []*http.Cookie, error) {
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, err
	}

	db := database.New(connStr)
	if db == nil {
		return nil, nil, err
	}

	s := server.NewServer(db, "supersecretpassword")
	if s == nil {
		return nil, nil, err
	}

	cookies, err := login(t, s)
	if err != nil {
		return nil, nil, err
	}

	return s, cookies, nil
}
