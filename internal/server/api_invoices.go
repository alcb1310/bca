package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllInvoices(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	invoices, err := s.DB.GetInvoices(ctx.CompanyId)
	if err != nil {
		slog.Error("ApiGetAllInvoices", "err", err)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(invoices)
}
