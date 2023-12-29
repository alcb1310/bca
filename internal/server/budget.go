package server

import (
	"bca-go-final/internal/types"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) AllBudgets(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctx, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	case http.MethodPost:
		budget := &types.CreateBudget{}
		err := json.NewDecoder(r.Body).Decode(budget)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		if budget.ProjectId == uuid.Nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "project_id cannot be empty"
			resp["field"] = "project_id"
			json.NewEncoder(w).Encode(resp)
			return
		}
		if budget.BudgetItemId == uuid.Nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "budget_item_id cannot be empty"
			resp["field"] = "budget_item_id"
			json.NewEncoder(w).Encode(resp)
			return
		}
		if budget.Quantity == nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "quantity cannot be empty"
			resp["field"] = "quantity"
			json.NewEncoder(w).Encode(resp)
			return
		}
		if budget.Cost == nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "cost cannot be empty"
			resp["field"] = "cost"
			json.NewEncoder(w).Encode(resp)
			return
		}

		budget.CompanyId = ctx.CompanyId
		b, err := s.DB.CreateBudget(budget)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(b)

	case http.MethodGet:
		budgets, err := s.DB.GetBudgets(ctx.CompanyId)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(budgets)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *Server) OneBudget(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctx, _ := getMyPaload(r)
	_ = ctx
	projectId := mux.Vars(r)["projectId"]
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	budgetItemId := mux.Vars(r)["budgetItemId"]
	budgetItemUuid, err := uuid.Parse(budgetItemId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}
	budget, err := s.DB.GetOneBudget(ctx.CompanyId, projectUuid, budgetItemUuid)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	switch r.Method {
	case http.MethodPut:
		// TODO: implement update
		w.WriteHeader(http.StatusNotImplemented)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(budget)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) AllBudgetsByProject(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctx, _ := getMyPaload(r)
	projectId := mux.Vars(r)["projectId"]
	projectUuid, err := uuid.Parse(projectId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	case http.MethodGet:
		budgets, err := s.DB.GetBudgetsByProjectId(ctx.CompanyId, projectUuid)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(budgets)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
