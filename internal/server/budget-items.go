package server

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
	"github.com/alcb1310/bca/internal/views/bca/settings/partials"
)

func (s *Server) BudgetItemsTable(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	if err := r.ParseForm(); err != nil {
		slog.Error("BudgetItemsTable error", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	x := r.Form.Get("accumulate") == "accumulate"
	p := r.Form.Get("parent")
	var u *uuid.UUID
	if p == "" {
		u = nil
	} else {
		z, err := uuid.Parse(p)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Error("BudgetItemsTable error", "error", err)
			w.Write([]byte("Código de la partida padre es inválido"))
			return
		}
		u = &z
	}
	acc := sql.NullBool{Valid: true, Bool: x}
	bi := &types.BudgetItem{
		CompanyId:  ctxPayload.CompanyId,
		Code:       r.Form.Get("code"),
		Name:       r.Form.Get("name"),
		ParentId:   u,
		Accumulate: acc,
	}
	if bi.Code == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Debe proporcionar un código de la partida"))
		return
	}
	if bi.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Debe proporcionar un nombre de la partida"))
		return
	}

	if err := s.DB.CreateBudgetItem(bi); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("Ya existe una partida con el mismo código: %s y/o el mismo nombre: %s", bi.Code, bi.Name)))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		slog.Error("BudgetItemsTable error", "error", err)
		return
	}

	search := r.URL.Query().Get("search")
	b, _ := s.DB.GetBudgetItems(ctxPayload.CompanyId, search)
	component := partials.BudgetItemTable(b)
	component.Render(r.Context(), w)
}

func (s *Server) BudgetItemsTableDisplay(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	search := r.URL.Query().Get("search")
	b, _ := s.DB.GetBudgetItems(ctxPayload.CompanyId, search)
	component := partials.BudgetItemTable(b)
	component.Render(r.Context(), w)
}

func (s *Server) BudgetItemAdd(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	p := s.DB.GetBudgetItemsByAccumulate(ctxPayload.CompanyId, true)
	component := partials.EditBudgetItem(nil, p)
	component.Render(r.Context(), w)
}

func (s *Server) BudgetItemEdit(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)
	budgetItem, err := s.DB.GetOneBudgetItem(parsedId, ctxPayload.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Partida no encontrada"))
		return
	}

	p := s.DB.GetBudgetItemsByAccumulate(ctxPayload.CompanyId, true)
	component := partials.EditBudgetItem(budgetItem, p)
	component.Render(r.Context(), w)
}

func (s *Server) EditBudgetItem(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)
	budgetItem, err := s.DB.GetOneBudgetItem(parsedId, ctxPayload.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Partida no encontrada"))
		return
	}

	r.ParseForm()
	biCode := r.Form.Get("code")
	if biCode != "" {
		budgetItem.Code = biCode
	}

	biName := r.Form.Get("name")
	if biName != "" {
		budgetItem.Name = biName
	}

	x := r.Form.Get("accumulate") == "accumulate"
	acc := sql.NullBool{Valid: true, Bool: x}
	budgetItem.Accumulate = acc

	var u *uuid.UUID
	p := r.Form.Get("parent")
	if p == "" {
		u = budgetItem.ParentId
	} else {
		z, err := uuid.Parse(p)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Código de la partida padre es inválido"))
			slog.Error("BudgetItemEdit error", "error", err)
			return
		}
		u = &z
	}
	budgetItem.ParentId = u

	if budgetItem.Code == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Debe proporcionar un código de la partida"))
		return
	}
	if budgetItem.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Debe proporcionar un nombre de la partida"))
		return
	}

	if err := s.DB.UpdateBudgetItem(budgetItem); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("Ya existe una partida con el mismo código: %s y/o el mismo nombre: %s", budgetItem.Code, budgetItem.Name)))
			return
		}

		if strings.Contains(err.Error(), "partida padre") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		slog.Error("BudgetItemEdit error", "error", err)
		return
	}

	b, _ := s.DB.GetBudgetItems(ctxPayload.CompanyId, "")
	component := partials.BudgetItemTable(b)
	component.Render(r.Context(), w)
}
