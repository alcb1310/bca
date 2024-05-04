package mocks

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) GetInvoices(companyId uuid.UUID) ([]types.InvoiceResponse, error) {
	args := s.Called(companyId)
	return args.Get(0).([]types.InvoiceResponse), args.Error(1)
}

func (s *ServiceMock) CreateInvoice(invoice *types.InvoiceCreate) error {
	args := s.Called(invoice)
	id := uuid.MustParse("cdefa321-9f2d-4673-9949-7cac744e941a")
	invoice.Id = &id
	return args.Error(0)
}

func (s *ServiceMock) GetOneInvoice(invoiceId, companyId uuid.UUID) (types.InvoiceResponse, error) {
	args := s.Called(invoiceId, companyId)
	return args.Get(0).(types.InvoiceResponse), args.Error(1)
}

func (s *ServiceMock) UpdateInvoice(invoice types.InvoiceCreate) error {
	args := s.Called(invoice)
	return args.Error(0)
}

func (s *ServiceMock) DeleteInvoice(invoiceId, companyId uuid.UUID) error {
	args := s.Called(invoiceId, companyId)
	return args.Error(0)
}

func (s *ServiceMock) BalanceInvoice(invoice types.InvoiceResponse) error {
	args := s.Called(invoice)
	return args.Error(0)
}
