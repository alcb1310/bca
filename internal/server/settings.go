package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
	"github.com/alcb1310/bca/internal/views/bca/settings"
	"github.com/alcb1310/bca/internal/views/bca/settings/partials"
)

func (s *Server) BudgetItems(w http.ResponseWriter, r *http.Request) {
	component := settings.BudgetItems()
	component.Render(r.Context(), w)
}

func (s *Server) Suppliers(w http.ResponseWriter, r *http.Request) {
	component := settings.SupplierView()
	component.Render(r.Context(), w)
}

func (s *Server) Projects(w http.ResponseWriter, r *http.Request) {
	component := settings.ProjectView()
	component.Render(r.Context(), w)
}

func (s *Server) Categories(w http.ResponseWriter, r *http.Request) {
	component := settings.CategoryView()
	component.Render(r.Context(), w)
}

func (s *Server) Materiales(w http.ResponseWriter, r *http.Request) {
	component := settings.MaterialsView()
	component.Render(r.Context(), w)
}

func (s *Server) Rubros(w http.ResponseWriter, r *http.Request) {
	component := settings.RubrosView()
	component.Render(r.Context(), w)
}

func (s *Server) RubrosAdd(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	redirectURL := "/bca/configuracion/rubros/crear"
	var rubro *types.Rubro = nil

	id := r.URL.Query().Get("id")

	if id != "" {
		parsedId, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		r, err := s.DB.GetOneRubro(parsedId, ctxPayload.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		rubro = &r
		redirectURL = fmt.Sprintf("%s?id=%s", redirectURL, id)
	}

	if r.Method == "POST" || r.Method == "PUT" {
		r.ParseForm()

		code := r.Form.Get("code")
		if code == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese un valor para el Código"))
			return
		}

		name := r.Form.Get("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese un valor para el Nombre"))
			return
		}

		unit := r.Form.Get("unit")
		if unit == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese un valor para la Unidad"))
			return
		}

		rubro = &types.Rubro{
			Code:      code,
			Name:      name,
			Unit:      unit,
			CompanyId: ctxPayload.CompanyId,
		}
	}

	if r.Method == "POST" {
		id, err := s.DB.CreateRubro(*rubro)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(fmt.Sprintf("El rubro con código %s y/o nombre %s ya existe", rubro.Code, rubro.Name)))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			slog.Error(err.Error())
			w.Write([]byte(err.Error()))
			return
		}

		redirectURL = fmt.Sprintf("%s?id=%s", redirectURL, id)
		rubro.Id = id

	} else if r.Method == "PUT" {
		rubro.Id, _ = uuid.Parse(id)

		if err := s.DB.UpdateRubro(*rubro); err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(fmt.Sprintf("El Código %s ya existe", rubro.Code)))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			slog.Error(err.Error())
			w.Write([]byte(err.Error()))
			return
		}
	}

	component := partials.EditRubros(rubro)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
	component.Render(r.Context(), w)
}
