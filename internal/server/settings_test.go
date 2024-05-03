package server_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bca-go-final/mocks"
)

func TestBudgetItems(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/configuracion/partidas", nil)

	srv.BudgetItems(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Partidas")
}

func TestSuppliers(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/configuracion/proveedores", nil)

	srv.Suppliers(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Proveedores")
}

func TestProjects(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/configuracion/proyectos", nil)

	srv.Projects(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Proyectos")
}

func TestCategories(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/configuracion/categorias", nil)

	srv.Categories(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Categorias")
}

func TestMateriales(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/configuracion/materiales", nil)

	srv.Materiales(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Materiales")
}

func TestRubros(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/bca/configuracion/rubros", nil)

	srv.Rubros(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Rubros")
}

func TestRubrosAdd(t *testing.T) {
	db := mocks.NewServiceMock()
	_, srv := server.NewServer(db)

	rubro := types.Rubro{
		Id:        uuid.UUID{},
		Code:      "code",
		Name:      "name",
		Unit:      "unit",
		CompanyId: uuid.UUID{},
	}

	t.Run("GET request", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/bca/configuracion/rubros/crear", nil)

		srv.RubrosAdd(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), "Nuevo Rubro")
	})

	t.Run("POST request", func(t *testing.T) {
		t.Run("valid data", func(t *testing.T) {
			t.Run("successful creation", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", rubro.Code)
				form.Add("name", rubro.Name)
				form.Add("unit", rubro.Unit)
				buf := strings.NewReader(form.Encode())
				db.On("CreateRubro", rubro).Return(uuid.New(), nil)

				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/bca/configuracion/rubros/crear", buf)
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				srv.RubrosAdd(response, request)

				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("failed creation", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					db := mocks.NewServiceMock()
					_, srv := server.NewServer(db)
					form := url.Values{}
					form.Add("code", rubro.Code)
					form.Add("name", rubro.Name)
					form.Add("unit", rubro.Unit)
					buf := strings.NewReader(form.Encode())
					db.On("CreateRubro", rubro).Return(uuid.UUID{}, errors.New("duplicate"))

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPost, "/bca/configuracion/rubros/crear", buf)
					request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

					srv.RubrosAdd(response, request)

					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Contains(t, response.Body.String(), fmt.Sprintf("El Código %s ya existe", rubro.Code))
				})

				t.Run("unknown", func(t *testing.T) {
					db := mocks.NewServiceMock()
					_, srv := server.NewServer(db)
					form := url.Values{}
					form.Add("code", rubro.Code)
					form.Add("name", rubro.Name)
					form.Add("unit", rubro.Unit)
					buf := strings.NewReader(form.Encode())
					db.On("CreateRubro", rubro).Return(uuid.UUID{}, errors.New("unknown"))

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPost, "/bca/configuracion/rubros/crear", buf)
					request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

					srv.RubrosAdd(response, request)

					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Contains(t, response.Body.String(), "unknown")
				})
			})
		})

		t.Run("invalid data", func(t *testing.T) {
			t.Run("code", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", "")
				buf := strings.NewReader(form.Encode())

				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/bca/configuracion/rubros/crear", buf)
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				srv.RubrosAdd(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un valor para el Código")
			})

			t.Run("name", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", "code")
				form.Add("name", "")
				buf := strings.NewReader(form.Encode())

				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/bca/configuracion/rubros/crear", buf)
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				srv.RubrosAdd(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un valor para el Nombre")
			})

			t.Run("unit", func(t *testing.T) {
				form := url.Values{}
				form.Add("code", "code")
				form.Add("name", "name")
				form.Add("unit", "")
				buf := strings.NewReader(form.Encode())

				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/bca/configuracion/rubros/crear", buf)
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				srv.RubrosAdd(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un valor para la Unidad")
			})
		})
	})

	t.Run("PUT request", func(t *testing.T) {
		rubroId := uuid.New()

		t.Run("valid data", func(t *testing.T) {
			t.Run("successful creation", func(t *testing.T) {
				db := mocks.NewServiceMock()
				_, srv := server.NewServer(db)
				urlQuery := fmt.Sprintf("/bca/configuracion/rubros?id=%s", rubroId)

				form := url.Values{}
				form.Add("code", rubro.Code)
				form.Add("name", rubro.Name)
				form.Add("unit", rubro.Unit)
				buf := strings.NewReader(form.Encode())

				rubro.Id = rubroId
				db.On("GetOneRubro", rubroId, uuid.UUID{}).Return(types.Rubro{}, nil)
				db.On("UpdateRubro", rubro).Return(nil)

				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPut, urlQuery, buf)
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				srv.RubrosAdd(response, request)

				assert.Equal(t, http.StatusOK, response.Code)
				rubro.Id = uuid.UUID{}
			})

			t.Run("failed creation", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					db := mocks.NewServiceMock()
					_, srv := server.NewServer(db)
					urlQuery := fmt.Sprintf("/bca/configuracion/rubros?id=%s", rubroId)

					form := url.Values{}
					form.Add("code", rubro.Code)
					form.Add("name", rubro.Name)
					form.Add("unit", rubro.Unit)
					buf := strings.NewReader(form.Encode())

					rubro.Id = rubroId
					db.On("GetOneRubro", rubroId, uuid.UUID{}).Return(types.Rubro{}, nil)
					db.On("UpdateRubro", rubro).Return(errors.New("duplicate"))

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPut, urlQuery, buf)
					request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

					srv.RubrosAdd(response, request)

					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Contains(t, response.Body.String(), fmt.Sprintf("El Código %s ya existe", rubro.Code))
					rubro.Id = uuid.UUID{}
				})

				t.Run("unknown", func(t *testing.T) {
					db := mocks.NewServiceMock()
					_, srv := server.NewServer(db)
					urlQuery := fmt.Sprintf("/bca/configuracion/rubros?id=%s", rubroId)

					form := url.Values{}
					form.Add("code", rubro.Code)
					form.Add("name", rubro.Name)
					form.Add("unit", rubro.Unit)
					buf := strings.NewReader(form.Encode())

					rubro.Id = rubroId
					db.On("GetOneRubro", rubroId, uuid.UUID{}).Return(types.Rubro{}, nil)
					db.On("UpdateRubro", rubro).Return(UnknownError)

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPut, urlQuery, buf)
					request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

					srv.RubrosAdd(response, request)

					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Contains(t, response.Body.String(), UnknownError.Error())
				})
			})
		})

		t.Run("invalid data", func(t *testing.T) {
			t.Run("ID", func(t *testing.T) {
				t.Run("invalid", func(t *testing.T) {
					db := mocks.NewServiceMock()
					_, srv := server.NewServer(db)
					url := fmt.Sprintf("/bca/configuracion/rubros?id=%s", "invalid")

					response := httptest.NewRecorder()
					request := httptest.NewRequest(http.MethodPut, url, nil)

					srv.RubrosAdd(response, request)

					assert.Equal(t, http.StatusBadRequest, response.Code)
				})
			})

			t.Run("code", func(t *testing.T) {
				db := mocks.NewServiceMock()
				_, srv := server.NewServer(db)
				urlQuery := fmt.Sprintf("/bca/configuracion/rubros?id=%s", rubroId)

				form := url.Values{}
				form.Add("code", "")
				buf := strings.NewReader(form.Encode())

				db.On("GetOneRubro", rubroId, uuid.UUID{}).Return(types.Rubro{}, nil)
				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPut, urlQuery, buf)
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				srv.RubrosAdd(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un valor para el Código")
			})

			t.Run("name", func(t *testing.T) {
				db := mocks.NewServiceMock()
				_, srv := server.NewServer(db)
				urlQuery := fmt.Sprintf("/bca/configuracion/rubros?id=%s", rubroId)

				form := url.Values{}
				form.Add("code", "code")
				form.Add("name", "")
				buf := strings.NewReader(form.Encode())

				db.On("GetOneRubro", rubroId, uuid.UUID{}).Return(types.Rubro{}, nil)
				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPut, urlQuery, buf)
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				srv.RubrosAdd(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un valor para el Nombre")
			})

			t.Run("unit", func(t *testing.T) {
				db := mocks.NewServiceMock()
				_, srv := server.NewServer(db)
				urlQuery := fmt.Sprintf("/bca/configuracion/rubros?id=%s", rubroId)

				form := url.Values{}
				form.Add("code", "code")
				form.Add("name", "name")
				form.Add("unit", "")
				buf := strings.NewReader(form.Encode())

				db.On("GetOneRubro", rubroId, uuid.UUID{}).Return(types.Rubro{}, nil)
				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPut, urlQuery, buf)
				request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

				srv.RubrosAdd(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un valor para la Unidad")
			})

			t.Run("not found", func(t *testing.T) {
				db := mocks.NewServiceMock()
				_, srv := server.NewServer(db)
				urlQuery := fmt.Sprintf("/bca/configuracion/rubros?id=%s", rubroId)

				form := url.Values{}
				form.Add("code", "code")
				form.Add("name", "name")
				form.Add("unit", "unit")
				buf := strings.NewReader(form.Encode())

				db.On("GetOneRubro", rubroId, uuid.UUID{}).Return(types.Rubro{}, sql.ErrNoRows)
				response := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPut, urlQuery, buf)

				srv.RubrosAdd(response, request)

				assert.Equal(t, http.StatusNotFound, response.Code)
			})
		})
	})
}
