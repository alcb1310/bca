package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetBudgets(companyId, project_id uuid.UUID, search string) ([]types.GetBudget, error) {
	return []types.GetBudget{}, nil
}
func (s ServiceMock) CreateBudget(b *types.CreateBudget) (types.Budget, error) {
	return types.Budget{}, nil
}

func (s ServiceMock) GetBudgetsByProjectId(companyId, projectId uuid.UUID, level *uint8) ([]types.GetBudget, error) {
	return []types.GetBudget{}, nil
}

func (s ServiceMock) GetOneBudget(companyId, projectId, budgetItemId uuid.UUID) (*types.GetBudget, error) {
	return &types.GetBudget{}, nil
}

func (s ServiceMock) UpdateBudget(b types.CreateBudget, budget types.Budget) error {
	return nil
}
