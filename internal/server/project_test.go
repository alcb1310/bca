package server_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
	"bca-go-final/mocks"
)

func TestCreateProject(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	testData := []struct {
		name           string
		form           url.Values
		status         int
		body           []string
		createProject  *mocks.Service_CreateProject_Call
		getAllProjects *mocks.Service_GetAllProjects_Call
	}{
		{
			name:           "should pass a form",
			form:           nil,
			status:         http.StatusBadRequest,
			body:           []string{},
			createProject:  nil,
			getAllProjects: nil,
		},
		{
			name:   "should pass a name",
			form:   url.Values{},
			status: http.StatusBadRequest,
			body: []string{
				"Ingrese un valor para el nombre",
			},
			createProject:  nil,
			getAllProjects: nil,
		},
		{
			name: "should pass a valid number for gross area",
			form: url.Values{
				"name":       []string{"test"},
				"gross_area": []string{"test"},
			},
			status: http.StatusBadRequest,
			body: []string{
				"El área bruta debe ser un número válido",
			},
			createProject:  nil,
			getAllProjects: nil,
		},
		{
			name: "should pass a valid number for net area",
			form: url.Values{
				"name":       []string{"test"},
				"gross_area": []string{"1"},
				"net_area":   []string{"test"},
			},
			status: http.StatusBadRequest,
			body: []string{
				"El área neta debe ser un número válido",
			},
			createProject:  nil,
			getAllProjects: nil,
		},
		{
			name: "should create a project",
			form: url.Values{
				"name":       []string{"test"},
				"gross_area": []string{"1"},
				"net_area":   []string{"1"},
			},
			status: http.StatusOK,
			body:   []string{},
			createProject: db.EXPECT().CreateProject(types.Project{
				Name:      "test",
				GrossArea: 1,
				NetArea:   1,
				IsActive:  new(bool),
				CompanyId: uuid.UUID{},
			}).Return(types.Project{
				ID:        uuid.UUID{},
				Name:      "test",
				GrossArea: 1,
				NetArea:   1,
				IsActive:  new(bool),
				CompanyId: uuid.UUID{},
			}, nil),
			getAllProjects: db.EXPECT().GetAllProjects(uuid.UUID{}).Return([]types.Project{
				{
					ID:        uuid.UUID{},
					Name:      "test",
					GrossArea: 1,
					NetArea:   1,
					IsActive:  new(bool),
					CompanyId: uuid.UUID{},
				},
			}, nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			req, res := createRequest(token, http.MethodPost, "/bca/partials/projects", strings.NewReader(tt.form.Encode()))

			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)
			if len(tt.body) != 0 {
				for _, b := range tt.body {
					assert.Contains(t, res.Body.String(), b)
				}
			}
		})
	}
}
