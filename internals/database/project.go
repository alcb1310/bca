package database

import (
	"log/slog"

	"github.com/google/uuid"

	"github.com/alcb1310/bca/internals/types"
)

func (s *service) GetAllProjects(companyID uuid.UUID) []types.Project {
	projects := []types.Project{}

	query := "select id, name, is_active, company_id from project where company_id = $1"
	rows, err := s.DB.Query(query, companyID)
	if err != nil {
		slog.Error("GetAllProjects: Error querying the projects", "err", err)
		return projects
	}
	defer rows.Close()

	for rows.Next() {
		project := types.Project{}
		if err := rows.Scan(&project.ID, &project.Name, &project.IsActive, &project.CompanyID); err != nil {
			slog.Error("GetAllProjects: Error scanning the rows", "err", err)
			return []types.Project{}
		}

		projects = append(projects, project)
	}

	return projects
}

func (s *service) CreateProject(project types.Project) (types.Project, error) {
    query := "insert into project (name, is_active, company_id) values ($1, $2, $3) returning id"
    if err := s.DB.QueryRow(query, project.Name, project.IsActive, project.CompanyID).Scan(&project.ID); err != nil {
        return project, err
    }

    return project, nil
}
