package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/settings/partials"
)

func (s *Server) ProjectsTable(w http.ResponseWriter, r *http.Request) {
	var err error
	ctxPayload, _ := utils.GetMyPaload(r)

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
			w.Write([]byte("Ingrese un valor para el nombre"))
			return
		}
		if r.Form.Get("gross_area") != "" {
			p.GrossArea, err = utils.ConvertFloat(r.Form.Get("gross_area"), "área bruta", false)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
		}
		if r.Form.Get("net_area") != "" {
			p.NetArea, err = utils.ConvertFloat(r.Form.Get("net_area"), "área útil", false)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
		}
		_, err := s.DB.CreateProject(p)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(fmt.Sprintf("El nombre %s ya existe", p.Name)))
				return
			}
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
	var err error
	ctx, _ := utils.GetMyPaload(r)
	parsedId, _ := utils.ValidateUUID(mux.Vars(r)["id"], "proyecto")
	p, _ := s.DB.GetProject(parsedId, ctx.CompanyId)

	r.ParseForm()
	if r.Form.Get("name") != "" {
		p.Name = r.Form.Get("name")
	}
	x := r.Form.Get("active") == "active"
	p.IsActive = &x

	if r.Form.Get("gross_area") != "" {
		p.GrossArea, err = utils.ConvertFloat(r.Form.Get("gross_area"), "área bruta", false)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	if r.Form.Get("net_area") != "" {
		p.NetArea, err = utils.ConvertFloat(r.Form.Get("net_area"), "área útil", false)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
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
		log.Println(err)
		return
	}

	projects, _ := s.DB.GetAllProjects(ctx.CompanyId)
	component := partials.ProjectsTable(projects)
	component.Render(r.Context(), w)
}

func (s *Server) ProjectEdit(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	parsedId, _ := utils.ValidateUUID(mux.Vars(r)["id"], "proyecto")
	p, _ := s.DB.GetProject(parsedId, ctx.CompanyId)

	component := partials.EditProject(&p)
	component.Render(r.Context(), w)
}
