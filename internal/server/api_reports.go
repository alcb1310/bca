package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

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
