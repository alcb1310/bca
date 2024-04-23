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
	"bca-go-final/internal/types"
)

var companyId = uuid.New()

var projectDuplicate = "bc39e850-0a1f-446f-a112-3e9a5b3134f0"
var budgetItemDuplicate = "cc5cbcb9-43cc-4062-b3d3-ea60a3c2e6d0"

var oldBudget = types.GetBudget{
	CompanyId:      companyId,
	InitialTotal:   1,
	SpentTotal:     0,
	RemainingTotal: 1,
	UpdatedBudget:  1,
}

func TestBudgetsTable(t *testing.T) {
	db := database.ServiceMock{}
	_, router := NewServer(db)

	t.Run("Successful budgets request", func(t *testing.T) {
		t.Run("Should return 200 when GET budgets", func(t *testing.T) {
			response := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, "/bca/partials/budget", nil)
			router.BudgetsTable(response, request)

			got := response.Code
			want := http.StatusOK
			if got != want {
				t.Errorf("got %d, want %d", got, want)
			}

			expected := `<table><thead><tr><th width="180px" rowspan="2">Proyecto</th><th width="380px" rowspan="2">Partida</th><th colspan="3">Por Gastar</th><th width="130px" rowspan="2">Actualizado</th><th width="30px" rowspan="2"></th></tr><tr><th width="130px">Cantidad</th><th width="130px">Unitario</th><th width="130px">Total</th></tr></thead> <tbody><tr><td colspan="8">No existen presupuestos</td></tr></tbody></table>`
			if response.Body.String() != expected {
				t.Errorf("got %s, want %s", response.Body.String(), expected)
			}
		})

		t.Run("Valid POST request", func(t *testing.T) {
			t.Run("Should return 200 when POST budgets", func(t *testing.T) {
				response := httptest.NewRecorder()

				form := url.Values{}
				form.Add("project", uuid.New().String())
				form.Add("budgetItem", uuid.New().String())
				form.Add("quantity", "1")
				form.Add("cost", "1")

				request := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Path: "/bca/partials/budget",
					},
					Form: form,
				}
				router.BudgetsTable(response, request)

				got := response.Code
				want := http.StatusOK
				if got != want {
					t.Errorf("got %d, want %d", got, want)
					t.Error(response.Body.String())
				}
			})

			t.Run("Should return 409 when duplicate budget", func(t *testing.T) {
				response := httptest.NewRecorder()

				form := url.Values{}
				form.Add("project", projectDuplicate)
				form.Add("budgetItem", budgetItemDuplicate)
				form.Add("quantity", "1")
				form.Add("cost", "1")

				request := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Path: "/bca/partials/budget",
					},
					Form: form,
				}
				router.BudgetsTable(response, request)

				got := response.Code
				want := http.StatusConflict
				if got != want {
					t.Errorf("got %d, want %d", got, want)
					t.Error(response.Body.String())
				}

				expected := fmt.Sprintf("Ya existe partida %s en el proyecto %s", budgetItemDuplicate, projectDuplicate)
				received := strings.Trim(response.Body.String(), "\n")
				if received != expected {
					t.Errorf("got %s, want %s", received, expected)
				}
			})
		})

		t.Run("Unsuccessful budgets request", func(t *testing.T) {
			t.Run("POST request", func(t *testing.T) {
				t.Run("Validate project id", func(t *testing.T) {
					response := httptest.NewRecorder()
					form := url.Values{}
					form.Add("budgetItem", uuid.New().String())
					form.Add("quantity", "1")
					form.Add("cost", "1")

					t.Run("no project id", func(t *testing.T) {
						request := &http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Path: "/bca/partials/budget",
							},
							Form: form,
						}
						router.BudgetsTable(response, request)

						got := response.Code
						want := http.StatusBadRequest
						if got != want {
							t.Errorf("got %d, want %d", got, want)
							t.Error(response.Body.String())
						}

						expected := "Seleccione un proyecto"
						received := strings.Trim(response.Body.String(), "\n")
						if received != expected {
							t.Errorf("got %s, want %s", received, expected)
						}
					})

					t.Run("invalid project id", func(t *testing.T) {
						form.Del("project")
						form.Add("project", "invalid")
						request := &http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Path: "/bca/partials/budget",
							},
							Form: form,
						}
						router.BudgetsTable(response, request)

						got := response.Code
						want := http.StatusBadRequest
						if got != want {
							t.Errorf("got %d, want %d", got, want)
							t.Error(response.Body.String())
						}

						expected := "invalid UUID"
						received := strings.Trim(response.Body.String(), "\n")
						if !strings.Contains(received, expected) {
							t.Errorf("got %s, want %s", received, expected)
						}
					})
				})

				t.Run("Validate budget item id", func(t *testing.T) {
					response := httptest.NewRecorder()
					form := url.Values{}
					form.Add("project", uuid.New().String())
					form.Add("quantity", "1")
					form.Add("cost", "1")

					t.Run("no budget item id", func(t *testing.T) {
						request := &http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Path: "/bca/partials/budget",
							},
							Form: form,
						}
						router.BudgetsTable(response, request)

						got := response.Code
						want := http.StatusBadRequest
						if got != want {
							t.Errorf("got %d, want %d", got, want)
							t.Error(response.Body.String())
						}

						expected := "Seleccione un partida"
						received := strings.Trim(response.Body.String(), "\n")
						if received != expected {
							t.Errorf("got %s, want %s", received, expected)
						}
					})

					t.Run("invalid budget item id", func(t *testing.T) {
						form.Del("budgetItem")
						form.Add("budgetItem", "invalid")
						request := &http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Path: "/bca/partials/budget",
							},
							Form: form,
						}
						router.BudgetsTable(response, request)

						got := response.Code
						want := http.StatusBadRequest
						if got != want {
							t.Errorf("got %d, want %d", got, want)
							t.Error(response.Body.String())
						}

						expected := "invalid UUID"
						received := strings.Trim(response.Body.String(), "\n")
						if !strings.Contains(received, expected) {
							t.Errorf("got %s, want %s", received, expected)
						}
					})
				})

				t.Run("Validate quantity", func(t *testing.T) {
					response := httptest.NewRecorder()
					form := url.Values{}
					form.Add("budgetItem", uuid.New().String())
					form.Add("project", uuid.New().String())
					form.Add("cost", "1")

					t.Run("no quantity", func(t *testing.T) {
						request := &http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Path: "/bca/partials/budget",
							},
							Form: form,
						}
						router.BudgetsTable(response, request)

						got := response.Code
						want := http.StatusBadRequest
						if got != want {
							t.Errorf("got %d, want %d", got, want)
							t.Error(response.Body.String())
						}

						expected := "cantidad es requerido"
						received := strings.Trim(response.Body.String(), "\n")
						if !strings.Contains(received, expected) {
							t.Errorf("got %s, want %s", received, expected)
						}
					})

					t.Run("invalid quantity", func(t *testing.T) {
						form.Add("quantity", "invalid")
						request := &http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Path: "/bca/partials/budget",
							},
							Form: form,
						}
						router.BudgetsTable(response, request)

						got := response.Code
						want := http.StatusBadRequest
						if got != want {
							t.Errorf("got %d, want %d", got, want)
							t.Error(response.Body.String())
						}

						expected := "cantidad debe ser un número válido"
						received := strings.Trim(response.Body.String(), "\n")
						if !strings.Contains(received, expected) {
							t.Errorf("got %s, want %s", received, expected)
						}
					})

					t.Run("negative quantity", func(t *testing.T) {
						form.Del("quantity")
						form.Add("quantity", "-1")
						request := &http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Path: "/bca/partials/budget",
							},
							Form: form,
						}
						router.BudgetsTable(response, request)

						got := response.Code
						want := http.StatusBadRequest
						if got != want {
							t.Errorf("got %d, want %d", got, want)
							t.Error(response.Body.String())
						}

						expected := "cantidad debe ser un número positivo"
						received := strings.Trim(response.Body.String(), "\n")
						if !strings.Contains(received, expected) {
							t.Errorf("got %s, want %s", received, expected)
						}
					})
				})

				t.Run("Validate cost", func(t *testing.T) {
					response := httptest.NewRecorder()
					form := url.Values{}
					form.Add("budgetItem", uuid.New().String())
					form.Add("project", uuid.New().String())
					form.Add("quantity", "1")

					t.Run("no cost", func(t *testing.T) {
						request := &http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Path: "/bca/partials/budget",
							},
							Form: form,
						}
						router.BudgetsTable(response, request)

						got := response.Code
						want := http.StatusBadRequest
						if got != want {
							t.Errorf("got %d, want %d", got, want)
							t.Error(response.Body.String())
						}

						expected := "costo es requerido"
						received := strings.Trim(response.Body.String(), "\n")
						if !strings.Contains(received, expected) {
							t.Errorf("got %s, want %s", received, expected)
						}
					})

					t.Run("invalid cost", func(t *testing.T) {
						form.Add("cost", "invalid")
						request := &http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Path: "/bca/partials/budget",
							},
							Form: form,
						}
						router.BudgetsTable(response, request)

						got := response.Code
						want := http.StatusBadRequest
						if got != want {
							t.Errorf("got %d, want %d", got, want)
							t.Error(response.Body.String())
						}

						expected := "costo debe ser un número válido"
						received := strings.Trim(response.Body.String(), "\n")
						if !strings.Contains(received, expected) {
							t.Errorf("got %s, want %s", received, expected)
						}
					})

					t.Run("negative cost", func(t *testing.T) {
						form.Del("cost")
						form.Add("cost", "-1")
						request := &http.Request{
							Method: http.MethodPost,
							URL: &url.URL{
								Path: "/bca/partials/budget",
							},
							Form: form,
						}
						router.BudgetsTable(response, request)

						got := response.Code
						want := http.StatusBadRequest
						if got != want {
							t.Errorf("got %d, want %d", got, want)
							t.Error(response.Body.String())
						}

						expected := "costo debe ser un número positivo"
						received := strings.Trim(response.Body.String(), "\n")
						if !strings.Contains(received, expected) {
							t.Errorf("got %s, want %s", received, expected)
						}
					})
				})
			})
		})
	})
}

func TestUpdateBudget(t *testing.T) {
	db := database.ServiceMock{}
	_, srv := NewServer(db)

	t.Run("valid budget data", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			form := url.Values{}

			form.Add("quantity", "1")
			form.Add("cost", "1")

			err := updateBudget(form, uuid.New(), uuid.New(), companyId, &oldBudget, srv)

			if err != nil {
				t.Errorf("Expected no error and got %v", err)
			}
		})
	})

	t.Run("invalid budget data", func(t *testing.T) {
		t.Run("invalid quantity", func(t *testing.T) {
			t.Run("empty quantity", func(t *testing.T) {
				form := url.Values{}
				form.Add("cost", "1")

				err := updateBudget(form, uuid.New(), uuid.New(), companyId, &oldBudget, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "cantidad es requerido" {
					t.Errorf("Expected 'cantidad es requerido' and got '%s'", err.Error())
				}
			})

			t.Run("invalid quantity", func(t *testing.T) {
				form := url.Values{}
				form.Add("quantity", "invalid")
				form.Add("cost", "1")

				err := updateBudget(form, uuid.New(), uuid.New(), companyId, &oldBudget, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "cantidad debe ser un número válido" {
					t.Errorf("Expected 'cantidad debe ser un número válido' and got '%s'", err.Error())
				}
			})

			t.Run("quantity must be a positive number", func(t *testing.T) {
				form := url.Values{}
				form.Add("quantity", "-4")
				form.Add("cost", "1")

				err := updateBudget(form, uuid.New(), uuid.New(), companyId, &oldBudget, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "cantidad debe ser un número positivo" {
					t.Errorf("Expected 'cantidad debe ser un número positivo' and got '%s'", err.Error())
				}
			})
		})

		t.Run("invalid cost", func(t *testing.T) {
			t.Run("empty cost", func(t *testing.T) {
				form := url.Values{}
				form.Add("quantity", "1")

				err := updateBudget(form, uuid.New(), uuid.New(), companyId, &oldBudget, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "costo es requerido" {
					t.Errorf("Expected 'costo es requerido' and got '%s'", err.Error())
				}
			})

			t.Run("invalid cost", func(t *testing.T) {
				form := url.Values{}
				form.Add("cost", "invalid")
				form.Add("quantity", "1")

				err := updateBudget(form, uuid.New(), uuid.New(), companyId, &oldBudget, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "costo debe ser un número válido" {
					t.Errorf("Expected 'costo debe ser un número válido' and got '%s'", err.Error())
				}
			})

			t.Run("cost must be a positive number", func(t *testing.T) {
				form := url.Values{}
				form.Add("cost", "-4")
				form.Add("quantity", "1")

				err := updateBudget(form, uuid.New(), uuid.New(), companyId, &oldBudget, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "costo debe ser un número positivo" {
					t.Errorf("Expected 'costo debe ser un número positivo' and got '%s'", err.Error())
				}
			})
		})
	})
}
