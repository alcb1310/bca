package database

import (
	"bca-go-final/internal/types"

	"github.com/google/uuid"
)

func (s *service) GetAllProjects(companyId uuid.UUID) ([]types.Project, error) {
	projects := []types.Project{}

	sql := "select id, name, is_active, company_id from project where company_id = $1"
	rows, err := s.db.Query(sql, companyId)
	if err != nil {
		return projects, err
	}
	defer rows.Close()

	for rows.Next() {
		p := types.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.IsActive, &p.CompanyId); err != nil {
			return projects, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}
