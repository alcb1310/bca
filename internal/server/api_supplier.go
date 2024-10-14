package server

import (
	"encoding/json"
	"net/http"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllSuppliers(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	queryParams := r.URL.Query()
	search := queryParams.Get("query")

	suppliers, _ := s.DB.GetAllSuppliers(ctx.CompanyId, search)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(suppliers)
}
