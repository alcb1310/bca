package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetBudgetItems(companyId uuid.UUID) ([]types.BudgetItem, error) {
	sql := "select id, code, name, level, accumulate, parent_id, company_id from vw_budget_item where company_id = $1"

	rows, err := s.db.Query(sql, companyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bis := []types.BudgetItem{}

	for rows.Next() {
		bi := types.BudgetItem{}
		if err := rows.Scan(&bi.ID, &bi.Code, &bi.Name, &bi.Level, &bi.Accumulate, &bi.ParentId, &bi.CompanyId); err != nil {
			return nil, err
		}
		bis = append(bis, bi)
	}

	return bis, nil
}
