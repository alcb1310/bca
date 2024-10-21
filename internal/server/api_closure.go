package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiPostClosure(w http.ResponseWriter, r *http.Request) {
	var data types.Closure

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	ctx, _ := utils.GetMyPaload(r)

	if err := s.DB.CreateClosure(ctx.CompanyId, data.ProjectId, data.Date); err != nil {
		if strings.Contains(err.Error(), "Ya existe un cierre para el proyecto:") {
			w.WriteHeader(http.StatusConflict)
			errorResponse := make(map[string]string)
			errorResponse["error"] = err.Error()
			_ = json.NewEncoder(w).Encode(errorResponse)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	w.WriteHeader(http.StatusOK)
}
