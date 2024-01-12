package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/settings/partials"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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

func (s *Server) SuppliersEdit(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sup, _ := s.DB.GetOneSupplier(parsedId, ctxPayload.CompanyId)
	component := partials.EditSupplier(&sup)
	component.Render(r.Context(), w)
}

func (s *Server) SuppliersEditSave(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := getMyPaload(r)
	id := mux.Vars(r)["id"]
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sup, _ := s.DB.GetOneSupplier(parsedId, ctxPayload.CompanyId)
	r.ParseForm()
	sup.SupplierId = r.Form.Get("supplier_id")
	sup.Name = r.Form.Get("name")
	email := r.Form.Get("contact_email")
	name := r.Form.Get("contact_name")
	phone := r.Form.Get("contact_phone")
	sup.ContactEmail = &email
	sup.ContactName = &name
	sup.ContactPhone = &phone

	if err := s.DB.UpdateSupplier(&sup); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	suppliers, _ := s.DB.GetAllSuppliers(ctxPayload.CompanyId)
	component := partials.SuppliersTable(suppliers)
	component.Render(r.Context(), w)
}
