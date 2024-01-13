package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/transaction/partials"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) BudgetsTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)

	if r.Method == http.MethodPost {
		r.ParseForm()
		p := r.Form.Get("project")
		pId, _ := uuid.Parse(p)
		bi := r.Form.Get("budgetItem")
		bId, _ := uuid.Parse(bi)
		q, err := strconv.ParseFloat(r.Form.Get("quantity"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		c, err := strconv.ParseFloat(r.Form.Get("cost"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	b, _ := s.DB.GetBudgets(ctx.CompanyId)
	component := partials.BudgetTable(b)
	component.Render(r.Context(), w)
}

func (s *Server) BudgetAdd(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)

	p := s.DB.GetActiveProjects(ctx.CompanyId, true)
	projectMap := make(map[string]string)
	for _, v := range p {
		projectMap[v.ID.String()] = v.Name
	}

	b := s.DB.GetBudgetItemsByAccumulate(ctx.CompanyId, false)
	budgetItemMap := make(map[string]string)
	for _, v := range b {
		budgetItemMap[v.ID.String()] = v.Name
	}

	component := partials.EditBudget(nil, projectMap, budgetItemMap)
	component.Render(r.Context(), w)
}

func (s *Server) BudgetEdit(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	projectId, _ := uuid.Parse(mux.Vars(r)["projectId"])
	budgetItemId, _ := uuid.Parse(mux.Vars(r)["budgetItemId"])

	p := s.DB.GetActiveProjects(ctx.CompanyId, true)
	projectMap := make(map[string]string)
	for _, v := range p {
		projectMap[v.ID.String()] = v.Name
	}

	b := s.DB.GetBudgetItemsByAccumulate(ctx.CompanyId, false)
	budgetItemMap := make(map[string]string)
	for _, v := range b {
		budgetItemMap[v.ID.String()] = v.Name
	}

	bd, _ := s.DB.GetOneBudget(ctx.CompanyId, projectId, budgetItemId)

	switch r.Method {
	case http.MethodPut:
		r.ParseForm()
		q, err := strconv.ParseFloat(r.Form.Get("quantity"), 10)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		c, err := strconv.ParseFloat(r.Form.Get("cost"), 10)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, _ := s.DB.GetBudgets(ctx.CompanyId)

		component := partials.BudgetTable(b)
		component.Render(r.Context(), w)

	case http.MethodGet:
		budget := &types.CreateBudget{
			ProjectId:    bd.Project.ID,
			BudgetItemId: bd.BudgetItem.ID,

			Quantity:  *bd.RemainingQuantity,
			Cost:      *bd.RemainingCost,
			CompanyId: ctx.CompanyId,
		}
		component := partials.EditBudget(budget, projectMap, budgetItemMap)
		component.Render(r.Context(), w)
	}

}
