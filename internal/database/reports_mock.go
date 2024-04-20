package database

import (
	"time"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetBalance(companyId, projectId uuid.UUID, date time.Time) types.BalanceResponse {
	return types.BalanceResponse{}
}

func (s ServiceMock) GetHistoricByProject(companyId, projectId uuid.UUID, date time.Time, level uint8) []types.GetBudget {
	return []types.GetBudget{}
}

func (s ServiceMock) GetSpentByBudgetItem(companyId, projectId, budgetItemId uuid.UUID, date time.Time, ids []uuid.UUID) float64 {
	return 0
}

func (s ServiceMock) GetDetailsByBudgetItem(companyId, projectId, budgetItemId uuid.UUID, date time.Time, ids []uuid.UUID) []types.InvoiceDetails {
	return []types.InvoiceDetails{}
}
