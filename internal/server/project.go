package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/settings/partials"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
		x := true
		p.IsActive = &x
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
	project, err := s.DB.GetProject(parsedId, ctxPayload.CompanyId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			resp["error"] = fmt.Sprintf("Project with ID: `%s` not found", id)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	switch r.Method {
	case http.MethodPut:
		var p types.Project

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		if p.Name == "" {
			p.Name = project.Name
		}

		if p.IsActive == nil {
			p.IsActive = project.IsActive
		}

		if err := s.DB.UpdateProject(p, parsedId, ctxPayload.CompanyId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)

	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(project)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) ProjectsTable(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)

	if r.Method == http.MethodPost {
		r.ParseForm()
		x := r.Form.Get("active") == "active"
		p := types.Project{
			Name:      r.Form.Get("name"),
			IsActive:  &x,
			CompanyId: ctxPayload.CompanyId,
		}
		if p.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("name cannot be empty")
			return
		}
		_, err := s.DB.CreateProject(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

	}
	projects, _ := s.DB.GetAllProjects(ctxPayload.CompanyId)
	component := partials.ProjectsTable(projects)
	component.Render(r.Context(), w)
}

func (s *Server) ProjectAdd(w http.ResponseWriter, r *http.Request) {
	component := partials.EditProject(nil)
	component.Render(r.Context(), w)
}

func (s *Server) ProjectEditSave(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]
	parsedId, _ := uuid.Parse(id)
	p, _ := s.DB.GetProject(parsedId, ctx.CompanyId)

	r.ParseForm()
	if r.Form.Get("name") != "" {
		p.Name = r.Form.Get("name")
	}
	x := r.Form.Get("active") == "active"
	p.IsActive = &x

	if err := s.DB.UpdateProject(p, parsedId, ctx.CompanyId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	projects, _ := s.DB.GetAllProjects(ctx.CompanyId)
	component := partials.ProjectsTable(projects)
	component.Render(r.Context(), w)
}

func (s *Server) ProjectEdit(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]
	parsedId, _ := uuid.Parse(id)
	p, _ := s.DB.GetProject(parsedId, ctx.CompanyId)

	component := partials.EditProject(&p)
	component.Render(r.Context(), w)
}
