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

func (s *Server) ApiUpdateBudget(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	ctx, _ := utils.GetMyPaload(r)
	var budget types.CreateBudget

	if err := json.NewDecoder(r.Body).Decode(&budget); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errResponse := make(map[string]string)
		errResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errResponse)
	}

	slog.Info("ApiUpdateBudget", "budget", budget)

	b, err := s.DB.GetOneBudget(ctx.CompanyId, budget.ProjectId, budget.BudgetItemId)
  if err!=nil {
    w.WriteHeader(http.StatusNotFound)
    errorResponse := make(map[string]string)
    errorResponse["error"] = err.Error()
    _ = json.NewEncoder(w).Encode(errorResponse)
    return
  }
  b.CompanyId = ctx.CompanyId
  
  var budgetToUpdate = types.Budget {
    ProjectId: b.Project.ID,
    BudgetItemId: b.BudgetItem.ID,

    InitialQuantity: b.InitialQuantity,
    InitialCost: b.InitialCost,
    InitialTotal: b.InitialTotal,

    SpentQuantity: b.SpentQuantity,
    SpentTotal: b.SpentTotal,

    RemainingQuantity: b.RemainingQuantity,
    RemainingCost: b.RemainingCost,
    RemainingTotal: b.RemainingTotal,

    UpdatedBudget: b.UpdatedBudget,
    CompanyId: b.CompanyId,
  }

  if err := s.DB.UpdateBudget(budget, budgetToUpdate); err != nil {
		var e *pgconn.PgError

		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			errResponse := make(map[string]string)
			errResponse["error"] = "Ya existe un presupuesto con ese nombre"
			_ = json.NewEncoder(w).Encode(errResponse)
			return
		}
  }

  w.WriteHeader(http.StatusNoContent)
}
