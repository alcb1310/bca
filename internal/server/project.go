package server

import "net/http"
import (
	"encoding/json"
	"net/http"
)

func (s *Server) AllProjects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		projects, err := s.DB.GetAllProjects(ctxPayload.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(projects)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
