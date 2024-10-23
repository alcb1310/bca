package server_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/alcb1310/bca/internal/server"
	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/mocks"
)

func TestCreateInvoice(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret", -5)
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

func TestUpdateInvoice(t *testing.T) {
	invoiceId := uuid.New()
	testUrl := fmt.Sprintf("/bca/partials/invoices/%s", invoiceId.String())

	projectId := uuid.New()
	supplierId := uuid.New()
	truePointer := true

	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret", -5)
	token := createToken(s.TokenAuth)

	invoiceNumber := "test"
	invoiceDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	testData := []struct {
		name            string
		form            url.Values
		status          int
		body            []string
		getInvoiceCalls int
		updateInvoice   *mocks.Service_UpdateInvoice_Call
	}{
		{
			name: "should pass a valid supplier id",
			form: url.Values{
				"supplier": []string{"test"},
			},
			status:          http.StatusBadRequest,
			body:            []string{"Código del proveedor inválido"},
			getInvoiceCalls: 1,
			updateInvoice:   nil,
		},
		{
			name: "should pass a non empty invoice number",
			form: url.Values{
				"supplier":      []string{supplierId.String()},
				"invoiceNumber": []string{""},
			},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese un número de factura"},
			getInvoiceCalls: 1,
			updateInvoice:   nil,
		},
		{
			name: "should pass a non empty date",
			form: url.Values{
				"supplier":      []string{supplierId.String()},
				"invoiceNumber": []string{"test"},
				"invoiceDate":   []string{""},
			},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese una fecha"},
			getInvoiceCalls: 1,
			updateInvoice:   nil,
		},
		{
			name: "should pass a valid date",
			form: url.Values{
				"supplier":      []string{supplierId.String()},
				"invoiceNumber": []string{"test"},
				"invoiceDate":   []string{"test"},
			},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese una fecha válida"},
			getInvoiceCalls: 1,
			updateInvoice:   nil,
		},
		{
			name: "should update an invoice",
			form: url.Values{
				"supplier":      []string{supplierId.String()},
				"invoiceNumber": []string{"test"},
				"invoiceDate":   []string{"2024-01-01"},
			},
			status:          http.StatusOK,
			body:            []string{},
			getInvoiceCalls: 2,
			updateInvoice: db.EXPECT().UpdateInvoice(types.InvoiceCreate{
				Id:            &invoiceId,
				SupplierId:    &supplierId,
				ProjectId:     &projectId,
				InvoiceNumber: &invoiceNumber,
				InvoiceDate:   &invoiceDate,
				CompanyId:     uuid.UUID{},
			}).Return(nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			db.EXPECT().GetActiveProjects(uuid.UUID{}, true).Return([]types.Project{}).Times(1)
			db.EXPECT().GetAllSuppliers(uuid.UUID{}, "").Return([]types.Supplier{}, nil).Times(1)

			if tt.updateInvoice != nil {
				tt.updateInvoice.Times(1)
			}

			if tt.getInvoiceCalls != 0 {
				db.EXPECT().GetOneInvoice(invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{
					Id: invoiceId,
					Project: types.Project{
						ID:        projectId,
						Name:      "test",
						IsActive:  &truePointer,
						CompanyId: uuid.UUID{},
						GrossArea: 0,
						NetArea:   0,
					},
					Supplier: types.Supplier{
						ID:         supplierId,
						Name:       "test",
						SupplierId: "test",
						CompanyId:  uuid.UUID{},
					},
					InvoiceNumber: "test",
					InvoiceDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					InvoiceTotal:  0,
					IsBalanced:    false,
					CompanyId:     uuid.UUID{},
				}, nil).Times(tt.getInvoiceCalls)
			}

			req, res := createRequest(token, http.MethodPut, testUrl, strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)

			for _, v := range tt.body {
				assert.Contains(t, res.Body.String(), v)
			}
		})
	}
}
