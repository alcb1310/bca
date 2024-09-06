package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
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

func TestBudgetItem(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:14.1-alpine",
		postgres.WithDatabase("testbudgetitem"),
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

	t.Run("it should display no budget items", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/bca/partials/budget-item", nil)
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "No hay Partidas")
	})

	t.Run("it should create a budget item", func(t *testing.T) {
		form := url.Values{
			"code":       []string{"500"},
			"name":       []string{"Costo Directo"},
			"accumulate": []string{"accumulate"},
		}
		req, err := http.NewRequest(http.MethodPost, "/bca/partials/budget-item", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", fmt.Sprintf("%s=%s", cookies[0].Name, cookies[0].Value))
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		assert.Contains(t, resp.Body.String(), "<td>500</td>")
		assert.Contains(t, resp.Body.String(), "<td>Costo Directo</td>")

		form = url.Values{
			"code":       []string{"600"},
			"name":       []string{"NonAccumulate"},
			"accumulate": []string{"false"},
		}
		req, err = http.NewRequest(http.MethodPost, "/bca/partials/budget-item", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Cookie", fmt.Sprintf("%s=%s", cookies[0].Name, cookies[0].Value))
		resp = httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		assert.Contains(t, resp.Body.String(), "<td>600</td>")
		assert.Contains(t, resp.Body.String(), "<td>NonAccumulate</td>")

		assert.Contains(t, resp.Body.String(), "<td>500</td>")
		assert.Contains(t, resp.Body.String(), "<td>Costo Directo</td>")
	})

	t.Run("it should conflict with same budget item code", func(t *testing.T) {
		form := url.Values{
			"code":       []string{"500"},
			"name":       []string{"Test"},
			"accumulate": []string{"true"},
		}
		req, err := http.NewRequest(http.MethodPost, "/bca/partials/budget-item", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookies[0])
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)

		assert.Equal(t, resp.Body.String(), "Ya existe una partida con el mismo código: 500 y/o el mismo nombre: Test")
	})

	t.Run("it should conflic with same budget item name", func(t *testing.T) {
		form := url.Values{
			"code":       []string{"501"},
			"name":       []string{"Costo Directo"},
			"accumulate": []string{"true"},
		}
		req, err := http.NewRequest(http.MethodPost, "/bca/partials/budget-item", strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(cookies[0])
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Code)

		assert.Equal(t, resp.Body.String(), "Ya existe una partida con el mismo código: 501 y/o el mismo nombre: Costo Directo")
	})

	t.Run("single budget item", func(t *testing.T) {
		var bID uuid.UUID
		companyId := getCompanyId(t, s, cookies)
		budgetItems := s.DB.GetBudgetItemsByAccumulate(companyId, true)
		assert.Equal(t, 1, len(budgetItems))
		bID = budgetItems[0].ID
		testUrl := fmt.Sprintf("/bca/partials/budget-item/%s", bID.String())

		t.Run("it should get all budget-items by accumulate", func(t *testing.T) {
			var parentUUID *uuid.UUID = nil
			budgetItems := s.DB.GetBudgetItemsByAccumulate(companyId, true)
			slog.Info("single budget item", "budgetItems", budgetItems)
			assert.Equal(t, 1, len(budgetItems))
			assert.Equal(t, budgetItems[0].Name, "Costo Directo")
			assert.Equal(t, budgetItems[0].Code, "500")
			assert.Equal(t, budgetItems[0].Accumulate, sql.NullBool{Bool: true, Valid: true})
			assert.Equal(t, budgetItems[0].ParentId, parentUUID)
			assert.Equal(t, budgetItems[0].CompanyId, companyId)

			budgetItems = s.DB.GetBudgetItemsByAccumulate(companyId, false)
			assert.Equal(t, 1, len(budgetItems))
			assert.Equal(t, budgetItems[0].Name, "NonAccumulate")
			assert.Equal(t, budgetItems[0].Code, "600")
			assert.Equal(t, budgetItems[0].Accumulate, sql.NullBool{Bool: false, Valid: true})
			assert.Equal(t, budgetItems[0].ParentId, parentUUID)
			assert.Equal(t, budgetItems[0].CompanyId, companyId)
		})

		t.Run("it should get one budget item", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, testUrl, nil)
			req.AddCookie(cookies[0])
			assert.NoError(t, err)
			resp := httptest.NewRecorder()
			s.Router.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Contains(t, resp.Body.String(), "Costo Directo")
			assert.Contains(t, resp.Body.String(), "500")
		})

		t.Run("it should return 404 if no buget item found", func(t *testing.T) {
			testUrl := fmt.Sprintf("/bca/partials/budget-item/%s", "00000000-0000-0000-0000-000000000000")
			req, err := http.NewRequest(http.MethodGet, testUrl, nil)
			req.AddCookie(cookies[0])
			assert.NoError(t, err)
			resp := httptest.NewRecorder()
			s.Router.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusNotFound, resp.Code)
			assert.Equal(t, resp.Body.String(), "Partida no encontrada")
		})

		t.Run("it should update a budget item", func(t *testing.T) {
			form := url.Values{
				"name": []string{"Another Budget Item"},
			}
			req, err := http.NewRequest(http.MethodPut, testUrl, strings.NewReader(form.Encode()))
			assert.NoError(t, err)
			req.AddCookie(cookies[0])
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp := httptest.NewRecorder()
			s.Router.ServeHTTP(resp, req)
			assert.NotNil(t, s.Router)
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Contains(t, resp.Body.String(), "Another Budget Item")
		})

		t.Run("it should not let update the parent id", func(t *testing.T) {
			form := url.Values{
				"name":   []string{"Another Budget Item"},
				"parent": []string{"420f8bb3-bc8e-4564-be99-75cd7c1a6ff8"},
			}
			testUrl := fmt.Sprintf("/bca/partials/budget-item/%s", bID.String())
			req, err := http.NewRequest(http.MethodPut, testUrl, strings.NewReader(form.Encode()))
			assert.NoError(t, err)
			req.AddCookie(cookies[0])
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			resp := httptest.NewRecorder()
			s.Router.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusBadRequest, resp.Code)
			assert.Contains(t, resp.Body.String(), "No se puede cambiar la partida padre")
		})
	})
}
