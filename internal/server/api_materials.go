package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllMaterials(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	materials := s.DB.GetAllMaterials(ctx.CompanyId)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(materials)
}

func (s *Server) ApiCreateMaterial(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody || r.Body == nil {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "Invalid request body"

		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ctx, _ := utils.GetMyPaload(r)

	var mat types.Material
	if err := json.NewDecoder(r.Body).Decode(&mat); err != nil {
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	mat.CompanyId = ctx.CompanyId
	mat.Category.CompanyId = ctx.CompanyId

	if err := s.DB.CreateMaterial(mat); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Material ya existe"
			w.WriteHeader(http.StatusNotAcceptable)
			_ = json.NewEncoder(w).Encode(errorResponse)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(mat)
}

func (s *Server) ApiUpdateMaterial(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if r.Body == http.NoBody || r.Body == nil {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "Invalid request body"

		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ctx, _ := utils.GetMyPaload(r)
	matToUpdate, err := s.DB.GetMaterial(parsedId, ctx.CompanyId)
	if err != nil {
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var mat types.Material
	if err := json.NewDecoder(r.Body).Decode(&mat); err != nil {
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if mat.Code != "" {
		matToUpdate.Code = mat.Code
	}
	if mat.Name != "" {
		matToUpdate.Name = mat.Name
	}
	if mat.Unit != "" {
		matToUpdate.Unit = mat.Unit
	}
	if mat.Category.Id != uuid.Nil {
		matToUpdate.Category.Id = mat.Category.Id
	}

	if err := s.DB.UpdateMaterial(matToUpdate); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Material ya existe"
			w.WriteHeader(http.StatusNotAcceptable)
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
	_ = json.NewEncoder(w).Encode(mat)
}
