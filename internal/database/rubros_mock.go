package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetAllRubros(companyId uuid.UUID) ([]types.Rubro, error) {
	return []types.Rubro{}, nil
}

func (s ServiceMock) CreateRubro(rubro types.Rubro) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (s ServiceMock) GetOneRubro(id, companyId uuid.UUID) (types.Rubro, error) {
	return types.Rubro{}, nil
}

func (s ServiceMock) UpdateRubro(rubro types.Rubro) error {
	return nil
}

func (s ServiceMock) GetMaterialsByItem(id, companyId uuid.UUID) []types.ACU {
	return []types.ACU{}
}

func (s ServiceMock) AddMaterialsByItem(itemId, materialId uuid.UUID, quantity float64, companyId uuid.UUID) error {
	return nil
}

func (s ServiceMock) DeleteMaterialsByItem(itemId, materialId, companyId uuid.UUID) error {
	return nil
}

func (s ServiceMock) UpdateMaterialByItem(itemId, materialId uuid.UUID, quantity float64, companyId uuid.UUID) error {
	return nil
}
