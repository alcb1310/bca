package server_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
)

func TestBudgetsTable(t *testing.T) {
	testURL := "/bca/partials/budgets"
	projectID := uuid.New()
	budgetItemID := uuid.New()
	budgetCreate := &types.CreateBudget{
		ProjectId:    projectID,
		BudgetItemId: budgetItemID,
		Quantity:     10,
		Cost:         0.5,
		CompanyId:    uuid.UUID{},
	}

	t.Run("method not allowed", func(t *testing.T) {
		srv, _ := server.MakeServer()
		request, response := server.MakeRequest(http.MethodPut, testURL, nil)
		srv.BudgetsTable(response, request)
		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})

	t.Run("method get", func(t *testing.T) {
		projectId := uuid.New()
		srv, db := server.MakeServer()
		db.On("GetBudgets", uuid.UUID{}, projectId, "").Return([]types.GetBudget{}, nil)

		testQuery := testURL + "?proyecto=" + projectId.String()

		request, response := server.MakeRequest(http.MethodGet, testQuery, nil)
		srv.BudgetsTable(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("method post", func(t *testing.T) {
		t.Run("data validation", func(t *testing.T) {
			t.Run("project", func(t *testing.T) {
				t.Run("empty project", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", "")
					buf := strings.NewReader(form.Encode())
					srv, _ := server.MakeServer()

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.BudgetsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "Seleccione un proyecto")
				})

				t.Run("invalid project", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", "invalid")
					buf := strings.NewReader(form.Encode())
					srv, _ := server.MakeServer()

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.BudgetsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "invalid UUID length: 7")
				})
			})

			t.Run("budgetItem", func(t *testing.T) {
				t.Run("empty budgetItem", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", budgetCreate.ProjectId.String())
					form.Add("budgetItem", "")
					buf := strings.NewReader(form.Encode())
					srv, _ := server.MakeServer()

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.BudgetsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "Seleccione un partida")
				})

				t.Run("invalid budgetItem", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", budgetCreate.ProjectId.String())
					form.Add("budgetItem", "invalid")
					buf := strings.NewReader(form.Encode())
					srv, _ := server.MakeServer()

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.BudgetsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "invalid UUID length: 7")
				})
			})

			t.Run("quantity", func(t *testing.T) {
				t.Run("empty quantity", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", budgetCreate.ProjectId.String())
					form.Add("budgetItem", budgetCreate.BudgetItemId.String())
					form.Add("quantity", "")
					buf := strings.NewReader(form.Encode())
					srv, _ := server.MakeServer()

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.BudgetsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "cantidad es requerido")
				})

				t.Run("invalid quantity", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", budgetCreate.ProjectId.String())
					form.Add("budgetItem", budgetCreate.BudgetItemId.String())
					form.Add("quantity", "invalid")
					buf := strings.NewReader(form.Encode())
					srv, _ := server.MakeServer()

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.BudgetsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "cantidad debe ser un número válido")
				})
			})

			t.Run("cost", func(t *testing.T) {
				t.Run("empty cost", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", budgetCreate.ProjectId.String())
					form.Add("budgetItem", budgetCreate.BudgetItemId.String())
					form.Add("quantity", fmt.Sprintf("%f", budgetCreate.Quantity))
					form.Add("cost", "")
					buf := strings.NewReader(form.Encode())
					srv, _ := server.MakeServer()

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.BudgetsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "costo es requerido")
				})

				t.Run("invalid cost", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", budgetCreate.ProjectId.String())
					form.Add("budgetItem", budgetCreate.BudgetItemId.String())
					form.Add("quantity", fmt.Sprintf("%f", budgetCreate.Quantity))
					form.Add("cost", "invalid")
					buf := strings.NewReader(form.Encode())
					srv, _ := server.MakeServer()

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.BudgetsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "costo debe ser un número válido")
				})
			})
		})

		t.Run("valid data", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				form := url.Values{}
				form.Add("project", budgetCreate.ProjectId.String())
				form.Add("budgetItem", budgetCreate.BudgetItemId.String())
				form.Add("quantity", fmt.Sprintf("%f", budgetCreate.Quantity))
				form.Add("cost", fmt.Sprintf("%f", budgetCreate.Cost))
				buf := strings.NewReader(form.Encode())
				srv, db := server.MakeServer()
				db.On("CreateBudget", budgetCreate).Return(types.Budget{}, nil)
				db.On("GetBudgets", uuid.UUID{}, uuid.UUID{}, "").Return([]types.GetBudget{}, nil)

				request, response := server.MakeRequest(http.MethodPost, testURL, buf)
				srv.BudgetsTable(response, request)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("fail", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", budgetCreate.ProjectId.String())
					form.Add("budgetItem", budgetCreate.BudgetItemId.String())
					form.Add("quantity", fmt.Sprintf("%f", budgetCreate.Quantity))
					form.Add("cost", fmt.Sprintf("%f", budgetCreate.Cost))
					buf := strings.NewReader(form.Encode())
					srv, db := server.MakeServer()
					db.On("CreateBudget", budgetCreate).Return(types.Budget{}, errors.New("duplicate"))

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.BudgetsTable(response, request)
					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Contains(t, response.Body.String(), fmt.Sprintf("Ya existe partida %s en el proyecto %s", budgetCreate.BudgetItemId.String(), budgetCreate.ProjectId.String()))
				})

				t.Run("unknown", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", budgetCreate.ProjectId.String())
					form.Add("budgetItem", budgetCreate.BudgetItemId.String())
					form.Add("quantity", fmt.Sprintf("%f", budgetCreate.Quantity))
					form.Add("cost", fmt.Sprintf("%f", budgetCreate.Cost))
					buf := strings.NewReader(form.Encode())
					srv, db := server.MakeServer()
					db.On("CreateBudget", budgetCreate).Return(types.Budget{}, UnknownError)

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.BudgetsTable(response, request)
					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Contains(t, response.Body.String(), UnknownError.Error())
				})
			})
		})
	})
}

func TestBudgetsAdd(t *testing.T) {
	srv, db := server.MakeServer()
	db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, false).Return([]types.BudgetItem{}, nil)
	db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{})

	request, response := server.MakeRequest(http.MethodPost, "/bca/budgets/add", nil)
	srv.BudgetAdd(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Agregar Presupuesto")
}

func TestBudgetEdit(t *testing.T) {
	projectId := uuid.New()
	budgetItemId := uuid.New()
	testURL := fmt.Sprintf("/bca/partials/budgets/%s/%s", projectId.String(), budgetItemId.String())
	_ = testURL
	muxVars := make(map[string]string)
	muxVars["projectId"] = projectId.String()
	muxVars["budgetItemId"] = budgetItemId.String()

	t.Run("method not allowed", func(t *testing.T) {
		srv, db := server.MakeServer()
		db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, false).Return([]types.BudgetItem{}, nil)
		db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
		db.On("GetOneBudget", uuid.UUID{}, projectId, budgetItemId).Return(&types.GetBudget{}, nil)

		request, response := server.MakeRequest(http.MethodPost, testURL, nil)
		request = mux.SetURLVars(request, muxVars)
		srv.BudgetEdit(response, request)
		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})

	t.Run("method GET", func(t *testing.T) {
		srv, db := server.MakeServer()
		db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, false).Return([]types.BudgetItem{}, nil)
		db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
		db.On("GetOneBudget", uuid.UUID{}, projectId, budgetItemId).Return(&types.GetBudget{}, nil)

		request, response := server.MakeRequest(http.MethodGet, testURL, nil)
		request = mux.SetURLVars(request, muxVars)
		srv.BudgetEdit(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), "Editar Presupuesto")
	})

	t.Run("method PUT", func(t *testing.T) {
		t.Run("data validation", func(t *testing.T) {
			t.Run("quantity", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("budgetItem", budgetItemId.String())
					form.Add("quantity", "")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, false).Return([]types.BudgetItem{}, nil)
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetOneBudget", uuid.UUID{}, projectId, budgetItemId).Return(&types.GetBudget{}, nil)

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.BudgetEdit(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "cantidad es requerido\n")
				})

				t.Run("invalid", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("budgetItem", budgetItemId.String())
					form.Add("quantity", "invalid")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, false).Return([]types.BudgetItem{}, nil)
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetOneBudget", uuid.UUID{}, projectId, budgetItemId).Return(&types.GetBudget{}, nil)

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.BudgetEdit(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "cantidad debe ser un número válido\n")
				})
			})

			t.Run("cost", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("budgetItem", budgetItemId.String())
					form.Add("quantity", fmt.Sprintf("%f", 1.0))
					form.Add("cost", "")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, false).Return([]types.BudgetItem{}, nil)
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetOneBudget", uuid.UUID{}, projectId, budgetItemId).Return(&types.GetBudget{}, nil)

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.BudgetEdit(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "costo es requerido\n")
				})

				t.Run("invalid", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("budgetItem", budgetItemId.String())
					form.Add("quantity", fmt.Sprintf("%f", 1.0))
					form.Add("cost", "invalid")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, false).Return([]types.BudgetItem{}, nil)
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetOneBudget", uuid.UUID{}, projectId, budgetItemId).Return(&types.GetBudget{}, nil)

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.BudgetEdit(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "costo debe ser un número válido\n")
				})
			})
		})

		t.Run("Valid data", func(t *testing.T) {
			budget := types.CreateBudget{
				ProjectId:    projectId,
				BudgetItemId: budgetItemId,
				Quantity:     1.0,
				Cost:         1.0,
				CompanyId:    uuid.UUID{},
			}

			bu := types.Budget{
				ProjectId:         projectId,
				BudgetItemId:      budgetItemId,
				InitialQuantity:   sql.NullFloat64{Float64: 0, Valid: false},
				InitialCost:       sql.NullFloat64{Float64: 0, Valid: false},
				InitialTotal:      0.0,
				SpentQuantity:     sql.NullFloat64{Float64: 0, Valid: false},
				SpentTotal:        0.0,
				RemainingQuantity: sql.NullFloat64{Float64: 0, Valid: false},
				RemainingCost:     sql.NullFloat64{Float64: 0, Valid: false},
				RemainingTotal:    0.0,
				UpdatedBudget:     0.0,
				CompanyId:         uuid.UUID{},
			}

			t.Run("success", func(t *testing.T) {
				form := url.Values{}
				form.Add("project", projectId.String())
				form.Add("budgetItem", budgetItemId.String())
				form.Add("quantity", fmt.Sprintf("%f", 1.0))
				form.Add("cost", fmt.Sprintf("%f", 1.0))
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, false).Return([]types.BudgetItem{}, nil)
				db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
				db.On("GetOneBudget", uuid.UUID{}, projectId, budgetItemId).Return(&types.GetBudget{}, nil)
				db.On("UpdateBudget", budget, bu).Return(nil)
				db.On("GetBudgets", uuid.UUID{}, uuid.UUID{}, "").Return([]types.GetBudget{}, nil)

				request, response := server.MakeRequest(http.MethodPut, testURL, buf)
				request = mux.SetURLVars(request, muxVars)
				srv.BudgetEdit(response, request)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("success", func(t *testing.T) {
				form := url.Values{}
				form.Add("project", projectId.String())
				form.Add("budgetItem", budgetItemId.String())
				form.Add("quantity", fmt.Sprintf("%f", 1.0))
				form.Add("cost", fmt.Sprintf("%f", 1.0))
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, false).Return([]types.BudgetItem{}, nil)
				db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
				db.On("GetOneBudget", uuid.UUID{}, projectId, budgetItemId).Return(&types.GetBudget{}, nil)
				db.On("UpdateBudget", budget, bu).Return(UnknownError)
				db.On("GetBudgets", uuid.UUID{}, uuid.UUID{}, "").Return([]types.GetBudget{}, nil)

				request, response := server.MakeRequest(http.MethodPut, testURL, buf)
				request = mux.SetURLVars(request, muxVars)
				srv.BudgetEdit(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
			})
		})
	})
}
