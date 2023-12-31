package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetInvoices(companyId uuid.UUID) ([]types.InvoiceResponse, error) {
	// TODO: Implement get all invoices method
	return []types.InvoiceResponse{}, nil
}

func (s *service) CreateInvoice(invoice types.InvoiceCreate) error {
	// TODO: implement create invoice method
	return nil
}

func (s *service) GetOneInvoice(invoiceId, companyId uuid.UUID) (types.InvoiceResponse, error) {
	// TODO: implement get one invoice method
	return types.InvoiceResponse{}, nil
}

func (s *service) UpdateInvoice(invoice types.InvoiceCreate) error {
	// TODO: implement update invoice method
	return nil
}
