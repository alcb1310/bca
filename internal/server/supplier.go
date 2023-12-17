package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"encoding/json"
	"net/http"
)

func (s *Server) AllSuppliers(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctxPayload, _ := getMyPaload(r)

	switch r.Method {
	case http.MethodPost:
		supplier := &types.Supplier{}

		if err := json.NewDecoder(r.Body).Decode(supplier); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		if supplier.SupplierId == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "supplier_id cannot be empty"
			resp["field"] = "supplier_id"
			json.NewEncoder(w).Encode(resp)
			return
		}

		if supplier.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "name cannot be empty"
			resp["field"] = "name"
			json.NewEncoder(w).Encode(resp)
			return
		}

		if supplier.ContactEmail != nil && !utils.IsValidEmail(*supplier.ContactEmail) {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "invalid email"
			resp["field"] = "contact_email"
			json.NewEncoder(w).Encode(resp)
			return
		}

		supplier.CompanyId = ctxPayload.CompanyId

		if err := s.DB.CreateSupplier(supplier); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(supplier)

	case http.MethodGet:
		suppliers, err := s.DB.GetAllSuppliers(ctxPayload.CompanyId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(suppliers)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
