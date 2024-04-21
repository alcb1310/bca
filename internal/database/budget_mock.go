package database

import (
	"errors"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetBudgets(companyId, project_id uuid.UUID, search string) ([]types.GetBudget, error) {
	return []types.GetBudget{}, nil
}
func (s ServiceMock) CreateBudget(b *types.CreateBudget) (types.Budget, error) {
	bId := uuid.MustParse("cc5cbcb9-43cc-4062-b3d3-ea60a3c2e6d0")
	pId := uuid.MustParse("bc39e850-0a1f-446f-a112-3e9a5b3134f0")
	if b.ProjectId == pId && b.BudgetItemId == bId {
		return types.Budget{}, errors.New("duplicate budget")
	}

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
