package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
	"github.com/alcb1310/bca/internal/views/bca/settings/partials"
)

func (s *Server) RubrosTable(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	rubros, _ := s.DB.GetAllRubros(ctxPayload.CompanyId)

	component := partials.RubrosTable(rubros)
	component.Render(r.Context(), w)
}

func (s *Server) MaterialsByItem(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
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

	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.ParseForm()

	materialId := r.Form.Get("material")
	if materialId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Seleccione un Material"))
		return
	}
	parsedMaterialId, err := uuid.Parse(materialId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Material Incorrecto"))
		return
	}

	quantityText := r.Form.Get("quantity")
	if quantityText == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese una Cantidad"))
		return
	}

	quantity, err := strconv.ParseFloat(quantityText, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cantidad debe ser un valor numérico"))
		return
	}

	if quantity <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("La Cantidad debe ser mayor a 0"))
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
}

func (s *Server) MaterialsByItemFormDisplay(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	materials := s.DB.GetAllMaterials(ctxPayload.CompanyId)
	materialsSelect := []types.Select{}
	for _, m := range materials {
		materialsSelect = append(materialsSelect, types.Select{Key: m.Id.String(), Value: m.Name})
	}

	w.WriteHeader(http.StatusOK)
	component := partials.MaterialsItemsForm(nil, parsedId, materialsSelect)
	component.Render(r.Context(), w)
}

func (s *Server) MaterialItemsOperationsUpdate(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	materialId := chi.URLParam(r, "materialId")
	parsedMaterialId, err := uuid.Parse(materialId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.ParseForm()

	quantityText := r.Form.Get("quantity")
	if quantityText == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese una Cantidad"))
		return
	}

	quantity, err := strconv.ParseFloat(quantityText, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cantidad debe ser un valor numérico"))
		return
	}

	if quantity <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("La Cantidad debe ser mayor a 0"))
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
}

func (s *Server) MaterialItemsOperationsDelete(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	materialId := chi.URLParam(r, "materialId")
	parsedMaterialId, err := uuid.Parse(materialId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.DB.DeleteMaterialsByItem(parsedId, parsedMaterialId, ctxPayload.CompanyId); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	acus := s.DB.GetMaterialsByItem(parsedId, ctxPayload.CompanyId)

	w.WriteHeader(http.StatusOK)
	component := partials.MaterialsItemsTable(acus)
	component.Render(r.Context(), w)
}

func (s *Server) MaterialItemsOperationsDisplay(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	materialId := chi.URLParam(r, "materialId")
	parsedMaterialId, err := uuid.Parse(materialId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	im := s.DB.GetQuantityByMaterialAndItem(parsedId, parsedMaterialId, ctxPayload.CompanyId)

	materials := s.DB.GetAllMaterials(ctxPayload.CompanyId)
	materialsSelect := []types.Select{}
	for _, m := range materials {
		materialsSelect = append(materialsSelect, types.Select{Key: m.Id.String(), Value: m.Name})
	}

	w.WriteHeader(http.StatusOK)
	component := partials.MaterialsItemsForm(&im, parsedId, materialsSelect)
	component.Render(r.Context(), w)
}
