package server

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

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
		rubros := s.getSelect("rubros", ctx.CompanyId)
		projects := s.getSelect("projects", ctx.CompanyId)

		component := partials.EditCantidades(nil, projects, rubros)
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
	projects := s.getSelect("projects", ctxPayload.CompanyId)

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

func (s *Server) CantidadesEdit(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	id := mux.Vars(r)["id"]
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error al parsear el ID"))
		return
	}

	switch r.Method {
	case http.MethodDelete:
		if err := s.DB.DeleteCantidades(parsedId, ctx.CompanyId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error al borrar la cantidad"))
			log.Println(err.Error())
			return
		}

		quantities := s.DB.CantidadesTable(ctx.CompanyId)

		component := partials.CantidadesTable(quantities)
		component.Render(r.Context(), w)

	case http.MethodGet:
		rubros := s.getSelect("rubros", ctx.CompanyId)
		projects := s.getSelect("projects", ctx.CompanyId)

		quantity := s.DB.GetOneQuantityById(parsedId, ctx.CompanyId)

		component := partials.EditCantidades(&quantity, projects, rubros)
		component.Render(r.Context(), w)

	case http.MethodPut:
		r.ParseForm()

		q := r.Form.Get("quantity")
		if q == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese una cantidad"))
			return
		}
		quantity, err := strconv.ParseFloat(q, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Cantidad debe ser numérica"))
			return
		}

		quan := s.DB.GetOneQuantityById(parsedId, ctx.CompanyId)

		quan.Quantity = quantity

		if err := s.DB.UpdateQuantity(quan, ctx.CompanyId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error al actualizar la cantidad"))
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
