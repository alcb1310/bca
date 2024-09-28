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

	"github.com/alcb1310/bca/internal/types"
)

func TestUpdateBudget(t *testing.T) {
	result := map[string]types.Budget{
		"500":   {ProjectId: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"), BudgetItemId: uuid.MustParse("439082ad-f1bd-4228-91f2-8e744894ffdc"), InitialQuantity: sql.NullFloat64{Float64: 0, Valid: false}, InitialCost: sql.NullFloat64{Float64: 0, Valid: false}, InitialTotal: 100, SpentQuantity: sql.NullFloat64{Float64: 0, Valid: false}, SpentTotal: 50, RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false}, RemainingCost: sql.NullFloat64{Float64: 0, Valid: false}, RemainingTotal: 100, UpdatedBudget: 150},
		"500.1": {ProjectId: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"), BudgetItemId: uuid.MustParse("b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb"), InitialQuantity: sql.NullFloat64{Float64: 4, Valid: true}, InitialCost: sql.NullFloat64{Float64: 25, Valid: true}, InitialTotal: 100, SpentQuantity: sql.NullFloat64{Float64: 2, Valid: true}, SpentTotal: 50, RemainingQuantity: sql.NullFloat64{Float64: 4, Valid: true}, RemainingCost: sql.NullFloat64{Float64: 25, Valid: true}, RemainingTotal: 100, UpdatedBudget: 150},
		"200":   {ProjectId: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"), BudgetItemId: uuid.MustParse("420f8bb3-bc8e-4564-be99-75cd7c1a6ff8"), InitialQuantity: sql.NullFloat64{Float64: 0, Valid: false}, InitialCost: sql.NullFloat64{Float64: 0, Valid: false}, InitialTotal: 4567.5, SpentQuantity: sql.NullFloat64{Float64: 0, Valid: false}, SpentTotal: 180, RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false}, RemainingCost: sql.NullFloat64{Float64: 0, Valid: false}, RemainingTotal: 200, UpdatedBudget: 380},
		"200.1": {ProjectId: uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa"), BudgetItemId: uuid.MustParse("9abc2426-a92b-46ef-b074-ddbc8ee2df1a"), InitialQuantity: sql.NullFloat64{Float64: 2537.5, Valid: true}, InitialCost: sql.NullFloat64{Float64: 1.8, Valid: true}, InitialTotal: 4567.5, SpentQuantity: sql.NullFloat64{Float64: 10, Valid: true}, SpentTotal: 180, RemainingQuantity: sql.NullFloat64{Float64: 40, Valid: true}, RemainingCost: sql.NullFloat64{Float64: 5, Valid: true}, RemainingTotal: 200, UpdatedBudget: 380},
	}

	ctx := context.Background()
	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testupdatebudget"),
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
		postgres.WithInitScripts(filepath.Join("scripts", "u001-budget-item.sql")),
		postgres.WithInitScripts(filepath.Join("scripts", "u003-budget.sql")),
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

	companyId := getCompanyId(t, s, cookie)
	_ = result

	t.Run("should edit a budget with spent", func(t *testing.T) {
		budgetId := uuid.MustParse("9abc2426-a92b-46ef-b074-ddbc8ee2df1a")
		projectId := uuid.MustParse("2118e27d-1ae5-4554-b0ba-2503917a31aa")
		testUrl := fmt.Sprintf("/bca/partials/budgets/%s/%s", projectId.String(), budgetId.String())
		form := url.Values{
			"project":    []string{"2118e27d-1ae5-4554-b0ba-2503917a31aa"},
			"budgetItem": []string{"9abc2426-a92b-46ef-b074-ddbc8ee2df1a"},
			"quantity":   []string{"40"},
			"cost":       []string{"5"},
		}
		req, err := http.NewRequest(http.MethodPut, testUrl, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		req.AddCookie(cookie[0])
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()
		s.Router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)

		budgets, err := s.DB.GetBudgets(companyId, projectId, "")
		assert.NoError(t, err)

		// fmt.Println(budgets)
		for _, budget := range budgets {
			b := result[budget.BudgetItem.Code]

			assert.Equal(t, b.ProjectId, budget.Project.ID)
			assert.Equal(t, b.BudgetItemId, budget.BudgetItem.ID)
			assert.Equal(t, b.InitialQuantity, budget.InitialQuantity)
			assert.Equal(t, b.InitialCost, budget.InitialCost)
			assert.Equal(t, b.InitialTotal, budget.InitialTotal)
			assert.Equal(t, b.SpentQuantity, budget.SpentQuantity)
			assert.Equal(t, b.SpentTotal, budget.SpentTotal)
			assert.Equal(t, b.RemainingQuantity, budget.RemainingQuantity)
			assert.Equal(t, b.RemainingCost, budget.RemainingCost)
			assert.Equal(t, b.RemainingTotal, budget.RemainingTotal)
			assert.Equal(t, b.UpdatedBudget, budget.UpdatedBudget)
		}
	})
}
