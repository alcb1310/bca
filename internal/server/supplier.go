package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) AllSuppliers(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctxPayload, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodGet:
		suppliers, err := s.DB.GetAllSuppliers(ctxPayload.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(suppliers)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
