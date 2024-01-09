package server

import (
	"bca-go-final/internal/views/bca/settings"
	"net/http"
)

func (s *Server) BudgetItems(w http.ResponseWriter, r *http.Request) {
	component := settings.BudgetItems()
	component.Render(r.Context(), w)

}

func (s *Server) Suppliers(w http.ResponseWriter, r *http.Request) {
	component := settings.SupplierView()
	component.Render(r.Context(), w)
}

func (s *Server) Projects(w http.ResponseWriter, r *http.Request) {
	component := settings.ProjectView()
	component.Render(r.Context(), w)
}
