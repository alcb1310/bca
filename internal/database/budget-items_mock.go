package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetBudgetItems(companyId uuid.UUID, search string) ([]types.BudgetItemResponse, error) {
	return []types.BudgetItemResponse{}, nil
}

func (s ServiceMock) CreateBudgetItem(bi *types.BudgetItem) error {
	return nil
}

func (s ServiceMock) GetOneBudgetItem(id uuid.UUID, companyId uuid.UUID) (*types.BudgetItem, error) {
	return &types.BudgetItem{}, nil
}

func (s ServiceMock) UpdateBudgetItem(bi *types.BudgetItem) error {
	return nil
}

func (s ServiceMock) GetBudgetItemsByAccumulate(companyId uuid.UUID, accum bool) []types.BudgetItem {
	return []types.BudgetItem{}
}

func (s ServiceMock) GetBudgetItemsByLevel(companyId uuid.UUID, level uint8) []types.BudgetItem {
	return []types.BudgetItem{}
}

func (s ServiceMock) GetNonAccumulateChildren(companyId, id *uuid.UUID, budgetItems []types.BudgetItem, results []uuid.UUID) []uuid.UUID {
	return []uuid.UUID{}
}
