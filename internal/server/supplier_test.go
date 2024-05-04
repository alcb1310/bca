package server_test

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
)

func TestSuppliersTable(t *testing.T) {
	srv, db := server.MakeServer()

	t.Run("method not allowed", func(t *testing.T) {
		request, response := server.MakeRequest(http.MethodPut, "/bca/partials/suppliers", nil)

		srv.SuppliersTable(response, request)

		assert.Equal(t, http.StatusMethodNotAllowed, response.Code)
	})

	t.Run("method GET", func(t *testing.T) {
		db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
		request, response := server.MakeRequest(http.MethodGet, "/bca/partials/suppliers", nil)

		srv.SuppliersTable(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("method POST", func(t *testing.T) {
		supplier := types.Supplier{
			SupplierId:   "1",
			Name:         "1",
			ContactEmail: sql.NullString{Valid: true, String: "1"},
			ContactName:  sql.NullString{Valid: true, String: "1"},
			ContactPhone: sql.NullString{Valid: true, String: "1"},
			CompanyId:    uuid.UUID{},
		}

		t.Run("validate data", func(t *testing.T) {
			t.Run("supplier_id", func(t *testing.T) {
				form := url.Values{}
				form.Add("supplier_id", "")
				buf := strings.NewReader(form.Encode())

				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/suppliers", buf)

				srv.SuppliersTable(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un valor para el RUC")
			})

			t.Run("name", func(t *testing.T) {
				form := url.Values{}
				form.Add("supplier_id", supplier.SupplierId)
				form.Add("name", "")
				buf := strings.NewReader(form.Encode())

				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/suppliers", buf)

				srv.SuppliersTable(response, request)

				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un valor para el nombre")
			})
		})

		t.Run("valid data", func(t *testing.T) {
			form := url.Values{}
			form.Add("supplier_id", supplier.SupplierId)
			form.Add("name", supplier.Name)
			form.Add("contact_email", supplier.ContactEmail.String)
			form.Add("contact_name", supplier.ContactName.String)
			form.Add("contact_phone", supplier.ContactPhone.String)
			buf := strings.NewReader(form.Encode())

			t.Run("create supplier", func(t *testing.T) {
				db.On("CreateSupplier", &supplier).Return(nil)
				db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
				request, response := server.MakeRequest(http.MethodPost, "/bca/partials/suppliers", buf)

				srv.SuppliersTable(response, request)

				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("create fail", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					srv, db := server.MakeServer()
					form := url.Values{}
					form.Add("supplier_id", supplier.SupplierId)
					form.Add("name", supplier.Name)
					form.Add("contact_email", supplier.ContactEmail.String)
					form.Add("contact_name", supplier.ContactName.String)
					form.Add("contact_phone", supplier.ContactPhone.String)
					buf := strings.NewReader(form.Encode())

					db.On("CreateSupplier", &supplier).Return(errors.New("duplicate"))

					request, response := server.MakeRequest(http.MethodPost, "/bca/partials/suppliers", buf)

					srv.SuppliersTable(response, request)

					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Equal(t, response.Body.String(), fmt.Sprintf("Proveedor con ruc %s y/o nombre %s ya existe", supplier.SupplierId, supplier.Name))
				})

				t.Run("unknown error", func(t *testing.T) {
					srv, db := server.MakeServer()
					form := url.Values{}
					form.Add("supplier_id", supplier.SupplierId)
					form.Add("name", supplier.Name)
					form.Add("contact_email", supplier.ContactEmail.String)
					form.Add("contact_name", supplier.ContactName.String)
					form.Add("contact_phone", supplier.ContactPhone.String)
					buf := strings.NewReader(form.Encode())

					db.On("CreateSupplier", &supplier).Return(UnknownError)

					request, response := server.MakeRequest(http.MethodPost, "/bca/partials/suppliers", buf)

					srv.SuppliersTable(response, request)

					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Contains(t, response.Body.String(), UnknownError.Error())
				})
			})
		})
	})
}

func TestSupplierAdd(t *testing.T) {
	srv, _ := server.MakeServer()

	request, response := server.MakeRequest(http.MethodGet, "/bca/partials/suppliers/add", nil)

	srv.SupplierAdd(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Agregar Proveedor")
}

func TestSuppliersEdit(t *testing.T) {
	supplierId := uuid.New()
	testURL := fmt.Sprintf("/bca/partials/suppliers/%s", supplierId.String())
	muxVars := make(map[string]string)
	muxVars["id"] = supplierId.String()

	srv, db := server.MakeServer()

	db.On("GetOneSupplier", supplierId, uuid.UUID{}).Return(types.Supplier{}, nil)

	request, response := server.MakeRequest(http.MethodGet, testURL, nil)
	request = mux.SetURLVars(request, muxVars)

	srv.SuppliersEdit(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Editar Proveedor")
}
