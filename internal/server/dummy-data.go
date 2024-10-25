package server

import (
	"encoding/json"
	"net/http"

	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) loadDummyDataHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctx, _ := utils.GetMyPaload(r)
	companyId := ctx.CompanyId

	err := s.DB.LoadDummyData(companyId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
