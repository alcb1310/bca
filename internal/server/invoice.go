package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var resp = make(map[string]string)

func (s *Server) AllInvoices(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodPost:
		// TODO: implment MethodPost
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

func (s *Server) OneInvoice(w http.ResponseWriter, r *http.Request) {
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
