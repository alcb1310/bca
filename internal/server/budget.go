package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/transaction/partials"
)

func (s *Server) BudgetsTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	if r.Method == http.MethodPost {
		r.ParseForm()
		p := r.Form.Get("project")
		if p == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Seleccione un proyecto"))
			return
		}
		pId, _ := utils.ValidateUUID(p, "proyecto")

		bi := r.Form.Get("budgetItem")
		if bi == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Seleccione una partida"))
			return
		}
		bId, _ := utils.ValidateUUID(bi, "partida")

		q, err := utils.ConvertFloat(r.Form.Get("quantity"), "cantidad", true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		c, err := utils.ConvertFloat(r.Form.Get("cost"), "costo", true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		b := &types.CreateBudget{
			ProjectId:    pId,
			BudgetItemId: bId,
			CompanyId:    ctx.CompanyId,
			Quantity:     q,
			Cost:         c,
		}

		if _, err := s.DB.CreateBudget(b); err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Ya existe partida %s en el proyecto %s", b.BudgetItemId, b.ProjectId)))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Println(err)
			return
		}
	}

	var err error
	project_id := uuid.Nil

	p := r.URL.Query().Get("proyecto")
	if p != "" {
		project_id, err = utils.ValidateUUID(p, "proyecto")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Println(err)
			return
		}
	}

	se := r.URL.Query().Get("buscar")

	b, _ := s.DB.GetBudgets(ctx.CompanyId, project_id, se)
	component := partials.BudgetTable(b)
	component.Render(r.Context(), w)
}

func (s *Server) BudgetAdd(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	data := []string{"projects", "budgetitems"}
	res := s.returnAllSelects(data, ctx.CompanyId)
	projects := res["projects"]
	budgetItems := res["budgetitems"]

	component := partials.EditBudget(nil, projects, budgetItems)
	component.Render(r.Context(), w)
}

func (s *Server) BudgetEdit(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	projectId, _ := utils.ValidateUUID(mux.Vars(r)["projectId"], "proyecto")
	budgetItemId, _ := utils.ValidateUUID(mux.Vars(r)["budgetItemId"], "partida")

	data := []string{"projects", "budgetitems"}
	res := s.returnAllSelects(data, ctx.CompanyId)
	projects := res["projects"]
	budgetItems := res["budgetitems"]

	bd, _ := s.DB.GetOneBudget(ctx.CompanyId, projectId, budgetItemId)

	switch r.Method {
	case http.MethodPut:
		r.ParseForm()
		q, err := utils.ConvertFloat(r.Form.Get("quantity"), "cantidad", true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		c, err := utils.ConvertFloat(r.Form.Get("cost"), "costo", true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		budget := types.CreateBudget{
			ProjectId:    projectId,
			BudgetItemId: budgetItemId,
			Quantity:     q,
			Cost:         c,
			CompanyId:    ctx.CompanyId,
		}

		bu := types.Budget{
			ProjectId:         projectId,
			BudgetItemId:      budgetItemId,
			InitialQuantity:   bd.InitialQuantity,
			InitialCost:       bd.InitialCost,
			InitialTotal:      bd.InitialTotal,
			SpentQuantity:     bd.SpentQuantity,
			SpentTotal:        bd.SpentTotal,
			RemainingQuantity: bd.RemainingQuantity,
			RemainingCost:     bd.RemainingCost,
			RemainingTotal:    bd.RemainingTotal,
			UpdatedBudget:     bd.UpdatedBudget,
			CompanyId:         ctx.CompanyId,
		}

		if err := s.DB.UpdateBudget(budget, bu); err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("Ya existe partida %s en el proyecto %s", budgetItemId, projectId)))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, _ := s.DB.GetBudgets(ctx.CompanyId, uuid.Nil, "")

		component := partials.BudgetTable(b)
		component.Render(r.Context(), w)

	case http.MethodGet:
		budget := &types.CreateBudget{
			ProjectId:    bd.Project.ID,
			BudgetItemId: bd.BudgetItem.ID,

			Quantity:  bd.RemainingQuantity.Float64,
			Cost:      bd.RemainingCost.Float64,
			CompanyId: ctx.CompanyId,
		}
		component := partials.EditBudget(budget, projects, budgetItems)
		component.Render(r.Context(), w)
	}

}
