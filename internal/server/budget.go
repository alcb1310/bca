package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/transaction/partials"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) BudgetsTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)

	if r.Method == http.MethodPost {
		r.ParseForm()
		p := r.Form.Get("project")
		if p == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Seleccione un proyecto"))
			return
		}
		pId, _ := uuid.Parse(p)
		bi := r.Form.Get("budgetItem")
		if bi == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Seleccione una partida"))
			return
		}
		bId, _ := uuid.Parse(bi)
		q, err := strconv.ParseFloat(r.Form.Get("quantity"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("La cantidad debe ser un número"))
			return
		}
		c, err := strconv.ParseFloat(r.Form.Get("cost"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("El costo debe ser un número"))
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

	b, _ := s.DB.GetBudgets(ctx.CompanyId)
	component := partials.BudgetTable(b)
	component.Render(r.Context(), w)
}

func (s *Server) BudgetAdd(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)

	p := s.DB.GetActiveProjects(ctx.CompanyId, true)
	projectMap := []types.Select{}
	for _, v := range p {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		projectMap = append(projectMap, x)
	}

	b := s.DB.GetBudgetItemsByAccumulate(ctx.CompanyId, false)
	budgetItemMap := []types.Select{}
	for _, v := range b {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		budgetItemMap = append(budgetItemMap, x)
	}

	component := partials.EditBudget(nil, projectMap, budgetItemMap)
	component.Render(r.Context(), w)
}

func (s *Server) BudgetEdit(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	projectId, _ := uuid.Parse(mux.Vars(r)["projectId"])
	budgetItemId, _ := uuid.Parse(mux.Vars(r)["budgetItemId"])

	p := s.DB.GetActiveProjects(ctx.CompanyId, true)
	projectMap := []types.Select{}
	for _, v := range p {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		projectMap = append(projectMap, x)
	}

	b := s.DB.GetBudgetItemsByAccumulate(ctx.CompanyId, false)
	budgetItemMap := []types.Select{}
	for _, v := range b {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		budgetItemMap = append(budgetItemMap, x)
	}

	bd, _ := s.DB.GetOneBudget(ctx.CompanyId, projectId, budgetItemId)

	switch r.Method {
	case http.MethodPut:
		r.ParseForm()
		q, err := strconv.ParseFloat(r.Form.Get("quantity"), 10)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("La cantidad debe ser un número"))
			return
		}
		c, err := strconv.ParseFloat(r.Form.Get("cost"), 10)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("El costo debe ser un número"))
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

		b, _ := s.DB.GetBudgets(ctx.CompanyId)

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
		component := partials.EditBudget(budget, projectMap, budgetItemMap)
		component.Render(r.Context(), w)
	}

}
