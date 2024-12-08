package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
	"github.com/alcb1310/bca/internal/views/bca/transaction/partials"
)

func (s *Server) BudgetsTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	var err error

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	projectId := r.Form.Get("project")
	if projectId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Seleccione un proyecto"))
		return
	}
	pId, err := uuid.Parse(projectId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Código del proyecto inválido"))
		return
	}

	bi := r.Form.Get("budgetItem")
	if bi == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Seleccione una partida"))
		return
	}
	bId, err := uuid.Parse(bi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Código de la partida inválido"))
		return
	}

	q, err := strconv.ParseFloat(r.Form.Get("quantity"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("La cantidad debe ser un número"))
		return
	}

	c, err := strconv.ParseFloat(r.Form.Get("cost"), 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("El costo debe ser un número"))
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
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(fmt.Sprintf("Ya existe partida %s en el proyecto %s", b.BudgetItemId, b.ProjectId)))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		slog.Error(err.Error())
		return
	}

	project_id := uuid.Nil

	p := r.URL.Query().Get("proyecto")
	if p != "" {
		project_id, err = uuid.Parse(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			slog.Error(err.Error())
			return
		}
	}

	se := r.URL.Query().Get("buscar")

	budgets, _ := s.DB.GetBudgets(ctx.CompanyId, project_id, se)
	component := partials.BudgetTable(budgets)
	_ = component.Render(r.Context(), w)
}

func (s *Server) BudgetsTableDisplay(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	var err error
	project_id := uuid.Nil

	p := r.URL.Query().Get("proyecto")
	if p != "" {
		project_id, err = uuid.Parse(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			slog.Error(err.Error())
			return
		}
	}

	se := r.URL.Query().Get("buscar")

	b, _ := s.DB.GetBudgets(ctx.CompanyId, project_id, se)
	component := partials.BudgetTable(b)
	_ = component.Render(r.Context(), w)
}

func (s *Server) BudgetAdd(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

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
	_ = component.Render(r.Context(), w)
}

func (s *Server) BudgetEdit(w http.ResponseWriter, r *http.Request) {
	var (
		q, c float64
		err  error
	)

	ctx, _ := utils.GetMyPaload(r)
	pId := chi.URLParam(r, "projectId")
	bId := chi.URLParam(r, "budgetItemId")
	projectId, _ := uuid.Parse(pId)
	budgetItemId, _ := uuid.Parse(bId)

	bd, _ := s.DB.GetOneBudget(ctx.CompanyId, projectId, budgetItemId)

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	fQuan := r.Form.Get("quantity")
	if fQuan != "" {
		q, err = strconv.ParseFloat(fQuan, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("La cantidad debe ser un número"))
			return
		}
	}

	fCost := r.Form.Get("cost")
	if fCost != "" {
		c, err = strconv.ParseFloat(r.Form.Get("cost"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("El costo debe ser un número"))
			return
		}
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
			_, _ = w.Write([]byte(fmt.Sprintf("Ya existe partida %s en el proyecto %s", budgetItemId, projectId)))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	budgets, _ := s.DB.GetBudgets(ctx.CompanyId, uuid.Nil, "")

	component := partials.BudgetTable(budgets)
	_ = component.Render(r.Context(), w)
}

func (s *Server) SingleBudget(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	pId := chi.URLParam(r, "projectId")
	bId := chi.URLParam(r, "budgetItemId")
	projectId, _ := uuid.Parse(pId)
	budgetItemId, _ := uuid.Parse(bId)

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

	budget := &types.CreateBudget{
		ProjectId:    bd.Project.ID,
		BudgetItemId: bd.BudgetItem.ID,

		Quantity:  bd.RemainingQuantity.Float64,
		Cost:      bd.RemainingCost.Float64,
		CompanyId: ctx.CompanyId,
	}
	component := partials.EditBudget(budget, projectMap, budgetItemMap)
	_ = component.Render(r.Context(), w)
}
