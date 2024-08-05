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

func TestCreateInvoice(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	projectId := uuid.New()
	supplierId := uuid.New()

	testData := []struct {
		name          string
		form          url.Values
		status        int
		body          []string
		createInvoice *mocks.Service_CreateInvoice_Call
		getOneInvoice *mocks.Service_GetOneInvoice_Call
	}{
		{
			name:          "should pass a form",
			form:          nil,
			status:        http.StatusBadRequest,
			body:          []string{},
			createInvoice: nil,
			getOneInvoice: nil,
		},
		{
			name:   "should pass a project id",
			form:   url.Values{},
			status: http.StatusBadRequest,
			body: []string{
				"Ingrese un proyecto",
			},
			createInvoice: nil,
			getOneInvoice: nil,
		},
		{
			name: "should pass a valid project id",
			form: url.Values{
				"project": []string{"test"},
			},
			status:        http.StatusBadRequest,
			body:          []string{"Código del proyecto inválido"},
			createInvoice: nil,
			getOneInvoice: nil,
		},
		{
			name: "should pass a supplier id",
			form: url.Values{
				"project": []string{projectId.String()},
			},
			status:        http.StatusBadRequest,
			body:          []string{"Ingrese un proveedor"},
			createInvoice: nil,
			getOneInvoice: nil,
		},
		{
			name: "should pass a valid supplier id",
			form: url.Values{
				"project":  []string{projectId.String()},
				"supplier": []string{"test"},
			},
			status:        http.StatusBadRequest,
			body:          []string{"Código del proveedor inválido"},
			createInvoice: nil,
			getOneInvoice: nil,
		},
		{
			name: "should pass a valid invoice number",
			form: url.Values{
				"project":  []string{projectId.String()},
				"supplier": []string{supplierId.String()},
			},
			status:        http.StatusBadRequest,
			body:          []string{"Ingrese un número de factura"},
			createInvoice: nil,
			getOneInvoice: nil,
		},
		{
			name: "should pass a date",
			form: url.Values{
				"project":       []string{projectId.String()},
				"supplier":      []string{supplierId.String()},
				"invoiceNumber": []string{"test"},
			},
			status:        http.StatusBadRequest,
			body:          []string{"Ingrese una fecha"},
			createInvoice: nil,
			getOneInvoice: nil,
		},
		{
			name: "should pass a valid date",
			form: url.Values{
				"project":       []string{projectId.String()},
				"supplier":      []string{supplierId.String()},
				"invoiceNumber": []string{"test"},
				"invoiceDate":   []string{"test"},
			},
			status:        http.StatusBadRequest,
			body:          []string{"Ingrese una fecha válida"},
			createInvoice: nil,
			getOneInvoice: nil,
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			db.EXPECT().GetActiveProjects(uuid.UUID{}, true).Return([]types.Project{}).Times(1)
			db.EXPECT().GetAllSuppliers(uuid.UUID{}, "").Return([]types.Supplier{}, nil).Times(1)

			if tt.createInvoice != nil {
				tt.createInvoice.Times(1)
			}

			if tt.getOneInvoice != nil {
				tt.getOneInvoice.Times(1)
			}

			req, res := createRequest(token, http.MethodPost, "/bca/transacciones/facturas/crear", strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)

			assert.Equal(t, tt.status, res.Code)
			for _, v := range tt.body {
				assert.Contains(t, res.Body.String(), v)
			}
		})
	}
}
