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

func TestSuppliers(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:14.1-alpine",
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
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err.Error())
		}
	})

	s, cookies, err := createServer(t, ctx, pgContainer)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, 1, len(cookies))
	assert.Equal(t, "jwt", cookies[0].Name)
	assert.NotEmpty(t, cookies[0].Value)

	t.Run("should have no suppliers", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/bca/partials/suppliers", nil)
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "No hay Proveedores")
	})

	t.Run("should be able to create a supplier", func(t *testing.T) {
		form := url.Values{
			"supplier_id": []string{"123456789"},
			"name":        []string{"Test Supplier"},
		}
		req, err := http.NewRequest("POST", "/bca/partials/suppliers", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "<td>123456789</td>")
		assert.Contains(t, res.Body.String(), "<td>Test Supplier</td>")
	})

	t.Run("should conflict on supplier_id", func(t *testing.T) {
		form := url.Values{
			"supplier_id": []string{"123456789"},
			"name":        []string{"Another Supplier"},
		}
		req, err := http.NewRequest("POST", "/bca/partials/suppliers", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusConflict, res.Code)
		assert.Contains(t, res.Body.String(), "Proveedor con ruc 123456789 y/o nombre Another Supplier ya existe")
	})

	t.Run("should conflict on name", func(t *testing.T) {
		form := url.Values{
			"supplier_id": []string{"987654321"},
			"name":        []string{"Test Supplier"},
		}
		req, err := http.NewRequest("POST", "/bca/partials/suppliers", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusConflict, res.Code)
		assert.Contains(t, res.Body.String(), "Proveedor con ruc 987654321 y/o nombre Test Supplier ya existe")
	})
}

func TestSingleSupplier(t *testing.T) {
	// TODO: implement
	supplierId := "b2fa8dc7-49e1-41d1-b1f8-0646eb1346b4"
	testUrl := fmt.Sprintf("/bca/partials/suppliers/%s", supplierId)
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:14.1-alpine",
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
		postgres.WithInitScripts(filepath.Join("scripts", "u001-supplier.sql")),
	)
	assert.NoError(t, err)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err.Error())
		}
	})

	s, cookies, err := createServer(t, ctx, pgContainer)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, 1, len(cookies))
	assert.Equal(t, "jwt", cookies[0].Name)
	assert.NotEmpty(t, cookies[0].Value)

	t.Run("it should return the selected supplier", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testUrl, nil)
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), supplierId)
		assert.Contains(t, res.Body.String(), "123456789")
		assert.Contains(t, res.Body.String(), "Test Supplier")
	})

  t.Run("it should return 404 if no supplier found", func(t *testing.T) {
    req, err := http.NewRequest(http.MethodGet, "/bca/partials/suppliers/00000000-0000-0000-0000-000000000000", nil)
    assert.NoError(t, err)
    req.AddCookie(cookies[0])
    res := httptest.NewRecorder()
    s.Router.ServeHTTP(res, req)
    assert.Equal(t, http.StatusNotFound, res.Code)
    assert.Equal(t, res.Body.String(), "Proveedor no encontrado")
  })

  t.Run("it should update a supplier", func(t *testing.T) {
    form := url.Values{
      "supplier_id": []string{"123987456"},
      "name":        []string{"Another Supplier"},
    }
    testUrl := fmt.Sprintf("/bca/partials/suppliers/edit/%s", supplierId)
    req, err := http.NewRequest(http.MethodPut, testUrl, strings.NewReader(form.Encode()))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    assert.NoError(t, err)
    req.AddCookie(cookies[0])
    res := httptest.NewRecorder()
    s.Router.ServeHTTP(res, req)
    assert.Equal(t, http.StatusOK, res.Code)
    assert.Contains(t, res.Body.String(), supplierId)
    assert.Contains(t, res.Body.String(), "<td>123987456</td>")
    assert.Contains(t, res.Body.String(), "<td>Another Supplier</td>")
  })
}
