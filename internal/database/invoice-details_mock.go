package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetAllDetails(invoiceId, companyId uuid.UUID) ([]types.InvoiceDetailsResponse, error) {
	return []types.InvoiceDetailsResponse{}, nil
}

func (s ServiceMock) DeleteDetail(invoiceId, budgetItemId, companyId uuid.UUID) error {
	return nil
}

func (s ServiceMock) AddDetail(detail types.InvoiceDetailCreate) error {
	return nil
}
