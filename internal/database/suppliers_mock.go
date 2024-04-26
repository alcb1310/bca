package database

import (
	"errors"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetAllSuppliers(companyId uuid.UUID, search string) ([]types.Supplier, error) {
	return []types.Supplier{}, nil
}

func (s ServiceMock) CreateSupplier(supplier *types.Supplier) error {
	if supplier.SupplierId == "0123456789" {
		return errors.New("duplicate key value violates unique constraint \"supplier_supplier_id_key\"")
	}

	if supplier.Name == "exists" {
		return errors.New("conflict duplicate key value violates unique constraint \"supplier_name_key\"")
	}
	return nil
}

func (s ServiceMock) GetOneSupplier(id, companyId uuid.UUID) (types.Supplier, error) {
	return types.Supplier{}, nil
}

func (s ServiceMock) UpdateSupplier(supplier *types.Supplier) error {
	return nil
}
