package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
	"github.com/alcb1310/bca/internal/views/bca/settings/partials"
)

func (s *Server) CategoriesTable(w http.ResponseWriter, r *http.Request) {
	var n string
	ctxPayload, _ := utils.GetMyPaload(r)

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if n = r.Form.Get("name"); n == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Ingrese un nombre de categoría"))
		return
	}

	c := types.Category{
		Name:      n,
		CompanyId: ctxPayload.CompanyId,
	}
	if err := s.DB.CreateCategory(c); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(fmt.Sprintf("La categoría %s ya existe", c.Name)))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("CategoriesTable error", "error", err)
		return
	}

	categories, _ := s.DB.GetAllCategories(ctxPayload.CompanyId)
	component := partials.CategoriesTable(categories)

	_ = component.Render(r.Context(), w)
}

func (s *Server) CategoriesTableDisplay(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	categories, _ := s.DB.GetAllCategories(ctxPayload.CompanyId)
	component := partials.CategoriesTable(categories)

	_ = component.Render(r.Context(), w)
}

func (s *Server) CategoryAdd(w http.ResponseWriter, r *http.Request) {
	component := partials.EditCategory(nil)
	_ = component.Render(r.Context(), w)
}

func (s *Server) EditCategory(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)
	c, err := s.DB.GetCategory(parsedId, ctxPayload.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Categoría no encontrada"))
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	cat := types.Category{
		Id:        parsedId,
		CompanyId: ctxPayload.CompanyId,
	}

	n := r.Form.Get("name")
	if n == "" {
		cat.Name = c.Name
	} else {
		cat.Name = n
	}

	if err := s.DB.UpdateCategory(cat); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(fmt.Sprintf("La categoria %s ya existe", cat.Name)))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("EditCategory error", "error", err)
		return
	}

	categories, _ := s.DB.GetAllCategories(ctxPayload.CompanyId)
	component := partials.CategoriesTable(categories)

	_ = component.Render(r.Context(), w)
}

func (s *Server) GetOneCategory(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)
	c, err := s.DB.GetCategory(parsedId, ctxPayload.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Categoría no encontrada"))
		return
	}
	component := partials.EditCategory(&c)
	_ = component.Render(r.Context(), w)
}
