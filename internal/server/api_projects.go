package server

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllProjects(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	projects, _ := s.DB.GetAllProjects(ctx.CompanyId)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(projects)
}
