package server_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/alcb1310/bca/internal/server"
	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/mocks"
)

func TestCreateBudget(t *testing.T) {
	pId := uuid.New()
	bId := uuid.New()
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret", -5)
	token := createToken(s.TokenAuth)
	testURL := "/bca/partials/budgets"

	testData := []struct {
		name         string
		form         url.Values
		status       int
		body         []string
		createBudget *mocks.Service_CreateBudget_Call
		getBudgets   *mocks.Service_GetBudgets_Call
	}{
		{
			name:         "should pass a form",
			form:         nil,
			status:       http.StatusBadRequest,
			body:         []string{},
			createBudget: nil,
			getBudgets:   nil,
		},
		{
			name:         "should pass a project",
			form:         url.Values{},
			status:       http.StatusBadRequest,
			body:         []string{"Seleccione un proyecto"},
			createBudget: nil,
			getBudgets:   nil,
		},
		{
			name: "should pass a valid project id",
			form: url.Values{
				"project": []string{"123"},
			},
			status:       http.StatusBadRequest,
			body:         []string{"Código del proyecto inválido"},
			createBudget: nil,
			getBudgets:   nil,
		},
		{
			name: "should pass a budget item",
			form: url.Values{
				"project": []string{pId.String()},
			},
			status:       http.StatusBadRequest,
			body:         []string{"Seleccione una partida"},
			createBudget: nil,
			getBudgets:   nil,
		},
		{
			name: "should pass a valid budget item id",
			form: url.Values{
				"project":    []string{pId.String()},
				"budgetItem": []string{"123"},
			},
			status:       http.StatusBadRequest,
			body:         []string{"Código de la partida inválido"},
			createBudget: nil,
			getBudgets:   nil,
		},
		{
			name: "should pass a quantity",
			form: url.Values{
				"project":    []string{pId.String()},
				"budgetItem": []string{bId.String()},
			},
			status:       http.StatusBadRequest,
			body:         []string{"La cantidad debe ser un número"},
			createBudget: nil,
			getBudgets:   nil,
		},
		{
			name: "should pass a number for quantity",
			form: url.Values{
				"project":    []string{pId.String()},
				"budgetItem": []string{bId.String()},
				"quantity":   []string{"test"},
			},
			status:       http.StatusBadRequest,
			body:         []string{"La cantidad debe ser un número"},
			createBudget: nil,
			getBudgets:   nil,
		},
		{
			name: "should pass a cost",
			form: url.Values{
				"project":    []string{pId.String()},
				"budgetItem": []string{bId.String()},
				"quantity":   []string{"10.0"},
			},
			status:       http.StatusBadRequest,
			body:         []string{"El costo debe ser un número"},
			createBudget: nil,
			getBudgets:   nil,
		},
		{
			name: "should pass a number for cost",
			form: url.Values{
				"project":    []string{pId.String()},
				"budgetItem": []string{bId.String()},
				"quantity":   []string{"10.0"},
				"cost":       []string{"test"},
			},
			status:       http.StatusBadRequest,
			body:         []string{"El costo debe ser un número"},
			createBudget: nil,
			getBudgets:   nil,
		},
		{
			name: "should create a budget",
			form: url.Values{
				"project":    []string{pId.String()},
				"budgetItem": []string{bId.String()},
				"quantity":   []string{"10.0"},
				"cost":       []string{"10.0"},
			},
			status: http.StatusOK,
			body:   []string{},
			createBudget: db.EXPECT().CreateBudget(&types.CreateBudget{
				ProjectId:    pId,
				BudgetItemId: bId,
				CompanyId:    uuid.UUID{},
				Quantity:     10.0,
				Cost:         10.0,
			}).Return(types.Budget{
				ProjectId:         pId,
				BudgetItemId:      bId,
				InitialQuantity:   sql.NullFloat64{Float64: 10.0, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 10.0, Valid: true},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0.0, Valid: true},
				SpentTotal:        0.0,
				RemainingQuantity: sql.NullFloat64{Float64: 10.0, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 10.0, Valid: true},
				RemainingTotal:    100.0,
				UpdatedBudget:     100.0,
				CompanyId:         uuid.UUID{},
			}, nil),
			getBudgets: db.EXPECT().GetBudgets(uuid.UUID{}, uuid.UUID{}, "").Return([]types.GetBudget{
				{},
			}, nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createBudget != nil {
				tt.createBudget.Times(1)
			}

			if tt.getBudgets != nil {
				tt.getBudgets.Times(1)
			}

			req, res := createRequest(
				token,
				http.MethodPost,
				testURL,
				strings.NewReader(tt.form.Encode()),
			)
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}

func TestBudgetEdit(t *testing.T) {
	projectId := uuid.New()
	budgetItemId := uuid.New()
	testURL := fmt.Sprintf("/bca/partials/budgets/%s/%s", projectId.String(), budgetItemId.String())

	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret", -5)
	token := createToken(s.TokenAuth)

	testData := []struct {
		name         string
		form         url.Values
		status       int
		body         []string
		updateBudget *mocks.Service_UpdateBudget_Call
		getBudgets   *mocks.Service_GetBudgets_Call
	}{
		{
			name: "should error on non numeric quantity",
			form: url.Values{
				"quantity": []string{"test"},
			},
			status:       http.StatusBadRequest,
			body:         []string{"La cantidad debe ser un número"},
			updateBudget: nil,
			getBudgets:   nil,
		},
		{
			name: "should error on non numeric cost",
			form: url.Values{
				"cost": []string{"test"},
			},
			status:       http.StatusBadRequest,
			body:         []string{"El costo debe ser un número"},
			updateBudget: nil,
			getBudgets:   nil,
		},
		{
			name: "should update a budget",
			form: url.Values{
				"quantity": []string{"10.0"},
				"cost":     []string{"10.0"},
			},
			status: http.StatusOK,
			body:   []string{},
			updateBudget: db.EXPECT().UpdateBudget(types.CreateBudget{
				ProjectId:    projectId,
				BudgetItemId: budgetItemId,
				Quantity:     10.0,
				Cost:         10.0,
			}, types.Budget{
				ProjectId:         projectId,
				BudgetItemId:      budgetItemId,
				InitialQuantity:   sql.NullFloat64{Float64: 10.0, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 10.0, Valid: true},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0.0, Valid: true},
				SpentTotal:        0.0,
				RemainingQuantity: sql.NullFloat64{Float64: 10.0, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 10.0, Valid: true},
				RemainingTotal:    100.0,
				UpdatedBudget:     100.0,
				CompanyId:         uuid.UUID{},
			}).Return(nil),
			getBudgets: db.EXPECT().GetBudgets(uuid.UUID{}, uuid.UUID{}, "").Return([]types.GetBudget{
				{
					Project: types.ProjectData{
						ID:        projectId,
						Name:      "test",
						NetArea:   10.0,
						GrossArea: 10.0,
					},
					BudgetItem: types.BudgetItemData{
						ID:         budgetItemId,
						Code:       "test",
						Name:       "test",
						Level:      0,
						Accumulate: false,
					},
					InitialQuantity:   sql.NullFloat64{Float64: 10.0, Valid: true},
					InitialCost:       sql.NullFloat64{Float64: 10.0, Valid: true},
					InitialTotal:      100.0,
					SpentQuantity:     sql.NullFloat64{Float64: 0.0, Valid: true},
					SpentTotal:        0.0,
					RemainingQuantity: sql.NullFloat64{Float64: 10.0, Valid: true},
					RemainingCost:     sql.NullFloat64{Float64: 10.0, Valid: true},
					RemainingTotal:    100.0,
					UpdatedBudget:     100.0,
					CompanyId:         uuid.UUID{},
				},
			}, nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			db.EXPECT().GetOneBudget(uuid.UUID{}, projectId, budgetItemId).Return(&types.GetBudget{
				Project: types.ProjectData{
					ID:        projectId,
					Name:      "test",
					NetArea:   10.0,
					GrossArea: 10.0,
				},
				BudgetItem: types.BudgetItemData{
					ID:         budgetItemId,
					Code:       "test",
					Name:       "test",
					Level:      0,
					Accumulate: false,
				},
				InitialQuantity:   sql.NullFloat64{Float64: 10.0, Valid: true},
				InitialCost:       sql.NullFloat64{Float64: 10.0, Valid: true},
				InitialTotal:      100.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0.0, Valid: true},
				SpentTotal:        0.0,
				RemainingQuantity: sql.NullFloat64{Float64: 10.0, Valid: true},
				RemainingCost:     sql.NullFloat64{Float64: 10.0, Valid: true},
				RemainingTotal:    100.0,
				UpdatedBudget:     100.0,
				CompanyId:         uuid.UUID{},
			}, nil).Times(1)

			if tt.updateBudget != nil {
				tt.updateBudget.Times(1)
			}

			if tt.getBudgets != nil {
				tt.getBudgets.Times(1)
			}

			req, res := createRequest(
				token,
				http.MethodPut,
				testURL,
				strings.NewReader(tt.form.Encode()),
			)
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}
