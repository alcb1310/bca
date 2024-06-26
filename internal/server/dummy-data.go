package server

import (
	"bca-go-final/internal/utils"
	"encoding/json"
	"net/http"
)

func (s *Server) loadDummyDataHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctx, _ := utils.GetMyPaload(r)
	companyId := ctx.CompanyId

	switch r.Method {
	case http.MethodPost:
		err := s.DB.LoadDummyData(companyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusCreated)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
