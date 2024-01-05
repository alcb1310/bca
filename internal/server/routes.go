package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.Use(s.authVerify)

	r.HandleFunc("/", s.HelloWorldHandler)
	r.HandleFunc("/health", s.healthHandler)
	r.HandleFunc("/login", s.Login)
	r.HandleFunc("/register", s.Register)

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

	// invoice details
	r.HandleFunc("/api/v1/invoices/{id}/details", s.AllInvoiceDetails)
	r.HandleFunc("/api/v1/invoices/{id}/details/{detailId}", s.OneInvoiceDetails)

	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "x-access-token"}),
		handlers.AllowCredentials(),
	)(r)
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.DB.Health())
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
