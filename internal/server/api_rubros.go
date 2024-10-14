package server

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllRubros(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	rubros, _ := s.DB.GetAllRubros(ctx.CompanyId)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&rubros)
}
