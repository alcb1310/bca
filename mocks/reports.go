package mocks

import (
	"time"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) GetBalance(companyId, projectId uuid.UUID, date time.Time) types.BalanceResponse {
	args := s.Called(companyId, projectId, date)
	return args.Get(0).(types.BalanceResponse)
}

func (s *ServiceMock) GetHistoricByProject(companyId, projectId uuid.UUID, date time.Time, level uint8) []types.GetBudget {
	args := s.Called(companyId, projectId, date, level)
	return args.Get(0).([]types.GetBudget)
}

func (s *ServiceMock) GetSpentByBudgetItem(companyId, projectId, budgetItemId uuid.UUID, date time.Time, ids []uuid.UUID) float64 {
	args := s.Called(companyId, projectId, budgetItemId, date, ids)
	return args.Get(0).(float64)
}

func (s *ServiceMock) GetDetailsByBudgetItem(companyId, projectId, budgetItemId uuid.UUID, date time.Time, ids []uuid.UUID) []types.InvoiceDetails {
	args := s.Called(companyId, projectId, budgetItemId, date, ids)
	return args.Get(0).([]types.InvoiceDetails)
}
