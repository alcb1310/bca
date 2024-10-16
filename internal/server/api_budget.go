package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllBudgets(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	queryParams := r.URL.Query()
	search := queryParams.Get("query")

	budgets, err := s.DB.GetBudgets(ctx.CompanyId, uuid.UUID{}, search)
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
