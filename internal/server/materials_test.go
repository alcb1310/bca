package server_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/alcb1310/bca/internal/server"
	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/mocks"
)

func TestCreateMaterial(t *testing.T) {
	categoryId := uuid.New()

	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret", -5)
	token := createToken(s.TokenAuth)

	testData := []struct {
		name            string
		form            url.Values
		status          int
		body            []string
		createMaterial  *mocks.Service_CreateMaterial_Call
		getAllMaterials *mocks.Service_GetAllMaterials_Call
	}{
		{
			name:            "should pass a form",
			form:            nil,
			status:          http.StatusBadRequest,
			body:            []string{},
			createMaterial:  nil,
			getAllMaterials: nil,
		},
		{
			name:            "should pass a code",
			form:            url.Values{},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese un valor para el Código"},
			createMaterial:  nil,
			getAllMaterials: nil,
		},
		{
			name: "should pass a name",
			form: url.Values{
				"code": {"123"},
			},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese un valor para el Nombre"},
			createMaterial:  nil,
			getAllMaterials: nil,
		},
		{
			name: "shoule pass a unit",
			form: url.Values{
				"code": {"123"},
				"name": {"name"},
			},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese un valor para la Unidad"},
			createMaterial:  nil,
			getAllMaterials: nil,
		},
		{
			name: "should pass a category",
			form: url.Values{
				"code": {"123"},
				"name": {"name"},
				"unit": {"unit"},
			},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese un valor para la Categoría"},
			createMaterial:  nil,
			getAllMaterials: nil,
		},
		{
			name: "should pass a valid uuid for a category",
			form: url.Values{
				"code":     {"123"},
				"name":     {"name"},
				"unit":     {"unit"},
				"category": {"invalid"},
			},
			status:          http.StatusBadRequest,
			body:            []string{},
			createMaterial:  nil,
			getAllMaterials: nil,
		},
		{
			name: "should create a category",
			form: url.Values{
				"code":     {"123"},
				"name":     {"name"},
				"unit":     {"unit"},
				"category": {categoryId.String()},
			},
			status: http.StatusOK,
			body:   []string{},
			createMaterial: db.EXPECT().CreateMaterial(types.Material{
				Code:      "123",
				Name:      "name",
				Unit:      "unit",
				Category:  types.Category{Id: categoryId},
				CompanyId: uuid.UUID{},
			}).Return(nil),
			getAllMaterials: db.EXPECT().GetAllMaterials(uuid.UUID{}).Return([]types.Material{}),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createMaterial != nil {
				tt.createMaterial.Times(1)
			}
			if tt.getAllMaterials != nil {
				tt.getAllMaterials.Times(1)
			}

			req, res := createRequest(token, http.MethodPost, "/bca/partials/materiales", strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}

func TestUpdateMaterials(t *testing.T) {
	materialId := uuid.New()
	categoryId := uuid.New()
	testUrl := fmt.Sprintf("/bca/partials/materiales/%s", materialId)

	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret", -5)
	token := createToken(s.TokenAuth)

	testData := []struct {
		name            string
		form            url.Values
		status          int
		body            []string
		updateMaterial  *mocks.Service_UpdateMaterial_Call
		getAllMaterials *mocks.Service_GetAllMaterials_Call
	}{
		{
			name: "should pass a valid category id",
			form: url.Values{
				"category": {"invalid"},
			},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese un valor para la Categoría"},
			updateMaterial:  nil,
			getAllMaterials: nil,
		},
		{
			name:   "should update a material",
			form:   url.Values{},
			status: http.StatusOK,
			body:   []string{},
			updateMaterial: db.EXPECT().UpdateMaterial(types.Material{
				Id:        materialId,
				Code:      "123",
				Name:      "name",
				Unit:      "unit",
				Category:  types.Category{Id: categoryId},
				CompanyId: uuid.UUID{},
			}).Return(nil),
			getAllMaterials: db.EXPECT().GetAllMaterials(uuid.UUID{}).Return([]types.Material{}),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			db.EXPECT().GetMaterial(materialId, uuid.UUID{}).Return(types.Material{
				Id:        materialId,
				Code:      "123",
				Name:      "name",
				Unit:      "unit",
				Category:  types.Category{Id: categoryId},
				CompanyId: uuid.UUID{},
			}, nil)

			if tt.updateMaterial != nil {
				tt.updateMaterial.Times(1)
			}

			if tt.getAllMaterials != nil {
				tt.getAllMaterials.Times(1)
			}

			req, res := createRequest(token, http.MethodPut, testUrl, strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}
