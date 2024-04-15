package server

import (
	"log"
	"net/http"
	"slices"
	"strings"

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
		data := []string{"rubros", "projects"}
		results := s.returnAllSelects(data, ctx.CompanyId)
		rubros := results["rubros"]
		projects := results["projects"]

		component := partials.EditCantidades(nil, projects, rubros)
		component.Render(r.Context(), w)

	case http.MethodPost:
		r.ParseForm()

		pId := r.Form.Get("project")
		parsedProjectId, err := utils.ValidateUUID(pId, "proyecto")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Seleccione un proyecto"))
			log.Println(err)
			return
		}

		rubroId := r.Form.Get("item")
		parsedRubroId, err := utils.ValidateUUID(rubroId, "rubro")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Seleccione un rubro"))
			log.Println(err)
			return
		}

		q := r.Form.Get("quantity")
		quantity, err := utils.ConvertFloat(q, "cantidad", true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
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
	projects := s.returnAllSelects([]string{"projects"}, ctxPayload.CompanyId)["projects"]

	component := unit_cost.Analysis(projects)
	component.Render(r.Context(), w)
}

func (s *Server) AnalysisTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	p := r.URL.Query().Get("project")
	projectId, err := utils.ValidateUUID(p, "proyecto")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
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

	parsedId, err := utils.ValidateUUID(mux.Vars(r)["id"], "cantidad")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
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
		data := []string{"rubros", "projects"}
		results := s.returnAllSelects(data, ctx.CompanyId)
		rubros := results["rubros"]
		projects := results["projects"]

		quantity := s.DB.GetOneQuantityById(parsedId, ctx.CompanyId)

		component := partials.EditCantidades(&quantity, projects, rubros)
		component.Render(r.Context(), w)

	case http.MethodPut:
		r.ParseForm()

		q := r.Form.Get("quantity")
		quantity, err := utils.ConvertFloat(q, "cantidad", true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
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
