package server_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
)

func TestUnitQuantity(t *testing.T) {
	srv, _ := server.MakeServer()
	request, response := server.MakeRequest(http.MethodGet, "/bca/costo-unitario/analisis", nil)

	srv.UnitQuantity(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Cantidades")
}

func TestUnitAnalysis(t *testing.T) {
	srv, db := server.MakeServer()
	request, response := server.MakeRequest(http.MethodGet, "/bca/costo-unitario/analisis", nil)

	db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{})

	srv.UnitAnalysis(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Analisis")
}

// r.HandleFunc("/bca/partials/cantidades", s.CantidadesTable)
func TestCantidadesTable(t *testing.T) {
	srv, db := server.MakeServer()
	db.On("CantidadesTable", uuid.UUID{}).Return([]types.Quantity{})

	request, response := server.MakeRequest(http.MethodGet, "/bca/partials/cantidades", nil)
	srv.CantidadesTable(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
}

// r.HandleFunc("/bca/partials/cantidades/add", s.CantidadesAdd)
// r.HandleFunc("/bca/partials/cantidades/{id}", s.CantidadesEdit)
