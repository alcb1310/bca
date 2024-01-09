package server

import (
	"bca-go-final/internal/views/bca/transaction"
	"net/http"
)

func (s *Server) Budget(w http.ResponseWriter, r *http.Request) {
	component := transaction.BudgetView()
	component.Render(r.Context(), w)
}

func (s *Server) Invoice(w http.ResponseWriter, r *http.Request) {
	component := transaction.InvoiceView()
	component.Render(r.Context(), w)
}

func (s *Server) Closure(w http.ResponseWriter, r *http.Request) {
	component := transaction.ClosureView()
	component.Render(r.Context(), w)
}
