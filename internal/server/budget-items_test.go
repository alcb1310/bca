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

func TestBudgetItemsTable(t *testing.T) {
	t.Run("method not allowed", func(t *testing.T) {
		srv, _ := server.MakeServer()
		request, response := server.MakeRequest(http.MethodPut, "/bca/partials/budget-items", nil)
		srv.BudgetItemsTable(response, request)
		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})

	t.Run("method GET", func(t *testing.T) {
		srv, db := server.MakeServer()
		db.On("GetBudgetItems", uuid.UUID{}, "").Return([]types.BudgetItemResponse{}, nil)

		request, response := server.MakeRequest(http.MethodGet, "/bca/partials/budget-items", nil)
		srv.BudgetItemsTable(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), "Partida")
	})

	t.Run("method POST", func(t *testing.T) {
		parentId := uuid.New()
		budgetItem := types.BudgetItem{
			ID:         uuid.UUID{},
			Code:       "1",
			Name:       "1",
			ParentId:   nil,
			Accumulate: sql.NullBool{Valid: true, Bool: false},
			CompanyId:  uuid.UUID{},
		}

		t.Run("data validation", func(t *testing.T) {
			t.Run("code", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", "")
				buf := strings.NewReader(form.Encode())

				srv, _ := server.MakeServer()
				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/budget-items", buf)
				srv.BudgetItemsTable(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Debe proporcionar un código de la partida")
			})

			t.Run("name", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", budgetItem.Code)
				form.Add("name", "")
				buf := strings.NewReader(form.Encode())

				srv, _ := server.MakeServer()
				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/budget-items", buf)
				srv.BudgetItemsTable(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Debe proporcionar un nombre de la partida")
			})

			t.Run("parent_id", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", budgetItem.Code)
				form.Add("name", budgetItem.Name)
				form.Add("parent", "invalid")
				buf := strings.NewReader(form.Encode())

				srv, _ := server.MakeServer()
				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/budget-items", buf)
				srv.BudgetItemsTable(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Código de la partida padre es inválido")
			})
		})

		t.Run("valid data", func(t *testing.T) {
			form := url.Values{}
			form.Add("code", budgetItem.Code)
			form.Add("name", budgetItem.Name)
			form.Add("parent", "")
			form.Add("accumulate", "")
			buf := strings.NewReader(form.Encode())

			t.Run("success", func(t *testing.T) {
				srv, db := server.MakeServer()
				db.On("GetBudgetItems", uuid.UUID{}, "").Return([]types.BudgetItemResponse{}, nil)
				db.On("CreateBudgetItem", &budgetItem).Return(nil)

				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/budget-items", buf)
				srv.BudgetItemsTable(response, request)
				assert.Equal(t, http.StatusOK, response.Code)
				assert.Contains(t, response.Body.String(), "Partida")
			})

			t.Run("fail", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					buf := strings.NewReader(form.Encode())
					srv, db := server.MakeServer()
					db.On("CreateBudgetItem", &budgetItem).Return(errors.New("duplicate"))

					request, response := server.MakeRequest(http.MethodPost, "/bca/partials/budget-items", buf)
					srv.BudgetItemsTable(response, request)
					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Equal(t, response.Body.String(), fmt.Sprintf("Ya existe una partida con el mismo código: %s y/o el mismo nombre: %s", budgetItem.Code, budgetItem.Name))
				})

				t.Run("unknown errorr", func(t *testing.T) {
					form := url.Values{}
					form.Add("code", budgetItem.Code)
					form.Add("name", budgetItem.Name)
					form.Add("parent", parentId.String())
					form.Add("accumulate", "")
					buf := strings.NewReader(form.Encode())
					srv, db := server.MakeServer()
					budgetItem.ParentId = &parentId

					db.On("CreateBudgetItem", &budgetItem).Return(UnknownError)

					request, response := server.MakeRequest(http.MethodPost, "/bca/partials/budget-items", buf)
					srv.BudgetItemsTable(response, request)
					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Contains(t, response.Body.String(), UnknownError.Error())
					budgetItem.ParentId = nil
				})
			})
		})
	})
}

func TestBudgetItemsAdd(t *testing.T) {
	srv, db := server.MakeServer()
	db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, true).Return([]types.BudgetItem{})

	request, response := server.MakeRequest(http.MethodGet, "/bca/partials/budget-items", nil)
	srv.BudgetItemAdd(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Agregar Partida")
}

func TestBudgetItemEdit(t *testing.T) {
	budgetItemId := uuid.New()
	testURL := fmt.Sprintf("/bca/partials/budget-item/%s", budgetItemId.String())
	muxVars := make(map[string]string)
	muxVars["id"] = budgetItemId.String()

	budgetItem := types.BudgetItem{
		CompanyId:  uuid.UUID{},
		ID:         uuid.UUID{},
		Code:       "code",
		Name:       "name",
		ParentId:   nil,
		Accumulate: sql.NullBool{Valid: true, Bool: false},
	}

	t.Run("method GET", func(t *testing.T) {
		srv, db := server.MakeServer()
		db.On("GetOneBudgetItem", budgetItemId, uuid.UUID{}).Return(&types.BudgetItem{}, nil)
		db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, true).Return([]types.BudgetItem{})

		request, response := server.MakeRequest(http.MethodGet, testURL, nil)
		request = mux.SetURLVars(request, muxVars)
		srv.BudgetItemEdit(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), "Editar Partida")
	})

	t.Run("method PUT", func(t *testing.T) {
		t.Run("data validation", func(t *testing.T) {
			t.Run("code", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", "")
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("GetOneBudgetItem", budgetItemId, uuid.UUID{}).Return(&types.BudgetItem{}, nil)

				request, response := server.MakeRequest(http.MethodPut, testURL, buf)
				request = mux.SetURLVars(request, muxVars)
				srv.BudgetItemEdit(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Debe proporcionar un código de la partida")
			})

			t.Run("name", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", budgetItem.Code)
				form.Add("name", "")
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("GetOneBudgetItem", budgetItemId, uuid.UUID{}).Return(&types.BudgetItem{}, nil)

				request, response := server.MakeRequest(http.MethodPut, testURL, buf)
				request = mux.SetURLVars(request, muxVars)
				srv.BudgetItemEdit(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Debe proporcionar un nombre de la partida")
			})
			t.Run("parent", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", budgetItem.Code)
				form.Add("name", budgetItem.Name)
				form.Add("parent", "invalid")
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("GetOneBudgetItem", budgetItemId, uuid.UUID{}).Return(&types.BudgetItem{}, nil)

				request, response := server.MakeRequest(http.MethodPut, testURL, buf)
				request = mux.SetURLVars(request, muxVars)
				srv.BudgetItemEdit(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Código de la partida padre es inválido")
			})
		})

		t.Run("valid data", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", budgetItem.Code)
				form.Add("name", budgetItem.Name)
				form.Add("parent", "")
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("GetOneBudgetItem", budgetItemId, uuid.UUID{}).Return(&types.BudgetItem{}, nil)
				db.On("UpdateBudgetItem", &budgetItem).Return(nil)
				db.On("GetBudgetItems", uuid.UUID{}, "").Return([]types.BudgetItemResponse{}, nil)

				request, response := server.MakeRequest(http.MethodPut, testURL, buf)
				request = mux.SetURLVars(request, muxVars)
				srv.BudgetItemEdit(response, request)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("fail", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					form := url.Values{}
					form.Add("code", budgetItem.Code)
					form.Add("name", budgetItem.Name)
					form.Add("parent", budgetItemId.String())
					buf := strings.NewReader(form.Encode())

					budgetItem.ParentId = &budgetItemId
					srv, db := server.MakeServer()
					db.On("GetOneBudgetItem", budgetItemId, uuid.UUID{}).Return(&types.BudgetItem{}, nil)
					db.On("UpdateBudgetItem", &budgetItem).Return(errors.New("duplicate"))

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.BudgetItemEdit(response, request)
					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Equal(t, response.Body.String(), fmt.Sprintf("Ya existe una partida con el mismo código: %s y/o el mismo nombre: %s", budgetItem.Code, budgetItem.Name))
					budgetItem.ParentId = nil
				})

				t.Run("unknown error", func(t *testing.T) {
					form := url.Values{}
					form.Add("code", budgetItem.Code)
					form.Add("name", budgetItem.Name)
					form.Add("parent", "")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetOneBudgetItem", budgetItemId, uuid.UUID{}).Return(&types.BudgetItem{}, nil)
					db.On("UpdateBudgetItem", &budgetItem).Return(UnknownError)

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.BudgetItemEdit(response, request)
					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Equal(t, response.Body.String(), UnknownError.Error())
				})
			})
		})
	})
}
