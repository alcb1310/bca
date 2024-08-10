package integration

import (
	"context"
	"database/sql"
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
			"accumulate": []string{"true"},
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
}

func TestSingleBudgetItem(t *testing.T) {
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
		postgres.WithInitScripts(filepath.Join("scripts", "u001-budget-item.sql")),
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

	t.Run("it should get all budget-items by accumulate", func(t *testing.T) {
		companyId := "3308a6e7-4060-4d7c-8490-f1ccddd9c411"
		companyUuid := uuid.MustParse(companyId)

		var parentUUID *uuid.UUID = nil
		budgetItems := s.DB.GetBudgetItemsByAccumulate(companyUuid, true)
		assert.Equal(t, 2, len(budgetItems))
		assert.Equal(t, budgetItems[0].Name, "Costo Directo")
		assert.Equal(t, budgetItems[0].Code, "500")
		assert.Equal(t, budgetItems[0].Accumulate, sql.NullBool{Bool: true, Valid: true})
		assert.Equal(t, budgetItems[0].ID, uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc"))
		assert.Equal(t, budgetItems[0].ParentId, parentUUID)
		assert.Equal(t, budgetItems[0].CompanyId, companyUuid)

		par := uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc")
		parentUUID = &par
		budgetItems = s.DB.GetBudgetItemsByAccumulate(companyUuid, false)
		assert.Equal(t, 2, len(budgetItems))
		assert.Equal(t, budgetItems[1].Name, "Obra Gruesa")
		assert.Equal(t, budgetItems[1].Code, "500.1")
		assert.Equal(t, budgetItems[1].Accumulate, sql.NullBool{Bool: false, Valid: true})
		assert.Equal(t, budgetItems[1].ID, uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"))
		assert.Equal(t, budgetItems[1].ParentId, parentUUID)
		assert.Equal(t, budgetItems[1].CompanyId, companyUuid)
	})

	t.Run("it should get one budget item", func(t *testing.T) {
		testUrl := fmt.Sprintf("/bca/partials/budget-item/%s", "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")
		req, err := http.NewRequest(http.MethodGet, testUrl, nil)
		req.AddCookie(cookies[0])
		assert.NoError(t, err)
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Obra Gruesa")
		assert.Contains(t, resp.Body.String(), "500.1")
		assert.Contains(t, resp.Body.String(), "439082ad-f1bd-4228-91f2-8e744894ffdc")
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
		testUrl := fmt.Sprintf("/bca/partials/budget-item/%s", "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")
		req, err := http.NewRequest(http.MethodPut, testUrl, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")
		assert.Contains(t, resp.Body.String(), "Another Budget Item")
	})

	t.Run("it should not let update the parent id", func(t *testing.T) {
		form := url.Values{
			"name":      []string{"Another Budget Item"},
			"parent": []string{"420f8bb3-bc8e-4564-be99-75cd7c1a6ff8"},
		}
		testUrl := fmt.Sprintf("/bca/partials/budget-item/%s", "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb")
		req, err := http.NewRequest(http.MethodPut, testUrl, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "No se puede cambiar la partida padre")
	})
}
