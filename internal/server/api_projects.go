package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllProjects(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	projects, _ := s.DB.GetAllProjects(ctx.CompanyId)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(projects)
}

func (s *Server) ApiCreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Body == http.NoBody || r.Body == nil {
		errorResponse := make(map[string]string)
		errorResponse["error"] = "Invalid request body"

		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	ctx, _ := utils.GetMyPaload(r)

	var createdProject types.Project
	var err error
	project := types.Project{}

	if err = json.NewDecoder(r.Body).Decode(&project); err != nil {
		errorResponse := make(map[string]string)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	project.CompanyId = ctx.CompanyId

	if createdProject, err = s.DB.CreateProject(project); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Ya existe un proyecto con ese nombre"
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
	_ = json.NewEncoder(w).Encode(createdProject)
}
