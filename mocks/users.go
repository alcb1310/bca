package mocks

import (
	"github.com/google/uuid"

	"bca-go-final/internal/types"
)

func (s *ServiceMock) GetAllUsers(companyId uuid.UUID) ([]types.User, error) {
	args := s.Called(companyId)
	return args.Get(0).([]types.User), args.Error(1)
}

func (s *ServiceMock) CreateUser(u *types.UserCreate) (types.User, error) {
	args := s.Called(u)
	return args.Get(0).(types.User), args.Error(1)
}

func (s *ServiceMock) GetUser(id, companyId uuid.UUID) (types.User, error) {
	args := s.Called(id, companyId)
	return args.Get(0).(types.User), args.Error(1)
}

func (s *ServiceMock) UpdateUser(u types.User, id, companyId uuid.UUID) (types.User, error) {
	args := s.Called(u, id, companyId)
	return args.Get(0).(types.User), args.Error(1)
}

func (s *ServiceMock) UpdatePassword(pass string, id, companyId uuid.UUID) (types.User, error) {
	args := s.Called(pass, id, companyId)
	return args.Get(0).(types.User), args.Error(1)
}

func (s *ServiceMock) DeleteUser(id, companyId uuid.UUID) error {
	args := s.Called(id, companyId)
	return args.Error(0)
}
