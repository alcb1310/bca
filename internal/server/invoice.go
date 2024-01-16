package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/transaction/partials"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

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

	invoice, err := s.DB.GetOneInvoice(invoiceId, ctx.CompanyId)
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
	case http.MethodDelete:
		if invoice.InvoiceTotal != 0 {
			w.WriteHeader(http.StatusNotAcceptable)
			resp["error"] = "invoice is not empty"
			json.NewEncoder(w).Encode(resp)
			return
		}

		if err := s.DB.DeleteInvoice(invoiceId, ctx.CompanyId); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	case http.MethodPut:
		i := types.InvoiceCreate{}

		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		i.Id = &invoice.Id
		i.CompanyId = ctx.CompanyId
		if i.SupplierId == nil || *i.SupplierId == uuid.Nil {
			i.SupplierId = &invoice.Supplier.ID
		}
		if i.ProjectId == nil || *i.ProjectId == uuid.Nil {
			i.ProjectId = &invoice.Project.ID
		}
		if i.InvoiceNumber == nil || *i.InvoiceNumber == "" {
			i.InvoiceNumber = &invoice.InvoiceNumber
		}
		if i.InvoiceDate == nil {
			i.InvoiceDate = &invoice.InvoiceDate
		}

		err := s.DB.UpdateInvoice(i)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp["error"] = err.Error()
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)

	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(invoice)

	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) InvoicesTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)

	invoices, _ := s.DB.GetInvoices(ctx.CompanyId)
	components := partials.InvoiceTable(invoices)
	components.Render(r.Context(), w)
}

func (s *Server) InvoiceAdd(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	var invoice *types.InvoiceResponse
	redirectURL := "/bca/transacciones/facturas/crear"
	invoice = nil

	projects := []types.Select{}
	suppliers := []types.Select{}

	id := r.URL.Query().Get("id")
	if id != "" {
		parsedId, _ := uuid.Parse(id)
		in, _ := s.DB.GetOneInvoice(parsedId, ctx.CompanyId)
		invoice = &in
	}

	p := s.DB.GetActiveProjects(ctx.CompanyId, true)
	for _, v := range p {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		projects = append(projects, x)
	}

	sx, _ := s.DB.GetAllSuppliers(ctx.CompanyId)
	for _, v := range sx {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		suppliers = append(suppliers, x)
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		pId := r.Form.Get("project")
		projectId, _ := uuid.Parse(pId)
		sId := r.Form.Get("supplier")
		supplierId, _ := uuid.Parse(sId)
		iNumber := r.Form.Get("invoiceNumber")
		iD := r.Form.Get("invoiceDate")
		iDate, err := time.Parse("2006-01-02", iD)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("error parsing date: %v", err)
			return
		}

		i := &types.InvoiceCreate{
			CompanyId:     ctx.CompanyId,
			ProjectId:     &projectId,
			SupplierId:    &supplierId,
			InvoiceNumber: &iNumber,
			InvoiceDate:   &iDate,
		}

		err = s.DB.CreateInvoice(i)
		in, _ := s.DB.GetOneInvoice(*i.Id, ctx.CompanyId)
		invoice = &in
		redirectURL += "?id=" + in.Id.String()
	}

	components := partials.EditInvoice(invoice, projects, suppliers)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
	components.Render(r.Context(), w)
}

func (s *Server) InvoiceEdit(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getMyPaload(r)
	redirectURL := "/bca/transacciones/facturas/crear"
	id := mux.Vars(r)["id"]
	parsedId, _ := uuid.Parse(id)
	invoice := &types.InvoiceResponse{}

	projects := []types.Select{}
	suppliers := []types.Select{}
	// projects := make(map[string]string)
	// suppliers := make(map[string]string)
	in, _ := s.DB.GetOneInvoice(parsedId, ctx.CompanyId)
	invoice = &in

	p := s.DB.GetActiveProjects(ctx.CompanyId, true)
	for _, v := range p {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		projects = append(projects, x)
	}

	sx, _ := s.DB.GetAllSuppliers(ctx.CompanyId)
	for _, v := range sx {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		suppliers = append(suppliers, x)
	}

	switch r.Method {
	case http.MethodPut:
		r.ParseForm()
		sId, _ := uuid.Parse(r.Form.Get("supplier"))
		pId := invoice.Project.ID
		iNumber := r.Form.Get("invoiceNumber")
		iDate, _ := time.Parse("2006-01-02", r.Form.Get("invoiceDate"))

		i := types.InvoiceCreate{
			CompanyId:     ctx.CompanyId,
			ProjectId:     &pId,
			SupplierId:    &sId,
			InvoiceNumber: &iNumber,
			InvoiceDate:   &iDate,
			Id:            &parsedId,
		}

		err := s.DB.UpdateInvoice(i)
		if err != nil {
			log.Printf("error updating invoice: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		in, _ := s.DB.GetOneInvoice(parsedId, ctx.CompanyId)
		invoice = &in
		redirectURL += "?id=" + in.Id.String()
	}

	components := partials.EditInvoice(invoice, projects, suppliers)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
	components.Render(r.Context(), w)
}
