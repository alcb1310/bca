package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"

	"bca-go-final/internal/database"
)

func TestSuppliersTable(t *testing.T) {
	db := database.ServiceMock{}
	_, router := NewServer(db)

	t.Run("Should return not allowed", func(t *testing.T) {
		t.Run("PUT method", func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPut, "/bca/partials/suppliers", nil)

			router.SuppliersTable(response, request)

			got := response.Code
			want := http.StatusMethodNotAllowed

			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
	})

	t.Run("Successful requests", func(t *testing.T) {
		t.Run("Should return 200 on GET", func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/bca/partials/suppliers", nil)

			router.SuppliersTable(response, request)

			got := response.Code
			want := http.StatusOK

			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}

			expected := "No hay Proveedores"
			if !strings.Contains(response.Body.String(), expected) {
				t.Errorf("expected %s, got %s", expected, response.Body.String())
			}
		})

		t.Run("Should return 200 on POST", func(t *testing.T) {
			response := httptest.NewRecorder()
			form := url.Values{
				"supplier_id":   {"1234567890"},
				"name":          {"prueba"},
				"contact_name":  {"prueba"},
				"contact_phone": {"1234567890"},
				"contact_email": {"prueba@example.com"},
			}
			request := &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Path: "/bca/partials/project",
				},
				Form: form,
			}

			router.SuppliersTable(response, request)

			got := response.Code
			want := http.StatusOK
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
	})

	t.Run("Invalid data", func(t *testing.T) {
		t.Run("Empty RUC", func(t *testing.T) {
			response := httptest.NewRecorder()
			form := url.Values{}

			request := &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Path: "/bca/partials/suppliers",
				},
				Form: form,
			}

			router.SuppliersTable(response, request)

			got := response.Code
			want := http.StatusBadRequest
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}

			expected := "Ingrese un valor para el RUC"
			if !strings.Contains(response.Body.String(), expected) {
				t.Errorf("expected %s, got %s", expected, response.Body.String())
			}
		})

		t.Run("Empty Name", func(t *testing.T) {
			response := httptest.NewRecorder()
			form := url.Values{
				"supplier_id": {"1234567890"},
			}

			request := &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Path: "/bca/partials/suppliers",
				},
				Form: form,
			}

			router.SuppliersTable(response, request)

			got := response.Code
			want := http.StatusBadRequest
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}

			expected := "Ingrese un valor para el nombre"
			if !strings.Contains(response.Body.String(), expected) {
				t.Errorf("expected %s, got %s", expected, response.Body.String())
			}
		})

		t.Run("Conflict", func(t *testing.T) {
			t.Run("Duplicate RUC", func(t *testing.T) {
				response := httptest.NewRecorder()
				form := url.Values{
					"supplier_id":   {"0123456789"},
					"name":          {"prueba"},
					"contact_name":  {"prueba"},
					"contact_phone": {"1234567890"},
					"contact_email": {"prueba@example.com"},
				}

				request := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Path: "/bca/partials/suppliers",
					},
					Form: form,
				}

				router.SuppliersTable(response, request)

				got := response.Code
				want := http.StatusConflict
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}

				expected := "Proveedor con ruc 0123456789 y/o nombre prueba ya existe"
				if !strings.Contains(response.Body.String(), expected) {
					t.Errorf("expected %s, got %s", expected, response.Body.String())
				}
			})

			t.Run("Duplicate Name", func(t *testing.T) {
				response := httptest.NewRecorder()
				form := url.Values{
					"supplier_id":   {"1234567890"},
					"name":          {"exists"},
					"contact_name":  {"prueba"},
					"contact_phone": {"1234567890"},
					"contact_email": {"prueba@example.com"},
				}

				request := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Path: "/bca/partials/suppliers",
					},
					Form: form,
				}

				router.SuppliersTable(response, request)

				got := response.Code
				want := http.StatusConflict
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}

				expected := "Proveedor con ruc 1234567890 y/o nombre exists ya existe"
				if !strings.Contains(response.Body.String(), expected) {
					t.Errorf("expected %s, got %s", expected, response.Body.String())
				}
			})
		})
	})
}

func TestSupplierAdd(t *testing.T) {
	db := database.ServiceMock{}
	_, router := NewServer(db)

	response := httptest.NewRecorder()
	request := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: "/bca/partials/suppliers/add",
		},
	}

	router.SupplierAdd(response, request)

	got := response.Code
	want := http.StatusOK
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	expected := "Agregar Proveedor"
	if !strings.Contains(response.Body.String(), expected) {
		t.Errorf("expected %s, got %s", expected, response.Body.String())
	}
}

func TestSupplierEdit(t *testing.T) {
	db := database.ServiceMock{}
	_, router := NewServer(db)
	supplierId := uuid.New().String()

	response := httptest.NewRecorder()
	request := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: fmt.Sprintf("/bca/partials/suppliers/%s", supplierId),
		},
	}

	router.SuppliersEdit(response, request)

	got := response.Code
	want := http.StatusOK
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	expected := "Editar Proveedor"
	if !strings.Contains(response.Body.String(), expected) {
		t.Errorf("expected %s, got %s", expected, response.Body.String())
	}
}

func TestSupplierEditSave(t *testing.T) {
	db := database.ServiceMock{}
	_, router := NewServer(db)

	t.Run("Successful request", func(t *testing.T) {
		response := httptest.NewRecorder()
		supplierId := uuid.New().String()

		form := url.Values{
			"supplier_id":   {"1234567890"},
			"name":          {"prueba"},
			"contact_name":  {"prueba"},
			"contact_phone": {"1234567890"},
			"contact_email": {"prueba@example.com"},
		}

		request := &http.Request{
			Method: http.MethodPut,
			URL: &url.URL{
				Path: fmt.Sprintf("/bca/partials/suppliers/edit/%s", supplierId),
			},
			Form: form,
		}

		router.SuppliersEditSave(response, request)

		got := response.Code
		want := http.StatusOK
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("Invalid data", func(t *testing.T) {
		t.Run("Empty RUC", func(t *testing.T) {
			response := httptest.NewRecorder()
			supplierId := uuid.New().String()

			form := url.Values{
				"name":          {"prueba"},
				"contact_name":  {"prueba"},
				"contact_phone": {"1234567890"},
			}

			request := &http.Request{
				Method: http.MethodPut,
				URL: &url.URL{
					Path: fmt.Sprintf("/bca/partials/suppliers/edit/%s", supplierId),
				},
				Form: form,
			}

			router.SuppliersEditSave(response, request)
			got := response.Code
			want := http.StatusBadRequest
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}

			expected := "Ingrese un valor para el RUC"
			responseBody := response.Body.String()
			if !strings.Contains(responseBody, expected) {
				t.Errorf("expected %s, got %s", expected, responseBody)
			}
		})

		t.Run("Empty Name", func(t *testing.T) {
			response := httptest.NewRecorder()
			supplierId := uuid.New().String()

			form := url.Values{
				"supplier_id":   {"0123456789"},
				"contact_name":  {"prueba"},
				"contact_phone": {"1234567890"},
			}

			request := &http.Request{
				Method: http.MethodPut,
				URL: &url.URL{
					Path: fmt.Sprintf("/bca/partials/suppliers/edit/%s", supplierId),
				},
				Form: form,
			}

			router.SuppliersEditSave(response, request)

			got := response.Code
			want := http.StatusBadRequest
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}

			expected := "Ingrese un valor para el nombre"
			responseBody := response.Body.String()
			if !strings.Contains(responseBody, expected) {
				t.Errorf("expected %s, got %s", expected, responseBody)
			}
		})

		t.Run("Conflict", func(t *testing.T) {
			t.Run("Existing supplier_id", func(t *testing.T) {
				response := httptest.NewRecorder()
				supplierId := uuid.New().String()

				form := url.Values{
					"supplier_id":   {"0123456789"},
					"name":          {"prueba"},
					"contact_name":  {"prueba"},
					"contact_phone": {"1234567890"},
					"contact_email": {"prueba@example.com"},
				}

				request := &http.Request{
					Method: http.MethodPut,
					URL: &url.URL{
						Path: fmt.Sprintf("/bca/partials/suppliers/edit/%s", supplierId),
					},
					Form: form,
				}

				router.SuppliersEditSave(response, request)

				got := response.Code
				want := http.StatusConflict
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}

				expected := "El ruc 0123456789 y/o nombre prueba ya existe"
				received := response.Body.String()

				if expected != received {
					t.Errorf("got %s, want %s", received, expected)
				}
			})

			t.Run("Existing supplier_id", func(t *testing.T) {
				response := httptest.NewRecorder()
				supplierId := uuid.New().String()

				form := url.Values{
					"supplier_id":   {"023456789"},
					"name":          {"exists"},
					"contact_name":  {"prueba"},
					"contact_phone": {"1234567890"},
					"contact_email": {"prueba@example.com"},
				}

				request := &http.Request{
					Method: http.MethodPut,
					URL: &url.URL{
						Path: fmt.Sprintf("/bca/partials/suppliers/edit/%s", supplierId),
					},
					Form: form,
				}

				router.SuppliersEditSave(response, request)

				got := response.Code
				want := http.StatusConflict
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}

				expected := "El ruc 023456789 y/o nombre exists ya existe"
				received := response.Body.String()

				if expected != received {
					t.Errorf("got %s, want %s", received, expected)
				}
			})
		})
	})
}
