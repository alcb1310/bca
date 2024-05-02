package server

import (
	"testing"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
	"bca-go-final/mocks"
)

var companyId = uuid.New()
var projectId = uuid.New()

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
	db := mocks.NewServiceMock()
	_, router := NewServer(db)
	_ = router

	t.Run("Successful budgets request", func(t *testing.T) {
		t.Run("Should return 200 when GET budgets", func(t *testing.T) {
			t.Run("No filter", func(t *testing.T) {
				t.Skip()
				// response := httptest.NewRecorder()
				// request, _ := http.NewRequest(http.MethodGet, "/bca/partials/budget", nil)
				//
				// db.On("GetBudgets", companyId, uuid.Nil, "").Return([]types.GetBudget{}, nil)
				// router.BudgetsTable(response, request)
				//
				// assert.Equal(t, response.Code, http.StatusOK)
				// assert.Contains(t, response.Body, "Por Gastar")
			})

			// t.Run("With filter", func(t *testing.T) {
			// 	t.Skip()
			// 	response := httptest.NewRecorder()
			// 	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/bca/partials/budget?proyecto=%s&buscar=%s", projectDuplicate, "term"), nil)
			// 	router.BudgetsTable(response, request)
			//
			// 	got := response.Code
			// 	want := http.StatusOK
			// 	if got != want {
			// 		t.Errorf("got %d, want %d", got, want)
			// 	}
			// })
		})

		// t.Run("Valid POST request", func(t *testing.T) {
		// 	t.Run("Should return 200 when POST budgets", func(t *testing.T) {
		// 		response := httptest.NewRecorder()
		//
		// 		form := url.Values{}
		// 		form.Add("project", uuid.New().String())
		// 		form.Add("budgetItem", uuid.New().String())
		// 		form.Add("quantity", "1")
		// 		form.Add("cost", "1")
		//
		// 		request := &http.Request{
		// 			Method: http.MethodPost,
		// 			URL: &url.URL{
		// 				Path: "/bca/partials/budget",
		// 			},
		// 			Form: form,
		// 		}
		// 		router.BudgetsTable(response, request)
		//
		// 		got := response.Code
		// 		want := http.StatusOK
		// 		if got != want {
		// 			t.Errorf("got %d, want %d", got, want)
		// 			t.Error(response.Body.String())
		// 		}
		// 	})
		//
		// 	t.Run("Should return 409 when duplicate budget", func(t *testing.T) {
		// 		response := httptest.NewRecorder()
		//
		// 		form := url.Values{}
		// 		form.Add("project", projectDuplicate)
		// 		form.Add("budgetItem", budgetItemDuplicate)
		// 		form.Add("quantity", "1")
		// 		form.Add("cost", "1")
		//
		// 		request := &http.Request{
		// 			Method: http.MethodPost,
		// 			URL: &url.URL{
		// 				Path: "/bca/partials/budget",
		// 			},
		// 			Form: form,
		// 		}
		// 		router.BudgetsTable(response, request)
		//
		// 		got := response.Code
		// 		want := http.StatusConflict
		// 		if got != want {
		// 			t.Errorf("got %d, want %d", got, want)
		// 			t.Error(response.Body.String())
		// 		}
		//
		// 		expected := fmt.Sprintf("Ya existe partida %s en el proyecto %s", budgetItemDuplicate, projectDuplicate)
		// 		received := strings.Trim(response.Body.String(), "\n")
		// 		if received != expected {
		// 			t.Errorf("got %s, want %s", received, expected)
		// 		}
		// 	})
		// })

		// t.Run("Unsuccessful budgets request", func(t *testing.T) {
		// 	t.Run("POST request", func(t *testing.T) {
		// 		t.Run("Validate project id", func(t *testing.T) {
		// 			response := httptest.NewRecorder()
		// 			form := url.Values{}
		// 			form.Add("budgetItem", uuid.New().String())
		// 			form.Add("quantity", "1")
		// 			form.Add("cost", "1")
		//
		// 			t.Run("no project id", func(t *testing.T) {
		// 				request := &http.Request{
		// 					Method: http.MethodPost,
		// 					URL: &url.URL{
		// 						Path: "/bca/partials/budget",
		// 					},
		// 					Form: form,
		// 				}
		// 				router.BudgetsTable(response, request)
		//
		// 				got := response.Code
		// 				want := http.StatusBadRequest
		// 				if got != want {
		// 					t.Errorf("got %d, want %d", got, want)
		// 					t.Error(response.Body.String())
		// 				}
		//
		// 				expected := "Seleccione un proyecto"
		// 				received := strings.Trim(response.Body.String(), "\n")
		// 				if received != expected {
		// 					t.Errorf("got %s, want %s", received, expected)
		// 				}
		// 			})
		//
		// 			t.Run("invalid project id", func(t *testing.T) {
		// 				form.Del("project")
		// 				form.Add("project", "invalid")
		// 				request := &http.Request{
		// 					Method: http.MethodPost,
		// 					URL: &url.URL{
		// 						Path: "/bca/partials/budget",
		// 					},
		// 					Form: form,
		// 				}
		// 				router.BudgetsTable(response, request)
		//
		// 				got := response.Code
		// 				want := http.StatusBadRequest
		// 				if got != want {
		// 					t.Errorf("got %d, want %d", got, want)
		// 					t.Error(response.Body.String())
		// 				}
		//
		// 				expected := "invalid UUID"
		// 				received := strings.Trim(response.Body.String(), "\n")
		// 				if !strings.Contains(received, expected) {
		// 					t.Errorf("got %s, want %s", received, expected)
		// 				}
		// 			})
		// 		})
		//
		// 		t.Run("Validate budget item id", func(t *testing.T) {
		// 			response := httptest.NewRecorder()
		// 			form := url.Values{}
		// 			form.Add("project", uuid.New().String())
		// 			form.Add("quantity", "1")
		// 			form.Add("cost", "1")
		//
		// 			t.Run("no budget item id", func(t *testing.T) {
		// 				request := &http.Request{
		// 					Method: http.MethodPost,
		// 					URL: &url.URL{
		// 						Path: "/bca/partials/budget",
		// 					},
		// 					Form: form,
		// 				}
		// 				router.BudgetsTable(response, request)
		//
		// 				got := response.Code
		// 				want := http.StatusBadRequest
		// 				if got != want {
		// 					t.Errorf("got %d, want %d", got, want)
		// 					t.Error(response.Body.String())
		// 				}
		//
		// 				expected := "Seleccione un partida"
		// 				received := strings.Trim(response.Body.String(), "\n")
		// 				if received != expected {
		// 					t.Errorf("got %s, want %s", received, expected)
		// 				}
		// 			})
		//
		// 			t.Run("invalid budget item id", func(t *testing.T) {
		// 				form.Del("budgetItem")
		// 				form.Add("budgetItem", "invalid")
		// 				request := &http.Request{
		// 					Method: http.MethodPost,
		// 					URL: &url.URL{
		// 						Path: "/bca/partials/budget",
		// 					},
		// 					Form: form,
		// 				}
		// 				router.BudgetsTable(response, request)
		//
		// 				got := response.Code
		// 				want := http.StatusBadRequest
		// 				if got != want {
		// 					t.Errorf("got %d, want %d", got, want)
		// 					t.Error(response.Body.String())
		// 				}
		//
		// 				expected := "invalid UUID"
		// 				received := strings.Trim(response.Body.String(), "\n")
		// 				if !strings.Contains(received, expected) {
		// 					t.Errorf("got %s, want %s", received, expected)
		// 				}
		// 			})
		// 		})
		//
		// 		t.Run("Validate quantity", func(t *testing.T) {
		// 			response := httptest.NewRecorder()
		// 			form := url.Values{}
		// 			form.Add("budgetItem", uuid.New().String())
		// 			form.Add("project", uuid.New().String())
		// 			form.Add("cost", "1")
		//
		// 			t.Run("no quantity", func(t *testing.T) {
		// 				request := &http.Request{
		// 					Method: http.MethodPost,
		// 					URL: &url.URL{
		// 						Path: "/bca/partials/budget",
		// 					},
		// 					Form: form,
		// 				}
		// 				router.BudgetsTable(response, request)
		//
		// 				got := response.Code
		// 				want := http.StatusBadRequest
		// 				if got != want {
		// 					t.Errorf("got %d, want %d", got, want)
		// 					t.Error(response.Body.String())
		// 				}
		//
		// 				expected := "cantidad es requerido"
		// 				received := strings.Trim(response.Body.String(), "\n")
		// 				if !strings.Contains(received, expected) {
		// 					t.Errorf("got %s, want %s", received, expected)
		// 				}
		// 			})
		//
		// 			t.Run("invalid quantity", func(t *testing.T) {
		// 				form.Add("quantity", "invalid")
		// 				request := &http.Request{
		// 					Method: http.MethodPost,
		// 					URL: &url.URL{
		// 						Path: "/bca/partials/budget",
		// 					},
		// 					Form: form,
		// 				}
		// 				router.BudgetsTable(response, request)
		//
		// 				got := response.Code
		// 				want := http.StatusBadRequest
		// 				if got != want {
		// 					t.Errorf("got %d, want %d", got, want)
		// 					t.Error(response.Body.String())
		// 				}
		//
		// 				expected := "cantidad debe ser un número válido"
		// 				received := strings.Trim(response.Body.String(), "\n")
		// 				if !strings.Contains(received, expected) {
		// 					t.Errorf("got %s, want %s", received, expected)
		// 				}
		// 			})
		//
		// 			t.Run("negative quantity", func(t *testing.T) {
		// 				form.Del("quantity")
		// 				form.Add("quantity", "-1")
		// 				request := &http.Request{
		// 					Method: http.MethodPost,
		// 					URL: &url.URL{
		// 						Path: "/bca/partials/budget",
		// 					},
		// 					Form: form,
		// 				}
		// 				router.BudgetsTable(response, request)
		//
		// 				got := response.Code
		// 				want := http.StatusBadRequest
		// 				if got != want {
		// 					t.Errorf("got %d, want %d", got, want)
		// 					t.Error(response.Body.String())
		// 				}
		//
		// 				expected := "cantidad debe ser un número positivo"
		// 				received := strings.Trim(response.Body.String(), "\n")
		// 				if !strings.Contains(received, expected) {
		// 					t.Errorf("got %s, want %s", received, expected)
		// 				}
		// 			})
		// 		})
		//
		// 		t.Run("Validate cost", func(t *testing.T) {
		// 			response := httptest.NewRecorder()
		// 			form := url.Values{}
		// 			form.Add("budgetItem", uuid.New().String())
		// 			form.Add("project", uuid.New().String())
		// 			form.Add("quantity", "1")
		//
		// 			t.Run("no cost", func(t *testing.T) {
		// 				request := &http.Request{
		// 					Method: http.MethodPost,
		// 					URL: &url.URL{
		// 						Path: "/bca/partials/budget",
		// 					},
		// 					Form: form,
		// 				}
		// 				router.BudgetsTable(response, request)
		//
		// 				got := response.Code
		// 				want := http.StatusBadRequest
		// 				if got != want {
		// 					t.Errorf("got %d, want %d", got, want)
		// 					t.Error(response.Body.String())
		// 				}
		//
		// 				expected := "costo es requerido"
		// 				received := strings.Trim(response.Body.String(), "\n")
		// 				if !strings.Contains(received, expected) {
		// 					t.Errorf("got %s, want %s", received, expected)
		// 				}
		// 			})
		//
		// 			t.Run("invalid cost", func(t *testing.T) {
		// 				form.Add("cost", "invalid")
		// 				request := &http.Request{
		// 					Method: http.MethodPost,
		// 					URL: &url.URL{
		// 						Path: "/bca/partials/budget",
		// 					},
		// 					Form: form,
		// 				}
		// 				router.BudgetsTable(response, request)
		//
		// 				got := response.Code
		// 				want := http.StatusBadRequest
		// 				if got != want {
		// 					t.Errorf("got %d, want %d", got, want)
		// 					t.Error(response.Body.String())
		// 				}
		//
		// 				expected := "costo debe ser un número válido"
		// 				received := strings.Trim(response.Body.String(), "\n")
		// 				if !strings.Contains(received, expected) {
		// 					t.Errorf("got %s, want %s", received, expected)
		// 				}
		// 			})
		//
		// 			t.Run("negative cost", func(t *testing.T) {
		// 				form.Del("cost")
		// 				form.Add("cost", "-1")
		// 				request := &http.Request{
		// 					Method: http.MethodPost,
		// 					URL: &url.URL{
		// 						Path: "/bca/partials/budget",
		// 					},
		// 					Form: form,
		// 				}
		// 				router.BudgetsTable(response, request)
		//
		// 				got := response.Code
		// 				want := http.StatusBadRequest
		// 				if got != want {
		// 					t.Errorf("got %d, want %d", got, want)
		// 					t.Error(response.Body.String())
		// 				}
		//
		// 				expected := "costo debe ser un número positivo"
		// 				received := strings.Trim(response.Body.String(), "\n")
		// 				if !strings.Contains(received, expected) {
		// 					t.Errorf("got %s, want %s", received, expected)
		// 				}
		// 			})
		// 		})
		// 	})
		// })
	})
	//	t.Run("Unsupported request methods", func(t *testing.T) {
	//		t.Run("PATCH", func(t *testing.T) {
	//			response := httptest.NewRecorder()
	//			request := &http.Request{
	//				Method: http.MethodPatch,
	//				URL: &url.URL{
	//					Path: "/bca/partials/budget",
	//				},
	//			}
	//			router.BudgetsTable(response, request)
	//			got := response.Code
	//			want := http.StatusMethodNotAllowed
	//			if got != want {
	//				t.Errorf("got %d, want %d", got, want)
	//			}
	//		})
	//
	//		t.Run("DELETE", func(t *testing.T) {
	//			response := httptest.NewRecorder()
	//			request := &http.Request{
	//				Method: http.MethodDelete,
	//				URL: &url.URL{
	//					Path: "/bca/partials/budget",
	//				},
	//			}
	//			router.BudgetsTable(response, request)
	//			got := response.Code
	//			want := http.StatusMethodNotAllowed
	//			if got != want {
	//				t.Errorf("got %d, want %d", got, want)
	//			}
	//		})
	//
	//		t.Run("PUT", func(t *testing.T) {
	//			response := httptest.NewRecorder()
	//			request := &http.Request{
	//				Method: http.MethodPut,
	//				URL: &url.URL{
	//					Path: "/bca/partials/budget",
	//				},
	//			}
	//			router.BudgetsTable(response, request)
	//			got := response.Code
	//			want := http.StatusMethodNotAllowed
	//			if got != want {
	//				t.Errorf("got %d, want %d", got, want)
	//			}
	//		})
	//	})
}

//
// func TestBudgetAdd(t *testing.T) {
// 	db := mocks.NewServiceMock()
// 	_, srv := NewServer(db)
//
// 	t.Run("display budget form", func(t *testing.T) {
// 		response := httptest.NewRecorder()
// 		request := &http.Request{
// 			Method: http.MethodGet,
// 			URL: &url.URL{
// 				Path: "/bca/partials/budget/add",
// 			},
// 		}
// 		srv.BudgetAdd(response, request)
// 		got := response.Code
// 		want := http.StatusOK
// 		if got != want {
// 			t.Errorf("got %d, want %d", got, want)
// 		}
// 		expected := "Agregar Presupuesto"
//
// 		if !strings.Contains(response.Body.String(), expected) {
// 			t.Errorf("got %s, want %s", response.Body.String(), expected)
// 		}
//
// 	})
// }
//
// func TestBudgetEdit(t *testing.T) {
// 	db := mocks.NewServiceMock()
// 	_, srv := NewServer(db)
// 	projectId := uuid.New().String()
// 	budgetItemId := uuid.New().String()
//
// 	t.Run("unimplemented method", func(t *testing.T) {
// 		response := httptest.NewRecorder()
// 		request := &http.Request{
// 			Method: http.MethodPost,
// 			URL: &url.URL{
// 				Path: fmt.Sprintf("/bca/partials/budget/%s/%s", projectId, budgetItemId),
// 			},
// 		}
//
// 		srv.BudgetEdit(response, request)
// 		got := response.Code
// 		want := http.StatusMethodNotAllowed
// 		if got != want {
// 			t.Errorf("got %d, want %d", got, want)
// 		}
// 	})
//
// 	t.Run("display edit form", func(t *testing.T) {
// 		response := httptest.NewRecorder()
// 		request := &http.Request{
// 			Method: http.MethodGet,
// 			URL: &url.URL{
// 				Path: fmt.Sprintf("/bca/partials/budget/%s/%s", projectId, budgetItemId),
// 			},
// 		}
// 		srv.BudgetEdit(response, request)
// 		got := response.Code
// 		want := http.StatusOK
// 		if got != want {
// 			t.Errorf("got %d, want %d", got, want)
// 		}
//
// 		expected := "Editar Presupuesto"
//
// 		if !strings.Contains(response.Body.String(), expected) {
// 			t.Errorf("got %s, want %s", response.Body.String(), expected)
// 		}
// 	})
//
// 	t.Run("update budget", func(t *testing.T) {
// 		t.Run("valid budget data", func(t *testing.T) {
// 			t.Run("success", func(t *testing.T) {
// 				response := httptest.NewRecorder()
// 				form := url.Values{}
// 				form.Add("quantity", "1")
// 				form.Add("cost", "1")
//
// 				request := &http.Request{
// 					Method: http.MethodPut,
// 					URL: &url.URL{
// 						Path: fmt.Sprintf("/bca/partials/budget/%s/%s", projectId, budgetItemId),
// 					},
// 					Form: form,
// 				}
//
// 				srv.BudgetEdit(response, request)
// 				got := response.Code
// 				want := http.StatusOK
// 				if got != want {
// 					t.Errorf("got %d, want %d", got, want)
// 				}
//
// 				expected := "<table"
//
// 				if !strings.Contains(response.Body.String(), expected) {
// 					t.Errorf("got %s, want %s", response.Body.String(), expected)
// 				}
// 			})
// 		})
//
// 		t.Run("invalid budget data", func(t *testing.T) {
// 			t.Run("invalid quantity", func(t *testing.T) {
// 				t.Run("empty quantity", func(t *testing.T) {
// 					response := httptest.NewRecorder()
// 					form := url.Values{}
// 					form.Add("cost", "1")
//
// 					request := &http.Request{
// 						Method: http.MethodPut,
// 						URL: &url.URL{
// 							Path: fmt.Sprintf("/bca/partials/budget/%s/%s", projectId, budgetItemId),
// 						},
// 						Form: form,
// 					}
//
// 					srv.BudgetEdit(response, request)
//
// 					got := response.Code
// 					want := http.StatusBadRequest
// 					if got != want {
// 						t.Errorf("got %d, want %d", got, want)
// 					}
// 					expected := "cantidad es requerido"
// 					if !strings.Contains(response.Body.String(), expected) {
// 						t.Errorf("got %s, want %s", response.Body.String(), expected)
// 					}
// 				})
//
// 				t.Run("invalid quantity", func(t *testing.T) {
// 					response := httptest.NewRecorder()
// 					form := url.Values{}
// 					form.Add("quantity", "invalid")
// 					form.Add("cost", "1")
//
// 					request := &http.Request{
// 						Method: http.MethodPut,
// 						URL: &url.URL{
// 							Path: fmt.Sprintf("/bca/partials/budget/%s/%s", projectId, budgetItemId),
// 						},
// 						Form: form,
// 					}
//
// 					srv.BudgetEdit(response, request)
//
// 					got := response.Code
// 					want := http.StatusBadRequest
// 					if got != want {
// 						t.Errorf("got %d, want %d", got, want)
// 					}
//
// 					expected := "cantidad debe ser un número válido"
// 					if !strings.Contains(response.Body.String(), expected) {
// 						t.Errorf("got %s, want %s", response.Body.String(), expected)
// 					}
// 				})
//
// 				t.Run("negative quantity", func(t *testing.T) {
// 					response := httptest.NewRecorder()
// 					form := url.Values{}
// 					form.Add("quantity", "-1")
// 					form.Add("cost", "1")
//
// 					request := &http.Request{
// 						Method: http.MethodPut,
// 						URL: &url.URL{
// 							Path: fmt.Sprintf("/bca/partials/budget/%s/%s", projectId, budgetItemId),
// 						},
// 						Form: form,
// 					}
//
// 					srv.BudgetEdit(response, request)
//
// 					got := response.Code
// 					want := http.StatusBadRequest
// 					if got != want {
// 						t.Errorf("got %d, want %d", got, want)
// 					}
//
// 					expected := "cantidad debe ser un número positivo"
// 					if !strings.Contains(response.Body.String(), expected) {
// 						t.Errorf("got %s, want %s", response.Body.String(), expected)
// 					}
// 				})
// 			})
//
// 			t.Run("invalid quantity", func(t *testing.T) {
// 				t.Run("empty cost", func(t *testing.T) {
// 					response := httptest.NewRecorder()
// 					form := url.Values{}
// 					form.Add("quantity", "1")
//
// 					request := &http.Request{
// 						Method: http.MethodPut,
// 						URL: &url.URL{
// 							Path: fmt.Sprintf("/bca/partials/budget/%s/%s", projectId, budgetItemId),
// 						},
// 						Form: form,
// 					}
//
// 					srv.BudgetEdit(response, request)
//
// 					got := response.Code
// 					want := http.StatusBadRequest
// 					if got != want {
// 						t.Errorf("got %d, want %d", got, want)
// 					}
// 					expected := "costo es requerido"
// 					if !strings.Contains(response.Body.String(), expected) {
// 						t.Errorf("got %s, want %s", response.Body.String(), expected)
// 					}
// 				})
//
// 				t.Run("invalid cost", func(t *testing.T) {
// 					response := httptest.NewRecorder()
// 					form := url.Values{}
// 					form.Add("cost", "invalid")
// 					form.Add("quantity", "1")
//
// 					request := &http.Request{
// 						Method: http.MethodPut,
// 						URL: &url.URL{
// 							Path: fmt.Sprintf("/bca/partials/budget/%s/%s", projectId, budgetItemId),
// 						},
// 						Form: form,
// 					}
//
// 					srv.BudgetEdit(response, request)
//
// 					got := response.Code
// 					want := http.StatusBadRequest
// 					if got != want {
// 						t.Errorf("got %d, want %d", got, want)
// 					}
//
// 					expected := "costo debe ser un número válido"
// 					if !strings.Contains(response.Body.String(), expected) {
// 						t.Errorf("got %s, want %s", response.Body.String(), expected)
// 					}
// 				})
//
// 				t.Run("negative cost", func(t *testing.T) {
// 					response := httptest.NewRecorder()
// 					form := url.Values{}
// 					form.Add("cost", "-1")
// 					form.Add("quantity", "1")
//
// 					request := &http.Request{
// 						Method: http.MethodPut,
// 						URL: &url.URL{
// 							Path: fmt.Sprintf("/bca/partials/budget/%s/%s", projectId, budgetItemId),
// 						},
// 						Form: form,
// 					}
//
// 					srv.BudgetEdit(response, request)
//
// 					got := response.Code
// 					want := http.StatusBadRequest
// 					if got != want {
// 						t.Errorf("got %d, want %d", got, want)
// 					}
//
// 					expected := "costo debe ser un número positivo"
// 					if !strings.Contains(response.Body.String(), expected) {
// 						t.Errorf("got %s, want %s", response.Body.String(), expected)
// 					}
// 				})
// 			})
// 		})
// 	})
// }
