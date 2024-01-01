package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetInvoiceDetails(invoiceId, companyId uuid.UUID) ([]types.InvoiceDetailResponse, error) {
	// TODO: implement get invoice details
	return []types.InvoiceDetailResponse{}, nil
}

func (s *service) CreateInvoiceDetail(detail *types.InvoiceDetail) (types.InvoiceDetail, error) {
	// TODO: implement create invoice details
	return types.InvoiceDetail{}, nil
}

func (s *service) UpdateInvoiceDetail(detail *types.InvoiceDetail) error {
	// TODO: implement update invoice details
	return nil
}

func (s *service) DeleteInvoiceDetail(detailId, companyId uuid.UUID) error {
	// TODO: implement delete invoice details
	return nil
}

func (s *service) GetOneInvoiceDetail(invoiceId, detailId, companyId uuid.UUID) (types.InvoiceDetailResponse, error) {
	// TODO: implement get invoice details
	return types.InvoiceDetailResponse{}, nil
}
