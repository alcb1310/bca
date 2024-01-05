package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) AllInvoiceDetails(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctx, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]
	invoiceId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	switch r.Method {
	case http.MethodPost:
		// TODO: implement create invoice details
		w.WriteHeader(http.StatusNotImplemented)

	case http.MethodGet:
		// TODO: implement get all invoice details
		invoiceDetails, err := s.DB.GetInvoiceDetails(invoiceId, ctx.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(invoiceDetails)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) OneInvoiceDetails(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		// TODO: implement delete invoice details
		w.WriteHeader(http.StatusNotImplemented)

	case http.MethodPut:
		// TODO: implement update invoice details
		w.WriteHeader(http.StatusNotImplemented)

	case http.MethodGet:
		// TODO: implement get one invoice details
		w.WriteHeader(http.StatusNotImplemented)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
