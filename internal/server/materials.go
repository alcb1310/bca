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

		categoryIdParsed, err := utils.ValidateUUID(r.Form.Get("category"), "categoria")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
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

	categories := s.returnAllSelects([]string{"categories"}, ctxPayload.CompanyId)["categories"]

	component := partials.EditMaterial(nil, categories)
	component.Render(r.Context(), w)
}

func (s *Server) MaterialsEdit(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	parsedId, err := utils.ValidateUUID(mux.Vars(r)["id"], "material")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	material, _ := s.DB.GetMaterial(parsedId, ctxPayload.CompanyId)

	switch r.Method {
	case http.MethodGet:
		categories := s.returnAllSelects([]string{"categories"}, ctxPayload.CompanyId)["categories"]

		component := partials.EditMaterial(&material, categories)
		component.Render(r.Context(), w)

	case http.MethodPut:
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

		categoryIdParsed, err := utils.ValidateUUID(r.Form.Get("category"), "categoria")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		material := types.Material{
			Id:        parsedId,
			Code:      code,
			Name:      name,
			Unit:      unit,
			Category:  types.Category{Id: categoryIdParsed},
			CompanyId: ctxPayload.CompanyId,
		}

		if err := s.DB.UpdateMaterial(material); err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(fmt.Sprintf("El material con código: %s o nombre: %s ya existe", material.Code, material.Name)))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			w.Write([]byte(err.Error()))
			return
		}

		materials := s.DB.GetAllMaterials(ctxPayload.CompanyId)

		component := partials.MaterialsTable(materials)
		component.Render(r.Context(), w)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
