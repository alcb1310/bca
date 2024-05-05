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

func TestDetailsTable(t *testing.T) {
	invoiceId := uuid.New()
	testURL := fmt.Sprintf("/bca/partials/invoices/%s/details", invoiceId)
	muxVars := make(map[string]string)
	muxVars["invoiceId"] = invoiceId.String()
	budgetItemId := uuid.New()
	quantity := 1.0
	cost := 1.0

	t.Run("method GET", func(t *testing.T) {
		t.Run("Successful Query", func(t *testing.T) {
			srv, db := server.MakeServer()
			db.On("GetAllDetails", invoiceId, uuid.UUID{}).Return([]types.InvoiceDetailsResponse{}, nil)

			request, response := server.MakeRequest(http.MethodGet, testURL, nil)
			request = mux.SetURLVars(request, muxVars)
			srv.DetailsTable(response, request)
			assert.Equal(t, http.StatusOK, response.Code)
		})

		t.Run("Error Query", func(t *testing.T) {
			srv, db := server.MakeServer()
			db.On("GetAllDetails", invoiceId, uuid.UUID{}).Return([]types.InvoiceDetailsResponse{}, UnknownError)

			request, response := server.MakeRequest(http.MethodGet, testURL, nil)
			request = mux.SetURLVars(request, muxVars)
			srv.DetailsTable(response, request)
			assert.Equal(t, http.StatusInternalServerError, response.Code)
			assert.Contains(t, response.Body.String(), UnknownError.Error())
		})
	})

	t.Run("method POST", func(t *testing.T) {
		t.Run("Data validation", func(t *testing.T) {
			t.Run("BudgetItem (item)", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					form := url.Values{}
					form.Add("item", "")
					buf := strings.NewReader(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.DetailsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "Seleccione un partida")
				})

				t.Run("invalid", func(t *testing.T) {
					form := url.Values{}
					form.Add("item", "invalid")
					buf := strings.NewReader(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.DetailsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "invalid UUID length: 7")
				})
			})

			t.Run("Quantity (quantity)", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					form := url.Values{}
					form.Add("item", budgetItemId.String())
					form.Add("quantity", "")
					buf := strings.NewReader(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.DetailsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "cantidad es requerido")
				})

				t.Run("invalid", func(t *testing.T) {
					form := url.Values{}
					form.Add("item", budgetItemId.String())
					form.Add("quantity", "invalid")
					buf := strings.NewReader(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.DetailsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "cantidad debe ser un número válido")
				})
			})

			t.Run("Cost (cost)", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					form := url.Values{}
					form.Add("item", budgetItemId.String())
					form.Add("quantity", fmt.Sprintf("%f", quantity))
					form.Add("cost", "")
					buf := strings.NewReader(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.DetailsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "costo es requerido")
				})

				t.Run("invalid", func(t *testing.T) {
					form := url.Values{}
					form.Add("item", budgetItemId.String())
					form.Add("quantity", fmt.Sprintf("%f", quantity))
					form.Add("cost", "invalid")
					buf := strings.NewReader(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.DetailsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "costo debe ser un número válido")
				})
			})
		})

		t.Run("Valid data", func(t *testing.T) {
			invoiceCreate := types.InvoiceDetailCreate{
				CompanyId:    uuid.UUID{},
				InvoiceId:    invoiceId,
				BudgetItemId: budgetItemId,
				Quantity:     quantity,
				Cost:         cost,
				Total:        quantity * cost,
			}

			t.Run("successfull", func(t *testing.T) {
				form := url.Values{}
				form.Add("item", budgetItemId.String())
				form.Add("quantity", fmt.Sprintf("%f", quantity))
				form.Add("cost", fmt.Sprintf("%f", cost))
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("AddDetail", invoiceCreate).Return(nil)
				db.On("GetAllDetails", invoiceId, uuid.UUID{}).Return([]types.InvoiceDetailsResponse{}, nil)

				request, response := server.MakeRequest(http.MethodPost, testURL, buf)
				request = mux.SetURLVars(request, muxVars)
				srv.DetailsTable(response, request)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("error", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					form := url.Values{}
					form.Add("item", budgetItemId.String())
					form.Add("quantity", fmt.Sprintf("%f", quantity))
					form.Add("cost", fmt.Sprintf("%f", cost))
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("AddDetail", invoiceCreate).Return(errors.New("duplicate"))

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.DetailsTable(response, request)
					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Equal(t, response.Body.String(), "Ya existe una partida con ese nombre en la factura")
				})

				t.Run("unknown", func(t *testing.T) {
					form := url.Values{}
					form.Add("item", budgetItemId.String())
					form.Add("quantity", fmt.Sprintf("%f", quantity))
					form.Add("cost", fmt.Sprintf("%f", cost))
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("AddDetail", invoiceCreate).Return(UnknownError)

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.DetailsTable(response, request)
					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Contains(t, response.Body.String(), UnknownError.Error())
				})

				t.Run("no rows", func(t *testing.T) {
					form := url.Values{}
					form.Add("item", budgetItemId.String())
					form.Add("quantity", fmt.Sprintf("%f", quantity))
					form.Add("cost", fmt.Sprintf("%f", cost))
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("AddDetail", invoiceCreate).Return(sql.ErrNoRows)

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.DetailsTable(response, request)
					assert.Equal(t, http.StatusNotFound, response.Code)
					assert.Contains(t, response.Body.String(), "No existe presupuesto para esa partida")
				})
			})
		})
	})
}

func TestDetailsAdd(t *testing.T) {
	invoiceId := uuid.New()
	testURL := fmt.Sprintf("/bca/partials/invoices/%s/details", invoiceId)
	muxVars := make(map[string]string)
	muxVars["invoiceId"] = invoiceId.String()

	srv, db := server.MakeServer()
	db.On("GetBudgetItemsByAccumulate", uuid.UUID{}, false).Return([]types.BudgetItem{})

	request, response := server.MakeRequest(http.MethodGet, testURL, nil)
	request = mux.SetURLVars(request, muxVars)
	srv.DetailsAdd(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Agregar Detalles")
}
