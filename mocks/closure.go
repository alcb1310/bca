package mocks

import (
	"time"

	"github.com/google/uuid"
)

func (s *ServiceMock) CreateClosure(companyId, projectId uuid.UUID, date time.Time) error {
	args := s.Called(companyId, projectId, date)
	return args.Error(0)
}
