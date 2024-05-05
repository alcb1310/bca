package server_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"bca-go-final/internal/server"
	"bca-go-final/internal/types"
)

func TestInvoicesTable(t *testing.T) {
	srv, db := server.MakeServer()
	db.On("GetInvoices", uuid.UUID{}).Return([]types.InvoiceResponse{}, nil)

	request, response := server.MakeRequest(http.MethodGet, "/bca/partials/invoices", nil)
	srv.InvoicesTable(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestInvoiceAdd(t *testing.T) {
	testURL := "/bca/transacciones/facturas/crear"
	t.Run("method GET", func(t *testing.T) {
		t.Run("no query parameters", func(t *testing.T) {
			srv, db := server.MakeServer()
			db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
			db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)

			request, response := server.MakeRequest(http.MethodGet, testURL, nil)
			srv.InvoiceAdd(response, request)
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Contains(t, response.Body.String(), "Nueva Factura")
		})

		t.Run("with query parameters", func(t *testing.T) {
			invoiceId := uuid.New()
			newTestURL := fmt.Sprintf("%s?id=%s", testURL, invoiceId.String())
			srv, db := server.MakeServer()
			db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
			db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
			db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)

			request, response := server.MakeRequest(http.MethodGet, newTestURL, nil)
			srv.InvoiceAdd(response, request)
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Contains(t, response.Body.String(), "Editar Factura")
		})
	})

	t.Run("method POST", func(t *testing.T) {
		projectId := uuid.New()
		supplierId := uuid.New()
		invoiceNumber := "S/N"
		invoiceDate := "0001-01-01"
		invoiceId := uuid.MustParse("cdefa321-9f2d-4673-9949-7cac744e941a")

		t.Run("Data validation", func(t *testing.T) {
			t.Run("project", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", "")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.InvoiceAdd(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "Seleccione un proyecto")
				})

				t.Run("invalid", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", "invalid")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.InvoiceAdd(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "invalid UUID length: 7")
				})
			})

			t.Run("supplier", func(t *testing.T) {
				t.Run("empty", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("supplier", "")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.InvoiceAdd(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "Seleccione un proveedor")
				})

				t.Run("invalid", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("supplier", "invalid")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.InvoiceAdd(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Contains(t, response.Body.String(), "invalid UUID length: 7")
				})
			})

			t.Run("invoice number", func(t *testing.T) {
				form := url.Values{}
				form.Add("project", projectId.String())
				form.Add("supplier", supplierId.String())
				form.Add("invoiceNumber", "")
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
				db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)

				request, response := server.MakeRequest(http.MethodPost, testURL, buf)
				srv.InvoiceAdd(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese un número de factura")
			})

			t.Run("invoice date", func(t *testing.T) {
				form := url.Values{}
				form.Add("project", projectId.String())
				form.Add("supplier", supplierId.String())
				form.Add("invoiceNumber", invoiceNumber)
				form.Add("invoiceDate", "")
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
				db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)

				request, response := server.MakeRequest(http.MethodPost, testURL, buf)
				srv.InvoiceAdd(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Contains(t, response.Body.String(), "Ingrese una fecha válida")
			})
		})

		t.Run("Valid data", func(t *testing.T) {
			t.Run("Create invoice", func(t *testing.T) {
				form := url.Values{}
				form.Add("project", projectId.String())
				form.Add("supplier", supplierId.String())
				form.Add("invoiceNumber", invoiceNumber)
				form.Add("invoiceDate", invoiceDate)
				buf := strings.NewReader(form.Encode())
				tm, _ := time.Parse("2020-01-02", invoiceDate)

				createInvoice := &types.InvoiceCreate{
					Id:            nil,
					ProjectId:     &projectId,
					SupplierId:    &supplierId,
					InvoiceNumber: &invoiceNumber,
					InvoiceDate:   &tm,
					IsBalanced:    false,
					CompanyId:     uuid.UUID{},
				}
				_ = createInvoice

				srv, db := server.MakeServer()
				db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
				db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
				db.On("CreateInvoice", createInvoice).Return(nil)
				db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{
					Id: invoiceId,
				}, nil)

				request, response := server.MakeRequest(http.MethodPost, testURL, buf)
				srv.InvoiceAdd(response, request)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("Falied to create invoice", func(t *testing.T) {
				t.Run("Conflict", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("supplier", supplierId.String())
					form.Add("invoiceNumber", invoiceNumber)
					form.Add("invoiceDate", invoiceDate)
					buf := strings.NewReader(form.Encode())
					tm, _ := time.Parse("2020-01-02", invoiceDate)

					createInvoice := &types.InvoiceCreate{
						Id:            nil,
						ProjectId:     &projectId,
						SupplierId:    &supplierId,
						InvoiceNumber: &invoiceNumber,
						InvoiceDate:   &tm,
						IsBalanced:    false,
						CompanyId:     uuid.UUID{},
					}
					_ = createInvoice

					srv, db := server.MakeServer()
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
					db.On("CreateInvoice", createInvoice).Return(errors.New("duplicate"))

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.InvoiceAdd(response, request)
					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Equal(t, "La Factura ya existe", response.Body.String())
				})

				t.Run("Internal server error", func(t *testing.T) {
					form := url.Values{}
					form.Add("project", projectId.String())
					form.Add("supplier", supplierId.String())
					form.Add("invoiceNumber", invoiceNumber)
					form.Add("invoiceDate", invoiceDate)
					buf := strings.NewReader(form.Encode())
					tm, _ := time.Parse("2020-01-02", invoiceDate)

					createInvoice := &types.InvoiceCreate{
						Id:            nil,
						ProjectId:     &projectId,
						SupplierId:    &supplierId,
						InvoiceNumber: &invoiceNumber,
						InvoiceDate:   &tm,
						IsBalanced:    false,
						CompanyId:     uuid.UUID{},
					}
					_ = createInvoice

					srv, db := server.MakeServer()
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
					db.On("CreateInvoice", createInvoice).Return(UnknownError)

					request, response := server.MakeRequest(http.MethodPost, testURL, buf)
					srv.InvoiceAdd(response, request)
					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Equal(t, UnknownError.Error(), response.Body.String())

				})
			})
		})
	})
}

func TestInvoiceEdit(t *testing.T) {
	invoiceId := uuid.New()
	testURL := fmt.Sprintf("/bca/partials/invoices/%s", invoiceId.String())
	muxVars := make(map[string]string)
	muxVars["id"] = invoiceId.String()

	supplierResponse := types.InvoiceResponse{
		Id:            uuid.UUID{},
		CompanyId:     uuid.UUID{},
		Project:       types.Project{},
		Supplier:      types.Supplier{},
		InvoiceNumber: "",
	}

	t.Run("method GET", func(t *testing.T) {
		srv, db := server.MakeServer()
		db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
		db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
		db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)

		request, response := server.MakeRequest(http.MethodGet, testURL, nil)
		request = mux.SetURLVars(request, muxVars)
		srv.InvoiceEdit(response, request)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), "Editar Factura")
	})

	t.Run("method DELETE", func(t *testing.T) {
		t.Run("successful delete", func(t *testing.T) {
			srv, db := server.MakeServer()
			db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
			db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
			db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)
			db.On("DeleteInvoice", invoiceId, uuid.UUID{}).Return(nil)
			db.On("GetInvoices", uuid.UUID{}).Return([]types.InvoiceResponse{}, nil)

			request, response := server.MakeRequest(http.MethodDelete, testURL, nil)
			request = mux.SetURLVars(request, muxVars)
			srv.InvoiceEdit(response, request)
			assert.Equal(t, http.StatusOK, response.Code)
		})

		t.Run("failed delete", func(t *testing.T) {
			srv, db := server.MakeServer()
			db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
			db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
			db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)
			db.On("DeleteInvoice", invoiceId, uuid.UUID{}).Return(UnknownError)
			db.On("GetInvoices", uuid.UUID{}).Return([]types.InvoiceResponse{}, nil)

			request, response := server.MakeRequest(http.MethodDelete, testURL, nil)
			request = mux.SetURLVars(request, muxVars)
			srv.InvoiceEdit(response, request)
			assert.Equal(t, http.StatusInternalServerError, response.Code)
			assert.Equal(t, UnknownError.Error(), response.Body.String())
		})
	})

	t.Run("method PUT", func(t *testing.T) {
		t.Run("data validation", func(t *testing.T) {
			t.Run("supplier", func(t *testing.T) {
				t.Run("empty supplier", func(t *testing.T) {
					form := url.Values{}
					form.Add("supplier", "")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.InvoiceEdit(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "Seleccione un proveedor")
				})

				t.Run("invalid supplier", func(t *testing.T) {
					form := url.Values{}
					form.Add("supplier", "invalid")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.InvoiceEdit(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "invalid UUID length: 7")
				})
			})

			t.Run("invoice number", func(t *testing.T) {
				form := url.Values{}
				form.Add("supplier", uuid.UUID{}.String())
				form.Add("invoiceNumber", "")
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
				db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
				db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)

				request, response := server.MakeRequest(http.MethodPut, testURL, buf)
				request = mux.SetURLVars(request, muxVars)
				srv.InvoiceEdit(response, request)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.Equal(t, response.Body.String(), "Ingrese un número de factura")
			})

			t.Run("invoice date", func(t *testing.T) {
				t.Run("empty date", func(t *testing.T) {
					form := url.Values{}
					form.Add("supplier", uuid.UUID{}.String())
					form.Add("invoiceNumber", "S/N")
					form.Add("invoiceDate", "")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.InvoiceEdit(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "Ingrese una fecha válida")
				})

				t.Run("invalid date", func(t *testing.T) {
					form := url.Values{}
					form.Add("supplier", uuid.UUID{}.String())
					form.Add("invoiceNumber", "S/N")
					form.Add("invoiceDate", "invalid")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.InvoiceEdit(response, request)
					assert.Equal(t, http.StatusBadRequest, response.Code)
					assert.Equal(t, response.Body.String(), "Ingrese una fecha válida")
				})
			})
		})

		t.Run("valid data", func(t *testing.T) {
			invoiceNumber := "S/N"
			date, _ := time.Parse("2006-01-02", "2020-10-10")
			invoiceCreate := types.InvoiceCreate{
				Id:            &invoiceId,
				ProjectId:     &uuid.UUID{},
				SupplierId:    &uuid.UUID{},
				InvoiceNumber: &invoiceNumber,
				InvoiceDate:   &date,
			}

			t.Run("successful update", func(t *testing.T) {
				form := url.Values{}
				form.Add("supplier", uuid.UUID{}.String())
				form.Add("invoiceNumber", "S/N")
				form.Add("invoiceDate", "2020-10-10")
				buf := strings.NewReader(form.Encode())

				srv, db := server.MakeServer()
				db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
				db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
				db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)
				db.On("UpdateInvoice", invoiceCreate).Return(nil)

				request, response := server.MakeRequest(http.MethodPut, testURL, buf)
				request = mux.SetURLVars(request, muxVars)
				srv.InvoiceEdit(response, request)
				assert.Equal(t, http.StatusOK, response.Code)
			})

			t.Run("failed update", func(t *testing.T) {
				t.Run("duplicate", func(t *testing.T) {
					form := url.Values{}
					form.Add("supplier", uuid.UUID{}.String())
					form.Add("invoiceNumber", "S/N")
					form.Add("invoiceDate", "2020-10-10")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)
					db.On("UpdateInvoice", invoiceCreate).Return(errors.New("duplicate"))

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.InvoiceEdit(response, request)
					assert.Equal(t, http.StatusConflict, response.Code)
					assert.Equal(t, response.Body.String(), "La Factura ya existe")
				})

				t.Run("unknown error", func(t *testing.T) {
					form := url.Values{}
					form.Add("supplier", uuid.UUID{}.String())
					form.Add("invoiceNumber", "S/N")
					form.Add("invoiceDate", "2020-10-10")
					buf := strings.NewReader(form.Encode())

					srv, db := server.MakeServer()
					db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
					db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
					db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)
					db.On("UpdateInvoice", invoiceCreate).Return(UnknownError)

					request, response := server.MakeRequest(http.MethodPut, testURL, buf)
					request = mux.SetURLVars(request, muxVars)
					srv.InvoiceEdit(response, request)
					assert.Equal(t, http.StatusInternalServerError, response.Code)
					assert.Equal(t, response.Body.String(), UnknownError.Error())
				})
			})
		})
	})

	t.Run("method PATCH", func(t *testing.T) {
		t.Run("successful update", func(t *testing.T) {
			srv, db := server.MakeServer()
			db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
			db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
			db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)
			db.On("BalanceInvoice", supplierResponse).Return(nil)

			request, response := server.MakeRequest(http.MethodPatch, testURL, nil)
			request = mux.SetURLVars(request, muxVars)
			srv.InvoiceEdit(response, request)
			assert.Equal(t, http.StatusOK, response.Code)
		})

		t.Run("failed update", func(t *testing.T) {
			srv, db := server.MakeServer()
			db.On("GetAllSuppliers", uuid.UUID{}, "").Return([]types.Supplier{}, nil)
			db.On("GetActiveProjects", uuid.UUID{}, true).Return([]types.Project{}, nil)
			db.On("GetOneInvoice", invoiceId, uuid.UUID{}).Return(types.InvoiceResponse{}, nil)
			db.On("BalanceInvoice", supplierResponse).Return(UnknownError)

			request, response := server.MakeRequest(http.MethodPatch, testURL, nil)
			request = mux.SetURLVars(request, muxVars)
			srv.InvoiceEdit(response, request)
			assert.Equal(t, http.StatusInternalServerError, response.Code)
			assert.Equal(t, UnknownError.Error(), response.Body.String())
		})
	})
}
