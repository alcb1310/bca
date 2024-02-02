package server

import (
	"bca-go-final/internal/types"
	"bca-go-final/internal/views/bca/transaction/partials"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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

	sx, _ := s.DB.GetAllSuppliers(ctx.CompanyId, "")
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
		if pId == "" {
			w.WriteHeader(http.StatusBadRequest)
      w.Header().Set("HX-Retarget", "#error")
			w.Write([]byte("Ingrese un proyecto"))
			return
		}
		projectId, _ := uuid.Parse(pId)
		sId := r.Form.Get("supplier")
		if sId == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ingrese un proveedor"))
			return
		}
		supplierId, _ := uuid.Parse(sId)
		iNumber := r.Form.Get("invoiceNumber")
    if iNumber == "" {
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("Ingrese un número de factura"))
      return
    }
		iD := r.Form.Get("invoiceDate")
		iDate, err := time.Parse("2006-01-02", iD)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("Ingrese una fecha válida"))
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
    if err != nil {
      if strings.Contains(err.Error(), "duplicate") {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("La Factura ya existe"))
        return
      }
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte(err.Error()))
      return
    }
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

	sx, _ := s.DB.GetAllSuppliers(ctx.CompanyId, "")
	for _, v := range sx {
		x := types.Select{
			Key:   v.ID.String(),
			Value: v.Name,
		}
		suppliers = append(suppliers, x)
	}

	switch r.Method {
	case http.MethodPatch:
		if err := s.DB.BalanceInvoice(in); err != nil {
			log.Printf("error updating invoice: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		in, _ := s.DB.GetOneInvoice(parsedId, ctx.CompanyId)

		comp := partials.BudgetRow(in)
		comp.Render(r.Context(), w)
		return

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
