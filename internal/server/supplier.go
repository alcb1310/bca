package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/utils"
	"bca-go-final/internal/views/bca/settings/partials"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

func (s *Server) OneSupplier(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	ctxPayload, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]

	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	supplier, err := s.DB.GetOneSupplier(parsedId, ctxPayload.CompanyId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		resp["error"] = err.Error()
		json.NewEncoder(w).Encode(resp)
		return
	}

	switch r.Method {
	case http.MethodPut:
		sup := &types.Supplier{}

		if err := json.NewDecoder(r.Body).Decode(sup); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		if sup.SupplierId == "" {
			sup.SupplierId = supplier.SupplierId
		}
		if sup.Name == "" {
			sup.Name = supplier.Name
		}
		if sup.ContactEmail == nil {
			sup.ContactEmail = supplier.ContactEmail
		} else if !utils.IsValidEmail(*sup.ContactEmail) {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = "invalid email"
			resp["field"] = "contact_email"
			json.NewEncoder(w).Encode(resp)
			return
		}
		if sup.ContactPhone == nil {
			sup.ContactPhone = supplier.ContactPhone
		}
		if sup.ContactName == nil {
			sup.ContactName = supplier.ContactName
		}
		sup.CompanyId = ctxPayload.CompanyId
		sup.ID = parsedId

		if err := s.DB.UpdateSupplier(sup); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sup)

	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(supplier)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) SuppliersTable(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)

	if r.Method == http.MethodPost {
		r.ParseForm()
		email := r.Form.Get("contact_email")
		name := r.Form.Get("contact_name")
		phone := r.Form.Get("contact_phone")
		sup := types.Supplier{
			SupplierId:   r.Form.Get("supplier_id"),
			Name:         r.Form.Get("name"),
			ContactEmail: &email,
			ContactName:  &name,
			ContactPhone: &phone,
			CompanyId:    ctxPayload.CompanyId,
		}

		err := s.DB.CreateSupplier(&sup)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	suppliers, _ := s.DB.GetAllSuppliers(ctxPayload.CompanyId)
	component := partials.SuppliersTable(suppliers)
	component.Render(r.Context(), w)
}

func (s *Server) SupplierAdd(w http.ResponseWriter, r *http.Request) {
	component := partials.EditSupplier(nil)
	component.Render(r.Context(), w)
}
