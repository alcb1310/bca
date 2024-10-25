package server

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
	"github.com/alcb1310/bca/internal/views/bca/reports"
	"github.com/alcb1310/bca/internal/views/bca/reports/partials"
)

func (s *Server) Actual(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

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
	_ = component.Render(r.Context(), w)
}

func (s *Server) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	p, _ := s.DB.GetAllProjects(ctx.CompanyId, "")
	projects := []types.Select{}
	for _, v := range p {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		projects = append(projects, x)
	}

	component := reports.BalanceView(projects)
	_ = component.Render(r.Context(), w)
}

func (s *Server) RetreiveBalance(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	_ = r.ParseForm()
	pId := r.Form.Get("project")
	parsedProjectId, _ := uuid.Parse(pId)
	d := r.Form.Get("date")
	date, _ := time.Parse("2006-01-02", d)

	invoices := s.DB.GetBalance(ctx.CompanyId, parsedProjectId, date)

	component := partials.BalanceView(invoices)
	_ = component.Render(r.Context(), w)
}

func (s *Server) Historic(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

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
		_ = component.Render(r.Context(), w)
		return
	}

	component := reports.HistoricView(projects, levels)
	_ = component.Render(r.Context(), w)
}

func (s *Server) Spent(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
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
		_ = component.Render(nCtx, w)
		return
	}

	levels := s.DB.Levels(ctx.CompanyId)
	component := reports.SpentView(projects, levels)
	_ = component.Render(r.Context(), w)
}

func (s *Server) ActualGenerate(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	_ = r.ParseForm()
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
			slog.Error(err.Error())
			return
		}
	}
	level := uint8(l)

	budgets, err := s.DB.GetBudgetsByProjectId(ctx.CompanyId, projectId, &level)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	component := partials.BudgetView(budgets)
	_ = component.Render(r.Context(), w)
}

func (s *Server) SpentByBudgetItem(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "budgetItemId")
	budgetItemId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	pId := chi.URLParam(r, "projectId")
	parsedProjectId, err := uuid.Parse(pId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	d := chi.URLParam(r, "date")
	date, _ := time.Parse("2006-01-02", d)

	budgetItem, _ := s.DB.GetOneBudgetItem(budgetItemId, ctx.CompanyId)

	x := []types.BudgetItem{*budgetItem}
	var res []uuid.UUID
	res = s.DB.GetNonAccumulateChildren(&ctx.CompanyId, &parsedProjectId, x, res)

	spent := s.DB.GetDetailsByBudgetItem(ctx.CompanyId, parsedProjectId, budgetItemId, date, res)

	component := partials.SpentDetails(spent, *budgetItem)
	_ = component.Render(r.Context(), w)
}
