package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
	"github.com/alcb1310/bca/internal/views/bca/settings/partials"
)

func (s *Server) ProjectsTable(w http.ResponseWriter, r *http.Request) {
	var err error
	ctxPayload, _ := utils.GetMyPaload(r)

	r.ParseForm()
	x := r.Form.Get("active") == "active"
	p := types.Project{
		Name:      r.Form.Get("name"),
		IsActive:  &x,
		CompanyId: ctxPayload.CompanyId,
	}
	if p.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese un valor para el nombre"))
		return
	}
	if r.Form.Get("gross_area") != "" {
		p.GrossArea, err = strconv.ParseFloat(r.Form.Get("gross_area"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("El área bruta debe ser un número válido"))
			return
		}
	}
	if r.Form.Get("net_area") != "" {
		p.NetArea, err = strconv.ParseFloat(r.Form.Get("net_area"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("El área neta debe ser un número válido"))
			return
		}
	}
	_, err = s.DB.CreateProject(p)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("El nombre %s ya existe", p.Name)))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("CreateProject error", "error", err)
		return
	}

	projects, _ := s.DB.GetAllProjects(ctxPayload.CompanyId)
	component := partials.ProjectsTable(projects)
	component.Render(r.Context(), w)
}

func (s *Server) ProjectsTableDisplay(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	projects, _ := s.DB.GetAllProjects(ctxPayload.CompanyId)
	component := partials.ProjectsTable(projects)
	component.Render(r.Context(), w)
}

func (s *Server) ProjectAdd(w http.ResponseWriter, r *http.Request) {
	component := partials.EditProject(nil)
	component.Render(r.Context(), w)
}

func (s *Server) ProjectEditSave(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)
	p, _ := s.DB.GetProject(parsedId, ctx.CompanyId)

	r.ParseForm()
	if r.Form.Get("name") != "" {
		p.Name = r.Form.Get("name")
	}
	x := r.Form.Get("active") == "active"
	p.IsActive = &x

	if r.Form.Get("gross_area") != "" {
		p.GrossArea, err = strconv.ParseFloat(r.Form.Get("gross_area"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("El área bruta debe ser un número válido"))
			return
		}
	}

	if r.Form.Get("net_area") != "" {
		p.NetArea, err = strconv.ParseFloat(r.Form.Get("net_area"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("El área neta debe ser un número válido"))
			return
		}
	}

	if err := s.DB.UpdateProject(p, parsedId, ctx.CompanyId); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("El nombre %s ya existe", p.Name)))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("UpdateProject error", "error", err)
		return
	}

	projects, _ := s.DB.GetAllProjects(ctx.CompanyId)
	component := partials.ProjectsTable(projects)
	component.Render(r.Context(), w)
}

func (s *Server) ProyectDisplay(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)
	p, err := s.DB.GetProject(parsedId, ctx.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Proyecto no encontrado"))
		return
	}

	component := partials.EditProject(&p)
	component.Render(r.Context(), w)
}
