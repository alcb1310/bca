package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alcb1310/bca/internal/types"
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

func (s *Server) ApiGetOneInvoice(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	id := chi.URLParam(r, "id")
	if id == "crear" {
		invoice := types.InvoiceCreate{}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(invoice)
		return
	}

func (s *Server) ApiCreateInvoice(w http.ResponseWriter, r *http.Request) {
	slog.Debug("ApiCreateInvoice")
	var data types.InvoiceCreate
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		slog.Error("ApiCreateInvoice", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	ctx, _ := utils.GetMyPaload(r)
	data.CompanyId = ctx.CompanyId

	if err := s.DB.CreateInvoice(&data); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Ya existe esa factura"
			_ = json.NewEncoder(w).Encode(errorResponse)
			return
		}

		slog.Error("ApiCreateInvoice", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(data)
}
