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

func TestCategoriesTable(t *testing.T) {
	t.Run("Method GET", func(t *testing.T) {
		srv, db := server.MakeServer()
		db.On("GetAllCategories", uuid.UUID{}).Return([]types.Category{}, nil)

		request, response := server.MakeRequest(http.MethodGet, "/bca/partials/categories", nil)
		srv.CategoriesTable(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Method POST", func(t *testing.T) {
		t.Run("data validation", func(t *testing.T) {
			t.Run("name", func(t *testing.T) {
				form := url.Values{}
				form.Add("name", "")
				buf := strings.NewReader(form.Encode())

				srv, _ := server.MakeServer()
				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/categories", buf)
				srv.CategoriesTable(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "nombre es requerido")
			})
		})

		t.Run("valid data", func(t *testing.T) {
			t.Run("successful", func(t *testing.T) {
				form := url.Values{}
				form.Add("name", "cat")
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("CreateCategory", types.Category{
					Name:      "cat",
					CompanyId: uuid.UUID{},
				}).Return(nil)
				db.On("GetAllCategories", uuid.UUID{}).Return([]types.Category{}, nil)

				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/categories", buf)
				srv.CategoriesTable(response, request)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("error", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					form := url.Values{}
					form.Add("name", "cat")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("CreateCategory", types.Category{
						Name:      "cat",
						CompanyId: uuid.UUID{},
					}).Return(errors.New("duplicate"))

					request, response := server.MakeRequest(http.MethodPost, "/bca/partials/categories", buf)
					srv.CategoriesTable(response, request)
					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Contains(t, response.Body.String(), "La categoria cat ya existe")
				})

				t.Run("other", func(t *testing.T) {
					form := url.Values{}
					form.Add("name", "cat")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("CreateCategory", types.Category{
						Name:      "cat",
						CompanyId: uuid.UUID{},
					}).Return(UnknownError)

					request, response := server.MakeRequest(http.MethodPost, "/bca/partials/categories", buf)
					srv.CategoriesTable(response, request)
					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Contains(t, response.Body.String(), UnknownError.Error())
				})
			})
		})
	})
}

func TestCategoryAdd(t *testing.T) {
	srv, _ := server.MakeServer()
	request, response := server.MakeRequest(http.MethodGet, "/bca/partials/categories/add", nil)
	srv.CategoryAdd(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Agregar Categoría")
}
