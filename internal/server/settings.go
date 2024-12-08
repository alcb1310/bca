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
	_ = component.Render(r.Context(), w)
}

func (s *Server) Suppliers(w http.ResponseWriter, r *http.Request) {
	component := settings.SupplierView()
	_ = component.Render(r.Context(), w)
}

func (s *Server) Projects(w http.ResponseWriter, r *http.Request) {
	component := settings.ProjectView()
	_ = component.Render(r.Context(), w)
}

func (s *Server) Categories(w http.ResponseWriter, r *http.Request) {
	component := settings.CategoryView()
	_ = component.Render(r.Context(), w)
}

func (s *Server) Materiales(w http.ResponseWriter, r *http.Request) {
	component := settings.MaterialsView()
	_ = component.Render(r.Context(), w)
}

func (s *Server) Rubros(w http.ResponseWriter, r *http.Request) {
	component := settings.RubrosView()
	_ = component.Render(r.Context(), w)
}

func (s *Server) RubrosAdd(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	redirectURL := "/bca/configuracion/rubros/crear"
	var rubro types.Rubro

	id := r.URL.Query().Get("id")

	if id != "" {
		parsedId, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rubro, err = s.DB.GetOneRubro(parsedId, ctxPayload.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		redirectURL = fmt.Sprintf("%s?id=%s", redirectURL, id)
	}

	_ = r.ParseForm()

	code := r.Form.Get("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Ingrese un valor para el Código"))
		return
	}

	name := r.Form.Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Ingrese un valor para el Nombre"))
		return
	}

	unit := r.Form.Get("unit")
	if unit == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Ingrese un valor para la Unidad"))
		return
	}

	rubro = types.Rubro{
		Code:      code,
		Name:      name,
		Unit:      unit,
		CompanyId: ctxPayload.CompanyId,
	}

	rubroId, err := s.DB.CreateRubro(rubro)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(fmt.Sprintf("El rubro con código %s y/o nombre %s ya existe", rubro.Code, rubro.Name)))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	redirectURL = fmt.Sprintf("%s?id=%s", redirectURL, rubroId)
	rubro.Id = rubroId

	component := partials.EditRubros(&rubro)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
	_ = component.Render(r.Context(), w)
}

func (s *Server) RubrosAddForm(w http.ResponseWriter, r *http.Request) {
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

	component := partials.EditRubros(rubro)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
	_ = component.Render(r.Context(), w)
}

func (s *Server) RubrosEdit(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	redirectURL := "/bca/configuracion/rubros/crear"
	var rubro types.Rubro

	id := r.URL.Query().Get("id")

	if id != "" {
		parsedId, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rubro, err = s.DB.GetOneRubro(parsedId, ctxPayload.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		redirectURL = fmt.Sprintf("%s?id=%s", redirectURL, id)
	}

	_ =r.ParseForm()
	code := r.Form.Get("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Ingrese un valor para el Código"))
		return
	}

	name := r.Form.Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Ingrese un valor para el Nombre"))
		return
	}

	unit := r.Form.Get("unit")
	if unit == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Ingrese un valor para la Unidad"))
		return
	}

	rubro = types.Rubro{
		Code:      code,
		Name:      name,
		Unit:      unit,
		CompanyId: ctxPayload.CompanyId,
	}

	rubro.Id, _ = uuid.Parse(id)

	if err := s.DB.UpdateRubro(rubro); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(fmt.Sprintf("El Código %s ya existe", rubro.Code)))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	component := partials.EditRubros(&rubro)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
	_ = component.Render(r.Context(), w)
}
