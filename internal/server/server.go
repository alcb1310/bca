package server

import (
	"bca-go-final/internal/database"
	"bca-go-final/internal/utils"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Server struct {
	DB        database.Service
	Router    *chi.Mux
	TokenAuth *jwtauth.JWTAuth
}

func NewServer(db database.Service, secret string) *Server {
	s := &Server{
		DB:        db,
		Router:    chi.NewRouter(),
		TokenAuth: jwtauth.New("HS256", []byte(secret), nil),
	}
	s.Router.Use(middleware.Logger)

	s.Router.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))

	s.Router.Get("/", s.HelloWorldHandler)
	s.RegisterRoutes(s.Router)

	s.Router.Get("/login", s.DisplayLogin) // fully tested
	s.Router.Post("/login", s.LoginView) // fully tested

	s.Router.Route("/bca", func(r chi.Router) {
		r.Use(jwtauth.Verifier(s.TokenAuth))
		r.Use(authenticator())

		r.HandleFunc("/", s.BcaView)
		r.HandleFunc("/logout", s.Logout)
		r.Post("/dummy-data", s.loadDummyDataHandler)

		r.Route("/transacciones", func(r chi.Router) {
			r.HandleFunc("/presupuesto", s.Budget)
			r.HandleFunc("/facturas", s.Invoice)
			r.HandleFunc("/facturas/crear", s.InvoiceAdd) // fullly unit tested
			r.HandleFunc("/cierre", s.Closure)            // fullly unit tested
		})

		r.Route("/reportes", func(r chi.Router) {
			r.HandleFunc("/actual", s.Actual)
			r.HandleFunc("/actual/generar", s.ActualGenerate)
			r.HandleFunc("/cuadre", s.Balance)
			r.HandleFunc("/historico", s.Historic)
			r.HandleFunc("/gastado", s.Spent)
			r.HandleFunc("/gastado/{projectId}/{budgetItemId}/{date}", s.SpentByBudgetItem) // convert

			r.Route("/excel", func(r chi.Router) {
				r.HandleFunc("/cuadre", s.BalanceExcel)
				r.HandleFunc("/actual", s.ActualExcel)
				r.HandleFunc("/historico", s.HistoricExcel)
				r.HandleFunc("/gastado", s.SpentExcel)
			})
		})

		r.Route("/configuracion", func(r chi.Router) {
			r.HandleFunc("/partidas", s.BudgetItems)
			r.HandleFunc("/proveedores", s.Suppliers)
			r.HandleFunc("/proyectos", s.Projects) // fully tested
			r.HandleFunc("/categorias", s.Categories)
			r.HandleFunc("/materiales", s.Materiales)
			r.HandleFunc("/rubros", s.Rubros)
			r.HandleFunc("/rubros/crear", s.RubrosAdd) // fully integration tested
		})

		r.Route("/user", func(r chi.Router) {
			r.HandleFunc("/perfil", s.Profile)
			r.HandleFunc("/admin", s.Admin)
			r.HandleFunc("/cambio", s.ChangePassword)
		})

		r.Route("/costo-unitario", func(r chi.Router) {
			r.HandleFunc("/cantidades", s.UnitQuantity)
			r.HandleFunc("/analisis", s.UnitAnalysis)
		})

		r.Route("/partials", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.HandleFunc("/", s.UsersTable) // fully unit tested
				r.HandleFunc("/add", s.UserAdd)
				r.HandleFunc("/edit/{id}", s.UserEdit) // convert
				r.HandleFunc("/{id}", s.SingleUser)    // convert
			})

			r.Route("/projects", func(r chi.Router) {
				r.HandleFunc("/", s.ProjectsTable) // fully tested
				r.HandleFunc("/add", s.ProjectAdd)
				r.HandleFunc("/edit/{id}", s.ProjectEditSave) // convert fully tested
				r.HandleFunc("/{id}", s.ProjectEdit)          // convert fully tested
			})

			r.Route("/suppliers", func(r chi.Router) {
				r.HandleFunc("/", s.SuppliersTable) // fully tested
				r.HandleFunc("/add", s.SupplierAdd) 
				r.HandleFunc("/edit/{id}", s.SuppliersEditSave) // convert fully tested
				r.HandleFunc("/{id}", s.SuppliersEdit)          // convert fully tested
			})

			r.Route("/budget-item", func(r chi.Router) {
				r.HandleFunc("/", s.BudgetItemsTable) // fully tested
				r.HandleFunc("/add", s.BudgetItemAdd)
				r.HandleFunc("/{id}", s.BudgetItemEdit) // convert fully unit tested
			})

			r.Route("/budgets", func(r chi.Router) {
				r.HandleFunc("/", s.BudgetsTable) // fully unit tested
				r.HandleFunc("/add", s.BudgetAdd)
				r.HandleFunc("/{projectId}/{budgetItemId}", s.BudgetEdit) // convert fully unit tested
			})

			r.Route("/invoices", func(r chi.Router) {
				r.HandleFunc("/", s.InvoicesTable)
				r.HandleFunc("/{id}", s.InvoiceEdit) // convert fully unit tested

				r.Route("/{invoiceId}/details", func(r chi.Router) {
					r.HandleFunc("/", s.DetailsTable)              // convert fully unit tested
					r.HandleFunc("/add", s.DetailsAdd)             // convert
					r.HandleFunc("/{budgetItemId}", s.DetailsEdit) // convert
				})
			})

			r.Route("/categories", func(r chi.Router) {
				r.HandleFunc("/", s.CategoriesTable) // fully tested
				r.HandleFunc("/add", s.CategoryAdd)
				r.HandleFunc("/{id}", s.EditCategory) // convert fully tested
			})

			r.Route("/materiales", func(r chi.Router) {
				r.HandleFunc("/", s.MaterialsTable) // fully tested
				r.HandleFunc("/add", s.MaterialsAdd)
				r.HandleFunc("/{id}", s.MaterialsEdit) // convert fully tested
			})

			r.Route("/rubros", func(r chi.Router) {
				r.HandleFunc("/", s.RubrosTable) // fully tested
				r.HandleFunc("/{id}", s.MaterialsByItem)                               // convert
				r.HandleFunc("/{id}/material", s.MaterialByItemForm)                   // convert fully unit tested
				r.HandleFunc("/{id}/material/{materialId}", s.MaterialItemsOperations) // convert fully unit tested
			})

			r.Route("/cantidades", func(r chi.Router) {
				r.HandleFunc("/", s.CantidadesTable)
				r.HandleFunc("/add", s.CantidadesAdd)   // fully unit tested
				r.HandleFunc("/{id}", s.CantidadesEdit) // convert fully unit tested
			})

			r.HandleFunc("/analisis", s.AnalysisTable)
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
