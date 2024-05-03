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

func TestUnitQuantity(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/costo-unitario/cantidades", nil)

	srv.UnitQuantity(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Cantidades")
}

func TestUnitAnalysis(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/costo-unitario/analisis", nil)

	db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{
		{
			ID:        uuid.New(),
			Name:      "Project 1",
			IsActive:  &trueValue,
			CompanyId: uuid.New(),
		},
	})

	srv.UnitAnalysis(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Analisis")
}
