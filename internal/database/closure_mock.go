package database

import (
	"time"

	"github.com/google/uuid"
)

func (s ServiceMock) CreateClosure(companyId, projectId uuid.UUID, date time.Time) error {
	return nil
}
