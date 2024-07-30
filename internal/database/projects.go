package database

import (
	"log/slog"

	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *service) GetAllProjects(companyId uuid.UUID) ([]types.Project, error) {
	projects := []types.Project{}

	sql := "select id, name, is_active, company_id, gross_area, net_area from project where company_id = $1 order by is_active desc, name"
	rows, err := s.db.Query(sql, companyId)
	if err != nil {
		return projects, err
	}
	defer rows.Close()

	for rows.Next() {
		p := types.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.IsActive, &p.CompanyId, &p.GrossArea, &p.NetArea); err != nil {
			return projects, err
		}
		projects = append(projects, p)
	}

	slog.Info("GetAllProjects", "projects", projects)
	return projects, nil
}

func (s *service) CreateProject(p types.Project) (types.Project, error) {
	sql := "insert into project (name, is_active, company_id, gross_area, net_area) values ($1, $2, $3, $4, $5) returning id"
	err := s.db.QueryRow(sql, p.Name, p.IsActive, p.CompanyId, p.GrossArea, p.NetArea).Scan(&p.ID)
	if err != nil {
		return types.Project{}, err
	}

	return p, nil
}

func (s *service) GetProject(id, companyId uuid.UUID) (types.Project, error) {
	p := types.Project{}
	sql := "select id, name, is_active, company_id, gross_area, net_area from project where id = $1 and company_id = $2"

	err := s.db.QueryRow(sql, id, companyId).Scan(&p.ID, &p.Name, &p.IsActive, &p.CompanyId, &p.GrossArea, &p.NetArea)
	if err != nil {
		return types.Project{}, err
	}

	return p, nil
}

func (s *service) UpdateProject(p types.Project, id, companyId uuid.UUID) error {
	sql := "update project set name = $1, is_active = $2, gross_area = $5, net_area = $6 where id = $3 and company_id = $4"
	_, err := s.db.Exec(sql, p.Name, p.IsActive, id, companyId, p.GrossArea, p.NetArea)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetActiveProjects(companyId uuid.UUID, active bool) []types.Project {
	projects := []types.Project{}

	sql := "select id, name, is_active, company_id from project where company_id = $1 and is_active = $2 order by name"
	rows, err := s.db.Query(sql, companyId, active)
	if err != nil {
		return projects
	}
	defer rows.Close()

	for rows.Next() {
		p := types.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.IsActive, &p.CompanyId); err != nil {
			return projects
		}
		projects = append(projects, p)
	}

	return projects
}
