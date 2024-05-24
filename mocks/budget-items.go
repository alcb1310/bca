package mocks

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) GetBudgetItems(companyId uuid.UUID, search string) ([]types.BudgetItemResponse, error) {
	args := s.Called(companyId, search)
	return args.Get(0).([]types.BudgetItemResponse), args.Error(1)
}

func (s *ServiceMock) CreateBudgetItem(bi *types.BudgetItem) error {
	args := s.Called(bi)
	return args.Error(0)
}

func (s *ServiceMock) GetOneBudgetItem(id uuid.UUID, companyId uuid.UUID) (*types.BudgetItem, error) {
	args := s.Called(id, companyId)
	return args.Get(0).(*types.BudgetItem), args.Error(1)
}

func (s *ServiceMock) UpdateBudgetItem(bi *types.BudgetItem) error {
	args := s.Called(bi)
	return args.Error(0)
}

func (s *ServiceMock) GetBudgetItemsByAccumulate(companyId uuid.UUID, accum bool) []types.BudgetItem {
	args := s.Called(companyId, accum)
	return args.Get(0).([]types.BudgetItem)
}

func (s *ServiceMock) GetBudgetItemsByLevel(companyId uuid.UUID, level uint8) []types.BudgetItem {
	args := s.Called(companyId, level)
	return args.Get(0).([]types.BudgetItem)
}

func (s *ServiceMock) GetNonAccumulateChildren(companyId, id *uuid.UUID, budgetItems []types.BudgetItem, results []uuid.UUID) []uuid.UUID {
	args := s.Called(companyId, id, budgetItems, results)
	return args.Get(0).([]uuid.UUID)
}
