package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/alcb1310/bca/internal/database"
	"github.com/alcb1310/bca/internal/utils"
)

type Server struct {
	DB        database.Service
	Router    *chi.Mux
	TokenAuth *jwtauth.JWTAuth
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func NewServer(db database.Service, secret string) *Server {
	s := &Server{
		DB:        db,
		Router:    chi.NewRouter(),
		TokenAuth: jwtauth.New("HS256", []byte(secret), nil),
	}
	s.Router.Use(middleware.Logger)
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	s.Router.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))

	s.Router.Get("/", s.HelloWorldHandler)
	s.RegisterRoutes(s.Router)

	s.Router.Get("/login", s.DisplayLogin)
	s.Router.Post("/login", s.LoginView)

	s.Router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(commonMiddleware)

		r.Post("/login", s.ApiLogin)

		r.Route("/users", func(r chi.Router) {
			r.Use(jwtauth.Verifier(s.TokenAuth))
			r.Use(authenticator())

			r.Get("/me", s.ApiGetCurrentUser)
			r.Get("/", s.ApiGetAllUsers)
			r.Post("/", s.ApiCreateUser)
			r.Put("/", s.ApiChangePassword)
			r.Delete("/{id}", s.ApiDeleteUser)
			r.Put("/{id}", s.ApiUpdateUser)
		})

    r.Route("/transacciones", func(r chi.Router) {
      r.Use(jwtauth.Verifier(s.TokenAuth))
      r.Use(authenticator())

      r.Route("/presupuestos", func(r chi.Router) {
        r.Get("/", s.ApiGetAllBudgets)
        r.Post("/", s.ApiCreateBudget)
        r.Put("/{projectId}/{budgetItemId}", s.ApiUpdateBudget)
      })

      r.Route("/facturas", func(r chi.Router) {
        r.Get("/", s.ApiGetAllInvoices)
      })
    })

		r.Route("/parametros", func(r chi.Router) {
			r.Use(jwtauth.Verifier(s.TokenAuth))
			r.Use(authenticator())

			r.Route("/partidas", func(r chi.Router) {
				r.Get("/", s.ApiGetAllBudgetItems)
				r.Post("/", s.ApiCreateBudgetItem)
				r.Put("/{id}", s.ApiUpdateBudgetItem)
			})

			r.Route("/categorias", func(r chi.Router) {
				r.Get("/", s.ApiGetAllCategories)
				r.Post("/", s.ApiCreateCategory)
				r.Put("/{id}", s.ApiUpdateCategory)
			})

			r.Route("/materiales", func(r chi.Router) {
				r.Get("/", s.ApiGetAllMaterials)
				r.Post("/", s.ApiCreateMaterial)
				r.Put("/{id}", s.ApiUpdateMaterial)
			})

			r.Route("/proyectos", func(r chi.Router) {
				r.Get("/", s.ApiGetAllProjects)
				r.Post("/", s.ApiCreateProject)
				r.Put("/{id}", s.ApiUpdateProject)
			})

			r.Route("/proveedores", func(r chi.Router) {
				r.Get("/", s.ApiGetAllSuppliers)
				r.Post("/", s.ApiCreateSupplier)
				r.Put("/{id}", s.ApiUpdateSupplier)
			})

			r.Route("/rubros", func(r chi.Router) {
				r.Get("/", s.ApiGetAllRubros)
				r.Post("/", s.ApiCreateRubros)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", s.ApiGetRubro)
					r.Put("/", s.ApiUpdateRubro)

					r.Route("/materiales", func(r chi.Router) {
						r.Get("/", s.ApiGetAllItemsMaterials)
						r.Post("/", s.ApiCreateItemsMaterials)
						r.Delete("/{materialId}", s.ApiDeleteItemsMaterials)
						r.Put("/{materialId}", s.ApiUpdateItemsMaterials)
					})
				})
			})
		})
	})

	s.Router.Route("/bca", func(r chi.Router) {
		r.Use(jwtauth.Verifier(s.TokenAuth))
		r.Use(authenticator())

		r.HandleFunc("/", s.BcaView)
		r.Post("/logout", s.Logout)
		r.Post("/dummy-data", s.loadDummyDataHandler)

		r.Route("/transacciones", func(r chi.Router) {
			r.Get("/presupuesto", s.Budget)
			r.Get("/facturas", s.Invoice)
			r.Post("/facturas/crear", s.InvoiceAdd)
			r.Get("/facturas/crear", s.InvoiceAddForm)
			r.Get("/cierre", s.ClosureForm)
			r.Post("/cierre", s.Closure)
		})

		r.Route("/reportes", func(r chi.Router) {
			r.Get("/actual", s.Actual)
			r.Get("/actual/generar", s.ActualGenerate)
			r.Post("/cuadre", s.RetreiveBalance)
			r.Get("/cuadre", s.GetBalance)
			r.Get("/historico", s.Historic)
			r.Get("/gastado", s.Spent)
			r.Get("/gastado/{projectId}/{budgetItemId}/{date}", s.SpentByBudgetItem)

			r.Route("/excel", func(r chi.Router) {
				r.Get("/cuadre", s.BalanceExcel)
				r.Get("/actual", s.ActualExcel)
				r.Get("/historico", s.HistoricExcel)
				r.Get("/gastado", s.SpentExcel)
			})
		})

		r.Route("/configuracion", func(r chi.Router) {
			r.Get("/partidas", s.BudgetItems)
			r.Get("/proveedores", s.Suppliers)
			r.Get("/proyectos", s.Projects)
			r.Get("/categorias", s.Categories)
			r.Get("/materiales", s.Materiales)
			r.Get("/rubros", s.Rubros)
			r.Get("/rubros/crear", s.RubrosAddForm)
			r.Put("/rubros/crear", s.RubrosEdit)
			r.Post("/rubros/crear", s.RubrosAdd)
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/perfil", s.Profile)
			r.Get("/admin", s.Admin)
			r.Get("/cambio", s.ChangePassword)
		})

		r.Route("/costo-unitario", func(r chi.Router) {
			r.Get("/cantidades", s.UnitQuantity)
			r.Get("/analisis", s.UnitAnalysis)
		})

		r.Route("/partials", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", s.UsersTableDisplay)
				r.Put("/", s.UsersChangePassword)
				r.Post("/", s.UsersCreate)
				r.Get("/add", s.UserAdd)
				r.Get("/edit/{id}", s.UserEdit)
				r.Get("/{id}", s.SingleUserGet)
				r.Delete("/{id}", s.DeleteUser)
				r.Put("/{id}", s.UpdateUser)
			})

			r.Route("/projects", func(r chi.Router) {
				r.Get("/", s.ProjectsTableDisplay)
				r.Post("/", s.ProjectsTable)
				r.Get("/add", s.ProjectAdd)
				r.Put("/edit/{id}", s.ProjectEditSave)
				r.Get("/{id}", s.ProyectDisplay)
			})

			r.Route("/suppliers", func(r chi.Router) {
				r.Get("/", s.SuppliersTableDisplay)
				r.Post("/", s.CreateSupplier)
				r.Get("/add", s.SupplierAdd)
				r.Put("/edit/{id}", s.EditSupplier)
				r.Get("/{id}", s.GetSupplier)
			})

			r.Route("/budget-item", func(r chi.Router) {
				r.Get("/", s.BudgetItemsTableDisplay)
				r.Post("/", s.BudgetItemsTable)
				r.Get("/add", s.BudgetItemAdd)
				r.Put("/{id}", s.EditBudgetItem)
				r.Get("/{id}", s.BudgetItemEdit)
			})

			r.Route("/budgets", func(r chi.Router) {
				r.Get("/", s.BudgetsTableDisplay)
				r.Post("/", s.BudgetsTable)
				r.Get("/add", s.BudgetAdd)
				r.Get("/{projectId}/{budgetItemId}", s.SingleBudget)
				r.Put("/{projectId}/{budgetItemId}", s.BudgetEdit)
			})

			r.Route("/invoices", func(r chi.Router) {
				r.Get("/", s.InvoicesTable)
				r.Get("/{id}", s.GetOneInvoice)
				r.Delete("/{id}", s.DeleteInvoice)
				r.Patch("/{id}", s.PatchInvoice)
				r.Put("/{id}", s.InvoiceEdit)

				r.Route("/{invoiceId}/details", func(r chi.Router) {
					r.Get("/", s.DetailsTableDisplay)
					r.Post("/", s.DetailsTable)
					r.Get("/add", s.DetailsAdd)
					r.Delete("/{budgetItemId}", s.DetailsEdit)
				})
			})

			r.Route("/categories", func(r chi.Router) {
				r.Get("/", s.CategoriesTableDisplay)
				r.Post("/", s.CategoriesTable)
				r.Get("/add", s.CategoryAdd)
				r.Get("/{id}", s.GetOneCategory)
				r.Put("/{id}", s.EditCategory)
			})

			r.Route("/materiales", func(r chi.Router) {
				r.Get("/", s.MaterialsTableDisplay)
				r.Post("/", s.MaterialsTable)
				r.Get("/add", s.MaterialsAdd)
				r.Get("/{id}", s.GetOneMaterial)
				r.Put("/{id}", s.MaterialsEdit)
			})

			r.Route("/rubros", func(r chi.Router) {
				r.Get("/", s.RubrosTable)
				r.Get("/{id}", s.MaterialsByItem)
				r.Get("/{id}/material", s.MaterialsByItemFormDisplay)
				r.Post("/{id}/material", s.MaterialByItemForm)
				r.Get("/{id}/material/{materialId}", s.MaterialItemsOperationsDisplay)
				r.Delete("/{id}/material/{materialId}", s.MaterialItemsOperationsDelete)
				r.Put("/{id}/material/{materialId}", s.MaterialItemsOperationsUpdate)
			})

			r.Route("/cantidades", func(r chi.Router) {
				r.Get("/", s.CantidadesTable)
				r.Get("/add", s.CantidadesAddDisplay)
				r.Post("/add", s.CantidadesAdd)
				r.Delete("/{id}", s.CantiadesDelete)
				r.Put("/{id}", s.CantidadesUpdate)
				r.Get("/{id}", s.CantidadesDisplay)
			})

			r.Get("/analisis", s.AnalysisTable)
		})
	})

	return s
}

func authenticator() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			if token == nil || jwt.Validate(token) != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Token is authenticated, pass it through
			marshalStr, _ := json.Marshal(token.PrivateClaims())
			ctxKey := utils.ContextKey("token")
			ctx := context.WithValue(r.Context(), ctxKey, marshalStr)
			r = r.Clone(ctx)

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
