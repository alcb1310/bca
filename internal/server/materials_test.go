package server_test

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
)

func TestMaterialsTable(t *testing.T) {
	testURL := "/bca/partials/materiales"

	t.Run("method GET", func(t *testing.T) {
		srv, db := server.MakeServer()
		db.On("GetAllMaterials", uuid.UUID{}).Return([]types.Material{})

		request, response := server.MakeRequest(http.MethodGet, testURL, nil)
		srv.MaterialsTable(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("method POST", func(t *testing.T) {
		material := types.Material{
			Code:      "1",
			Name:      "1",
			Unit:      "1",
			Category:  types.Category{Id: uuid.UUID{}},
			CompanyId: uuid.UUID{},
		}

		t.Run("data validation", func(t *testing.T) {
			t.Run("code", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", "")
				buf := strings.NewReader(form.Encode())

				srv, _ := server.MakeServer()
				request, response := server.MakeRequest(http.MethodPost, testURL, buf)
				srv.MaterialsTable(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un valor para el Código")
			})

			t.Run("name", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", material.Code)
				form.Add("name", "")
				buf := strings.NewReader(form.Encode())

				srv, _ := server.MakeServer()
				request, response := server.MakeRequest(http.MethodPost, testURL, buf)
				srv.MaterialsTable(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un valor para el Nombre")
			})

			t.Run("unit", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", material.Code)
				form.Add("name", material.Name)
				form.Add("unit", "")
				buf := strings.NewReader(form.Encode())

				srv, _ := server.MakeServer()
				request, response := server.MakeRequest(http.MethodPost, testURL, buf)
				srv.MaterialsTable(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un valor para la Unidad")
			})

			t.Run("category", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					form := url.Values{}
					form.Add("code", material.Code)
					form.Add("name", material.Name)
					form.Add("unit", material.Unit)
					form.Add("category", "")
					buf := strings.NewReader(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.MaterialsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "Seleccione un categoria")
				})

				t.Run("invalid", func(t *testing.T) {
					form := url.Values{}
					form.Add("code", material.Code)
					form.Add("name", material.Name)
					form.Add("unit", material.Unit)
					form.Add("category", "invalid")
					buf := strings.NewReader(form.Encode())

					srv, _ := server.MakeServer()
					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.MaterialsTable(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "invalid UUID length: 7")
				})
			})
		})

		t.Run("valid data", func(t *testing.T) {
			t.Run("successful", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", material.Code)
				form.Add("name", material.Name)
				form.Add("unit", material.Unit)
				form.Add("category", material.Category.Id.String())
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("CreateMaterial", material).Return(nil)
				db.On("GetAllMaterials", uuid.UUID{}).Return([]types.Material{material})

				request, response := server.MakeRequest(http.MethodPost, testURL, buf)
				srv.MaterialsTable(response, request)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("error", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					form := url.Values{}
					form.Add("code", material.Code)
					form.Add("name", material.Name)
					form.Add("unit", material.Unit)
					form.Add("category", material.Category.Id.String())
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("CreateMaterial", material).Return(errors.New("duplicate"))

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.MaterialsTable(response, request)
					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Equal(t, response.Body.String(), "El Código 1 ya existe")
				})

				t.Run("general", func(t *testing.T) {
					form := url.Values{}
					form.Add("code", material.Code)
					form.Add("name", material.Name)
					form.Add("unit", material.Unit)
					form.Add("category", material.Category.Id.String())
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("CreateMaterial", material).Return(UnknownError)

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.MaterialsTable(response, request)
					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Equal(t, response.Body.String(), UnknownError.Error())
				})
			})
		})
	})
}

func TestMaterialsAdd(t *testing.T) {
	testURL := "/bca/partials/materiales/add"
	srv, db := server.MakeServer()
	db.On("GetAllCategories", uuid.UUID{}).Return([]types.Category{}, nil)

	request, response := server.MakeRequest(http.MethodGet, testURL, nil)
	srv.MaterialsAdd(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Agregar Material")
}
