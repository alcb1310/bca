package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/settings/partials"
)

func (s *Server) RubrosTable(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	rubros, _ := s.DB.GetAllRubros(ctxPayload.CompanyId)

	component := partials.RubrosTable(rubros)
	component.Render(r.Context(), w)
}

func (s *Server) MaterialsByItem(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	parsedId, err := utils.ValidateUUID(mux.Vars(r)["id"], "rubro")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acus := s.DB.GetMaterialsByItem(parsedId, ctxPayload.CompanyId)

	w.WriteHeader(http.StatusOK)
	component := partials.MaterialsItemsTable(acus)
	component.Render(r.Context(), w)
}

func (s *Server) MaterialByItemForm(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	parsedId, err := utils.ValidateUUID(mux.Vars(r)["id"], "material")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		r.ParseForm()

		parsedMaterialId, err := utils.ValidateUUID(r.Form.Get("material"), "material")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		quantityText := r.Form.Get("quantity")
		quantity, err := utils.ConvertFloat(quantityText, "cantidad", true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if err := s.DB.AddMaterialsByItem(parsedId, parsedMaterialId, quantity, ctxPayload.CompanyId); err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("Ya existe un material con ese Código"))
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		acus := s.DB.GetMaterialsByItem(parsedId, ctxPayload.CompanyId)

		w.WriteHeader(http.StatusOK)
		component := partials.MaterialsItemsTable(acus)
		component.Render(r.Context(), w)

	case http.MethodGet:
		materials := s.returnAllSelects([]string{"materials"}, ctxPayload.CompanyId)["materials"]

		w.WriteHeader(http.StatusOK)
		component := partials.MaterialsItemsForm(nil, parsedId, materials)
		component.Render(r.Context(), w)

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *Server) MaterialItemsOperations(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	parsedId, err := utils.ValidateUUID(mux.Vars(r)["id"], "rubro")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parsedMaterialId, err := utils.ValidateUUID(mux.Vars(r)["materialId"], "rubro")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		if err := s.DB.DeleteMaterialsByItem(parsedId, parsedMaterialId, ctxPayload.CompanyId); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		acus := s.DB.GetMaterialsByItem(parsedId, ctxPayload.CompanyId)

		w.WriteHeader(http.StatusOK)
		component := partials.MaterialsItemsTable(acus)
		component.Render(r.Context(), w)

	case http.MethodPut:
		r.ParseForm()

		quantityText := r.Form.Get("quantity")
		quantity, err := utils.ConvertFloat(quantityText, "cantidad", true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if err := s.DB.UpdateMaterialByItem(parsedId, parsedMaterialId, quantity, ctxPayload.CompanyId); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		acus := s.DB.GetMaterialsByItem(parsedId, ctxPayload.CompanyId)

		w.WriteHeader(http.StatusOK)
		component := partials.MaterialsItemsTable(acus)
		component.Render(r.Context(), w)

	case http.MethodGet:
		im := s.DB.GetQuantityByMaterialAndItem(parsedId, parsedMaterialId, ctxPayload.CompanyId)
		materials := s.returnAllSelects([]string{"materials"}, ctxPayload.CompanyId)["materials"]

		w.WriteHeader(http.StatusOK)
		component := partials.MaterialsItemsForm(&im, parsedId, materials)
		component.Render(r.Context(), w)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

}
