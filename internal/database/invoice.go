package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetInvoices(companyId uuid.UUID) ([]types.InvoiceResponse, error) {
	invoices := []types.InvoiceResponse{}
	query := `
		 select
			  id, supplier_id, supplier_number, supplier_name, supplier_contact_name, supplier_contact_email, supplier_contact_phone,
			  project_id, project_name, project_is_active,
			  invoice_number, invoice_date, invoice_total
		 from
			   vw_invoice
		 where
			   company_id = $1
	`

	rows, err := s.db.Query(query, companyId)
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
	query := "insert into invoice (company_id, supplier_id, project_id, invoice_number, invoice_date) values ($1, $2, $3, $4, $5) returning id"
	err := s.db.QueryRow(query, invoice.CompanyId, invoice.SupplierId, invoice.ProjectId, invoice.InvoiceNumber, invoice.InvoiceDate).Scan(&invoice.Id)
	return err
}

func (s *service) GetOneInvoice(invoiceId, companyId uuid.UUID) (types.InvoiceResponse, error) {
	// TODO: implement get one invoice method
	return types.InvoiceResponse{}, nil
}

func (s *service) UpdateInvoice(invoice types.InvoiceCreate) error {
	// TODO: implement update invoice method
	return nil
}
