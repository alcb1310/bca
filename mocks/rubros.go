package mocks

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) GetAllRubros(companyId uuid.UUID) ([]types.Rubro, error) {
	args := s.Called(companyId)
	return args.Get(0).([]types.Rubro), args.Error(1)
}

func (s *ServiceMock) CreateRubro(rubro types.Rubro) (uuid.UUID, error) {
	args := s.Called(rubro)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (s *ServiceMock) GetOneRubro(id, companyId uuid.UUID) (types.Rubro, error) {
	args := s.Called(id, companyId)
	return args.Get(0).(types.Rubro), args.Error(1)
}

func (s *ServiceMock) UpdateRubro(rubro types.Rubro) error {
	args := s.Called(rubro)
	return args.Error(0)
}
