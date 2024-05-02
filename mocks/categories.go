package mocks

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) GetAllCategories(companyId uuid.UUID) ([]types.Category, error) {
	args := s.Called(companyId)
	return args.Get(0).([]types.Category), args.Error(1)
}

func (s *ServiceMock) CreateCategory(category types.Category) error {
	args := s.Called(category)
	return args.Error(0)
}

func (s *ServiceMock) GetCategory(id, companyId uuid.UUID) (types.Category, error) {
	args := s.Called(id, companyId)
	return args.Get(0).(types.Category), args.Error(1)
}

func (s *ServiceMock) UpdateCategory(category types.Category) error {
	args := s.Called(category)
	return args.Error(0)
}
