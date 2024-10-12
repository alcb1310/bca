package server

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllBudgetItems(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	queryParams := r.URL.Query()
	search := queryParams.Get("query")

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
