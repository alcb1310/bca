package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"bca-go-final/internal/database"
)

func TestProjectsTable(t *testing.T) {
	db := database.ServiceMock{}
	_, router := NewServer(db)

	t.Run("Unsupported request methods", func(t *testing.T) {
		t.Run("PATCH", func(t *testing.T) {
			response := httptest.NewRecorder()
			request := &http.Request{
				Method: http.MethodPatch,
				URL: &url.URL{
					Path: "/bca/partials/project",
				},
			}
			router.ProjectsTable(response, request)
			got := response.Code
			want := http.StatusMethodNotAllowed
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})

		t.Run("DELETE", func(t *testing.T) {
			response := httptest.NewRecorder()
			request := &http.Request{
				Method: http.MethodDelete,
				URL: &url.URL{
					Path: "/bca/partials/project",
				},
			}
			router.ProjectsTable(response, request)
			got := response.Code
			want := http.StatusMethodNotAllowed
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})

		t.Run("PUT", func(t *testing.T) {
			response := httptest.NewRecorder()
			request := &http.Request{
				Method: http.MethodPut,
				URL: &url.URL{
					Path: "/bca/partials/project",
				},
			}
			router.ProjectsTable(response, request)
			got := response.Code
			want := http.StatusMethodNotAllowed
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}
		})
	})

	t.Run("Successful projects request", func(t *testing.T) {
		t.Run("GET", func(t *testing.T) {
			response := httptest.NewRecorder()
			request := &http.Request{
				Method: http.MethodGet,
				URL: &url.URL{
					Path: "/bca/partials/project",
				},
			}
			router.ProjectsTable(response, request)
			got := response.Code
			want := http.StatusOK
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}

			expected := "No existen proyectos"
			if !strings.Contains(response.Body.String(), expected) {
				t.Errorf("expected %s, got %s", expected, response.Body.String())
			}
		})

		t.Run("POST", func(t *testing.T) {
			t.Run("Valid data", func(t *testing.T) {
				response := httptest.NewRecorder()
				formData := url.Values{
					"name":       {"prueba"},
					"gross_area": {"100"},
					"net_area":   {"100"},
				}
				request := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Path: "/bca/partials/project",
					},
					Form: formData,
				}
				router.ProjectsTable(response, request)
				got := response.Code
				want := http.StatusOK
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}
			})

			t.Run("Duplicate data", func(t *testing.T) {
				response := httptest.NewRecorder()
				formData := url.Values{
					"name":       {"exists"},
					"gross_area": {"100"},
					"net_area":   {"100"},
				}
				request := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Path: "/bca/partials/project",
					},
					Form: formData,
				}
				router.ProjectsTable(response, request)
				got := response.Code
				want := http.StatusConflict
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}
			})
		})
	})

	t.Run("Invalid projects request", func(t *testing.T) {
		t.Run("Empty project name", func(t *testing.T) {
			response := httptest.NewRecorder()
			formData := url.Values{}
			request := &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Path: "/bca/partials/project",
				},
				Form: formData,
			}
			router.ProjectsTable(response, request)
			got := response.Code
			want := http.StatusBadRequest
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}

			expected := "El nombre del proyecto es requerido"
			if !strings.Contains(response.Body.String(), expected) {
				t.Errorf("expected %s, got %s", expected, response.Body.String())
			}
		})

		t.Run("gross area", func(t *testing.T) {
			t.Run("Invalid gross area", func(t *testing.T) {
				response := httptest.NewRecorder()
				formData := url.Values{
					"name":       {"prueba"},
					"gross_area": {"invalid"},
				}

				request := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Path: "/bca/partials/project",
					},
					Form: formData,
				}

				router.ProjectsTable(response, request)

				got := response.Code
				want := http.StatusBadRequest
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}

				expected := "área bruta debe ser un número válido"
				if !strings.Contains(response.Body.String(), expected) {
					t.Errorf("expected %s, got %s", expected, response.Body.String())
				}
			})

			t.Run("Negative gross area", func(t *testing.T) {
				response := httptest.NewRecorder()
				formData := url.Values{
					"name":       {"prueba"},
					"gross_area": {"-100"},
				}

				request := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Path: "/bca/partials/project",
					},
					Form: formData,
				}

				router.ProjectsTable(response, request)

				got := response.Code
				want := http.StatusBadRequest
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}

				expected := "área bruta debe ser un número positivo"
				if !strings.Contains(response.Body.String(), expected) {
					t.Errorf("expected %s, got %s", expected, response.Body.String())
				}

			})
		})

		t.Run("net area", func(t *testing.T) {
			t.Run("Invalid net area", func(t *testing.T) {
				response := httptest.NewRecorder()
				formData := url.Values{
					"name":     {"prueba"},
					"net_area": {"invalid"},
				}

				request := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Path: "/bca/partials/project",
					},
					Form: formData,
				}

				router.ProjectsTable(response, request)

				got := response.Code
				want := http.StatusBadRequest
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}

				expected := "área útil debe ser un número válido"
				if !strings.Contains(response.Body.String(), expected) {
					t.Errorf("expected %s, got %s", expected, response.Body.String())
				}
			})

			t.Run("Negative net area", func(t *testing.T) {
				response := httptest.NewRecorder()
				formData := url.Values{
					"name":     {"prueba"},
					"net_area": {"-100"},
				}

				request := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Path: "/bca/partials/project",
					},
					Form: formData,
				}

				router.ProjectsTable(response, request)

				got := response.Code
				want := http.StatusBadRequest
				if got != want {
					t.Errorf("got %d, want %d", got, want)
				}

				expected := "área útil debe ser un número positivo"
				if !strings.Contains(response.Body.String(), expected) {
					t.Errorf("expected %s, got %s", expected, response.Body.String())
				}
			})
		})
	})
}

func TestProjectAdd(t *testing.T) {
	db := database.ServiceMock{}
	_, router := NewServer(db)

	response := httptest.NewRecorder()
	request := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: "/bca/partials/projects/add",
		},
	}

	router.ProjectAdd(response, request)

	got := response.Code
	want := http.StatusOK
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	expected := "Agregar Proyecto"
	if !strings.Contains(response.Body.String(), expected) {
		t.Errorf("expected %s, got %s", expected, response.Body.String())
	}
}
