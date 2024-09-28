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

func TestInvoice(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testunitcost"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(6*time.Second),
		),
		postgres.WithInitScripts(filepath.Join("..", "..", "internal", "database", "tables.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u000-company.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u001-project.sql")),
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

	t.Run("it should be able to list all the invoices", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/bca/partials/invoices", nil)
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "No hay facturas")
	})

	t.Run("it should be able to create an invoice", func(t *testing.T) {
		form := url.Values{
			"supplier":      []string{"b2fa8dc7-49e1-41d1-b1f8-0646eb1346b4"},
			"project":       []string{"2118e27d-1ae5-4554-b0ba-2503917a31aa"},
			"invoiceNumber": []string{"123"},
			"invoiceDate":   []string{"2021-01-01"},
		}

		req, err := http.NewRequest(http.MethodPost, "/bca/transacciones/facturas/crear", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		res := httptest.NewRecorder()
		s.Router.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "Test Supplier")
		assert.Contains(t, res.Body.String(), "Project 1")
		assert.Contains(t, res.Body.String(), "123")
		assert.Contains(t, res.Body.String(), "2021-01-01")
	})

	t.Run("single invoice", func(t *testing.T) {
		companyId := getCompanyId(t, s, cookies)
		invoices, err := s.DB.GetInvoices(companyId)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(invoices))
		testUrl := fmt.Sprintf("/bca/partials/invoices/%s", invoices[0].Id.String())

		t.Run("it should be able to show an invoice", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, testUrl, nil)
			assert.NoError(t, err)
			req.AddCookie(cookies[0])
			res := httptest.NewRecorder()
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, http.StatusOK, res.Code)
			assert.Contains(t, res.Body.String(), "Test Supplier")
			assert.Contains(t, res.Body.String(), "Project 1")
			assert.Contains(t, res.Body.String(), "123")
			assert.Contains(t, res.Body.String(), "2021-01-01")
		})

		t.Run("it should be able to update an invoice", func(t *testing.T) {
			form := url.Values{
				"supplier":      []string{"b2fa8dc7-49e1-41d1-b1f8-0646eb1346b4"},
				"project":       []string{"2118e27d-1ae5-4554-b0ba-2503917a31aa"},
				"invoiceNumber": []string{"124"},
				"invoiceDate":   []string{"2022-01-01"},
			}

			req, err := http.NewRequest(http.MethodPut, testUrl, strings.NewReader(form.Encode()))
			assert.NoError(t, err)
			req.AddCookie(cookies[0])
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			res := httptest.NewRecorder()
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, http.StatusOK, res.Code)
			assert.Contains(t, res.Body.String(), "Test Supplier")
			assert.Contains(t, res.Body.String(), "Project 1")
			assert.Contains(t, res.Body.String(), "124")
			assert.Contains(t, res.Body.String(), "2022-01-01")
		})

		t.Run("it should be able to delete an invoice", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, testUrl, nil)
			assert.NoError(t, err)
			req.AddCookie(cookies[0])
			res := httptest.NewRecorder()
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, http.StatusOK, res.Code)
			assert.Contains(t, res.Body.String(), "No hay facturas")
		})
	})
}
