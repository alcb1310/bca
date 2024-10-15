package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"

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

func (s *Server) ApiCreateRubros(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody || r.Body == nil {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "Invalid request body"

		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ctx, _ := utils.GetMyPaload(r)
	var rubro types.Rubro
	if err := json.NewDecoder(r.Body).Decode(&rubro); err != nil {
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	rubro.CompanyId = ctx.CompanyId

	rubroId, err := s.DB.CreateRubro(rubro)
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Rubro con ese código y/ nombre ya existe"
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(errorResponse)
			return
		}

		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	rubro.Id = rubroId
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(rubro)
}
