package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllRubros(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	rubros, _ := s.DB.GetAllRubros(ctx.CompanyId)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&rubros)
}

func (s *Server) ApiGetRubro(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	rubroId := chi.URLParam(r, "id")
	if strings.ToLower(rubroId) == "crear" {
		selectedRubro := types.Rubro{
			Code:      "",
			Name:      "",
			Unit:      "",
			CompanyId: ctx.CompanyId,
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&selectedRubro)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}
