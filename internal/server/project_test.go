package server_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
)

func TestProjectsTable(t *testing.T) {
	t.Run("method not allowed", func(t *testing.T) {
		srv, _ := server.MakeServer()

		request, response := server.MakeRequest(http.MethodPut, "/bca/partials/projects", nil)

		srv.ProjectsTable(response, request)

		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})

	t.Run("method GET", func(t *testing.T) {
		srv, db := server.MakeServer()

		db.On("GetAllProjects", uuid.UUID{}).Return([]types.Project{}, nil)

		request, response := server.MakeRequest(http.MethodGet, "/bca/partials/projects", nil)

		srv.ProjectsTable(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("method POST", func(t *testing.T) {
		falseValue := false
		srv, db := server.MakeServer()
		project := types.Project{
			ID:        uuid.UUID{},
			Name:      "proyecto",
			CompanyId: uuid.UUID{},
			IsActive:  &falseValue,
			GrossArea: 100.0,
			NetArea:   80.0,
		}

		t.Run("validate data", func(t *testing.T) {
			t.Run("Nombre", func(t *testing.T) {
				form := url.Values{}
				form.Add("name", "")
				request, response := server.MakeRequest(http.MethodPost, "/bca/projects", strings.NewReader(form.Encode()))

				srv.ProjectsTable(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Equal(t, "El nombre del proyecto es requerido", response.Body.String())
			})

			t.Run("Area bruta", func(t *testing.T) {
				form := url.Values{}
				form.Add("name", "proyecto")
				form.Add("gross_area", "invalid")
				request, response := server.MakeRequest(http.MethodPost, "/bca/projects", strings.NewReader(form.Encode()))

				srv.ProjectsTable(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Equal(t, response.Body.String(), "área bruta debe ser un número válido")
			})

			t.Run("Area net", func(t *testing.T) {
				form := url.Values{}
				form.Add("name", "proyecto")
				form.Add("gross_area", "100")
				form.Add("net_area", "invalid")
				request, response := server.MakeRequest(http.MethodPost, "/bca/projects", strings.NewReader(form.Encode()))

				srv.ProjectsTable(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Equal(t, response.Body.String(), "área útil debe ser un número válido")
			})
		})

		t.Run("valid data", func(t *testing.T) {
			form := url.Values{}
			form.Add("name", project.Name)
			form.Add("gross_area", fmt.Sprintf("%f", project.GrossArea))
			form.Add("net_area", fmt.Sprintf("%f", project.NetArea))

			t.Run("successfull creation", func(t *testing.T) {
				db.On("CreateProject", project).Return(types.Project{}, nil)
				db.On("GetAllProjects", uuid.UUID{}).Return([]types.Project{}, nil)

				request, response := server.MakeRequest(http.MethodPost, "/bca/projects", strings.NewReader(form.Encode()))

				srv.ProjectsTable(response, request)

				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("failed creation", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					srv, db := server.MakeServer()
					db.On("CreateProject", project).Return(types.Project{}, fmt.Errorf("duplicate"))

					request, response := server.MakeRequest(http.MethodPost, "/bca/projects", strings.NewReader(form.Encode()))

					srv.ProjectsTable(response, request)

					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Equal(t, response.Body.String(), fmt.Sprintf("El nombre %s ya existe", project.Name))
				})

				t.Run("other error", func(t *testing.T) {
					srv, db := server.MakeServer()
					db.On("CreateProject", project).Return(types.Project{}, UnknownError)

					request, response := server.MakeRequest(http.MethodPost, "/bca/projects", strings.NewReader(form.Encode()))

					srv.ProjectsTable(response, request)

					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Equal(t, response.Body.String(), fmt.Sprintf("%s", UnknownError.Error()))
				})
			})
		})
	})
}

func TestProjectAdd(t *testing.T) {
	srv, _ := server.MakeServer()

	request, response := server.MakeRequest(http.MethodGet, "/bca/projects/add", nil)

	srv.ProjectAdd(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Agregar Proyecto")
}
