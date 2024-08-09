package integration

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"bca-go-final/internal/database"
	"bca-go-final/internal/server"
)

func TestLogin(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.4-alpine"),
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
		postgres.WithInitScripts(filepath.Join("..", "..", "internal", "database", "tables.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u000-company.sql")),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err.Error())
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)
	assert.NotEmpty(t, connStr)

	db := database.New(connStr)
	assert.NotNil(t, db)

	h := db.Health()
	assert.Equal(t, "It's healthy", h["message"])
	slog.Info("Database is healthy")

	s := server.NewServer(db, "supersecretpassword")
	assert.NotNil(t, s)

	t.Run("should display the login page", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/login", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		assert.Contains(t, res.Body.String(), "Login")
		assert.Contains(t, res.Body.String(), "Email")
		assert.Contains(t, res.Body.String(), "Password")
	})

	t.Run("should not login on invalid credentials", func(t *testing.T) {
		t.Run("wrong password", func(t *testing.T) {
			form := url.Values{
				"email":    {"test@test.com"},
				"password": {"wrongpassword"},
			}
			req, err := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			res := httptest.NewRecorder()
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, http.StatusBadRequest, res.Code)
			assert.Contains(t, res.Body.String(), "credenciales inválidas")
		})

		t.Run("wrong email", func(t *testing.T) {
			form := url.Values{
				"email":    {"wrong@test.com"},
				"password": {"password123"},
			}
			req, err := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			res := httptest.NewRecorder()
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, http.StatusBadRequest, res.Code)
			assert.Contains(t, res.Body.String(), "credenciales inválidas")
		})
	})

	t.Run("should login the user", func(t *testing.T) {
		form := url.Values{
			"email":    {"test@test.com"},
			"password": {"password123"},
		}
		req, err := http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusSeeOther, res.Code)

		loc, err := res.Result().Location()
		assert.NoError(t, err)
		assert.Equal(t, "/bca", loc.Path)

		cookies := res.Result().Cookies()
		assert.Equal(t, 1, len(cookies))
		assert.Equal(t, "jwt", cookies[0].Name)
		assert.NotEmpty(t, cookies[0].Value)

		req, err = http.NewRequest("GET", loc.Path, nil)
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		res = httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		assert.Contains(t, res.Body.String(), "Bienvenido")
		assert.Contains(t, res.Body.String(), "Test User")
	})
}
