package server_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
)

func TestRubrosTable(t *testing.T) {
	testURL := "/bca/partials/rubros"
	srv, db := server.MakeServer()
	db.On("GetAllRubros", uuid.UUID{}).Return([]types.Rubro{}, nil)

	request, response := server.MakeRequest(http.MethodGet, testURL, nil)
	srv.RubrosTable(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
}
