package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetInvoices(companyId uuid.UUID) ([]types.InvoiceResponse, error) {
	// TODO: Implement get all invoices method
	return []types.InvoiceResponse{}, nil
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
