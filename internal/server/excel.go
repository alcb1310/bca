package server

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/excel"
	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) BalanceExcel(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := r.URL.Query().Get("project")
	parsedProjectId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Error parsing projectId", "err", err)
		return
	}
	d := r.URL.Query().Get("date")
	dateVal, _ := time.Parse("2006-01-02", d)

	f := excel.Balance(ctx.CompanyId, parsedProjectId, dateVal, s.DB)
	fName := strings.Trim(f.Path, ".")

	w.Header().Set("HX-Redirect", fName)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", "attachment; filename=cuadre.xlsx")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	if err := f.Write(w); err != nil {
		slog.Error(err.Error())
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err := os.Remove(f.Path); err != nil {
			slog.Error(err.Error())
		}
	}()
}

func (s *Server) ActualExcel(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	parsedProjectId, err := uuid.Parse(r.URL.Query().Get("proyecto"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Error parsing projectId", "Err", err)
		return
	}

	l, err := strconv.ParseUint(r.URL.Query().Get("nivel"), 10, 8)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Error parsing level", "Err", err)
		return
	}
	level := uint8(l)
	budgets, _ := s.DB.GetBudgetsByProjectId(ctx.CompanyId, parsedProjectId, &level)

	f := excel.Actual(ctx.CompanyId, parsedProjectId, budgets, nil, s.DB)
	fName := strings.Trim(f.Path, ".")

	w.Header().Set("HX-Redirect", fName)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", "attachment; filename=actual.xlsx")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	if err := f.Write(w); err != nil {
		slog.Error(err.Error())
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err := os.Remove(f.Path); err != nil {
			slog.Error(err.Error())
		}
	}()
}

func (s *Server) HistoricExcel(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	parsedProjectId, err := uuid.Parse(r.URL.Query().Get("proyecto"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Error parsing projectId", "err", err)
		return
	}

	l, err := strconv.ParseUint(r.URL.Query().Get("nivel"), 10, 8)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Error parsing level", "Err", err)
		return
	}
	level := uint8(l)
	d := r.URL.Query().Get("fecha")
	dateVal, _ := time.Parse("2006-01-02", d)

	budgets := s.DB.GetHistoricByProject(ctx.CompanyId, parsedProjectId, dateVal, level)

	f := excel.Actual(ctx.CompanyId, parsedProjectId, budgets, &dateVal, s.DB)
	fName := strings.Trim(f.Path, ".")

	w.Header().Set("HX-Redirect", fName)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", "attachment; filename=actual.xlsx")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	if err := f.Write(w); err != nil {
		slog.Error(err.Error())
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err := os.Remove(f.Path); err != nil {
			slog.Error(err.Error())
		}
	}()
}

func (s *Server) SpentExcel(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	parsedProjectId, err := uuid.Parse(r.URL.Query().Get("proyecto"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Error parsing projectId", "err", err)
		return
	}

	l, err := strconv.ParseUint(r.URL.Query().Get("nivel"), 10, 8)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("Error parsing level", "Err", err)
		return
	}
	level := uint8(l)

	d := r.URL.Query().Get("fecha")
	dateVal, _ := time.Parse("2006-01-02", d)

	budgetItems := s.DB.GetBudgetItemsByLevel(ctx.CompanyId, level)
	reportData := []types.Spent{}
	var grandTotal float64 = 0

	for _, bi := range budgetItems {
		x := []types.BudgetItem{bi}
		res := []uuid.UUID{}
		res = s.DB.GetNonAccumulateChildren(&ctx.CompanyId, &parsedProjectId, x, res)

		total := s.DB.GetSpentByBudgetItem(ctx.CompanyId, parsedProjectId, bi.ID, dateVal, res)
		grandTotal += total
		if total > 0 {
			reportData = append(reportData, types.Spent{
				Spent:      total,
				BudgetItem: bi,
			})
		}
	}

	project, _ := s.DB.GetProject(parsedProjectId, ctx.CompanyId)
	f := excel.Spent(project, reportData, dateVal)
	fName := "/" + strings.Trim(f.Path, "./public")

	w.Header().Set("HX-Redirect", fName)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", "attachment; filename=actual.xlsx")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	if err := f.Write(w); err != nil {
		slog.Error(err.Error())
	}

	go func() {
		time.Sleep(1 * time.Second)
		if err := os.Remove(f.Path); err != nil {
			slog.Error(err.Error())
		}
	}()
}
