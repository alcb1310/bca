package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetBudgets(companyId uuid.UUID) ([]types.Budget, error) {
	query := "SELECT project_id, budget_item_id, initial_quantity, initial_cost, initial_total, spent_quantity, spent_total, remaining_quantity, remaining_cost, remaining_total, updated_budget, company_id FROM budget WHERE company_id = $1"
	rows, err := s.db.Query(query, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	budgets := []types.Budget{}

	for rows.Next() {
		b := types.Budget{}
		if err := rows.Scan(&b.ProjectId, &b.BudgetItemId, &b.InitialQuantity, &b.InitialCost, &b.InitialTotal, &b.SpentQuantity, &b.SpentTotal, &b.RemainingQuantity, &b.RemainingCost, &b.RemainingTotal, &b.UpdatedBudget, &b.CompanyId); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}

	return budgets, nil
}
