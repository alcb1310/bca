package server

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllMaterials(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	materials := s.DB.GetAllMaterials(ctx.CompanyId)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(materials)
}
