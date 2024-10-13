package server

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	categories, _ := s.DB.GetAllCategories(ctx.CompanyId)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(categories)
}
