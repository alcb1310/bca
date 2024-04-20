package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetAllMaterials(companyId uuid.UUID) []types.Material {
	return []types.Material{}
}

func (s ServiceMock) CreateMaterial(material types.Material) error {
	return nil
}

func (s ServiceMock) GetMaterial(id, companyId uuid.UUID) (types.Material, error) {
	return types.Material{}, nil
}

func (s ServiceMock) UpdateMaterial(material types.Material) error {
	return nil
}
