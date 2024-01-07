package server

import (
	"bca-go-final/internal/views"
	"bca-go-final/internal/views/base"
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

	// budget routes
	r.HandleFunc("/api/v1/budgets", s.AllBudgets)
	r.HandleFunc("/api/v1/budgets/{projectId}", s.AllBudgetsByProject)
	r.HandleFunc("/api/v1/budgets/{projectId}/{budgetItemId}", s.OneBudget)

	// invoice options
	r.HandleFunc("/api/v1/invoices", s.AllInvoices)
	r.HandleFunc("/api/v1/invoices/{id}", s.OneInvoice)

	// views

	r.HandleFunc("/login", s.LoginView)

	// This should be the last route for static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	component := views.WelcomeView()
	base := base.Layout(component)
	base.Render(r.Context(), w)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.DB.Health())
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
