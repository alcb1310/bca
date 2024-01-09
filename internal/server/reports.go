package server

import (
	"bca-go-final/internal/views/bca/reports"
	"net/http"
)

func (s *Server) Actual(w http.ResponseWriter, r *http.Request) {
	component := reports.ActualView()
	component.Render(r.Context(), w)
}

func (s *Server) Balance(w http.ResponseWriter, r *http.Request) {
	component := reports.BalanceView()
	component.Render(r.Context(), w)
}

func (s *Server) Historic(w http.ResponseWriter, r *http.Request) {
	component := reports.HistoricView()
	component.Render(r.Context(), w)
}

func (s *Server) Spent(w http.ResponseWriter, r *http.Request) {
	component := reports.SpentView()
	component.Render(r.Context(), w)
}
