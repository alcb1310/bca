package mocks

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) GetBudgets(companyId, project_id uuid.UUID, search string) ([]types.GetBudget, error) {
	args := s.Called(companyId, project_id, search)
	return args.Get(0).([]types.GetBudget), args.Error(1)
}

func (s *ServiceMock) CreateBudget(b *types.CreateBudget) (types.Budget, error) {
	args := s.Called(b)
	return args.Get(0).(types.Budget), args.Error(1)
}

func (s *ServiceMock) GetBudgetsByProjectId(companyId, projectId uuid.UUID, level *uint8) ([]types.GetBudget, error) {
	args := s.Called(companyId, projectId, level)
	return args.Get(0).([]types.GetBudget), args.Error(1)
}

func (s *ServiceMock) GetOneBudget(companyId, projectId, budgetItemId uuid.UUID) (*types.GetBudget, error) {
	args := s.Called(companyId, projectId, budgetItemId)
	return args.Get(0).(*types.GetBudget), args.Error(1)
}

func (s *ServiceMock) UpdateBudget(b types.CreateBudget, budget types.Budget) error {
	args := s.Called(b, budget)
	return args.Error(0)
}
