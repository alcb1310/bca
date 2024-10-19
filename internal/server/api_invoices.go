package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

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
	id := chi.URLParam(r, "id")
	if id == "crear" {
		invoice := types.InvoiceCreate{}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(invoice)
		return
	}
	parsedInvoiceId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	ctx, _ := utils.GetMyPaload(r)

	invoiceResponse, err := s.DB.GetOneInvoice(parsedInvoiceId, ctx.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	dt := invoiceResponse.InvoiceDate.In(time.Local)
	dt = dt.Add(time.Duration(s.Timezone) * -1 * time.Hour)

	errorCreate := types.InvoiceCreate{
		Id:            &invoiceResponse.Id,
		SupplierId:    &invoiceResponse.Supplier.ID,
		ProjectId:     &invoiceResponse.Project.ID,
		InvoiceNumber: &invoiceResponse.InvoiceNumber,
		InvoiceDate:   &dt,
		InvoiceTotal:  invoiceResponse.InvoiceTotal,
		IsBalanced:    invoiceResponse.IsBalanced,
		CompanyId:     invoiceResponse.CompanyId,
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(errorCreate)
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

func (s *Server) ApiUpdateInvoice(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementar
	id := chi.URLParam(r, "id")
	parsedInvoiceId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	ctx, _ := utils.GetMyPaload(r)

	invoiceToEdit, err := s.DB.GetOneInvoice(parsedInvoiceId, ctx.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := make(map[string]string)
		errorResponse["error"] = "No se encontro la factura"
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var data types.InvoiceCreate
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		slog.Error("ApiUpdateInvoice", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	data.Id = &invoiceToEdit.Id
	data.SupplierId = &invoiceToEdit.Supplier.ID
	data.ProjectId = &invoiceToEdit.Project.ID
	data.IsBalanced = invoiceToEdit.IsBalanced
	data.InvoiceTotal = invoiceToEdit.InvoiceTotal
	data.IsBalanced = invoiceToEdit.IsBalanced
	data.CompanyId = ctx.CompanyId

	if err := s.DB.UpdateInvoice(data); err != nil {
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

	w.WriteHeader(http.StatusNoContent)
}
