package server

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllBudgetItems(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	queryParams := r.URL.Query()
	search := queryParams.Get("query")
	accumulate := queryParams.Get("accumulate")
	slog.Info("ApiGetAllBudgetItems", "query", search, "accumulate", accumulate)
	if accumulate != "" {
		acc := accumulate == "true"

		budgetItems := s.DB.GetBudgetItemsByAccumulate(ctx.CompanyId, acc)

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(budgetItems)
		return
	}

	budgetItems, _ := s.DB.GetBudgetItems(ctx.CompanyId, search)
	var response []types.BudgetItemJsonResponse

	for _, bi := range budgetItems {
		var parent *types.BudgetItem
		if !bi.ParentId.Valid {
			parent = nil
		} else {
			parent = &types.BudgetItem{
				ID:        bi.ParentId.UUID,
				Code:      bi.ParentCode.String,
				Name:      bi.ParentName.String,
				CompanyId: bi.CompanyId,
			}
		}

		var acc bool

		if bi.Accumulate.Bool {
			acc = true
		} else {
			acc = false
		}

		x := types.BudgetItemJsonResponse{
			ID:         bi.ID,
			Code:       bi.Code,
			Name:       bi.Name,
			Level:      bi.Level,
			Accumulate: acc,
			Parent:     parent,
			CompanyId:  bi.CompanyId,
		}
		response = append(response, x)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (s *Server) ApiCreateBudgetItem(w http.ResponseWriter, r *http.Request) {
	// TODO: implemente method
	if r.Body == http.NoBody || r.Body == nil {
		slog.Info("ApiCreateBudgetItem", "body", "No body received")
		w.WriteHeader(http.StatusNotAcceptable)
		errorResponse := make(map[string]string)
		errorResponse["error"] = "Invalid request body"

		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ctx, _ := utils.GetMyPaload(r)
	var data types.BudgetItemCreate
	biToCreate := types.BudgetItem{
		CompanyId: ctx.CompanyId,
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		slog.Error("ApiCreateBudgetItem", "body", r.Body, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	slog.Info("ApiCreateBudgetItem", "data", data)

	errorReseponse := make(map[string]string)

	if data.Code == "" {
		errorReseponse["code"] = "El código es obligatorio"
	}

	if data.Name == "" {
		errorReseponse["name"] = "El nombre es obligatorio"
	}

	if len(errorReseponse) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorReseponse)
		return
	}

	biToCreate.Code = data.Code
	biToCreate.Name = data.Name
	biToCreate.Accumulate = sql.NullBool{Bool: data.Accumulate, Valid: true}

	if err := s.DB.CreateBudgetItem(&biToCreate); err != nil {
		slog.Error("ApiCreateBudgetItem", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(biToCreate)
}
