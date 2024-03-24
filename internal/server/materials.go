package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/settings/partials"
)

func (s *Server) MaterialsTable(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	if r.Method == "POST" {
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
				w.Write([]byte(fmt.Sprintf("El Código %s ya existe", material.Code)))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			w.Write([]byte(err.Error()))
			return
		}

		_ = categoryIdParsed

	}

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
