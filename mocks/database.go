package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"bca-go-final/internal/types"
)

type ServiceMock struct {
	mock.Mock
}

func NewServiceMock() *ServiceMock {
	return &ServiceMock{}
}

func (s *ServiceMock) Health() map[string]string {
	args := s.Called()
	return args.Get(0).(map[string]string)
}

func (s *ServiceMock) Levels(companyId uuid.UUID) []types.Select {
	args := s.Called(companyId)
	return args.Get(0).([]types.Select)
}

func (s *ServiceMock) CreateCompany(company *types.CompanyCreate) error {
	args := s.Called(company)
	return args.Error(0)
}

func (s *ServiceMock) Login(l *types.Login) (string, error) {
	args := s.Called(l)
	return args.Get(0).(string), args.Error(1)
}

func (s *ServiceMock) RegenerateToken(token string, user uuid.UUID) error {
	args := s.Called(token, user)
	return args.Error(0)
}

func (s *ServiceMock) IsLoggedIn(token string, user uuid.UUID) bool {
	args := s.Called(token, user)
	return args.Get(0).(bool)
}

func (s *ServiceMock) LoadDummyData(companyId uuid.UUID) error {
	args := s.Called(companyId)
	return args.Error(0)
}
