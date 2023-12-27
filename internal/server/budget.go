package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) AllBudgets(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctx, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		budgets, err := s.DB.GetBudgets(ctx.CompanyId)
		if err != nil {
			resp["error"] = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(budgets)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
