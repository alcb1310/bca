package server

import (
	"bca-go-final/internal/types"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) AllInvoices(w http.ResponseWriter, r *http.Request) {
	var resp = make(map[string]string)
	ctx, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodPost:
		i := &types.InvoiceCreate{
			CompanyId: ctx.CompanyId,
		}
		err := json.NewDecoder(r.Body).Decode(i)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		if i.SupplierId == nil || *i.SupplierId == uuid.Nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "supplier_id cannot be empty"
			resp["field"] = "supplier_id"
			json.NewEncoder(w).Encode(resp)
			return
		}

		if i.ProjectId == nil || *i.ProjectId == uuid.Nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "project_id cannot be empty"
			resp["field"] = "project_id"
			json.NewEncoder(w).Encode(resp)
			return
		}

		if i.InvoiceNumber == nil || *i.InvoiceNumber == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "invoice_number cannot be empty"
			resp["field"] = "invoice_number"
			json.NewEncoder(w).Encode(resp)
			return
		}

		if i.InvoiceDate == nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "invoice_date cannot be empty"
			resp["field"] = "invoice_date"
			json.NewEncoder(w).Encode(resp)
			return
		}

		err = s.DB.CreateInvoice(i)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(i)

	case http.MethodGet:
		invoices, err := s.DB.GetInvoices(ctx.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(invoices)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) OneInvoice(w http.ResponseWriter, r *http.Request) {
	var resp = make(map[string]string)
	ctx, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]
	invoiceId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	_ = invoiceId

	switch r.Method {
	case http.MethodDelete:
		// TODO: implement MethodDelete
		w.WriteHeader(http.StatusNotImplemented)

	case http.MethodPut:
		// TODO: implement MethodPut
		w.WriteHeader(http.StatusNotImplemented)

	case http.MethodGet:
		// TODO: implement MethodGet
		w.WriteHeader(http.StatusNotImplemented)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	_ = ctx
}
