package server

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
	"github.com/alcb1310/bca/internal/utils"
	"github.com/alcb1310/bca/internal/views/bca/transaction/partials"
)

func (s *Server) InvoicesTable(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)

	invoices, _ := s.DB.GetInvoices(ctx.CompanyId)
	components := partials.InvoiceTable(invoices)
	components.Render(r.Context(), w)
}

func (s *Server) InvoiceAdd(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	var invoice types.InvoiceResponse
	redirectURL := "/bca/transacciones/facturas/crear"

	projects := []types.Select{}
	suppliers := []types.Select{}

	id := r.URL.Query().Get("id")
	if id != "" {
		parsedId, _ := uuid.Parse(id)
		invoice, _ = s.DB.GetOneInvoice(parsedId, ctx.CompanyId)
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

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	pId := r.Form.Get("project")
	if pId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese un proyecto"))
		return
	}
	projectId, err := uuid.Parse(pId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Código del proyecto inválido"))
		return
	}

	sId := r.Form.Get("supplier")
	if sId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese un proveedor"))
		return
	}
	supplierId, err := uuid.Parse(sId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Código del proveedor inválido"))
		return
	}

	iNumber := r.Form.Get("invoiceNumber")
	if iNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese un número de factura"))
		return
	}

	iD := r.Form.Get("invoiceDate")
	if iD == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese una fecha"))
		return
	}

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
	invoice, _ = s.DB.GetOneInvoice(*i.Id, ctx.CompanyId)
	redirectURL += "?id=" + invoice.Id.String()

	components := partials.EditInvoice(&invoice, projects, suppliers)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
	components.Render(r.Context(), w)
}

func (s *Server) InvoiceAddForm(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
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

	components := partials.EditInvoice(invoice, projects, suppliers)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
	components.Render(r.Context(), w)
}

func (s *Server) InvoiceEdit(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	redirectURL := "/bca/transacciones/facturas/crear"
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)

	projects := []types.Select{}
	suppliers := []types.Select{}
	invoice, _ := s.DB.GetOneInvoice(parsedId, ctx.CompanyId)

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

	r.ParseForm()
	pId := invoice.Project.ID

	sId, err := uuid.Parse(r.Form.Get("supplier"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Código del proveedor inválido"))
		return
	}
	iNumber := r.Form.Get("invoiceNumber")
	if iNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese un número de factura"))
		return
	}

	fDate := r.Form.Get("invoiceDate")
	if fDate == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese una fecha"))
		return
	}

	iDate, err := time.Parse("2006-01-02", fDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Ingrese una fecha válida"))
		return
	}

	i := types.InvoiceCreate{
		CompanyId:     ctx.CompanyId,
		ProjectId:     &pId,
		SupplierId:    &sId,
		InvoiceNumber: &iNumber,
		InvoiceDate:   &iDate,
		Id:            &parsedId,
	}

	if err := s.DB.UpdateInvoice(i); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("La Factura ya existe"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		slog.Error("Error updating invoice", "error", err)
		return
	}

	invoice, _ = s.DB.GetOneInvoice(parsedId, ctx.CompanyId)
	redirectURL += "?id=" + invoice.Id.String()

	components := partials.EditInvoice(&invoice, projects, suppliers)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
	components.Render(r.Context(), w)
}

func (s *Server) GetOneInvoice(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	redirectURL := "/bca/transacciones/facturas/crear"
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)

	projects := []types.Select{}
	suppliers := []types.Select{}
	invoice, _ := s.DB.GetOneInvoice(parsedId, ctx.CompanyId)

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

	components := partials.EditInvoice(&invoice, projects, suppliers)
	w.Header().Set("HX-Redirect", redirectURL)
	w.WriteHeader(http.StatusOK)
	components.Render(r.Context(), w)
}

func (s *Server) DeleteInvoice(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)

	_, err := s.DB.GetOneInvoice(parsedId, ctx.CompanyId)
	if err != nil {
		slog.Error("Error getting invoice", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.DB.DeleteInvoice(parsedId, ctx.CompanyId); err != nil {
		slog.Error("Error deleting invoice", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	invoices, _ := s.DB.GetInvoices(ctx.CompanyId)

	components := partials.InvoiceTable(invoices)
	w.WriteHeader(http.StatusOK)
	components.Render(r.Context(), w)
}

func (s *Server) PatchInvoice(w http.ResponseWriter, r *http.Request) {
	ctx, _ := utils.GetMyPaload(r)
	id := chi.URLParam(r, "id")
	parsedId, _ := uuid.Parse(id)

	invoice, _ := s.DB.GetOneInvoice(parsedId, ctx.CompanyId)

	if err := s.DB.BalanceInvoice(invoice); err != nil {
		slog.Error("Error updating invoice", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	in, _ := s.DB.GetOneInvoice(parsedId, ctx.CompanyId)

	comp := partials.BudgetRow(in)
	comp.Render(r.Context(), w)
}
