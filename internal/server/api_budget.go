package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllBudgets(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	queryParams := r.URL.Query()
	search := queryParams.Get("query")
	project := queryParams.Get("project")

	projectId, err := uuid.Parse(project)
	if err != nil {
		if project != "" {
			slog.Info("GetAllBudgets: invalid project id", "error", err)
			w.WriteHeader(http.StatusNotAcceptable)
			errorResponse := make(map[string]string)
			errorResponse["error"] = err.Error()
			_ = json.NewEncoder(w).Encode(errorResponse)
			return
		}
    projectId = uuid.UUID{}
	}


	budgets, err := s.DB.GetBudgets(ctx.CompanyId, projectId, search)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errResponse := make(map[string]string)
		errResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(budgets)
}

func (s *Server) ApiCreateBudget(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody || r.Body != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errResponse := make(map[string]string)
		errResponse["error"] = "Invalid request body"
		_ = json.NewEncoder(w).Encode(errResponse)
	}
	slog.Info("ApiCreateBudget")
	var budget types.CreateBudget

	if err := json.NewDecoder(r.Body).Decode(&budget); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errResponse := make(map[string]string)
		errResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errResponse)
	}

	ctx, _ := utils.GetMyPaload(r)
	budget.CompanyId = ctx.CompanyId

	if _, err := s.DB.CreateBudget(&budget); err != nil {
		var e *pgconn.PgError

		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			errResponse := make(map[string]string)
			errResponse["error"] = "Ya existe un presupuesto con ese nombre"
			_ = json.NewEncoder(w).Encode(errResponse)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(budget)
}
