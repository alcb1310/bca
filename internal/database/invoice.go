package database

import (
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
)

func (s *service) GetInvoices(companyId uuid.UUID) ([]types.InvoiceResponse, error) {
	invoices := []types.InvoiceResponse{}
	query := `
		 select
			  id, supplier_id, supplier_number, supplier_name, supplier_contact_name, supplier_contact_email, supplier_contact_phone,
			  project_id, project_name, project_is_active,
			  invoice_number, invoice_date, invoice_total, is_balanced
		 from
			   vw_invoice
		 where
			   company_id = $1 and is_balanced = $2
		 order by
			   invoice_date desc, supplier_name, invoice_number
	`

	rows, err := s.db.Query(query, companyId, false)
	if err != nil {
		return invoices, err
	}
	defer rows.Close()

	for rows.Next() {
		invoice := types.InvoiceResponse{}
		if err := rows.Scan(
			&invoice.Id,
			&invoice.Supplier.ID,
			&invoice.Supplier.SupplierId,
			&invoice.Supplier.Name,
			&invoice.Supplier.ContactName,
			&invoice.Supplier.ContactEmail,
			&invoice.Supplier.ContactPhone,
			&invoice.Project.ID,
			&invoice.Project.Name,
			&invoice.Project.IsActive,
			&invoice.InvoiceNumber,
			&invoice.InvoiceDate,
			&invoice.InvoiceTotal,
			&invoice.IsBalanced,
		); err != nil {
			return invoices, err
		}
		invoice.CompanyId = companyId
		invoice.Supplier.CompanyId = companyId
		invoice.Project.CompanyId = companyId
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (s *service) CreateInvoice(invoice *types.InvoiceCreate) error {
	var cl, closure_date time.Time

	query := "select last_closure from project where id = $1 and company_id = $2"

	err := s.db.QueryRow(query, invoice.ProjectId, invoice.CompanyId).Scan(&cl)
	if err == nil {
		closure_date = time.Date(cl.Year(), cl.Month()+1, 1-1, 1, 0, 0, 0, time.UTC)

		if invoice.InvoiceDate.Before(closure_date) {
			return errors.New("La fecha indicada es menor al último cierre")
		}
	}

	query = "insert into invoice (company_id, supplier_id, project_id, invoice_number, invoice_date) values ($1, $2, $3, $4, $5) returning id"
	err = s.db.QueryRow(query, invoice.CompanyId, invoice.SupplierId, invoice.ProjectId, invoice.InvoiceNumber, invoice.InvoiceDate).Scan(&invoice.Id)
	return err
}

func (s *service) GetOneInvoice(invoiceId, companyId uuid.UUID) (types.InvoiceResponse, error) {
	i := &types.InvoiceResponse{}
	query := `
		 select
			  id, supplier_id, supplier_number, supplier_name, supplier_contact_name, supplier_contact_email, supplier_contact_phone,
			  project_id, project_name, project_is_active,
			  invoice_number, invoice_date, invoice_total, is_balanced
		 from
			   vw_invoice
		 where
			   company_id = $1 and id = $2
	`
	err := s.db.QueryRow(query, companyId, invoiceId).Scan(
		&i.Id, &i.Supplier.ID, &i.Supplier.SupplierId, &i.Supplier.Name, &i.Supplier.ContactName, &i.Supplier.ContactEmail, &i.Supplier.ContactPhone,
		&i.Project.ID, &i.Project.Name, &i.Project.IsActive,
		&i.InvoiceNumber, &i.InvoiceDate, &i.InvoiceTotal, &i.IsBalanced,
	)
	i.CompanyId = companyId
	i.Supplier.CompanyId = companyId
	i.Project.CompanyId = companyId

	return *i, err
}

func (s *service) UpdateInvoice(invoice types.InvoiceCreate) error {
	// Can not update invoice if it is already balanced
	if invoice.IsBalanced {
		return errors.New("La factura ya se encuentra cuadrada")
	}

	// Can not update invoice date to a date before last_closure
	var cl, closure_date time.Time

	query := "select last_closure from project where id = $1 and company_id = $2"

	err := s.db.QueryRow(query, invoice.ProjectId, invoice.CompanyId).Scan(&cl)
	if err == nil {
		closure_date = time.Date(cl.Year(), cl.Month()+1, 1-1, 1, 0, 0, 0, time.UTC)

		if invoice.InvoiceDate.Before(closure_date) {
			return errors.New("La fecha indicada es menor al último cierre")
		}
	}

	query = `
		update invoice
		set supplier_id = $1, project_id = $2, invoice_number = $3, invoice_date = $4
		where id = $5 and company_id = $6
	`
	_, err = s.db.Exec(query, invoice.SupplierId, invoice.ProjectId, invoice.InvoiceNumber, invoice.InvoiceDate, invoice.Id, invoice.CompanyId)

	return err
}

func (s *service) DeleteInvoice(invoiceId, companyId uuid.UUID) error {
	var total float64
	isBalanced := false

	query := "select invoice_total, is_balanced from invoice where id = $1 and company_id = $2"
	if err := s.db.QueryRow(query, invoiceId, companyId).Scan(&total, &isBalanced); err != nil {
		return errors.New("La Factura no existe")
	}

	if isBalanced {
		return errors.New("No se puede borrar una Factura cuadrada")
	}

	if total > 0 {
		return errors.New("No se puede borrar una Factura que ya tiene pagos")
	}

	query = `
		delete from invoice
		where id = $1 and company_id = $2
	`
	_, err := s.db.Exec(query, invoiceId, companyId)
	return err
}

func (s *service) BalanceInvoice(invoice types.InvoiceResponse) error {
	var bal bool
	var num string
	query := "select is_balanced, invoice_number from invoice where id = $1 and company_id = $2"

	err := s.db.QueryRow(query, invoice.Id, invoice.CompanyId).Scan(&bal, &num)
	if err != nil {
		return err
	}

	newBal := !bal

	slog.Error("balance", "newBal", newBal, "bal", bal, "num", num)

	query = `
		update invoice
		set is_balanced = $3 
		where id = $1 and company_id = $2
	`
	_, err = s.db.Exec(query, invoice.Id, invoice.CompanyId, newBal)
	slog.Error("balance", "err", err)
	return err
}
