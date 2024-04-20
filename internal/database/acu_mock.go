package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) CreateCantidades(projectId, rubroId uuid.UUID, quantity float64, companyId uuid.UUID) error {
	return nil
}

func (s ServiceMock) DeleteCantidades(id, companyId uuid.UUID) error {
	return nil
}

func (s ServiceMock) CantidadesTable(companyId uuid.UUID) []types.Quantity {
	return []types.Quantity{}
}

func (s ServiceMock) AnalysisReport(project_id, company_id uuid.UUID) map[string][]types.AnalysisReport {
	return map[string][]types.AnalysisReport{}
}

func (s ServiceMock) GetQuantityByMaterialAndItem(itemId, materialId, companyId uuid.UUID) types.ItemMaterialType {
	return types.ItemMaterialType{}
}

func (s ServiceMock) GetOneQuantityById(id, companyId uuid.UUID) types.Quantity {
	return types.Quantity{}
}

func (s ServiceMock) UpdateQuantity(q types.Quantity, companyId uuid.UUID) error {
	return nil
}
