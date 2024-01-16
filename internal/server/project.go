package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/settings/partials"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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
