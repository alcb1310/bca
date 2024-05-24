package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/settings/partials"
)

func (s *Server) BudgetItemsTable(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		x := r.Form.Get("accumulate") == "accumulate"
		p := r.Form.Get("parent")
		var u *uuid.UUID
		if p == "" {
			u = nil
		} else {
			z, err := utils.ValidateUUID(p, "proveedor")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Código de la partida padre es inválido"))
				w.Write([]byte(err.Error()))
				return
			}
			u = &z
		}
		bi := &types.BudgetItem{
			CompanyId:  ctxPayload.CompanyId,
			Code:       r.Form.Get("code"),
			Name:       r.Form.Get("name"),
			ParentId:   u,
			Accumulate: sql.NullBool{Valid: true, Bool: x},
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
			log.Println(err)
			return
		}
	}

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
	id := mux.Vars(r)["id"]
	parsedId, _ := utils.ValidateUUID(id, "partida")
	budgetItem, _ := s.DB.GetOneBudgetItem(parsedId, ctxPayload.CompanyId)

	switch r.Method {
	case http.MethodPut:
		r.ParseForm()
		budgetItem.Code = r.Form.Get("code")
		budgetItem.Name = r.Form.Get("name")
		x := r.Form.Get("accumulate") == "accumulate"
		acc := sql.NullBool{Valid: true, Bool: x}
		budgetItem.Accumulate = acc
		p := r.Form.Get("parent")
		var u *uuid.UUID
		if p == "" {
			u = nil
		} else {
			z, err := utils.ValidateUUID(p, "partida")
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Código de la partida padre es inválido"))
				log.Println(err)
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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Println(err)
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
