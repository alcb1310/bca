package server

import (
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

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
