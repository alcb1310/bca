package integration

import (
	"context"
	"fmt"
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

	"github.com/alcb1310/bca/internal/types"
)

func TestUsers(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testusers"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second),
		),
		postgres.WithInitScripts(filepath.Join("..", "..", "internal", "database", "tables.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u000-company.sql")),
	)
	assert.NoError(t, err)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err.Error())
		}
	})

	s, cookie, err := createServer(t, ctx, pgContainer)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, 1, len(cookie))
	assert.Equal(t, "jwt", cookie[0].Name)
	assert.NotEmpty(t, cookie[0].Value)

	t.Run("it should return the user's profile", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/bca/user/perfil", nil)
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Test User")
		assert.Contains(t, res.Body.String(), "test@test.com")
		assert.Contains(t, res.Body.String(), "Mi Perfil")
	})

	t.Run("it should display a table with the users", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/bca/partials/users", nil)
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "<td>Test User</td>")
		assert.Contains(t, res.Body.String(), "<td>test@test.com</td>")
	})

	t.Run("it should be able to create a new user", func(t *testing.T) {
		form := url.Values{
			"name":     []string{"Automated Test"},
			"email":    []string{"automated@test.com"},
			"password": []string{"automatedtestpassword"},
		}

		req, err := http.NewRequest("POST", "/bca/partials/users", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "<td>Test User</td>")
		assert.Contains(t, res.Body.String(), "<td>test@test.com</td>")
		assert.Contains(t, res.Body.String(), "<td>Automated Test</td>")
		assert.Contains(t, res.Body.String(), "<td>automated@test.com</td>")
		assert.NotContains(t, res.Body.String(), "automatedtestpassword")
	})

	t.Run("it should be able to get a user", func(t *testing.T) {
		var user types.User
		companyId := getCompanyId(t, s, cookie)
		users, err := s.DB.GetAllUsers(companyId)
		assert.NoError(t, err)

		for _, u := range users {
			if u.Email == "automated@test.com" {
				user = u
				break
			}
		}
		assert.Equal(t, "Automated Test", user.Name)
		assert.Equal(t, "automated@test.com", user.Email)
	})

	t.Run("it should be able to update a user", func(t *testing.T) {
		var user types.User
		companyId := getCompanyId(t, s, cookie)
		users, err := s.DB.GetAllUsers(companyId)
		assert.NoError(t, err)

		for _, u := range users {
			if u.Email == "automated@test.com" {
				user = u
				break
			}
		}
		testUrl := fmt.Sprintf("/bca/partials/users/%s", user.Id.String())

		form := url.Values{
			"name":  []string{"Automated Test Updated"},
			"email": []string{"updated@test.com"},
		}

		req, err := http.NewRequest(http.MethodPut, testUrl, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "<td>Test User</td>")
		assert.Contains(t, res.Body.String(), "<td>test@test.com</td>")
		assert.Contains(t, res.Body.String(), "<td>Automated Test Updated</td>")
		assert.Contains(t, res.Body.String(), "<td>updated@test.com</td>")
		assert.NotContains(t, res.Body.String(), "<td>Automated Test</td>")
		assert.NotContains(t, res.Body.String(), "<td>automated@test.com</td>")
		assert.NotContains(t, res.Body.String(), "automatedtestpassword")
	})

	t.Run("it should be able to delete a user", func(t *testing.T) {
		var user types.User
		companyId := getCompanyId(t, s, cookie)
		users, err := s.DB.GetAllUsers(companyId)
		assert.NoError(t, err)

		for _, u := range users {
			if u.Email == "updated@test.com" {
				user = u
				break
			}
		}
		testUrl := fmt.Sprintf("/bca/partials/users/%s", user.Id.String())

		req, err := http.NewRequest(http.MethodDelete, testUrl, nil)
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "<td>Test User</td>")
		assert.Contains(t, res.Body.String(), "<td>test@test.com</td>")
		assert.NotContains(t, res.Body.String(), "<td>Automated Test Updated</td>")
		assert.NotContains(t, res.Body.String(), "<td>updated@test.com</td>")
		assert.NotContains(t, res.Body.String(), "<td>Automated Test</td>")
		assert.NotContains(t, res.Body.String(), "<td>automated@test.com</td>")
		assert.NotContains(t, res.Body.String(), "automatedtestpassword")
	})
}
