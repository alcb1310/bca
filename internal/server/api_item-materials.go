package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

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
	var rubroId uuid.UUID
	var err error
	id := chi.URLParam(r, "id")
	if rubroId, err = uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
	}

	if r.Body == http.NoBody || r.Body == nil {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "Invalid request body"

		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var material types.ItemMaterialType
	ctx, _ := utils.GetMyPaload(r)

	if err := json.NewDecoder(r.Body).Decode(&material); err != nil {
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	material.CompanyId = ctx.CompanyId
	if rubroId != material.ItemId {
		w.WriteHeader(http.StatusTeapot)
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["error"] = "Invalid request"
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err := s.DB.AddMaterialsByItem(material.ItemId, material.MaterialId, material.Quantity, ctx.CompanyId); err != nil {
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

	w.WriteHeader(http.StatusCreated)
}
