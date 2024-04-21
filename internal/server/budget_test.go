package server

import (
	"net/url"
	"testing"

	"github.com/google/uuid"

	"bca-go-final/internal/database"
	"bca-go-final/internal/types"
)

var companyId = uuid.New()

var oldBudget = types.GetBudget{
	CompanyId:      companyId,
	InitialTotal:   1,
	SpentTotal:     0,
	RemainingTotal: 1,
	UpdatedBudget:  1,
}

func TestCreateBudget(t *testing.T) {
	db := database.ServiceMock{}
	_, srv := NewServer(db)

	t.Run("create budget", func(t *testing.T) {
		t.Run("succesfully create budget", func(t *testing.T) {
			form := url.Values{}

			form.Add("project", uuid.New().String())
			form.Add("budgetItem", uuid.New().String())
			form.Add("quantity", "1")
			form.Add("cost", "1")

			_, err := createBudget(form, companyId, srv)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("duplicate budget", func(t *testing.T) {
			form := url.Values{}

			form.Add("project", "bc39e850-0a1f-446f-a112-3e9a5b3134f0")
			form.Add("budgetItem", "cc5cbcb9-43cc-4062-b3d3-ea60a3c2e6d0")
			form.Add("quantity", "1")
			form.Add("cost", "1")

			_, err := createBudget(form, companyId, srv)
			if err == nil {
				t.Error("Expected an error and got none")
			}

			if err.Error() != "Ya existe partida cc5cbcb9-43cc-4062-b3d3-ea60a3c2e6d0 en el proyecto bc39e850-0a1f-446f-a112-3e9a5b3134f0" {
				t.Errorf("Expected 'Ya existe partida cc5cbcb9-43cc-4062-b3d3-ea60a3c2e6d0 en el proyecto bc39e850-0a1f-446f-a112-3e9a5b3134f0' and got '%s'", err.Error())
			}
		})
	})

	t.Run("create budget error", func(t *testing.T) {
		t.Run("invalid project id", func(t *testing.T) {
			t.Run("empty project id", func(t *testing.T) {
				form := url.Values{}

				form.Add("budgetItem", uuid.New().String())
				form.Add("quantity", "1")
				form.Add("cost", "1")

				_, err := createBudget(form, companyId, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "Seleccione un proyecto" {
					t.Errorf("Expected 'Seleccione un proyecto' and got '%s'", err.Error())
				}
			})

			t.Run("invalid project id", func(t *testing.T) {
				form := url.Values{}

				form.Add("project", "invalid-project-id")
				form.Add("budgetItem", uuid.New().String())
				form.Add("quantity", "1")
				form.Add("cost", "1")

				_, err := createBudget(form, companyId, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "invalid UUID length: 18" {
					t.Errorf("Expected 'invalid UUID length: 18' and got '%s'", err.Error())
				}
			})
		})

		t.Run("invalid budget item", func(t *testing.T) {
			t.Run("empty budget item id", func(t *testing.T) {
				form := url.Values{}

				form.Add("project", uuid.New().String())
				form.Add("quantity", "1")
				form.Add("cost", "1")

				_, err := createBudget(form, companyId, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "Seleccione un partida" {
					t.Errorf("Expected 'Seleccione un partida' and got '%s'", err.Error())
				}
			})

			t.Run("invalid project id", func(t *testing.T) {
				form := url.Values{}

				form.Add("budgetItem", "invalid-project-id")
				form.Add("project", uuid.New().String())
				form.Add("quantity", "1")
				form.Add("cost", "1")

				_, err := createBudget(form, companyId, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "invalid UUID length: 18" {
					t.Errorf("Expected 'invalid UUID length: 18' and got '%s'", err.Error())
				}
			})
		})

		t.Run("invalid quantity", func(t *testing.T) {
			t.Run("empty quantity", func(t *testing.T) {
				form := url.Values{}

				form.Add("project", uuid.New().String())
				form.Add("budgetItem", uuid.New().String())
				form.Add("cost", "1")

				_, err := createBudget(form, companyId, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "cantidad es requerido" {
					t.Errorf("Expected 'cantidad es requerido' and got '%s'", err.Error())
				}
			})

			t.Run("invalid quantity", func(t *testing.T) {
				form := url.Values{}

				form.Add("project", uuid.New().String())
				form.Add("budgetItem", uuid.New().String())
				form.Add("quantity", "invalid")
				form.Add("cost", "1")

				_, err := createBudget(form, companyId, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "cantidad debe ser un número válido" {
					t.Errorf("Expected 'cantidad debe ser un número válido' and got '%s'", err.Error())
				}
			})

			t.Run("quantity must be a positive number", func(t *testing.T) {
				form := url.Values{}

				form.Add("project", uuid.New().String())
				form.Add("budgetItem", uuid.New().String())
				form.Add("quantity", "-4")
				form.Add("cost", "1")

				_, err := createBudget(form, companyId, srv)
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

				form.Add("project", uuid.New().String())
				form.Add("budgetItem", uuid.New().String())
				form.Add("quantity", "1")

				_, err := createBudget(form, companyId, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "costo es requerido" {
					t.Errorf("Expected 'costo es requerido' and got '%s'", err.Error())
				}
			})

			t.Run("invalid cost", func(t *testing.T) {
				form := url.Values{}

				form.Add("project", uuid.New().String())
				form.Add("budgetItem", uuid.New().String())
				form.Add("quantity", "1")
				form.Add("cost", "invalid")

				_, err := createBudget(form, companyId, srv)
				if err == nil {
					t.Error("Expected an error and got none")
				}

				if err.Error() != "costo debe ser un número válido" {
					t.Errorf("Expected 'costo debe ser un número válido' and got '%s'", err.Error())
				}
			})

			t.Run("cost must be a positive number", func(t *testing.T) {
				form := url.Values{}

				form.Add("project", uuid.New().String())
				form.Add("budgetItem", uuid.New().String())
				form.Add("quantity", "1")
				form.Add("cost", "-4")

				_, err := createBudget(form, companyId, srv)
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
