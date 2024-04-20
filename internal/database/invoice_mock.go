package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetInvoices(companyId uuid.UUID) ([]types.InvoiceResponse, error) {
	return []types.InvoiceResponse{}, nil
}

func (s ServiceMock) CreateInvoice(invoice *types.InvoiceCreate) error {
	return nil
}

func (s ServiceMock) GetOneInvoice(id, companyId uuid.UUID) (types.InvoiceResponse, error) {
	return types.InvoiceResponse{}, nil
}

func (s ServiceMock) UpdateInvoice(invoice types.InvoiceCreate) error {
	return nil
}

func (s ServiceMock) DeleteInvoice(id, companyId uuid.UUID) error {
	return nil
}

func (s ServiceMock) BalanceInvoice(invoice types.InvoiceResponse) error {
	return nil
}
