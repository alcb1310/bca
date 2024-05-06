package server_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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
func TestCantidadesAdd(t *testing.T) {
	testURL := "/bca/partials/cantidades/add"

	t.Run("method not allowed", func(t *testing.T) {
		srv, _ := server.MakeServer()

		request, response := server.MakeRequest(http.MethodPatch, testURL, nil)
		srv.CantidadesAdd(response, request)
		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})

	t.Run("GET", func(t *testing.T) {
		srv, db := server.MakeServer()
		db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{})
		db.On("GetAllRubros", uuid.UUID{}).Return([]types.Rubro{}, nil)

		request, response := server.MakeRequest(http.MethodGet, testURL, nil)
		srv.CantidadesAdd(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), "Crear Cantidad")
	})

	t.Run("POST", func(t *testing.T) {
		projectId := uuid.New()
		rubroId := uuid.New()
		quantity := 1.0

		t.Run("validate data", func(t *testing.T) {
			t.Run("project", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					srv, _ := server.MakeServer()
					form := url.Values{}
					form.Add("project", "")
					buf := bytes.NewBufferString(form.Encode())

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.CantidadesAdd(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "Seleccione un proyecto")
				})

				t.Run("invalid", func(t *testing.T) {
					srv, _ := server.MakeServer()
					form := url.Values{}
					form.Add("project", "invalid")
					buf := bytes.NewBufferString(form.Encode())

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.CantidadesAdd(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "Seleccione un proyecto")
				})
			})

			t.Run("item", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					srv, _ := server.MakeServer()
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("item", "")
					buf := bytes.NewBufferString(form.Encode())

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.CantidadesAdd(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "Seleccione un rubro")
				})

				t.Run("invalid", func(t *testing.T) {
					srv, _ := server.MakeServer()
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("item", "invalid")
					buf := bytes.NewBufferString(form.Encode())

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.CantidadesAdd(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "Seleccione un rubro")
				})
			})

			t.Run("quantity", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					srv, _ := server.MakeServer()
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("item", rubroId.String())
					form.Add("quantity", "")
					buf := bytes.NewBufferString(form.Encode())

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.CantidadesAdd(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "cantidad es requerido")
				})

				t.Run("invalid", func(t *testing.T) {
					srv, _ := server.MakeServer()
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("item", rubroId.String())
					form.Add("quantity", "invalid")
					buf := bytes.NewBufferString(form.Encode())

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.CantidadesAdd(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "cantidad debe ser un número válido")
				})
			})
		})

		t.Run("valid data", func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				srv, db := server.MakeServer()
				form := url.Values{}
				form.Add("project", projectId.String())
				form.Add("item", rubroId.String())
				form.Add("quantity", fmt.Sprintf("%f", quantity))
				buf := bytes.NewBufferString(form.Encode())

				db.On("CreateCantidades", projectId, rubroId, quantity, uuid.UUID{}).Return(nil)
				db.On("CantidadesTable", uuid.UUID{}).Return([]types.Quantity{})

				request, response := server.MakeRequest(http.MethodPost, testURL, buf)
				srv.CantidadesAdd(response, request)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("error", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					srv, db := server.MakeServer()
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("item", rubroId.String())
					form.Add("quantity", fmt.Sprintf("%f", quantity))
					buf := bytes.NewBufferString(form.Encode())

					db.On("CreateCantidades", projectId, rubroId, quantity, uuid.UUID{}).Return(errors.New("duplicate"))

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.CantidadesAdd(response, request)
					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Equal(t, response.Body.String(), "La cantidad ya existe")
				})

				t.Run("Unknown", func(t *testing.T) {
					srv, db := server.MakeServer()
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("item", rubroId.String())
					form.Add("quantity", fmt.Sprintf("%f", quantity))
					buf := bytes.NewBufferString(form.Encode())

					db.On("CreateCantidades", projectId, rubroId, quantity, uuid.UUID{}).Return(UnknownError)

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.CantidadesAdd(response, request)
					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Equal(t, response.Body.String(), UnknownError.Error())
				})
			})
		})
	})
}
// r.HandleFunc("/bca/partials/cantidades/{id}", s.CantidadesEdit)
