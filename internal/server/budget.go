package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/transaction/partials"
	"net/http"
	"strconv"

	"github.com/google/uuid"
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
