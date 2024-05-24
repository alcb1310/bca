package mocks

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) GetAllMaterials(companyId uuid.UUID) []types.Material {
	args := s.Called(companyId)
	return args.Get(0).([]types.Material)
}

func (s *ServiceMock) CreateMaterial(material types.Material) error {
	args := s.Called(material)
	return args.Error(0)
}

func (s *ServiceMock) GetMaterial(id, companyId uuid.UUID) (types.Material, error) {
	args := s.Called(id, companyId)
	return args.Get(0).(types.Material), args.Error(1)
}

func (s *ServiceMock) UpdateMaterial(material types.Material) error {
	args := s.Called(material)
	return args.Error(0)
}

func (s *ServiceMock) GetMaterialsByItem(id, companyId uuid.UUID) []types.ACU {
	args := s.Called(id, companyId)
	return args.Get(0).([]types.ACU)
}

func (s *ServiceMock) AddMaterialsByItem(itemId, materialId uuid.UUID, quantity float64, companyId uuid.UUID) error {
	args := s.Called(itemId, materialId, quantity, companyId)
	return args.Error(0)
}

func (s *ServiceMock) DeleteMaterialsByItem(itemId, materialId, companyId uuid.UUID) error {
	args := s.Called(itemId, materialId, companyId)
	return args.Error(0)
}

func (s *ServiceMock) UpdateMaterialByItem(itemId, materialId uuid.UUID, quantity float64, companyId uuid.UUID) error {
	args := s.Called(itemId, materialId, quantity, companyId)
	return args.Error(0)
}
