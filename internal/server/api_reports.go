package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiActualReport(w http.ResponseWriter, r *http.Request) {
	project_id := r.URL.Query().Get("project_id")
	l := r.URL.Query().Get("level")

	if project_id == "" || l == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	parsedProjectId, err := uuid.Parse(project_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	ctx, _ := utils.GetMyPaload(r)

	le, err := strconv.Atoi(l)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	level := uint8(le)
	budgets, _ := s.DB.GetBudgetsByProjectId(ctx.CompanyId, parsedProjectId, &level)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(budgets)
}

func (s *Server) ApiLevels(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	data := s.DB.Levels(ctx.CompanyId)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

func (s *Server) ApiHistoricReport(w http.ResponseWriter, r *http.Request) {
	project_id := r.URL.Query().Get("project_id")
	l := r.URL.Query().Get("level")
	d := r.URL.Query().Get("date")

	if project_id == "" || l == "" || d == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parsedProjectId, err := uuid.Parse(project_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	le, err := strconv.Atoi(l)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	level := uint8(le)

	layout := "2006-01-02"
	selectedDate, err := time.Parse(layout, d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	ctx, _ := utils.GetMyPaload(r)
	hitoricBudgets := s.DB.GetHistoricByProject(ctx.CompanyId, parsedProjectId, selectedDate, level)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(hitoricBudgets)
}

func (s *Server) ApiBalanceReport(w http.ResponseWriter, r *http.Request) {
	var balanceReport types.BalanceResponse
	balanceReport.Invoices = []types.InvoiceResponse{}
	balanceReport.Total = 0

	project_id := r.URL.Query().Get("project_id")
	d := r.URL.Query().Get("date")
	if project_id == "" || d == "" {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(balanceReport)
		return
	}

	parsedProjectId, err := uuid.Parse(project_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	layout := "2006-01-02"
	selectedDate, err := time.Parse(layout, d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	ctx, _ := utils.GetMyPaload(r)

	balanceReport = s.DB.GetBalance(ctx.CompanyId, parsedProjectId, selectedDate)
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(balanceReport)
}

func (s *Server) ApiUpdateBalanceReport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		responseError := make(map[string]string)
		responseError["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(responseError)
		return
	}
	ctx, _ := utils.GetMyPaload(r)

	i, err := s.DB.GetOneInvoice(parsedId, ctx.CompanyId)
	if err != nil {
		responseError := make(map[string]string)
		responseError["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(responseError)
		return
	}

	if err := s.DB.BalanceInvoice(i); err != nil {
		responseError := make(map[string]string)
		responseError["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(responseError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) ApiSpentReport(w http.ResponseWriter, r *http.Request) {
	project_id := r.URL.Query().Get("project_id")
	l := r.URL.Query().Get("level")
	d := r.URL.Query().Get("date")

	if project_id == "" || l == "" || d == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parsedProjectId, err := uuid.Parse(project_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	le, err := strconv.Atoi(l)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	level := uint8(le)

	layout := "2006-01-02"
	selectedDate, err := time.Parse(layout, d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	ctx, _ := utils.GetMyPaload(r)

	budgetItems := s.DB.GetBudgetItemsByLevel(ctx.CompanyId, level)
	reportData := []types.Spent{}
	var grandTotal float64 = 0

	for _, bi := range budgetItems {
		x := []types.BudgetItem{bi}
		res := []uuid.UUID{}
		res = s.DB.GetNonAccumulateChildren(&ctx.CompanyId, &parsedProjectId, x, res)

		total := s.DB.GetSpentByBudgetItem(ctx.CompanyId, parsedProjectId, bi.ID, selectedDate, res)
		if total != 0 {
			grandTotal += total

			reportData = append(reportData, types.Spent{
				Spent:      total,
				BudgetItem: bi,
			})
		}
	}

	responseData := types.SpentResponse{
		Spent:   reportData,
		Total:   grandTotal,
		Project: parsedProjectId,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(responseData)
}

func (s *Server) ApiSpentByBudgetItem(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	projectId := chi.URLParam(r, "projectId")
	parsedProjectId, err := uuid.Parse(projectId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	budgetItemId := chi.URLParam(r, "budgetItemId")
	parsedBudgetItemId, err := uuid.Parse(budgetItemId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	d := chi.URLParam(r, "date")
	date, _ := time.Parse("2006-01-02", d)

	budgetItem, err := s.DB.GetOneBudgetItem(parsedBudgetItemId, ctx.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		slog.Error(err.Error())
		return
	}

	x := []types.BudgetItem{*budgetItem}
	var res []uuid.UUID
	res = s.DB.GetNonAccumulateChildren(&ctx.CompanyId, &parsedProjectId, x, res)

	spent := s.DB.GetDetailsByBudgetItem(ctx.CompanyId, parsedProjectId, parsedBudgetItemId, date, res)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(spent)
}
