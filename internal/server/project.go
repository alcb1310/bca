package server

import (
	"bca-go-final/internal/types"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) AllProjects(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctxPayload, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodPost:
		var p types.Project

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		if p.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "name cannot be empty"
			resp["field"] = "name"
			json.NewEncoder(w).Encode(resp)
			return
		}
		p.IsActive = true
		p.CompanyId = ctxPayload.CompanyId

		project, err := s.DB.CreateProject(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(project)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	case http.MethodGet:
		projects, err := s.DB.GetAllProjects(ctxPayload.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(projects)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) OneProject(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctxPayload, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]

	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	switch r.Method {
	case http.MethodGet:
		project, err := s.DB.GetProject(parsedId, ctxPayload.CompanyId)

		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				resp["error"] = fmt.Sprintf("Project with ID: `%s` not found", id)
				json.NewEncoder(w).Encode(resp)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(project)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
