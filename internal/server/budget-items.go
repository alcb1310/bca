package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/settings/partials"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) BudgetItemsTable(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)

	if r.Method == http.MethodPost {
		r.ParseForm()
		x := r.Form.Get("accumulate") == "accumulate"
		p := r.Form.Get("parent")
		var u *uuid.UUID
		if p == "" {
			u = nil
		} else {
			z, err := uuid.Parse(p)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			u = &z
		}
		bi := &types.BudgetItem{
			CompanyId:  ctxPayload.CompanyId,
			Code:       r.Form.Get("code"),
			Name:       r.Form.Get("name"),
			ParentId:   u,
			Accumulate: &x,
		}

		if err := s.DB.CreateBudgetItem(bi); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	search := r.URL.Query().Get("search")
	b, _ := s.DB.GetBudgetItems(ctxPayload.CompanyId, search)
	component := partials.BudgetItemTable(b)
	component.Render(r.Context(), w)
}

func (s *Server) BudgetItemAdd(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	p := s.DB.GetBudgetItemsByAccumulate(ctxPayload.CompanyId, true)
	component := partials.EditBudgetItem(nil, p)
	component.Render(r.Context(), w)
}

func (s *Server) BudgetItemEdit(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]
	parsedId, _ := uuid.Parse(id)
	budgetItem, _ := s.DB.GetOneBudgetItem(parsedId, ctxPayload.CompanyId)

	switch r.Method {
	case http.MethodPut:
		r.ParseForm()
		budgetItem.Code = r.Form.Get("code")
		budgetItem.Name = r.Form.Get("name")
		x := r.Form.Get("accumulate") == "accumulate"
		budgetItem.Accumulate = &x
		p := r.Form.Get("parent")
		var u *uuid.UUID
		if p == "" {
			u = nil
		} else {
			z, err := uuid.Parse(p)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			u = &z
		}
		budgetItem.ParentId = u

		if err := s.DB.UpdateBudgetItem(budgetItem); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, _ := s.DB.GetBudgetItems(ctxPayload.CompanyId, "")
		component := partials.BudgetItemTable(b)
		component.Render(r.Context(), w)

	case http.MethodGet:
		p := s.DB.GetBudgetItemsByAccumulate(ctxPayload.CompanyId, true)
		component := partials.EditBudgetItem(budgetItem, p)
		component.Render(r.Context(), w)
	}
}
