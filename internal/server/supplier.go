package server

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
	"github.com/alcb1310/bca/internal/views/bca/settings/partials"
)

func (s *Server) SuppliersTable(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)

	if r.Method == http.MethodPost {
		r.ParseForm()
		e := r.Form.Get("contact_email")
		if e != "" && !utils.IsValidEmail(e) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese un correo válido"))
			return
		}

		email := sql.NullString{Valid: true, String: e}
		n := r.Form.Get("contact_name")
		name := sql.NullString{Valid: true, String: n}
		p := r.Form.Get("contact_phone")
		phone := sql.NullString{Valid: true, String: p}
		sup := types.Supplier{
			SupplierId:   r.Form.Get("supplier_id"),
			Name:         r.Form.Get("name"),
			ContactEmail: email,
			ContactName:  name,
			ContactPhone: phone,
			CompanyId:    ctxPayload.CompanyId,
		}

		if sup.SupplierId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese un valor para el RUC"))
			return
		}

		if sup.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese un valor para el nombre"))
			return
		}

		err := s.DB.CreateSupplier(&sup)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(fmt.Sprintf("Proveedor con ruc %s y/o nombre %s ya existe", sup.SupplierId, sup.Name)))
				return
			}
			slog.Error("Error creating supplier", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("<p>%s</p>", err.Error())))
			return
		}
	}

	search := r.URL.Query().Get("search")
	suppliers, _ := s.DB.GetAllSuppliers(ctxPayload.CompanyId, search)
	component := partials.SuppliersTable(suppliers)
	component.Render(r.Context(), w)
}

func (s *Server) SupplierAdd(w http.ResponseWriter, r *http.Request) {
	component := partials.EditSupplier(nil)
	component.Render(r.Context(), w)
}

func (s *Server) SuppliersEdit(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sup, err := s.DB.GetOneSupplier(parsedId, ctxPayload.CompanyId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Proveedor no encontrado"))
		return
	}

	component := partials.EditSupplier(&sup)
	component.Render(r.Context(), w)
}

func (s *Server) SuppliersEditSave(w http.ResponseWriter, r *http.Request) {
	ctxPayload, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sup, _ := s.DB.GetOneSupplier(parsedId, ctxPayload.CompanyId)
	r.ParseForm()
	sp := r.Form.Get("supplier_id")
	if sp != "" {
		sup.SupplierId = sp
	}

	nm := r.Form.Get("name")
	if nm != "" {
		sup.Name = nm
	}

	e := r.Form.Get("contact_email")
	if e != "" && !utils.IsValidEmail(e) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese un correo válido"))
		return
	}
	email := sql.NullString{Valid: true, String: e}

	n := r.Form.Get("contact_name")
	name := sql.NullString{Valid: true, String: n}

	p := r.Form.Get("contact_phone")
	phone := sql.NullString{Valid: true, String: p}

	sup.ContactEmail = email
	sup.ContactName = name
	sup.ContactPhone = phone

	if err := s.DB.UpdateSupplier(&sup); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(fmt.Sprintf("El ruc %s y/o nombre %s ya existe", sup.SupplierId, sup.Name)))
			return
		}
		log.Println(fmt.Sprintf("ERROR: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("<p>%s</p>", err.Error())))
		return
	}

	suppliers, _ := s.DB.GetAllSuppliers(ctxPayload.CompanyId, "")
	component := partials.SuppliersTable(suppliers)
	component.Render(r.Context(), w)
}
