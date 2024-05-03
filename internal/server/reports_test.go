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
