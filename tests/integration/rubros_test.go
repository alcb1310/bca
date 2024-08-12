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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestRubros(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:14.1-alpine",
		postgres.WithDatabase("testrubros"),
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

	t.Run("it should get all rubros", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/bca/partials/rubros", nil)
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "No existen rubros")
	})

	t.Run("it should create a new rubro", func(t *testing.T) {
		form := url.Values{
			"code": {"001"},
			"name": {"Rubro"},
			"unit": {"Unit"},
		}

		req, err := http.NewRequest(http.MethodPost, "/bca/configuracion/rubros/crear", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "001")
		assert.Contains(t, res.Body.String(), "Rubro")
		assert.Contains(t, res.Body.String(), "Unit")
	})

	t.Run("it should conflict with same code", func(t *testing.T) {
		form := url.Values{
			"code": {"001"},
			"name": {"New Rubro"},
			"unit": {"Unit"},
		}

		req, err := http.NewRequest(http.MethodPost, "/bca/configuracion/rubros/crear", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusConflict, res.Code)
		assert.Contains(t, res.Body.String(), "El rubro con código 001 y/o nombre New Rubro ya existe")
	})

	t.Run("it should conflict with same name", func(t *testing.T) {
		form := url.Values{
			"code": {"002"},
			"name": {"Rubro"},
			"unit": {"Unit"},
		}

		req, err := http.NewRequest(http.MethodPost, "/bca/configuracion/rubros/crear", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusConflict, res.Code)
		assert.Contains(t, res.Body.String(), "El rubro con código 002 y/o nombre Rubro ya existe")
	})
}

func TestSingleRubro(t *testing.T) {
	rubroId := uuid.MustParse("0e4b9bd3-5239-4e03-a41a-2fa0ef6363e6")
	testUrl := fmt.Sprintf("/bca/configuracion/rubros/crear?id=%s", rubroId.String())

	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:14.1-alpine",
		postgres.WithDatabase("testrubros"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
		postgres.WithInitScripts(filepath.Join("..", "..", "internal", "database", "tables.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u000-company.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u001-rubros.sql")),
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

  t.Run("it should get a single rubro", func(t *testing.T) {
    req, err := http.NewRequest(http.MethodGet, testUrl, nil)
    assert.NoError(t, err)
    req.AddCookie(cookie[0])
    res := httptest.NewRecorder()
    s.Router.ServeHTTP(res, req)
    assert.Equal(t, http.StatusOK, res.Code)
    assert.Contains(t, res.Body.String(), "C001")
    assert.Contains(t, res.Body.String(), "Test Item")
    assert.Contains(t, res.Body.String(), "unit")
  })

  t.Run("it should update a rubro", func(t *testing.T) {
    form := url.Values{
      "code": {"m001"},
      "name": {"Updated Rubro"},
      "unit": {"Updated unit"},
    }

    req, err := http.NewRequest(http.MethodPost, testUrl, strings.NewReader(form.Encode()))
    assert.NoError(t, err)
    req.AddCookie(cookie[0])
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    res := httptest.NewRecorder()
    s.Router.ServeHTTP(res, req)
    assert.Equal(t, http.StatusOK, res.Code)
    assert.Contains(t, res.Body.String(), "m001")
    assert.Contains(t, res.Body.String(), "Updated Rubro")
    assert.Contains(t, res.Body.String(), "Updated unit")
  })
}
