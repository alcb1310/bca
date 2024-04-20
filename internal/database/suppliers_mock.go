package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetAllSuppliers(companyId uuid.UUID, search string) ([]types.Supplier, error) {
	return []types.Supplier{}, nil
}

func (s ServiceMock) CreateSupplier(supplier *types.Supplier) error {
	return nil
}

func (s ServiceMock) GetOneSupplier(id, companyId uuid.UUID) (types.Supplier, error) {
	return types.Supplier{}, nil
}

func (s ServiceMock) UpdateSupplier(supplier *types.Supplier) error {
	return nil
}
