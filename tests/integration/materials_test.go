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

func TestMaterials(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:14.1-alpine",
		postgres.WithDatabase("testmaterials"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
		postgres.WithInitScripts(filepath.Join("..", "..", "internal", "database", "tables.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u000-company.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u001-category.sql")),
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

	t.Run("should retrieve all the materials", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/bca/partials/materiales", nil)
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "No existen materiales")
	})

	t.Run("should create a new material", func(t *testing.T) {
		form := url.Values{
			"code":     []string{"mat001"},
			"name":     []string{"Material 001"},
			"unit":     []string{"Unidad"},
			"category": []string{"a0f9b5b0-7f9b-4a9e-9f9b-7f9b5b0a9e7f"},
		}
		req, err := http.NewRequest(http.MethodPost, "/bca/partials/materiales", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie[0])
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "<td>Material 001</td>")
		assert.Contains(t, resp.Body.String(), "mat001</td>")
		assert.Contains(t, resp.Body.String(), "Unidad</td>")
		assert.Contains(t, resp.Body.String(), "<td>CATEGORÍA 1</td>")
	})

	t.Run("should display a conflict on existing code", func(t *testing.T) {
		form := url.Values{
			"code":     []string{"mat001"},
			"name":     []string{"Material 002"},
			"unit":     []string{"Unidad"},
			"category": []string{"a0f9b5b0-7f9b-4a9e-9f9b-7f9b5b0a9e7f"},
		}
		req, err := http.NewRequest(http.MethodPost, "/bca/partials/materiales", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie[0])
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Contains(t, resp.Body.String(), "El material con código mat001 y/o nombre Material 002 ya existe")
	})

	t.Run("should display a conflict on existing name", func(t *testing.T) {
		form := url.Values{
			"code":     []string{"mat002"},
			"name":     []string{"Material 001"},
			"unit":     []string{"Unidad"},
			"category": []string{"a0f9b5b0-7f9b-4a9e-9f9b-7f9b5b0a9e7f"},
		}
		req, err := http.NewRequest(http.MethodPost, "/bca/partials/materiales", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie[0])
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Contains(t, resp.Body.String(), "El material con código mat002 y/o nombre Material 001 ya existe")
	})
}

func TestSingleMaterial(t *testing.T) {
	materialId := uuid.MustParse("7f9b5b0a-7f9b-4a9e-9f9b-7f9b5b0a9e7f")
	testUrl := fmt.Sprintf("/bca/partials/materiales/%s", materialId.String())

	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:14.1-alpine",
		postgres.WithDatabase("testsinglematerial"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
		postgres.WithInitScripts(filepath.Join("..", "..", "internal", "database", "tables.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u000-company.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u001-category.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u002-material.sql")),
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

	t.Run("it should be able to display a single material", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testUrl, nil)
		assert.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookie[0])
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "0001")
		assert.Contains(t, resp.Body.String(), "PRODUCTO 1")
		assert.Contains(t, resp.Body.String(), "UN")
		assert.Contains(t, resp.Body.String(), "CATEGORÍA 1")
	})

  t.Run("it should be able to update a material", func(t *testing.T) {
    form := url.Values{
      "code":     []string{"mat001"},
      "name":     []string{"Material 001 Updated"},
      "unit":     []string{"Unidad Updated"},
      "category": []string{"a0f9b5b0-7f9b-4a9e-9f9b-7f9b5b0a9e7f"},
    }
    req, err := http.NewRequest(http.MethodPut, testUrl, strings.NewReader(form.Encode()))
    assert.NoError(t, err)
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    req.AddCookie(cookie[0])
    resp := httptest.NewRecorder()
    s.Router.ServeHTTP(resp, req)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.Code)
    assert.Contains(t, resp.Body.String(), "mat001")
    assert.Contains(t, resp.Body.String(), "Material 001 Updated")
    assert.Contains(t, resp.Body.String(), "Unidad Updated")
    assert.Contains(t, resp.Body.String(), "CATEGORÍA 1")
  })
}
