package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
)

func (s *Server) ApiGetAllInvoiceDetails(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	invoiceResponse, err := s.DB.GetAllDetails(parsedId, ctx.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(invoiceResponse)
}

func (s *Server) ApiCreateInvoiceDetails(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var newInvoiceDetail types.InvoiceDetailCreate
	err = json.NewDecoder(r.Body).Decode(&newInvoiceDetail)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	newInvoiceDetail.InvoiceId = parsedId
	newInvoiceDetail.CompanyId = ctx.CompanyId

	if err := s.DB.AddDetail(newInvoiceDetail); err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == "23505" {
			w.WriteHeader(http.StatusConflict)
			errorResponse := make(map[string]string)
			errorResponse["error"] = "Detalle ya existe"
			_ = json.NewEncoder(w).Encode(errorResponse)
			return
		}

		if err.Error() == "sql: no rows in result set" {
			w.WriteHeader(http.StatusNotFound)
			errorResponse := make(map[string]string)
			errorResponse["error"] = "No se encontro la cuenta en el presupuesto"
			_ = json.NewEncoder(w).Encode(errorResponse)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) ApiDeleteInvoiceDetails(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}
	biId := chi.URLParam(r, "budgetItemId")
	parsedBudgetItemId, err := uuid.Parse(biId)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if err := s.DB.DeleteDetail(parsedId, parsedBudgetItemId, ctx.CompanyId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := make(map[string]string)
		errorResponse["error"] = err.Error()
		_ = json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
