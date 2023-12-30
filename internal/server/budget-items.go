package server

import (
	"bca-go-final/internal/types"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) AllBudgetItems(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctxPayload, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodGet:
		bis, err := s.DB.GetBudgetItems(ctxPayload.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bis)

	case http.MethodPost:
		bi := &types.BudgetItem{}
		bi.CompanyId = ctxPayload.CompanyId

		if err := json.NewDecoder(r.Body).Decode(bi); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		if bi.Code == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "code cannot be empty"
			resp["field"] = "code"
			json.NewEncoder(w).Encode(resp)
			return
		}
		if bi.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "name cannot be empty"
			resp["field"] = "name"
			json.NewEncoder(w).Encode(resp)
			return
		}
		if bi.Accumulate == nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "accumulate cannot be empty"
			resp["field"] = "accumulate"
			json.NewEncoder(w).Encode(resp)
			return
		}

		if err := s.DB.CreateBudgetItem(bi); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(bi)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) OneBudgetItem(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctxPayload, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]

	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	bi, err := s.DB.GetOneBudgetItem(parsedId, ctxPayload.CompanyId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	switch r.Method {
	case http.MethodPut:
		b := &types.BudgetItem{}

		updated := false
		if err := json.NewDecoder(r.Body).Decode(b); err != nil {
			if strings.Contains(err.Error(), "invalid UUID length: 0") {
				b.ParentId = nil
				updated = true
			} else {
				w.WriteHeader(http.StatusBadRequest)
				resp["error"] = err.Error()
				json.NewEncoder(w).Encode(resp)
				return
			}
		}

		if b.Code == "" {
			b.Code = bi.Code
		}
		if b.Name == "" {
			b.Name = bi.Name
		}
		if b.Accumulate == nil {
			b.Accumulate = bi.Accumulate
		}

		if !updated && b.ParentId == nil {
			b.ParentId = bi.ParentId
		}
		b.CompanyId = ctxPayload.CompanyId
		b.ID = parsedId

		err := s.DB.UpdateBudgetItem(b)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(b)

	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bi)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
