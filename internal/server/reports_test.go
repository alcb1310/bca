package server_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bca-go-final/mocks"
)

func TestActual(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{
		{
			ID:        uuid.UUID{},
			Name:      "1",
			CompanyId: companyId,
			IsActive:  &trueValue,
		},
	})

	db.On("Levels", uuid.UUID{}).Return([]types.Select{
		{
			Key:   "1",
			Value: "1",
		},
	}, nil)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/reportes/actual", nil)

	srv.Actual(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Actual")
}

func TestActualGenerate(t *testing.T) {
	var lev uint8 = 0
	t.Run("valid data", func(t *testing.T) {
		db := mocks.NewServiceMock()
		_, srv := server.NewServer(db)

		form := url.Values{}
		form.Add("proyecto", projectId.String())
		form.Add("nivel", strconv.Itoa(int(lev)))
		reader := strings.NewReader(form.Encode())

		db.On("GetBudgetsByProjectId", uuid.UUID{}, projectId, &lev).Return([]types.GetBudget{}, nil)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/bca/reportes/actual/generar", reader)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		srv.ActualGenerate(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("data validation", func(t *testing.T) {
		t.Run("level", func(t *testing.T) {
			db := mocks.NewServiceMock()
			_, srv := server.NewServer(db)

			form := url.Values{}
			form.Add("proyecto", projectId.String())
			form.Add("nivel", "nivel")
			reader := strings.NewReader(form.Encode())
			db.On("GetBudgetsByProjectId", uuid.UUID{}, uuid.UUID{}, &lev).Return([]types.GetBudget{}, nil)

			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/bca/reportes/actual/generar", reader)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			srv.ActualGenerate(response, request)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Contains(t, response.Body.String(), "nivel debe ser un número válido")
		})

	})

	t.Run("Database Error", func(t *testing.T) {
		db := mocks.NewServiceMock()
		_, srv := server.NewServer(db)

		form := url.Values{}
		form.Add("proyecto", projectId.String())
		form.Add("nivel", strconv.Itoa(int(lev)))
		reader := strings.NewReader(form.Encode())

		db.On("GetBudgetsByProjectId", uuid.UUID{}, uuid.UUID{}, &lev).Return([]types.GetBudget{}, UnknownError)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/bca/reportes/actual/generar", reader)

		srv.ActualGenerate(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
	})
}

func TestBalance(t *testing.T) {
	t.Run("GET Method", func(t *testing.T) {
		db := mocks.NewServiceMock()
		_, srv := server.NewServer(db)

		db.On("Levels", uuid.UUID{}).Return([]types.Select{
			{
				Key:   "1",
				Value: "1",
			},
		})

		db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{
			{
				ID:        uuid.UUID{},
				Name:      "1",
				CompanyId: companyId,
				IsActive:  &trueValue,
			},
		}, nil)

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/bca/reportes/cuadre", nil)

		srv.Balance(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), "Cuadre")
	})

	t.Run("POST Method", func(t *testing.T) {
		db := mocks.NewServiceMock()
		_, srv := server.NewServer(db)

		form := url.Values{}
		form.Add("proyecto", projectId.String())
		form.Add("date", "2022-01-01")
		reader := strings.NewReader(form.Encode())

		date := time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC)

		db.On("GetBalance", uuid.UUID{}, uuid.UUID{}, date).Return(types.BalanceResponse{})

		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/bca/reportes/cuadre", reader)

		srv.Balance(response, request)

		assert.Equal(t, http.StatusOK, response.Code)

	})
}
