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

func (s *Server) MaterialsTable(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

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

	categoryId := r.Form.Get("category")
	categoryIdParsed, err := uuid.Parse(categoryId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese un valor para la Categoría"))
		return
	}

	material := types.Material{
		Code:      code,
		Name:      name,
		Unit:      unit,
		Category:  types.Category{Id: categoryIdParsed},
		CompanyId: ctxPayload.CompanyId,
	}

	if err := s.DB.CreateMaterial(material); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("El material con código %s y/o nombre %s ya existe", material.Code, material.Name)))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("MaterialsTable error", "error", err)
		w.Write([]byte(err.Error()))
		return
	}

	materials := s.DB.GetAllMaterials(ctxPayload.CompanyId)
	component := partials.MaterialsTable(materials)
	component.Render(r.Context(), w)
}

func (s *Server) MaterialsTableDisplay(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	materials := s.DB.GetAllMaterials(ctxPayload.CompanyId)

	component := partials.MaterialsTable(materials)
	component.Render(r.Context(), w)
}

func (s *Server) MaterialsAdd(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	categories, _ := s.DB.GetAllCategories(ctxPayload.CompanyId)

	categoriesSelect := []types.Select{}
	for _, c := range categories {
		categoriesSelect = append(categoriesSelect, types.Select{Key: c.Id.String(), Value: c.Name})
	}

	component := partials.EditMaterial(nil, categoriesSelect)
	component.Render(r.Context(), w)
}

func (s *Server) MaterialsEdit(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	material, _ := s.DB.GetMaterial(parsedId, ctxPayload.CompanyId)
	r.ParseForm()
	updatedMaterial := types.Material{
		Id:        parsedId,
		CompanyId: ctxPayload.CompanyId,
	}

	code := r.Form.Get("code")
	if code == "" {
		updatedMaterial.Code = material.Code
	} else {
		updatedMaterial.Code = code
	}

	name := r.Form.Get("name")
	if name == "" {
		updatedMaterial.Name = material.Name
	} else {
		updatedMaterial.Name = name
	}

	unit := r.Form.Get("unit")
	if unit == "" {
		updatedMaterial.Unit = material.Unit
	} else {
		updatedMaterial.Unit = unit
	}

	categoryId := r.Form.Get("category")
	if categoryId == "" {
		updatedMaterial.Category = material.Category
	} else {
		categoryIdParsed, err := uuid.Parse(categoryId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese un valor para la Categoría"))
			return
		}
		updatedMaterial.Category = types.Category{Id: categoryIdParsed}
	}

	if err := s.DB.UpdateMaterial(updatedMaterial); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("El material con código: %s o nombre: %s ya existe", material.Code, material.Name)))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("MaterialsTable error", "error", err)
		w.Write([]byte(err.Error()))
		return
	}

	materials := s.DB.GetAllMaterials(ctxPayload.CompanyId)

	component := partials.MaterialsTable(materials)
	component.Render(r.Context(), w)
}

func (s *Server) GetOneMaterial(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	material, _ := s.DB.GetMaterial(parsedId, ctxPayload.CompanyId)

	categories, _ := s.DB.GetAllCategories(ctxPayload.CompanyId)

	categoriesSelect := []types.Select{}
	for _, c := range categories {
		categoriesSelect = append(categoriesSelect, types.Select{Key: c.Id.String(), Value: c.Name})
	}

	component := partials.EditMaterial(&material, categoriesSelect)
	component.Render(r.Context(), w)
}
