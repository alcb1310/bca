package database

import (
	"bca-go-final/internal/types"
	"database/sql"

	"github.com/google/uuid"
)

func (s *service) GetBudgets(companyId uuid.UUID) ([]types.GetBudget, error) {
	query := `
        SELECT
            project_id, project_name,
            budget_item_id, budget_item_code, budget_item_name, budget_item_level, budget_item_accumulate,
            initial_quantity, initial_cost, initial_total,
            spent_quantity, spent_total,
            remaining_quantity, remaining_cost, remaining_total,
            updated_budget, company_id
        FROM vw_budget
        WHERE company_id = $1
        ORDER BY project_name, budget_item_code
    `
	rows, err := s.db.Query(query, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	budgets := []types.GetBudget{}

	for rows.Next() {
		b := types.GetBudget{}
		if err := rows.Scan(
			&b.Project.ID, &b.Project.Name,
			&b.BudgetItem.ID, &b.BudgetItem.Code, &b.BudgetItem.Name, &b.BudgetItem.Level, &b.BudgetItem.Accumulate,
			&b.InitialQuantity, &b.InitialCost, &b.InitialTotal,
			&b.SpentQuantity, &b.SpentTotal,
			&b.RemainingQuantity, &b.RemainingCost, &b.RemainingTotal,
			&b.UpdatedBudget, &b.CompanyId,
		); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}

	return budgets, nil
}

func (s *service) CreateBudget(budget *types.CreateBudget) (types.Budget, error) {
	tx, _ := s.db.Begin()
	defer tx.Commit()

	if err := saveBudget(budget, tx); err != nil {
		return types.Budget{}, err
	}

	var z float64 = 0
	total := *budget.Quantity * *budget.Cost
	b := types.Budget{
		ProjectId:         budget.ProjectId,
		BudgetItemId:      budget.BudgetItemId,
		InitialQuantity:   budget.Quantity,
		InitialCost:       budget.Cost,
		InitialTotal:      total,
		SpentQuantity:     &z,
		SpentTotal:        0,
		RemainingQuantity: budget.Quantity,
		RemainingCost:     budget.Cost,
		RemainingTotal:    total,
		UpdatedBudget:     total,
		CompanyId:         budget.CompanyId,
	}

	return b, nil
}

func saveBudget(b *types.CreateBudget, s *sql.Tx) error {
	if b == nil || b.BudgetItemId == uuid.Nil {
		return nil
	}

	budget := &types.Budget{}
	budget.BudgetItemId = b.BudgetItemId
	budget.CompanyId = b.CompanyId
	budget.ProjectId = b.ProjectId

	total := *b.Quantity * *b.Cost
	budget.InitialTotal = total
	budget.SpentTotal = 0
	budget.RemainingTotal = total
	budget.UpdatedBudget = total
	budget.InitialQuantity = b.Quantity
	budget.InitialCost = b.Cost
	budget.RemainingQuantity = b.Quantity
	budget.RemainingCost = b.Cost
	var z float64 = 0
	budget.SpentQuantity = &z

	query := "select accumulate, parent_id from budget_item where id = $1 and company_id = $2"
	var accumulate bool
	var parentId uuid.UUID

	err := s.QueryRow(query, budget.BudgetItemId, budget.CompanyId).Scan(&accumulate, &parentId)
	if err != nil {
		return err
	}

	if !accumulate {
		query := `
            insert into budget
            (project_id, budget_item_id, initial_quantity, initial_cost,
            initial_total, spent_quantity, spent_total, remaining_quantity,
            remaining_cost, remaining_total, updated_budget, company_id)
            values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := s.Exec(query, budget.ProjectId, budget.BudgetItemId,
			budget.InitialQuantity, budget.InitialCost, budget.InitialTotal,
			budget.SpentQuantity, budget.SpentTotal, budget.InitialQuantity,
			budget.InitialCost, budget.InitialTotal, budget.UpdatedBudget,
			budget.CompanyId)
		if err != nil {
			return err
		}

	} else {
		query := `
            insert into budget (project_id, budget_item_id, initial_total,
            spent_total, remaining_total, updated_budget, company_id)
            values ($1, $2, $3, $4, $5, $6, $7)
            on conflict (project_id, budget_item_id, company_id)
            do update set initial_total = budget.initial_total + $3,
            spent_total = budget.spent_total + $4,
            remaining_total = budget.remaining_total + $5,
            updated_budget = budget.updated_budget + $6`

		_, err := s.Exec(query, budget.ProjectId, budget.BudgetItemId, budget.InitialTotal, budget.SpentTotal, budget.RemainingTotal, budget.UpdatedBudget, budget.CompanyId)
		if err != nil {
			return err
		}
	}

	b.BudgetItemId = parentId

	return saveBudget(b, s)
}

func (s *service) GetBudgetsByProjectId(companyId, projectId uuid.UUID) ([]types.GetBudget, error) {
	query := `
        SELECT
            project_id, project_name,
            budget_item_id, budget_item_code, budget_item_name, budget_item_level, budget_item_accumulate,
            initial_quantity, initial_cost, initial_total,
            spent_quantity, spent_total,
            remaining_quantity, remaining_cost, remaining_total,
            updated_budget, company_id
        FROM vw_budget
        WHERE company_id = $1 and project_id = $2
        ORDER BY budget_item_code
    `
	rows, err := s.db.Query(query, companyId, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	budgets := []types.GetBudget{}

	for rows.Next() {
		b := types.GetBudget{}
		if err := rows.Scan(
			&b.Project.ID, &b.Project.Name,
			&b.BudgetItem.ID, &b.BudgetItem.Code, &b.BudgetItem.Name, &b.BudgetItem.Level, &b.BudgetItem.Accumulate,
			&b.InitialQuantity, &b.InitialCost, &b.InitialTotal,
			&b.SpentQuantity, &b.SpentTotal,
			&b.RemainingQuantity, &b.RemainingCost, &b.RemainingTotal,
			&b.UpdatedBudget, &b.CompanyId,
		); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}

	return budgets, nil
}

func (s *service) GetOneBudget(companyId, projectId, budgetItemId uuid.UUID) (*types.GetBudget, error) {
	b := &types.GetBudget{}
	query := `
        SELECT
            project_id, project_name,
            budget_item_id, budget_item_code, budget_item_name, budget_item_level, budget_item_accumulate,
            initial_quantity, initial_cost, initial_total,
            spent_quantity, spent_total,
            remaining_quantity, remaining_cost, remaining_total,
            updated_budget, company_id
        FROM vw_budget
        WHERE company_id = $1 and project_id = $2 and budget_item_id = $3
    `

	err := s.db.QueryRow(query, companyId, projectId, budgetItemId).Scan(
		&b.Project.ID, &b.Project.Name,
		&b.BudgetItem.ID, &b.BudgetItem.Code, &b.BudgetItem.Name, &b.BudgetItem.Level, &b.BudgetItem.Accumulate,
		&b.InitialQuantity, &b.InitialCost, &b.InitialTotal,
		&b.SpentQuantity, &b.SpentTotal,
		&b.RemainingQuantity, &b.RemainingCost, &b.RemainingTotal,
		&b.UpdatedBudget, &b.CompanyId,
	)
	if err != nil {
		return b, err
	}

	return b, nil
}
