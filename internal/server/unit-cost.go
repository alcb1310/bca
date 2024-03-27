package server

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/unit_cost"
	"bca-go-final/internal/views/bca/unit_cost/partials"
)

func (s *Server) UnitQuantity(w http.ResponseWriter, r *http.Request) {
	component := unit_cost.UnitCostQuantity()
	component.Render(r.Context(), w)
}

func (s *Server) CantidadesTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	quantities := s.DB.CantidadesTable(ctx.CompanyId)

	component := partials.CantidadesTable(quantities)
	component.Render(r.Context(), w)

}

func (s *Server) CantidadesAdd(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	switch r.Method {
	case http.MethodGet:
		rub, _ := s.DB.GetAllRubros(ctx.CompanyId)
		rubSelect := []types.Select{}
		for _, v := range rub {
			x := types.Select{
				Key:   v.Id.String(),
				Value: v.Name,
			}
			rubSelect = append(rubSelect, x)
		}

		p := s.DB.GetActiveProjects(ctx.CompanyId, true)
		projects := []types.Select{}
		for _, v := range p {
			x := types.Select{
				Key:   v.ID.String(),
				Value: v.Name,
			}
			projects = append(projects, x)
		}

		component := partials.EditCantidades(nil, projects, rubSelect)
		component.Render(r.Context(), w)

	case http.MethodPost:
		r.ParseForm()

		pId := r.Form.Get("project")
		parsedProjectId, err := uuid.Parse(pId)
		if err != nil {
			if strings.Contains(err.Error(), "length: 0") {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Seleccione un proyecto"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		rubroId := r.Form.Get("item")
		parsedRubroId, err := uuid.Parse(rubroId)
		if err != nil {
			if strings.Contains(err.Error(), "length: 0") {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Seleccione un rubro"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		q := r.Form.Get("quantity")
		if q == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("La cantidad es requerida"))
			return
		}

		quantity, err := strconv.ParseFloat(q, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("La cantidad debe ser un número"))
			return
		}
		if quantity <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("La cantidad debe ser mayor a 0"))
			return
		}

		if err := s.DB.CreateCantidades(parsedProjectId, parsedRubroId, quantity, ctx.CompanyId); err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("La cantidad ya existe"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error al crear la cantidad"))
			log.Println(err.Error())
			return
		}

		quantities := s.DB.CantidadesTable(ctx.CompanyId)

		component := partials.CantidadesTable(quantities)
		component.Render(r.Context(), w)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func (s *Server) UnitAnalysis(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	p := s.DB.GetActiveProjects(ctxPayload.CompanyId, true)
	projects := []types.Select{}
	for _, v := range p {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		projects = append(projects, x)
	}

	component := unit_cost.Analysis(projects)
	component.Render(r.Context(), w)
}

func (s *Server) AnalysisTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	p := r.URL.Query().Get("project")
	projectId, err := uuid.Parse(p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Seleccione un proyecto"))
		return
	}

	analysis := s.DB.AnalysisReport(projectId, ctx.CompanyId)
	keys := []string{}
	for k := range analysis {
		keys = append(keys, k)
	}

	slices.Sort(keys)
	component := partials.AnalysisTable(analysis, keys)
	component.Render(r.Context(), w)
}
