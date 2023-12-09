package server

import (
	"bca-go-final/internal/types"
	"encoding/json"
	"net/http"
)

func (s *Server) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		users, err := s.DB.GetAllUsers(ctxPayload.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
