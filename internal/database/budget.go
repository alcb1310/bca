package database

import (
	"database/sql"
	"log/slog"

	"github.com/google/uuid"

	"github.com/alcb1310/bca/internal/types"
)

func (s *service) GetBudgets(companyId, project_id uuid.UUID, search string) ([]types.GetBudget, error) {
	var rows *sql.Rows
	var err error

	term := "%" + search + "%"

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
            AND (budget_item_code like $2 OR budget_item_name like $2)
    `

	if project_id != uuid.Nil {
		query += " AND project_id = $3"
		query += `
            ORDER BY project_name, budget_item_code
        `
		rows, err = s.db.Query(query, companyId, term, project_id)
	} else {
		query += `
            ORDER BY project_name, budget_item_code
        `
		rows, err = s.db.Query(query, companyId, term)
	}

	if err != nil {
		slog.Error("Query Error: ", "err", err)
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
			slog.Error("Scan Error: ", "err", err)
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

	z := sql.NullFloat64{Float64: 0, Valid: true}
	q := sql.NullFloat64{Float64: budget.Quantity, Valid: true}
	c := sql.NullFloat64{Float64: budget.Cost, Valid: true}
	total := budget.Quantity * budget.Cost
	b := types.Budget{
		ProjectId:         budget.ProjectId,
		BudgetItemId:      budget.BudgetItemId,
		InitialQuantity:   q,
		InitialCost:       c,
		InitialTotal:      total,
		SpentQuantity:     z,
		SpentTotal:        0,
		RemainingQuantity: q,
		RemainingCost:     c,
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

	q := sql.NullFloat64{Float64: b.Quantity, Valid: true}
	c := sql.NullFloat64{Float64: b.Cost, Valid: true}

	total := b.Quantity * b.Cost
	budget.InitialTotal = total
	budget.SpentTotal = 0
	budget.RemainingTotal = total
	budget.UpdatedBudget = total
	budget.InitialQuantity = q
	budget.InitialCost = c
	budget.RemainingQuantity = q
	budget.RemainingCost = c
	z := sql.NullFloat64{Float64: 0, Valid: true}
	budget.SpentQuantity = z

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

func (s *service) GetBudgetsByProjectId(companyId, projectId uuid.UUID, level *uint8) ([]types.GetBudget, error) {
	var rows *sql.Rows
	var err error
	query := `
        SELECT
            project_id, project_name, project_net_area, project_gross_area,
            budget_item_id, budget_item_code, budget_item_name, budget_item_level, budget_item_accumulate,
            initial_quantity, initial_cost, initial_total,
            spent_quantity, spent_total,
            remaining_quantity, remaining_cost, remaining_total,
            updated_budget, company_id
        FROM vw_budget
        WHERE company_id = $1 and project_id = $2
		`
	if level == nil {
		query += "ORDER BY budget_item_code"
		rows, err = s.db.Query(query, companyId, projectId)
	} else {
		query += "AND budget_item_level <= $3 ORDER BY budget_item_code"
		rows, err = s.db.Query(query, companyId, projectId, *level)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	budgets := []types.GetBudget{}

	for rows.Next() {
		b := types.GetBudget{}
		if err := rows.Scan(
			&b.Project.ID, &b.Project.Name, &b.Project.NetArea, &b.Project.GrossArea,
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

func (s *service) UpdateBudget(b types.CreateBudget, budget types.Budget) error {
	total := b.Quantity * b.Cost
	updated := total - budget.UpdatedBudget
	diff := total - budget.RemainingTotal

	q := sql.NullFloat64{Float64: b.Quantity, Valid: true}
	c := sql.NullFloat64{Float64: b.Cost, Valid: true}
	toUpdate := types.Budget{
		ProjectId:         budget.ProjectId,
		BudgetItemId:      budget.BudgetItemId,
		InitialQuantity:   budget.InitialQuantity,
		InitialCost:       budget.InitialCost,
		InitialTotal:      budget.InitialTotal,
		SpentQuantity:     budget.SpentQuantity,
		SpentTotal:        budget.SpentTotal,
		RemainingQuantity: q,
		RemainingCost:     c,
		RemainingTotal:    total,
		UpdatedBudget:     updated,
		CompanyId:         budget.CompanyId,
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.executeUpdateBudget(&toUpdate, tx, diff)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (s *service) executeUpdateBudget(budget *types.Budget, tx *sql.Tx, diff float64) error {
	bi, err := s.getBudgetItem(budget.BudgetItemId, budget.CompanyId)
	if err != nil {
		return err
	}
	if bi == nil {
		return nil
	}

	if bi.Accumulate.Bool {
		query := `
			 UPDATE budget
			 SET remaining_total = budget.remaining_total + $1, updated_budget = budget.updated_budget + $1
			 WHERE project_id = $2 and budget_item_id = $3 and company_id = $4
		  `
		_, err = tx.Exec(
			query, diff,
			budget.ProjectId, budget.BudgetItemId, budget.CompanyId,
		)
	} else {
		query := `
			 UPDATE budget
			 SET remaining_quantity = $1, remaining_cost = $2, remaining_total = $3, updated_budget = budget.updated_budget + $4
			 WHERE project_id = $5 and budget_item_id = $6 and company_id = $7
		  `
		_, err = tx.Exec(
			query, budget.RemainingQuantity, budget.RemainingCost, budget.RemainingTotal, budget.UpdatedBudget,
			budget.ProjectId, budget.BudgetItemId, budget.CompanyId,
		)
	}

	if err != nil {
		return err
	}

	if bi.ParentId == nil {
		return nil
	}
	budget.BudgetItemId = *bi.ParentId

	return s.executeUpdateBudget(budget, tx, diff)
}
