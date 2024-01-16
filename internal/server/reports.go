package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/reports"
	"bca-go-final/internal/views/bca/reports/partials"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (s *Server) Actual(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)

	p := s.DB.GetActiveProjects(ctx.CompanyId, true)
	projects := []types.Select{}
	for _, v := range p {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		projects = append(projects, x)
	}

	levels := s.DB.Levels(ctx.CompanyId)

	component := reports.ActualView(projects, levels)
	component.Render(r.Context(), w)
}

func (s *Server) Balance(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		pId := r.Form.Get("project")
		parsedProjectId, _ := uuid.Parse(pId)
		d := r.Form.Get("date")
		date, _ := time.Parse("2006-01-02", d)

		invoices := s.DB.GetBalance(ctx.CompanyId, parsedProjectId, date)

		component := partials.BalanceView(invoices)
		component.Render(r.Context(), w)

	case http.MethodGet:
		p, _ := s.DB.GetAllProjects(ctx.CompanyId)
		projects := []types.Select{}
		for _, v := range p {
			x := types.Select{
				Key:   v.ID.String(),
				Value: v.Name,
			}
			projects = append(projects, x)
		}

		component := reports.BalanceView(projects)
		component.Render(r.Context(), w)
	}

}

func (s *Server) Historic(w http.ResponseWriter, r *http.Request) {
	component := reports.HistoricView()
	component.Render(r.Context(), w)
}

func (s *Server) Spent(w http.ResponseWriter, r *http.Request) {
	component := reports.SpentView()
	component.Render(r.Context(), w)
}

func (s *Server) ActualGenerate(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	r.ParseForm()
	p := r.Form.Get("proyecto")
	projectId, _ := uuid.Parse(p)
	z := r.Form.Get("nivel")
	var l uint64
	var err error
	if z == "" {
		l = 0
	} else {
		l, err = strconv.ParseUint(z, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	level := uint8(l)

	budgets, err := s.DB.GetBudgetsByProjectId(ctx.CompanyId, projectId, &level)
	component := partials.BudgetView(budgets)
	component.Render(r.Context(), w)
}
