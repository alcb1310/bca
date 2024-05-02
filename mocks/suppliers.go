package mocks

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) GetAllSuppliers(companyId uuid.UUID, search string) ([]types.Supplier, error) {
	args := s.Called(companyId)
	return args.Get(0).([]types.Supplier), args.Error(1)
}

func (s *ServiceMock) CreateSupplier(supplier *types.Supplier) error {
	args := s.Called(supplier)
	return args.Error(0)
}

func (s *ServiceMock) GetOneSupplier(id, companyId uuid.UUID) (types.Supplier, error) {
	args := s.Called(id, companyId)
	return args.Get(0).(types.Supplier), args.Error(1)
}

func (s *ServiceMock) UpdateSupplier(supplier *types.Supplier) error {
	args := s.Called(supplier)
	return args.Error(0)
}
