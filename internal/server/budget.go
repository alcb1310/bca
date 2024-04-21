package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/transaction/partials"
)

func createBudget(form url.Values, companyId uuid.UUID, s *Server) (types.Budget, error) {

	projectID, err := utils.ValidateUUID(form.Get("project"), "proyecto")
	if err != nil {
		return types.Budget{}, err
	}
	budgetItemID, err := utils.ValidateUUID(form.Get("budgetItem"), "partida")
	if err != nil {
		return types.Budget{}, err
	}

	quantity, err := utils.ConvertFloat(form.Get("quantity"), "cantidad", true)
	if err != nil {
		return types.Budget{}, err
	}

	cost, err := utils.ConvertFloat(form.Get("cost"), "costo", true)
	if err != nil {
		return types.Budget{}, err
	}
	budgetInfo := &types.CreateBudget{
		ProjectId:    projectID,
		BudgetItemId: budgetItemID,
		CompanyId:    companyId,
		Quantity:     quantity,
		Cost:         cost,
	}

	budget, err := s.DB.CreateBudget(budgetInfo)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return types.Budget{}, errors.New(fmt.Sprintf("Ya existe partida %s en el proyecto %s", budgetInfo.BudgetItemId, budgetInfo.ProjectId))
		}
		log.Println(fmt.Sprintf("Error creating budget. Err: %v", err))
		return types.Budget{}, err
	}

	return budget, nil
}

func (s *Server) BudgetsTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	if r.Method == http.MethodPost {
		r.ParseForm()

		form := r.Form
		_, err := createBudget(form, ctx.CompanyId, s)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	searchTerm := r.URL.Query().Get("buscar")
	searchProjectID := uuid.Nil
	project := r.URL.Query().Get("proyecto")
	if project != "" {
		searchProjectID, _ = utils.ValidateUUID(project, "proyecto")
	}

	budgets, _ := s.DB.GetBudgets(ctx.CompanyId, searchProjectID, searchTerm)
	component := partials.BudgetTable(budgets)
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
