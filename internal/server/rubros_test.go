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
	"bca-go-final/mocks"
)

func TestCreateRubros(t *testing.T) {
	rubroId := uuid.New()
	materialId := uuid.New()
	testUrl := fmt.Sprintf("/bca/partials/rubros/%s/material", rubroId.String())

	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	testData := []struct {
		name            string
		form            url.Values
		status          int
		body            []string
		createRubro     *mocks.Service_AddMaterialsByItem_Call
		getAllMaterials *mocks.Service_GetMaterialsByItem_Call
	}{
		{
			name:            "should pass a form",
			form:            nil,
			status:          http.StatusBadRequest,
			body:            []string{},
			createRubro:     nil,
			getAllMaterials: nil,
		},
		{
			name:            "should pass a material",
			form:            url.Values{},
			status:          http.StatusBadRequest,
			body:            []string{"Seleccione un Material"},
			createRubro:     nil,
			getAllMaterials: nil,
		},
		{
			name:            "the material should be a valid uuid",
			form:            url.Values{"material": {"invalid"}},
			status:          http.StatusBadRequest,
			body:            []string{"Material Incorrecto"},
			createRubro:     nil,
			getAllMaterials: nil,
		},
		{
			name:            "should pass a quantity",
			form:            url.Values{"material": {materialId.String()}},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese una Cantidad"},
			createRubro:     nil,
			getAllMaterials: nil,
		},
		{
			name:            "should pass a number for quantity",
			form:            url.Values{"material": {materialId.String()}, "quantity": {"invalid"}},
			status:          http.StatusBadRequest,
			body:            []string{"Cantidad debe ser un valor numérico"},
			createRubro:     nil,
			getAllMaterials: nil,
		},
		{
			name:            "should pass a number greater than 0 for quantity",
			form:            url.Values{"material": {materialId.String()}, "quantity": {"-1.4"}},
			status:          http.StatusBadRequest,
			body:            []string{"La Cantidad debe ser mayor a 0"},
			createRubro:     nil,
			getAllMaterials: nil,
		},
		{
			name:   "should creaate a rubro",
			form:   url.Values{"material": {materialId.String()}, "quantity": {"1.4"}},
			status: http.StatusOK,
			body:   []string{""},
			createRubro: db.EXPECT().AddMaterialsByItem(
				rubroId,
				materialId,
				1.4,
				uuid.UUID{},
			).Return(nil),
			getAllMaterials: db.EXPECT().GetMaterialsByItem(rubroId, uuid.UUID{}).Return([]types.ACU{}),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createRubro != nil {
				tt.createRubro.Times(1)
			}

			if tt.getAllMaterials != nil {
				tt.getAllMaterials.Times(1)
			}

			req, res := createRequest(token, http.MethodPost, testUrl, strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)
			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}

func TestUpdateRubros(t *testing.T) {
	rubroId := uuid.New()
	materialId := uuid.New()
	testUrl := fmt.Sprintf("/bca/partials/rubros/%s/material/%s", rubroId.String(), materialId.String())

	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	testData := []struct {
		name            string
		form            url.Values
		status          int
		body            []string
		updateRubro     *mocks.Service_UpdateMaterialByItem_Call
		getAllMaterials *mocks.Service_GetMaterialsByItem_Call
	}{
		{
			name:            "should have a quantity",
			form:            url.Values{},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese una Cantidad"},
			updateRubro:     nil,
			getAllMaterials: nil,
		},
		{
			name:            "should pass a number for quantity",
			form:            url.Values{"quantity": {"invalid"}},
			status:          http.StatusBadRequest,
			body:            []string{"Cantidad debe ser un valor numérico"},
			updateRubro:     nil,
			getAllMaterials: nil,
		},
		{
			name:            "should pass a number greater than 0 for quantity",
			form:            url.Values{"quantity": {"-1.35"}},
			status:          http.StatusBadRequest,
			body:            []string{"La Cantidad debe ser mayor a 0"},
			updateRubro:     nil,
			getAllMaterials: nil,
		},
		{
			name:            "should update a rubro",
			form:            url.Values{"quantity": {"1.35"}},
			status:          http.StatusOK,
			body:            []string{""},
			updateRubro:     db.EXPECT().UpdateMaterialByItem(rubroId, materialId, 1.35, uuid.UUID{}).Return(nil),
			getAllMaterials: db.EXPECT().GetMaterialsByItem(rubroId, uuid.UUID{}).Return([]types.ACU{}),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.updateRubro != nil {
				tt.updateRubro.Times(1)
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
