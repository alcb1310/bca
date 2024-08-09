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

func TestCreateCantidades(t *testing.T) {
	testUrl := "/bca/partials/cantidades/add"

	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	projectId := uuid.New()
	itemId := uuid.New()

	testData := []struct {
		name             string
		form             url.Values
		status           int
		body             []string
		createCantidades *mocks.Service_CreateCantidades_Call
		cantidadesTable  *mocks.Service_CantidadesTable_Call
	}{
		{
			name:             "should pass a form",
			form:             nil,
			status:           http.StatusBadRequest,
			body:             []string{},
			createCantidades: nil,
			cantidadesTable:  nil,
		},
		{
			name:             "should pass a project",
			form:             url.Values{},
			status:           http.StatusBadRequest,
			body:             []string{"Seleccione un proyecto"},
			createCantidades: nil,
			cantidadesTable:  nil,
		},
		{
			name:             "should pass a valid project id",
			form:             url.Values{"project": []string{"invalid"}},
			status:           http.StatusBadRequest,
			body:             []string{"Código de proyecto inválido"},
			createCantidades: nil,
			cantidadesTable:  nil,
		},
		{
			name:             "should pass an item",
			form:             url.Values{"project": []string{projectId.String()}},
			status:           http.StatusBadRequest,
			body:             []string{"Seleccione un rubro"},
			createCantidades: nil,
			cantidadesTable:  nil,
		},
		{
			name: "should pass a valid item",
			form: url.Values{
				"project": []string{projectId.String()},
				"item":    []string{"invalid"},
			},
			status:           http.StatusBadRequest,
			body:             []string{"Código de rubro inválido"},
			createCantidades: nil,
			cantidadesTable:  nil,
		},
		{
			name: "should pass a quantity",
			form: url.Values{
				"project": []string{projectId.String()},
				"item":    []string{itemId.String()},
			},
			status:           http.StatusBadRequest,
			body:             []string{"La cantidad es requerida"},
			createCantidades: nil,
			cantidadesTable:  nil,
		},
		{
			name: "Quantity should be a number",
			form: url.Values{
				"project":  []string{projectId.String()},
				"item":     []string{itemId.String()},
				"quantity": []string{"invalid"},
			},
			status:           http.StatusBadRequest,
			body:             []string{"La cantidad debe ser un número"},
			createCantidades: nil,
			cantidadesTable:  nil,
		},
		{
			name: "Quantity should be a positive number",
			form: url.Values{
				"project":  []string{projectId.String()},
				"item":     []string{itemId.String()},
				"quantity": []string{"-1.65"},
			},
			status:           http.StatusBadRequest,
			body:             []string{"La cantidad debe ser mayor a 0"},
			createCantidades: nil,
			cantidadesTable:  nil,
		},
		{
			name: "should create a unit cost",
			form: url.Values{
				"project":  []string{projectId.String()},
				"item":     []string{itemId.String()},
				"quantity": []string{"1.65"},
			},
			status:           http.StatusOK,
			body:             []string{""},
			createCantidades: db.EXPECT().CreateCantidades(projectId, itemId, 1.65, uuid.UUID{}).Return(nil),
			cantidadesTable:  db.EXPECT().CantidadesTable(uuid.UUID{}).Return([]types.Quantity{}),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createCantidades != nil {
				tt.createCantidades.Times(1)
			}

			if tt.cantidadesTable != nil {
				tt.cantidadesTable.Times(1)
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

func TestUpdateCantidades(t *testing.T) {
	unitCostId := uuid.New()
	projectId := uuid.New()
	rubroId := uuid.New()
	testUrl := fmt.Sprintf("/bca/partials/cantidades/%s", unitCostId.String())

	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	testData := []struct {
		name             string
		form             url.Values
		status           int
		body             []string
		getOneQuantity   *mocks.Service_GetOneQuantityById_Call
		updateQuantity   *mocks.Service_UpdateQuantity_Call
		getAllQuantities *mocks.Service_CantidadesTable_Call
	}{
		{
			name:             "should pass a valid form",
			form:             nil,
			status:           http.StatusBadRequest,
			body:             []string{},
			getOneQuantity:   nil,
			updateQuantity:   nil,
			getAllQuantities: nil,
		},
		{
			name:             "should pass a quantity",
			form:             url.Values{},
			status:           http.StatusBadRequest,
			body:             []string{"Ingrese una cantidad"},
			getOneQuantity:   nil,
			updateQuantity:   nil,
			getAllQuantities: nil,
		},
		{
			name:             "should pass a quantity",
			form:             url.Values{"quantity": []string{"invalid"}},
			status:           http.StatusBadRequest,
			body:             []string{"Cantidad debe ser numérica"},
			getOneQuantity:   nil,
			updateQuantity:   nil,
			getAllQuantities: nil,
		},
		{
			name:             "Quantity should be greater than zero",
			form:             url.Values{"quantity": []string{"-2"}},
			status:           http.StatusBadRequest,
			body:             []string{"La cantidad debe ser mayor a 0"},
			getOneQuantity:   nil,
			updateQuantity:   nil,
			getAllQuantities: nil,
		},
		{
			name:   "should update a unit cost",
			form:   url.Values{"quantity": []string{"2"}},
			status: http.StatusOK,
			body:   []string{""},
			getOneQuantity: db.EXPECT().GetOneQuantityById(unitCostId, uuid.UUID{}).Return(types.Quantity{
				Id: unitCostId,
				Project: types.Project{
					ID:        projectId,
					Name:      "test",
					CompanyId: uuid.UUID{},
				},
				Rubro: types.Rubro{
					Id:        rubroId,
					Name:      "test",
					CompanyId: uuid.UUID{},
				},
				Quantity:  1,
				CompanyId: uuid.UUID{},
			}),
			updateQuantity: db.EXPECT().UpdateQuantity(types.Quantity{
				Id: unitCostId,
				Project: types.Project{
					ID:        projectId,
					Name:      "test",
					CompanyId: uuid.UUID{},
				},
				Rubro: types.Rubro{
					Id:        rubroId,
					Name:      "test",
					CompanyId: uuid.UUID{},
				},
				Quantity:  2,
				CompanyId: uuid.UUID{},
			}, uuid.UUID{}).Return(nil),
			getAllQuantities: db.EXPECT().CantidadesTable(uuid.UUID{}).Return([]types.Quantity{}),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.getOneQuantity != nil {
				tt.getOneQuantity.Times(1)
			}

			if tt.updateQuantity != nil {
				tt.updateQuantity.Times(1)
			}

			if tt.getAllQuantities != nil {
				tt.getAllQuantities.Times(1)
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
