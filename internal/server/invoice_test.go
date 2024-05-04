package server_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
)

func TestInvoicesTable(t *testing.T) {
	srv, db := server.MakeServer()
	db.On("GetInvoices", uuid.UUID{}).Return([]types.InvoiceResponse{}, nil)

	request, response := server.MakeRequest(http.MethodGet, "/bca/partials/invoices", nil)
	srv.InvoicesTable(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}
