package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllItemsMaterials(w http.ResponseWriter, r *http.Request) {
	var rubroId uuid.UUID
	var err error
	ctx, _ := utils.GetMyPaload(r)

	id := chi.URLParam(r, "id")
	if strings.ToLower(id) == "crear" {
		selected := types.ACU{}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(selected)
		return
	}

	if rubroId, err = uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
	}

	rubroMaterials := s.DB.GetMaterialsByItem(rubroId, ctx.CompanyId)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(rubroMaterials)
}

func (s *Server) ApiCreateItemsMaterials(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	w.WriteHeader(http.StatusNotImplemented)
}
