package mocks

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) GetAllDetails(invoiceId, companyId uuid.UUID) ([]types.InvoiceDetailsResponse, error) {
	args := s.Called(invoiceId, companyId)
	return args.Get(0).([]types.InvoiceDetailsResponse), args.Error(1)
}

func (s *ServiceMock) AddDetail(detail types.InvoiceDetailCreate) error {
	args := s.Called(detail)
	return args.Error(0)
}

func (s *ServiceMock) DeleteDetail(invoiceId, budgetItemId, companyId uuid.UUID) error {
	args := s.Called(invoiceId, budgetItemId, companyId)
	return args.Error(0)
}
