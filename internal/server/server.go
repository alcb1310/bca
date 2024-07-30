package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"bca-go-final/internal/database"
)

type Server struct {
	DB     database.Service
	Router *chi.Mux
}

func NewServer(db database.Service) *Server {
	s := &Server{
		DB:     db,
		Router: chi.NewRouter(),
	}
	s.Router.Use(middleware.Logger)
	// s.Router.Use(s.authVerify)

	s.Router.Handle("/public/*", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))

	s.Router.Get("/", s.HelloWorldHandler)
	s.RegisterRoutes(s.Router)

	s.Router.Get("/login", s.DisplayLogin)
	s.Router.Post("/login", s.LoginView) // TODO: test form validation
	// TODO: properly migrate to chi router
	// TODO: change all mux.Vars(r)["id"] to chi.URLParam(r, "id")

	s.Router.Route("/bca", func(r chi.Router) {
		r.Use(s.authVerify)
		r.Get("/dummy-data", s.loadDummyDataHandler)
		r.HandleFunc("/", s.BcaView)
		r.HandleFunc("/logout", s.Logout)

		r.Route("/transacciones", func(r chi.Router) {
			r.HandleFunc("/presupuesto", s.Budget)
			r.HandleFunc("/facturas", s.Invoice)
			r.HandleFunc("/facturas/crear", s.InvoiceAdd)
			r.HandleFunc("/cierre", s.Closure)
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
			r.HandleFunc("/proyectos", s.Projects)
			r.HandleFunc("/categorias", s.Categories)
			r.HandleFunc("/materiales", s.Materiales)
			r.HandleFunc("/rubros", s.Rubros)
			r.HandleFunc("/rubros/crear", s.RubrosAdd)
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
				r.HandleFunc("/", s.UsersTable)
				r.HandleFunc("/add", s.UserAdd)
				r.HandleFunc("/edit/{id}", s.UserEdit) // convert
				r.HandleFunc("/{id}", s.SingleUser)    // convert
			})

			r.Route("/projects", func(r chi.Router) {
				r.HandleFunc("/", s.ProjectsTable)
				r.HandleFunc("/add", s.ProjectAdd)
				r.HandleFunc("/edit/{id}", s.ProjectEditSave) // convert
				r.HandleFunc("/{id}", s.ProjectEdit)          // convert
			})

			r.Route("/suppliers", func(r chi.Router) {
				r.HandleFunc("/", s.SuppliersTable)
				r.HandleFunc("/add", s.SupplierAdd)
				r.HandleFunc("/edit/{id}", s.SuppliersEditSave) // convert
				r.HandleFunc("/{id}", s.SuppliersEdit)          // convert
			})

			r.Route("/budget-item", func(r chi.Router) {
				r.HandleFunc("/", s.BudgetItemsTable)
				r.HandleFunc("/add", s.BudgetItemAdd)
				r.HandleFunc("/{id}", s.BudgetItemEdit) // convert
			})

			r.Route("/budgets", func(r chi.Router) {
				r.HandleFunc("/", s.BudgetsTable)
				r.HandleFunc("/add", s.BudgetAdd)
				r.HandleFunc("/{projectId}/{budgetItemId}", s.BudgetEdit) // convert
			})

			r.Route("/invoices", func(r chi.Router) {
				r.HandleFunc("/", s.InvoicesTable)
				r.HandleFunc("/{id}", s.InvoiceEdit) // convert

				r.Route("/{invoiceId}/details", func(r chi.Router) {
					r.HandleFunc("/", s.DetailsTable)              // convert
					r.HandleFunc("/add", s.DetailsAdd)             // convert
					r.HandleFunc("/{budgetItemId}", s.DetailsEdit) // convert
				})
			})

			r.Route("/categories", func(r chi.Router) {
				r.HandleFunc("/", s.CategoriesTable)
				r.HandleFunc("/add", s.CategoryAdd)
				r.HandleFunc("/{id}", s.EditCategory) // convert
			})

			r.Route("/materiales", func(r chi.Router) {
				r.HandleFunc("/", s.MaterialsTable)
				r.HandleFunc("/add", s.MaterialsAdd)
				r.HandleFunc("/{id}", s.MaterialsEdit) // convert
			})

			r.Route("/rubros", func(r chi.Router) {
				r.HandleFunc("/", s.RubrosTable)
				r.HandleFunc("/{id}", s.MaterialsByItem)
				r.HandleFunc("/{id}/material", s.MaterialByItemForm)
				r.HandleFunc("/{id}/material/{materialId}", s.MaterialItemsOperations)
			})

			r.Route("/cantidades", func(r chi.Router) {
				r.HandleFunc("/", s.CantidadesTable)
				r.HandleFunc("/add", s.CantidadesAdd)
				r.HandleFunc("/{id}", s.CantidadesEdit)
			})

			r.HandleFunc("/analisis", s.AnalysisTable)
		})
	})

	return s
}
