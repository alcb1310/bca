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

	"github.com/alcb1310/bca/internal/types"
)

func TestInvoiceDetails(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:14.1-alpine",
		postgres.WithDatabase("testinvoicedetails"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
		postgres.WithInitScripts(filepath.Join("..", "..", "internal", "database", "tables.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u000-company.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u001-project.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u001-supplier.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u001-budget-item.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u002-budget.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u002-invoice.sql")),
	)
	assert.NoError(t, err)
	if err != nil {
		slog.Error("TestInvoiceDetails, failed to run pgContainer", "err", err)
		panic(err)
	}

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

	invoiceId := uuid.MustParse("c3be2956-1c3c-46f7-af14-d28420116f14")
	testUrl := fmt.Sprintf("/bca/partials/invoices/%s/details", invoiceId)

	t.Run("should display invoice details", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testUrl, nil)
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "</thead> <tbody></tbody></table>")
	})

	t.Run("should be able to create an invoice detail", func(t *testing.T) {
		invoiceDetails := []types.InvoiceDetails{
			{
				BudgetItemId:   uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
				BudgetItemCode: "500.1",
				BudgetItemName: "Obra Gruesa",
				Quantity:       1,
				Cost:           25,
				Total:          25,
				InvoiceTotal:   25,
				CompanyId:      uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
		}

		budgetResponse := []types.Budget{
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("420f8bb3-bc8e-4564-be99-75cd7c1a6ff8"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("9abc2426-a92b-46ef-b074-ddbc8ee2df1a"),
				InitialQuantity:   sql.NullFloat64{Float64: 2537.5, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 1.8, Valid: true},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 2537.5, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 1.8, Valid: true},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        25.0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    75.0,
				UpdatedBudget:     100.0,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 1, Valid: true},
				SpentTotal:        25.0,
				RemainingQuantity: sql.NullFloat64{Float64: 3, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 25, Valid: true},
				RemainingTotal:    75.0,
				UpdatedBudget:     100.0,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828"),
				BudgetItemId:      uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    100,
				UpdatedBudget:     100,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828"),
				BudgetItemId:      uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 4, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 25, Valid: true},
				RemainingTotal:    100,
				UpdatedBudget:     100,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
		}

		form := url.Values{
			"item":     {"b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"},
			"quantity": {"1"},
			"cost":     {"25"},
		}
		req, err := http.NewRequest(http.MethodPost, testUrl, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)

		savedDetails, err := s.DB.GetAllDetails(invoiceId, uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"))
		assert.NoError(t, err)
		assert.Equal(t, 1, len(savedDetails))

		assert.Equal(t, savedDetails[0].BudgetItemId, invoiceDetails[0].BudgetItemId)
		assert.Equal(t, savedDetails[0].BudgetItemCode, invoiceDetails[0].BudgetItemCode)
		assert.Equal(t, savedDetails[0].BudgetItemName, invoiceDetails[0].BudgetItemName)
		assert.Equal(t, savedDetails[0].Quantity, invoiceDetails[0].Quantity)
		assert.Equal(t, savedDetails[0].Cost, invoiceDetails[0].Cost)
		assert.Equal(t, savedDetails[0].Total, invoiceDetails[0].Total)
		assert.Equal(t, savedDetails[0].InvoiceTotal, invoiceDetails[0].InvoiceTotal)
		assert.Equal(t, savedDetails[0].CompanyId, invoiceDetails[0].CompanyId)

		budgets, err := s.DB.GetBudgets(
			uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			uuid.UUID{},
			"",
		)
		assert.NoError(t, err)
		assert.Equal(t, 6, len(budgets))

		for i, b := range budgets {
			assert.Equal(t, b.Project.ID, budgetResponse[i].ProjectId)
			assert.Equal(t, b.BudgetItem.ID, budgetResponse[i].BudgetItemId)
			assert.Equal(t, b.InitialQuantity, budgetResponse[i].InitialQuantity)
			assert.Equal(t, b.InitialCost, budgetResponse[i].InitialCost)
			assert.Equal(t, b.InitialTotal, budgetResponse[i].InitialTotal)
			assert.Equal(t, b.SpentQuantity, budgetResponse[i].SpentQuantity)
			assert.Equal(t, b.SpentTotal, budgetResponse[i].SpentTotal)
			assert.Equal(t, b.RemainingQuantity, budgetResponse[i].RemainingQuantity)
			assert.Equal(t, b.RemainingCost, budgetResponse[i].RemainingCost)
			assert.Equal(t, b.RemainingTotal, budgetResponse[i].RemainingTotal)
			assert.Equal(t, b.UpdatedBudget, budgetResponse[i].UpdatedBudget)
			assert.Equal(t, b.CompanyId, budgetResponse[i].CompanyId)
		}
	})

	t.Run("should create a conflict when creating the same detail twice", func(t *testing.T) {
		form := url.Values{
			"item":     {"b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"},
			"quantity": {"1"},
			"cost":     {"25"},
		}
		req, err := http.NewRequest(http.MethodPost, testUrl, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Contains(t, resp.Body.String(), "Ya existe una partida con ese nombre en la factura")
	})

	t.Run("should delete an invoice detail", func(t *testing.T) {
		budgetResponse := []types.Budget{
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("420f8bb3-bc8e-4564-be99-75cd7c1a6ff8"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("9abc2426-a92b-46ef-b074-ddbc8ee2df1a"),
				InitialQuantity:   sql.NullFloat64{Float64: 2537.5, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 1.8, Valid: true},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 2537.5, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 1.8, Valid: true},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0.0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    100.0,
				UpdatedBudget:     100.0,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0.0,
				RemainingQuantity: sql.NullFloat64{Float64: 4, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 25, Valid: true},
				RemainingTotal:    100.0,
				UpdatedBudget:     100.0,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828"),
				BudgetItemId:      uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    100,
				UpdatedBudget:     100,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828"),
				BudgetItemId:      uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 4, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 25, Valid: true},
				RemainingTotal:    100,
				UpdatedBudget:     100,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
		}

		form := url.Values{}
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", testUrl, "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"), strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "<tbody></tbody>")

		budgets, err := s.DB.GetBudgets(
			uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			uuid.UUID{},
			"",
		)
		assert.NoError(t, err)
		assert.Equal(t, 6, len(budgets))

		for i, b := range budgets {
			assert.Equal(t, b.Project.ID, budgetResponse[i].ProjectId)
			assert.Equal(t, b.BudgetItem.ID, budgetResponse[i].BudgetItemId)
			assert.Equal(t, b.InitialQuantity, budgetResponse[i].InitialQuantity)
			assert.Equal(t, b.InitialCost, budgetResponse[i].InitialCost)
			assert.Equal(t, b.InitialTotal, budgetResponse[i].InitialTotal)
			assert.Equal(t, b.SpentQuantity, budgetResponse[i].SpentQuantity)
			assert.Equal(t, b.SpentTotal, budgetResponse[i].SpentTotal)
			assert.Equal(t, b.RemainingQuantity, budgetResponse[i].RemainingQuantity)
			assert.Equal(t, b.RemainingCost, budgetResponse[i].RemainingCost)
			assert.Equal(t, b.RemainingTotal, budgetResponse[i].RemainingTotal)
			assert.Equal(t, b.UpdatedBudget, budgetResponse[i].UpdatedBudget)
			assert.Equal(t, b.CompanyId, budgetResponse[i].CompanyId)
		}
	})

	t.Run("should update correctly the budget", func(t *testing.T) {
		invoiceDetails := []types.InvoiceDetails{
			{
				BudgetItemId:   uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
				BudgetItemCode: "500.1",
				BudgetItemName: "Obra Gruesa",
				Quantity:       1,
				Cost:           30,
				Total:          30,
				InvoiceTotal:   30,
				CompanyId:      uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
		}

		budgetResponse := []types.Budget{
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("420f8bb3-bc8e-4564-be99-75cd7c1a6ff8"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("9abc2426-a92b-46ef-b074-ddbc8ee2df1a"),
				InitialQuantity:   sql.NullFloat64{Float64: 2537.5, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 1.8, Valid: true},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 2537.5, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 1.8, Valid: true},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        30.0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    90.0,
				UpdatedBudget:     120.0,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 1, Valid: true},
				SpentTotal:        30.0,
				RemainingQuantity: sql.NullFloat64{Float64: 3, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 30, Valid: true},
				RemainingTotal:    90.0,
				UpdatedBudget:     120.0,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828"),
				BudgetItemId:      uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    100,
				UpdatedBudget:     100,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828"),
				BudgetItemId:      uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 4, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 25, Valid: true},
				RemainingTotal:    100,
				UpdatedBudget:     100,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
		}

		form := url.Values{
			"item":     {"b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"},
			"quantity": {"1"},
			"cost":     {"30"},
		}
		req, err := http.NewRequest(http.MethodPost, testUrl, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)

		savedDetails, err := s.DB.GetAllDetails(invoiceId, uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"))
		assert.NoError(t, err)
		assert.Equal(t, 1, len(savedDetails))

		assert.Equal(t, savedDetails[0].BudgetItemId, invoiceDetails[0].BudgetItemId)
		assert.Equal(t, savedDetails[0].BudgetItemCode, invoiceDetails[0].BudgetItemCode)
		assert.Equal(t, savedDetails[0].BudgetItemName, invoiceDetails[0].BudgetItemName)
		assert.Equal(t, savedDetails[0].Quantity, invoiceDetails[0].Quantity)
		assert.Equal(t, savedDetails[0].Cost, invoiceDetails[0].Cost)
		assert.Equal(t, savedDetails[0].Total, invoiceDetails[0].Total)
		assert.Equal(t, savedDetails[0].InvoiceTotal, invoiceDetails[0].InvoiceTotal)
		assert.Equal(t, savedDetails[0].CompanyId, invoiceDetails[0].CompanyId)

		budgets, err := s.DB.GetBudgets(
			uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			uuid.UUID{},
			"",
		)
		assert.NoError(t, err)
		assert.Equal(t, 6, len(budgets))

		for i, b := range budgets {
			assert.Equal(t, b.Project.ID, budgetResponse[i].ProjectId)
			assert.Equal(t, b.BudgetItem.ID, budgetResponse[i].BudgetItemId)
			assert.Equal(t, b.InitialQuantity, budgetResponse[i].InitialQuantity)
			assert.Equal(t, b.InitialCost, budgetResponse[i].InitialCost)
			assert.Equal(t, b.InitialTotal, budgetResponse[i].InitialTotal)
			assert.Equal(t, b.SpentQuantity, budgetResponse[i].SpentQuantity)
			assert.Equal(t, b.SpentTotal, budgetResponse[i].SpentTotal)
			assert.Equal(t, b.RemainingQuantity, budgetResponse[i].RemainingQuantity)
			assert.Equal(t, b.RemainingCost, budgetResponse[i].RemainingCost)
			assert.Equal(t, b.RemainingTotal, budgetResponse[i].RemainingTotal)
			assert.Equal(t, b.UpdatedBudget, budgetResponse[i].UpdatedBudget)
			assert.Equal(t, b.CompanyId, budgetResponse[i].CompanyId)
		}
	})

	t.Run("should reduce the budget if less cost", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", testUrl, "b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"), nil)
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "<tbody></tbody>")

		invoiceDetails := []types.InvoiceDetails{
			{
				BudgetItemId:   uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
				BudgetItemCode: "500.1",
				BudgetItemName: "Obra Gruesa",
				Quantity:       1,
				Cost:           20,
				Total:          20,
				InvoiceTotal:   20,
				CompanyId:      uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
		}

		budgetResponse := []types.Budget{
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("420f8bb3-bc8e-4564-be99-75cd7c1a6ff8"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("9abc2426-a92b-46ef-b074-ddbc8ee2df1a"),
				InitialQuantity:   sql.NullFloat64{Float64: 2537.5, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 1.8, Valid: true},
				InitialTotal:      4567.5,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 2537.5, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 1.8, Valid: true},
				RemainingTotal:    4567.5,
				UpdatedBudget:     4567.5,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        20.0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    60.0,
				UpdatedBudget:     80.0,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"),
				BudgetItemId:      uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 1, Valid: true},
				SpentTotal:        20.0,
				RemainingQuantity: sql.NullFloat64{Float64: 3, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 20, Valid: true},
				RemainingTotal:    60.0,
				UpdatedBudget:     80.0,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828"),
				BudgetItemId:      uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc"),
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      100,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    100,
				UpdatedBudget:     100,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
			{
				ProjectId:         uuid.MustParse("1c6020db-39a0-451d-89ee-fdd20d519828"),
				BudgetItemId:      uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"),
				InitialQuantity:   sql.NullFloat64{Float64: 4, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 25, Valid: true},
				InitialTotal:      100,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: true},
				SpentTotal:        0,
				RemainingQuantity: sql.NullFloat64{Float64: 4, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 25, Valid: true},
				RemainingTotal:    100,
				UpdatedBudget:     100,
				CompanyId:         uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			},
		}

		form := url.Values{
			"item":     {"b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"},
			"quantity": {"1"},
			"cost":     {"20"},
		}
		req, err = http.NewRequest(http.MethodPost, testUrl, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookies[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp = httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.NoError(t, err)

		savedDetails, err := s.DB.GetAllDetails(invoiceId, uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"))
		assert.NoError(t, err)
		assert.Equal(t, 1, len(savedDetails))

		assert.Equal(t, savedDetails[0].BudgetItemId, invoiceDetails[0].BudgetItemId)
		assert.Equal(t, savedDetails[0].BudgetItemCode, invoiceDetails[0].BudgetItemCode)
		assert.Equal(t, savedDetails[0].BudgetItemName, invoiceDetails[0].BudgetItemName)
		assert.Equal(t, savedDetails[0].Quantity, invoiceDetails[0].Quantity)
		assert.Equal(t, savedDetails[0].Cost, invoiceDetails[0].Cost)
		assert.Equal(t, savedDetails[0].Total, invoiceDetails[0].Total)
		assert.Equal(t, savedDetails[0].InvoiceTotal, invoiceDetails[0].InvoiceTotal)
		assert.Equal(t, savedDetails[0].CompanyId, invoiceDetails[0].CompanyId)

		budgets, err := s.DB.GetBudgets(
			uuid.MustParse("3308a6e7-4060-4d7c-8490-f1ccddd9c411"),
			uuid.UUID{},
			"",
		)
		assert.NoError(t, err)
		assert.Equal(t, 6, len(budgets))

		for i, b := range budgets {
			assert.Equal(t, b.Project.ID, budgetResponse[i].ProjectId)
			assert.Equal(t, b.BudgetItem.ID, budgetResponse[i].BudgetItemId)
			assert.Equal(t, b.InitialQuantity, budgetResponse[i].InitialQuantity)
			assert.Equal(t, b.InitialCost, budgetResponse[i].InitialCost)
			assert.Equal(t, b.InitialTotal, budgetResponse[i].InitialTotal)
			assert.Equal(t, b.SpentQuantity, budgetResponse[i].SpentQuantity)
			assert.Equal(t, b.SpentTotal, budgetResponse[i].SpentTotal)
			assert.Equal(t, b.RemainingQuantity, budgetResponse[i].RemainingQuantity)
			assert.Equal(t, b.RemainingCost, budgetResponse[i].RemainingCost)
			assert.Equal(t, b.RemainingTotal, budgetResponse[i].RemainingTotal)
			assert.Equal(t, b.UpdatedBudget, budgetResponse[i].UpdatedBudget)
			assert.Equal(t, b.CompanyId, budgetResponse[i].CompanyId)
		}
	})
}
