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
)

func TestCategories(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testcategories"),
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

	t.Run("it should get all categories", func(t *testing.T) {
		testUrl := "/bca/partials/categories"
		req, err := http.NewRequest(http.MethodGet, testUrl, nil)
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "No existen categorías")
	})

	t.Run("it should be able to create a category", func(t *testing.T) {
		testUrl := "/bca/partials/categories"
		form := url.Values{
			"name": {"test"},
		}
		req, err := http.NewRequest(http.MethodPost, testUrl, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "<td>test</td>")
	})

	t.Run("it should conflict when creating a category with the same name", func(t *testing.T) {
		testUrl := "/bca/partials/categories"
		form := url.Values{
			"name": {"test"},
		}
		req, err := http.NewRequest(http.MethodPost, testUrl, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Contains(t, resp.Body.String(), "La categoría test ya existe")
	})

	t.Run("single category", func(t *testing.T) {
		companyId := getCompanyId(t, s, cookie)
		categories, err := s.DB.GetAllCategories(companyId)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(categories))
		categoryId := categories[0].Id
		testUrl := fmt.Sprintf("/bca/partials/categories/%s", categoryId)

		t.Run("it should get a single category", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, testUrl, nil)
			assert.NoError(t, err)
			req.AddCookie(cookie[0])
			resp := httptest.NewRecorder()
			s.Router.ServeHTTP(resp, req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Contains(t, resp.Body.String(), "Editar Categoría")
			assert.Contains(t, resp.Body.String(), "test")
		})

		t.Run("it should return 404 if not found", func(t *testing.T) {
			testUrl := "/bca/partials/categories/00000000-0000-0000-0000-000000000000"
			req, err := http.NewRequest(http.MethodGet, testUrl, nil)
			assert.NoError(t, err)
			req.AddCookie(cookie[0])
			resp := httptest.NewRecorder()
			s.Router.ServeHTTP(resp, req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.Code)
			assert.Equal(t, resp.Body.String(), "Categoría no encontrada")
		})

		t.Run("it should update a category", func(t *testing.T) {
			form := url.Values{
				"name": {"CATEGORÍA 2"},
			}
			req, err := http.NewRequest(http.MethodPut, testUrl, strings.NewReader(form.Encode()))
			assert.NoError(t, err)
			req.AddCookie(cookie[0])
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			resp := httptest.NewRecorder()
			s.Router.ServeHTTP(resp, req)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Contains(t, resp.Body.String(), "CATEGORÍA 2")
		})
	})
}
