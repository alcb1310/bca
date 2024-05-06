package server_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

func TestMaterialByItemForm(t *testing.T) {
	rubroId := uuid.New()
	testURL := fmt.Sprintf("/bca/partials/rubros/%s/material", rubroId.String())
	muxVars := make(map[string]string)
	muxVars["id"] = rubroId.String()

	t.Run("invalid ID", func(t *testing.T) {
		testURL := fmt.Sprintf("/bca/partials/rubros/invalid/material")
		muxVars := make(map[string]string)
		muxVars["id"] = "invalid"

		srv, _ := server.MakeServer()
		request, response := server.MakeRequest(http.MethodGet, testURL, nil)
		request = mux.SetURLVars(request, muxVars)
		srv.MaterialByItemForm(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, "invalid UUID length: 7", response.Body.String())
	})

	t.Run("method not allowed", func(t *testing.T) {
		srv, _ := server.MakeServer()
		request, response := server.MakeRequest(http.MethodPut, testURL, nil)
		request = mux.SetURLVars(request, muxVars)
		srv.MaterialByItemForm(response, request)
		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})

	t.Run("method GET", func(t *testing.T) {
		srv, db := server.MakeServer()
		db.On("GetAllMaterials", uuid.UUID{}).Return([]types.Material{}, nil)

		request, response := server.MakeRequest(http.MethodGet, testURL, nil)
		request = mux.SetURLVars(request, muxVars)
		srv.MaterialByItemForm(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("method POST", func(t *testing.T) {
		materialId := uuid.New()
		quantity := 1.0

		t.Run("invalid data", func(t *testing.T) {
			t.Run("material", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					form := url.Values{}
					form.Add("material", "")
					buf := bytes.NewBufferString(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.MaterialByItemForm(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, "Seleccione un material", response.Body.String())
				})

				t.Run("invalid", func(t *testing.T) {
					form := url.Values{}
					form.Add("material", "invaild")
					buf := bytes.NewBufferString(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.MaterialByItemForm(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, "invalid UUID length: 7", response.Body.String())
				})
			})

			t.Run("quantity", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					form := url.Values{}
					form.Add("material", materialId.String())
					form.Add("quantity", "")
					buf := bytes.NewBufferString(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.MaterialByItemForm(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, "cantidad es requerido", response.Body.String())
				})

				t.Run("invalid", func(t *testing.T) {
					form := url.Values{}
					form.Add("material", materialId.String())
					form.Add("quantity", "invalid")
					buf := bytes.NewBufferString(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.MaterialByItemForm(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, "cantidad debe ser un número válido", response.Body.String())
				})
			})
		})

		t.Run("valid data", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				form := url.Values{}
				form.Add("material", materialId.String())
				form.Add("quantity", fmt.Sprintf("%f", quantity))
				buf := bytes.NewBufferString(form.Encode())

				srv, db := server.MakeServer()
				db.On("AddMaterialsByItem", rubroId, materialId, quantity, uuid.UUID{}).Return(nil)
				db.On("GetMaterialsByItem", rubroId, uuid.UUID{}).Return([]types.ACU{}, nil)

				request, response := server.MakeRequest(http.MethodPost, testURL, buf)
				request = mux.SetURLVars(request, muxVars)
				srv.MaterialByItemForm(response, request)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("error", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					form := url.Values{}
					form.Add("material", materialId.String())
					form.Add("quantity", fmt.Sprintf("%f", quantity))
					buf := bytes.NewBufferString(form.Encode())

					srv, db := server.MakeServer()
					db.On("AddMaterialsByItem", rubroId, materialId, quantity, uuid.UUID{}).Return(errors.New("duplicate"))

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.MaterialByItemForm(response, request)
					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Equal(t, "Ya existe un material con ese Código", response.Body.String())
				})

				t.Run("unknown", func(t *testing.T) {
					form := url.Values{}
					form.Add("material", materialId.String())
					form.Add("quantity", fmt.Sprintf("%f", quantity))
					buf := bytes.NewBufferString(form.Encode())

					srv, db := server.MakeServer()
					db.On("AddMaterialsByItem", rubroId, materialId, quantity, uuid.UUID{}).Return(UnknownError)

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.MaterialByItemForm(response, request)
					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Equal(t, UnknownError.Error(), response.Body.String())
				})
			})
		})
	})
}

// TEST: r.HandleFunc("/bca/partials/rubros/{id}/material/{materialId}", s.MaterialItemsOperations)
