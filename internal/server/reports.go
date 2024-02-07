package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/reports"
	"bca-go-final/internal/views/bca/reports/partials"
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

	if r.URL.Query().Get("proyecto") != "" && r.URL.Query().Get("fecha") != "" && r.URL.Query().Get("nivel") != "" {
		pId := r.URL.Query().Get("proyecto")
		parsedProjectId, _ := uuid.Parse(pId)
		d := r.URL.Query().Get("fecha")
		date, _ := time.Parse("2006-01-02", d)
		l, _ := strconv.ParseUint(r.URL.Query().Get("nivel"), 10, 64)
		nivel := uint8(l)

		budgets := s.DB.GetHistoricByProject(ctx.CompanyId, parsedProjectId, date, nivel)
		component := partials.BudgetView(budgets)
		component.Render(r.Context(), w)
		return
	}

	component := reports.HistoricView(projects, levels)
	component.Render(r.Context(), w)
}

func (s *Server) Spent(w http.ResponseWriter, r *http.Request) {
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

	if r.URL.Query().Get("proyecto") != "" && r.URL.Query().Get("fecha") != "" && r.URL.Query().Get("nivel") != "" {
		pId := r.URL.Query().Get("proyecto")
		parsedProjectId, _ := uuid.Parse(pId)
		d := r.URL.Query().Get("fecha")
		date, _ := time.Parse("2006-01-02", d)
		l, _ := strconv.ParseUint(r.URL.Query().Get("nivel"), 10, 64)
		nivel := uint8(l)

		budgetItems := s.DB.GetBudgetItemsByLevel(ctx.CompanyId, nivel)
		reportData := []types.Spent{}
		var grandTotal float64 = 0

		for _, bi := range budgetItems {
			x := []types.BudgetItem{bi}
			res := []uuid.UUID{}
			res = s.DB.GetNonAccumulateChildren(&ctx.CompanyId, &parsedProjectId, x, res)

			total := s.DB.GetSpentByBudgetItem(ctx.CompanyId, parsedProjectId, bi.ID, date, res)
			grandTotal += total
			if total > 0 {
				reportData = append(reportData, types.Spent{
					Spent:      total,
					BudgetItem: bi,
				})
			}
		}

		nCtx := context.WithValue(r.Context(), "date", utils.ConvertDate(date))

		component := partials.SpentView(types.SpentResponse{
			Spent:   reportData,
			Total:   grandTotal,
			Project: parsedProjectId,
		})
		component.Render(nCtx, w)
		return
	}

	levels := s.DB.Levels(ctx.CompanyId)
	component := reports.SpentView(projects, levels)
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

func (s *Server) SpentByBudgetItem(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	id := mux.Vars(r)["budgetItemId"]
	budgetItemId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	pId := mux.Vars(r)["projectId"]
	parsedProjectId, err := uuid.Parse(pId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	d := mux.Vars(r)["date"]
	date, _ := time.Parse("2006-01-02", d)

	budgetItem, _ := s.DB.GetOneBudgetItem(budgetItemId, ctx.CompanyId)

	x := []types.BudgetItem{*budgetItem}
	var res []uuid.UUID
	res = s.DB.GetNonAccumulateChildren(&ctx.CompanyId, &parsedProjectId, x, res)

	spent := s.DB.GetDetailsByBudgetItem(ctx.CompanyId, parsedProjectId, budgetItemId, date, res)

	component := partials.SpentDetails(spent, *budgetItem)
	component.Render(r.Context(), w)
}
