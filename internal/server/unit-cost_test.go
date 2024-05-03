package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
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
