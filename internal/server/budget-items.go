package server

import (
	"bca-go-final/internal/types"
	"encoding/json"
	"net/http"
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
