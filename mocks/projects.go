package mocks

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) GetAllProjects(companyId uuid.UUID) ([]types.Project, error) {
	args := s.Called(companyId)
	return args.Get(0).([]types.Project), args.Error(1)
}

func (s *ServiceMock) CreateProject(p types.Project) (types.Project, error) {
	args := s.Called(p)
	return args.Get(0).(types.Project), args.Error(1)
}

func (s *ServiceMock) GetProject(id, companyId uuid.UUID) (types.Project, error) {
	args := s.Called(id, companyId)
	return args.Get(0).(types.Project), args.Error(1)
}

func (s *ServiceMock) UpdateProject(p types.Project, id, companyId uuid.UUID) error {
	args := s.Called(p, id, companyId)
	return args.Error(0)
}

func (s *ServiceMock) GetActiveProjects(companyId uuid.UUID, active bool) []types.Project {
	args := s.Called(companyId, active)
	return args.Get(0).([]types.Project)
}
