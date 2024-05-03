package server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bca-go-final/mocks"
)

var (
	trueValue = true
	companyId = uuid.New()
	projectId = uuid.New()
)

var UnknownError = errors.New("unknown error")

func TestBudget(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{
		{
			ID:        uuid.New(),
			CompanyId: companyId,
			Name:      "Proyecto 1",
			IsActive:  &trueValue,
			GrossArea: 1000.0,
			NetArea:   500.0,
		},
	}, nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/transacciones/presupuesto", nil)

	srv.Budget(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Presupuesto")
	assert.Contains(t, response.Body.String(), "Proyecto 1")
}

func TestInvoice(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/transacciones/facturas", nil)

	srv.Invoice(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Facturas")
}

func TestClosure(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{
		{
			ID:        projectId,
			CompanyId: companyId,
			Name:      "Proyecto 1",
			IsActive:  &trueValue,
			GrossArea: 1000.0,
			NetArea:   500.0,
		},
	})

	t.Run("GET Method", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/bca/transacciones/cierre", nil)

		srv.Closure(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), "Cierre Mensual")
		assert.Contains(t, response.Body.String(), "Proyecto 1")
	})

	t.Run("POST Method", func(t *testing.T) {
		t.Run("valid input", func(t *testing.T) {
			form := url.Values{}
			form.Add("proyecto", projectId.String())
			form.Add("date", "2022-01-01")
			reader := strings.NewReader(form.Encode())

			date, _ := time.Parse("2006-01-02", "2022-01-01")

			db.On("CreateClosure", uuid.UUID{}, projectId, date).Return(nil)

			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/bca/transacciones/cierre", reader)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			srv.Closure(response, request)

			assert.Equal(t, http.StatusOK, response.Code)

		})

		t.Run("invalid input", func(t *testing.T) {
			t.Run("project_id", func(t *testing.T) {
				t.Run("empty project_id", func(t *testing.T) {
					form := url.Values{}
					form.Add("proyecto", "")
					reader := strings.NewReader(form.Encode())

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPost, "/bca/transacciones/cierre", reader)
					request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

					srv.Closure(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "Seleccione un proyecto")
				})

				t.Run("invalid project id", func(t *testing.T) {
					form := url.Values{}
					form.Add("proyecto", "invalid")
					reader := strings.NewReader(form.Encode())

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPost, "/bca/transacciones/cierre", reader)
					request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

					srv.Closure(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "invalid UUID")
				})
			})

			t.Run("date", func(t *testing.T) {
				t.Run("empty date", func(t *testing.T) {
					form := url.Values{}
					form.Add("proyecto", companyId.String())
					form.Add("date", "")
					reader := strings.NewReader(form.Encode())

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPost, "/bca/transacciones/cierre", reader)
					request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

					srv.Closure(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "Seleccione una fecha")
				})

				t.Run("invalid date", func(t *testing.T) {
					form := url.Values{}
					form.Add("proyecto", companyId.String())
					form.Add("date", "invalid")
					reader := strings.NewReader(form.Encode())

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPost, "/bca/transacciones/cierre", reader)
					request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

					srv.Closure(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "Fecha inválida")
				})
			})
		})

		t.Run("DB error", func(t *testing.T) {
			db := mocks.NewServiceMock()
			_, srv := server.NewServer(db)

			db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{
				{
					ID:        projectId,
					CompanyId: companyId,
					Name:      "Proyecto 1",
					IsActive:  &trueValue,
					GrossArea: 1000.0,
					NetArea:   500.0,
				},
			})

			form := url.Values{}
			form.Add("proyecto", projectId.String())
			form.Add("date", "2022-01-01")
			reader := strings.NewReader(form.Encode())

			date, _ := time.Parse("2006-01-02", "2022-01-01")
			_ = date

			db.On("CreateClosure", uuid.UUID{}, projectId, date).Return(UnknownError)

			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/bca/transacciones/cierre", reader)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			srv.Closure(response, request)

			assert.Equal(t, http.StatusInternalServerError, response.Code)
			assert.Contains(t, response.Body.String(), "No se pudo cerrar el proyecto")
		})
	})
}
