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

func TestProjects(t *testing.T) {
	ctx := context.Background()
	pgContaineer, err := postgres.Run(ctx,
		"postgres:16.4-alpine",
		postgres.WithDatabase("testproject"),
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
		if err := pgContaineer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err.Error())
		}
	})

	s, cookies, err := createServer(t, ctx, pgContaineer)
	assert.NoError(t, err)

	t.Run("should have no projects", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/bca/partials/projects", nil)
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "No existen proyectos")
	})

	t.Run("should be able to create a project", func(t *testing.T) {
		form := url.Values{
			"name":       []string{"Test Project"},
			"gross_area": []string{"1"},
			"net_area":   []string{"1"},
		}
		req, err := http.NewRequest(http.MethodPost, "/bca/partials/projects", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookies[0])
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "<td>Test Project</td>")

		form = url.Values{
			"name":       []string{"Test Project 2"},
			"gross_area": []string{"1"},
			"net_area":   []string{"1"},
		}
		req, err = http.NewRequest(http.MethodPost, "/bca/partials/projects", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookies[0])
		res = httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "<td>Test Project</td>")
		assert.Contains(t, res.Body.String(), "<td>Test Project 2</td>")
	})

	t.Run("should conflict when creating an existing project", func(t *testing.T) {
		form := url.Values{
			"name":       []string{"Test Project"},
			"gross_area": []string{"1"},
			"net_area":   []string{"1"},
		}
		req, err := http.NewRequest(http.MethodPost, "/bca/partials/projects", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookies[0])
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
		assert.Contains(t, res.Body.String(), "El nombre Test Project ya existe")
	})

	t.Run("single project", func(t *testing.T) {
    companyId := getCompanyId(t, s, cookies)

		projects, err := s.DB.GetAllProjects(companyId)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(projects))
		projectId := projects[0].ID
		testUrl := fmt.Sprintf("/bca/partials/projects/%s", projectId)

		t.Run("should be able to get a project by id", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, testUrl, nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.AddCookie(cookies[0])
			res := httptest.NewRecorder()
			s.Router.ServeHTTP(res, req)

			assert.Equal(t, http.StatusOK, res.Code)
			assert.Contains(t, res.Body.String(), "Test Project")
		})

		t.Run("should return not found when getting a project that doesn't exist", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/bca/partials/projects/00000000-0000-0000-0000-000000000000", nil)
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.AddCookie(cookies[0])
			res := httptest.NewRecorder()
			s.Router.ServeHTTP(res, req)

			assert.Equal(t, http.StatusNotFound, res.Code)
			assert.Contains(t, res.Body.String(), "Proyecto no encontrado")
		})

		t.Run("should update a project", func(t *testing.T) {
			testUrl := fmt.Sprintf("/bca/partials/projects/edit/%s", projectId)
			form := url.Values{
				"name":       []string{"Updated Project"},
				"gross_area": []string{"1"},
				"net_area":   []string{"1"},
			}
			req, err := http.NewRequest(http.MethodPut, testUrl, strings.NewReader(form.Encode()))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.AddCookie(cookies[0])
			res := httptest.NewRecorder()
			s.Router.ServeHTTP(res, req)

			assert.Equal(t, http.StatusOK, res.Code)
			assert.Contains(t, res.Body.String(), "<td>Updated Project</td>")
		})
	})
}
