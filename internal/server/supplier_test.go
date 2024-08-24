package server_test

import (
	"database/sql"
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

func TestCreateSupplier(t *testing.T) {
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)

	testData := []struct {
		name            string
		form            url.Values
		status          int
		body            []string
		createSupplier  *mocks.Service_CreateSupplier_Call
		getAllSuppliers *mocks.Service_GetAllSuppliers_Call
	}{
		{
			name:            "should pass a form",
			form:            nil,
			status:          http.StatusBadRequest,
			body:            []string{},
			createSupplier:  nil,
			getAllSuppliers: nil,
		},
		{
			name:            "should pass a supplier id",
			form:            url.Values{},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese un valor para el RUC"},
			createSupplier:  nil,
			getAllSuppliers: nil,
		},
		{
			name: "should pass a supplier name",
			form: url.Values{
				"supplier_id": {"test"},
			},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese un valor para el nombre"},
			createSupplier:  nil,
			getAllSuppliers: nil,
		},
		{
			name: "should pass a valid email",
			form: url.Values{
				"supplier_id":   {"test"},
				"name":          {"test"},
				"contact_email": {"test"},
			},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese un correo válido"},
			createSupplier:  nil,
			getAllSuppliers: nil,
		},
		{
			name: "should create a supplier",
			form: url.Values{
				"supplier_id": {"test"},
				"name":        {"test"},
			},
			status: http.StatusOK,
			body:   []string{"test"},
			createSupplier: db.EXPECT().CreateSupplier(&types.Supplier{
				SupplierId:   "test",
				Name:         "test",
				ContactEmail: sql.NullString{Valid: true, String: ""},
				ContactName:  sql.NullString{Valid: true, String: ""},
				ContactPhone: sql.NullString{Valid: true, String: ""},
				CompanyId:    uuid.UUID{},
			}).Return(nil),
			getAllSuppliers: db.EXPECT().GetAllSuppliers(uuid.UUID{}, "").Return([]types.Supplier{
				{
					SupplierId:   "test",
					Name:         "test",
					ContactEmail: sql.NullString{Valid: true, String: "test"},
					ContactName:  sql.NullString{Valid: true, String: ""},
					ContactPhone: sql.NullString{Valid: true, String: ""},
					CompanyId:    uuid.UUID{},
				},
			}, nil),
		},
	}

	for _, tt := range testData {
		if tt.createSupplier != nil {
			tt.createSupplier.Times(1)
		}
		if tt.getAllSuppliers != nil {
			tt.getAllSuppliers.Times(1)
		}

		t.Run(tt.name, func(t *testing.T) {
			req, res := createRequest(token, http.MethodPost, "/bca/partials/suppliers", strings.NewReader(tt.form.Encode()))
			s.Router.ServeHTTP(res, req)
			assert.Equal(t, tt.status, res.Code)
			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}

func TestEditSupplier(t *testing.T) {
	id := uuid.New()
	db := mocks.NewService(t)
	s := server.NewServer(db, "supersecret")
	token := createToken(s.TokenAuth)
	getOne := db.EXPECT().GetOneSupplier(id, uuid.UUID{}).Return(types.Supplier{
		ID:           id,
		SupplierId:   "test",
		Name:         "test",
		ContactEmail: sql.NullString{Valid: true, String: ""},
		ContactName:  sql.NullString{Valid: true, String: ""},
		ContactPhone: sql.NullString{Valid: true, String: ""},
		CompanyId:    uuid.UUID{},
	}, nil)

	testData := []struct {
		name            string
		form            url.Values
		status          int
		body            []string
		getOneSupplier  *mocks.Service_GetOneSupplier_Call
		updateSupplier  *mocks.Service_UpdateSupplier_Call
		getAllSuppliers *mocks.Service_GetAllSuppliers_Call
	}{
		{
			name:            "should pass a form",
			form:            nil,
			status:          http.StatusOK,
			body:            []string{},
			getOneSupplier:  getOne,
			updateSupplier:  nil,
			getAllSuppliers: nil,
		},
		{
			name: "should pass a valid email",
			form: url.Values{
				"contact_email": {"test"},
			},
			status:          http.StatusBadRequest,
			body:            []string{"Ingrese un correo válido"},
			getOneSupplier:  getOne,
			updateSupplier:  nil,
			getAllSuppliers: nil,
		},
		{
			name:           "should edit supplier",
			form:           url.Values{},
			status:         http.StatusOK,
			body:           []string{},
			getOneSupplier: getOne,
			updateSupplier: db.EXPECT().UpdateSupplier(&types.Supplier{
				ID:           id,
				SupplierId:   "test",
				Name:         "test",
				ContactEmail: sql.NullString{Valid: true, String: ""},
				ContactName:  sql.NullString{Valid: true, String: ""},
				ContactPhone: sql.NullString{Valid: true, String: ""},
				CompanyId:    uuid.UUID{},
			}).Return(nil),
			getAllSuppliers: db.EXPECT().GetAllSuppliers(uuid.UUID{}, "").Return([]types.Supplier{
				{
					ID:           id,
					SupplierId:   "test",
					Name:         "test",
					ContactEmail: sql.NullString{Valid: true, String: ""},
					ContactName:  sql.NullString{Valid: true, String: ""},
					ContactPhone: sql.NullString{Valid: true, String: ""},
					CompanyId:    uuid.UUID{},
				},
			}, nil),
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			if tt.getOneSupplier != nil {
				tt.getOneSupplier.Times(1)
			}

			if tt.updateSupplier != nil {
				tt.updateSupplier.Times(1)
			}

			if tt.getAllSuppliers != nil {
				tt.getAllSuppliers.Times(1)
			}

			req, res := createRequest(
				token,
				http.MethodPut,
				fmt.Sprintf("/bca/partials/suppliers/edit/%s", id.String()),
				strings.NewReader(tt.form.Encode()),
			)
			s.Router.ServeHTTP(res, req)
			for _, b := range tt.body {
				assert.Contains(t, res.Body.String(), b)
			}
		})
	}
}
