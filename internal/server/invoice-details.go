package server

import "net/http"

func (s *Server) AllInvoiceDetails(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// TODO: implement create invoice details
		w.WriteHeader(http.StatusNotImplemented)

	case http.MethodGet:
		// TODO: implement get all invoice details
		w.WriteHeader(http.StatusNotImplemented)

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
