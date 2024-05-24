package mocks

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) CreateCantidades(projectId, rubroId uuid.UUID, quantity float64, companyId uuid.UUID) error {
	args := s.Called(projectId, rubroId, quantity, companyId)
	return args.Error(0)
}

func (s *ServiceMock) DeleteCantidades(id, companyId uuid.UUID) error {
	args := s.Called(id, companyId)
	return args.Error(0)
}

func (s *ServiceMock) CantidadesTable(companyId uuid.UUID) []types.Quantity {
	args := s.Called(companyId)
	return args.Get(0).([]types.Quantity)
}

func (s *ServiceMock) AnalysisReport(project_id, company_id uuid.UUID) map[string][]types.AnalysisReport {
	args := s.Called(project_id, company_id)
	return args.Get(0).(map[string][]types.AnalysisReport)
}

func (s *ServiceMock) GetQuantityByMaterialAndItem(itemId, materialId, companyId uuid.UUID) types.ItemMaterialType {
	args := s.Called(itemId, materialId, companyId)
	return args.Get(0).(types.ItemMaterialType)
}

func (s *ServiceMock) GetOneQuantityById(id, companyId uuid.UUID) types.Quantity {
	args := s.Called(id, companyId)
	return args.Get(0).(types.Quantity)
}

func (s *ServiceMock) UpdateQuantity(q types.Quantity, companyId uuid.UUID) error {
	args := s.Called(q, companyId)
	return args.Error(0)
}
