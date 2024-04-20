package database

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s ServiceMock) GetAllProjects(companyId uuid.UUID) ([]types.Project, error) {
	return []types.Project{}, nil
}

func (s ServiceMock) GetProject(id, companyId uuid.UUID) (types.Project, error) {
	return types.Project{}, nil
}

func (s ServiceMock) UpdateProject(project types.Project, id, companyId uuid.UUID) error {
	return nil
}

func (s ServiceMock) GetActiveProjects(companyId uuid.UUID, active bool) []types.Project {
	return []types.Project{}
}

func (s ServiceMock) CreateProject(p types.Project) (types.Project, error) {
	return types.Project{}, nil
}
