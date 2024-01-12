package server

import (
	"bca-go-final/internal/views"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.Use(s.authVerify)

	r.HandleFunc("/", s.HelloWorldHandler)
	r.HandleFunc("/health", s.healthHandler)
	r.HandleFunc("/api/login", s.Login)
	r.HandleFunc("/api/register", s.Register)

	// load dummy data
	r.HandleFunc("/api/v1/load-dummy-data", s.loadDummyDataHandler)

	// users routes
	r.HandleFunc("/api/v1/users", s.AllUsers)
	r.HandleFunc("/api/v1/users/{id}", s.OneUser)

	// projects routes
	r.HandleFunc("/api/v1/projects", s.AllProjects)
	r.HandleFunc("/api/v1/projects/{id}", s.OneProject)

	// suppliers routes
	r.HandleFunc("/api/v1/suppliers", s.AllSuppliers)
	r.HandleFunc("/api/v1/suppliers/{id}", s.OneSupplier)

	// budget-items routes
	r.HandleFunc("/api/v1/budget-items", s.AllBudgetItems)
	r.HandleFunc("/api/v1/budget-items/{id}", s.OneBudgetItem)

	// invoice options
	r.HandleFunc("/api/v1/invoices", s.AllInvoices)
	r.HandleFunc("/api/v1/invoices/{id}", s.OneInvoice)

	// views

	r.HandleFunc("/login", s.LoginView)
	r.HandleFunc("/bca", s.BcaView)
	r.HandleFunc("/bca/logout", s.Logout)
	r.HandleFunc("/bca/transacciones/presupuesto", s.Budget)
	r.HandleFunc("/bca/transacciones/facturas", s.Invoice)
	r.HandleFunc("/bca/transacciones/cierre", s.Closure)

	r.HandleFunc("/bca/reportes/actual", s.Actual)
	r.HandleFunc("/bca/reportes/cuadre", s.Balance)
	r.HandleFunc("/bca/reportes/historico", s.Historic)
	r.HandleFunc("/bca/reportes/gastado", s.Spent)

	r.HandleFunc("/bca/configuracion/partidas", s.BudgetItems)
	r.HandleFunc("/bca/configuracion/proveedores", s.Suppliers)
	r.HandleFunc("/bca/configuracion/proyectos", s.Projects)

	r.HandleFunc("/bca/user/perfil", s.Profile)
	r.HandleFunc("/bca/user/admin", s.Admin)
	r.HandleFunc("/bca/user/cambio", s.ChangePassword)

	// partials

	r.HandleFunc("/bca/partials/users", s.UsersTable)
	r.HandleFunc("/bca/partials/users/add", s.UserAdd)
	r.HandleFunc("/bca/partials/users/edit/{id}", s.UserEdit)
	r.HandleFunc("/bca/partials/users/{id}", s.SingleUser)

	r.HandleFunc("/bca/partials/projects", s.ProjectsTable)
	r.HandleFunc("/bca/partials/projects/add", s.ProjectAdd)
	r.HandleFunc("/bca/partials/projects/edit/{id}", s.ProjectEditSave)
	r.HandleFunc("/bca/partials/projects/{id}", s.ProjectEdit)

	r.HandleFunc("/bca/partials/suppliers", s.SuppliersTable)
	r.HandleFunc("/bca/partials/suppliers/add", s.SupplierAdd)
	r.HandleFunc("/bca/partials/suppliers/edit/{id}", s.SuppliersEditSave)
	r.HandleFunc("/bca/partials/suppliers/{id}", s.SuppliersEdit)

	r.HandleFunc("/bca/partials/budget-item", s.BudgetItemsTable)
	r.HandleFunc("/bca/partials/budget-item/add", s.BudgetItemAdd)
	r.HandleFunc("/bca/partials/budget-item/{id}", s.BudgetItemEdit)

	r.HandleFunc("/bca/partials/budgets", s.BudgetsTable)
	r.HandleFunc("/bca/partials/budgets/add", s.BudgetAdd)

	// This should be the last route for static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	component := views.WelcomeView()
	component.Render(r.Context(), w)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.DB.Health())
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
