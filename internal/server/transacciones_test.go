package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bca-go-final/mocks"
)

var (
	trueValue = true
)

func TestBudget(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{
		{
			ID:        uuid.UUID{},
			CompanyId: uuid.UUID{},
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
