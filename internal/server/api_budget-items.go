package server

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllBudgetItems(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	queryParams := r.URL.Query()
	search := queryParams.Get("query")

	budgetItems, _ := s.DB.GetBudgetItems(ctx.CompanyId, search)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(budgetItems)
}

