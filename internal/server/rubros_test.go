package server_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func TestMaterialsByItem(t *testing.T) {
	rubroId := uuid.New()
	testURL := fmt.Sprintf("/bca/partials/rubros/%s", rubroId.String())
	muxVars := make(map[string]string)
	muxVars["id"] = rubroId.String()

	t.Run("invalid ID", func(t *testing.T) {
		testURL := fmt.Sprintf("/bca/partials/rubros/invalid")
		muxVars := make(map[string]string)
		muxVars["id"] = "invalid"

		srv, _ := server.MakeServer()
		request, response := server.MakeRequest(http.MethodGet, testURL, nil)
		request = mux.SetURLVars(request, muxVars)
		srv.MaterialsByItem(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("method GET", func(t *testing.T) {
		srv, db := server.MakeServer()
		db.On("GetMaterialsByItem", rubroId, uuid.UUID{}).Return([]types.ACU{}, nil)

		request, response := server.MakeRequest(http.MethodGet, testURL, nil)
		request = mux.SetURLVars(request, muxVars)
		srv.MaterialsByItem(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})
}
