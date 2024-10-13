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

func (s *Server) ApiGetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	categories, _ := s.DB.GetAllCategories(ctx.CompanyId)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(categories)
}

func (s *Server) ApiCreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody || r.Body == nil {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "Invalid request body"

		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ctx, _ := utils.GetMyPaload(r)

	mat := types.Category{}

	if err := json.NewDecoder(r.Body).Decode(&mat); err != nil {
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if mat.Name == "" {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "El nombre es obligatorio"
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	mat.CompanyId = ctx.CompanyId

	if err := s.DB.CreateCategory(mat); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Ya existe una categoria con ese nombre"
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
}

func (s *Server) ApiUpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)

	catToUpdate, err := s.DB.GetCategory(parsedId, ctx.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var data types.Category
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if data.Name != "" {
		catToUpdate.Name = data.Name
	}

	if err := s.DB.UpdateCategory(catToUpdate); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Ya existe una categoria con ese nombre"
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
